package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	travelmemory "platform-go/internal/domain/travel_memory"
	travelmemoryday "platform-go/internal/domain/travel_memory_day"
	travelmemoryphoto "platform-go/internal/domain/travel_memory_photo"
	"platform-go/internal/infrastructure/repository"
)

const UpdateTravelMemoryUseCaseSuccessMessage = "Update travel memory successfully"

var (
	ErrUpdateTravelMemoryInvalidID = errors.New(
		"travel memory ID must be greater than zero",
	)
	ErrUpdateTravelMemoryInvalidDayID = errors.New(
		"day ID does not belong to this travel memory",
	)
	ErrUpdateTravelMemoryDuplicateDayID = errors.New(
		"duplicate day ID in request",
	)
	ErrUpdateTravelMemoryInvalidPhotoID = errors.New(
		"photo ID does not belong to this day",
	)
	ErrUpdateTravelMemoryDuplicatePhotoID = errors.New(
		"duplicate photo ID in request",
	)
)

type UpdateTravelMemoryUseCase struct {
	db                    repository.DBInterface
	travelMemoryRepo      travelmemory.TravelMemoryRepository
	travelMemoryDayRepo   travelmemoryday.TravelMemoryDayRepository
	travelMemoryPhotoRepo travelmemoryphoto.TravelMemoryPhotoRepository
}

func NewUpdateTravelMemoryUseCase(
	db repository.DBInterface,
	travelMemoryRepo travelmemory.TravelMemoryRepository,
	travelMemoryDayRepo travelmemoryday.TravelMemoryDayRepository,
	travelMemoryPhotoRepo travelmemoryphoto.TravelMemoryPhotoRepository,
) *UpdateTravelMemoryUseCase {
	return &UpdateTravelMemoryUseCase{
		db:                    db,
		travelMemoryRepo:      travelMemoryRepo,
		travelMemoryDayRepo:   travelMemoryDayRepo,
		travelMemoryPhotoRepo: travelMemoryPhotoRepo,
	}
}

type pendingNewPhoto struct {
	newDayIndex int
	memoryDayID int64
	imagePath   string
	sortOrder   int32
}

type syncPlan struct {
	dayIDsToDelete        []int64
	daysToUpdate          []travelmemoryday.TravelMemoryDay
	daysToCreate          []travelmemoryday.TravelMemoryDay
	photoIDsToDelete      []int64
	photosToUpdate        []travelmemoryphoto.TravelMemoryPhoto
	pendingPhotosToCreate []pendingNewPhoto
}

func (uc *UpdateTravelMemoryUseCase) Execute(
	ctx context.Context,
	travelMemoryID int64,
	request api.UpdateTravelMemoryRequest,
) (*api.MessageOnlyResponse, error) {
	if travelMemoryID <= 0 {
		return nil, ErrUpdateTravelMemoryInvalidID
	}

	err := uc.db.Transaction(ctx, func(tx *gorm.DB) error {
		txCtx := repository.WithTx(ctx, tx)

		existingMemory, err := uc.travelMemoryRepo.FindTravelMemoryByID(txCtx, travelMemoryID)
		if err != nil {
			return err
		}
		if existingMemory == nil {
			return ErrTravelMemoryNotFound
		}

		if err := existingMemory.Update(
			request.LocationId,
			request.Title,
			request.StartDate.Time,
			request.EndDate.Time,
			request.Description,
			request.CoverImagePath,
		); err != nil {
			return err
		}

		if request.Days == nil {
			return uc.travelMemoryRepo.UpdateTravelMemory(txCtx, existingMemory)
		}

		plan, err := uc.buildAndValidateSyncPlan(txCtx, existingMemory, *request.Days)
		if err != nil {
			return err
		}

		if err := uc.travelMemoryRepo.UpdateTravelMemory(txCtx, existingMemory); err != nil {
			return err
		}

		return uc.applySyncPlan(txCtx, plan)
	})
	if err != nil {
		return nil, err
	}

	return &api.MessageOnlyResponse{
		Message: UpdateTravelMemoryUseCaseSuccessMessage,
	}, nil
}

func (uc *UpdateTravelMemoryUseCase) buildAndValidateSyncPlan(
	txCtx context.Context,
	memory *travelmemory.TravelMemory,
	rawDays []api.UpdateTravelMemoryDayInput,
) (*syncPlan, error) {
	existingDays, err := uc.travelMemoryDayRepo.FindDaysByMemoryID(txCtx, memory.ID)
	if err != nil {
		return nil, err
	}

	existingDayMap := make(map[int64]travelmemoryday.TravelMemoryDay, len(existingDays))
	existingDayIDs := make([]int64, len(existingDays))
	for i, day := range existingDays {
		existingDayMap[day.ID] = day
		existingDayIDs[i] = day.ID
	}

	if len(rawDays) == 0 {
		return &syncPlan{
			dayIDsToDelete: existingDayIDs,
		}, nil
	}

	var existingPhotos []travelmemoryphoto.TravelMemoryPhoto
	if len(existingDayIDs) > 0 {
		existingPhotos, err = uc.travelMemoryPhotoRepo.FindPhotosByDayIDs(txCtx, existingDayIDs)
		if err != nil {
			return nil, err
		}
	}

	existingPhotoMapByDayID := make(map[int64]map[int64]travelmemoryphoto.TravelMemoryPhoto, len(existingDays))
	existingPhotoListByDayID := make(map[int64][]travelmemoryphoto.TravelMemoryPhoto, len(existingDays))
	for _, photo := range existingPhotos {
		if existingPhotoMapByDayID[photo.MemoryDayID] == nil {
			existingPhotoMapByDayID[photo.MemoryDayID] = make(map[int64]travelmemoryphoto.TravelMemoryPhoto)
		}
		existingPhotoMapByDayID[photo.MemoryDayID][photo.ID] = photo
		existingPhotoListByDayID[photo.MemoryDayID] = append(existingPhotoListByDayID[photo.MemoryDayID], photo)
	}

	plan := &syncPlan{
		daysToUpdate: make([]travelmemoryday.TravelMemoryDay, 0, len(rawDays)),
		daysToCreate: make([]travelmemoryday.TravelMemoryDay, 0, len(rawDays)),
	}

	requestedDayIDs := make(map[int64]struct{}, len(rawDays))
	requestedPhotoIDs := make(map[int64]struct{})
	seenDates := make(map[string]struct{}, len(rawDays))

	for _, dayReq := range rawDays {
		var isExistingDay bool
		var existingDayID int64

		if dayReq.Id != nil {
			dayID := *dayReq.Id
			if dayID <= 0 {
				return nil, fmt.Errorf("%w: %d", ErrUpdateTravelMemoryInvalidDayID, dayID)
			}
			if _, exists := requestedDayIDs[dayID]; exists {
				return nil, fmt.Errorf("%w: %d", ErrUpdateTravelMemoryDuplicateDayID, dayID)
			}
			if _, exists := existingDayMap[dayID]; !exists {
				return nil, fmt.Errorf("%w: %d", ErrUpdateTravelMemoryInvalidDayID, dayID)
			}
			requestedDayIDs[dayID] = struct{}{}
			isExistingDay = true
			existingDayID = dayID
		}

		dayEntity, err := travelmemoryday.NewTravelMemoryDay(
			memory.ID,
			dayReq.Date.Time,
			memory.StartDate,
			memory.EndDate,
			dayReq.Description,
		)
		if err != nil {
			return nil, err
		}

		dateKey := dayEntity.Date.Format(time.DateOnly)
		if _, exists := seenDates[dateKey]; exists {
			return nil, travelmemoryday.ErrDuplicateDayDate
		}
		seenDates[dateKey] = struct{}{}

		newDayIndex := -1
		if isExistingDay {
			dayEntity.ID = existingDayID
			plan.daysToUpdate = append(plan.daysToUpdate, dayEntity)
		} else {
			newDayIndex = len(plan.daysToCreate)
			plan.daysToCreate = append(plan.daysToCreate, dayEntity)
		}

		if dayReq.Photos == nil {
			continue
		}

		if err := validateAndClassifyDayPhotos(
			*dayReq.Photos,
			isExistingDay,
			existingDayID,
			newDayIndex,
			existingPhotoMapByDayID[existingDayID],
			existingPhotoListByDayID[existingDayID],
			requestedPhotoIDs,
			plan,
		); err != nil {
			return nil, err
		}
	}

	for _, existingDay := range existingDays {
		if _, exists := requestedDayIDs[existingDay.ID]; !exists {
			plan.dayIDsToDelete = append(plan.dayIDsToDelete, existingDay.ID)
		}
	}

	return plan, nil
}

func validateAndClassifyDayPhotos(
	rawPhotos []api.UpdateTravelMemoryPhotoInput,
	isExistingDay bool,
	existingDayID int64,
	newDayIndex int,
	dayExistingPhotoMap map[int64]travelmemoryphoto.TravelMemoryPhoto,
	dayExistingPhotoList []travelmemoryphoto.TravelMemoryPhoto,
	requestedPhotoIDs map[int64]struct{},
	plan *syncPlan,
) error {
	dayRequestedPhotoIDs := make(map[int64]struct{}, len(rawPhotos))

	for photoIdx, photoReq := range rawPhotos {
		var isExistingPhoto bool
		var existingPhotoID int64

		if photoReq.Id != nil {
			photoID := *photoReq.Id
			if photoID <= 0 {
				return fmt.Errorf("%w: %d", ErrUpdateTravelMemoryInvalidPhotoID, photoID)
			}
			if _, exists := requestedPhotoIDs[photoID]; exists {
				return fmt.Errorf("%w: %d", ErrUpdateTravelMemoryDuplicatePhotoID, photoID)
			}
			if !isExistingDay {
				return fmt.Errorf("%w: %d", ErrUpdateTravelMemoryInvalidPhotoID, photoID)
			}
			if _, exists := dayExistingPhotoMap[photoID]; !exists {
				return fmt.Errorf("%w: %d", ErrUpdateTravelMemoryInvalidPhotoID, photoID)
			}
			requestedPhotoIDs[photoID] = struct{}{}
			dayRequestedPhotoIDs[photoID] = struct{}{}
			isExistingPhoto = true
			existingPhotoID = photoID
		}

		sortOrder := int32(photoIdx + 1)

		if isExistingPhoto {
			photoEntity, err := travelmemoryphoto.NewTravelMemoryPhoto(
				existingDayID,
				photoReq.ImagePath,
				sortOrder,
			)
			if err != nil {
				return err
			}
			photoEntity.ID = existingPhotoID
			plan.photosToUpdate = append(plan.photosToUpdate, photoEntity)
		} else {
			plan.pendingPhotosToCreate = append(plan.pendingPhotosToCreate, pendingNewPhoto{
				newDayIndex: newDayIndex,
				memoryDayID: existingDayID,
				imagePath:   photoReq.ImagePath,
				sortOrder:   sortOrder,
			})
		}
	}

	if isExistingDay {
		for _, existingPhoto := range dayExistingPhotoList {
			if _, exists := dayRequestedPhotoIDs[existingPhoto.ID]; !exists {
				plan.photoIDsToDelete = append(plan.photoIDsToDelete, existingPhoto.ID)
			}
		}
	}

	return nil
}

func (uc *UpdateTravelMemoryUseCase) applySyncPlan(
	txCtx context.Context,
	plan *syncPlan,
) error {
	if len(plan.dayIDsToDelete) > 0 {
		if err := uc.travelMemoryPhotoRepo.DeletePhotosByDayIDs(txCtx, plan.dayIDsToDelete); err != nil {
			return err
		}
		if err := uc.travelMemoryDayRepo.DeleteDaysByIDs(txCtx, plan.dayIDsToDelete); err != nil {
			return err
		}
	}

	if len(plan.daysToUpdate) > 0 {
		if err := uc.travelMemoryDayRepo.UpdateTravelMemoryDays(txCtx, plan.daysToUpdate); err != nil {
			return err
		}
	}

	if len(plan.daysToCreate) > 0 {
		if err := uc.travelMemoryDayRepo.CreateTravelMemoryDays(txCtx, plan.daysToCreate); err != nil {
			return err
		}
	}

	if len(plan.photoIDsToDelete) > 0 {
		if err := uc.travelMemoryPhotoRepo.DeletePhotosByIDs(txCtx, plan.photoIDsToDelete); err != nil {
			return err
		}
	}

	if len(plan.photosToUpdate) > 0 {
		if err := uc.travelMemoryPhotoRepo.UpdateTravelMemoryPhotos(txCtx, plan.photosToUpdate); err != nil {
			return err
		}
	}

	if len(plan.pendingPhotosToCreate) > 0 {
		photosToCreate := make([]travelmemoryphoto.TravelMemoryPhoto, len(plan.pendingPhotosToCreate))
		for i, pending := range plan.pendingPhotosToCreate {
			dayID := pending.memoryDayID
			if pending.newDayIndex >= 0 {
				dayID = plan.daysToCreate[pending.newDayIndex].ID
			}

			photoEntity, err := travelmemoryphoto.NewTravelMemoryPhoto(
				dayID,
				pending.imagePath,
				pending.sortOrder,
			)
			if err != nil {
				return err
			}
			photosToCreate[i] = photoEntity
		}

		if err := uc.travelMemoryPhotoRepo.CreateTravelMemoryPhotos(txCtx, photosToCreate); err != nil {
			return err
		}
	}

	return nil
}

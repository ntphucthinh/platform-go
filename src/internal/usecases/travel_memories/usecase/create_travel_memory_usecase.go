package usecase

import (
	"context"
	"time"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	travelmemory "platform-go/internal/domain/travel_memory"
	travelmemoryday "platform-go/internal/domain/travel_memory_day"
	travelmemoryphoto "platform-go/internal/domain/travel_memory_photo"
	"platform-go/internal/infrastructure/repository"
)

const CreateTravelMemoryUseCaseSuccessMessage = "Create new travel memory successfully"

type CreateTravelMemoryUseCase struct {
	db                    repository.DBInterface
	travelMemoryRepo      travelmemory.TravelMemoryRepository
	travelMemoryDayRepo   travelmemoryday.TravelMemoryDayRepository
	travelMemoryPhotoRepo travelmemoryphoto.TravelMemoryPhotoRepository
}

func NewCreateTravelMemoryUseCase(
	db repository.DBInterface,
	travelMemoryRepo travelmemory.TravelMemoryRepository,
	travelMemoryDayRepo travelmemoryday.TravelMemoryDayRepository,
	travelMemoryPhotoRepo travelmemoryphoto.TravelMemoryPhotoRepository,
) *CreateTravelMemoryUseCase {
	return &CreateTravelMemoryUseCase{
		db:                    db,
		travelMemoryRepo:      travelMemoryRepo,
		travelMemoryDayRepo:   travelMemoryDayRepo,
		travelMemoryPhotoRepo: travelMemoryPhotoRepo,
	}
}

func (uc *CreateTravelMemoryUseCase) Execute(
	ctx context.Context,
	request api.CreateTravelMemoryRequest,
) (*api.MessageOnlyResponse, error) {
	err := uc.db.Transaction(ctx, func(tx *gorm.DB) error {
		txCtx := repository.WithTx(ctx, tx)

		memory, err := travelmemory.NewTravelMemory(
			request.LocationId,
			request.Title,
			request.StartDate.Time,
			request.EndDate.Time,
			request.Description,
			request.CoverImagePath,
		)
		if err != nil {
			return err
		}

		if err := uc.travelMemoryRepo.CreateTravelMemory(txCtx, &memory); err != nil {
			return err
		}

		if request.Days == nil || len(*request.Days) == 0 {
			return nil
		}

		rawDays := *request.Days

		days, err := uc.buildAndValidateDays(memory.ID, memory.StartDate, memory.EndDate, rawDays)
		if err != nil {
			return err
		}

		if err := uc.travelMemoryDayRepo.CreateTravelMemoryDays(txCtx, days); err != nil {
			return err
		}

		photos, err := uc.buildAndValidatePhotos(days, rawDays)
		if err != nil {
			return err
		}

		if len(photos) > 0 {
			if err := uc.travelMemoryPhotoRepo.CreateTravelMemoryPhotos(txCtx, photos); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &api.MessageOnlyResponse{
		Message: CreateTravelMemoryUseCaseSuccessMessage,
	}, nil
}

func (uc *CreateTravelMemoryUseCase) buildAndValidateDays(
	memoryID int64,
	startDate, endDate time.Time,
	rawDays []api.CreateTravelMemoryDayInput,
) ([]travelmemoryday.TravelMemoryDay, error) {
	days := make([]travelmemoryday.TravelMemoryDay, len(rawDays))
	seenDates := make(map[string]struct{}, len(rawDays))

	for i, d := range rawDays {
		day, err := travelmemoryday.NewTravelMemoryDay(
			memoryID,
			d.Date.Time,
			startDate,
			endDate,
			d.Description,
		)
		if err != nil {
			return nil, err
		}

		dateKey := day.Date.Format(time.DateOnly)
		if _, exists := seenDates[dateKey]; exists {
			return nil, travelmemoryday.ErrDuplicateDayDate
		}
		seenDates[dateKey] = struct{}{}

		days[i] = day
	}

	return days, nil
}

func (uc *CreateTravelMemoryUseCase) buildAndValidatePhotos(
	days []travelmemoryday.TravelMemoryDay,
	rawDays []api.CreateTravelMemoryDayInput,
) ([]travelmemoryphoto.TravelMemoryPhoto, error) {
	var totalPhotos int
	for _, d := range rawDays {
		if d.Photos != nil {
			totalPhotos += len(*d.Photos)
		}
	}
	if totalPhotos == 0 {
		return nil, nil
	}

	allPhotos := make([]travelmemoryphoto.TravelMemoryPhoto, 0, totalPhotos)

	for i, d := range rawDays {
		if d.Photos == nil || len(*d.Photos) == 0 {
			continue
		}

		dayID := days[i].ID

		for photoIdx, p := range *d.Photos {
			sortOrder := int32(photoIdx + 1)
			photo, err := travelmemoryphoto.NewTravelMemoryPhoto(
				dayID,
				p.ImagePath,
				sortOrder,
			)
			if err != nil {
				return nil, err
			}

			allPhotos = append(allPhotos, photo)
		}
	}

	return allPhotos, nil
}

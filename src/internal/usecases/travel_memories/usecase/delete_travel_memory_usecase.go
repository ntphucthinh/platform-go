package usecase

import (
	"context"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	travelmemory "platform-go/internal/domain/travel_memory"
	travelmemoryday "platform-go/internal/domain/travel_memory_day"
	travelmemoryphoto "platform-go/internal/domain/travel_memory_photo"
	"platform-go/internal/infrastructure/repository"
)

const DeleteTravelMemoryUseCaseSuccessMessage = "Delete travel memory successfully"

type DeleteTravelMemoryUseCase struct {
	db                    repository.DBInterface
	travelMemoryRepo      travelmemory.TravelMemoryRepository
	travelMemoryDayRepo   travelmemoryday.TravelMemoryDayRepository
	travelMemoryPhotoRepo travelmemoryphoto.TravelMemoryPhotoRepository
}

func NewDeleteTravelMemoryUseCase(
	db repository.DBInterface,
	travelMemoryRepo travelmemory.TravelMemoryRepository,
	travelMemoryDayRepo travelmemoryday.TravelMemoryDayRepository,
	travelMemoryPhotoRepo travelmemoryphoto.TravelMemoryPhotoRepository,
) *DeleteTravelMemoryUseCase {
	return &DeleteTravelMemoryUseCase{
		db:                    db,
		travelMemoryRepo:      travelMemoryRepo,
		travelMemoryDayRepo:   travelMemoryDayRepo,
		travelMemoryPhotoRepo: travelMemoryPhotoRepo,
	}
}

func (uc *DeleteTravelMemoryUseCase) Execute(ctx context.Context, travelMemoryID int64) (*api.MessageOnlyResponse, error) {
	err := uc.db.Transaction(ctx, func(tx *gorm.DB) error {
		txCtx := repository.WithTx(ctx, tx)

		existing, err := uc.travelMemoryRepo.FindTravelMemoryByID(txCtx, travelMemoryID)
		if err != nil {
			return err
		}
		if existing == nil {
			return ErrTravelMemoryNotFound
		}

		days, err := uc.travelMemoryDayRepo.FindDaysByMemoryID(txCtx, travelMemoryID)
		if err != nil {
			return err
		}

		if len(days) > 0 {
			dayIDs := make([]int64, len(days))
			for i, d := range days {
				dayIDs[i] = d.ID
			}
			if err := uc.travelMemoryPhotoRepo.DeletePhotosByDayIDs(txCtx, dayIDs); err != nil {
				return err
			}
			if err := uc.travelMemoryDayRepo.DeleteDaysByIDs(txCtx, dayIDs); err != nil {
				return err
			}
		}

		if err := uc.travelMemoryRepo.DeleteTravelMemoryByID(txCtx, travelMemoryID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &api.MessageOnlyResponse{
		Message: DeleteTravelMemoryUseCaseSuccessMessage,
	}, nil
}

package usecase

import (
	"context"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/domain/itinerary"
	itineraryactivity "platform-go/internal/domain/itinerary_activity"
	itineraryday "platform-go/internal/domain/itinerary_day"
	itinerarytheme "platform-go/internal/domain/itinerary_theme"
	"platform-go/internal/infrastructure/repository"
)

const DeleteItineraryUseCaseSuccessMessage = "Delete itinerary successfully"

type DeleteItineraryUseCase struct {
	db                    repository.DBInterface
	itineraryRepo         itinerary.ItineraryRepository
	itineraryThemeRepo    itinerarytheme.ItineraryThemeRepository
	itineraryDayRepo      itineraryday.ItineraryDayRepository
	itineraryActivityRepo itineraryactivity.ItineraryActivityRepository
}

func NewDeleteItineraryUseCase(
	db repository.DBInterface,
	itineraryRepo itinerary.ItineraryRepository,
	itineraryThemeRepo itinerarytheme.ItineraryThemeRepository,
	itineraryDayRepo itineraryday.ItineraryDayRepository,
	itineraryActivityRepo itineraryactivity.ItineraryActivityRepository,
) *DeleteItineraryUseCase {
	return &DeleteItineraryUseCase{
		db:                    db,
		itineraryRepo:         itineraryRepo,
		itineraryThemeRepo:    itineraryThemeRepo,
		itineraryDayRepo:      itineraryDayRepo,
		itineraryActivityRepo: itineraryActivityRepo,
	}
}

func (uc *DeleteItineraryUseCase) Execute(ctx context.Context, itineraryID int64) (*api.MessageOnlyResponse, error) {
	err := uc.db.Transaction(ctx, func(tx *gorm.DB) error {
		txCtx := repository.WithTx(ctx, tx)

		existing, err := uc.itineraryRepo.FindItineraryByID(txCtx, itineraryID)
		if err != nil {
			return err
		}
		if existing == nil {
			return ErrItineraryNotFound
		}

		days, err := uc.itineraryDayRepo.FindDaysByItineraryID(txCtx, itineraryID)
		if err != nil {
			return err
		}

		if len(days) > 0 {
			dayIDs := make([]int64, len(days))
			for i, d := range days {
				dayIDs[i] = d.ID
			}
			if err := uc.itineraryActivityRepo.DeleteActivitiesByDayIDs(txCtx, dayIDs); err != nil {
				return err
			}
			if err := uc.itineraryDayRepo.DeleteDaysByIDs(txCtx, dayIDs); err != nil {
				return err
			}
		}

		if err := uc.itineraryThemeRepo.DeleteThemesByItineraryID(txCtx, itineraryID); err != nil {
			return err
		}

		if err := uc.itineraryRepo.DeleteItineraryByID(txCtx, itineraryID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &api.MessageOnlyResponse{
		Message: DeleteItineraryUseCaseSuccessMessage,
	}, nil
}

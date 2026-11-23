package usecase

import (
	"context"
	"errors"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	travellocation "platform-go/internal/domain/travel_location"
	"platform-go/internal/infrastructure/repository"
)

const DeleteTravelLocationUseCaseSuccessMessage = "Delete travel location successfully"

var (
	ErrTravelLocationHasMemories    = errors.New("cannot delete travel location because it has associated travel memories")
	ErrTravelLocationHasItineraries = errors.New("cannot delete travel location because it has associated itineraries")
	ErrTravelLocationInUse          = errors.New("cannot delete travel location because it has associated travel memories or itineraries")
)

type DeleteTravelLocationUsecase struct {
	db                 repository.DBInterface
	travelLocationRepo travellocation.TravelLocationRepository
}

func NewDeleteTravelLocationUsecase(
	db repository.DBInterface,
	travelLocationRepo travellocation.TravelLocationRepository,
) *DeleteTravelLocationUsecase {
	return &DeleteTravelLocationUsecase{
		db:                 db,
		travelLocationRepo: travelLocationRepo,
	}
}

func (uc *DeleteTravelLocationUsecase) Execute(ctx context.Context, travelLocationID int64) (*api.MessageOnlyResponse, error) {
	err := uc.db.Transaction(ctx, func(tx *gorm.DB) error {
		txCtx := repository.WithTx(ctx, tx)

		existing, err := uc.travelLocationRepo.FindTravelLocationByID(txCtx, travelLocationID)
		if err != nil {
			return err
		}
		if existing == nil {
			return ErrTravelLocationNotFound
		}

		hasMemories, err := uc.travelLocationRepo.HasAssociatedTravelMemories(txCtx, travelLocationID)
		if err != nil {
			return err
		}
		if hasMemories {
			return ErrTravelLocationHasMemories
		}

		hasItineraries, err := uc.travelLocationRepo.HasAssociatedItineraries(txCtx, travelLocationID)
		if err != nil {
			return err
		}
		if hasItineraries {
			return ErrTravelLocationHasItineraries
		}

		return uc.travelLocationRepo.DeleteTravelLocationByID(txCtx, travelLocationID)
	})
	if err != nil {
		return nil, err
	}

	return &api.MessageOnlyResponse{
		Message: DeleteTravelLocationUseCaseSuccessMessage,
	}, nil
}

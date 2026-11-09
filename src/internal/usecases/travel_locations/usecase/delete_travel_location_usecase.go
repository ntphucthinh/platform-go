package usecase

import (
	"context"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	travellocation "platform-go/internal/domain/travel_location"
	"platform-go/internal/infrastructure/repository"
)

const DeleteTravelLocationUseCaseSuccessMessage = "Delete travel location successfully"

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

		return uc.travelLocationRepo.DeleteTravelLocationByID(txCtx, travelLocationID)
	})
	if err != nil {
		return nil, err
	}

	return &api.MessageOnlyResponse{
		Message: DeleteTravelLocationUseCaseSuccessMessage,
	}, nil
}


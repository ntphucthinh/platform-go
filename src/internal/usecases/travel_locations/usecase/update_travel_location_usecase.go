package usecase

import (
	"context"
	"errors"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	travellocation "platform-go/internal/domain/travel_location"
	"platform-go/internal/infrastructure/repository"
)

const UpdateTravelLocationUseCaseSuccessMessage = "Update travel location successfully"

var ErrUpdateTravelLocationInvalidID = errors.New("travel location ID must be greater than zero")

type UpdateTravelLocationUsecase struct {
	db                 repository.DBInterface
	travelLocationRepo travellocation.TravelLocationRepository
}

func NewUpdateTravelLocationUsecase(
	db repository.DBInterface,
	travelLocationRepo travellocation.TravelLocationRepository,
) *UpdateTravelLocationUsecase {
	return &UpdateTravelLocationUsecase{
		db:                 db,
		travelLocationRepo: travelLocationRepo,
	}
}

func (uc *UpdateTravelLocationUsecase) Execute(
	ctx context.Context,
	travelLocationID int64,
	res api.UpdateTravelLocationRequest,
) (*api.MessageOnlyResponse, error) {
	if travelLocationID <= 0 {
		return nil, ErrUpdateTravelLocationInvalidID
	}

	err := uc.db.Transaction(ctx, func(tx *gorm.DB) error {
		txCtx := repository.WithTx(ctx, tx)

		existing, err := uc.travelLocationRepo.FindTravelLocationByID(txCtx, travelLocationID)
		if err != nil {
			return err
		}
		if existing == nil {
			return ErrTravelLocationNotFound
		}

		if err := existing.Update(
			res.CountryCode,
			res.CountryName,
			res.LocationName,
			string(res.Status),
			res.Latitude,
			res.Longitude,
			res.ImagePath,
			res.Description,
			res.Notes,
			res.GoogleMapsUrl,
		); err != nil {
			return err
		}

		return uc.travelLocationRepo.UpdateTravelLocation(txCtx, existing)
	})
	if err != nil {
		return nil, err
	}

	return &api.MessageOnlyResponse{
		Message: UpdateTravelLocationUseCaseSuccessMessage,
	}, nil
}

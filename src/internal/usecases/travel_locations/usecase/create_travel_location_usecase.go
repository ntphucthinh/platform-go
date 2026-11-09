package usecase

import (
	"context"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	travellocation "platform-go/internal/domain/travel_location"
	"platform-go/internal/infrastructure/repository"
)

const CreateTravelLocationUsecaseSuccessMessage = "Create new travel location successfully"

type CreateTravelLocationUsecase struct {
	db                 repository.DBInterface
	travelLocationRepo travellocation.TravelLocationRepository
}

func NewCreateTravelLocationUsecase(
	db repository.DBInterface,
	travelLocationRepo travellocation.TravelLocationRepository,
) *CreateTravelLocationUsecase {
	return &CreateTravelLocationUsecase{
		db:                 db,
		travelLocationRepo: travelLocationRepo,
	}
}

func (uc *CreateTravelLocationUsecase) Execute(
	ctx context.Context,
	res api.CreateTravelLocationRequest,
) (*api.MessageOnlyResponse, error) {
	err := uc.db.Transaction(ctx, func(tx *gorm.DB) error {
		txCtx := repository.WithTx(ctx, tx)

		loc, err := travellocation.NewTravelLocation(
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
		)
		if err != nil {
			return err
		}

		return uc.travelLocationRepo.CreateTravelLocation(txCtx, &loc)
	})
	if err != nil {
		return nil, err
	}

	return &api.MessageOnlyResponse{
		Message: CreateTravelLocationUsecaseSuccessMessage,
	}, nil
}

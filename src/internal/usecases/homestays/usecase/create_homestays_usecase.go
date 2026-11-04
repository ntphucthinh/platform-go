package usecase

import (
	"context"
	"errors"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/domain/homestay"
	"platform-go/internal/infrastructure/repository"
)

const CreateHomestaysUseCaseSuccessMessage = "Create new homestay successfully"

var ErrInvalidHomestayInput = errors.New("invalid homestay input")

type CreateHomestayUsecase struct {
	db           repository.DBInterface
	homestayRepo homestay.HomestayRepository
}

func NewCreateHomestayUsecase(
	db repository.DBInterface,
	homestayRepo homestay.HomestayRepository,
) *CreateHomestayUsecase {
	return &CreateHomestayUsecase{
		db:           db,
		homestayRepo: homestayRepo,
	}
}

func (uc *CreateHomestayUsecase) Execute(ctx context.Context, request api.CreateHomestayRequest) (*api.MessageOnlyResponse, error) {
	err := uc.db.Transaction(ctx, func(tx *gorm.DB) error {
		txCtx := repository.WithTx(ctx, tx)

		h := homestay.NewHomestay(
			request.Name,
			request.Address,
			request.Description,
			request.GoogleMapsUrl,
		)
		if err := uc.homestayRepo.CreateHomestay(txCtx, &h); err != nil {
			return err
		}

		if request.Rooms != nil {
			for _, r := range *request.Rooms {
				room := homestay.NewHomestayRoom(
					h.ID,
					r.Name,
					r.Description,
					r.PriceAmount,
					(*string)(r.PriceUnit),
					r.SortOrder,
				)
				if err := uc.homestayRepo.CreateHomestayRoom(txCtx, &room); err != nil {
					return err
				}

				if r.Images != nil && len(*r.Images) > 0 {
					images := make([]homestay.HomestayImage, len(*r.Images))
					for i, img := range *r.Images {
						images[i] = homestay.NewHomestayImage(
							room.ID,
							img.ImagePath,
							img.SortOrder,
							img.IsCover,
						)
					}
					if err := uc.homestayRepo.CreateHomestayImage(txCtx, images); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &api.MessageOnlyResponse{
		Message: CreateHomestaysUseCaseSuccessMessage,
	}, nil
}

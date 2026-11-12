package usecase

import (
	"context"
	"errors"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/domain/homestay"
	homestayimage "platform-go/internal/domain/homestay_image"
	homestayroom "platform-go/internal/domain/homestay_room"
	"platform-go/internal/infrastructure/repository"
)

const CreateHomestaysUseCaseSuccessMessage = "Create new homestay successfully"

var ErrInvalidHomestayInput = errors.New("invalid homestay input")

type CreateHomestayUsecase struct {
	db           repository.DBInterface
	homestayRepo homestay.HomestayRepository
	roomRepo     homestayroom.HomestayRoomRepository
	imageRepo    homestayimage.HomestayImageRepository
}

func NewCreateHomestayUsecase(
	db repository.DBInterface,
	homestayRepo homestay.HomestayRepository,
	roomRepo homestayroom.HomestayRoomRepository,
	imageRepo homestayimage.HomestayImageRepository,
) *CreateHomestayUsecase {
	return &CreateHomestayUsecase{
		db:           db,
		homestayRepo: homestayRepo,
		roomRepo:     roomRepo,
		imageRepo:    imageRepo,
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
				room := homestayroom.NewHomestayRoom(
					h.ID,
					r.Name,
					r.Description,
					r.PriceAmount,
					(*string)(r.PriceUnit),
					r.SortOrder,
				)
				if err := uc.roomRepo.CreateHomestayRoom(txCtx, &room); err != nil {
					return err
				}

				if r.Images != nil && len(*r.Images) > 0 {
					images := make([]homestayimage.HomestayImage, len(*r.Images))
					for i, img := range *r.Images {
						images[i] = homestayimage.NewHomestayImage(
							room.ID,
							img.ImagePath,
							img.SortOrder,
							img.IsCover,
						)
					}
					if err := uc.imageRepo.CreateHomestayImages(txCtx, images); err != nil {
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

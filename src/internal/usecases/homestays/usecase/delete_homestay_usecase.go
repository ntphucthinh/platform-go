package usecase

import (
	"context"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/domain/homestay"
	homestayimage "platform-go/internal/domain/homestay_image"
	homestayroom "platform-go/internal/domain/homestay_room"
	"platform-go/internal/infrastructure/repository"
)

const DeleteHomestaysUseCaseSuccessMessage = "Delete homestay successfully"

type DeleteHomestayUsecase struct {
	db           repository.DBInterface
	homestayRepo homestay.HomestayRepository
	roomRepo     homestayroom.HomestayRoomRepository
	imageRepo    homestayimage.HomestayImageRepository
}

func NewDeleteHomestayUsecase(
	db repository.DBInterface,
	homestayRepo homestay.HomestayRepository,
	roomRepo homestayroom.HomestayRoomRepository,
	imageRepo homestayimage.HomestayImageRepository,
) *DeleteHomestayUsecase {
	return &DeleteHomestayUsecase{
		db:           db,
		homestayRepo: homestayRepo,
		roomRepo:     roomRepo,
		imageRepo:    imageRepo,
	}
}

func (uc *DeleteHomestayUsecase) Execute(ctx context.Context, homestayID int64) (*api.MessageOnlyResponse, error) {
	err := uc.db.Transaction(ctx, func(tx *gorm.DB) error {
		txCtx := repository.WithTx(ctx, tx)

		existing, err := uc.homestayRepo.FindHomestayByID(txCtx, homestayID)
		if err != nil {
			return err
		}
		if existing == nil {
			return ErrHomestayNotFound
		}

		rooms, err := uc.roomRepo.FindRoomsByHomestayID(txCtx, homestayID)
		if err != nil {
			return err
		}

		if len(rooms) > 0 {
			roomIDs := make([]int64, len(rooms))
			for i, r := range rooms {
				roomIDs[i] = r.ID
			}
			if err := uc.imageRepo.DeleteImagesByRoomIDs(txCtx, roomIDs); err != nil {
				return err
			}
			if err := uc.roomRepo.DeleteRoomsByIDs(txCtx, roomIDs); err != nil {
				return err
			}
		}

		if err := uc.homestayRepo.DeleteHomestayByID(txCtx, homestayID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &api.MessageOnlyResponse{
		Message: DeleteHomestaysUseCaseSuccessMessage,
	}, nil
}

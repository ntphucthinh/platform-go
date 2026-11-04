package usecase

import (
	"context"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/domain/homestay"
	"platform-go/internal/infrastructure/repository"
)

const DeleteHomestaysUseCaseSuccessMessage = "Delete homestay successfully"

type DeleteHomestayUsecase struct {
	db           repository.DBInterface
	homestayRepo homestay.HomestayRepository
}

func NewDeleteHomestayUsecase(
	db repository.DBInterface,
	homestayRepo homestay.HomestayRepository,
) *DeleteHomestayUsecase {
	return &DeleteHomestayUsecase{
		db:           db,
		homestayRepo: homestayRepo,
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

		rooms, err := uc.homestayRepo.FindRoomsByHomestayID(txCtx, homestayID)
		if err != nil {
			return err
		}

		if len(rooms) > 0 {
			roomIDs := make([]int64, len(rooms))
			for i, r := range rooms {
				roomIDs[i] = r.ID
			}
			if err := uc.homestayRepo.DeleteImagesByRoomIDs(txCtx, roomIDs); err != nil {
				return err
			}
			if err := uc.homestayRepo.DeleteRoomsByIDs(txCtx, roomIDs); err != nil {
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

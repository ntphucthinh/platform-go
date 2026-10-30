package usecase

import (
	"context"
	"errors"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
)

const DeleteHomestaysUseCaseSuccessMessage = "Delete homestay successfully"

type DeleteHomestayUsecase struct {
	db repository.DBInterface
}

func NewDeleteHomestayUsecase(db repository.DBInterface) *DeleteHomestayUsecase {
	return &DeleteHomestayUsecase{db: db}
}

func (uc *DeleteHomestayUsecase) Execute(ctx context.Context, homestayID int64) (*api.MessageOnlyResponse, error) {
	err := uc.db.Transaction(ctx, func(tx *gorm.DB) error {
		var homestay dbmodel.Homestay
		if err := tx.Where("id = ?", homestayID).First(&homestay).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrHomestayNotFound
			}
			return err
		}

		var roomIDs []int64
		if err := tx.Model(&dbmodel.HomestayRoom{}).Where("homestay_id = ?", homestayID).Pluck("id", &roomIDs).Error; err != nil {
			return err
		}

		if len(roomIDs) > 0 {
			if err := tx.Where("room_id IN ?", roomIDs).Delete(&dbmodel.HomestayImage{}).Error; err != nil {
				return err
			}
			if err := tx.Where("id IN ?", roomIDs).Delete(&dbmodel.HomestayRoom{}).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("id = ?", homestayID).Delete(&dbmodel.Homestay{}).Error; err != nil {
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

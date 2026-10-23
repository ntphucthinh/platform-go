package usecase

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
)

const CreateHomestaysUseCaseSuccessMessage = "Create new homestay successfully"

var ErrInvalidHomestayInput = errors.New("invalid homestay input")

type CreateHomestayUsecase struct {
	db repository.DBInterface
}

func NewCreateHomestayUsecase(db repository.DBInterface) *CreateHomestayUsecase {
	return &CreateHomestayUsecase{db: db}
}

func (uc *CreateHomestayUsecase) Execute(ctx context.Context, request api.CreateHomestayRequest) (*api.MessageOnlyResponse, error) {
	// if strings.TrimSpace(request.Name) == "" {
	// 	return nil, fmt.Errorf("%w: name is required", ErrInvalidHomestayInput)
	// }
	// if strings.TrimSpace(request.Address) == "" {
	// 	return nil, fmt.Errorf("%w: address is required", ErrInvalidHomestayInput)
	// }

	// for _, room := range request.Rooms {
	// 	if strings.TrimSpace(room.Name) == "" {
	// 		return nil, fmt.Errorf("%w: room name is required", ErrInvalidHomestayInput)
	// 	}
	// 	for _, img := range room.Images {
	// 		if strings.TrimSpace(img.ImagePath) == "" {
	// 			return nil, fmt.Errorf("%w: image path is required", ErrInvalidHomestayInput)
	// 		}
	// 	}
	// }

	// var homestayID int64
	err := uc.db.Transaction(ctx, func(tx *gorm.DB) error {
		homestay := dbmodel.Homestay{
			Name:          strings.TrimSpace(request.Name),
			Address:       strings.TrimSpace(request.Address),
			Description:   request.Description,
			GoogleMapsURL: request.GoogleMapsUrl,
		}
		if err := tx.Create(&homestay).Error; err != nil {
			return err
		}

		if request.Rooms != nil {
			for _, r := range *request.Rooms {
				room := dbmodel.HomestayRoom{
					HomestayID:  homestay.ID,
					Name:        strings.TrimSpace(r.Name),
					Description: r.Description,
					PriceAmount: r.PriceAmount,
					PriceUnit:   r.PriceUnit,
					SortOrder:   r.SortOrder,
				}
				if err := tx.Create(&room).Error; err != nil {
					return err
				}

				if r.Images != nil && len(*r.Images) > 0 {
					images := make([]dbmodel.HomestayImage, len(*r.Images))
					for i, img := range *r.Images {
						imagePath := ""
						if img.ImagePath != nil {
							imagePath = strings.TrimSpace(*img.ImagePath)
						}
						images[i] = dbmodel.HomestayImage{
							RoomID:    room.ID,
							ImagePath: imagePath,
							SortOrder: img.SortOrder,
							IsCover:   img.IsCover,
						}
					}
					if err := tx.Create(&images).Error; err != nil {
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

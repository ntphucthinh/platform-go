package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
)

const UpdateHomestaysUseCaseSuccessMessage = "Update homestay successfully"

var (
	ErrUpdateHomestayInvalidRoomID = errors.New(
		"room ID does not belong to this homestay",
	)
	ErrUpdateHomestayInvalidImageID = errors.New(
		"image ID does not belong to this room",
	)
	ErrUpdateHomestayDuplicateRoomID = errors.New(
		"duplicate room ID in request",
	)
	ErrUpdateHomestayDuplicateImageID = errors.New(
		"duplicate image ID in request",
	)
	ErrUpdateHomestayInvalidID = errors.New(
		"homestay ID must be greater than zero",
	)
)

type UpdateHomestayUsecase struct {
	db repository.DBInterface
}

func NewUpdateHomestayUsecase(
	db repository.DBInterface,
) *UpdateHomestayUsecase {
	return &UpdateHomestayUsecase{db: db}
}

func (uc *UpdateHomestayUsecase) Execute(
	ctx context.Context,
	homestayID int64,
	request api.UpdateHomestayRequest,
) (*api.MessageOnlyResponse, error) {
	if homestayID <= 0 {
		return nil, ErrUpdateHomestayInvalidID
	}

	err := uc.db.Transaction(ctx, func(tx *gorm.DB) error {
		// Find the homestay to update.
		var homestay dbmodel.Homestay
		if err := tx.
			Where("id = ?", homestayID).
			First(&homestay).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrHomestayNotFound
			}
			return err
		}

		// Update homestay information.
		homestay.Name = strings.TrimSpace(request.Name)
		homestay.Address = strings.TrimSpace(request.Address)
		homestay.Description = request.Description
		homestay.GoogleMapsURL = request.GoogleMapsUrl

		if err := tx.Save(&homestay).Error; err != nil {
			return err
		}

		// A nil Rooms field means that room data should remain unchanged.
		// An empty Rooms list means that all existing rooms should be removed.
		if request.Rooms == nil {
			return nil
		}

		// Load existing rooms belonging to this homestay.
		var existingRooms []dbmodel.HomestayRoom
		if err := tx.
			Where("homestay_id = ?", homestayID).
			Find(&existingRooms).Error; err != nil {
			return err
		}

		existingRoomMap := make(map[int64]dbmodel.HomestayRoom, len(existingRooms))
		for _, room := range existingRooms {
			existingRoomMap[room.ID] = room
		}

		// Validate requested room IDs before deleting existing records.
		requestedRoomIDs := make(map[int64]struct{})

		for _, room := range *request.Rooms {
			if room.Id == nil || *room.Id == 0 {
				continue
			}

			if *room.Id < 0 {
				return fmt.Errorf(
					"%w: %d",
					ErrUpdateHomestayInvalidRoomID,
					*room.Id,
				)
			}

			roomID := *room.Id

			if _, exists := requestedRoomIDs[roomID]; exists {
				return fmt.Errorf(
					"%w: %d",
					ErrUpdateHomestayDuplicateRoomID,
					roomID,
				)
			}

			if _, exists := existingRoomMap[roomID]; !exists {
				return fmt.Errorf(
					"%w: %d",
					ErrUpdateHomestayInvalidRoomID,
					roomID,
				)
			}

			requestedRoomIDs[roomID] = struct{}{}
		}

		// Delete rooms omitted from the request, together with their images.
		var roomIDsToDelete []int64

		for roomID := range existingRoomMap {
			if _, exists := requestedRoomIDs[roomID]; !exists {
				roomIDsToDelete = append(roomIDsToDelete, roomID)
			}
		}

		if len(roomIDsToDelete) > 0 {
			if err := tx.
				Where("room_id IN ?", roomIDsToDelete).
				Delete(&dbmodel.HomestayImage{}).Error; err != nil {
				return err
			}

			if err := tx.
				Where("id IN ?", roomIDsToDelete).
				Delete(&dbmodel.HomestayRoom{}).Error; err != nil {
				return err
			}
		}

		// Create or update rooms.
		for _, roomRequest := range *request.Rooms {
			var roomID int64

			if roomRequest.Id != nil && *roomRequest.Id > 0 {
				roomID = *roomRequest.Id

				roomUpdates := map[string]interface{}{
					"name":         strings.TrimSpace(roomRequest.Name),
					"description":  roomRequest.Description,
					"price_amount": roomRequest.PriceAmount,
					"price_unit":   (*string)(roomRequest.PriceUnit),
					"sort_order":   roomRequest.SortOrder,
				}

				if err := tx.
					Model(&dbmodel.HomestayRoom{}).
					Where(
						"id = ? AND homestay_id = ?",
						roomID,
						homestayID,
					).
					Updates(roomUpdates).Error; err != nil {
					return err
				}
			} else {
				newRoom := dbmodel.HomestayRoom{
					HomestayID:  homestayID,
					Name:        strings.TrimSpace(roomRequest.Name),
					Description: roomRequest.Description,
					PriceAmount: roomRequest.PriceAmount,
					PriceUnit:   (*string)(roomRequest.PriceUnit),
					SortOrder:   roomRequest.SortOrder,
				}

				if err := tx.Create(&newRoom).Error; err != nil {
					return err
				}

				roomID = newRoom.ID
			}

			// A nil Images field means that the room's images remain unchanged.
			// An empty Images list means that all images of this room are removed.
			if roomRequest.Images == nil {
				continue
			}

			// Load images belonging to the current room only.
			var existingImages []dbmodel.HomestayImage
			if err := tx.
				Where("room_id = ?", roomID).
				Find(&existingImages).Error; err != nil {
				return err
			}

			existingImageMap := make(
				map[int64]dbmodel.HomestayImage,
				len(existingImages),
			)

			for _, image := range existingImages {
				existingImageMap[image.ID] = image
			}

			// Validate image IDs and detect duplicates.
			requestedImageIDs := make(map[int64]struct{})

			for _, imageRequest := range *roomRequest.Images {
				if imageRequest.Id == nil || *imageRequest.Id == 0 {
					continue
				}

				if *imageRequest.Id < 0 {
					return fmt.Errorf(
						"%w: %d",
						ErrUpdateHomestayInvalidImageID,
						*imageRequest.Id,
					)
				}

				imageID := *imageRequest.Id

				if _, exists := requestedImageIDs[imageID]; exists {
					return fmt.Errorf(
						"%w: %d",
						ErrUpdateHomestayDuplicateImageID,
						imageID,
					)
				}

				if _, exists := existingImageMap[imageID]; !exists {
					return fmt.Errorf(
						"%w: %d",
						ErrUpdateHomestayInvalidImageID,
						imageID,
					)
				}

				requestedImageIDs[imageID] = struct{}{}
			}

			// Delete images omitted from the request.
			var imageIDsToDelete []int64

			for imageID := range existingImageMap {
				if _, exists := requestedImageIDs[imageID]; !exists {
					imageIDsToDelete = append(imageIDsToDelete, imageID)
				}
			}

			if len(imageIDsToDelete) > 0 {
				if err := tx.
					Where("id IN ?", imageIDsToDelete).
					Delete(&dbmodel.HomestayImage{}).Error; err != nil {
					return err
				}
			}

			// Create or update images.
			for _, imageRequest := range *roomRequest.Images {
				imagePath := ""
				if imageRequest.ImagePath != nil {
					imagePath = strings.TrimSpace(*imageRequest.ImagePath)
				}

				if imageRequest.Id != nil && *imageRequest.Id > 0 {
					imageID := *imageRequest.Id

					imageUpdates := map[string]interface{}{
						"image_path": imagePath,
						"sort_order": imageRequest.SortOrder,
						"is_cover":   imageRequest.IsCover,
					}

					if err := tx.
						Model(&dbmodel.HomestayImage{}).
						Where(
							"id = ? AND room_id = ?",
							imageID,
							roomID,
						).
						Updates(imageUpdates).Error; err != nil {
						return err
					}
				} else {
					newImage := dbmodel.HomestayImage{
						RoomID:    roomID,
						ImagePath: imagePath,
						SortOrder: imageRequest.SortOrder,
						IsCover:   imageRequest.IsCover,
					}

					if err := tx.Create(&newImage).Error; err != nil {
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
		Message: UpdateHomestaysUseCaseSuccessMessage,
	}, nil
}

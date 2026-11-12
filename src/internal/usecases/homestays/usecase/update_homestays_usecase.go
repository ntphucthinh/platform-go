package usecase

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/domain/homestay"
	homestayimage "platform-go/internal/domain/homestay_image"
	homestayroom "platform-go/internal/domain/homestay_room"
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
	db           repository.DBInterface
	homestayRepo homestay.HomestayRepository
	roomRepo     homestayroom.HomestayRoomRepository
	imageRepo    homestayimage.HomestayImageRepository
}

func NewUpdateHomestayUsecase(
	db repository.DBInterface,
	homestayRepo homestay.HomestayRepository,
	roomRepo homestayroom.HomestayRoomRepository,
	imageRepo homestayimage.HomestayImageRepository,
) *UpdateHomestayUsecase {
	return &UpdateHomestayUsecase{
		db:           db,
		homestayRepo: homestayRepo,
		roomRepo:     roomRepo,
		imageRepo:    imageRepo,
	}
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
		txCtx := repository.WithTx(ctx, tx)

		// Find the homestay to update.
		existingHomestay, err := uc.homestayRepo.FindHomestayByID(txCtx, homestayID)
		if err != nil {
			return err
		}
		if existingHomestay == nil {
			return ErrHomestayNotFound
		}

		// Update homestay information.
		existingHomestay.Update(
			request.Name,
			request.Address,
			request.Description,
			request.GoogleMapsUrl,
		)
		if err := uc.homestayRepo.UpdateHomestay(txCtx, existingHomestay); err != nil {
			return err
		}

		// A nil Rooms field means that room data should remain unchanged.
		// An empty Rooms list means that all existing rooms should be removed.
		if request.Rooms == nil {
			return nil
		}

		// Load existing rooms belonging to this homestay.
		existingRooms, err := uc.roomRepo.FindRoomsByHomestayID(txCtx, homestayID)
		if err != nil {
			return err
		}

		existingRoomMap := make(map[int64]homestayroom.HomestayRoom, len(existingRooms))
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
			if err := uc.imageRepo.DeleteImagesByRoomIDs(txCtx, roomIDsToDelete); err != nil {
				return err
			}

			if err := uc.roomRepo.DeleteRoomsByIDs(txCtx, roomIDsToDelete); err != nil {
				return err
			}
		}

		// Create or update rooms.
		for _, roomRequest := range *request.Rooms {
			room := homestayroom.NewHomestayRoom(
				homestayID,
				roomRequest.Name,
				roomRequest.Description,
				roomRequest.PriceAmount,
				(*string)(roomRequest.PriceUnit),
				roomRequest.SortOrder,
			)

			if roomRequest.Id != nil && *roomRequest.Id > 0 {
				room.ID = *roomRequest.Id
				if err := uc.roomRepo.UpdateHomestayRoom(txCtx, &room); err != nil {
					return err
				}
			} else {
				if err := uc.roomRepo.CreateHomestayRoom(txCtx, &room); err != nil {
					return err
				}
			}

			roomID := room.ID

			// A nil Images field means that the room's images remain unchanged.
			// An empty Images list means that all images of this room are removed.
			if roomRequest.Images == nil {
				continue
			}

			// Load images belonging to the current room only.
			existingImages, err := uc.imageRepo.FindImagesByRoomID(txCtx, roomID)
			if err != nil {
				return err
			}

			existingImageMap := make(
				map[int64]homestayimage.HomestayImage,
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
				if err := uc.imageRepo.DeleteImagesByIDs(txCtx, imageIDsToDelete); err != nil {
					return err
				}
			}

			// Create or update images.
			for _, imageRequest := range *roomRequest.Images {
				img := homestayimage.NewHomestayImage(
					roomID,
					imageRequest.ImagePath,
					imageRequest.SortOrder,
					imageRequest.IsCover,
				)

				if imageRequest.Id != nil && *imageRequest.Id > 0 {
					img.ID = *imageRequest.Id
					if err := uc.imageRepo.UpdateHomestayImage(txCtx, &img); err != nil {
						return err
					}
				} else {
					if err := uc.imageRepo.CreateHomestayImages(txCtx, []homestayimage.HomestayImage{img}); err != nil {
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

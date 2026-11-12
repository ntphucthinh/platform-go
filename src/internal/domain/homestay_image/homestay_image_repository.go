package homestayimage

import "context"

// HomestayImageRepository defines persistence operations for homestay room images.
type HomestayImageRepository interface {
	// FindImagesByRoomID retrieves all images belonging to a given room ID.
	FindImagesByRoomID(ctx context.Context, roomID int64) ([]HomestayImage, error)
	// CreateHomestayImages batch-inserts homestay room image records and populates their generated IDs and timestamps.
	CreateHomestayImages(ctx context.Context, images []HomestayImage) error
	// UpdateHomestayImage updates an existing image's path, sort order, and cover flag scoped by image ID and room ID.
	UpdateHomestayImage(ctx context.Context, image *HomestayImage) error
	// DeleteImagesByRoomIDs deletes all homestay image records belonging to the given room IDs.
	DeleteImagesByRoomIDs(ctx context.Context, roomIDs []int64) error
	// DeleteImagesByIDs deletes multiple homestay image records matching the given image IDs.
	DeleteImagesByIDs(ctx context.Context, ids []int64) error
}


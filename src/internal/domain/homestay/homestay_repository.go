package homestay

import "context"

// HomestayRepository defines persistence operations for homestays, rooms, and images.
type HomestayRepository interface {
	// FindHomestayByID retrieves a homestay by its ID, returning (nil, nil) if not found.
	FindHomestayByID(ctx context.Context, id int64) (*Homestay, error)
	// FindRoomsByHomestayID retrieves all rooms belonging to a given homestay ID.
	FindRoomsByHomestayID(ctx context.Context, homestayID int64) ([]HomestayRoom, error)
	// FindImagesByRoomID retrieves all images belonging to a given room ID.
	FindImagesByRoomID(ctx context.Context, roomID int64) ([]HomestayImage, error)

	// CreateHomestay inserts a new homestay record and populates its generated ID and timestamps.
	CreateHomestay(ctx context.Context, homestay *Homestay) error
	// CreateHomestayRoom inserts a new homestay room record and populates its generated ID and timestamps.
	CreateHomestayRoom(ctx context.Context, room *HomestayRoom) error
	// CreateHomestayImage batch-inserts homestay room image records and populates their generated IDs and timestamps.
	CreateHomestayImage(ctx context.Context, images []HomestayImage) error

	// UpdateHomestay updates an existing homestay's general information by ID.
	UpdateHomestay(ctx context.Context, homestay *Homestay) error
	// UpdateHomestayRoom updates an existing room's details scoped by room ID and homestay ID.
	UpdateHomestayRoom(ctx context.Context, room *HomestayRoom) error
	// UpdateHomestayImage updates an existing image's path, sort order, and cover flag scoped by image ID and room ID.
	UpdateHomestayImage(ctx context.Context, image *HomestayImage) error

	// DeleteHomestayByID deletes a homestay record by its ID.
	DeleteHomestayByID(ctx context.Context, id int64) error
	// DeleteRoomsByIDs deletes multiple homestay room records matching the given room IDs.
	DeleteRoomsByIDs(ctx context.Context, ids []int64) error
	// DeleteImagesByRoomIDs deletes all homestay image records belonging to the given room IDs.
	DeleteImagesByRoomIDs(ctx context.Context, roomIDs []int64) error
	// DeleteImagesByIDs deletes multiple homestay image records matching the given image IDs.
	DeleteImagesByIDs(ctx context.Context, ids []int64) error
}

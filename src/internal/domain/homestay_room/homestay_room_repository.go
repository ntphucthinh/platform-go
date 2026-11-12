package homestayroom

import "context"

// HomestayRoomRepository defines persistence operations for homestay rooms.
type HomestayRoomRepository interface {
	// FindRoomsByHomestayID retrieves all rooms belonging to a given homestay ID.
	FindRoomsByHomestayID(ctx context.Context, homestayID int64) ([]HomestayRoom, error)
	// CreateHomestayRoom inserts a new homestay room record and populates its generated ID and timestamps.
	CreateHomestayRoom(ctx context.Context, room *HomestayRoom) error
	// UpdateHomestayRoom updates an existing room's details scoped by room ID and homestay ID.
	UpdateHomestayRoom(ctx context.Context, room *HomestayRoom) error
	// DeleteRoomsByIDs deletes multiple homestay room records matching the given room IDs.
	DeleteRoomsByIDs(ctx context.Context, ids []int64) error
}


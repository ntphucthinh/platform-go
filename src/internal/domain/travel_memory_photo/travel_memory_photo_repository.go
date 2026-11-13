package travelmemoryphoto

import "context"

// TravelMemoryPhotoRepository defines persistence operations for travel memory photos.
type TravelMemoryPhotoRepository interface {
	// FindPhotosByDayIDs retrieves all photos belonging to the given memory day IDs.
	FindPhotosByDayIDs(ctx context.Context, dayIDs []int64) ([]TravelMemoryPhoto, error)
	// CreateTravelMemoryPhotos batch-inserts travel memory photo records and populates their generated IDs and timestamps.
	CreateTravelMemoryPhotos(ctx context.Context, photos []TravelMemoryPhoto) error
	// UpdateTravelMemoryPhotos updates existing photo records' image paths and sort orders.
	UpdateTravelMemoryPhotos(ctx context.Context, photos []TravelMemoryPhoto) error
	// DeletePhotosByDayIDs deletes all photos belonging to the given memory day IDs.
	DeletePhotosByDayIDs(ctx context.Context, dayIDs []int64) error
	// DeletePhotosByIDs deletes multiple photo records matching the given photo IDs.
	DeletePhotosByIDs(ctx context.Context, ids []int64) error
}

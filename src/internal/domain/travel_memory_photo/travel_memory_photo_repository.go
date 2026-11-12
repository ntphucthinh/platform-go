package travelmemoryphoto

import "context"

// TravelMemoryPhotoRepository defines persistence operations for travel memory photos.
type TravelMemoryPhotoRepository interface {
	// CreateTravelMemoryPhotos batch-inserts travel memory photo records and populates their generated IDs and timestamps.
	CreateTravelMemoryPhotos(ctx context.Context, photos []TravelMemoryPhoto) error
}

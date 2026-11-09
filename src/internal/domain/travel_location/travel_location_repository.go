package travellocation

import "context"

// TravelLocationRepository defines persistence operations for travel locations.
type TravelLocationRepository interface {
	// FindTravelLocationByID retrieves a travel location by its ID, returning (nil, nil) if not found.
	FindTravelLocationByID(ctx context.Context, id int64) (*TravelLocation, error)
	// CreateTravelLocation inserts a new travel location record and populates its generated ID and timestamps.
	CreateTravelLocation(ctx context.Context, location *TravelLocation) error
	// UpdateTravelLocation updates an existing travel location record by its ID.
	UpdateTravelLocation(ctx context.Context, location *TravelLocation) error
	// DeleteTravelLocationByID soft-deletes a travel location record by its ID.
	DeleteTravelLocationByID(ctx context.Context, id int64) error
}


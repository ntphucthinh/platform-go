package itinerary

import "context"

// ItineraryRepository defines persistence operations for itineraries.
type ItineraryRepository interface {
	// FindItineraryByID retrieves an itinerary by its ID, returning (nil, nil) if not found.
	FindItineraryByID(ctx context.Context, id int64) (*Itinerary, error)
	// CreateItinerary inserts a new itinerary record and populates its generated ID and timestamps.
	CreateItinerary(ctx context.Context, itinerary *Itinerary) error
	// UpdateItinerary updates an existing itinerary's information by ID.
	UpdateItinerary(ctx context.Context, itinerary *Itinerary) error
	// DeleteItineraryByID soft-deletes an itinerary record by its ID.
	DeleteItineraryByID(ctx context.Context, id int64) error
}

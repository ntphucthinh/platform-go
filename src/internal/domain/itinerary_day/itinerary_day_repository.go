package itineraryday

import "context"

// ItineraryDayRepository defines persistence operations for itinerary days.
type ItineraryDayRepository interface {
	// FindDaysByItineraryID retrieves all itinerary days belonging to a given itinerary ID.
	FindDaysByItineraryID(ctx context.Context, itineraryID int64) ([]ItineraryDay, error)
	// CreateItineraryDays batch-inserts itinerary day records and populates their generated IDs and timestamps.
	CreateItineraryDays(ctx context.Context, days []ItineraryDay) error
	// UpdateItineraryDays updates existing itinerary day records.
	UpdateItineraryDays(ctx context.Context, days []ItineraryDay) error
	// DeleteDaysByIDs deletes multiple itinerary day records matching the given day IDs.
	DeleteDaysByIDs(ctx context.Context, ids []int64) error
}

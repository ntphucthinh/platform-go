package itineraryactivity

import "context"

// ItineraryActivityRepository defines persistence operations for itinerary activities.
type ItineraryActivityRepository interface {
	// FindActivitiesByDayIDs retrieves all activities belonging to the given itinerary day IDs.
	FindActivitiesByDayIDs(ctx context.Context, dayIDs []int64) ([]ItineraryActivity, error)
	// CreateItineraryActivities batch-inserts itinerary activity records and populates their generated IDs and timestamps.
	CreateItineraryActivities(ctx context.Context, activities []ItineraryActivity) error
	// UpdateItineraryActivities updates existing itinerary activity records.
	UpdateItineraryActivities(ctx context.Context, activities []ItineraryActivity) error
	// DeleteActivitiesByDayIDs deletes all activities belonging to the given itinerary day IDs.
	DeleteActivitiesByDayIDs(ctx context.Context, dayIDs []int64) error
	// DeleteActivitiesByIDs deletes multiple activity records matching the given activity IDs.
	DeleteActivitiesByIDs(ctx context.Context, ids []int64) error
}

package travelmemoryday

import "context"

// TravelMemoryDayRepository defines persistence operations for travel memory days.
type TravelMemoryDayRepository interface {
	// CreateTravelMemoryDays batch-inserts travel memory day records and populates their generated IDs and timestamps.
	CreateTravelMemoryDays(ctx context.Context, days []TravelMemoryDay) error
}

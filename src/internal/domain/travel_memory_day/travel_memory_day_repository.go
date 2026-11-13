package travelmemoryday

import "context"

// TravelMemoryDayRepository defines persistence operations for travel memory days.
type TravelMemoryDayRepository interface {
	// FindDaysByMemoryID retrieves all travel memory days belonging to a given memory ID.
	FindDaysByMemoryID(ctx context.Context, memoryID int64) ([]TravelMemoryDay, error)
	// CreateTravelMemoryDays batch-inserts travel memory day records and populates their generated IDs and timestamps.
	CreateTravelMemoryDays(ctx context.Context, days []TravelMemoryDay) error
	// UpdateTravelMemoryDays updates existing travel memory day records.
	UpdateTravelMemoryDays(ctx context.Context, days []TravelMemoryDay) error
	// DeleteDaysByIDs deletes multiple travel memory day records matching the given day IDs.
	DeleteDaysByIDs(ctx context.Context, ids []int64) error
}

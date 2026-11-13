package travelmemory

import "context"

// TravelMemoryRepository defines persistence operations for travel memories.
type TravelMemoryRepository interface {
	// FindTravelMemoryByID retrieves a travel memory by its ID, returning (nil, nil) if not found.
	FindTravelMemoryByID(ctx context.Context, id int64) (*TravelMemory, error)
	// CreateTravelMemory inserts a new travel memory record and populates its generated ID and timestamps.
	CreateTravelMemory(ctx context.Context, memory *TravelMemory) error
	// UpdateTravelMemory updates an existing travel memory's information by ID.
	UpdateTravelMemory(ctx context.Context, memory *TravelMemory) error
	// DeleteTravelMemoryByID deletes a travel memory record by its ID.
	DeleteTravelMemoryByID(ctx context.Context, id int64) error
}

package travelmemory

import "context"

// TravelMemoryRepository defines persistence operations for travel memories.
type TravelMemoryRepository interface {
	// CreateTravelMemory inserts a new travel memory record and populates its generated ID and timestamps.
	CreateTravelMemory(ctx context.Context, memory *TravelMemory) error
}

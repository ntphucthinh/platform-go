package homestay

import "context"

// HomestayRepository defines persistence operations for homestays.
type HomestayRepository interface {
	// FindHomestayByID retrieves a homestay by its ID, returning (nil, nil) if not found.
	FindHomestayByID(ctx context.Context, id int64) (*Homestay, error)
	// CreateHomestay inserts a new homestay record and populates its generated ID and timestamps.
	CreateHomestay(ctx context.Context, homestay *Homestay) error
	// UpdateHomestay updates an existing homestay's general information by ID.
	UpdateHomestay(ctx context.Context, homestay *Homestay) error
	// DeleteHomestayByID deletes a homestay record by its ID.
	DeleteHomestayByID(ctx context.Context, id int64) error
}

package queryservice

import (
	"context"

	"platform-go/internal/infrastructure/dbutil"
	"platform-go/internal/infrastructure/qsdto"
)

// ItineraryQueryService defines the query interface for fetching itineraries.
type ItineraryQueryService interface {
	FetchItineraries(ctx context.Context, page, perPage int, locationID *int64) (*dbutil.Paginator[qsdto.ItineraryDto], error)
	FetchItineraryByID(ctx context.Context, id int64) (*qsdto.ItineraryDetailDto, error)
}

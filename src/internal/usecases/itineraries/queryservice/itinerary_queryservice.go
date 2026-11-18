package queryservice

import (
	"context"

	"platform-go/internal/infrastructure/dbutil"
	"platform-go/internal/infrastructure/qsdto"
)

type ItineraryQueryService interface {
	FetchItineraries(ctx context.Context, page, perPage int, locationID *int64) (*dbutil.Paginator[qsdto.ItineraryDto], error)
}

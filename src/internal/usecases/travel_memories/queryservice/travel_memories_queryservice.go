package queryservice

import (
	"context"

	"platform-go/internal/infrastructure/dbutil"
	"platform-go/internal/infrastructure/qsdto"
)

type TravelMemoryQueryService interface {
	FetchTravelMemories(ctx context.Context, page, perPage int, countryCode *string, locationID *int64, year *int) (*dbutil.Paginator[qsdto.TravelMemoryDto], error)
	FetchTravelMemoryByID(ctx context.Context, id int64) (*qsdto.TravelMemoryDto, error)
}

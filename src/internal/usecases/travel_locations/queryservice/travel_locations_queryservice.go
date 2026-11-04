package queryservice

import (
	"context"

	"platform-go/internal/infrastructure/qsdto"
)

type TravelLocationQueryService interface {
	FetchTravelLocations(ctx context.Context, filter qsdto.TravelLocationFilter) ([]qsdto.TravelLocationDto, int64, error)
	FetchTravelLocationByID(ctx context.Context, id int64) (*qsdto.TravelLocationDto, error)
}


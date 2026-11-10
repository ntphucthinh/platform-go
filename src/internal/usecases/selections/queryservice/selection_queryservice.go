package queryservice

import (
	"context"

	"platform-go/internal/infrastructure/qsdto"
)

type SelectionQueryService interface {
	FetchTravelLocations(ctx context.Context, countryCode *string) ([]qsdto.SelectionTravelLocationDto, error)
}


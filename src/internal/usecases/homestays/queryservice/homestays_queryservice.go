package queryservice

import (
	"context"

	qsdto "platform-go/internal/infrastructure/qsdto"
)

type HomestayQueryService interface {
	FetchHomestays(ctx context.Context) ([]qsdto.HomestayDto, error)
}

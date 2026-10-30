package queryservice

import (
	"context"

	"platform-go/internal/infrastructure/dbutil"
	qsdto "platform-go/internal/infrastructure/qsdto"
)

type HomestayQueryService interface {
	FetchHomestays(ctx context.Context, page, perPage int) (*dbutil.Paginator[qsdto.HomestayDto], error)
	FetchHomestayByID(ctx context.Context, id int64) (*qsdto.HomestayDto, error)
}

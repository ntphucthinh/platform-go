package repository

import (
	"context"

	"platform-go/internal/infrastructure/dbmodel"
	travelmemoryphoto "platform-go/internal/domain/travel_memory_photo"
)

type TravelMemoryPhotoRepository = travelmemoryphoto.TravelMemoryPhotoRepository

type travelMemoryPhotoRepository struct {
	db DBInterface
}

func NewTravelMemoryPhotoRepository(db DBInterface) travelmemoryphoto.TravelMemoryPhotoRepository {
	return &travelMemoryPhotoRepository{db: db}
}

func (r *travelMemoryPhotoRepository) CreateTravelMemoryPhotos(ctx context.Context, photos []travelmemoryphoto.TravelMemoryPhoto) error {
	if len(photos) == 0 {
		return nil
	}

	const batchSize = 100
	models := make([]dbmodel.TravelMemoryPhoto, len(photos))
	for i, p := range photos {
		models[i] = dbmodel.TravelMemoryPhoto{
			MemoryDayID: p.MemoryDayID,
			ImagePath:   p.ImagePath,
			SortOrder:   p.SortOrder,
		}
	}

	if err := r.db.GetDBByCtx(ctx).CreateInBatches(&models, batchSize).Error; err != nil {
		return err
	}

	for i := range photos {
		photos[i].ID = models[i].ID
		photos[i].CreatedAt = models[i].CreatedAt
		photos[i].UpdatedAt = models[i].UpdatedAt
	}

	return nil
}

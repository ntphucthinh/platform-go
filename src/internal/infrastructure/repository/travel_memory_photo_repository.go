package repository

import (
	"context"

	travelmemoryphoto "platform-go/internal/domain/travel_memory_photo"
	"platform-go/internal/infrastructure/dbmodel"
)

type TravelMemoryPhotoRepository = travelmemoryphoto.TravelMemoryPhotoRepository

type travelMemoryPhotoRepository struct {
	db DBInterface
}

func NewTravelMemoryPhotoRepository(db DBInterface) travelmemoryphoto.TravelMemoryPhotoRepository {
	return &travelMemoryPhotoRepository{db: db}
}

func (r *travelMemoryPhotoRepository) FindPhotosByDayIDs(ctx context.Context, dayIDs []int64) ([]travelmemoryphoto.TravelMemoryPhoto, error) {
	if len(dayIDs) == 0 {
		return nil, nil
	}

	var models []dbmodel.TravelMemoryPhoto
	if err := r.db.GetDBByCtx(ctx).
		Where("memory_day_id IN ?", dayIDs).
		Order("memory_day_id ASC, sort_order ASC, id ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	photos := make([]travelmemoryphoto.TravelMemoryPhoto, len(models))
	for i, m := range models {
		photos[i] = travelmemoryphoto.TravelMemoryPhoto{
			ID:          m.ID,
			MemoryDayID: m.MemoryDayID,
			ImagePath:   m.ImagePath,
			SortOrder:   m.SortOrder,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		}
	}
	return photos, nil
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

func (r *travelMemoryPhotoRepository) UpdateTravelMemoryPhotos(ctx context.Context, photos []travelmemoryphoto.TravelMemoryPhoto) error {
	if len(photos) == 0 {
		return nil
	}

	db := r.db.GetDBByCtx(ctx)
	for _, photo := range photos {
		updates := map[string]interface{}{
			"image_path": photo.ImagePath,
			"sort_order": photo.SortOrder,
		}

		if err := db.Model(&dbmodel.TravelMemoryPhoto{}).
			Where("id = ? AND memory_day_id = ?", photo.ID, photo.MemoryDayID).
			Updates(updates).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *travelMemoryPhotoRepository) DeletePhotosByDayIDs(ctx context.Context, dayIDs []int64) error {
	if len(dayIDs) == 0 {
		return nil
	}
	return r.db.GetDBByCtx(ctx).
		Where("memory_day_id IN ?", dayIDs).
		Delete(&dbmodel.TravelMemoryPhoto{}).Error
}

func (r *travelMemoryPhotoRepository) DeletePhotosByIDs(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.GetDBByCtx(ctx).
		Where("id IN ?", ids).
		Delete(&dbmodel.TravelMemoryPhoto{}).Error
}

package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	travelmemory "platform-go/internal/domain/travel_memory"
	"platform-go/internal/infrastructure/dbmodel"
)

type TravelMemoryRepository = travelmemory.TravelMemoryRepository

type travelMemoryRepository struct {
	db DBInterface
}

func NewTravelMemoryRepository(db DBInterface) travelmemory.TravelMemoryRepository {
	return &travelMemoryRepository{db: db}
}

func (r *travelMemoryRepository) FindTravelMemoryByID(ctx context.Context, id int64) (*travelmemory.TravelMemory, error) {
	var model dbmodel.TravelMemory
	if err := r.db.GetDBByCtx(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &travelmemory.TravelMemory{
		ID:             model.ID,
		LocationID:     model.LocationID,
		Title:          model.Title,
		StartDate:      model.StartDate,
		EndDate:        model.EndDate,
		Description:    model.Description,
		CoverImagePath: model.CoverImagePath,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}, nil
}

func (r *travelMemoryRepository) CreateTravelMemory(ctx context.Context, memory *travelmemory.TravelMemory) error {
	model := dbmodel.TravelMemory{
		LocationID:     memory.LocationID,
		Title:          memory.Title,
		StartDate:      memory.StartDate,
		EndDate:        memory.EndDate,
		Description:    memory.Description,
		CoverImagePath: memory.CoverImagePath,
	}

	if err := r.db.GetDBByCtx(ctx).Create(&model).Error; err != nil {
		return err
	}

	memory.ID = model.ID
	memory.CreatedAt = model.CreatedAt
	memory.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *travelMemoryRepository) UpdateTravelMemory(ctx context.Context, memory *travelmemory.TravelMemory) error {
	updates := map[string]interface{}{
		"location_id":      memory.LocationID,
		"title":            memory.Title,
		"start_date":       memory.StartDate,
		"end_date":         memory.EndDate,
		"description":      memory.Description,
		"cover_image_path": memory.CoverImagePath,
	}

	return r.db.GetDBByCtx(ctx).
		Model(&dbmodel.TravelMemory{}).
		Where("id = ?", memory.ID).
		Updates(updates).Error
}

func (r *travelMemoryRepository) DeleteTravelMemoryByID(ctx context.Context, id int64) error {
	return r.db.GetDBByCtx(ctx).
		Where("id = ?", id).
		Delete(&dbmodel.TravelMemory{}).Error
}

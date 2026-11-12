package repository

import (
	"context"

	"platform-go/internal/infrastructure/dbmodel"
	travelmemory "platform-go/internal/domain/travel_memory"
)

type TravelMemoryRepository = travelmemory.TravelMemoryRepository

type travelMemoryRepository struct {
	db DBInterface
}

func NewTravelMemoryRepository(db DBInterface) travelmemory.TravelMemoryRepository {
	return &travelMemoryRepository{db: db}
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

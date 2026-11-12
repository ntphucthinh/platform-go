package repository

import (
	"context"

	"platform-go/internal/infrastructure/dbmodel"
	travelmemoryday "platform-go/internal/domain/travel_memory_day"
)

type TravelMemoryDayRepository = travelmemoryday.TravelMemoryDayRepository

type travelMemoryDayRepository struct {
	db DBInterface
}

func NewTravelMemoryDayRepository(db DBInterface) travelmemoryday.TravelMemoryDayRepository {
	return &travelMemoryDayRepository{db: db}
}

func (r *travelMemoryDayRepository) CreateTravelMemoryDays(ctx context.Context, days []travelmemoryday.TravelMemoryDay) error {
	if len(days) == 0 {
		return nil
	}

	const batchSize = 100
	models := make([]dbmodel.TravelMemoryDay, len(days))
	for i, d := range days {
		models[i] = dbmodel.TravelMemoryDay{
			MemoryID:    d.MemoryID,
			Date:        d.Date,
			Description: d.Description,
		}
	}

	if err := r.db.GetDBByCtx(ctx).CreateInBatches(&models, batchSize).Error; err != nil {
		return err
	}

	for i := range days {
		days[i].ID = models[i].ID
		days[i].CreatedAt = models[i].CreatedAt
		days[i].UpdatedAt = models[i].UpdatedAt
	}

	return nil
}

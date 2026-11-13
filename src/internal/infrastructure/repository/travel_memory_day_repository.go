package repository

import (
	"context"

	travelmemoryday "platform-go/internal/domain/travel_memory_day"
	"platform-go/internal/infrastructure/dbmodel"
)

type TravelMemoryDayRepository = travelmemoryday.TravelMemoryDayRepository

type travelMemoryDayRepository struct {
	db DBInterface
}

func NewTravelMemoryDayRepository(db DBInterface) travelmemoryday.TravelMemoryDayRepository {
	return &travelMemoryDayRepository{db: db}
}

func (r *travelMemoryDayRepository) FindDaysByMemoryID(ctx context.Context, memoryID int64) ([]travelmemoryday.TravelMemoryDay, error) {
	var models []dbmodel.TravelMemoryDay
	if err := r.db.GetDBByCtx(ctx).
		Where("memory_id = ?", memoryID).
		Order("date ASC, id ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	days := make([]travelmemoryday.TravelMemoryDay, len(models))
	for i, m := range models {
		days[i] = travelmemoryday.TravelMemoryDay{
			ID:          m.ID,
			MemoryID:    m.MemoryID,
			Date:        m.Date,
			Description: m.Description,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		}
	}
	return days, nil
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

func (r *travelMemoryDayRepository) UpdateTravelMemoryDays(ctx context.Context, days []travelmemoryday.TravelMemoryDay) error {
	if len(days) == 0 {
		return nil
	}

	db := r.db.GetDBByCtx(ctx)
	for _, day := range days {
		updates := map[string]interface{}{
			"date":        day.Date,
			"description": day.Description,
		}

		if err := db.Model(&dbmodel.TravelMemoryDay{}).
			Where("id = ? AND memory_id = ?", day.ID, day.MemoryID).
			Updates(updates).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *travelMemoryDayRepository) DeleteDaysByIDs(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.GetDBByCtx(ctx).
		Where("id IN ?", ids).
		Delete(&dbmodel.TravelMemoryDay{}).Error
}

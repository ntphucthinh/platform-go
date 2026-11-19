package repository

import (
	"context"

	itineraryday "platform-go/internal/domain/itinerary_day"
	"platform-go/internal/infrastructure/dbmodel"
)

type ItineraryDayRepository = itineraryday.ItineraryDayRepository

type itineraryDayRepository struct {
	db DBInterface
}

func NewItineraryDayRepository(db DBInterface) itineraryday.ItineraryDayRepository {
	return &itineraryDayRepository{db: db}
}

func (r *itineraryDayRepository) FindDaysByItineraryID(ctx context.Context, itineraryID int64) ([]itineraryday.ItineraryDay, error) {
	var models []dbmodel.ItineraryDay
	if err := r.db.GetDBByCtx(ctx).
		Where("itinerary_id = ?", itineraryID).
		Order("day_number ASC, id ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	days := make([]itineraryday.ItineraryDay, len(models))
	for i, m := range models {
		dayTips, err := unmarshalOptionalStringSlice(m.DayTips)
		if err != nil {
			return nil, err
		}

		days[i] = itineraryday.ItineraryDay{
			ID:          m.ID,
			ItineraryID: m.ItineraryID,
			DayNumber:   m.DayNumber,
			Title:       m.Title,
			Subtitle:    m.Subtitle,
			Summary:     m.Summary,
			ImagePath:   m.ImagePath,
			DayTips:     dayTips,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		}
	}
	return days, nil
}

func (r *itineraryDayRepository) CreateItineraryDays(ctx context.Context, days []itineraryday.ItineraryDay) error {
	if len(days) == 0 {
		return nil
	}

	const batchSize = 100
	models := make([]dbmodel.ItineraryDay, len(days))
	for i, d := range days {
		dayTipsJSON, err := marshalOptionalStringSlice(d.DayTips)
		if err != nil {
			return err
		}

		models[i] = dbmodel.ItineraryDay{
			ItineraryID: d.ItineraryID,
			DayNumber:   d.DayNumber,
			Title:       d.Title,
			Subtitle:    d.Subtitle,
			Summary:     d.Summary,
			ImagePath:   d.ImagePath,
			DayTips:     dayTipsJSON,
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

func (r *itineraryDayRepository) UpdateItineraryDays(ctx context.Context, days []itineraryday.ItineraryDay) error {
	if len(days) == 0 {
		return nil
	}

	db := r.db.GetDBByCtx(ctx)
	for _, day := range days {
		dayTipsJSON, err := marshalOptionalStringSlice(day.DayTips)
		if err != nil {
			return err
		}

		updates := map[string]interface{}{
			"day_number": day.DayNumber,
			"title":      day.Title,
			"subtitle":   day.Subtitle,
			"summary":    day.Summary,
			"image_path": day.ImagePath,
			"day_tips":   dayTipsJSON,
		}

		res := db.Model(&dbmodel.ItineraryDay{}).
			Where("id = ? AND itinerary_id = ?", day.ID, day.ItineraryID).
			Updates(updates)
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			var count int64
			if err := db.Model(&dbmodel.ItineraryDay{}).
				Where("id = ? AND itinerary_id = ?", day.ID, day.ItineraryID).
				Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				return itineraryday.ErrDayNotFound
			}
		}
	}

	return nil
}

func (r *itineraryDayRepository) DeleteDaysByIDs(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.GetDBByCtx(ctx).
		Where("id IN ?", ids).
		Delete(&dbmodel.ItineraryDay{}).Error
}

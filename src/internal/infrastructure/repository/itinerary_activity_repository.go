package repository

import (
	"context"

	itineraryactivity "platform-go/internal/domain/itinerary_activity"
	"platform-go/internal/infrastructure/dbmodel"
)

type ItineraryActivityRepository = itineraryactivity.ItineraryActivityRepository

type itineraryActivityRepository struct {
	db DBInterface
}

func NewItineraryActivityRepository(db DBInterface) itineraryactivity.ItineraryActivityRepository {
	return &itineraryActivityRepository{db: db}
}

func (r *itineraryActivityRepository) FindActivitiesByDayIDs(ctx context.Context, dayIDs []int64) ([]itineraryactivity.ItineraryActivity, error) {
	if len(dayIDs) == 0 {
		return nil, nil
	}

	var models []dbmodel.ItineraryActivity
	if err := r.db.GetDBByCtx(ctx).
		Where("itinerary_day_id IN ?", dayIDs).
		Order("itinerary_day_id ASC, start_time ASC, id ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	activities := make([]itineraryactivity.ItineraryActivity, len(models))
	for i, m := range models {
		activities[i] = itineraryactivity.ItineraryActivity{
			ID:               m.ID,
			ItineraryDayID:   m.ItineraryDayID,
			StartTime:        m.StartTime,
			Title:            m.Title,
			ActivityCategory: m.ActivityCategory,
			DurationMinutes:  m.DurationMinutes,
			Location:         m.Location,
			GoogleMapsURL:    m.GoogleMapsURL,
			Description:      m.Description,
			Notes:            m.Notes,
			ImagePath:        m.ImagePath,
			CreatedAt:        m.CreatedAt,
			UpdatedAt:        m.UpdatedAt,
		}
	}
	return activities, nil
}

func (r *itineraryActivityRepository) CreateItineraryActivities(ctx context.Context, activities []itineraryactivity.ItineraryActivity) error {
	if len(activities) == 0 {
		return nil
	}

	const batchSize = 100
	models := make([]dbmodel.ItineraryActivity, len(activities))
	for i, a := range activities {
		models[i] = dbmodel.ItineraryActivity{
			ItineraryDayID:   a.ItineraryDayID,
			StartTime:        a.StartTime,
			Title:            a.Title,
			ActivityCategory: a.ActivityCategory,
			DurationMinutes:  a.DurationMinutes,
			Location:         a.Location,
			GoogleMapsURL:    a.GoogleMapsURL,
			Description:      a.Description,
			Notes:            a.Notes,
			ImagePath:        a.ImagePath,
		}
	}

	if err := r.db.GetDBByCtx(ctx).CreateInBatches(&models, batchSize).Error; err != nil {
		return err
	}

	for i := range activities {
		activities[i].ID = models[i].ID
		activities[i].CreatedAt = models[i].CreatedAt
		activities[i].UpdatedAt = models[i].UpdatedAt
	}

	return nil
}

func (r *itineraryActivityRepository) UpdateItineraryActivities(ctx context.Context, activities []itineraryactivity.ItineraryActivity) error {
	if len(activities) == 0 {
		return nil
	}

	db := r.db.GetDBByCtx(ctx)
	for _, a := range activities {
		updates := map[string]interface{}{
			"start_time":        a.StartTime,
			"title":             a.Title,
			"activity_category": a.ActivityCategory,
			"duration_minutes":  a.DurationMinutes,
			"location":          a.Location,
			"google_maps_url":   a.GoogleMapsURL,
			"description":       a.Description,
			"notes":             a.Notes,
			"image_path":        a.ImagePath,
		}

		res := db.Model(&dbmodel.ItineraryActivity{}).
			Where("id = ? AND itinerary_day_id = ?", a.ID, a.ItineraryDayID).
			Updates(updates)
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			var count int64
			if err := db.Model(&dbmodel.ItineraryActivity{}).
				Where("id = ? AND itinerary_day_id = ?", a.ID, a.ItineraryDayID).
				Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				return itineraryactivity.ErrActivityNotFound
			}
		}
	}

	return nil
}

func (r *itineraryActivityRepository) DeleteActivitiesByDayIDs(ctx context.Context, dayIDs []int64) error {
	if len(dayIDs) == 0 {
		return nil
	}
	return r.db.GetDBByCtx(ctx).
		Where("itinerary_day_id IN ?", dayIDs).
		Delete(&dbmodel.ItineraryActivity{}).Error
}

func (r *itineraryActivityRepository) DeleteActivitiesByIDs(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.GetDBByCtx(ctx).
		Where("id IN ?", ids).
		Delete(&dbmodel.ItineraryActivity{}).Error
}

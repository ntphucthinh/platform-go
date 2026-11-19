package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"gorm.io/gorm"

	"platform-go/internal/domain/itinerary"
	"platform-go/internal/infrastructure/dbmodel"
)

type ItineraryRepository = itinerary.ItineraryRepository

type itineraryRepository struct {
	db DBInterface
}

func NewItineraryRepository(db DBInterface) itinerary.ItineraryRepository {
	return &itineraryRepository{db: db}
}

func (r *itineraryRepository) FindItineraryByID(ctx context.Context, id int64) (*itinerary.Itinerary, error) {
	var model dbmodel.Itinerary
	if err := r.db.GetDBByCtx(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	travelTips, err := unmarshalOptionalStringSlice(model.TravelTips)
	if err != nil {
		return nil, err
	}

	return &itinerary.Itinerary{
		ID:             model.ID,
		ItineraryName:  model.ItineraryName,
		LocationID:     model.LocationID,
		SuitableFor:    model.SuitableFor,
		Description:    model.Description,
		CoverImagePath: model.CoverImagePath,
		TravelTips:     travelTips,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}, nil
}

func (r *itineraryRepository) CreateItinerary(ctx context.Context, it *itinerary.Itinerary) error {
	travelTipsJSON, err := marshalOptionalStringSlice(it.TravelTips)
	if err != nil {
		return err
	}

	model := dbmodel.Itinerary{
		ItineraryName:  it.ItineraryName,
		LocationID:     it.LocationID,
		SuitableFor:    it.SuitableFor,
		Description:    it.Description,
		CoverImagePath: it.CoverImagePath,
		TravelTips:     travelTipsJSON,
	}

	if err := r.db.GetDBByCtx(ctx).Create(&model).Error; err != nil {
		return err
	}

	it.ID = model.ID
	it.CreatedAt = model.CreatedAt
	it.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *itineraryRepository) UpdateItinerary(ctx context.Context, it *itinerary.Itinerary) error {
	travelTipsJSON, err := marshalOptionalStringSlice(it.TravelTips)
	if err != nil {
		return err
	}

	updates := map[string]interface{}{
		"itinerary_name":   it.ItineraryName,
		"location_id":      it.LocationID,
		"suitable_for":     it.SuitableFor,
		"description":      it.Description,
		"cover_image_path": it.CoverImagePath,
		"travel_tips":      travelTipsJSON,
	}

	db := r.db.GetDBByCtx(ctx)
	res := db.Model(&dbmodel.Itinerary{}).
		Where("id = ?", it.ID).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		var count int64
		if err := db.Model(&dbmodel.Itinerary{}).
			Where("id = ?", it.ID).
			Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return itinerary.ErrItineraryNotFound
		}
	}

	return nil
}

func (r *itineraryRepository) DeleteItineraryByID(ctx context.Context, id int64) error {
	return r.db.GetDBByCtx(ctx).
		Where("id = ?", id).
		Delete(&dbmodel.Itinerary{}).Error
}

func marshalOptionalStringSlice(items []string) (*string, error) {
	if len(items) == 0 {
		return nil, nil
	}

	b, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}

	s := string(b)
	return &s, nil
}

func unmarshalOptionalStringSlice(raw *string) ([]string, error) {
	if raw == nil {
		return nil, nil
	}

	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return nil, nil
	}

	var items []string
	if err := json.Unmarshal([]byte(trimmed), &items); err != nil {
		return nil, err
	}

	return items, nil
}

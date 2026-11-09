package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	travellocation "platform-go/internal/domain/travel_location"
	"platform-go/internal/infrastructure/dbmodel"
)

type travelLocationRepository struct {
	db DBInterface
}

func NewTravelLocationRepository(db DBInterface) travellocation.TravelLocationRepository {
	return &travelLocationRepository{db: db}
}

func (r *travelLocationRepository) FindTravelLocationByID(ctx context.Context, id int64) (*travellocation.TravelLocation, error) {
	var model dbmodel.TravelLocation
	if err := r.db.GetDBByCtx(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &travellocation.TravelLocation{
		ID:            model.ID,
		CountryCode:   model.CountryCode,
		CountryName:   model.CountryName,
		LocationName:  model.LocationName,
		Latitude:      model.Latitude,
		Longitude:     model.Longitude,
		ImagePath:     model.ImagePath,
		Status:        model.Status,
		Description:   model.Description,
		Notes:         model.Notes,
		GoogleMapsURL: model.GoogleMapsURL,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}, nil
}

func (r *travelLocationRepository) CreateTravelLocation(ctx context.Context, loc *travellocation.TravelLocation) error {
	model := dbmodel.TravelLocation{
		CountryCode:   loc.CountryCode,
		CountryName:   loc.CountryName,
		LocationName:  loc.LocationName,
		Latitude:      loc.Latitude,
		Longitude:     loc.Longitude,
		ImagePath:     loc.ImagePath,
		Status:        loc.Status,
		Description:   loc.Description,
		Notes:         loc.Notes,
		GoogleMapsURL: loc.GoogleMapsURL,
	}
	if err := r.db.GetDBByCtx(ctx).Create(&model).Error; err != nil {
		return err
	}
	loc.ID = model.ID
	loc.CreatedAt = model.CreatedAt
	loc.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *travelLocationRepository) UpdateTravelLocation(ctx context.Context, loc *travellocation.TravelLocation) error {
	updates := map[string]interface{}{
		"country_code":    loc.CountryCode,
		"country_name":    loc.CountryName,
		"location_name":   loc.LocationName,
		"latitude":        loc.Latitude,
		"longitude":       loc.Longitude,
		"image_path":      loc.ImagePath,
		"status":          loc.Status,
		"description":     loc.Description,
		"notes":           loc.Notes,
		"google_maps_url": loc.GoogleMapsURL,
	}
	return r.db.GetDBByCtx(ctx).
		Model(&dbmodel.TravelLocation{}).
		Where("id = ?", loc.ID).
		Updates(updates).Error
}

func (r *travelLocationRepository) DeleteTravelLocationByID(ctx context.Context, id int64) error {
	return r.db.GetDBByCtx(ctx).
		Where("id = ?", id).
		Delete(&dbmodel.TravelLocation{}).Error
}


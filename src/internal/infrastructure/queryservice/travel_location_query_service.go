package queryservice

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/qsdto"
)

type TravelLocationQueryService struct {
	db *gorm.DB
}

func NewTravelLocationQueryService(db *gorm.DB) *TravelLocationQueryService {
	return &TravelLocationQueryService{db: db}
}

func (qs *TravelLocationQueryService) FetchTravelLocations(ctx context.Context, filter qsdto.TravelLocationFilter) ([]qsdto.TravelLocationDto, int64, error) {
	query := qs.db.WithContext(ctx).Model(&dbmodel.TravelLocation{})

	if filter.Status != nil && strings.TrimSpace(*filter.Status) != "" {
		query = query.Where("status = ?", strings.TrimSpace(*filter.Status))
	}
	if filter.Name != nil && strings.TrimSpace(*filter.Name) != "" {
		query = query.Where("name LIKE ?", "%"+strings.TrimSpace(*filter.Name)+"%")
	}
	if filter.StateCode != nil && strings.TrimSpace(*filter.StateCode) != "" {
		query = query.Where("state_code = ?", strings.TrimSpace(*filter.StateCode))
	}
	if filter.CountryCode != nil && strings.TrimSpace(*filter.CountryCode) != "" {
		query = query.Where("country_code = ?", strings.TrimSpace(*filter.CountryCode))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var locations []qsdto.TravelLocationDto
	if err := query.Order("id ASC").Find(&locations).Error; err != nil {
		return nil, 0, err
	}

	return locations, total, nil
}

func (qs *TravelLocationQueryService) FetchTravelLocationByID(ctx context.Context, id int64) (*qsdto.TravelLocationDto, error) {
	var location qsdto.TravelLocationDto
	err := qs.db.WithContext(ctx).Model(&dbmodel.TravelLocation{}).Where("id = ?", id).First(&location).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &location, nil
}


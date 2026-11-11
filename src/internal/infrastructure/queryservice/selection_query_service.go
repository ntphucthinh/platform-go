package queryservice

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/qsdto"
)

type SelectionQueryService struct {
	db *gorm.DB
}

func NewSelectionQueryService(db *gorm.DB) *SelectionQueryService {
	return &SelectionQueryService{db: db}
}

func (qs *SelectionQueryService) FetchTravelLocations(ctx context.Context, countryCode *string) ([]qsdto.SelectionTravelLocationDto, error) {
	query := qs.db.WithContext(ctx).
		Model(&dbmodel.TravelLocation{}).
		Select("id, location_name AS name, country_code")

	if countryCode != nil && strings.TrimSpace(*countryCode) != "" {
		query = query.Where("country_code = ?", strings.ToUpper(strings.TrimSpace(*countryCode)))
	}

	var results []qsdto.SelectionTravelLocationDto
	if err := query.Order("location_name ASC, id ASC").Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}


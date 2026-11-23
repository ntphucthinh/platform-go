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

func (qs *TravelLocationQueryService) FetchTravelLocations(ctx context.Context, filter qsdto.TravelLocationFilter) ([]qsdto.TravelLocationDto, int64, qsdto.TravelLocationStatusCountsDto, error) {
	baseCountQuery := qs.db.WithContext(ctx).Model(&dbmodel.TravelLocation{})
	if filter.CountryCode != nil && strings.TrimSpace(*filter.CountryCode) != "" {
		baseCountQuery = baseCountQuery.Where("country_code = ?", strings.ToUpper(strings.TrimSpace(*filter.CountryCode)))
	}

	type statusCountsResult struct {
		Total    int64 `gorm:"column:total"`
		Visited  int64 `gorm:"column:visited"`
		Wishlist int64 `gorm:"column:wishlist"`
	}
	var counts statusCountsResult
	err := baseCountQuery.
		Select("COUNT(*) as total, COUNT(CASE WHEN status = 'visited' THEN 1 END) as visited, COUNT(CASE WHEN status = 'wishlist' THEN 1 END) as wishlist").
		Scan(&counts).Error
	if err != nil {
		return nil, 0, qsdto.TravelLocationStatusCountsDto{}, err
	}

	dataQuery := qs.db.WithContext(ctx).Model(&dbmodel.TravelLocation{})
	if filter.CountryCode != nil && strings.TrimSpace(*filter.CountryCode) != "" {
		dataQuery = dataQuery.Where("country_code = ?", strings.ToUpper(strings.TrimSpace(*filter.CountryCode)))
	}
	if filter.Status != nil && strings.TrimSpace(*filter.Status) != "" {
		dataQuery = dataQuery.Where("status = ?", strings.TrimSpace(*filter.Status))
	}

	var locations []qsdto.TravelLocationDto
	if err := dataQuery.Order("id ASC").Find(&locations).Error; err != nil {
		return nil, 0, qsdto.TravelLocationStatusCountsDto{}, err
	}

	statusCounts := qsdto.TravelLocationStatusCountsDto{
		Visited:  counts.Visited,
		Wishlist: counts.Wishlist,
	}

	return locations, counts.Total, statusCounts, nil
}

func (qs *TravelLocationQueryService) FetchTravelLocationByID(ctx context.Context, id int64) (*qsdto.TravelLocationDetailDto, error) {
	var location qsdto.TravelLocationDetailDto
	err := qs.db.WithContext(ctx).Model(&dbmodel.TravelLocation{}).Where("id = ?", id).First(&location).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	memories := make([]qsdto.LocationMemoryDto, 0)
	err = qs.db.WithContext(ctx).
		Table("travel_memories m").
		Select("m.id, m.title, m.cover_image_path, m.start_date, m.end_date, COUNT(p.id) as photo_count").
		Joins("LEFT JOIN travel_memory_days d ON d.memory_id = m.id").
		Joins("LEFT JOIN travel_memory_photos p ON p.memory_day_id = d.id").
		Where("m.location_id = ? AND m.deleted_at IS NULL", id).
		Group("m.id, m.title, m.cover_image_path, m.start_date, m.end_date").
		Order("m.start_date DESC, m.id DESC").
		Scan(&memories).Error
	if err != nil {
		return nil, err
	}
	location.Memories = memories

	type locationItineraryRow struct {
		ID            int64  `gorm:"column:id"`
		ItineraryName string `gorm:"column:itinerary_name"`
		DaysCount     int    `gorm:"column:days_count"`
	}

	var itineraryRows []locationItineraryRow
	err = qs.db.WithContext(ctx).
		Table("itineraries i").
		Select("i.id, i.itinerary_name, COUNT(DISTINCT d.id) as days_count").
		Joins("LEFT JOIN itinerary_days d ON d.itinerary_id = i.id").
		Where("i.location_id = ? AND i.deleted_at IS NULL", id).
		Group("i.id, i.itinerary_name").
		Order("i.id DESC").
		Scan(&itineraryRows).Error
	if err != nil {
		return nil, err
	}

	itineraries := make([]qsdto.LocationItineraryDto, len(itineraryRows))
	for i, row := range itineraryRows {
		_, duration := formatDuration(row.DaysCount)
		itineraries[i] = qsdto.LocationItineraryDto{
			ID:       row.ID,
			Title:    row.ItineraryName,
			Duration: duration,
		}
	}
	location.Itineraries = itineraries
	location.VisitCount = int64(len(memories))
	location.ItineraryCount = int64(len(itineraries))

	return &location, nil
}

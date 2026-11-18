package queryservice

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/dbutil"
	"platform-go/internal/infrastructure/qsdto"
)

type ItineraryQueryService struct {
	db *gorm.DB
}

func NewItineraryQueryService(db *gorm.DB) *ItineraryQueryService {
	return &ItineraryQueryService{
		db: db,
	}
}

const itinerarySelectColumns = `
	itineraries.id,
	itineraries.location_id,
	travel_locations.location_name,
	itineraries.itinerary_name,
	itineraries.description,
	itineraries.cover_image_path
`

const itineraryLocationJoin = "LEFT JOIN travel_locations ON travel_locations.id = itineraries.location_id"

func (qs *ItineraryQueryService) FetchItineraries(ctx context.Context, page, perPage int, locationID *int64) (*dbutil.Paginator[qsdto.ItineraryDto], error) {
	query := qs.db.WithContext(ctx).
		Model(&dbmodel.Itinerary{}).
		Select(itinerarySelectColumns).
		Joins(itineraryLocationJoin).
		Where("itineraries.deleted_at IS NULL")

	if locationID != nil {
		query = query.Where("itineraries.location_id = ?", *locationID)
	}

	// Default fixed sort
	query = query.Order("itineraries.id DESC")

	paginated, err := dbutil.ModelPaginate[qsdto.ItineraryDto](query, perPage, page)
	if err != nil {
		return nil, err
	}
	if len(paginated.Data) == 0 {
		return paginated, nil
	}

	if err := qs.populateDurationAndActivities(ctx, paginated.Data); err != nil {
		return nil, err
	}

	return paginated, nil
}

type itineraryStats struct {
	ItineraryID     int64 `gorm:"column:itinerary_id"`
	DaysCount       int   `gorm:"column:days_count"`
	ActivitiesCount int   `gorm:"column:activities_count"`
}

func (qs *ItineraryQueryService) populateDurationAndActivities(ctx context.Context, itineraries []qsdto.ItineraryDto) error {
	if len(itineraries) == 0 {
		return nil
	}

	itineraryIDs := make([]int64, 0, len(itineraries))
	for _, it := range itineraries {
		itineraryIDs = append(itineraryIDs, it.ID)
	}

	var stats []itineraryStats
	err := qs.db.WithContext(ctx).
		Table("itinerary_days d").
		Select("d.itinerary_id, COUNT(DISTINCT d.id) as days_count, COUNT(a.id) as activities_count").
		Joins("LEFT JOIN itinerary_activities a ON a.itinerary_day_id = d.id").
		Where("d.itinerary_id IN ?", itineraryIDs).
		Group("d.itinerary_id").
		Scan(&stats).Error
	if err != nil {
		return err
	}

	statsMap := make(map[int64]itineraryStats, len(stats))
	for _, s := range stats {
		statsMap[s.ItineraryID] = s
	}

	for i := range itineraries {
		s := statsMap[itineraries[i].ID]

		days := s.DaysCount
		activities := s.ActivitiesCount

		nights, duration := formatDuration(days)

		itineraries[i].DurationDays = days
		itineraries[i].DurationNights = nights
		itineraries[i].Duration = duration
		itineraries[i].ActivityCount = activities
	}

	return nil
}

func formatDuration(days int) (int, string) {
	nights := 0

	if days > 1 {
		nights = days - 1
	}

	duration := fmt.Sprintf("%d ngày", days)

	if nights > 0 {
		duration = fmt.Sprintf("%d ngày %d đêm", days, nights)
	}

	return nights, duration
}

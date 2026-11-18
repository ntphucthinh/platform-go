package queryservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

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

const itineraryDetailSelectColumns = `
	itineraries.id,
	itineraries.itinerary_name,
	itineraries.location_id,
	travel_locations.location_name,
	itineraries.suitable_for,
	itineraries.description,
	itineraries.cover_image_path,
	itineraries.travel_tips
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

func (qs *ItineraryQueryService) FetchItineraryByID(ctx context.Context, id int64) (*qsdto.ItineraryDetailDto, error) {
	var detail qsdto.ItineraryDetailDto
	err := qs.db.WithContext(ctx).
		Model(&dbmodel.Itinerary{}).
		Select(itineraryDetailSelectColumns).
		Joins(itineraryLocationJoin).
		Where("itineraries.id = ? AND itineraries.deleted_at IS NULL", id).
		First(&detail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	detail.TravelTips = parseJSONStringSlice(detail.TravelTipsRaw)

	// Fetch themes
	themes := make([]qsdto.ItineraryThemeDto, 0)
	if err := qs.db.WithContext(ctx).
		Table("itinerary_themes it").
		Select("tt.id, tt.name").
		Joins("INNER JOIN travel_themes tt ON tt.id = it.theme_id AND tt.deleted_at IS NULL").
		Where("it.itinerary_id = ?", id).
		Order("it.id ASC").
		Scan(&themes).Error; err != nil {
		return nil, err
	}

	detail.Themes = themes

	// Fetch featured images
	featuredImages := make([]qsdto.ItineraryFeaturedImageDto, 0)
	if err := qs.db.WithContext(ctx).
		Table("itinerary_featured_images").
		Select("id, itinerary_id, image_path").
		Where("itinerary_id = ?", id).
		Order("id ASC").
		Find(&featuredImages).Error; err != nil {
		return nil, err
	}

	detail.FeaturedImages = featuredImages

	// Fetch days
	days := make([]qsdto.ItineraryDayDto, 0)
	if err := qs.db.WithContext(ctx).
		Table("itinerary_days").
		Select("id, itinerary_id, day_number, title, subtitle, summary, image_path, day_tips").
		Where("itinerary_id = ?", id).
		Order("day_number ASC, id ASC").
		Find(&days).Error; err != nil {
		return nil, err
	}

	dayIDs := make([]int64, 0, len(days))
	for _, day := range days {
		dayIDs = append(dayIDs, day.ID)
	}

	activitiesByDayID := make(map[int64][]qsdto.ItineraryActivityDto, len(days))
	totalActivities := 0

	if len(dayIDs) > 0 {
		var activities []qsdto.ItineraryActivityDto
		if err := qs.db.WithContext(ctx).
			Table("itinerary_activities").
			Select("id, itinerary_day_id, start_time, title, activity_category, duration_minutes, location, google_maps_url, description, notes, image_path").
			Where("itinerary_day_id IN ?", dayIDs).
			Order("start_time ASC, id ASC").
			Find(&activities).Error; err != nil {
			return nil, err
		}

		totalActivities = len(activities)
		for _, act := range activities {
			act.StartTime = formatStartTime(act.StartTime)
			activitiesByDayID[act.ItineraryDayID] = append(activitiesByDayID[act.ItineraryDayID], act)
		}
	}

	for i := range days {
		days[i].DayTips = parseJSONStringSlice(days[i].DayTipsRaw)
		acts := activitiesByDayID[days[i].ID]
		if acts == nil {
			acts = []qsdto.ItineraryActivityDto{}
		}
		days[i].Activities = acts
	}

	detail.Days = days

	daysCount := len(days)
	nights, duration := formatDuration(daysCount)
	detail.DurationDays = daysCount
	detail.DurationNights = nights
	detail.Duration = duration
	detail.ActivityCount = totalActivities

	return &detail, nil
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

func parseJSONStringSlice(raw *string) []string {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return []string{}
	}
	var tips []string
	if err := json.Unmarshal([]byte(*raw), &tips); err != nil || tips == nil {
		return []string{}
	}
	return tips
}

func formatStartTime(raw string) string {
	trimmed := strings.TrimSpace(raw)

	parsed, err := time.Parse("15:04:05", trimmed)
	if err == nil {
		return parsed.Format("15:04")
	}

	parsed, err = time.Parse("15:04", trimmed)
	if err == nil {
		return parsed.Format("15:04")
	}

	return trimmed
}

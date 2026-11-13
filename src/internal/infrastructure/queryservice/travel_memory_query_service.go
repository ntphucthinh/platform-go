package queryservice

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/dbutil"
	"platform-go/internal/infrastructure/qsdto"
)

type TravelMemoryQueryService struct {
	db *gorm.DB
}

func NewTravelMemoryQueryService(db *gorm.DB) *TravelMemoryQueryService {
	return &TravelMemoryQueryService{
		db: db,
	}
}

const travelMemorySelectColumns = `
	travel_memories.id,
	travel_memories.location_id,
	travel_locations.location_name,
	travel_locations.country_code,
	travel_locations.country_name,
	travel_memories.title,
	travel_memories.start_date,
	travel_memories.end_date,
	travel_memories.description,
	travel_memories.cover_image_path
`

const travelMemoryLocationJoin = "LEFT JOIN travel_locations ON travel_locations.id = travel_memories.location_id"

func (qs *TravelMemoryQueryService) FetchTravelMemories(ctx context.Context, page, perPage int, countryCode *string, locationID *int64, year *int) (*dbutil.Paginator[qsdto.TravelMemoryDto], error) {
	query := qs.db.WithContext(ctx).
		Model(&dbmodel.TravelMemory{}).
		Select(travelMemorySelectColumns).
		Joins(travelMemoryLocationJoin)

	if countryCode != nil && strings.TrimSpace(*countryCode) != "" {
		query = query.Where("travel_locations.country_code = ?", strings.ToUpper(strings.TrimSpace(*countryCode)))
	}

	if locationID != nil {
		query = query.Where("travel_memories.location_id = ?", *locationID)
	}

	if year != nil {
		startOfYear := fmt.Sprintf("%04d-01-01", *year)
		endOfYear := fmt.Sprintf("%04d-12-31", *year)
		query = query.Where("travel_memories.start_date <= ? AND travel_memories.end_date >= ?", endOfYear, startOfYear)
	}

	query = query.Order("travel_memories.start_date DESC, travel_memories.id DESC")

	paginated, err := dbutil.ModelPaginate[qsdto.TravelMemoryDto](query, perPage, page)
	if err != nil {
		return nil, err
	}
	if len(paginated.Data) == 0 {
		return paginated, nil
	}

	if err := qs.populateDaysAndPhotos(ctx, paginated.Data); err != nil {
		return nil, err
	}

	return paginated, nil
}

func (qs *TravelMemoryQueryService) FetchTravelMemoryByID(ctx context.Context, id int64) (*qsdto.TravelMemoryDto, error) {
	var memory qsdto.TravelMemoryDto
	err := qs.db.WithContext(ctx).
		Model(&dbmodel.TravelMemory{}).
		Select(travelMemorySelectColumns).
		Joins(travelMemoryLocationJoin).
		Where("travel_memories.id = ?", id).
		First(&memory).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	memories := []qsdto.TravelMemoryDto{memory}
	if err := qs.populateDaysAndPhotos(ctx, memories); err != nil {
		return nil, err
	}

	return &memories[0], nil
}

func (qs *TravelMemoryQueryService) populateDaysAndPhotos(ctx context.Context, memories []qsdto.TravelMemoryDto) error {
	if len(memories) == 0 {
		return nil
	}

	memoryIDs := make([]int64, 0, len(memories))
	for _, memory := range memories {
		memoryIDs = append(memoryIDs, memory.ID)
	}

	// Fetch days
	var days []qsdto.TravelMemoryDayDto
	if err := qs.db.WithContext(ctx).
		Table("travel_memory_days").
		Where("memory_id IN ?", memoryIDs).
		Order("date ASC, id ASC").
		Find(&days).Error; err != nil {
		return err
	}

	// Fetch photos
	dayIDs := make([]int64, 0, len(days))
	for _, day := range days {
		dayIDs = append(dayIDs, day.ID)
	}

	photosByDayID := make(map[int64][]qsdto.TravelMemoryPhotoDto)

	if len(dayIDs) > 0 {
		var photos []qsdto.TravelMemoryPhotoDto

		if err := qs.db.WithContext(ctx).
			Table("travel_memory_photos").
			Where("memory_day_id IN ?", dayIDs).
			Order("sort_order ASC, id ASC").
			Find(&photos).Error; err != nil {
			return err
		}

		for _, photo := range photos {
			photosByDayID[photo.MemoryDayID] = append(
				photosByDayID[photo.MemoryDayID],
				photo,
			)
		}
	}

	// Group days by memory
	daysByMemoryID := make(map[int64][]qsdto.TravelMemoryDayDto)

	for i := range days {
		days[i].Photos = photosByDayID[days[i].ID]

		daysByMemoryID[days[i].MemoryID] = append(
			daysByMemoryID[days[i].MemoryID],
			days[i],
		)
	}

	// Assign days to memories
	for i := range memories {
		memories[i].Days = daysByMemoryID[memories[i].ID]
	}

	return nil
}

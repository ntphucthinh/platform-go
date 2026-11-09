package factory

import (
	"time"

	"platform-go/internal/infrastructure/dbmodel"
)

// TravelMemoryFactory generates TravelMemory dbmodel instances for seeders and tests.
type TravelMemoryFactory struct {
	locationID     int64
	title          string
	startDate      time.Time
	endDate        time.Time
	description    *string
	coverImagePath string
}

func NewTravelMemoryFactory() *TravelMemoryFactory {
	desc := "Hành trình khám phá tuyệt vời"
	start := time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 11, 3, 0, 0, 0, 0, time.UTC)

	return &TravelMemoryFactory{
		locationID:     1,
		title:          "Chuyến đi đáng nhớ",
		startDate:      start,
		endDate:        end,
		description:    &desc,
		coverImagePath: "/images/memories/default_cover.jpg",
	}
}

func (f *TravelMemoryFactory) WithLocationID(id int64) *TravelMemoryFactory {
	f.locationID = id
	return f
}

func (f *TravelMemoryFactory) WithTitle(title string) *TravelMemoryFactory {
	f.title = title
	return f
}

func (f *TravelMemoryFactory) WithDates(start, end time.Time) *TravelMemoryFactory {
	f.startDate = start
	f.endDate = end
	return f
}

func (f *TravelMemoryFactory) WithDescription(desc string) *TravelMemoryFactory {
	f.description = &desc
	return f
}

func (f *TravelMemoryFactory) WithCoverImagePath(path string) *TravelMemoryFactory {
	f.coverImagePath = path
	return f
}

func (f *TravelMemoryFactory) Build() *dbmodel.TravelMemory {
	return &dbmodel.TravelMemory{
		LocationID:     f.locationID,
		Title:          f.title,
		StartDate:      f.startDate,
		EndDate:        f.endDate,
		Description:    f.description,
		CoverImagePath: f.coverImagePath,
	}
}

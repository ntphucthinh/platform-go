package factory

import (
	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
)

// ItineraryDayFactory generates ItineraryDay dbmodel instances for seeders and tests.
type ItineraryDayFactory struct {
	itineraryID int64
	dayNumber   int
	title       string
	subtitle    *string
	summary     *string
	imagePath   *string
	dayTips     *string
}

func NewItineraryDayFactory() *ItineraryDayFactory {
	title := "Ngày 1: Khám phá điểm đến"
	subtitle := "Bắt đầu hành trình mới"
	summary := "Một ngày trải nghiệm đầy thú vị với nhiều hoạt động đa dạng"
	imagePath := "itineraries/days/default-day-1.jpg"
	tips := `["Nên dậy sớm", "Chuẩn bị giày đi bộ"]`

	return &ItineraryDayFactory{
		itineraryID: 1,
		dayNumber:   1,
		title:       title,
		subtitle:    &subtitle,
		summary:     &summary,
		imagePath:   &imagePath,
		dayTips:     &tips,
	}
}

func (f *ItineraryDayFactory) WithItineraryID(id int64) *ItineraryDayFactory {
	f.itineraryID = id
	return f
}

func (f *ItineraryDayFactory) WithDayNumber(num int) *ItineraryDayFactory {
	f.dayNumber = num
	return f
}

func (f *ItineraryDayFactory) WithTitle(title string) *ItineraryDayFactory {
	f.title = title
	return f
}

func (f *ItineraryDayFactory) WithSubtitle(sub string) *ItineraryDayFactory {
	f.subtitle = &sub
	return f
}

func (f *ItineraryDayFactory) WithSummary(summary string) *ItineraryDayFactory {
	f.summary = &summary
	return f
}

func (f *ItineraryDayFactory) WithImagePath(path string) *ItineraryDayFactory {
	f.imagePath = &path
	return f
}

func (f *ItineraryDayFactory) WithDayTips(tips string) *ItineraryDayFactory {
	f.dayTips = &tips
	return f
}

func (f *ItineraryDayFactory) Build() *dbmodel.ItineraryDay {
	return &dbmodel.ItineraryDay{
		ItineraryID: f.itineraryID,
		DayNumber:   f.dayNumber,
		Title:       f.title,
		Subtitle:    f.subtitle,
		Summary:     f.summary,
		ImagePath:   f.imagePath,
		DayTips:     f.dayTips,
	}
}

func (f *ItineraryDayFactory) Insert(db *gorm.DB) (*dbmodel.ItineraryDay, error) {
	item := f.Build()
	if err := db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (f *ItineraryDayFactory) MustInsert(db *gorm.DB) *dbmodel.ItineraryDay {
	item, err := f.Insert(db)
	if err != nil {
		panic(err)
	}
	return item
}

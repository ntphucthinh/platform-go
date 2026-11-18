package factory

import (
	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
)

// ItineraryActivityFactory generates ItineraryActivity dbmodel instances for seeders and tests.
type ItineraryActivityFactory struct {
	itineraryDayID   int64
	startTime        string
	title            string
	activityCategory *string
	durationMinutes  *uint16
	location         *string
	googleMapsURL    *string
	description      *string
	notes            *string
	imagePath        *string
}

func NewItineraryActivityFactory() *ItineraryActivityFactory {
	category := "sightseeing"
	duration := uint16(60)
	location := "Địa điểm tham quan"
	gmaps := "https://maps.google.com"
	desc := "Hoạt động trải nghiệm thú vị"
	notes := "Ghi chú hoạt động"
	image := "itineraries/activities/default.jpg"

	return &ItineraryActivityFactory{
		itineraryDayID:   1,
		startTime:        "09:00:00",
		title:            "Tham quan danh lam thắng cảnh",
		activityCategory: &category,
		durationMinutes:  &duration,
		location:         &location,
		googleMapsURL:    &gmaps,
		description:      &desc,
		notes:            &notes,
		imagePath:        &image,
	}
}

func (f *ItineraryActivityFactory) WithItineraryDayID(id int64) *ItineraryActivityFactory {
	f.itineraryDayID = id
	return f
}

func (f *ItineraryActivityFactory) WithStartTime(timeStr string) *ItineraryActivityFactory {
	f.startTime = timeStr
	return f
}

func (f *ItineraryActivityFactory) WithTitle(title string) *ItineraryActivityFactory {
	f.title = title
	return f
}

func (f *ItineraryActivityFactory) WithCategory(cat string) *ItineraryActivityFactory {
	f.activityCategory = &cat
	return f
}

func (f *ItineraryActivityFactory) WithDurationMinutes(mins uint16) *ItineraryActivityFactory {
	f.durationMinutes = &mins
	return f
}

func (f *ItineraryActivityFactory) WithLocation(loc string) *ItineraryActivityFactory {
	f.location = &loc
	return f
}

func (f *ItineraryActivityFactory) WithGoogleMapsURL(url string) *ItineraryActivityFactory {
	f.googleMapsURL = &url
	return f
}

func (f *ItineraryActivityFactory) WithDescription(desc string) *ItineraryActivityFactory {
	f.description = &desc
	return f
}

func (f *ItineraryActivityFactory) WithNotes(notes string) *ItineraryActivityFactory {
	f.notes = &notes
	return f
}

func (f *ItineraryActivityFactory) WithImagePath(path string) *ItineraryActivityFactory {
	f.imagePath = &path
	return f
}

func (f *ItineraryActivityFactory) Build() *dbmodel.ItineraryActivity {
	return &dbmodel.ItineraryActivity{
		ItineraryDayID:   f.itineraryDayID,
		StartTime:        f.startTime,
		Title:            f.title,
		ActivityCategory: f.activityCategory,
		DurationMinutes:  f.durationMinutes,
		Location:         f.location,
		GoogleMapsURL:    f.googleMapsURL,
		Description:      f.description,
		Notes:            f.notes,
		ImagePath:        f.imagePath,
	}
}

func (f *ItineraryActivityFactory) Insert(db *gorm.DB) (*dbmodel.ItineraryActivity, error) {
	item := f.Build()
	if err := db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (f *ItineraryActivityFactory) MustInsert(db *gorm.DB) *dbmodel.ItineraryActivity {
	item, err := f.Insert(db)
	if err != nil {
		panic(err)
	}
	return item
}

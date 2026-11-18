package factory

import (
	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
)

// ItineraryFactory generates Itinerary dbmodel instances for seeders and tests.
type ItineraryFactory struct {
	itineraryName  string
	locationID     *int64
	suitableFor    *string
	description    *string
	coverImagePath *string
	travelTips     *string
}

func NewItineraryFactory() *ItineraryFactory {
	name := "Hành trình khám phá"
	suitable := "Mọi đối tượng"
	desc := "Lịch trình du lịch chi tiết và hấp dẫn"
	cover := "itineraries/default-cover.jpg"
	tips := `["Chuẩn bị giấy tờ tùy thân", "Mang trang phục phù hợp"]`

	return &ItineraryFactory{
		itineraryName:  name,
		suitableFor:    &suitable,
		description:    &desc,
		coverImagePath: &cover,
		travelTips:     &tips,
	}
}

func (f *ItineraryFactory) WithName(name string) *ItineraryFactory {
	f.itineraryName = name
	return f
}

func (f *ItineraryFactory) WithLocationID(id *int64) *ItineraryFactory {
	f.locationID = id
	return f
}

func (f *ItineraryFactory) WithSuitableFor(suitable string) *ItineraryFactory {
	f.suitableFor = &suitable
	return f
}

func (f *ItineraryFactory) WithDescription(desc string) *ItineraryFactory {
	f.description = &desc
	return f
}

func (f *ItineraryFactory) WithCoverImagePath(path string) *ItineraryFactory {
	f.coverImagePath = &path
	return f
}

func (f *ItineraryFactory) WithTravelTips(tips string) *ItineraryFactory {
	f.travelTips = &tips
	return f
}

func (f *ItineraryFactory) Build() *dbmodel.Itinerary {
	return &dbmodel.Itinerary{
		ItineraryName:  f.itineraryName,
		LocationID:     f.locationID,
		SuitableFor:    f.suitableFor,
		Description:    f.description,
		CoverImagePath: f.coverImagePath,
		TravelTips:     f.travelTips,
	}
}

func (f *ItineraryFactory) Insert(db *gorm.DB) (*dbmodel.Itinerary, error) {
	item := f.Build()
	if err := db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (f *ItineraryFactory) MustInsert(db *gorm.DB) *dbmodel.Itinerary {
	item, err := f.Insert(db)
	if err != nil {
		panic(err)
	}
	return item
}

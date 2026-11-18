package factory

import (
	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
)

// ItineraryFeaturedImageFactory generates ItineraryFeaturedImage dbmodel instances for seeders and tests.
type ItineraryFeaturedImageFactory struct {
	itineraryID int64
	imagePath   string
}

func NewItineraryFeaturedImageFactory() *ItineraryFeaturedImageFactory {
	return &ItineraryFeaturedImageFactory{
		itineraryID: 1,
		imagePath:   "itineraries/featured/default.jpg",
	}
}

func (f *ItineraryFeaturedImageFactory) WithItineraryID(id int64) *ItineraryFeaturedImageFactory {
	f.itineraryID = id
	return f
}

func (f *ItineraryFeaturedImageFactory) WithImagePath(path string) *ItineraryFeaturedImageFactory {
	f.imagePath = path
	return f
}

func (f *ItineraryFeaturedImageFactory) Build() *dbmodel.ItineraryFeaturedImage {
	return &dbmodel.ItineraryFeaturedImage{
		ItineraryID: f.itineraryID,
		ImagePath:   f.imagePath,
	}
}

func (f *ItineraryFeaturedImageFactory) Insert(db *gorm.DB) (*dbmodel.ItineraryFeaturedImage, error) {
	item := f.Build()
	if err := db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (f *ItineraryFeaturedImageFactory) MustInsert(db *gorm.DB) *dbmodel.ItineraryFeaturedImage {
	item, err := f.Insert(db)
	if err != nil {
		panic(err)
	}
	return item
}

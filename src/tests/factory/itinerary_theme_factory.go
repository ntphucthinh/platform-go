package factory

import (
	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
)

// ItineraryThemeFactory generates ItineraryTheme dbmodel instances for seeders and tests.
type ItineraryThemeFactory struct {
	itineraryID int64
	themeID     int64
}

func NewItineraryThemeFactory() *ItineraryThemeFactory {
	return &ItineraryThemeFactory{
		itineraryID: 1,
		themeID:     1,
	}
}

func (f *ItineraryThemeFactory) WithItineraryID(id int64) *ItineraryThemeFactory {
	f.itineraryID = id
	return f
}

func (f *ItineraryThemeFactory) WithThemeID(id int64) *ItineraryThemeFactory {
	f.themeID = id
	return f
}

func (f *ItineraryThemeFactory) Build() *dbmodel.ItineraryTheme {
	return &dbmodel.ItineraryTheme{
		ItineraryID: f.itineraryID,
		ThemeID:     f.themeID,
	}
}

func (f *ItineraryThemeFactory) Insert(db *gorm.DB) (*dbmodel.ItineraryTheme, error) {
	item := f.Build()
	if err := db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (f *ItineraryThemeFactory) MustInsert(db *gorm.DB) *dbmodel.ItineraryTheme {
	item, err := f.Insert(db)
	if err != nil {
		panic(err)
	}
	return item
}

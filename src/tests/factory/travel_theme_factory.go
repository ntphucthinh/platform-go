package factory

import (
	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
)

// TravelThemeFactory generates TravelTheme dbmodel instances for tests and seeders.
type TravelThemeFactory struct {
	name string
}

func NewTravelThemeFactory() *TravelThemeFactory {
	return &TravelThemeFactory{
		name: "Thiên nhiên",
	}
}

func (f *TravelThemeFactory) WithName(name string) *TravelThemeFactory {
	f.name = name
	return f
}

func (f *TravelThemeFactory) Build() *dbmodel.TravelTheme {
	return &dbmodel.TravelTheme{
		Name: f.name,
	}
}

func (f *TravelThemeFactory) Insert(db *gorm.DB) (*dbmodel.TravelTheme, error) {
	theme := f.Build()
	if err := db.Create(theme).Error; err != nil {
		return nil, err
	}
	return theme, nil
}

func (f *TravelThemeFactory) MustInsert(db *gorm.DB) *dbmodel.TravelTheme {
	theme, err := f.Insert(db)
	if err != nil {
		panic(err)
	}
	return theme
}

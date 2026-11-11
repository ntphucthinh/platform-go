package qsdto

type SelectionTravelLocationDto struct {
	ID          int64  `gorm:"column:id"`
	Name        string `gorm:"column:name"`
	CountryCode string `gorm:"column:country_code"`
}

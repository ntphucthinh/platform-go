package qsdto

type TravelLocationDto struct {
	ID             int64    `gorm:"column:id"`
	CountryCode    string   `gorm:"column:country_code"`
	CountryName    string   `gorm:"column:country_name"`
	StateCode      string   `gorm:"column:state_code"`
	Name           string   `gorm:"column:name"`
	Address        string   `gorm:"column:address"`
	Latitude       *float64 `gorm:"column:latitude"`
	Longitude      *float64 `gorm:"column:longitude"`
	ImagePath      *string  `gorm:"column:image_path"`
	Status         string   `gorm:"column:status"`
	Description    *string  `gorm:"column:description"`
	Notes          *string  `gorm:"column:notes"`
	GoogleMapsURL  *string  `gorm:"column:google_maps_url"`
	VisitCount     int64    `gorm:"-"`
	ItineraryCount int64    `gorm:"-"`
}

type TravelLocationFilter struct {
	Status      *string
	Name        *string
	StateCode   *string
	CountryCode *string
}


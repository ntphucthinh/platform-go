package qsdto

import "time"

type TravelLocationDto struct {
	ID             int64    `gorm:"column:id"`
	CountryCode    string   `gorm:"column:country_code"`
	CountryName    string   `gorm:"column:country_name"`
	LocationName   string   `gorm:"column:location_name"`
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

type TravelLocationDetailDto struct {
	ID             int64                  `gorm:"column:id"`
	CountryCode    string                 `gorm:"column:country_code"`
	CountryName    string                 `gorm:"column:country_name"`
	LocationName   string                 `gorm:"column:location_name"`
	Latitude       *float64               `gorm:"column:latitude"`
	Longitude      *float64               `gorm:"column:longitude"`
	ImagePath      *string                `gorm:"column:image_path"`
	Status         string                 `gorm:"column:status"`
	Description    *string                `gorm:"column:description"`
	Notes          *string                `gorm:"column:notes"`
	GoogleMapsURL  *string                `gorm:"column:google_maps_url"`
	VisitCount     int64                  `gorm:"-"`
	ItineraryCount int64                  `gorm:"-"`
	Memories       []LocationMemoryDto    `gorm:"-"`
	Itineraries    []LocationItineraryDto `gorm:"-"`
}

type LocationMemoryDto struct {
	ID             int64     `gorm:"column:id"`
	Title          string    `gorm:"column:title"`
	CoverImagePath string    `gorm:"column:cover_image_path"`
	StartDate      time.Time `gorm:"column:start_date"`
	EndDate        time.Time `gorm:"column:end_date"`
	PhotoCount     int       `gorm:"column:photo_count"`
}

type LocationItineraryDto struct {
	ID       int64  `gorm:"column:id"`
	Title    string `gorm:"column:title"`
	Duration string `gorm:"column:duration"`
}

type TravelLocationStatusCountsDto struct {
	Visited  int64
	Wishlist int64
}

type TravelLocationFilter struct {
	Status      *string
	CountryCode *string
}

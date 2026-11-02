package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

// TravelLocation represents the travel_locations table in the database
type TravelLocation struct {
	ID            int64          `gorm:"primaryKey;autoIncrement"`
	CountryCode   string         `gorm:"column:country_code;size:2;not null"`
	CountryName   string         `gorm:"column:country_name;size:255;not null"`
	StateCode     string         `gorm:"column:state_code;size:20;not null"`
	Name          string         `gorm:"size:255;not null"`
	Address       string         `gorm:"size:255;not null"`
	Latitude      *float64       `gorm:"column:latitude;type:decimal(10,7)"`
	Longitude     *float64       `gorm:"column:longitude;type:decimal(10,7)"`
	ImagePath     *string        `gorm:"column:image_path;size:255"`
	Status        string         `gorm:"column:status;type:enum('visited','wishlist');default:visited;not null"`
	Description   *string        `gorm:"type:text"`
	Notes         *string        `gorm:"type:text"`
	GoogleMapsURL *string        `gorm:"column:google_maps_url;size:255"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP;autoUpdateLine"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// TableName returns the table name for TravelLocation
func (TravelLocation) TableName() string {
	return "travel_locations"
}


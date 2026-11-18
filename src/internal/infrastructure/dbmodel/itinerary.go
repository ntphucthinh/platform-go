package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

// Itinerary represents the itineraries table
type Itinerary struct {
	ID             int64          `gorm:"primaryKey;autoIncrement"`
	ItineraryName  string         `gorm:"column:itinerary_name;size:255;not null"`
	LocationID     *int64         `gorm:"column:location_id;index"`
	SuitableFor    *string        `gorm:"column:suitable_for;size:255"`
	Description    *string        `gorm:"column:description;size:500"`
	CoverImagePath *string        `gorm:"column:cover_image_path;size:255"`
	TravelTips     *string        `gorm:"column:travel_tips;type:json"`
	CreatedAt      time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	Location       *TravelLocation          `gorm:"foreignKey:LocationID"`
	Themes         []ItineraryTheme         `gorm:"foreignKey:ItineraryID"`
	Days           []ItineraryDay           `gorm:"foreignKey:ItineraryID"`
	FeaturedImages []ItineraryFeaturedImage `gorm:"foreignKey:ItineraryID"`
}

func (Itinerary) TableName() string {
	return "itineraries"
}

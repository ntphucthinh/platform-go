package dbmodel

import (
	"time"
)

// ItineraryFeaturedImage represents the itinerary_featured_images table
type ItineraryFeaturedImage struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	ItineraryID int64     `gorm:"column:itinerary_id;not null;index"`
	ImagePath   string    `gorm:"column:image_path;size:500;not null"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (ItineraryFeaturedImage) TableName() string {
	return "itinerary_featured_images"
}

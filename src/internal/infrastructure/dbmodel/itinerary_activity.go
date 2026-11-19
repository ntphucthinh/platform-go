package dbmodel

import (
	"time"
)

// ItineraryActivity represents the itinerary_activities table
type ItineraryActivity struct {
	ID               int64     `gorm:"primaryKey;autoIncrement"`
	ItineraryDayID   int64     `gorm:"column:itinerary_day_id;not null;index"`
	StartTime        string    `gorm:"column:start_time;type:time;not null"`
	Title            string    `gorm:"column:title;size:255;not null"`
	ActivityCategory *string   `gorm:"column:activity_category;type:enum('dining','sightseeing','nature','cafe','shopping','rest','other')"`
	DurationMinutes  *uint16   `gorm:"column:duration_minutes"`
	Location         *string   `gorm:"column:location;size:255"`
	GoogleMapsURL    *string   `gorm:"column:google_maps_url;size:500"`
	Description      *string   `gorm:"column:description;size:500"`
	Notes            *string   `gorm:"column:notes;size:255"`
	ImagePath        *string   `gorm:"column:image_path;size:255"`
	CreatedAt        time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (ItineraryActivity) TableName() string {
	return "itinerary_activities"
}

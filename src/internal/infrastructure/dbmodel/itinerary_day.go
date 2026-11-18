package dbmodel

import (
	"time"
)

// ItineraryDay represents the itinerary_days table
type ItineraryDay struct {
	ID          int64               `gorm:"primaryKey;autoIncrement"`
	ItineraryID int64               `gorm:"column:itinerary_id;not null;uniqueIndex:uq_itinerary_day_number,priority:1"`
	DayNumber   int                 `gorm:"column:day_number;not null;uniqueIndex:uq_itinerary_day_number,priority:2"`
	Title       string              `gorm:"column:title;size:255;not null"`
	Subtitle    *string             `gorm:"column:subtitle;size:255"`
	Summary     *string             `gorm:"column:summary;size:500"`
	ImagePath   *string             `gorm:"column:image_path;size:255"`
	DayTips     *string             `gorm:"column:day_tips;type:json"`
	CreatedAt   time.Time           `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time           `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`

	Activities  []ItineraryActivity `gorm:"foreignKey:ItineraryDayID"`
}

func (ItineraryDay) TableName() string {
	return "itinerary_days"
}

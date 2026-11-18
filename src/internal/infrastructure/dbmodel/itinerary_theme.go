package dbmodel

import (
	"time"
)

// ItineraryTheme represents the itinerary_themes table
type ItineraryTheme struct {
	ID          int64        `gorm:"primaryKey;autoIncrement"`
	ItineraryID int64        `gorm:"column:itinerary_id;not null;uniqueIndex:uq_itinerary_themes,priority:1"`
	ThemeID     int64        `gorm:"column:theme_id;not null;uniqueIndex:uq_itinerary_themes,priority:2;index"`
	CreatedAt   time.Time    `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`

	Itinerary   *Itinerary   `gorm:"foreignKey:ItineraryID"`
	Theme       *TravelTheme `gorm:"foreignKey:ThemeID"`
}

func (ItineraryTheme) TableName() string {
	return "itinerary_themes"
}

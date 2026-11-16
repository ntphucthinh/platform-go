package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

// TravelTheme represents the travel_themes table in the database
type TravelTheme struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"column:name;size:100;not null"`
	CreatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName returns the table name for TravelTheme
func (TravelTheme) TableName() string {
	return "travel_themes"
}

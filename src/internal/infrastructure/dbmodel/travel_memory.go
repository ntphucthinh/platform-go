package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

// TravelMemory represents the travel_memories table
type TravelMemory struct {
	ID             int64          `gorm:"primaryKey;autoIncrement"`
	LocationID     int64          `gorm:"column:location_id;not null;index"`
	Title          string         `gorm:"column:title;size:255;not null"`
	StartDate      time.Time      `gorm:"column:start_date;type:date;not null"`
	EndDate        time.Time      `gorm:"column:end_date;type:date;not null"`
	Description    *string        `gorm:"column:description;type:text"`
	CoverImagePath string         `gorm:"column:cover_image_path;size:255;not null"`
	CreatedAt      time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	Days []TravelMemoryDay `gorm:"foreignKey:MemoryID"`
}

func (TravelMemory) TableName() string {
	return "travel_memories"
}

// TravelMemoryDay represents the travel_memory_days table
type TravelMemoryDay struct {
	ID          int64               `gorm:"primaryKey;autoIncrement"`
	MemoryID    int64               `gorm:"column:memory_id;not null;index"`
	Date        time.Time           `gorm:"column:date;type:date;not null"`
	Description *string             `gorm:"column:description;size:500"`
	CreatedAt   time.Time           `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time           `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	Photos      []TravelMemoryPhoto `gorm:"foreignKey:MemoryDayID"`
}

func (TravelMemoryDay) TableName() string {
	return "travel_memory_days"
}

// TravelMemoryPhoto represents the travel_memory_photos table
type TravelMemoryPhoto struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	MemoryDayID int64     `gorm:"column:memory_day_id;not null;index"`
	ImagePath   string    `gorm:"column:image_path;size:255;not null"`
	SortOrder   int32     `gorm:"column:sort_order;default:0"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (TravelMemoryPhoto) TableName() string {
	return "travel_memory_photos"
}


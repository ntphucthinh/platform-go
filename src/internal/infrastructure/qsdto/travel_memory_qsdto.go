package qsdto

import "time"

type TravelMemoryDto struct {
	ID             int64     `gorm:"column:id"`
	LocationID     int64     `gorm:"column:location_id"`
	LocationName   string    `gorm:"column:location_name"`
	CountryCode    string    `gorm:"column:country_code"`
	CountryName    string    `gorm:"column:country_name"`
	Title          string    `gorm:"column:title"`
	StartDate      time.Time `gorm:"column:start_date"`
	EndDate        time.Time `gorm:"column:end_date"`
	Description    *string   `gorm:"column:description"`
	CoverImagePath string    `gorm:"column:cover_image_path"`

	Days []TravelMemoryDayDto `gorm:"-"`
}

type TravelMemoryDayDto struct {
	ID          int64     `gorm:"column:id"`
	MemoryID    int64     `gorm:"column:memory_id"`
	Date        time.Time `gorm:"column:date"`
	Description *string   `gorm:"column:description"`

	Photos []TravelMemoryPhotoDto `gorm:"-"`
}

type TravelMemoryPhotoDto struct {
	ID          int64  `gorm:"column:id"`
	MemoryDayID int64  `gorm:"column:memory_day_id"`
	ImagePath   string `gorm:"column:image_path"`
	SortOrder   int32  `gorm:"column:sort_order"`
}


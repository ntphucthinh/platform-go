package dbmodel

import "time"

type Homestay struct {
	ID            int64          `gorm:"primaryKey;autoIncrement"`
	Name          string         `gorm:"size:255;not null"`
	Address       string         `gorm:"type:text;not null"`
	Description   *string        `gorm:"type:text"`
	GoogleMapsURL *string        `gorm:"column:google_maps_url;type:text"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP;autoUpdateLine"`
	Rooms         []HomestayRoom `gorm:"foreignKey:HomestayID"`
}

func (Homestay) TableName() string {
	return "homestays"
}

type HomestayRoom struct {
	ID          int64           `gorm:"primaryKey;autoIncrement"`
	HomestayID  int64           `gorm:"column:homestay_id;not null;index"`
	Name        string          `gorm:"size:255;not null"`
	Description *string         `gorm:"type:text"`
	PriceAmount *float64        `gorm:"column:price_amount;type:decimal(15,2)"`
	PriceUnit   *string         `gorm:"column:price_unit;size:50"`
	SortOrder   int32           `gorm:"column:sort_order;default:0"`
	CreatedAt   time.Time       `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time       `gorm:"default:CURRENT_TIMESTAMP;autoUpdateLine"`
	Images      []HomestayImage `gorm:"foreignKey:RoomID"`
}

func (HomestayRoom) TableName() string {
	return "homestay_rooms"
}

type HomestayImage struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	RoomID    int64     `gorm:"column:room_id;not null;index"`
	ImagePath string    `gorm:"column:image_path;type:text;not null"`
	SortOrder int32     `gorm:"column:sort_order;default:0"`
	IsCover   bool      `gorm:"column:is_cover;default:false"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;autoUpdateLine"`
}

func (HomestayImage) TableName() string {
	return "homestay_images"
}

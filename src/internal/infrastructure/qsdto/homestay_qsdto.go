package qsdto

type HomestayDto struct {
	ID            int64             `gorm:"column:id"`
	Name          string            `gorm:"column:name"`
	Address       string            `gorm:"column:address"`
	Description   *string           `gorm:"column:description"`
	GoogleMapsURL *string           `gorm:"column:google_maps_url"`
	Rooms         []HomestayRoomDto `gorm:"-"`
}

type HomestayRoomDto struct {
	ID          int64              `gorm:"column:id"`
	HomestayID  int64              `gorm:"column:homestay_id"`
	Name        string             `gorm:"column:name"`
	Description *string            `gorm:"column:description"`
	PriceAmount *float64           `gorm:"column:price_amount"`
	PriceUnit   *string            `gorm:"column:price_unit"`
	SortOrder   int32              `gorm:"column:sort_order"`
	Images      []HomestayImageDto `gorm:"-"`
}

type HomestayImageDto struct {
	ID         int64  `gorm:"column:id"`
	HomestayID int64  `gorm:"column:homestay_id"`
	RoomID     *int64 `gorm:"column:room_id"`
	ImagePath  string `gorm:"column:image_path"`
	SortOrder  int32  `gorm:"column:sort_order"`
	IsCover    bool   `gorm:"column:is_cover"`
}

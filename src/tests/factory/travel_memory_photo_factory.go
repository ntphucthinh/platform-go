package factory

import (
	"platform-go/internal/infrastructure/dbmodel"
)

// TravelMemoryPhotoFactory generates TravelMemoryPhoto dbmodel instances for seeders and tests.
type TravelMemoryPhotoFactory struct {
	memoryDayID int64
	imagePath   string
	sortOrder   int32
}

func NewTravelMemoryPhotoFactory() *TravelMemoryPhotoFactory {
	return &TravelMemoryPhotoFactory{
		memoryDayID: 1,
		imagePath:   "/images/memories/photo_default.jpg",
		sortOrder:   0,
	}
}

func (f *TravelMemoryPhotoFactory) WithMemoryDayID(id int64) *TravelMemoryPhotoFactory {
	f.memoryDayID = id
	return f
}

func (f *TravelMemoryPhotoFactory) WithImagePath(path string) *TravelMemoryPhotoFactory {
	f.imagePath = path
	return f
}

func (f *TravelMemoryPhotoFactory) WithSortOrder(order int32) *TravelMemoryPhotoFactory {
	f.sortOrder = order
	return f
}

func (f *TravelMemoryPhotoFactory) Build() *dbmodel.TravelMemoryPhoto {
	return &dbmodel.TravelMemoryPhoto{
		MemoryDayID: f.memoryDayID,
		ImagePath:   f.imagePath,
		SortOrder:   f.sortOrder,
	}
}


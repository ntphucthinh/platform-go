package factory

import (
	"gorm.io/gorm"

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
		imagePath:   "travel_memories/6ba7b810-9dad-11d1-80b4-00c04fd43000.jpg",
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

// Insert persists the TravelMemoryPhoto to the database
func (f *TravelMemoryPhotoFactory) Insert(db *gorm.DB) (*dbmodel.TravelMemoryPhoto, error) {
	photo := f.Build()
	if err := db.Create(photo).Error; err != nil {
		return nil, err
	}
	return photo, nil
}

// MustInsert persists the TravelMemoryPhoto to the database, panicking on error
func (f *TravelMemoryPhotoFactory) MustInsert(db *gorm.DB) *dbmodel.TravelMemoryPhoto {
	photo, err := f.Insert(db)
	if err != nil {
		panic(err)
	}
	return photo
}


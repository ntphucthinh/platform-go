package factory

import (
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
)

// HomestayImageFactory is a factory for generating HomestayImage dbmodel instances for tests and seeders.
//
// Usage:
//
//	// Create a HomestayImage with default values
//	image := factory.NewHomestayImageFactory().Build()
//
//	// Create a HomestayImage with specific room ID and path
//	image := factory.NewHomestayImageFactory().
//	    WithRoomID(roomID).
//	    WithImagePath("homestays/room-1.jpg").
//	    WithCover(true).
//	    Build()
//
//	// Insert into DB
//	image, err := factory.NewHomestayImageFactory().Insert(db)
//	// or
//	image := factory.NewHomestayImageFactory().MustInsert(db)
type HomestayImageFactory struct {
	roomID    int64
	imagePath string
	sortOrder int32
	isCover   bool
}

// NewHomestayImageFactory creates a new HomestayImageFactory with default values
func NewHomestayImageFactory() *HomestayImageFactory {
	randomID := rand.Int63n(1000000)
	imagePath := fmt.Sprintf("homestays/seed_image_%d.jpg", randomID)

	return &HomestayImageFactory{
		roomID:    0,
		imagePath: imagePath,
		sortOrder: 1,
		isCover:   true,
	}
}

// WithRoomID sets the room ID
func (f *HomestayImageFactory) WithRoomID(roomID int64) *HomestayImageFactory {
	f.roomID = roomID
	return f
}

// WithImagePath sets the image path
func (f *HomestayImageFactory) WithImagePath(path string) *HomestayImageFactory {
	f.imagePath = path
	return f
}

// WithSortOrder sets the sort order
func (f *HomestayImageFactory) WithSortOrder(sortOrder int32) *HomestayImageFactory {
	f.sortOrder = sortOrder
	return f
}

// WithCover sets whether the image is a cover image
func (f *HomestayImageFactory) WithCover(isCover bool) *HomestayImageFactory {
	f.isCover = isCover
	return f
}

// Build creates a new dbmodel.HomestayImage instance
func (f *HomestayImageFactory) Build() *dbmodel.HomestayImage {
	now := time.Now()
	return &dbmodel.HomestayImage{
		RoomID:    f.roomID,
		ImagePath: f.imagePath,
		SortOrder: f.sortOrder,
		IsCover:   f.isCover,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// BuildWithID creates a new dbmodel.HomestayImage instance with a specific ID
func (f *HomestayImageFactory) BuildWithID(id int64) *dbmodel.HomestayImage {
	now := time.Now()
	return &dbmodel.HomestayImage{
		ID:        id,
		RoomID:    f.roomID,
		ImagePath: f.imagePath,
		SortOrder: f.sortOrder,
		IsCover:   f.isCover,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Insert persists the HomestayImage to the database
func (f *HomestayImageFactory) Insert(db *gorm.DB) (*dbmodel.HomestayImage, error) {
	image := f.Build()
	if err := db.Create(image).Error; err != nil {
		return nil, err
	}
	return image, nil
}

// MustInsert persists the HomestayImage to the database, panicking on error
func (f *HomestayImageFactory) MustInsert(db *gorm.DB) *dbmodel.HomestayImage {
	image, err := f.Insert(db)
	if err != nil {
		panic(err)
	}
	return image
}

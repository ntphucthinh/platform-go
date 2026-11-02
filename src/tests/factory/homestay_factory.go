package factory

import (
	"time"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
)

// HomestayFactory is a factory for generating Homestay dbmodel instances for tests and seeders.
//
// Usage:
//
//	// Create a Homestay with default values
//	homestay := factory.NewHomestayFactory().Build()
//
//	// Create a Homestay with specific values and rooms
//	homestay := factory.NewHomestayFactory().
//	    WithName("Pine Hill Homestay").
//	    WithAddress("10 Dang Thai Than, Da Lat").
//	    WithRooms([]dbmodel.HomestayRoom{...}).
//	    Build()
//
//	// Insert into DB
//	homestay, err := factory.NewHomestayFactory().Insert(db)
//	// or
//	homestay := factory.NewHomestayFactory().MustInsert(db)
type HomestayFactory struct {
	name          string
	address       string
	description   *string
	googleMapsURL *string
	rooms         []dbmodel.HomestayRoom
}

// NewHomestayFactory creates a new HomestayFactory with default values
func NewHomestayFactory() *HomestayFactory {
	desc := "Cozy homestay with garden and mountain view"
	mapsURL := "https://maps.google.com/?q=dalat"

	return &HomestayFactory{
		name:          "Default Homestay",
		address:       "12 Tran Phu, Da Lat",
		description:   &desc,
		googleMapsURL: &mapsURL,
		rooms:         []dbmodel.HomestayRoom{},
	}
}

// WithHomestay sets values from an existing dbmodel.Homestay
func (f *HomestayFactory) WithHomestay(homestay dbmodel.Homestay) *HomestayFactory {
	f.name = homestay.Name
	f.address = homestay.Address
	f.description = homestay.Description
	f.googleMapsURL = homestay.GoogleMapsURL
	f.rooms = homestay.Rooms
	return f
}

// WithName sets the homestay name
func (f *HomestayFactory) WithName(name string) *HomestayFactory {
	f.name = name
	return f
}

// WithAddress sets the homestay address
func (f *HomestayFactory) WithAddress(address string) *HomestayFactory {
	f.address = address
	return f
}

// WithDescription sets the homestay description
func (f *HomestayFactory) WithDescription(desc string) *HomestayFactory {
	f.description = &desc
	return f
}

// WithGoogleMapsURL sets the Google Maps URL
func (f *HomestayFactory) WithGoogleMapsURL(url string) *HomestayFactory {
	f.googleMapsURL = &url
	return f
}

// WithRooms sets the rooms slice
func (f *HomestayFactory) WithRooms(rooms []dbmodel.HomestayRoom) *HomestayFactory {
	f.rooms = rooms
	return f
}

// AddRoom appends a room to the rooms slice
func (f *HomestayFactory) AddRoom(room dbmodel.HomestayRoom) *HomestayFactory {
	f.rooms = append(f.rooms, room)
	return f
}

// Build creates a new dbmodel.Homestay instance
func (f *HomestayFactory) Build() *dbmodel.Homestay {
	now := time.Now()
	return &dbmodel.Homestay{
		Name:          f.name,
		Address:       f.address,
		Description:   f.description,
		GoogleMapsURL: f.googleMapsURL,
		Rooms:         f.rooms,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// BuildWithID creates a new dbmodel.Homestay instance with a specific ID
func (f *HomestayFactory) BuildWithID(id int64) *dbmodel.Homestay {
	now := time.Now()
	return &dbmodel.Homestay{
		ID:            id,
		Name:          f.name,
		Address:       f.address,
		Description:   f.description,
		GoogleMapsURL: f.googleMapsURL,
		Rooms:         f.rooms,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// Insert persists the Homestay (and associated rooms) to the database
func (f *HomestayFactory) Insert(db *gorm.DB) (*dbmodel.Homestay, error) {
	homestay := f.Build()
	if err := db.Create(homestay).Error; err != nil {
		return nil, err
	}
	return homestay, nil
}

// MustInsert persists the Homestay to the database, panicking on error
func (f *HomestayFactory) MustInsert(db *gorm.DB) *dbmodel.Homestay {
	homestay, err := f.Insert(db)
	if err != nil {
		panic(err)
	}
	return homestay
}

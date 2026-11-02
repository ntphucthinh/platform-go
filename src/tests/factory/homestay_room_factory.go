package factory

import (
	"time"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
)

// HomestayRoomFactory is a factory for generating HomestayRoom dbmodel instances for tests and seeders.
//
// Usage:
//
//	// Create a HomestayRoom with default values
//	room := factory.NewHomestayRoomFactory().Build()
//
//	// Create a HomestayRoom with specific values
//	room := factory.NewHomestayRoomFactory().
//	    WithHomestayID(homestayID).
//	    WithName("Deluxe Room").
//	    WithPrice(600000.0, "night_per_room").
//	    Build()
//
//	// Insert into DB
//	room, err := factory.NewHomestayRoomFactory().Insert(db)
//	// or
//	room := factory.NewHomestayRoomFactory().MustInsert(db)
type HomestayRoomFactory struct {
	homestayID  int64
	name        string
	description *string
	priceAmount *float64
	priceUnit   *string
	sortOrder   int32
}

// NewHomestayRoomFactory creates a new HomestayRoomFactory with default values
func NewHomestayRoomFactory() *HomestayRoomFactory {
	desc := "Comfortable room with garden view"
	price := 500000.0
	unit := "night_per_room"

	return &HomestayRoomFactory{
		homestayID:  0,
		name:        "Standard Room",
		description: &desc,
		priceAmount: &price,
		priceUnit:   &unit,
		sortOrder:   1,
	}
}

// WithHomestayID sets the homestay ID
func (f *HomestayRoomFactory) WithHomestayID(id int64) *HomestayRoomFactory {
	f.homestayID = id
	return f
}

// WithName sets the room name
func (f *HomestayRoomFactory) WithName(name string) *HomestayRoomFactory {
	f.name = name
	return f
}

// WithDescription sets the room description
func (f *HomestayRoomFactory) WithDescription(desc string) *HomestayRoomFactory {
	f.description = &desc
	return f
}

// WithPrice sets the room price amount and unit
func (f *HomestayRoomFactory) WithPrice(amount float64, unit string) *HomestayRoomFactory {
	f.priceAmount = &amount
	f.priceUnit = &unit
	return f
}

// WithPriceAmount sets the price amount
func (f *HomestayRoomFactory) WithPriceAmount(amount float64) *HomestayRoomFactory {
	f.priceAmount = &amount
	return f
}

// WithPriceUnit sets the price unit
func (f *HomestayRoomFactory) WithPriceUnit(unit string) *HomestayRoomFactory {
	f.priceUnit = &unit
	return f
}

// WithSortOrder sets the sort order
func (f *HomestayRoomFactory) WithSortOrder(order int32) *HomestayRoomFactory {
	f.sortOrder = order
	return f
}

// Build creates a new dbmodel.HomestayRoom instance
func (f *HomestayRoomFactory) Build() *dbmodel.HomestayRoom {
	now := time.Now()
	return &dbmodel.HomestayRoom{
		HomestayID:  f.homestayID,
		Name:        f.name,
		Description: f.description,
		PriceAmount: f.priceAmount,
		PriceUnit:   f.priceUnit,
		SortOrder:   f.sortOrder,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// BuildWithID creates a new dbmodel.HomestayRoom instance with a specific ID
func (f *HomestayRoomFactory) BuildWithID(id int64) *dbmodel.HomestayRoom {
	now := time.Now()
	return &dbmodel.HomestayRoom{
		ID:          id,
		HomestayID:  f.homestayID,
		Name:        f.name,
		Description: f.description,
		PriceAmount: f.priceAmount,
		PriceUnit:   f.priceUnit,
		SortOrder:   f.sortOrder,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Insert persists the HomestayRoom to the database
func (f *HomestayRoomFactory) Insert(db *gorm.DB) (*dbmodel.HomestayRoom, error) {
	room := f.Build()
	if err := db.Create(room).Error; err != nil {
		return nil, err
	}
	return room, nil
}

// MustInsert persists the HomestayRoom to the database, panicking on error
func (f *HomestayRoomFactory) MustInsert(db *gorm.DB) *dbmodel.HomestayRoom {
	room, err := f.Insert(db)
	if err != nil {
		panic(err)
	}
	return room
}

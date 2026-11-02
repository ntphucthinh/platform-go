package seeder

import (
	"fmt"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
	"platform-go/tests/factory"
)

type roomSeedDefinition struct {
	name        string
	description string
	price       float64
	unit        string
	sortOrder   int32
}

// CreateHomestayRooms creates initial homestay room seed data based on the provided homestays.
func CreateHomestayRooms(db repository.DBInterface, homestays []dbmodel.Homestay) []dbmodel.HomestayRoom {
	// Check if homestay room data already exists
	var existingRoom dbmodel.HomestayRoom
	result := db.GetDB().First(&existingRoom)
	if result.Error == nil {
		fmt.Println("Homestay room data already exists. Skipping homestay room seed.")
		var rooms []dbmodel.HomestayRoom
		db.GetDB().Find(&rooms)
		return rooms
	}

	// Room definitions per homestay name
	roomTemplates := map[string][]roomSeedDefinition{
		"Pine Hill Homestay": {
			{
				name:        "Pine View Double Room",
				description: "Spacious double room with balcony overlooking the pine forest.",
				price:       650000.0,
				unit:        "night_per_room",
				sortOrder:   1,
			},
			{
				name:        "Family Hill Suite",
				description: "Two-bedroom suite with wooden interior and private fireplace.",
				price:       1200000.0,
				unit:        "night_per_room",
				sortOrder:   2,
			},
		},
		"Old Town Riverside Homestay": {
			{
				name:        "Riverfront Balcony Room",
				description: "Deluxe room with river view and private balcony.",
				price:       850000.0,
				unit:        "night_per_room",
				sortOrder:   1,
			},
			{
				name:        "Entire Ancient Villa",
				description: "Traditional townhouse with private courtyard and kitchen.",
				price:       2500000.0,
				unit:        "night_per_house",
				sortOrder:   2,
			},
		},
		"Misty Valley Homestay": {
			{
				name:        "Mountain Terrace Room",
				description: "Cozy room with scenic views across terraced rice fields.",
				price:       550000.0,
				unit:        "night_per_room",
				sortOrder:   1,
			},
			{
				name:        "Whole Wooden Bungalow",
				description: "Independent eco-bungalow perched above the clouds.",
				price:       1800000.0,
				unit:        "night_per_house",
				sortOrder:   2,
			},
		},
		"Sea Breeze Retreat": {
			{
				name:        "Deluxe Ocean Studio",
				description: "Bright studio room with panoramic sea views and kitchenette.",
				price:       950000.0,
				unit:        "night_per_room",
				sortOrder:   1,
			},
			{
				name:        "Cozy Garden Room",
				description: "Quiet ground-floor room opening onto a tropical courtyard.",
				price:       700000.0,
				unit:        "night_per_room",
				sortOrder:   2,
			},
		},
		"Tam Coc Lotus Homestay": {
			{
				name:        "Karst View Bungalow",
				description: "Traditional brick bungalow facing spectacular limestone hills.",
				price:       600000.0,
				unit:        "night_per_room",
				sortOrder:   1,
			},
			{
				name:        "Lotus House",
				description: "Whole private house overlooking the blooming lotus lake.",
				price:       1500000.0,
				unit:        "night_per_house",
				sortOrder:   2,
			},
		},
	}

	var models []dbmodel.HomestayRoom

	for _, homestay := range homestays {
		definitions, exists := roomTemplates[homestay.Name]
		if !exists {
			definitions = []roomSeedDefinition{
				{
					name:        "Standard Room",
					description: "Comfortable room with standard amenities.",
					price:       500000.0,
					unit:        "night_per_room",
					sortOrder:   1,
				},
				{
					name:        "Deluxe Room",
					description: "Upgraded room with scenic view.",
					price:       800000.0,
					unit:        "night_per_room",
					sortOrder:   2,
				},
			}
		}

		for _, def := range definitions {
			models = append(models, *factory.NewHomestayRoomFactory().
				WithHomestayID(homestay.ID).
				WithName(def.name).
				WithDescription(def.description).
				WithPrice(def.price, def.unit).
				WithSortOrder(def.sortOrder).
				Build())
		}
	}

	// Batch insert into DB
	if len(models) > 0 {
		if err := db.GetDB().Session(&gorm.Session{SkipHooks: true}).Create(&models).Error; err != nil {
			panic(fmt.Sprintf("Failed to create homestay room seed data: %v", err))
		}
		fmt.Printf("Seeded %d homestay rooms successfully.\n", len(models))
	}

	return models
}


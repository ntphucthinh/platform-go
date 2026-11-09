package seeder

import (
	"fmt"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
	"platform-go/tests/factory"
)

// CreateTravelMemoryPhotos creates initial travel memory photo seed data based on the provided travel memory days.
func CreateTravelMemoryPhotos(db repository.DBInterface, days []dbmodel.TravelMemoryDay) []dbmodel.TravelMemoryPhoto {
	// Check if travel memory photo data already exists
	var existingPhoto dbmodel.TravelMemoryPhoto
	result := db.GetDB().First(&existingPhoto)
	if result.Error == nil {
		fmt.Println("Travel memory photo data already exists. Skipping travel memory photo seed.")
		var photos []dbmodel.TravelMemoryPhoto
		db.GetDB().Find(&photos)
		return photos
	}

	photoPool := []string{
		"travel_memories/b1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c01.jpg",
		"travel_memories/b1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c02.jpg",
		"travel_memories/b1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c03.jpg",
		"travel_memories/b1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c04.jpg",
		"travel_memories/b1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c05.jpg",
		"travel_memories/b1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c06.jpg",
		"travel_memories/b1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c07.jpg",
		"travel_memories/b1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c08.jpg",
	}

	var models []dbmodel.TravelMemoryPhoto
	for idx, day := range days {
		// Provide 2 photos per day
		img1 := photoPool[(idx*2)%len(photoPool)]
		img2 := photoPool[(idx*2+1)%len(photoPool)]

		models = append(models, *factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day.ID).
			WithImagePath(img1).
			WithSortOrder(1).
			Build())

		models = append(models, *factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day.ID).
			WithImagePath(img2).
			WithSortOrder(2).
			Build())
	}

	if len(models) > 0 {
		if err := db.GetDB().Session(&gorm.Session{SkipHooks: true}).Create(&models).Error; err != nil {
			panic(fmt.Sprintf("Failed to create travel memory photo seed data: %v", err))
		}
		fmt.Printf("Seeded %d travel memory photos successfully.\n", len(models))
	}

	return models
}

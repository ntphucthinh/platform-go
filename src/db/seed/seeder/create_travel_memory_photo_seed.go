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
		"https://images.unsplash.com/photo-1507525428034-b723cf961d3e",
		"https://images.unsplash.com/photo-1544644181-1484b3fdfc62",
		"https://images.unsplash.com/photo-1528127269322-539801943592",
		"https://images.unsplash.com/photo-1503899036084-c55cdd92da26",
		"https://images.unsplash.com/photo-1493976040374-85c8e12f0c0e",
		"https://images.unsplash.com/photo-1534447677768-be436bb09401",
		"https://images.unsplash.com/photo-1506744038136-46273834b3fb",
		"https://images.unsplash.com/photo-1512453979798-5ea266f8880c",
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

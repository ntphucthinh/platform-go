package seeder

import (
	"fmt"
	"math/rand"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
	"platform-go/tests/factory"
)

// CreateHomestayImages creates 1 image per room with a random path.
func CreateHomestayImages(db repository.DBInterface, rooms []dbmodel.HomestayRoom) []dbmodel.HomestayImage {
	// Check if homestay image data already exists
	var existingImage dbmodel.HomestayImage
	result := db.GetDB().First(&existingImage)
	if result.Error == nil {
		fmt.Println("Homestay image data already exists. Skipping homestay image seed.")
		var images []dbmodel.HomestayImage
		db.GetDB().Find(&images)
		return images
	}

	var models []dbmodel.HomestayImage

	for _, room := range rooms {
		randomSuffix := rand.Int63n(1000000)
		randomPath := fmt.Sprintf("homestays/%d/%d_%d.jpg", room.HomestayID, room.ID, randomSuffix)

		models = append(models, *factory.NewHomestayImageFactory().
			WithRoomID(room.ID).
			WithImagePath(randomPath).
			WithSortOrder(1).
			WithCover(true).
			Build())
	}

	// Batch insert into DB
	if len(models) > 0 {
		if err := db.GetDB().Session(&gorm.Session{SkipHooks: true}).Create(&models).Error; err != nil {
			panic(fmt.Sprintf("Failed to create homestay image seed data: %v", err))
		}
		fmt.Printf("Seeded %d homestay images successfully.\n", len(models))
	}

	return models
}


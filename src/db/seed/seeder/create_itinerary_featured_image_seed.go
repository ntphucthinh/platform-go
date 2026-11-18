package seeder

import (
	"fmt"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
)

// CreateItineraryFeaturedImages creates initial featured images for itineraries.
func CreateItineraryFeaturedImages(db repository.DBInterface, itineraries []dbmodel.Itinerary) []dbmodel.ItineraryFeaturedImage {
	var existingImages []dbmodel.ItineraryFeaturedImage
	db.GetDB().Find(&existingImages)
	if len(existingImages) >= 13 {
		fmt.Printf("Itinerary featured image seed data already contains %d records. Skipping.\n", len(existingImages))
		return existingImages
	}

	imagesConfig := map[string][]string{
		"Khám phá xứ sở sương mù & Săn mây Đà Lạt": {
			"itineraries/featured/dalat-1.jpg",
			"itineraries/featured/dalat-2.jpg",
			"itineraries/featured/dalat-3.jpg",
		},
		"Kỳ nghỉ biển sôi động & Khám phá Bà Nà Hills": {
			"itineraries/featured/danang-1.jpg",
			"itineraries/featured/danang-2.jpg",
			"itineraries/featured/danang-3.jpg",
		},
		"Hành trình 36 phố phường & Tinh hoa Hà Thành": {
			"itineraries/featured/hanoi-1.jpg",
			"itineraries/featured/hanoi-2.jpg",
		},
		"Nghỉ dưỡng du thuyền trên vịnh di sản Hạ Long": {
			"itineraries/featured/halong-1.jpg",
			"itineraries/featured/halong-2.jpg",
		},
		"Chạm vào tương lai và hoài niệm Tokyo": {
			"itineraries/featured/tokyo-1.jpg",
			"itineraries/featured/tokyo-2.jpg",
			"itineraries/featured/tokyo-3.jpg",
		},
	}

	existingMap := make(map[string]bool)
	for _, img := range existingImages {
		existingMap[fmt.Sprintf("%d_%s", img.ItineraryID, img.ImagePath)] = true
	}

	var models []dbmodel.ItineraryFeaturedImage
	for _, it := range itineraries {
		paths, ok := imagesConfig[it.ItineraryName]
		if !ok {
			continue
		}
		for _, p := range paths {
			key := fmt.Sprintf("%d_%s", it.ID, p)
			if existingMap[key] {
				continue
			}
			models = append(models, dbmodel.ItineraryFeaturedImage{
				ItineraryID: it.ID,
				ImagePath:   p,
			})
		}
	}

	if len(models) > 0 {
		if err := db.GetDB().Create(&models).Error; err != nil {
			fmt.Printf("failed to seed itinerary featured images: %v\n", err)
		}
	}

	var allImages []dbmodel.ItineraryFeaturedImage
	db.GetDB().Find(&allImages)
	return allImages
}

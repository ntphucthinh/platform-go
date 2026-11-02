package seeder

import (
	"fmt"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
	"platform-go/tests/factory"
)

// CreateHomestays creates initial homestay seed data (5 items).
func CreateHomestays(db repository.DBInterface) []dbmodel.Homestay {
	// Check if homestay data already exists
	var existingHomestay dbmodel.Homestay
	result := db.GetDB().Where("name = ?", "Pine Hill Homestay").First(&existingHomestay)
	if result.Error == nil {
		fmt.Println("Homestay data already exists. Skipping homestay seed.")
		var homestays []dbmodel.Homestay
		db.GetDB().Find(&homestays)
		return homestays
	}

	var models []dbmodel.Homestay

	// 1. Pine Hill Homestay (Da Lat)
	models = append(models, *factory.NewHomestayFactory().
		WithName("Pine Hill Homestay").
		WithAddress("10 Dang Thai Than, Ward 3, Da Lat, Lam Dong").
		WithDescription("Peaceful hillside retreat surrounded by pine forests.").
		WithGoogleMapsURL("https://maps.google.com/?q=Pine+Hill+Homestay+Da+Lat").
		Build())

	// 2. Old Town Riverside Homestay (Hoi An)
	models = append(models, *factory.NewHomestayFactory().
		WithName("Old Town Riverside Homestay").
		WithAddress("45 Nguyen Du, Minh An, Hoi An, Quang Nam").
		WithDescription("Charming traditional homestay overlooking the Hoai River, walking distance to ancient town.").
		WithGoogleMapsURL("https://maps.google.com/?q=Old+Town+Riverside+Homestay+Hoi+An").
		Build())

	// 3. Misty Valley Homestay (Sa Pa)
	models = append(models, *factory.NewHomestayFactory().
		WithName("Misty Valley Homestay").
		WithAddress("Ta Van Village, Sa Pa, Lao Cai").
		WithDescription("Authentic wooden homestay with breathtaking views of Muong Hoa valley and terraced fields.").
		WithGoogleMapsURL("https://maps.google.com/?q=Misty+Valley+Homestay+Sa+Pa").
		Build())

	// 4. Sea Breeze Retreat (Da Nang)
	models = append(models, *factory.NewHomestayFactory().
		WithName("Sea Breeze Retreat").
		WithAddress("88 Vo Nguyen Giap, Ngu Hanh Son, Da Nang").
		WithDescription("Modern coastal homestay just two minutes walk from My Khe Beach.").
		WithGoogleMapsURL("https://maps.google.com/?q=Sea+Breeze+Retreat+Da+Nang").
		Build())

	// 5. Tam Coc Lotus Homestay (Ninh Binh)
	models = append(models, *factory.NewHomestayFactory().
		WithName("Tam Coc Lotus Homestay").
		WithAddress("Ninh Hai, Hoa Lu, Ninh Binh").
		WithDescription("Quiet countryside escape situated beside lotus ponds and limestone karsts.").
		WithGoogleMapsURL("https://maps.google.com/?q=Tam+Coc+Lotus+Homestay+Ninh+Binh").
		Build())

	// Batch insert into DB
	if len(models) > 0 {
		if err := db.GetDB().Session(&gorm.Session{SkipHooks: true}).Create(&models).Error; err != nil {
			panic(fmt.Sprintf("Failed to create homestay seed data: %v", err))
		}
		fmt.Printf("Seeded %d homestays successfully.\n", len(models))
	}

	return models
}

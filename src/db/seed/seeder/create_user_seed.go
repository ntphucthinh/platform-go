package seeder

import (
	"fmt"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
	"platform-go/tests/factory"
)

// CreateUsers creates initial user seed data.
func CreateUsers(db repository.DBInterface) []dbmodel.User {
	// Check if admin user already exists
	var existingUser dbmodel.User
	result := db.GetDB().Where("email = ?", "admin@platform.com").First(&existingUser)
	if result.Error == nil {
		fmt.Println("User data already exists. Skipping user seed.")
		var users []dbmodel.User
		db.GetDB().Find(&users)
		return users
	}

	var models []dbmodel.User

	// Default administrator
	models = append(models, *factory.NewUserFactory().
		WithName("Administrator").
		WithEmail("admin@platform.com").
		WithPassword("123456").
		WithAccountTypeAdmin().
		WithActive().
		Build())

	// Batch insert into DB
	if len(models) > 0 {
		if err := db.GetDB().Session(&gorm.Session{SkipHooks: true}).Create(&models).Error; err != nil {
			panic(fmt.Sprintf("Failed to create user seed data: %v", err))
		}
		fmt.Printf("Seeded %d users successfully.\n", len(models))
	}

	return models
}

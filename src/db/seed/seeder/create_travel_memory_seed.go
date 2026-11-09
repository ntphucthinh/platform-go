package seeder

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
	"platform-go/tests/factory"
)

type memorySeedDefinition struct {
	locationName   string
	title          string
	startDate      time.Time
	endDate        time.Time
	description    string
	coverImagePath string
}

// CreateTravelMemories creates initial travel memory seed data based on seeded travel locations.
func CreateTravelMemories(db repository.DBInterface, locations []dbmodel.TravelLocation) []dbmodel.TravelMemory {
	// Check if travel memory data already exists
	var existingMemory dbmodel.TravelMemory
	result := db.GetDB().First(&existingMemory)
	if result.Error == nil {
		fmt.Println("Travel memory data already exists. Skipping travel memory seed.")
		var memories []dbmodel.TravelMemory
		db.GetDB().Find(&memories)
		return memories
	}

	locationMap := make(map[string]int64)
	for _, loc := range locations {
		locationMap[loc.LocationName] = loc.ID
	}

	// Fallback location ID if map is empty
	var defaultLocationID int64
	if len(locations) > 0 {
		defaultLocationID = locations[0].ID
	} else {
		var firstLoc dbmodel.TravelLocation
		if err := db.GetDB().First(&firstLoc).Error; err == nil {
			defaultLocationID = firstLoc.ID
		}
	}

	getLocationID := func(name string) int64 {
		if id, exists := locationMap[name]; exists && id > 0 {
			return id
		}
		return defaultLocationID
	}

	items := []memorySeedDefinition{
		{
			locationName:   "Đà Lạt",
			title:          "Chuyến đi mùa hoa dã quỳ Đà Lạt",
			startDate:      time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC),
			description:    "Hành trình săn hoa dã quỳ, ngắm bình minh trên đồi thông và thưởng thức cà phê se lạnh Đà Lạt.",
			coverImagePath: "https://images.unsplash.com/photo-1507525428034-b723cf961d3e",
		},
		{
			locationName:   "Đà Nẵng",
			title:          "Kỳ nghỉ hè sôi động tại thành phố biển Đà Nẵng",
			startDate:      time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2024, 7, 22, 0, 0, 0, 0, time.UTC),
			description:    "Thư giãn tại bãi biển Mỹ Khê, check-in Cầu Vàng Bà Nà Hills và xem Cầu Rồng phun lửa.",
			coverImagePath: "https://images.unsplash.com/photo-1544644181-1484b3fdfc62",
		},
		{
			locationName:   "Huế",
			title:          "Khám phá nét cổ kính và di sản Cố đô Huế",
			startDate:      time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC),
			description:    "Trải nghiệm không gian hoài niệm tại Đại Nội, các lăng tẩm triều Nguyễn và ẩm thực cung đình Huế.",
			coverImagePath: "https://images.unsplash.com/photo-1528127269322-539801943592",
		},
		{
			locationName:   "Tokyo",
			title:          "Hành trình ngắm mùa lá đỏ rực rỡ tại Tokyo",
			startDate:      time.Date(2024, 11, 20, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2024, 11, 23, 0, 0, 0, 0, time.UTC),
			description:    "Dạo bước dưới những tán lá vàng lá đỏ tại vườn Shinjuku Gyoen và đền Asakusa cổ kính.",
			coverImagePath: "https://images.unsplash.com/photo-1503899036084-c55cdd92da26",
		},
	}

	var models []dbmodel.TravelMemory
	for _, item := range items {
		locID := getLocationID(item.locationName)
		if locID == 0 {
			continue
		}

		models = append(models, *factory.NewTravelMemoryFactory().
			WithLocationID(locID).
			WithTitle(item.title).
			WithDates(item.startDate, item.endDate).
			WithDescription(item.description).
			WithCoverImagePath(item.coverImagePath).
			Build())
	}

	if len(models) > 0 {
		if err := db.GetDB().Session(&gorm.Session{SkipHooks: true}).Create(&models).Error; err != nil {
			panic(fmt.Sprintf("Failed to create travel memory seed data: %v", err))
		}
		fmt.Printf("Seeded %d travel memories successfully.\n", len(models))
	}

	return models
}


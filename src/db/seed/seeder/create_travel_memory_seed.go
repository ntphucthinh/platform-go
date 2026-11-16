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
	items := []memorySeedDefinition{
		// 1. Đà Lạt
		{
			locationName:   "Đà Lạt",
			title:          "Chuyến đi mùa hoa dã quỳ Đà Lạt",
			startDate:      time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC),
			description:    "Hành trình săn hoa dã quỳ, ngắm bình minh trên đồi thông và thưởng thức cà phê se lạnh Đà Lạt.",
			coverImagePath: "travel_memories/a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c50.jpg",
		},
		// 2. Đà Nẵng
		{
			locationName:   "Đà Nẵng",
			title:          "Kỳ nghỉ hè sôi động tại thành phố biển Đà Nẵng",
			startDate:      time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2024, 7, 22, 0, 0, 0, 0, time.UTC),
			description:    "Thư giãn tại bãi biển Mỹ Khê, check-in Cầu Vàng Bà Nà Hills và xem Cầu Rồng phun lửa.",
			coverImagePath: "travel_memories/a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c51.jpg",
		},
		// 3. Huế
		{
			locationName:   "Huế",
			title:          "Khám phá nét cổ kính và di sản Cố đô Huế",
			startDate:      time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC),
			description:    "Trải nghiệm không gian hoài niệm tại Đại Nội, các lăng tẩm triều Nguyễn và ẩm thực cung đình Huế.",
			coverImagePath: "travel_memories/a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c52.jpg",
		},
		// 4. Tokyo
		{
			locationName:   "Tokyo",
			title:          "Hành trình ngắm mùa lá đỏ rực rỡ tại Tokyo",
			startDate:      time.Date(2024, 11, 20, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2024, 11, 23, 0, 0, 0, 0, time.UTC),
			description:    "Dạo bước dưới những tán lá vàng lá đỏ tại vườn Shinjuku Gyoen và đền Asakusa cổ kính.",
			coverImagePath: "travel_memories/a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c53.jpg",
		},
		// 5. Hà Nội (No image)
		{
			locationName:   "Hà Nội",
			title:          "Hà Nội mùa thu lá bay và 36 phố phường",
			startDate:      time.Date(2024, 10, 8, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2024, 10, 10, 0, 0, 0, 0, time.UTC),
			description:    "Dạo quanh Hồ Gươm, thưởng thức cà phê trứng và ngắm phố cổ lúc sáng sớm.",
			coverImagePath: "",
		},
		// 6. TP. Hồ Chí Minh (No image)
		{
			locationName:   "TP. Hồ Chí Minh",
			title:          "Sài Gòn hoa lệ và nhịp sống sôi động về đêm",
			startDate:      time.Date(2025, 2, 14, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2025, 2, 16, 0, 0, 0, 0, time.UTC),
			description:    "Khám phá phố đi bộ Nguyễn Huệ, check-in Landmark 81 và ẩm thực đường phố Chợ Lớn.",
			coverImagePath: "",
		},
		// 7. Nha Trang (No image)
		{
			locationName:   "Nha Trang",
			title:          "Thiên đường biển xanh và lặn ngắm san hô Nha Trang",
			startDate:      time.Date(2024, 8, 12, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC),
			description:    "Trải nghiệm cano tour 3 đảo, lặn biển ngắm san hô tại Hòn Mun và tắm bùn khoáng nóng.",
			coverImagePath: "",
		},
		// 8. Cần Thơ (No image)
		{
			locationName:   "Cần Thơ",
			title:          "Hành trình miền Tây sông nước và chợ nổi Cái Răng",
			startDate:      time.Date(2023, 9, 20, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2023, 9, 22, 0, 0, 0, 0, time.UTC),
			description:    "Dậy sớm đi chợ nổi Cái Răng, tham quan vườn trái cây miệt vườn và ăn lẩu mắm miền Tây.",
			coverImagePath: "",
		},
		// 9. Hạ Long (No image)
		{
			locationName:   "Hạ Long",
			title:          "Nghỉ dưỡng du thuyền trên vịnh Hạ Long kỳ vĩ",
			startDate:      time.Date(2024, 4, 29, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			description:    "Trải nghiệm ngủ đêm trên vịnh, chèo kayak qua hang Luồn và ngắm hoàng hôn trên boong tàu.",
			coverImagePath: "",
		},
		// 10. Phú Quốc (No image)
		{
			locationName:   "Phú Quốc",
			title:          "Kỳ nghỉ hè rực rỡ tại đảo ngọc Phú Quốc",
			startDate:      time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2025, 3, 18, 0, 0, 0, 0, time.UTC),
			description:    "Ngắm hoàng hôn Bãi Sao, đi cáp treo Hòn Thơm và thưởng thức hải sản Làng chài Hàm Ninh.",
			coverImagePath: "",
		},
		// 11. Sa Pa (No image)
		{
			locationName:   "Sa Pa",
			title:          "Chinh phục đỉnh Fansipan và săn mây Sa Pa mùa lúa chín",
			startDate:      time.Date(2023, 10, 18, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2023, 10, 21, 0, 0, 0, 0, time.UTC),
			description:    "Check-in nóc nhà Đông Dương Fansipan, ngắm thung lũng Mường Hoa và bản Cát Cát.",
			coverImagePath: "",
		},
		// 12. Paris (No image)
		{
			locationName:   "Paris",
			title:          "Kinh đô ánh sáng và sự lãng mạn tại Paris",
			startDate:      time.Date(2023, 5, 10, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2023, 5, 14, 0, 0, 0, 0, time.UTC),
			description:    "Chiêm ngưỡng tháp Eiffel rực sáng, ghé thăm bảo tàng Louvre và dạo bước dọc sông Seine.",
			coverImagePath: "",
		},
		// 13. New York (No image)
		{
			locationName:   "New York",
			title:          "Khám phá nhịp sống không ngủ tại New York",
			startDate:      time.Date(2024, 9, 18, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC),
			description:    "Dạo bước Central Park mùa thu, ngắm toàn cảnh Manhattan từ Empire State và xem kịch Broadway.",
			coverImagePath: "",
		},
		// 14. London (No image)
		{
			locationName:   "London",
			title:          "Dạo bước bên bờ sông Thames và tháp Big Ben London",
			startDate:      time.Date(2025, 5, 2, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2025, 5, 6, 0, 0, 0, 0, time.UTC),
			description:    "Thăm bảo tàng British Museum, đi vòng quay London Eye và trải nghiệm trà chiều kiểu Anh.",
			coverImagePath: "",
		},
		// 15. Seoul (No image, crosses years 2024 - 2025)
		{
			locationName:   "Seoul",
			title:          "Chuyến đi đón Giao thừa Countdown rực rỡ tại Seoul",
			startDate:      time.Date(2024, 12, 29, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			description:    "Đón năm mới tại tháp N Seoul, ngắm tuyết rơi ở Cung điện Gyeongbokgung và ẩm thực Myeongdong.",
			coverImagePath: "",
		},
	}

	// Check if travel memory data already reached required count
	var existingMemories []dbmodel.TravelMemory
	db.GetDB().Find(&existingMemories)
	if len(existingMemories) >= len(items) {
		fmt.Printf("Travel memory seed data already contains %d records. Skipping travel memory seed.\n", len(existingMemories))
		return existingMemories
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

	existingTitles := make(map[string]bool)
	for _, m := range existingMemories {
		existingTitles[m.Title] = true
	}

	var models []dbmodel.TravelMemory
	for _, item := range items {
		if existingTitles[item.title] {
			continue
		}

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

	var allMemories []dbmodel.TravelMemory
	db.GetDB().Find(&allMemories)
	return allMemories
}

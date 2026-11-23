package seeder

import (
	"fmt"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
	"platform-go/tests/factory"
)

// CreateTravelThemes creates initial travel theme seed data (approx 30 themes).
func CreateTravelThemes(db repository.DBInterface) []dbmodel.TravelTheme {
	themeNames := []string{
		// Mandatory themes requested by user
		"Thiên nhiên",
		"Ẩm thực",
		"Cà phê",
		"Văn hóa địa phương",
		"Biển đảo",
		"Nghỉ dưỡng",
		"Khám phá",
		"Check-in",
		"Lịch sử",

		// Additional popular travel themes
		"Phượt & Trekking",
		"Săn mây",
		"Cắm trại (Camping & Glamping)",
		"Chữa lành (Healing & Wellness)",
		"Du lịch tâm linh",
		"Chụp ảnh nghệ thuật",
		"Leo núi",
		"Lặn ngắm san hô",
		"Chèo SUP & Kayak",
		"Trải nghiệm mạo hiểm",
		"Chợ đêm & Mua sắm",
		"Du lịch bụi (Backpacking)",
		"Du lịch gia đình",
		"Tuần trăng mật (Honeymoon)",
		"Lễ hội & Sự kiện",
		"Du lịch sinh thái",
		"Miệt vườn sông nước",
		"Kiến trúc & Nghệ thuật",
		"Đường phố (Street Life)",
		"Ngắm hoàng hôn & bình minh",
		"Du thuyền (Cruise)",
	}

	var existingThemes []dbmodel.TravelTheme
	if err := db.GetDB().Find(&existingThemes).Error; err != nil {
		fmt.Printf("Failed to query existing travel themes: %v\n", err)
		return nil
	}

	existingMap := make(map[string]bool, len(existingThemes))
	for _, t := range existingThemes {
		existingMap[t.Name] = true
	}

	var models []dbmodel.TravelTheme
	for _, name := range themeNames {
		if existingMap[name] {
			continue
		}
		models = append(models, *factory.NewTravelThemeFactory().WithName(name).Build())
	}

	if len(models) == 0 {
		fmt.Printf("All %d travel theme master data records already exist. Skipping seed.\n", len(themeNames))
		return existingThemes
	}

	if err := db.GetDB().Session(&gorm.Session{SkipHooks: true}).Create(&models).Error; err != nil {
		fmt.Printf("Failed to create travel theme seed data: %v\n", err)
		return existingThemes
	}
	fmt.Printf("Seeded %d missing travel themes successfully.\n", len(models))

	var allThemes []dbmodel.TravelTheme
	db.GetDB().Find(&allThemes)
	return allThemes
}

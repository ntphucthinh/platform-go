package seeder

import (
	"fmt"

	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
	"platform-go/tests/factory"
)

type daySeedDefinition struct {
	dayOffset   int
	description string
}

// CreateTravelMemoryDays creates initial travel memory day seed data based on the provided travel memories.
func CreateTravelMemoryDays(db repository.DBInterface, memories []dbmodel.TravelMemory) []dbmodel.TravelMemoryDay {
	// Check if travel memory day data already exists
	var existingDay dbmodel.TravelMemoryDay
	result := db.GetDB().First(&existingDay)
	if result.Error == nil {
		fmt.Println("Travel memory day data already exists. Skipping travel memory day seed.")
		var days []dbmodel.TravelMemoryDay
		db.GetDB().Find(&days)
		return days
	}

	dayTemplates := map[string][]daySeedDefinition{
		"Chuyến đi mùa hoa dã quỳ Đà Lạt": {
			{
				dayOffset:   0,
				description: "Ngày 1: Đặt chân đến Đà Lạt, dạo quanh bờ Hồ Xuân Hương và ngắm hoàng hôn đồi thông.",
			},
			{
				dayOffset:   1,
				description: "Ngày 2: Săn mây sớm tại đồi chè Cầu Đất và chạy xe ngắm những cung đường hoa dã quỳ vàng rực.",
			},
			{
				dayOffset:   2,
				description: "Ngày 3: Chinh phục đỉnh Langbiang, thưởng thức cà phê và mua quà lưu niệm ở chợ đêm.",
			},
		},
		"Kỳ nghỉ hè sôi động tại thành phố biển Đà Nẵng": {
			{
				dayOffset:   0,
				description: "Ngày 1: Check-in khách sạn gần biển Mỹ Khê, tắm biển chiều và ăn hải sản tươi ngon.",
			},
			{
				dayOffset:   1,
				description: "Ngày 2: Khám phá Cầu Vàng Bà Nà Hills trong mây sớm, tối xem Cầu Rồng phun lửa.",
			},
			{
				dayOffset:   2,
				description: "Ngày 3: Viếng Chùa Linh Ứng trên Bán đảo Sơn Trà và ngắm toàn cảnh vịnh Đà Nẵng.",
			},
		},
		"Khám phá nét cổ kính và di sản Cố đô Huế": {
			{
				dayOffset:   0,
				description: "Ngày 1: Thăm Đại Nội Huế (Hoàng Thành), ngắm Ngọ Môn và điện Thái Hòa.",
			},
			{
				dayOffset:   1,
				description: "Ngày 2: Khám phá Lăng Khải Định, Lăng Tự Đức và đi thuyền rồng ngắm hoàng hôn sông Hương.",
			},
			{
				dayOffset:   2,
				description: "Ngày 3: Thưởng thức bún bò Huế mụ Rơi và ghé chợ Đông Ba mua đặc sản kẹo mè xửng.",
			},
		},
		"Hành trình ngắm mùa lá đỏ rực rỡ tại Tokyo": {
			{
				dayOffset:   0,
				description: "Ngày 1: Hạ cánh sân bay Haneda, check-in khách sạn ở Shinjuku và ngắm giao lộ Shibuya về đêm.",
			},
			{
				dayOffset:   1,
				description: "Ngày 2: Dạo bước dưới những tán cây ngân hạnh vàng rực tại công viên Meiji Jingu Gaien.",
			},
			{
				dayOffset:   2,
				description: "Ngày 3: Chiêm ngưỡng lá phong đỏ tại vườn Shinjuku Gyoen và ghé thăm đền Senso-ji cổ kính.",
			},
			{
				dayOffset:   3,
				description: "Ngày 4: Mua sắm tại Akihabara và Ginza trước khi chuẩn bị bay về Việt Nam.",
			},
		},
	}

	var models []dbmodel.TravelMemoryDay
	for _, memory := range memories {
		definitions, exists := dayTemplates[memory.Title]
		if !exists {
			// Default days based on StartDate and EndDate
			daysCount := int(memory.EndDate.Sub(memory.StartDate).Hours()/24) + 1
			if daysCount <= 0 {
				daysCount = 1
			}
			for i := 0; i < daysCount; i++ {
				definitions = append(definitions, daySeedDefinition{
					dayOffset:   i,
					description: fmt.Sprintf("Ngày %d: Khám phá điểm đến và trải nghiệm văn hóa địa phương.", i+1),
				})
			}
		}

		for _, def := range definitions {
			dayDate := memory.StartDate.AddDate(0, 0, def.dayOffset)
			models = append(models, *factory.NewTravelMemoryDayFactory().
				WithMemoryID(memory.ID).
				WithDate(dayDate).
				WithDescription(def.description).
				Build())
		}
	}

	if len(models) > 0 {
		if err := db.GetDB().Session(&gorm.Session{SkipHooks: true}).Create(&models).Error; err != nil {
			panic(fmt.Sprintf("Failed to create travel memory day seed data: %v", err))
		}
		fmt.Printf("Seeded %d travel memory days successfully.\n", len(models))
	}

	return models
}

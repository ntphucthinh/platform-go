package seeder

import (
	"fmt"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
)

// CreateItineraryThemes creates initial itinerary theme links.
func CreateItineraryThemes(db repository.DBInterface, itineraries []dbmodel.Itinerary, themes []dbmodel.TravelTheme) []dbmodel.ItineraryTheme {
	var existingLinks []dbmodel.ItineraryTheme
	db.GetDB().Find(&existingLinks)
	if len(existingLinks) >= 15 {
		fmt.Printf("Itinerary theme seed data already contains %d records. Skipping.\n", len(existingLinks))
		return existingLinks
	}

	themeMap := make(map[string]int64)
	for _, t := range themes {
		themeMap[t.Name] = t.ID
	}
	if len(themeMap) == 0 {
		var dbThemes []dbmodel.TravelTheme
		db.GetDB().Find(&dbThemes)
		for _, t := range dbThemes {
			themeMap[t.Name] = t.ID
		}
	}

	itineraryThemesConfig := map[string][]string{
		"Khám phá xứ sở sương mù & Săn mây Đà Lạt":      {"Thiên nhiên", "Cà phê", "Check-in"},
		"Kỳ nghỉ biển sôi động & Khám phá Bà Nà Hills":  {"Biển đảo", "Check-in", "Nghỉ dưỡng"},
		"Hành trình 36 phố phường & Tinh hoa Hà Thành":  {"Văn hóa địa phương", "Ẩm thực", "Lịch sử"},
		"Nghỉ dưỡng du thuyền trên vịnh di sản Hạ Long": {"Thiên nhiên", "Nghỉ dưỡng", "Khám phá"},
		"Chạm vào tương lai và hoài niệm Tokyo":         {"Khám phá", "Ẩm thực", "Check-in"},
	}

	existingMap := make(map[string]bool)
	for _, l := range existingLinks {
		existingMap[fmt.Sprintf("%d_%d", l.ItineraryID, l.ThemeID)] = true
	}

	var models []dbmodel.ItineraryTheme
	for _, it := range itineraries {
		names, ok := itineraryThemesConfig[it.ItineraryName]
		if !ok {
			continue
		}
		for _, tName := range names {
			tID, ok := themeMap[tName]
			if !ok {
				continue
			}
			key := fmt.Sprintf("%d_%d", it.ID, tID)
			if existingMap[key] {
				continue
			}
			models = append(models, dbmodel.ItineraryTheme{
				ItineraryID: it.ID,
				ThemeID:     tID,
			})
		}
	}

	if len(models) > 0 {
		if err := db.GetDB().Create(&models).Error; err != nil {
			fmt.Printf("failed to seed itinerary themes: %v\n", err)
		}
	}

	var allLinks []dbmodel.ItineraryTheme
	db.GetDB().Find(&allLinks)
	return allLinks
}

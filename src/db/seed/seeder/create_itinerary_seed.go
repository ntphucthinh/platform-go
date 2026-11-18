package seeder

import (
	"encoding/json"
	"fmt"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
)

type itinerarySeedDefinition struct {
	itineraryName  string
	locationName   string
	suitableFor    string
	description    string
	coverImagePath string
	travelTips     []string
}

// CreateItineraries creates initial itinerary seed data based on travel locations.
func CreateItineraries(db repository.DBInterface, locations []dbmodel.TravelLocation) []dbmodel.Itinerary {
	var existingItineraries []dbmodel.Itinerary
	db.GetDB().Find(&existingItineraries)
	if len(existingItineraries) >= 5 {
		fmt.Printf("Itinerary seed data already contains %d records. Skipping itinerary seed.\n", len(existingItineraries))
		return existingItineraries
	}

	locationMap := make(map[string]int64)
	for _, l := range locations {
		locationMap[l.LocationName] = l.ID
	}
	if len(locationMap) == 0 {
		var dbLocations []dbmodel.TravelLocation
		db.GetDB().Find(&dbLocations)
		for _, l := range dbLocations {
			locationMap[l.LocationName] = l.ID
		}
	}

	items := []itinerarySeedDefinition{
		{
			itineraryName:  "Khám phá xứ sở sương mù & Săn mây Đà Lạt",
			locationName:   "Đà Lạt",
			suitableFor:    "Cặp đôi, Bạn bè",
			description:    "Hành trình 3 ngày 2 đêm thư giãn tại Đà Lạt mộng mơ, săn mây đồi thông, trải nghiệm cà phê view thung lũng và thưởng thức ẩm thực se lạnh.",
			coverImagePath: "itineraries/dalat-cover.jpg",
			travelTips: []string{
				"Nên dậy sớm từ 4h30 nếu muốn săn mây ở đồi Đa Phú hoặc Cầu Đất",
				"Nhiệt độ sáng sớm và tối khá lạnh, hãy chuẩn bị áo ấm và găng tay",
				"Thuê xe máy nên kiểm tra kỹ phanh và đèn vì đường nhiều dốc",
			},
		},
		{
			itineraryName:  "Kỳ nghỉ biển sôi động & Khám phá Bà Nà Hills",
			locationName:   "Đà Nẵng",
			suitableFor:    "Gia đình, Nhóm bạn",
			description:    "Hành trình 3 ngày 2 đêm tận hưởng bãi biển Mỹ Khê tuyệt đẹp, check-in Cầu Vàng Bà Nà Hills huyền ảo và chiêm ngưỡng Cầu Rồng phun lửa.",
			coverImagePath: "itineraries/danang-cover.jpg",
			travelTips: []string{
				"Đặt vé Bà Nà Hills trực tuyến trước để không phải xếp hàng mua vé tại cổng",
				"Cầu Rồng chỉ phun lửa và nước vào đúng 21:00 các tối thứ 7 và Chủ nhật",
				"Biển Mỹ Khê sóng khá êm vào sáng sớm, rất thích hợp tắm biển và ngắm bình minh",
			},
		},
		{
			itineraryName:  "Hành trình 36 phố phường & Tinh hoa Hà Thành",
			locationName:   "Hà Nội",
			suitableFor:    "Một mình, Cặp đôi",
			description:    "Khám phá chiều sâu văn hóa Thăng Long nghìn năm văn hiến, len lỏi qua từng góc phố rêu phong và thưởng thức cà phê trứng trứ danh.",
			coverImagePath: "itineraries/hanoi-cover.jpg",
			travelTips: []string{
				"Phố đi bộ Hồ Gươm chỉ hoạt động từ tối thứ Sáu đến hết ngày Chủ nhật",
				"Khi viếng Lăng Bác và Văn Miếu, cần mặc quần/váy dài qua đầu gối",
				"Nên đi bộ hoặc thuê xích lô để ngắm nhìn phố cổ một cách trọn vẹn nhất",
			},
		},
		{
			itineraryName:  "Nghỉ dưỡng du thuyền trên vịnh di sản Hạ Long",
			locationName:   "Hạ Long",
			suitableFor:    "Cặp đôi, Nghỉ dưỡng",
			description:    "Hành trình 2 ngày 1 đêm thư thái trên du thuyền 5 sao vịnh Hạ Long, chèo thuyền kayak qua các hang động kỳ vĩ và ngắm hoàng hôn biển.",
			coverImagePath: "itineraries/halong-cover.jpg",
			travelTips: []string{
				"Chuẩn bị túi chống nước cho điện thoại khi tham gia chèo thuyền kayak",
				"Buổi tối trên boong tàu có gió mát se lạnh, nên đem theo áo khoác nhẹ",
				"Hãy dậy sớm lúc 6h sáng để tham gia lớp học Thái Cực Quyền đón bình minh trên sundeck",
			},
		},
		{
			itineraryName:  "Chạm vào tương lai và hoài niệm Tokyo",
			locationName:   "Tokyo",
			suitableFor:    "Bạn bè, Khám phá",
			description:    "Hành trình 3 ngày trải nghiệm nhịp sống sôi động tại Tokyo: chiêm ngưỡng ngôi đền cổ Sensoji, thỏa sức mua sắm tại Akihabara và ngắm toàn cảnh từ Shibuya Sky.",
			coverImagePath: "itineraries/tokyo-cover.jpg",
			travelTips: []string{
				"Mua thẻ Suica hoặc Pasmo để đi tàu điện ngầm và thanh toán tiện lợi tại các cửa hàng tiện lợi",
				"Nên đặt trước vé đài quan sát Shibuya Sky khung giờ hoàng hôn ít nhất 2-3 tuần",
				"Nhớ mang theo hộ chiếu để được hoàn thuế Tax-Free trực tiếp khi mua sắm tại các cửa hàng",
			},
		},
	}

	existingMap := make(map[string]bool)
	for _, it := range existingItineraries {
		existingMap[it.ItineraryName] = true
	}

	for _, item := range items {
		if existingMap[item.itineraryName] {
			continue
		}

		var locID *int64
		if id, ok := locationMap[item.locationName]; ok {
			locID = &id
		}

		travelTipsJSON, _ := json.Marshal(item.travelTips)
		travelTipsStr := string(travelTipsJSON)

		itinerary := dbmodel.Itinerary{
			ItineraryName:  item.itineraryName,
			LocationID:     locID,
			SuitableFor:    &item.suitableFor,
			Description:    &item.description,
			CoverImagePath: &item.coverImagePath,
			TravelTips:     &travelTipsStr,
		}

		if err := db.GetDB().Create(&itinerary).Error; err != nil {
			fmt.Printf("failed to create itinerary %s: %v\n", item.itineraryName, err)
		}
	}

	var allItineraries []dbmodel.Itinerary
	db.GetDB().Find(&allItineraries)
	return allItineraries
}

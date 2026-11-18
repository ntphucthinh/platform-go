package seeder

import (
	"encoding/json"
	"fmt"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
)

type dayTemplateDef struct {
	dayNumber int
	title     string
	subtitle  string
	summary   string
	imagePath string
	dayTips   []string
}

// CreateItineraryDays creates initial itinerary days for the given itineraries.
func CreateItineraryDays(db repository.DBInterface, itineraries []dbmodel.Itinerary) []dbmodel.ItineraryDay {
	var existingDays []dbmodel.ItineraryDay
	db.GetDB().Find(&existingDays)
	if len(existingDays) >= 13 {
		fmt.Printf("Itinerary day seed data already contains %d records. Skipping.\n", len(existingDays))
		return existingDays
	}

	daysConfig := map[string][]dayTemplateDef{
		"Khám phá xứ sở sương mù & Săn mây Đà Lạt": {
			{
				dayNumber: 1,
				title:     "Ngày 1: Check-in Đà Lạt & Dạo quanh Hồ Xuân Hương",
				subtitle:  "Chào đón thành phố ngàn hoa và phố núi về đêm",
				summary:   "Bắt đầu hành trình với bánh căn giòn rụm, check-in Ga xe lửa cổ kính và thưởng thức cà phê ngắm hoàng hôn thung lũng.",
				imagePath: "itineraries/days/dalat-day-1.jpg",
				dayTips: []string{
					"Nên ghé quán cà phê trước 16:30 để chọn được bàn view hoàng hôn đẹp nhất",
					"Chợ đêm đông đúc, hãy chú ý tư trang và ví tiền",
				},
			},
			{
				dayNumber: 2,
				title:     "Ngày 2: Săn mây Cầu Đất & Đồi chè ngát xanh",
				subtitle:  "Đắm chìm vào biển mây bồng bềnh và hoàng hôn đồi thông",
				summary:   "Hành trình sáng sớm đón bình minh trên thảm gỗ săn mây, khám phá đồi chè Cầu Đất và thưởng thức tiệc nướng BBQ ngói buổi tối.",
				imagePath: "itineraries/days/dalat-day-2.jpg",
				dayTips: []string{
					"Thời gian mây đẹp nhất từ 5h15 đến 6h30 sáng",
					"Đường đèo Cầu Đất sáng sớm có sương mù, lái xe cần bật đèn pha",
				},
			},
			{
				dayNumber: 3,
				title:     "Ngày 3: Thác Datanla & Mua quà lưu niệm",
				subtitle:  "Trải nghiệm cảm giác mạnh và tạm biệt thành phố ngàn hoa",
				summary:   "Thử thách máng trượt xuyên rừng thông tại thác Datanla, thưởng thức bữa trưa cơm niêu truyền thống và mua sắm đặc sản mứt dâu tây.",
				imagePath: "itineraries/days/dalat-day-3.jpg",
				dayTips: []string{
					"Có thể gửi hành lý tại khách sạn trước khi đi chơi buổi sáng",
					"Nên mua dâu tây tại các vườn dâu uy tín thay vì mua ở lề đường",
				},
			},
		},
		"Kỳ nghỉ biển sôi động & Khám phá Bà Nà Hills": {
			{
				dayNumber: 1,
				title:     "Ngày 1: Chào đón biển Mỹ Khê & Bán đảo Sơn Trà",
				subtitle:  "Hòa mình vào làn nước trong xanh và viếng Chùa Linh Ứng",
				summary:   "Bắt đầu với tô bún chả cá thơm ngon, dạo biển Mỹ Khê, chiêm bái tượng Phật Bà tại chùa Linh Ứng và ngắm Cầu Rồng lung linh rực sáng.",
				imagePath: "itineraries/days/danang-day-1.jpg",
				dayTips: []string{
					"Buổi chiều từ 15h30 lên Sơn Trà thời tiết mát mẻ và có thể thấy voọc chà vá chân nâu",
					"Nên chọn vị trí đứng gần đuôi hoặc đầu Cầu Rồng để tránh bị ướt khi phun nước",
				},
			},
			{
				dayNumber: 2,
				title:     "Ngày 2: Thiên đường Sun World Bà Nà Hills",
				subtitle:  "Lạc bước vào chốn bồng lai tiên cảnh và Làng Pháp cổ",
				summary:   "Trải nghiệm cáp treo đạt nhiều kỷ lục Guinness, dạo bước trên Cầu Vàng trứ danh và vui chơi tại công viên giải trí Fantasy Park.",
				imagePath: "itineraries/days/danang-day-2.jpg",
				dayTips: []string{
					"Khí hậu trên đỉnh Bà Nà thay đổi 4 mùa trong một ngày, nên đem theo áo khoác mỏng",
					"Nên đi chuyến cáp treo sớm trước 8h30 để vắng khách check-in Cầu Vàng",
				},
			},
			{
				dayNumber: 3,
				title:     "Ngày 3: Danh thắng Ngũ Hành Sơn & Chợ Cồn",
				subtitle:  "Khám phá hang động kỳ bí và mua sắm quà miền Trung",
				summary:   "Chinh phục đỉnh Thủy Sơn ngắm cảnh sông Hàn, chiêm bái Động Huyền Không linh thiêng và ghé chợ Cồn thưởng thức ẩm thực đường phố.",
				imagePath: "itineraries/days/danang-day-3.jpg",
				dayTips: []string{
					"Bậc đá ở Ngũ Hành Sơn khá trơn, nên mang giày thể thao có độ bám tốt",
					"Chợ Cồn đóng cửa khu ẩm thực vào khoảng 18h",
				},
			},
		},
		"Hành trình 36 phố phường & Tinh hoa Hà Thành": {
			{
				dayNumber: 1,
				title:     "Ngày 1: Dấu ấn lịch sử & Phố cổ nghìn năm",
				subtitle:  "Lăng Bác, Văn Miếu và không khí náo nhiệt khu phố cổ",
				summary:   "Khởi đầu ngày mới với bát phở gia truyền, viếng Lăng Bác, tìm hiểu cội nguồn đạo học tại Văn Miếu và nhâm nhi cà phê trứng béo ngậy.",
				imagePath: "itineraries/days/hanoi-day-1.jpg",
				dayTips: []string{
					"Lăng Bác đóng cửa vào thứ Hai và thứ Sáu hàng tuần",
					"Uống cà phê Giảng nên chọn vị trí ngồi trên gác lửng hoài cổ",
				},
			},
			{
				dayNumber: 2,
				title:     "Ngày 2: Bình yên Hồ Tây & Trải nghiệm Làng gốm Bát Tràng",
				subtitle:  "Gió mát Tây Hồ và bàn xoay gốm thủ công truyền thống",
				summary:   "Ngắm bình minh Tây Hồ, thăm chùa Trấn Quốc cổ kính, tự tay nặn gốm tại làng Bát Tràng và thưởng thức bún chả Hà Nội.",
				imagePath: "itineraries/days/hanoi-day-2.jpg",
				dayTips: []string{
					"Gốm tự làm sau khi nung và tráng men mất khoảng 1-2 tiếng để hoàn thiện",
					"Buổi chiều ngắm hoàng hôn trên Cầu Long Biên rất thơ mộng",
				},
			},
		},
		"Nghỉ dưỡng du thuyền trên vịnh di sản Hạ Long": {
			{
				dayNumber: 1,
				title:     "Ngày 1: Check-in Du thuyền & Chèo kayak Hang Luồn",
				subtitle:  "Khởi hành vào lòng kỳ quan di sản thế giới",
				summary:   "Làm thủ tục lên tàu, thưởng thức tiệc trưa hải sản giữa biển khơi, khám phá Hang Sửng Sốt và chèo kayak qua cổng đá Hang Luồn.",
				imagePath: "itineraries/days/halong-day-1.jpg",
				dayTips: []string{
					"Có mặt tại Cảng tàu quốc tế Tuần Châu muộn nhất 11h30 để làm thủ tục hải quan",
					"Hang Sửng Sốt có khoảng 50 bậc đá leo, hãy đi giày mềm chống trượt",
				},
			},
			{
				dayNumber: 2,
				title:     "Ngày 2: Bình minh đảo Ti Tốp & Tạm biệt Vịnh",
				subtitle:  "Tập Thái Cực Quyền và ngắm toàn cảnh kỳ quan",
				summary:   "Đón bình minh trên biển với bài tập Taichi, leo 400 bậc đá lên đỉnh đảo Ti Tốp ngắm trọn vẹn vịnh Hạ Long và dùng bữa brunch chia tay.",
				imagePath: "itineraries/days/halong-day-2.jpg",
				dayTips: []string{
					"Đỉnh đảo Ti Tốp có góc chụp ảnh vịnh Hạ Long kinh điển nhất, không nên bỏ lỡ",
					"Chuẩn bị hành lý và thanh toán các dịch vụ phát sinh trước 10h00",
				},
			},
		},
		"Chạm vào tương lai và hoài niệm Tokyo": {
			{
				dayNumber: 1,
				title:     "Ngày 1: Nét cổ kính Asakusa & Biểu tượng Tokyo Skytree",
				subtitle:  "Từ truyền thống nghìn năm đến biểu tượng hiện đại",
				summary:   "Dạo bước dưới chiếc đèn lồng Kaminarimon khổng lồ tại đền Sensoji, nếm thử đồ ăn vặt phố Nakamise và ngắm toàn cảnh Tokyo từ tháp Skytree.",
				imagePath: "itineraries/days/tokyo-day-1.jpg",
				dayTips: []string{
					"Nên thử dịch vụ thuê trang phục Kimono dạo phố Asakusa để có những bức ảnh kỷ niệm ấn tượng",
					"Vé lên Skytree mua online sẽ rẻ hơn mua trực tiếp tại quầy",
				},
			},
			{
				dayNumber: 2,
				title:     "Ngày 2: Thiên đường Akihabara & Hoàng hôn Shibuya Sky",
				subtitle:  "Công nghệ, anime và ngã tư nhộn nhịp nhất hành tinh",
				summary:   "Khám phá các trung tâm điện tử và anime tại Akihabara, chiêm ngưỡng ngã tư Shibuya Scramble từ đài quan sát Shibuya Sky và trải nghiệm quán nhậu Izakaya.",
				imagePath: "itineraries/days/tokyo-day-2.jpg",
				dayTips: []string{
					"Không được mang gậy selfie, chân máy ảnh lên sàn vọng cảnh ngoài trời của Shibuya Sky",
					"Khu phố đèn lồng Omoide Yokocho đường nhỏ hẹp, thích hợp đi nhóm 2-3 người",
				},
			},
			{
				dayNumber: 3,
				title:     "Ngày 3: Đền Meiji Jingu & Phố thời trang Harajuku",
				subtitle:  "Rừng thiêng giữa lòng phố thị và phong cách Harajuku",
				summary:   "Tản bộ trong không gian tĩnh mịch của rừng cây đền Meiji, hòa mình vào sắc màu trẻ trung phố Takeshita và thưởng thức sushi băng chuyền.",
				imagePath: "itineraries/days/tokyo-day-3.jpg",
				dayTips: []string{
					"Sân bay Narita và Haneda đều có cửa hàng miễn thuế bán bánh Tokyo Banana và socola Royce",
					"Nên ra sân bay trước giờ bay ít nhất 3 tiếng để làm thủ tục",
				},
			},
		},
	}

	existingMap := make(map[string]bool)
	for _, d := range existingDays {
		existingMap[fmt.Sprintf("%d_%d", d.ItineraryID, d.DayNumber)] = true
	}

	for _, it := range itineraries {
		defs, ok := daysConfig[it.ItineraryName]
		if !ok {
			continue
		}
		for _, d := range defs {
			key := fmt.Sprintf("%d_%d", it.ID, d.dayNumber)
			if existingMap[key] {
				continue
			}

			tipsJSON, _ := json.Marshal(d.dayTips)
			tipsStr := string(tipsJSON)

			day := dbmodel.ItineraryDay{
				ItineraryID: it.ID,
				DayNumber:   d.dayNumber,
				Title:       d.title,
				Subtitle:    &d.subtitle,
				Summary:     &d.summary,
				ImagePath:   &d.imagePath,
				DayTips:     &tipsStr,
			}
			if err := db.GetDB().Create(&day).Error; err != nil {
				fmt.Printf("failed to seed itinerary day %d for itinerary %d: %v\n", d.dayNumber, it.ID, err)
			}
		}
	}

	var allDays []dbmodel.ItineraryDay
	db.GetDB().Find(&allDays)
	return allDays
}

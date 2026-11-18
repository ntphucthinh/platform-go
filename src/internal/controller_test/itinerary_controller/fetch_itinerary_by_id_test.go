package itinerary_controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/tests/factory"
)

func TestItineraryController_FetchItineraryByID(t *testing.T) {
	t.Run("Normal_GetItineraryDetail", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Lâm Đồng").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		theme1 := factory.NewTravelThemeFactory().
			WithName("Thiên nhiên").
			MustInsert(shared.TestDB)

		theme2 := factory.NewTravelThemeFactory().
			WithName("Ẩm thực").
			MustInsert(shared.TestDB)

		theme3 := factory.NewTravelThemeFactory().
			WithName("Cà phê").
			MustInsert(shared.TestDB)

		itinerary := factory.NewItineraryFactory().
			WithLocationID(&loc.ID).
			WithName("Khám phá Đà Lạt").
			WithSuitableFor("Cặp đôi, nhóm bạn, gia đình").
			WithDescription("Một hành trình nhẹ nhàng để cảm nhận trọn vẹn vẻ đẹp của Đà Lạt.").
			WithCoverImagePath("itineraries/dalat-cover.jpg").
			WithTravelTips(`["Thời tiết Đà Lạt se lạnh vào sáng sớm và ban đêm, nên chuẩn bị áo ấm.", "Thuê xe máy sẽ giúp chủ động di chuyển."]`).
			MustInsert(shared.TestDB)

		factory.NewItineraryThemeFactory().
			WithItineraryID(itinerary.ID).
			WithThemeID(theme1.ID).
			MustInsert(shared.TestDB)
		factory.NewItineraryThemeFactory().
			WithItineraryID(itinerary.ID).
			WithThemeID(theme2.ID).
			MustInsert(shared.TestDB)
		factory.NewItineraryThemeFactory().
			WithItineraryID(itinerary.ID).
			WithThemeID(theme3.ID).
			MustInsert(shared.TestDB)

		featImg1 := factory.NewItineraryFeaturedImageFactory().
			WithItineraryID(itinerary.ID).
			WithImagePath("itineraries/featured/dalat-1.jpg").
			MustInsert(shared.TestDB)
		featImg2 := factory.NewItineraryFeaturedImageFactory().
			WithItineraryID(itinerary.ID).
			WithImagePath("itineraries/featured/dalat-2.jpg").
			MustInsert(shared.TestDB)

		day1 := factory.NewItineraryDayFactory().
			WithItineraryID(itinerary.ID).
			WithDayNumber(1).
			WithTitle("Ngày 1 – Khám phá trung tâm Đà Lạt").
			WithSubtitle("Làm quen với thành phố ngàn hoa.").
			WithSummary("Cảm nhận không khí mát lành của trung tâm thành phố.").
			WithImagePath("itineraries/days/dalat-day-1.jpg").
			WithDayTips(`["Trời se lạnh (15–22°C), nên mang theo áo khoác nhẹ.", "Nên mang giày thể thao đế bám tốt."]`).
			MustInsert(shared.TestDB)

		day2 := factory.NewItineraryDayFactory().
			WithItineraryID(itinerary.ID).
			WithDayNumber(2).
			WithTitle("Ngày 2 – Thiên nhiên và ngoại ô").
			WithSubtitle("Khám phá cảnh quan ngoại ô.").
			WithSummary("Tham quan đồi chè và thưởng thức cà phê săn mây.").
			WithImagePath("itineraries/day-2.jpg").
			WithDayTips(`[]`).
			MustInsert(shared.TestDB)

		act1 := factory.NewItineraryActivityFactory().
			WithItineraryDayID(day1.ID).
			WithStartTime("07:30:00").
			WithTitle("Ăn sáng tại chợ Đà Lạt").
			WithCategory("dining").
			WithDurationMinutes(60).
			WithLocation("Chợ Đà Lạt").
			WithGoogleMapsURL("https://www.google.com/maps?q=11.9427,108.4373").
			WithDescription("Thưởng thức bánh căn và sữa đậu nành nóng.").
			WithNotes("Nên đi sớm trước 8h00 để không phải chờ bàn.").
			WithImagePath("itineraries/activities/dalat-banh-can.jpg").
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, fmt.Sprintf("/itineraries/%d", itinerary.ID), nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var body api.FetchItineraryByIDResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

		data := body.Data
		assert.Equal(t, itinerary.ID, data.Id)
		assert.Equal(t, "Khám phá Đà Lạt", data.ItineraryName)
		require.NotNil(t, data.LocationId)
		assert.Equal(t, loc.ID, *data.LocationId)
		require.NotNil(t, data.LocationName)
		assert.Equal(t, "Lâm Đồng", *data.LocationName)
		assert.Equal(t, 2, data.DurationDays)
		assert.Equal(t, 1, data.DurationNights)
		assert.Equal(t, "2 ngày 1 đêm", data.Duration)
		require.NotNil(t, data.SuitableFor)
		assert.Equal(t, "Cặp đôi, nhóm bạn, gia đình", *data.SuitableFor)
		assert.Equal(t, 1, data.ActivityCount)
		require.NotNil(t, data.Description)
		assert.Equal(t, "Một hành trình nhẹ nhàng để cảm nhận trọn vẹn vẻ đẹp của Đà Lạt.", *data.Description)
		require.NotNil(t, data.CoverImagePath)
		assert.Equal(t, "itineraries/dalat-cover.jpg", *data.CoverImagePath)

		require.Len(t, data.Themes, 3)
		assert.Equal(t, theme1.ID, data.Themes[0].Id)
		assert.Equal(t, "Thiên nhiên", data.Themes[0].Name)
		assert.Equal(t, theme2.ID, data.Themes[1].Id)
		assert.Equal(t, "Ẩm thực", data.Themes[1].Name)
		assert.Equal(t, theme3.ID, data.Themes[2].Id)
		assert.Equal(t, "Cà phê", data.Themes[2].Name)

		assert.Equal(t, []string{
			"Thời tiết Đà Lạt se lạnh vào sáng sớm và ban đêm, nên chuẩn bị áo ấm.",
			"Thuê xe máy sẽ giúp chủ động di chuyển.",
		}, data.TravelTips)

		require.Len(t, data.FeaturedImages, 2)
		assert.Equal(t, featImg1.ID, data.FeaturedImages[0].Id)
		assert.Equal(t, "itineraries/featured/dalat-1.jpg", data.FeaturedImages[0].ImagePath)
		assert.Equal(t, featImg2.ID, data.FeaturedImages[1].Id)
		assert.Equal(t, "itineraries/featured/dalat-2.jpg", data.FeaturedImages[1].ImagePath)

		require.Len(t, data.Days, 2)

		// Day 1
		d1 := data.Days[0]
		assert.Equal(t, day1.ID, d1.Id)
		assert.Equal(t, 1, d1.DayNumber)
		assert.Equal(t, "Ngày 1 – Khám phá trung tâm Đà Lạt", d1.Title)
		require.NotNil(t, d1.Subtitle)
		assert.Equal(t, "Làm quen với thành phố ngàn hoa.", *d1.Subtitle)
		require.NotNil(t, d1.Summary)
		assert.Equal(t, "Cảm nhận không khí mát lành của trung tâm thành phố.", *d1.Summary)
		require.NotNil(t, d1.ImagePath)
		assert.Equal(t, "itineraries/days/dalat-day-1.jpg", *d1.ImagePath)
		assert.Equal(t, []string{
			"Trời se lạnh (15–22°C), nên mang theo áo khoác nhẹ.",
			"Nên mang giày thể thao đế bám tốt.",
		}, d1.DayTips)

		require.Len(t, d1.Activities, 1)
		a1 := d1.Activities[0]
		assert.Equal(t, act1.ID, a1.Id)
		assert.Equal(t, "07:30", a1.StartTime)
		assert.Equal(t, "Ăn sáng tại chợ Đà Lạt", a1.Title)
		require.NotNil(t, a1.ActivityCategory)
		assert.Equal(t, api.ItineraryActivityCategoryDining, *a1.ActivityCategory)
		require.NotNil(t, a1.DurationMinutes)
		assert.Equal(t, 60, *a1.DurationMinutes)
		require.NotNil(t, a1.Location)
		assert.Equal(t, "Chợ Đà Lạt", *a1.Location)
		require.NotNil(t, a1.GoogleMapsUrl)
		assert.Equal(t, "https://www.google.com/maps?q=11.9427,108.4373", *a1.GoogleMapsUrl)
		require.NotNil(t, a1.Description)
		assert.Equal(t, "Thưởng thức bánh căn và sữa đậu nành nóng.", *a1.Description)
		require.NotNil(t, a1.Notes)
		assert.Equal(t, "Nên đi sớm trước 8h00 để không phải chờ bàn.", *a1.Notes)
		require.NotNil(t, a1.ImagePath)
		assert.Equal(t, "itineraries/activities/dalat-banh-can.jpg", *a1.ImagePath)

		// Day 2
		d2 := data.Days[1]
		assert.Equal(t, day2.ID, d2.Id)
		assert.Equal(t, 2, d2.DayNumber)
		assert.Equal(t, "Ngày 2 – Thiên nhiên và ngoại ô", d2.Title)
		assert.Empty(t, d2.DayTips)
		assert.Empty(t, d2.Activities)
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodGet, "/itineraries/999999", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		req := shared.NewUnAuthRequest(t, http.MethodGet, "/itineraries/1", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

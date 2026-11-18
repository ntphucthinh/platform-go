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

func TestItineraryController_FetchItineraries(t *testing.T) {
	t.Run("Normal_GetItinerariesList", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Lâm Đồng").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		desc := "Một hành trình nhẹ nhàng để cảm nhận trọn vẹn vẻ đẹp của Đà Lạt."
		cover := "itineraries/dalat-cover.jpg"

		// Itinerary 1: 3 days, 11 activities
		itinerary1 := factory.NewItineraryFactory().
			WithLocationID(&loc.ID).
			WithName("Khám phá Đà Lạt").
			WithDescription(desc).
			WithCoverImagePath(cover).
			MustInsert(shared.TestDB)

		day1 := factory.NewItineraryDayFactory().
			WithItineraryID(itinerary1.ID).
			WithDayNumber(1).
			WithTitle("Ngày 1").
			MustInsert(shared.TestDB)

		day2 := factory.NewItineraryDayFactory().
			WithItineraryID(itinerary1.ID).
			WithDayNumber(2).
			WithTitle("Ngày 2").
			MustInsert(shared.TestDB)

		day3 := factory.NewItineraryDayFactory().
			WithItineraryID(itinerary1.ID).
			WithDayNumber(3).
			WithTitle("Ngày 3").
			MustInsert(shared.TestDB)

		// 5 activities on day 1
		for i := 1; i <= 5; i++ {
			factory.NewItineraryActivityFactory().
				WithItineraryDayID(day1.ID).
				WithTitle(fmt.Sprintf("Hoạt động 1.%d", i)).
				MustInsert(shared.TestDB)
		}
		// 4 activities on day 2
		for i := 1; i <= 4; i++ {
			factory.NewItineraryActivityFactory().
				WithItineraryDayID(day2.ID).
				WithTitle(fmt.Sprintf("Hoạt động 2.%d", i)).
				MustInsert(shared.TestDB)
		}
		// 2 activities on day 3 (total = 11 activities)
		for i := 1; i <= 2; i++ {
			factory.NewItineraryActivityFactory().
				WithItineraryDayID(day3.ID).
				WithTitle(fmt.Sprintf("Hoạt động 3.%d", i)).
				MustInsert(shared.TestDB)
		}

		// Itinerary 2: 1 day, 2 activities
		itinerary2 := factory.NewItineraryFactory().
			WithLocationID(&loc.ID).
			WithName("Đà Lạt 1 ngày").
			MustInsert(shared.TestDB)

		daySingle := factory.NewItineraryDayFactory().
			WithItineraryID(itinerary2.ID).
			WithDayNumber(1).
			WithTitle("Ngày duy nhất").
			MustInsert(shared.TestDB)

		factory.NewItineraryActivityFactory().
			WithItineraryDayID(daySingle.ID).
			WithTitle("Hoạt động A").
			MustInsert(shared.TestDB)
		factory.NewItineraryActivityFactory().
			WithItineraryDayID(daySingle.ID).
			WithTitle("Hoạt động B").
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, "/itineraries?page=1&perPage=10", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var body api.FetchItinerariesResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

		require.Len(t, body.Data, 2)

		// Verify default sorting: newest id first (itinerary2, then itinerary1)
		assert.Equal(t, itinerary2.ID, body.Data[0].Id)
		assert.Equal(t, "Đà Lạt 1 ngày", body.Data[0].ItineraryName)
		assert.Equal(t, 1, body.Data[0].DurationDays)
		assert.Equal(t, 0, body.Data[0].DurationNights)
		assert.Equal(t, "1 ngày", body.Data[0].Duration)
		assert.Equal(t, 2, body.Data[0].ActivityCount)

		assert.Equal(t, itinerary1.ID, body.Data[1].Id)
		assert.Equal(t, "Khám phá Đà Lạt", body.Data[1].ItineraryName)
		assert.Equal(t, loc.ID, *body.Data[1].LocationId)
		assert.Equal(t, "Lâm Đồng", *body.Data[1].LocationName)
		assert.Equal(t, 3, body.Data[1].DurationDays)
		assert.Equal(t, 2, body.Data[1].DurationNights)
		assert.Equal(t, "3 ngày 2 đêm", body.Data[1].Duration)
		assert.Equal(t, 11, body.Data[1].ActivityCount)
		assert.Equal(t, desc, *body.Data[1].Description)
		assert.Equal(t, cover, *body.Data[1].CoverImagePath)

		// Paginator verification
		assert.Equal(t, 1, body.PaginatorInfo.CurrentPage)
		assert.Equal(t, 10, body.PaginatorInfo.PerPage)
		assert.Equal(t, int64(2), body.PaginatorInfo.Total)
	})

	t.Run("Normal_FilterByLocationId", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc1 := factory.NewTravelLocationFactory().
			WithLocationName("Lâm Đồng").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		loc2 := factory.NewTravelLocationFactory().
			WithLocationName("Đà Nẵng").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		it1 := factory.NewItineraryFactory().
			WithLocationID(&loc1.ID).
			WithName("Tour Lâm Đồng").
			MustInsert(shared.TestDB)

		factory.NewItineraryFactory().
			WithLocationID(&loc2.ID).
			WithName("Tour Đà Nẵng").
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, fmt.Sprintf("/itineraries?page=1&perPage=10&locationId=%d", loc1.ID), nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var body api.FetchItinerariesResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

		require.Len(t, body.Data, 1)
		assert.Equal(t, it1.ID, body.Data[0].Id)
		assert.Equal(t, "Tour Lâm Đồng", body.Data[0].ItineraryName)
	})

	t.Run("Normal_Pagination", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		for i := 1; i <= 3; i++ {
			factory.NewItineraryFactory().
				WithName(fmt.Sprintf("Itinerary %d", i)).
				MustInsert(shared.TestDB)
		}

		req := shared.NewRequest(t, admin, http.MethodGet, "/itineraries?page=2&perPage=2", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var body api.FetchItinerariesResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

		assert.Len(t, body.Data, 1)
		assert.Equal(t, 2, body.PaginatorInfo.CurrentPage)
		assert.Equal(t, 2, body.PaginatorInfo.PerPage)
		assert.Equal(t, int64(3), body.PaginatorInfo.Total)
		assert.Equal(t, 2, body.PaginatorInfo.LastPage)
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		req := shared.NewUnAuthRequest(t, http.MethodGet, "/itineraries", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

package travel_location_controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/tests/factory"
)

func TestTravelLocationController_FetchTravelLocationByID(t *testing.T) {
	t.Run("Normal_GetTravelLocationByIDSuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		lat := 10.776889
		lng := 106.700806
		img := "/images/locations/hcmc.jpg"
		desc := "Dynamic economic hub"
		notes := "Visited Ben Thanh Market"
		mapsURL := "https://maps.google.com/?q=Ho+Chi+Minh+City"

		loc := factory.NewTravelLocationFactory().
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			WithLocationName("TP. Hồ Chí Minh").
			WithLatitude(lat).
			WithLongitude(lng).
			WithImagePath(img).
			WithVisited().
			WithDescription(desc).
			WithNotes(notes).
			WithGoogleMapsURL(mapsURL).
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, fmt.Sprintf("/travel-locations/%d", loc.ID), nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.FetchTravelLocationByIDResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))

		assert.Equal(t, int64(loc.ID), res.Data.Id)
		assert.Equal(t, "TP. Hồ Chí Minh", res.Data.LocationName)
		assert.Equal(t, int64(0), res.Data.VisitCount)
		assert.Equal(t, int64(0), res.Data.ItineraryCount)
		assert.Empty(t, res.Data.Memories)
		assert.Empty(t, res.Data.Itineraries)

		rawJSON := w.Body.String()
		for _, camelField := range []string{
			`"countryCode"`,
			`"countryName"`,
			`"locationName"`,
			`"imagePath"`,
			`"googleMapsUrl"`,
			`"visitCount"`,
			`"itineraryCount"`,
			`"memories"`,
			`"itineraries"`,
		} {
			assert.True(t, strings.Contains(rawJSON, camelField), "expected camelField %s in response", camelField)
		}
	})

	t.Run("Normal_GetTravelLocationByIDWithMemoriesAndItineraries", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			WithLocationName("Đà Lạt").
			WithVisited().
			MustInsert(shared.TestDB)

		otherLoc := factory.NewTravelLocationFactory().
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			WithLocationName("Hà Nội").
			WithVisited().
			MustInsert(shared.TestDB)

		// Memory 1: 2 days, 3 photos total
		mem1Start := time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)
		mem1End := time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)
		mem1 := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			WithTitle("Chuyến đi Đà Lạt mùa hoa dã quỳ").
			WithDates(mem1Start, mem1End).
			WithCoverImagePath("travel_memories/cover1.jpg").
			MustInsert(shared.TestDB)

		day1 := factory.NewTravelMemoryDayFactory().
			WithMemoryID(mem1.ID).
			WithDate(mem1Start).
			MustInsert(shared.TestDB)

		day2 := factory.NewTravelMemoryDayFactory().
			WithMemoryID(mem1.ID).
			WithDate(mem1End).
			MustInsert(shared.TestDB)

		factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day1.ID).
			WithImagePath("travel_memories/photo1.jpg").
			MustInsert(shared.TestDB)
		factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day1.ID).
			WithImagePath("travel_memories/photo2.jpg").
			MustInsert(shared.TestDB)
		factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day2.ID).
			WithImagePath("travel_memories/photo3.jpg").
			MustInsert(shared.TestDB)

		// Memory 2: 1 day, 0 photos
		mem2Start := time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
		mem2End := time.Date(2024, 11, 3, 0, 0, 0, 0, time.UTC)
		mem2 := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			WithTitle("Chuyến đi Đà Lạt đầu tháng").
			WithDates(mem2Start, mem2End).
			WithCoverImagePath("travel_memories/cover2.jpg").
			MustInsert(shared.TestDB)

		factory.NewTravelMemoryDayFactory().
			WithMemoryID(mem2.ID).
			WithDate(mem2Start).
			MustInsert(shared.TestDB)

		// Memory for other location
		factory.NewTravelMemoryFactory().
			WithLocationID(otherLoc.ID).
			WithTitle("Chuyến đi Hà Nội").
			MustInsert(shared.TestDB)

		// Soft-deleted memory for this location
		softDeletedMem := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			WithTitle("Chuyến đi đã xóa").
			MustInsert(shared.TestDB)
		shared.TestDB.Delete(&dbmodel.TravelMemory{}, softDeletedMem.ID)

		// Itinerary 1: 3 days -> "3 ngày 2 đêm"
		locID := loc.ID
		itin1 := factory.NewItineraryFactory().
			WithName("Khám phá Đà Lạt 3N2Đ").
			WithLocationID(&locID).
			MustInsert(shared.TestDB)

		factory.NewItineraryDayFactory().WithItineraryID(itin1.ID).WithDayNumber(1).MustInsert(shared.TestDB)
		factory.NewItineraryDayFactory().WithItineraryID(itin1.ID).WithDayNumber(2).MustInsert(shared.TestDB)
		factory.NewItineraryDayFactory().WithItineraryID(itin1.ID).WithDayNumber(3).MustInsert(shared.TestDB)

		// Itinerary 2: 1 day -> "1 ngày"
		itin2 := factory.NewItineraryFactory().
			WithName("Tour Đà Lạt trong ngày").
			WithLocationID(&locID).
			MustInsert(shared.TestDB)

		factory.NewItineraryDayFactory().WithItineraryID(itin2.ID).WithDayNumber(1).MustInsert(shared.TestDB)

		// Itinerary for other location
		otherLocID := otherLoc.ID
		factory.NewItineraryFactory().
			WithName("Tour Hà Nội").
			WithLocationID(&otherLocID).
			MustInsert(shared.TestDB)

		// Soft-deleted itinerary for this location
		softDeletedItin := factory.NewItineraryFactory().
			WithName("Tour đã xóa").
			WithLocationID(&locID).
			MustInsert(shared.TestDB)
		shared.TestDB.Delete(&dbmodel.Itinerary{}, softDeletedItin.ID)

		req := shared.NewRequest(t, admin, http.MethodGet, fmt.Sprintf("/travel-locations/%d", loc.ID), nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.FetchTravelLocationByIDResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))

		assert.Equal(t, int64(loc.ID), res.Data.Id)
		assert.Equal(t, "Đà Lạt", res.Data.LocationName)
		assert.Equal(t, int64(2), res.Data.VisitCount)
		assert.Equal(t, int64(2), res.Data.ItineraryCount)

		// Assert memories: ordered by start_date DESC, id DESC
		require.Len(t, res.Data.Memories, 2)
		assert.Equal(t, mem1.ID, res.Data.Memories[0].Id)
		assert.Equal(t, "Chuyến đi Đà Lạt mùa hoa dã quỳ", res.Data.Memories[0].Title)
		assert.Equal(t, "travel_memories/cover1.jpg", res.Data.Memories[0].CoverImagePath)
		assert.Equal(t, "2024-11-10", res.Data.Memories[0].StartDate.String())
		assert.Equal(t, "2024-11-12", res.Data.Memories[0].EndDate.String())
		assert.Equal(t, 3, res.Data.Memories[0].PhotoCount)

		assert.Equal(t, mem2.ID, res.Data.Memories[1].Id)
		assert.Equal(t, "Chuyến đi Đà Lạt đầu tháng", res.Data.Memories[1].Title)
		assert.Equal(t, "travel_memories/cover2.jpg", res.Data.Memories[1].CoverImagePath)
		assert.Equal(t, "2024-11-01", res.Data.Memories[1].StartDate.String())
		assert.Equal(t, "2024-11-03", res.Data.Memories[1].EndDate.String())
		assert.Equal(t, 0, res.Data.Memories[1].PhotoCount)

		// Assert itineraries: ordered by id DESC
		require.Len(t, res.Data.Itineraries, 2)
		assert.Equal(t, itin2.ID, res.Data.Itineraries[0].Id)
		assert.Equal(t, "Tour Đà Lạt trong ngày", res.Data.Itineraries[0].Title)
		assert.Equal(t, "1 ngày", res.Data.Itineraries[0].Duration)

		assert.Equal(t, itin1.ID, res.Data.Itineraries[1].Id)
		assert.Equal(t, "Khám phá Đà Lạt 3N2Đ", res.Data.Itineraries[1].Title)
		assert.Equal(t, "3 ngày 2 đêm", res.Data.Itineraries[1].Duration)
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodGet, "/travel-locations/999999", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Error_InvalidID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodGet, "/travel-locations/not-a-number", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

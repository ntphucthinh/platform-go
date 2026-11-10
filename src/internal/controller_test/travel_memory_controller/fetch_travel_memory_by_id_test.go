package travel_memory_controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/tests/factory"
)

func TestTravelMemoryController_FetchTravelMemoryByID(t *testing.T) {
	t.Run("Normal_Success", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			MustInsert(shared.TestDB)

		start := time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC)
		end := time.Date(2025, 3, 13, 0, 0, 0, 0, time.UTC)
		mem := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			WithTitle("Mùa hoa dã quỳ trên cao nguyên Đà Lạt").
			WithDates(start, end).
			WithDescription("Kỷ niệm tuyệt vời tại Đà Lạt").
			WithCoverImagePath("travel_memories/a1b2c3d4-e5f6-7890-abcd-ef1234567890.jpg").
			MustInsert(shared.TestDB)

		day := factory.NewTravelMemoryDayFactory().
			WithMemoryID(mem.ID).
			WithDate(start).
			WithDescription("Ngày đầu tiên đặt chân đến Đà Lạt").
			MustInsert(shared.TestDB)

		photo := factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day.ID).
			WithImagePath("travel_memories/f1e2d3c4-b5a6-7890-1234-567890abcdef.jpg").
			WithSortOrder(1).
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, fmt.Sprintf("/travel-memories/%d", mem.ID), nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.FetchTravelMemoryByIDResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))

		assert.Equal(t, mem.ID, res.Data.Id)
		assert.Equal(t, mem.Title, res.Data.Title)
		assert.Equal(t, loc.ID, res.Data.Location.Id)
		assert.Equal(t, loc.LocationName, res.Data.Location.LocationName)
		assert.Equal(t, loc.CountryCode, res.Data.Location.CountryCode)
		assert.Equal(t, loc.CountryName, res.Data.Location.CountryName)
		assert.Equal(t, "travel_memories/a1b2c3d4-e5f6-7890-abcd-ef1234567890.jpg", res.Data.CoverImagePath)

		require.Len(t, res.Data.Days, 1)
		assert.Equal(t, day.ID, res.Data.Days[0].Id)
		require.Len(t, res.Data.Days[0].Photos, 1)
		assert.Equal(t, photo.ID, res.Data.Days[0].Photos[0].Id)
		assert.Equal(t, photo.ImagePath, res.Data.Days[0].Photos[0].ImagePath)

		rawJSON := w.Body.String()
		for _, camelField := range []string{
			`"locationId"`,
			`"startDate"`,
			`"endDate"`,
			`"coverImagePath"`,
			`"location"`,
			`"locationName"`,
			`"countryCode"`,
			`"countryName"`,
			`"days"`,
			`"memoryId"`,
			`"photos"`,
			`"memoryDayId"`,
			`"imagePath"`,
			`"sortOrder"`,
		} {
			assert.Contains(t, rawJSON, camelField)
		}
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodGet, "/travel-memories/999999", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Error_InvalidID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodGet, "/travel-memories/not-a-number", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		req := shared.NewUnAuthRequest(t, http.MethodGet, "/travel-memories/1", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

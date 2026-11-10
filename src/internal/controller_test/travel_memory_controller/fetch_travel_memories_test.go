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

func TestTravelMemoryController_FetchTravelMemories(t *testing.T) {
	t.Run("Normal_GetTravelMemoriesList", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
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

		req := shared.NewRequest(t, admin, http.MethodGet, "/travel-memories?page=1", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var body api.FetchTravelMemoriesResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

		require.Len(t, body.Data, 1)
		assert.Equal(t, mem.ID, body.Data[0].Id)
		assert.Equal(t, loc.ID, body.Data[0].LocationId)
		assert.Equal(t, loc.LocationName, body.Data[0].LocationName)
		assert.Equal(t, loc.CountryCode, body.Data[0].CountryCode)
		assert.Equal(t, mem.Title, body.Data[0].Title)
		assert.Equal(t, "travel_memories/a1b2c3d4-e5f6-7890-abcd-ef1234567890.jpg", body.Data[0].CoverImagePath)

		require.Len(t, body.Data[0].Days, 1)
		assert.Equal(t, day.ID, body.Data[0].Days[0].Id)
		require.Len(t, body.Data[0].Days[0].Photos, 1)
		assert.Equal(t, photo.ID, body.Data[0].Days[0].Photos[0].Id)
		assert.Equal(t, photo.ImagePath, body.Data[0].Days[0].Photos[0].ImagePath)

		assert.Equal(t, 1, body.PaginatorInfo.CurrentPage)
		assert.Equal(t, 10, body.PaginatorInfo.PerPage)
		assert.Equal(t, int64(1), body.PaginatorInfo.Total)

		rawJSON := w.Body.String()
		for _, camelField := range []string{
			`"locationId"`,
			`"locationName"`,
			`"countryCode"`,
			`"startDate"`,
			`"endDate"`,
			`"coverImagePath"`,
			`"days"`,
			`"memoryId"`,
			`"photos"`,
			`"memoryDayId"`,
			`"imagePath"`,
			`"sortOrder"`,
			`"paginatorInfo"`,
		} {
			assert.Contains(t, rawJSON, camelField)
		}
	})

	t.Run("Normal_FilterByLocationID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc1 := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		loc2 := factory.NewTravelLocationFactory().
			WithLocationName("Tokyo").
			WithCountryCode("JP").
			MustInsert(shared.TestDB)

		mem1 := factory.NewTravelMemoryFactory().
			WithLocationID(loc1.ID).
			WithTitle("Chuyến đi Đà Lạt").
			MustInsert(shared.TestDB)

		_ = factory.NewTravelMemoryFactory().
			WithLocationID(loc2.ID).
			WithTitle("Chuyến đi Tokyo").
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, fmt.Sprintf("/travel-memories?page=1&locationId=%d", loc1.ID), nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var body api.FetchTravelMemoriesResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

		require.Len(t, body.Data, 1)
		assert.Equal(t, mem1.ID, body.Data[0].Id)
		assert.Equal(t, loc1.ID, body.Data[0].LocationId)
	})

	t.Run("Normal_FilterByCountryCode", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		locVN := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		locJP := factory.NewTravelLocationFactory().
			WithLocationName("Tokyo").
			WithCountryCode("JP").
			MustInsert(shared.TestDB)

		memVN := factory.NewTravelMemoryFactory().
			WithLocationID(locVN.ID).
			WithTitle("Chuyến đi Việt Nam").
			MustInsert(shared.TestDB)

		_ = factory.NewTravelMemoryFactory().
			WithLocationID(locJP.ID).
			WithTitle("Chuyến đi Nhật Bản").
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, "/travel-memories?page=1&countryCode=VN", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var body api.FetchTravelMemoriesResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

		require.Len(t, body.Data, 1)
		assert.Equal(t, memVN.ID, body.Data[0].Id)
		assert.Equal(t, "VN", body.Data[0].CountryCode)
	})

	t.Run("Normal_FilterByCountryCode_SnakeCase", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		locJP := factory.NewTravelLocationFactory().
			WithLocationName("Tokyo").
			WithCountryCode("JP").
			MustInsert(shared.TestDB)

		memJP := factory.NewTravelMemoryFactory().
			WithLocationID(locJP.ID).
			WithTitle("Tokyo Trip").
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, "/travel-memories?page=1&country_code=jp", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var body api.FetchTravelMemoriesResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

		require.Len(t, body.Data, 1)
		assert.Equal(t, memJP.ID, body.Data[0].Id)
	})

	t.Run("Normal_FilterByBothCountryAndLocation", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		mem := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, fmt.Sprintf("/travel-memories?page=1&countryCode=VN&locationId=%d", loc.ID), nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var body api.FetchTravelMemoriesResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

		require.Len(t, body.Data, 1)
		assert.Equal(t, mem.ID, body.Data[0].Id)
	})

	t.Run("Error_MissingRequiredParams", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodGet, "/travel-memories", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		req := shared.NewUnAuthRequest(t, http.MethodGet, "/travel-memories?page=1", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

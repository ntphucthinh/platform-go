package travel_memory_controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	travelmemory "platform-go/internal/domain/travel_memory"
	travelmemoryday "platform-go/internal/domain/travel_memory_day"
	travelmemoryphoto "platform-go/internal/domain/travel_memory_photo"
	"platform-go/internal/infrastructure/dbmodel"
	travelmemoryusecase "platform-go/internal/usecases/travel_memories/usecase"
	"platform-go/tests/factory"
)

func TestTravelMemoryController_CreateTravelMemory(t *testing.T) {
	t.Run("Normal_CreateTravelMemorySuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			MustInsert(shared.TestDB)

		startDate := openapi_types.Date{Time: time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)}
		day2Date := openapi_types.Date{Time: time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC)}
		endDate := openapi_types.Date{Time: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)}

		desc := "Hành trình săn hoa dã quỳ và thưởng thức cà phê đồi thông."
		coverPath := "travel_memories/550e8400-e29b-41d4-a716-446655440000.jpg"
		day1Desc := "Ngày 1: Đặt chân đến Đà Lạt, dạo quanh bờ Hồ Xuân Hương."
		day2Desc := "Ngày 2: Săn mây Cầu Đất."

		payload := api.CreateTravelMemoryRequest{
			LocationId:     loc.ID,
			Title:          "Chuyến đi mùa hoa dã quỳ Đà Lạt",
			StartDate:      startDate,
			EndDate:        endDate,
			Description:    &desc,
			CoverImagePath: &coverPath,
			Days: &[]api.CreateTravelMemoryDayInput{
				{
					Date:        startDate,
					Description: &day1Desc,
					Photos: &[]api.TravelMemoryPhotoInput{
						{
							ImagePath: "travel_memories/6ba7b810-9dad-11d1-80b4-00c04fd430c8.png",
						},
						{
							ImagePath: "travel_memories/6ba7b810-9dad-11d1-80b4-00c04fd430c9.jpg",
						},
					},
				},
				{
					Date:        day2Date,
					Description: &day2Desc,
				},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/travel-memories", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.MessageOnlyResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, travelmemoryusecase.CreateTravelMemoryUseCaseSuccessMessage, res.Message)

		var mem dbmodel.TravelMemory
		require.NoError(t, shared.TestDB.Where("title = ?", "Chuyến đi mùa hoa dã quỳ Đà Lạt").First(&mem).Error)
		assert.Equal(t, loc.ID, mem.LocationID)
		assert.Equal(t, coverPath, mem.CoverImagePath)
		require.NotNil(t, mem.Description)
		assert.Equal(t, desc, *mem.Description)

		var days []dbmodel.TravelMemoryDay
		require.NoError(t, shared.TestDB.Where("memory_id = ?", mem.ID).Order("date ASC").Find(&days).Error)
		require.Len(t, days, 2)
		require.NotNil(t, days[0].Description)
		assert.Equal(t, day1Desc, *days[0].Description)

		var photos []dbmodel.TravelMemoryPhoto
		require.NoError(t, shared.TestDB.Where("memory_day_id = ?", days[0].ID).Order("sort_order ASC").Find(&photos).Error)
		require.Len(t, photos, 2)
		assert.Equal(t, "travel_memories/6ba7b810-9dad-11d1-80b4-00c04fd430c8.png", photos[0].ImagePath)
		assert.Equal(t, int32(1), photos[0].SortOrder)
		assert.Equal(t, "travel_memories/6ba7b810-9dad-11d1-80b4-00c04fd430c9.jpg", photos[1].ImagePath)
		assert.Equal(t, int32(2), photos[1].SortOrder)
	})

	t.Run("Normal_CreateTravelMemoryWithoutOptionalFields", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Huế").
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			MustInsert(shared.TestDB)

		startDate := openapi_types.Date{Time: time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)}
		endDate := openapi_types.Date{Time: time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC)}

		payload := api.CreateTravelMemoryRequest{
			LocationId: loc.ID,
			Title:      "Khám phá Cố đô Huế",
			StartDate:  startDate,
			EndDate:    endDate,
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/travel-memories", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var mem dbmodel.TravelMemory
		require.NoError(t, shared.TestDB.Where("title = ?", "Khám phá Cố đô Huế").First(&mem).Error)
		assert.Equal(t, loc.ID, mem.LocationID)
		assert.Equal(t, "", mem.CoverImagePath)
		assert.Nil(t, mem.Description)
	})

	t.Run("Error_InvalidLocationID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		startDate := openapi_types.Date{Time: time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)}
		endDate := openapi_types.Date{Time: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)}

		payload := api.CreateTravelMemoryRequest{
			LocationId: 999999,
			Title:      "Chuyến đi Đà Lạt",
			StartDate:  startDate,
			EndDate:    endDate,
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/travel-memories", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), "The specified travel location ID is invalid.")
	})

	t.Run("Error_InvalidDateRange", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			MustInsert(shared.TestDB)

		startDate := openapi_types.Date{Time: time.Date(2024, 11, 15, 0, 0, 0, 0, time.UTC)}
		endDate := openapi_types.Date{Time: time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)}

		payload := api.CreateTravelMemoryRequest{
			LocationId: loc.ID,
			Title:      "Chuyến đi Đà Lạt",
			StartDate:  startDate,
			EndDate:    endDate,
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/travel-memories", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), travelmemory.ErrInvalidDateRange.Error())
	})

	t.Run("Error_DayDateOutOfRange", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			MustInsert(shared.TestDB)

		startDate := openapi_types.Date{Time: time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)}
		endDate := openapi_types.Date{Time: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)}
		outOfRangeDate := openapi_types.Date{Time: time.Date(2024, 11, 15, 0, 0, 0, 0, time.UTC)}

		payload := api.CreateTravelMemoryRequest{
			LocationId: loc.ID,
			Title:      "Chuyến đi Đà Lạt",
			StartDate:  startDate,
			EndDate:    endDate,
			Days: &[]api.CreateTravelMemoryDayInput{
				{
					Date: outOfRangeDate,
				},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/travel-memories", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), travelmemoryday.ErrInvalidDayDate.Error())
	})

	t.Run("Error_DuplicateDayDate", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			MustInsert(shared.TestDB)

		startDate := openapi_types.Date{Time: time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)}
		endDate := openapi_types.Date{Time: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)}

		payload := api.CreateTravelMemoryRequest{
			LocationId: loc.ID,
			Title:      "Chuyến đi Đà Lạt",
			StartDate:  startDate,
			EndDate:    endDate,
			Days: &[]api.CreateTravelMemoryDayInput{
				{Date: startDate},
				{Date: startDate},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/travel-memories", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), travelmemoryday.ErrDuplicateDayDate.Error())
	})

	t.Run("Error_InvalidCoverImagePath", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			MustInsert(shared.TestDB)

		startDate := openapi_types.Date{Time: time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)}
		endDate := openapi_types.Date{Time: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)}
		badCover := "homestays/550e8400-e29b-41d4-a716-446655440000.jpg"

		payload := api.CreateTravelMemoryRequest{
			LocationId:     loc.ID,
			Title:          "Chuyến đi Đà Lạt",
			StartDate:      startDate,
			EndDate:        endDate,
			CoverImagePath: &badCover,
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/travel-memories", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), travelmemory.ErrInvalidImagePath.Error())
	})

	t.Run("Error_InvalidPhotoImagePath", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			MustInsert(shared.TestDB)

		startDate := openapi_types.Date{Time: time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)}
		endDate := openapi_types.Date{Time: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)}

		payload := api.CreateTravelMemoryRequest{
			LocationId: loc.ID,
			Title:      "Chuyến đi Đà Lạt",
			StartDate:  startDate,
			EndDate:    endDate,
			Days: &[]api.CreateTravelMemoryDayInput{
				{
					Date: startDate,
					Photos: &[]api.TravelMemoryPhotoInput{
						{
							ImagePath: "travel_location/550e8400-e29b-41d4-a716-446655440000.jpg",
						},
					},
				},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/travel-memories", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), travelmemoryphoto.ErrInvalidImagePath.Error())
	})

	t.Run("Error_InvalidBody", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodPost, "/travel-memories", strings.NewReader("{invalid-json"))
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		startDate := openapi_types.Date{Time: time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)}
		endDate := openapi_types.Date{Time: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)}

		payload := api.CreateTravelMemoryRequest{
			LocationId: 1,
			Title:      "Chuyến đi Đà Lạt",
			StartDate:  startDate,
			EndDate:    endDate,
		}

		req := shared.NewUnAuthRequest(t, http.MethodPost, "/travel-memories", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}


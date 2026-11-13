package travel_memory_controller_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	travelmemory "platform-go/internal/domain/travel_memory"
	travelmemoryday "platform-go/internal/domain/travel_memory_day"
	travelmemoryphoto "platform-go/internal/domain/travel_memory_photo"
	"platform-go/internal/infrastructure/dbmodel"
	travelmemoryusecase "platform-go/internal/usecases/travel_memories/usecase"
	"platform-go/tests/factory"
)

func TestTravelMemoryController_UpdateTravelMemory(t *testing.T) {
	t.Run("Normal_UpdateTravelMemorySuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc1 := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			MustInsert(shared.TestDB)

		loc2 := factory.NewTravelLocationFactory().
			WithLocationName("Nha Trang").
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			MustInsert(shared.TestDB)

		start := time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)

		initialMemory := factory.NewTravelMemoryFactory().
			WithLocationID(loc1.ID).
			WithTitle("Chuyến đi Đà Lạt ban đầu").
			WithDates(start, end).
			MustInsert(shared.TestDB)

		day1 := factory.NewTravelMemoryDayFactory().
			WithMemoryID(initialMemory.ID).
			WithDate(start).
			WithDescription("Ngày 1 cũ").
			MustInsert(shared.TestDB)

		day2 := factory.NewTravelMemoryDayFactory().
			WithMemoryID(initialMemory.ID).
			WithDate(start.AddDate(0, 0, 1)).
			WithDescription("Ngày 2 sẽ bị xóa").
			MustInsert(shared.TestDB)

		photo1 := factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day1.ID).
			WithImagePath("travel_memories/photo1.jpg").
			WithSortOrder(1).
			MustInsert(shared.TestDB)

		photo2 := factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day1.ID).
			WithImagePath("travel_memories/photo2-to-delete.jpg").
			WithSortOrder(2).
			MustInsert(shared.TestDB)

		photo3 := factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day2.ID).
			WithImagePath("travel_memories/photo3-will-be-deleted.jpg").
			WithSortOrder(1).
			MustInsert(shared.TestDB)

		newStart := time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)
		newEnd := time.Date(2024, 11, 15, 0, 0, 0, 0, time.UTC)
		newDay3Date := time.Date(2024, 11, 14, 0, 0, 0, 0, time.UTC)
		newDay4Date := time.Date(2024, 11, 15, 0, 0, 0, 0, time.UTC)
		newDesc := "Chuyến đi Nha Trang đã cập nhật"
		newCover := "travel_memories/cover-updated.jpg"
		day1UpdatedDesc := "Ngày 1 đã cập nhật"
		day3Desc := "Ngày 3 mới tạo"
		day4Desc := "Ngày 4 mới tạo"

		payload := api.UpdateTravelMemoryRequest{
			LocationId:     loc2.ID,
			Title:          "Chuyến đi đã cập nhật",
			StartDate:      openapi_types.Date{Time: newStart},
			EndDate:        openapi_types.Date{Time: newEnd},
			Description:    &newDesc,
			CoverImagePath: &newCover,
			Days: &[]api.UpdateTravelMemoryDayInput{
				{
					Id:          &day1.ID,
					Date:        openapi_types.Date{Time: newStart},
					Description: &day1UpdatedDesc,
					Photos: &[]api.UpdateTravelMemoryPhotoInput{
						{
							ImagePath: "travel_memories/photo4-new.jpg",
						},
						{
							Id:        &photo1.ID,
							ImagePath: "travel_memories/photo1-updated.jpg",
						},
					},
				},
				{
					Date:        openapi_types.Date{Time: newDay3Date},
					Description: &day3Desc,
					Photos: &[]api.UpdateTravelMemoryPhotoInput{
						{
							ImagePath: "travel_memories/photo5-new.jpg",
						},
						{
							ImagePath: "travel_memories/photo6-new.jpg",
						},
					},
				},
				{
					Date:        openapi_types.Date{Time: newDay4Date},
					Description: &day4Desc,
					Photos: &[]api.UpdateTravelMemoryPhotoInput{
						{
							ImagePath: "travel_memories/photo7-dup.jpg",
						},
						{
							ImagePath: "travel_memories/photo7-dup.jpg",
						},
					},
				},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/travel-memories/%d", initialMemory.ID), payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.MessageOnlyResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, travelmemoryusecase.UpdateTravelMemoryUseCaseSuccessMessage, res.Message)

		// Verify DB changes
		var updatedMemory dbmodel.TravelMemory
		require.NoError(t, shared.TestDB.First(&updatedMemory, initialMemory.ID).Error)
		assert.Equal(t, loc2.ID, updatedMemory.LocationID)
		assert.Equal(t, "Chuyến đi đã cập nhật", updatedMemory.Title)
		require.NotNil(t, updatedMemory.Description)
		assert.Equal(t, newDesc, *updatedMemory.Description)
		assert.Equal(t, newCover, updatedMemory.CoverImagePath)

		// Day 1 updated
		var updatedDay1 dbmodel.TravelMemoryDay
		require.NoError(t, shared.TestDB.First(&updatedDay1, day1.ID).Error)
		require.NotNil(t, updatedDay1.Description)
		assert.Equal(t, day1UpdatedDesc, *updatedDay1.Description)

		// Day 2 deleted
		var deletedDay2 dbmodel.TravelMemoryDay
		assert.True(t, errors.Is(shared.TestDB.First(&deletedDay2, day2.ID).Error, gorm.ErrRecordNotFound))

		// Verify Day 1 Photos & SortOrder (new photo at position 1, updated photo1 at position 2)
		var day1Photos []dbmodel.TravelMemoryPhoto
		require.NoError(t, shared.TestDB.Where("memory_day_id = ?", day1.ID).Order("sort_order ASC").Find(&day1Photos).Error)
		require.Len(t, day1Photos, 2)
		assert.Equal(t, day1.ID, day1Photos[0].MemoryDayID)
		assert.Equal(t, "travel_memories/photo4-new.jpg", day1Photos[0].ImagePath)
		assert.Equal(t, int32(1), day1Photos[0].SortOrder)
		assert.Equal(t, photo1.ID, day1Photos[1].ID)
		assert.Equal(t, day1.ID, day1Photos[1].MemoryDayID)
		assert.Equal(t, "travel_memories/photo1-updated.jpg", day1Photos[1].ImagePath)
		assert.Equal(t, int32(2), day1Photos[1].SortOrder)

		// Photo 2 deleted
		var deletedPhoto2 dbmodel.TravelMemoryPhoto
		assert.True(t, errors.Is(shared.TestDB.First(&deletedPhoto2, photo2.ID).Error, gorm.ErrRecordNotFound))

		// Photo 3 (of day 2) deleted
		var deletedPhoto3 dbmodel.TravelMemoryPhoto
		assert.True(t, errors.Is(shared.TestDB.First(&deletedPhoto3, photo3.ID).Error, gorm.ErrRecordNotFound))

		// Verify multiple newly created Days (Day 3 and Day 4) and their Photos receive the correct generated Day IDs (> 0)
		var createdDays []dbmodel.TravelMemoryDay
		require.NoError(t, shared.TestDB.Where("memory_id = ? AND id != ?", initialMemory.ID, day1.ID).Order("date ASC").Find(&createdDays).Error)
		require.Len(t, createdDays, 2)

		createdDay3 := createdDays[0]
		require.Greater(t, createdDay3.ID, int64(0))
		require.NotNil(t, createdDay3.Description)
		assert.Equal(t, day3Desc, *createdDay3.Description)

		var day3Photos []dbmodel.TravelMemoryPhoto
		require.NoError(t, shared.TestDB.Where("memory_day_id = ?", createdDay3.ID).Order("sort_order ASC").Find(&day3Photos).Error)
		require.Len(t, day3Photos, 2)
		assert.Equal(t, createdDay3.ID, day3Photos[0].MemoryDayID)
		assert.Equal(t, "travel_memories/photo5-new.jpg", day3Photos[0].ImagePath)
		assert.Equal(t, int32(1), day3Photos[0].SortOrder)
		assert.Equal(t, createdDay3.ID, day3Photos[1].MemoryDayID)
		assert.Equal(t, "travel_memories/photo6-new.jpg", day3Photos[1].ImagePath)
		assert.Equal(t, int32(2), day3Photos[1].SortOrder)

		createdDay4 := createdDays[1]
		require.Greater(t, createdDay4.ID, int64(0))
		assert.NotEqual(t, createdDay3.ID, createdDay4.ID)
		require.NotNil(t, createdDay4.Description)
		assert.Equal(t, day4Desc, *createdDay4.Description)

		// Verify Day 4 Photos (allowing duplicate ImagePath values within the same Day)
		var day4Photos []dbmodel.TravelMemoryPhoto
		require.NoError(t, shared.TestDB.Where("memory_day_id = ?", createdDay4.ID).Order("sort_order ASC").Find(&day4Photos).Error)
		require.Len(t, day4Photos, 2)
		assert.Equal(t, createdDay4.ID, day4Photos[0].MemoryDayID)
		assert.Equal(t, "travel_memories/photo7-dup.jpg", day4Photos[0].ImagePath)
		assert.Equal(t, int32(1), day4Photos[0].SortOrder)
		assert.Equal(t, createdDay4.ID, day4Photos[1].MemoryDayID)
		assert.Equal(t, "travel_memories/photo7-dup.jpg", day4Photos[1].ImagePath)
		assert.Equal(t, int32(2), day4Photos[1].SortOrder)
	})

	t.Run("Normal_DaysNil_EmptyDays_PhotosNil_EmptyPhotos_SameImageAcrossDays_ZeroValues", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Nẵng").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		start := time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 11, 15, 0, 0, 0, 0, time.UTC)

		memory := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			WithTitle("Initial Title").
			WithDescription("Initial Description").
			WithCoverImagePath("travel_memories/initial-cover.jpg").
			WithDates(start, end).
			MustInsert(shared.TestDB)

		day1 := factory.NewTravelMemoryDayFactory().
			WithMemoryID(memory.ID).
			WithDate(start).
			WithDescription("Day 1 Desc").
			MustInsert(shared.TestDB)

		day2 := factory.NewTravelMemoryDayFactory().
			WithMemoryID(memory.ID).
			WithDate(start.AddDate(0, 0, 1)).
			WithDescription("Day 2 Desc").
			MustInsert(shared.TestDB)

		photo1 := factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day1.ID).
			WithImagePath("travel_memories/day1-p1.jpg").
			WithSortOrder(1).
			MustInsert(shared.TestDB)

		photo2 := factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day2.ID).
			WithImagePath("travel_memories/day2-p1.jpg").
			WithSortOrder(1).
			MustInsert(shared.TestDB)

		// 1. Days == nil: existing Days and Photos remain unchanged, while optional Memory fields are cleared to nil/empty
		payloadDaysNil := api.UpdateTravelMemoryRequest{
			LocationId:     loc.ID,
			Title:          "Updated Title Without Touching Days",
			StartDate:      openapi_types.Date{Time: start},
			EndDate:        openapi_types.Date{Time: end},
			Description:    nil,
			CoverImagePath: nil,
			Days:           nil,
		}
		req1 := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/travel-memories/%d", memory.ID), payloadDaysNil)
		w1 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w1, req1)
		require.Equal(t, http.StatusOK, w1.Code, w1.Body.String())

		var memAfter1 dbmodel.TravelMemory
		require.NoError(t, shared.TestDB.First(&memAfter1, memory.ID).Error)
		assert.Nil(t, memAfter1.Description)
		assert.Equal(t, "", memAfter1.CoverImagePath)

		var daysAfter1 []dbmodel.TravelMemoryDay
		require.NoError(t, shared.TestDB.Where("memory_id = ?", memory.ID).Find(&daysAfter1).Error)
		assert.Len(t, daysAfter1, 2)

		// 2. Photos == nil on Day 1 (keeps photo1), Empty Photos on Day 2 (deletes photo2),
		// plus clear Day 1 Description to nil, and add Day 3 with same ImagePath as Day 1
		emptyPhotos := []api.UpdateTravelMemoryPhotoInput{}
		payloadPhotosSync := api.UpdateTravelMemoryRequest{
			LocationId: loc.ID,
			Title:      "Sync Photos States",
			StartDate:  openapi_types.Date{Time: start},
			EndDate:    openapi_types.Date{Time: end},
			Days: &[]api.UpdateTravelMemoryDayInput{
				{
					Id:          &day1.ID,
					Date:        openapi_types.Date{Time: start},
					Description: nil,
					Photos:      nil,
				},
				{
					Id:     &day2.ID,
					Date:   openapi_types.Date{Time: start.AddDate(0, 0, 1)},
					Photos: &emptyPhotos,
				},
				{
					Date: openapi_types.Date{Time: start.AddDate(0, 0, 2)},
					Photos: &[]api.UpdateTravelMemoryPhotoInput{
						{ImagePath: "travel_memories/day1-p1.jpg"},
					},
				},
			},
		}
		req2 := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/travel-memories/%d", memory.ID), payloadPhotosSync)
		w2 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w2, req2)
		require.Equal(t, http.StatusOK, w2.Code, w2.Body.String())

		var updatedDay1 dbmodel.TravelMemoryDay
		require.NoError(t, shared.TestDB.First(&updatedDay1, day1.ID).Error)
		assert.Nil(t, updatedDay1.Description)

		var day1Photos []dbmodel.TravelMemoryPhoto
		require.NoError(t, shared.TestDB.Where("memory_day_id = ?", day1.ID).Find(&day1Photos).Error)
		require.Len(t, day1Photos, 1)
		assert.Equal(t, photo1.ID, day1Photos[0].ID)

		var day2Photos []dbmodel.TravelMemoryPhoto
		require.NoError(t, shared.TestDB.Where("memory_day_id = ?", day2.ID).Find(&day2Photos).Error)
		assert.Empty(t, day2Photos)
		var deletedPhoto2 dbmodel.TravelMemoryPhoto
		assert.True(t, errors.Is(shared.TestDB.First(&deletedPhoto2, photo2.ID).Error, gorm.ErrRecordNotFound))

		// 3. Empty Days slice: deletes all existing Days and their Photos
		emptyDays := []api.UpdateTravelMemoryDayInput{}
		payloadEmptyDays := api.UpdateTravelMemoryRequest{
			LocationId: loc.ID,
			Title:      "Clear All Days",
			StartDate:  openapi_types.Date{Time: start},
			EndDate:    openapi_types.Date{Time: end},
			Days:       &emptyDays,
		}
		req3 := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/travel-memories/%d", memory.ID), payloadEmptyDays)
		w3 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w3, req3)
		require.Equal(t, http.StatusOK, w3.Code, w3.Body.String())

		var remainingDays []dbmodel.TravelMemoryDay
		require.NoError(t, shared.TestDB.Where("memory_id = ?", memory.ID).Find(&remainingDays).Error)
		assert.Empty(t, remainingDays)

		var remainingPhotos []dbmodel.TravelMemoryPhoto
		require.NoError(t, shared.TestDB.Find(&remainingPhotos).Error)
		assert.Empty(t, remainingPhotos)
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		start := openapi_types.Date{Time: time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)}
		end := openapi_types.Date{Time: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)}

		payload := api.UpdateTravelMemoryRequest{
			LocationId: loc.ID,
			Title:      "Not found test",
			StartDate:  start,
			EndDate:    end,
		}

		req := shared.NewRequest(t, admin, http.MethodPut, "/travel-memories/999999", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), "Travel memory not found")
	})

	t.Run("Error_InvalidID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		start := openapi_types.Date{Time: time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)}
		end := openapi_types.Date{Time: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)}

		payload := api.UpdateTravelMemoryRequest{
			LocationId: 1,
			Title:      "Invalid ID test",
			StartDate:  start,
			EndDate:    end,
		}

		req := shared.NewRequest(t, admin, http.MethodPut, "/travel-memories/0", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
	})

	t.Run("Error_InvalidDateRange", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		start := time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)

		memory := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			WithDates(start, end).
			MustInsert(shared.TestDB)

		payload := api.UpdateTravelMemoryRequest{
			LocationId: loc.ID,
			Title:      "Invalid range test",
			StartDate:  openapi_types.Date{Time: end},
			EndDate:    openapi_types.Date{Time: start},
		}

		req := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/travel-memories/%d", memory.ID), payload)
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
			MustInsert(shared.TestDB)

		start := time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)

		memory := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			WithDates(start, end).
			MustInsert(shared.TestDB)

		outOfRangeDate := openapi_types.Date{Time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}

		payload := api.UpdateTravelMemoryRequest{
			LocationId: loc.ID,
			Title:      "Day out of range",
			StartDate:  openapi_types.Date{Time: start},
			EndDate:    openapi_types.Date{Time: end},
			Days: &[]api.UpdateTravelMemoryDayInput{
				{Date: outOfRangeDate},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/travel-memories/%d", memory.ID), payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), travelmemoryday.ErrInvalidDayDate.Error())
	})

	t.Run("Error_DuplicateDayDates", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		start := time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)

		memory := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			WithDates(start, end).
			MustInsert(shared.TestDB)

		payload := api.UpdateTravelMemoryRequest{
			LocationId: loc.ID,
			Title:      "Duplicate dates",
			StartDate:  openapi_types.Date{Time: start},
			EndDate:    openapi_types.Date{Time: end},
			Days: &[]api.UpdateTravelMemoryDayInput{
				{Date: openapi_types.Date{Time: start}},
				{Date: openapi_types.Date{Time: start}},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/travel-memories/%d", memory.ID), payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), travelmemoryday.ErrDuplicateDayDate.Error())
	})

	t.Run("Error_InvalidAndNegativeAndOtherMemoryDayID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		start := time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)

		memory := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			WithDates(start, end).
			MustInsert(shared.TestDB)

		otherMemory := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			WithDates(start, end).
			MustInsert(shared.TestDB)

		otherMemoryDay := factory.NewTravelMemoryDayFactory().
			WithMemoryID(otherMemory.ID).
			WithDate(start).
			MustInsert(shared.TestDB)

		for _, invalidID := range []int64{88888, 0, -5, otherMemoryDay.ID} {
			idCopy := invalidID
			payload := api.UpdateTravelMemoryRequest{
				LocationId: loc.ID,
				Title:      "Wrong day ID",
				StartDate:  openapi_types.Date{Time: start},
				EndDate:    openapi_types.Date{Time: end},
				Days: &[]api.UpdateTravelMemoryDayInput{
					{
						Id:   &idCopy,
						Date: openapi_types.Date{Time: start},
					},
				},
			}

			req := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/travel-memories/%d", memory.ID), payload)
			w := httptest.NewRecorder()
			shared.TestRouter.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
			assert.Contains(t, w.Body.String(), travelmemoryusecase.ErrUpdateTravelMemoryInvalidDayID.Error())
		}
	})

	t.Run("Error_DuplicateDayID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		start := time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)

		memory := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			WithDates(start, end).
			MustInsert(shared.TestDB)

		day1 := factory.NewTravelMemoryDayFactory().
			WithMemoryID(memory.ID).
			WithDate(start).
			MustInsert(shared.TestDB)

		payload := api.UpdateTravelMemoryRequest{
			LocationId: loc.ID,
			Title:      "Duplicate day ID",
			StartDate:  openapi_types.Date{Time: start},
			EndDate:    openapi_types.Date{Time: end},
			Days: &[]api.UpdateTravelMemoryDayInput{
				{
					Id:   &day1.ID,
					Date: openapi_types.Date{Time: start},
				},
				{
					Id:   &day1.ID,
					Date: openapi_types.Date{Time: start.AddDate(0, 0, 1)},
				},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/travel-memories/%d", memory.ID), payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), travelmemoryusecase.ErrUpdateTravelMemoryDuplicateDayID.Error())
	})

	t.Run("Error_InvalidAndNegativeAndOtherDayPhotoID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		start := time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)

		memory := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			WithDates(start, end).
			MustInsert(shared.TestDB)

		day1 := factory.NewTravelMemoryDayFactory().
			WithMemoryID(memory.ID).
			WithDate(start).
			MustInsert(shared.TestDB)

		day2 := factory.NewTravelMemoryDayFactory().
			WithMemoryID(memory.ID).
			WithDate(start.AddDate(0, 0, 1)).
			MustInsert(shared.TestDB)

		photoOfDay2 := factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day2.ID).
			WithImagePath("travel_memories/day2-photo.jpg").
			MustInsert(shared.TestDB)

		for _, invalidPhotoID := range []int64{99999, 0, -10, photoOfDay2.ID} {
			idCopy := invalidPhotoID
			payload := api.UpdateTravelMemoryRequest{
				LocationId: loc.ID,
				Title:      "Wrong photo ID",
				StartDate:  openapi_types.Date{Time: start},
				EndDate:    openapi_types.Date{Time: end},
				Days: &[]api.UpdateTravelMemoryDayInput{
					{
						Id:   &day1.ID,
						Date: openapi_types.Date{Time: start},
						Photos: &[]api.UpdateTravelMemoryPhotoInput{
							{
								Id:        &idCopy,
								ImagePath: "travel_memories/test.jpg",
							},
						},
					},
				},
			}

			req := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/travel-memories/%d", memory.ID), payload)
			w := httptest.NewRecorder()
			shared.TestRouter.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
			assert.Contains(t, w.Body.String(), travelmemoryusecase.ErrUpdateTravelMemoryInvalidPhotoID.Error())
		}
	})

	t.Run("Error_DuplicatePhotoID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		start := time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)

		memory := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			WithDates(start, end).
			MustInsert(shared.TestDB)

		day1 := factory.NewTravelMemoryDayFactory().
			WithMemoryID(memory.ID).
			WithDate(start).
			MustInsert(shared.TestDB)

		photo1 := factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day1.ID).
			WithImagePath("travel_memories/photo1.jpg").
			MustInsert(shared.TestDB)

		payload := api.UpdateTravelMemoryRequest{
			LocationId: loc.ID,
			Title:      "Duplicate photo ID",
			StartDate:  openapi_types.Date{Time: start},
			EndDate:    openapi_types.Date{Time: end},
			Days: &[]api.UpdateTravelMemoryDayInput{
				{
					Id:   &day1.ID,
					Date: openapi_types.Date{Time: start},
					Photos: &[]api.UpdateTravelMemoryPhotoInput{
						{
							Id:        &photo1.ID,
							ImagePath: "travel_memories/photo1.jpg",
						},
						{
							Id:        &photo1.ID,
							ImagePath: "travel_memories/photo2.jpg",
						},
					},
				},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/travel-memories/%d", memory.ID), payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), travelmemoryusecase.ErrUpdateTravelMemoryDuplicatePhotoID.Error())
	})

	t.Run("Error_InvalidImagePath_AndVerifyRollback", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		start := time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)

		memory := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			WithTitle("Original Memory Title").
			WithDates(start, end).
			MustInsert(shared.TestDB)

		day1 := factory.NewTravelMemoryDayFactory().
			WithMemoryID(memory.ID).
			WithDate(start).
			MustInsert(shared.TestDB)

		day2ToKeepOnRollback := factory.NewTravelMemoryDayFactory().
			WithMemoryID(memory.ID).
			WithDate(start.AddDate(0, 0, 1)).
			MustInsert(shared.TestDB)

		payload := api.UpdateTravelMemoryRequest{
			LocationId: loc.ID,
			Title:      "Should Not Be Saved Due To Rollback",
			StartDate:  openapi_types.Date{Time: start},
			EndDate:    openapi_types.Date{Time: end},
			Days: &[]api.UpdateTravelMemoryDayInput{
				{
					Id:   &day1.ID,
					Date: openapi_types.Date{Time: start},
				},
				{
					Date: openapi_types.Date{Time: end},
					Photos: &[]api.UpdateTravelMemoryPhotoInput{
						{
							ImagePath: "wrong_folder/photo.jpg",
						},
					},
				},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/travel-memories/%d", memory.ID), payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), travelmemoryphoto.ErrInvalidImagePath.Error())

		// Verify DB state remained completely intact (no partial updates, deletions, or new Day insertions)
		var dbMem dbmodel.TravelMemory
		require.NoError(t, shared.TestDB.First(&dbMem, memory.ID).Error)
		assert.Equal(t, "Original Memory Title", dbMem.Title)

		var dbDay2 dbmodel.TravelMemoryDay
		require.NoError(t, shared.TestDB.First(&dbDay2, day2ToKeepOnRollback.ID).Error)

		var dbDays []dbmodel.TravelMemoryDay
		require.NoError(t, shared.TestDB.Where("memory_id = ?", memory.ID).Find(&dbDays).Error)
		assert.Len(t, dbDays, 2)
	})

	t.Run("Error_InvalidBody", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodPut, "/travel-memories/1", strings.NewReader("invalid-json"))
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		start := openapi_types.Date{Time: time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)}
		end := openapi_types.Date{Time: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)}

		payload := api.UpdateTravelMemoryRequest{
			LocationId: 1,
			Title:      "Unauthorized test",
			StartDate:  start,
			EndDate:    end,
		}

		req := shared.NewUnAuthRequest(t, http.MethodPut, "/travel-memories/1", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, w.Body.String())
	})
}

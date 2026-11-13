package travel_memory_controller_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/internal/infrastructure/dbmodel"
	travelmemoryusecase "platform-go/internal/usecases/travel_memories/usecase"
	"platform-go/tests/factory"
)

func TestTravelMemoryController_DeleteTravelMemory(t *testing.T) {
	t.Run("Normal_DeleteTravelMemorySuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			MustInsert(shared.TestDB)

		start := time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC)

		memory := factory.NewTravelMemoryFactory().
			WithLocationID(loc.ID).
			WithTitle("Chuyến đi Đà Lạt").
			WithDates(start, end).
			MustInsert(shared.TestDB)

		day1 := factory.NewTravelMemoryDayFactory().
			WithMemoryID(memory.ID).
			WithDate(start).
			WithDescription("Ngày 1").
			MustInsert(shared.TestDB)

		day2 := factory.NewTravelMemoryDayFactory().
			WithMemoryID(memory.ID).
			WithDate(start.AddDate(0, 0, 1)).
			WithDescription("Ngày 2").
			MustInsert(shared.TestDB)

		photo1 := factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day1.ID).
			WithImagePath("travel_memories/photo1.jpg").
			WithSortOrder(1).
			MustInsert(shared.TestDB)

		photo2 := factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day1.ID).
			WithImagePath("travel_memories/photo2.jpg").
			WithSortOrder(2).
			MustInsert(shared.TestDB)

		photo3 := factory.NewTravelMemoryPhotoFactory().
			WithMemoryDayID(day2.ID).
			WithImagePath("travel_memories/photo3.jpg").
			WithSortOrder(1).
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodDelete, fmt.Sprintf("/travel-memories/%d", memory.ID), nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.MessageOnlyResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, travelmemoryusecase.DeleteTravelMemoryUseCaseSuccessMessage, res.Message)

		// Verify database state: memory, days, photos deleted
		var deletedMemory dbmodel.TravelMemory
		err := shared.TestDB.First(&deletedMemory, memory.ID).Error
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "memory should be soft deleted")

		var deletedDay dbmodel.TravelMemoryDay
		assert.True(t, errors.Is(shared.TestDB.First(&deletedDay, day1.ID).Error, gorm.ErrRecordNotFound))
		assert.True(t, errors.Is(shared.TestDB.First(&deletedDay, day2.ID).Error, gorm.ErrRecordNotFound))

		var deletedPhoto dbmodel.TravelMemoryPhoto
		assert.True(t, errors.Is(shared.TestDB.First(&deletedPhoto, photo1.ID).Error, gorm.ErrRecordNotFound))
		assert.True(t, errors.Is(shared.TestDB.First(&deletedPhoto, photo2.ID).Error, gorm.ErrRecordNotFound))
		assert.True(t, errors.Is(shared.TestDB.First(&deletedPhoto, photo3.ID).Error, gorm.ErrRecordNotFound))
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodDelete, "/travel-memories/999999", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), "Travel memory not found")
	})

	t.Run("Error_InvalidID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodDelete, "/travel-memories/0", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		req := shared.NewUnAuthRequest(t, http.MethodDelete, "/travel-memories/1", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, w.Body.String())
	})
}


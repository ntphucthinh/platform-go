package itinerary_controller_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/internal/infrastructure/dbmodel"
	itineraryusecase "platform-go/internal/usecases/itineraries/usecase"
	"platform-go/tests/factory"
)

func TestItineraryController_DeleteItinerary(t *testing.T) {
	t.Run("Normal_DeleteItinerarySuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Lâm Đồng").
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			MustInsert(shared.TestDB)

		theme1 := factory.NewTravelThemeFactory().
			WithName("Thiên nhiên").
			MustInsert(shared.TestDB)

		theme2 := factory.NewTravelThemeFactory().
			WithName("Ẩm thực").
			MustInsert(shared.TestDB)

		it := factory.NewItineraryFactory().
			WithLocationID(&loc.ID).
			WithName("Khám phá Đà Lạt").
			WithCoverImagePath("itineraries/dalat-cover.jpg").
			MustInsert(shared.TestDB)

		itTheme1 := factory.NewItineraryThemeFactory().
			WithItineraryID(it.ID).
			WithThemeID(theme1.ID).
			MustInsert(shared.TestDB)
		itTheme2 := factory.NewItineraryThemeFactory().
			WithItineraryID(it.ID).
			WithThemeID(theme2.ID).
			MustInsert(shared.TestDB)

		day1 := factory.NewItineraryDayFactory().
			WithItineraryID(it.ID).
			WithDayNumber(1).
			WithTitle("Ngày 1").
			MustInsert(shared.TestDB)

		day2 := factory.NewItineraryDayFactory().
			WithItineraryID(it.ID).
			WithDayNumber(2).
			WithTitle("Ngày 2").
			MustInsert(shared.TestDB)

		act1 := factory.NewItineraryActivityFactory().
			WithItineraryDayID(day1.ID).
			WithStartTime("07:30:00").
			WithTitle("Ăn sáng").
			MustInsert(shared.TestDB)

		act2 := factory.NewItineraryActivityFactory().
			WithItineraryDayID(day2.ID).
			WithStartTime("09:00:00").
			WithTitle("Tham quan").
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodDelete, fmt.Sprintf("/itineraries/%d", it.ID), nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.MessageOnlyResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, itineraryusecase.DeleteItineraryUseCaseSuccessMessage, res.Message)

		// Verify database state: itinerary soft-deleted, themes, days, activities deleted
		var deletedItinerary dbmodel.Itinerary
		err := shared.TestDB.First(&deletedItinerary, it.ID).Error
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "itinerary should be soft deleted")

		var deletedTheme dbmodel.ItineraryTheme
		assert.True(t, errors.Is(shared.TestDB.First(&deletedTheme, itTheme1.ID).Error, gorm.ErrRecordNotFound))
		assert.True(t, errors.Is(shared.TestDB.First(&deletedTheme, itTheme2.ID).Error, gorm.ErrRecordNotFound))

		var deletedDay dbmodel.ItineraryDay
		assert.True(t, errors.Is(shared.TestDB.First(&deletedDay, day1.ID).Error, gorm.ErrRecordNotFound))
		assert.True(t, errors.Is(shared.TestDB.First(&deletedDay, day2.ID).Error, gorm.ErrRecordNotFound))

		var deletedAct dbmodel.ItineraryActivity
		assert.True(t, errors.Is(shared.TestDB.First(&deletedAct, act1.ID).Error, gorm.ErrRecordNotFound))
		assert.True(t, errors.Is(shared.TestDB.First(&deletedAct, act2.ID).Error, gorm.ErrRecordNotFound))
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodDelete, "/itineraries/999999", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), "Itinerary not found")
	})

	t.Run("Error_InvalidID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodDelete, "/itineraries/0", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		req := shared.NewUnAuthRequest(t, http.MethodDelete, "/itineraries/1", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, w.Body.String())
	})
}

package travel_location_controller_test

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
	"platform-go/internal/infrastructure/dbmodel"
	travellocationusecase "platform-go/internal/usecases/travel_locations/usecase"
	"platform-go/tests/factory"
)

func TestTravelLocationController_DeleteTravelLocation(t *testing.T) {
	t.Run("Normal_DeleteTravelLocationSuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			WithLocationName("Đà Lạt").
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodDelete, fmt.Sprintf("/travel-locations/%d", loc.ID), nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.MessageOnlyResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, travellocationusecase.DeleteTravelLocationUseCaseSuccessMessage, res.Message)

		var deleted dbmodel.TravelLocation
		err := shared.TestDB.First(&deleted, loc.ID).Error
		assert.Error(t, err, "record should be deleted")
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodDelete, "/travel-locations/999999", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Error_InvalidID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodDelete, "/travel-locations/not-a-number", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		req := shared.NewUnAuthRequest(t, http.MethodDelete, "/travel-locations/1", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

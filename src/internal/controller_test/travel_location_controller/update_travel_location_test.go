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
	travellocation "platform-go/internal/domain/travel_location"
	"platform-go/internal/infrastructure/dbmodel"
	travellocationusecase "platform-go/internal/usecases/travel_locations/usecase"
	"platform-go/tests/factory"
)

func TestTravelLocationController_UpdateTravelLocation(t *testing.T) {
	t.Run("Normal_UpdateTravelLocationSuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			WithLocationName("Đà Lạt").
			MustInsert(shared.TestDB)

		lat := 10.776889
		lng := 106.700806
		imgPath := "travel_location/550e8400-e29b-41d4-a716-446655440000.png"
		desc := "Updated description"
		payload := api.UpdateTravelLocationRequest{
			CountryCode:  "VN",
			CountryName:  "Vietnam",
			LocationName: "TP. Hồ Chí Minh",
			Latitude:     &lat,
			Longitude:    &lng,
			ImagePath:    &imgPath,
			Status:       api.TravelLocationStatusVisited,
			Description:  &desc,
		}

		req := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/travel-locations/%d", loc.ID), payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.MessageOnlyResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, travellocationusecase.UpdateTravelLocationUseCaseSuccessMessage, res.Message)

		var updated dbmodel.TravelLocation
		require.NoError(t, shared.TestDB.First(&updated, loc.ID).Error)
		assert.Equal(t, "TP. Hồ Chí Minh", updated.LocationName)
		assert.Equal(t, desc, *updated.Description)
	})

	t.Run("Error_InvalidCountryCode", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().MustInsert(shared.TestDB)

		payload := api.UpdateTravelLocationRequest{
			CountryCode:  "ZZ",
			CountryName:  "Invalid Country",
			LocationName: "Test",
			Status:       api.TravelLocationStatusVisited,
		}

		req := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/travel-locations/%d", loc.ID), payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), travellocation.ErrInvalidCountryCode.Error())
	})

	t.Run("Error_InvalidImagePath", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().MustInsert(shared.TestDB)

		badImgPath := "homestays/550e8400-e29b-41d4-a716-446655440000.png"
		payload := api.UpdateTravelLocationRequest{
			CountryCode:  "VN",
			CountryName:  "Vietnam",
			LocationName: "TP. Hồ Chí Minh",
			ImagePath:    &badImgPath,
			Status:       api.TravelLocationStatusVisited,
		}

		req := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/travel-locations/%d", loc.ID), payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), travellocation.ErrInvalidImagePath.Error())
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		payload := api.UpdateTravelLocationRequest{
			CountryCode:  "VN",
			CountryName:  "Vietnam",
			LocationName: "TP. Hồ Chí Minh",
			Status:       api.TravelLocationStatusVisited,
		}

		req := shared.NewRequest(t, admin, http.MethodPut, "/travel-locations/999999", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, w.Body.String())
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		payload := api.UpdateTravelLocationRequest{
			CountryCode:  "VN",
			CountryName:  "Vietnam",
			LocationName: "TP. Hồ Chí Minh",
			Status:       api.TravelLocationStatusVisited,
		}

		req := shared.NewUnAuthRequest(t, http.MethodPut, "/travel-locations/1", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

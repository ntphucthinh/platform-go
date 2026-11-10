package travel_location_controller_test

import (
	"encoding/json"
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
)

func TestTravelLocationController_CreateTravelLocation(t *testing.T) {
	t.Run("Normal_CreateTravelLocationSuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		lat := 11.9404192
		lng := 108.4583132
		imgPath := "travel_location/550e8400-e29b-41d4-a716-446655440000.png"
		desc := "City of Eternal Spring"
		payload := api.CreateTravelLocationRequest{
			CountryCode:  "VN",
			CountryName:  "Vietnam",
			LocationName: "Đà Lạt",
			Latitude:     &lat,
			Longitude:    &lng,
			ImagePath:    &imgPath,
			Status:       api.TravelLocationStatusVisited,
			Description:  &desc,
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/travel-locations", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.MessageOnlyResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, travellocationusecase.CreateTravelLocationUsecaseSuccessMessage, res.Message)

		var loc dbmodel.TravelLocation
		require.NoError(t, shared.TestDB.Where("location_name = ?", "Đà Lạt").First(&loc).Error)
		assert.Equal(t, "VN", loc.CountryCode)
		assert.Equal(t, "Vietnam", loc.CountryName)
		assert.Equal(t, "visited", loc.Status)
	})

	t.Run("Error_InvalidCountryCode", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		payload := api.CreateTravelLocationRequest{
			CountryCode:  "ZZ",
			CountryName:  "Invalid Country",
			LocationName: "Test",
			Status:       api.TravelLocationStatusVisited,
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/travel-locations", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), travellocation.ErrInvalidCountryCode.Error())
	})

	t.Run("Error_InvalidImagePath", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		badImgPath := "homestays/550e8400-e29b-41d4-a716-446655440000.png"
		payload := api.CreateTravelLocationRequest{
			CountryCode:  "VN",
			CountryName:  "Vietnam",
			LocationName: "Đà Lạt",
			ImagePath:    &badImgPath,
			Status:       api.TravelLocationStatusVisited,
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/travel-locations", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), travellocation.ErrInvalidImagePath.Error())
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		payload := api.CreateTravelLocationRequest{
			CountryCode:  "VN",
			CountryName:  "Vietnam",
			LocationName: "Đà Lạt",
			Status:       api.TravelLocationStatusVisited,
		}

		req := shared.NewUnAuthRequest(t, http.MethodPost, "/travel-locations", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

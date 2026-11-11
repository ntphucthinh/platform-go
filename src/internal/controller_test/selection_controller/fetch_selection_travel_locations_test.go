package selection_controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/tests/factory"
)

func TestSelectionController_FetchSelectionTravelLocations(t *testing.T) {
	t.Run("Normal_GetTravelLocationsList", func(t *testing.T) {
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

		req := shared.NewRequest(t, admin, http.MethodGet, "/selections/travel-locations", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.SelectionTravelLocationResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))

		require.Len(t, res.Data, 2)
		assert.Equal(t, loc1.ID, res.Data[0].Id)
		assert.Equal(t, loc1.LocationName, res.Data[0].Name)
		assert.Equal(t, "VN", res.Data[0].CountryCode)
		assert.Equal(t, loc2.ID, res.Data[1].Id)
		assert.Equal(t, loc2.LocationName, res.Data[1].Name)
		assert.Equal(t, "JP", res.Data[1].CountryCode)

		rawJSON := w.Body.String()
		assert.Contains(t, rawJSON, `"countryCode"`)
	})

	t.Run("Normal_FilterByCountryCode", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		locVN := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		_ = factory.NewTravelLocationFactory().
			WithLocationName("Tokyo").
			WithCountryCode("JP").
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, "/selections/travel-locations?countryCode=VN", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.SelectionTravelLocationResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))

		require.Len(t, res.Data, 1)
		assert.Equal(t, locVN.ID, res.Data[0].Id)
		assert.Equal(t, locVN.LocationName, res.Data[0].Name)
	})

	t.Run("Normal_FilterByCountryCode_SnakeCase", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		_ = factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		locJP := factory.NewTravelLocationFactory().
			WithLocationName("Tokyo").
			WithCountryCode("JP").
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, "/selections/travel-locations?country_code=jp", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.SelectionTravelLocationResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))

		require.Len(t, res.Data, 1)
		assert.Equal(t, locJP.ID, res.Data[0].Id)
		assert.Equal(t, locJP.LocationName, res.Data[0].Name)
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		req := shared.NewUnAuthRequest(t, http.MethodGet, "/selections/travel-locations", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

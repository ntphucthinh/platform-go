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
	"platform-go/tests/factory"
)

func TestTravelLocationController_FetchTravelLocations(t *testing.T) {
	t.Run("Normal_GetTravelLocationsList", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		lat := 11.9404192
		lng := 108.4583132
		img := "/images/locations/da_lat.jpg"
		desc := "City of Eternal Spring"
		notes := "Visited Xuan Huong Lake"
		mapsURL := "https://maps.google.com/?q=Da+Lat"

		loc := factory.NewTravelLocationFactory().
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			WithLocationName("Đà Lạt").
			WithLatitude(lat).
			WithLongitude(lng).
			WithImagePath(img).
			WithVisited().
			WithDescription(desc).
			WithNotes(notes).
			WithGoogleMapsURL(mapsURL).
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, "/travel-locations?status=visited&locationName=Đà+Lạt&countryCode=VN", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var body api.FetchTravelLocationsResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

		require.Len(t, body.Data, 1)
		assert.Equal(t, int64(1), body.Total)
		assert.Equal(t, loc.LocationName, body.Data[0].LocationName)
		assert.Equal(t, loc.CountryCode, body.Data[0].CountryCode)
		assert.Equal(t, int64(0), body.Data[0].VisitCount)
		assert.Equal(t, int64(0), body.Data[0].ItineraryCount)

		rawJSON := w.Body.String()
		for _, camelField := range []string{
			`"countryCode"`,
			`"countryName"`,
			`"locationName"`,
			`"imagePath"`,
			`"googleMapsUrl"`,
			`"visitCount"`,
			`"itineraryCount"`,
			`"total"`,
		} {
			assert.Contains(t, rawJSON, camelField)
		}
	})

	t.Run("Normal_CountryCodeQuerySnakeCase", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithCountryCode("VN").
			WithLocationName("Hà Nội").
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, "/travel-locations?country_code=vn", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var body api.FetchTravelLocationsResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		require.Len(t, body.Data, 1)
		assert.Equal(t, loc.LocationName, body.Data[0].LocationName)
	})

	t.Run("Error_InvalidStatus", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodGet, "/travel-locations?status=invalid", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		req := shared.NewUnAuthRequest(t, http.MethodGet, "/travel-locations", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

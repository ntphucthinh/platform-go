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

		req := shared.NewRequest(t, admin, http.MethodGet, "/travel-locations?status=visited&countryCode=VN", nil)
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

	t.Run("Normal_FilterByCountryCodeAndStatusOnly", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		factory.NewTravelLocationFactory().
			WithCountryCode("VN").
			WithLocationName("Hà Nội").
			WithVisited().
			MustInsert(shared.TestDB)

		factory.NewTravelLocationFactory().
			WithCountryCode("VN").
			WithLocationName("Đà Nẵng").
			WithWishlist().
			MustInsert(shared.TestDB)

		factory.NewTravelLocationFactory().
			WithCountryCode("JP").
			WithLocationName("Tokyo").
			WithVisited().
			MustInsert(shared.TestDB)

		// Filter by status=visited: should return VN (Hà Nội) and JP (Tokyo)
		reqVisited := shared.NewRequest(t, admin, http.MethodGet, "/travel-locations?status=visited", nil)
		wVisited := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(wVisited, reqVisited)
		require.Equal(t, http.StatusOK, wVisited.Code)
		var bodyVisited api.FetchTravelLocationsResponse
		require.NoError(t, json.Unmarshal(wVisited.Body.Bytes(), &bodyVisited))
		assert.Equal(t, int64(2), bodyVisited.Total)

		// Filter by countryCode=VN: should return Hà Nội and Đà Nẵng
		reqVN := shared.NewRequest(t, admin, http.MethodGet, "/travel-locations?countryCode=VN", nil)
		wVN := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(wVN, reqVN)
		require.Equal(t, http.StatusOK, wVN.Code)
		var bodyVN api.FetchTravelLocationsResponse
		require.NoError(t, json.Unmarshal(wVN.Body.Bytes(), &bodyVN))
		assert.Equal(t, int64(2), bodyVN.Total)

		// Filter by both countryCode=VN and status=wishlist: should return only Đà Nẵng
		reqVNWishlist := shared.NewRequest(t, admin, http.MethodGet, "/travel-locations?countryCode=VN&status=wishlist", nil)
		wVNWishlist := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(wVNWishlist, reqVNWishlist)
		require.Equal(t, http.StatusOK, wVNWishlist.Code)
		var bodyVNWishlist api.FetchTravelLocationsResponse
		require.NoError(t, json.Unmarshal(wVNWishlist.Body.Bytes(), &bodyVNWishlist))
		require.Len(t, bodyVNWishlist.Data, 1)
		assert.Equal(t, "Đà Nẵng", bodyVNWishlist.Data[0].LocationName)

		// Passing locationName query parameter should not filter by location name anymore
		reqWithIgnoredParam := shared.NewRequest(t, admin, http.MethodGet, "/travel-locations?countryCode=VN&status=visited&locationName=NonExistent", nil)
		wWithIgnoredParam := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(wWithIgnoredParam, reqWithIgnoredParam)
		require.Equal(t, http.StatusOK, wWithIgnoredParam.Code)
		var bodyIgnored api.FetchTravelLocationsResponse
		require.NoError(t, json.Unmarshal(wWithIgnoredParam.Body.Bytes(), &bodyIgnored))
		require.Len(t, bodyIgnored.Data, 1)
		assert.Equal(t, "Hà Nội", bodyIgnored.Data[0].LocationName)
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

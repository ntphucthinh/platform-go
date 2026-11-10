package travel_location_controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/tests/factory"
)

func TestTravelLocationController_FetchTravelLocationByID(t *testing.T) {
	t.Run("Normal_GetTravelLocationByIDSuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		lat := 10.776889
		lng := 106.700806
		img := "/images/locations/hcmc.jpg"
		desc := "Dynamic economic hub"
		notes := "Visited Ben Thanh Market"
		mapsURL := "https://maps.google.com/?q=Ho+Chi+Minh+City"

		loc := factory.NewTravelLocationFactory().
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			WithLocationName("TP. Hồ Chí Minh").
			WithLatitude(lat).
			WithLongitude(lng).
			WithImagePath(img).
			WithVisited().
			WithDescription(desc).
			WithNotes(notes).
			WithGoogleMapsURL(mapsURL).
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, fmt.Sprintf("/travel-locations/%d", loc.ID), nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.FetchTravelLocationByIDResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))

		assert.Equal(t, int64(loc.ID), res.Data.Id)
		assert.Equal(t, "TP. Hồ Chí Minh", res.Data.LocationName)
		assert.Equal(t, int64(0), res.Data.VisitCount)
		assert.Equal(t, int64(0), res.Data.ItineraryCount)

		rawJSON := w.Body.String()
		for _, camelField := range []string{
			`"countryCode"`,
			`"countryName"`,
			`"locationName"`,
			`"imagePath"`,
			`"googleMapsUrl"`,
			`"visitCount"`,
			`"itineraryCount"`,
		} {
			assert.True(t, strings.Contains(rawJSON, camelField), "expected camelField %s in response", camelField)
		}
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodGet, "/travel-locations/999999", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Error_InvalidID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodGet, "/travel-locations/not-a-number", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

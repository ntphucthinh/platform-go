package travel_location_controller_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	api "platform-go/generated/api"
	"platform-go/internal/infrastructure/qsdto"
)

func TestFetchTravelLocationByID_Success(t *testing.T) {
	lat := 10.776889
	lng := 106.700806
	img := "/images/locations/hcmc.jpg"
	desc := "Dynamic economic hub"
	notes := "Visited Ben Thanh Market"
	mapsURL := "https://maps.google.com/?q=Ho+Chi+Minh+City"

	qs := &travelLocationQueryService{
		location: &qsdto.TravelLocationDto{
			ID:            2,
			CountryCode:   "VN",
			CountryName:   "Vietnam",
			StateCode:     "SG",
			Name:          "TP. Hồ Chí Minh",
			Latitude:      &lat,
			Longitude:     &lng,
			ImagePath:     &img,
			Status:        "visited",
			Description:   &desc,
			Notes:         &notes,
			GoogleMapsURL: &mapsURL,
		},
	}

	router := newTravelLocationRouter(t, qs)
	request := httptest.NewRequest(http.MethodGet, "/travel-locations/2", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", response.Code, response.Body.String())
	}

	var res api.FetchTravelLocationByIDResponse
	if err := json.Unmarshal(response.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Data.Id != 2 || res.Data.Name != "TP. Hồ Chí Minh" || res.Data.StateCode != "SG" {
		t.Fatalf("unexpected travel location data: %#v", res.Data)
	}
	if res.Data.VisitCount != 0 || res.Data.ItineraryCount != 0 {
		t.Fatalf("expected visitCount=0 and itineraryCount=0, got %d and %d", res.Data.VisitCount, res.Data.ItineraryCount)
	}

	rawJSON := response.Body.String()
	for _, camelField := range []string{
		`"countryCode"`,
		`"countryName"`,
		`"stateCode"`,
		`"imagePath"`,
		`"googleMapsUrl"`,
		`"visitCount"`,
		`"itineraryCount"`,
	} {
		if !strings.Contains(rawJSON, camelField) {
			t.Fatalf("expected JSON response to contain camelCase field %s, got: %s", camelField, rawJSON)
		}
	}
}

func TestFetchTravelLocationByID_NotFound(t *testing.T) {
	qs := &travelLocationQueryService{location: nil}
	router := newTravelLocationRouter(t, qs)
	request := httptest.NewRequest(http.MethodGet, "/travel-locations/999", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestFetchTravelLocationByID_InvalidID(t *testing.T) {
	qs := &travelLocationQueryService{}
	router := newTravelLocationRouter(t, qs)
	request := httptest.NewRequest(http.MethodGet, "/travel-locations/not-a-number", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestFetchTravelLocationByID_DatabaseError(t *testing.T) {
	qs := &travelLocationQueryService{err: errors.New("db query failed")}
	router := newTravelLocationRouter(t, qs)
	request := httptest.NewRequest(http.MethodGet, "/travel-locations/1", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d, body: %s", response.Code, response.Body.String())
	}
}

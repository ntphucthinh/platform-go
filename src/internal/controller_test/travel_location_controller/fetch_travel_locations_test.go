package travel_location_controller_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/internal/infrastructure/qsdto"
	"platform-go/internal/interfaces/controller"
	travellocationqueryservice "platform-go/internal/usecases/travel_locations/queryservice"
	travellocationusecase "platform-go/internal/usecases/travel_locations/usecase"
)

type travelLocationQueryService struct {
	locations  []qsdto.TravelLocationDto
	location   *qsdto.TravelLocationDto
	lastFilter qsdto.TravelLocationFilter
	err        error
}

func (s *travelLocationQueryService) FetchTravelLocations(ctx context.Context, filter qsdto.TravelLocationFilter) ([]qsdto.TravelLocationDto, int64, error) {
	s.lastFilter = filter
	if s.err != nil {
		return nil, 0, s.err
	}
	return s.locations, int64(len(s.locations)), nil
}

func (s *travelLocationQueryService) FetchTravelLocationByID(ctx context.Context, id int64) (*qsdto.TravelLocationDto, error) {
	return s.location, s.err
}

var _ travellocationqueryservice.TravelLocationQueryService = (*travelLocationQueryService)(nil)

func newTravelLocationRouter(t *testing.T, queryService travellocationqueryservice.TravelLocationQueryService) http.Handler {
	t.Helper()
	fetchUsecase := travellocationusecase.NewFetchTravelLocationsUsecase(queryService)
	createUsecase := travellocationusecase.NewCreateTravelLocationUsecase(nil, nil)
	fetchByIDUsecase := travellocationusecase.NewFetchTravelLocationByIDUsecase(queryService)
	updateUsecase := travellocationusecase.NewUpdateTravelLocationUsecase(nil, nil)
	deleteUsecase := travellocationusecase.NewDeleteTravelLocationUsecase(nil, nil)
	handler := controller.NewTravelLocationController(fetchUsecase, createUsecase, fetchByIDUsecase, updateUsecase, deleteUsecase)
	serverController := controller.NewServerController(nil, nil, nil, handler, nil, nil)
	router := shared.NewSetupRouter(t)
	api.RegisterHandlers(router, serverController)
	return router
}

func TestFetchTravelLocations(t *testing.T) {
	lat := 11.9404192
	lng := 108.4583132
	img := "/images/locations/da_lat.jpg"
	desc := "City of Eternal Spring"
	notes := "Visited Xuan Huong Lake"
	mapsURL := "https://maps.google.com/?q=Da+Lat"

	qs := &travelLocationQueryService{
		locations: []qsdto.TravelLocationDto{
			{
				ID:            1,
				CountryCode:   "VN",
				CountryName:   "Vietnam",
				LocationName:  "Đà Lạt",
				Latitude:      &lat,
				Longitude:     &lng,
				ImagePath:     &img,
				Status:        "visited",
				Description:   &desc,
				Notes:         &notes,
				GoogleMapsURL: &mapsURL,
			},
		},
	}

	router := newTravelLocationRouter(t, qs)
	request := httptest.NewRequest(http.MethodGet, "/travel-locations?status=visited&locationName=Đà+Lạt&countryCode=VN", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}

	var body api.FetchTravelLocationsResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Data) != 1 || body.Total != 1 {
		t.Fatalf("response = %#v", body)
	}
	if body.Data[0].LocationName != "Đà Lạt" || body.Data[0].CountryCode != "VN" {
		t.Fatalf("unexpected location: %#v", body.Data[0])
	}
	if body.Data[0].VisitCount != 0 || body.Data[0].ItineraryCount != 0 {
		t.Fatalf("expected visitCount=0 and itineraryCount=0, got %d and %d", body.Data[0].VisitCount, body.Data[0].ItineraryCount)
	}

	rawJSON := response.Body.String()
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
		if !strings.Contains(rawJSON, camelField) {
			t.Fatalf("expected JSON response to contain camelCase field %s, got: %s", camelField, rawJSON)
		}
	}

	if qs.lastFilter.CountryCode == nil || *qs.lastFilter.CountryCode != "VN" {
		t.Fatalf("expected lastFilter.CountryCode = VN, got: %v", qs.lastFilter.CountryCode)
	}
	if qs.lastFilter.LocationName == nil || *qs.lastFilter.LocationName != "Đà Lạt" {
		t.Fatalf("expected lastFilter.LocationName = Đà Lạt, got: %v", qs.lastFilter.LocationName)
	}
}

func TestFetchTravelLocations_CountryCodeQuery(t *testing.T) {
	qs := &travelLocationQueryService{}
	router := newTravelLocationRouter(t, qs)
	request := httptest.NewRequest(http.MethodGet, "/travel-locations?country_code=vn", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if qs.lastFilter.CountryCode == nil || *qs.lastFilter.CountryCode != "vn" {
		t.Fatalf("expected lastFilter.CountryCode = vn, got: %v", qs.lastFilter.CountryCode)
	}
}

func TestFetchTravelLocations_InvalidStatus(t *testing.T) {
	router := newTravelLocationRouter(t, &travelLocationQueryService{})
	request := httptest.NewRequest(http.MethodGet, "/travel-locations?status=invalid", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body = %s", response.Code, response.Body.String())
	}
}

func TestFetchTravelLocations_DatabaseError(t *testing.T) {
	router := newTravelLocationRouter(t, &travelLocationQueryService{err: errors.New("database unavailable")})
	request := httptest.NewRequest(http.MethodGet, "/travel-locations", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d, body = %s", response.Code, response.Body.String())
	}
}

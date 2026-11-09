package travel_location_controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	api "platform-go/generated/api"
	travellocation "platform-go/internal/domain/travel_location"
	travellocationusecase "platform-go/internal/usecases/travel_locations/usecase"
)

func performPost(t *testing.T, router http.Handler, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/travel-locations", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func TestCreateTravelLocation_Success(t *testing.T) {
	repo := &fakeTravelLocationRepository{}
	router := newWriteTravelLocationRouter(t, &fakeDB{}, repo)

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
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	response := performPost(t, router, bodyBytes)
	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", response.Code, response.Body.String())
	}

	var res api.MessageOnlyResponse
	if err := json.Unmarshal(response.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Message != travellocationusecase.CreateTravelLocationUsecaseSuccessMessage {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if repo.created == nil || repo.created.CountryCode != "VN" || repo.created.LocationName != "Đà Lạt" {
		t.Fatalf("expected created entity with CountryCode=VN and LocationName=Đà Lạt, got: %#v", repo.created)
	}
}

func TestCreateTravelLocation_InvalidCountryCode(t *testing.T) {
	repo := &fakeTravelLocationRepository{}
	router := newWriteTravelLocationRouter(t, &fakeDB{}, repo)

	payload := api.CreateTravelLocationRequest{
		CountryCode:  "ZZ",
		CountryName:  "Invalid Country",
		LocationName: "Test",
		Status:       api.TravelLocationStatusVisited,
	}
	bodyBytes, _ := json.Marshal(payload)

	response := performPost(t, router, bodyBytes)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), travellocation.ErrInvalidCountryCode.Error()) {
		t.Fatalf("expected error %q, got: %s", travellocation.ErrInvalidCountryCode.Error(), response.Body.String())
	}
}

func TestCreateTravelLocation_InvalidImagePath(t *testing.T) {
	repo := &fakeTravelLocationRepository{}
	router := newWriteTravelLocationRouter(t, &fakeDB{}, repo)

	badImgPath := "homestays/550e8400-e29b-41d4-a716-446655440000.png"
	payload := api.CreateTravelLocationRequest{
		CountryCode:  "VN",
		CountryName:  "Vietnam",
		LocationName: "Đà Lạt",
		ImagePath:    &badImgPath,
		Status:       api.TravelLocationStatusVisited,
	}
	bodyBytes, _ := json.Marshal(payload)

	response := performPost(t, router, bodyBytes)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), travellocation.ErrInvalidImagePath.Error()) {
		t.Fatalf("expected error %q, got: %s", travellocation.ErrInvalidImagePath.Error(), response.Body.String())
	}
}

func TestCreateTravelLocation_DatabaseError(t *testing.T) {
	repo := &fakeTravelLocationRepository{createErr: errors.New("insert failed")}
	router := newWriteTravelLocationRouter(t, &fakeDB{}, repo)

	payload := api.CreateTravelLocationRequest{
		CountryCode:  "VN",
		CountryName:  "Vietnam",
		LocationName: "Đà Lạt",
		Status:       api.TravelLocationStatusVisited,
	}
	bodyBytes, _ := json.Marshal(payload)

	response := performPost(t, router, bodyBytes)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d, body: %s", response.Code, response.Body.String())
	}
}

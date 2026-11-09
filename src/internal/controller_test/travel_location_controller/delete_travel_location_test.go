package travel_location_controller_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	api "platform-go/generated/api"
	travellocation "platform-go/internal/domain/travel_location"
	travellocationusecase "platform-go/internal/usecases/travel_locations/usecase"
)

func performDelete(t *testing.T, router http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodDelete, path, nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func TestDeleteTravelLocation_Success(t *testing.T) {
	repo := &fakeTravelLocationRepository{
		existing: &travellocation.TravelLocation{
			ID:           1,
			CountryCode:  "VN",
			CountryName:  "Vietnam",
			LocationName: "Đà Lạt",
		},
	}
	router := newWriteTravelLocationRouter(t, &fakeDB{}, repo)

	response := performDelete(t, router, "/travel-locations/1")
	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", response.Code, response.Body.String())
	}

	var res api.MessageOnlyResponse
	if err := json.Unmarshal(response.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Message != travellocationusecase.DeleteTravelLocationUseCaseSuccessMessage {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if repo.deletedID != 1 {
		t.Fatalf("expected deletedID=1, got %d", repo.deletedID)
	}
}

func TestDeleteTravelLocation_NotFound(t *testing.T) {
	repo := &fakeTravelLocationRepository{existing: nil}
	router := newWriteTravelLocationRouter(t, &fakeDB{}, repo)

	response := performDelete(t, router, "/travel-locations/999")
	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestDeleteTravelLocation_InvalidID(t *testing.T) {
	repo := &fakeTravelLocationRepository{}
	router := newWriteTravelLocationRouter(t, &fakeDB{}, repo)

	response := performDelete(t, router, "/travel-locations/not-a-number")
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestDeleteTravelLocation_DatabaseError(t *testing.T) {
	repo := &fakeTravelLocationRepository{
		existing:  &travellocation.TravelLocation{ID: 1},
		deleteErr: errors.New("delete failed"),
	}
	router := newWriteTravelLocationRouter(t, &fakeDB{}, repo)

	response := performDelete(t, router, "/travel-locations/1")
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d, body: %s", response.Code, response.Body.String())
	}
}


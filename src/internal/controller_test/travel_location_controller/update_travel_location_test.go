package travel_location_controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	travellocation "platform-go/internal/domain/travel_location"
	"platform-go/internal/infrastructure/repository"
	"platform-go/internal/interfaces/controller"
	travellocationusecase "platform-go/internal/usecases/travel_locations/usecase"
)

type fakeDB struct {
	txFunc func(ctx context.Context, fn func(tx *gorm.DB) error) error
}

func (f *fakeDB) GetDB() *gorm.DB                         { return nil }
func (f *fakeDB) GetDBByCtx(ctx context.Context) *gorm.DB { return nil }
func (f *fakeDB) Close() error                            { return nil }
func (f *fakeDB) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	if f.txFunc != nil {
		return f.txFunc(ctx, fn)
	}
	return fn(nil)
}

type fakeTravelLocationRepository struct {
	existing  *travellocation.TravelLocation
	created   *travellocation.TravelLocation
	updated   *travellocation.TravelLocation
	deletedID int64
	findErr   error
	createErr error
	updateErr error
	deleteErr error
}

func (r *fakeTravelLocationRepository) FindTravelLocationByID(ctx context.Context, id int64) (*travellocation.TravelLocation, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	return r.existing, nil
}

func (r *fakeTravelLocationRepository) CreateTravelLocation(ctx context.Context, location *travellocation.TravelLocation) error {
	if r.createErr != nil {
		return r.createErr
	}
	copied := *location
	copied.ID = 1
	r.created = &copied
	return nil
}

func (r *fakeTravelLocationRepository) UpdateTravelLocation(ctx context.Context, location *travellocation.TravelLocation) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	copied := *location
	r.updated = &copied
	return nil
}

func (r *fakeTravelLocationRepository) DeleteTravelLocationByID(ctx context.Context, id int64) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}
	r.deletedID = id
	return nil
}

func newWriteTravelLocationRouter(
	t *testing.T,
	db repository.DBInterface,
	repo travellocation.TravelLocationRepository,
) http.Handler {
	t.Helper()
	fetchUsecase := travellocationusecase.NewFetchTravelLocationsUsecase(nil)
	createUsecase := travellocationusecase.NewCreateTravelLocationUsecase(db, repo)
	fetchByIDUsecase := travellocationusecase.NewFetchTravelLocationByIDUsecase(nil)
	updateUsecase := travellocationusecase.NewUpdateTravelLocationUsecase(db, repo)
	deleteUsecase := travellocationusecase.NewDeleteTravelLocationUsecase(db, repo)
	handler := controller.NewTravelLocationController(fetchUsecase, createUsecase, fetchByIDUsecase, updateUsecase, deleteUsecase)
	serverController := controller.NewServerController(nil, nil, nil, handler, nil)
	router := shared.NewSetupRouter(t)
	api.RegisterHandlers(router, serverController)
	return router
}

func performPut(t *testing.T, router http.Handler, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPut, path, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func TestUpdateTravelLocation_Success(t *testing.T) {
	repo := &fakeTravelLocationRepository{
		existing: &travellocation.TravelLocation{
			ID:           1,
			CountryCode:  "VN",
			CountryName:  "Vietnam",
			LocationName: "Đà Lạt",
			Status:       "visited",
		},
	}
	router := newWriteTravelLocationRouter(t, &fakeDB{}, repo)

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
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	response := performPut(t, router, "/travel-locations/1", bodyBytes)
	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", response.Code, response.Body.String())
	}

	var res api.MessageOnlyResponse
	if err := json.Unmarshal(response.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Message != travellocationusecase.UpdateTravelLocationUseCaseSuccessMessage {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if repo.updated == nil || repo.updated.LocationName != "TP. Hồ Chí Minh" {
		t.Fatalf("expected updated entity with LocationName=TP. Hồ Chí Minh, got: %#v", repo.updated)
	}
}

func TestUpdateTravelLocation_InvalidCountryCode(t *testing.T) {
	repo := &fakeTravelLocationRepository{
		existing: &travellocation.TravelLocation{ID: 1, CountryCode: "VN", LocationName: "Đà Lạt"},
	}
	router := newWriteTravelLocationRouter(t, &fakeDB{}, repo)

	payload := api.UpdateTravelLocationRequest{
		CountryCode:  "ZZ",
		CountryName:  "Invalid Country",
		LocationName: "Test",
		Status:       api.TravelLocationStatusVisited,
	}
	bodyBytes, _ := json.Marshal(payload)

	response := performPut(t, router, "/travel-locations/1", bodyBytes)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), travellocation.ErrInvalidCountryCode.Error()) {
		t.Fatalf("expected error %q, got: %s", travellocation.ErrInvalidCountryCode.Error(), response.Body.String())
	}
}

func TestUpdateTravelLocation_InvalidImagePath(t *testing.T) {
	repo := &fakeTravelLocationRepository{
		existing: &travellocation.TravelLocation{ID: 1, CountryCode: "VN", LocationName: "Đà Lạt"},
	}
	router := newWriteTravelLocationRouter(t, &fakeDB{}, repo)

	badImgPath := "homestays/550e8400-e29b-41d4-a716-446655440000.png"
	payload := api.UpdateTravelLocationRequest{
		CountryCode:  "VN",
		CountryName:  "Vietnam",
		LocationName: "TP. Hồ Chí Minh",
		ImagePath:    &badImgPath,
		Status:       api.TravelLocationStatusVisited,
	}
	bodyBytes, _ := json.Marshal(payload)

	response := performPut(t, router, "/travel-locations/1", bodyBytes)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), travellocation.ErrInvalidImagePath.Error()) {
		t.Fatalf("expected error %q, got: %s", travellocation.ErrInvalidImagePath.Error(), response.Body.String())
	}
}

func TestUpdateTravelLocation_NotFound(t *testing.T) {
	repo := &fakeTravelLocationRepository{existing: nil}
	router := newWriteTravelLocationRouter(t, &fakeDB{}, repo)

	payload := api.UpdateTravelLocationRequest{
		CountryCode:  "VN",
		CountryName:  "Vietnam",
		LocationName: "TP. Hồ Chí Minh",
		Status:       api.TravelLocationStatusVisited,
	}
	bodyBytes, _ := json.Marshal(payload)

	response := performPut(t, router, "/travel-locations/999", bodyBytes)
	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestUpdateTravelLocation_DatabaseError(t *testing.T) {
	repo := &fakeTravelLocationRepository{
		findErr: errors.New("db error"),
	}
	router := newWriteTravelLocationRouter(t, &fakeDB{}, repo)

	payload := api.UpdateTravelLocationRequest{
		CountryCode:  "VN",
		CountryName:  "Vietnam",
		LocationName: "TP. Hồ Chí Minh",
		Status:       api.TravelLocationStatusVisited,
	}
	bodyBytes, _ := json.Marshal(payload)

	response := performPut(t, router, "/travel-locations/1", bodyBytes)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d, body: %s", response.Code, response.Body.String())
	}
}

package homestay_controller_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/internal/infrastructure/qsdto"
	"platform-go/internal/interfaces/controller"
	homestayqueryservice "platform-go/internal/usecases/homestays/queryservice"
	homestayusecase "platform-go/internal/usecases/homestays/usecase"
)

func newFetchHomestayByIDRouter(t *testing.T, queryService homestayqueryservice.HomestayQueryService) http.Handler {
	t.Helper()
	fetchUsecase := homestayusecase.NewFetchHomestaysUseCase(queryService)
	createUsecase := homestayusecase.NewCreateHomestayUsecase(nil)
	fetchByIDUsecase := homestayusecase.NewFetchHomestayByIDUseCase(queryService)
	updateUsecase := homestayusecase.NewUpdateHomestayUsecase(nil)
	deleteUsecase := homestayusecase.NewDeleteHomestayUsecase(nil)
	handler := controller.NewHomestayController(fetchUsecase, createUsecase, fetchByIDUsecase, updateUsecase, deleteUsecase)
	serverController := controller.NewServerController(nil, nil, handler, nil)
	router := shared.NewSetupRouter(t)
	api.RegisterHandlers(router, serverController)
	return router
}

func performGet(t *testing.T, router http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func TestFetchHomestayByID_Success(t *testing.T) {
	desc := "A lovely homestay in Da Nang"
	mapsURL := "https://maps.google.com/?q=danang"
	roomDesc := "Deluxe Ocean View"
	price := 1200000.0
	unit := "per_room_per_night"
	roomID := int64(10)

	qs := homestayQueryService{
		homestay: &qsdto.HomestayDto{
			ID:            1,
			Name:          "Da Nang Homestay",
			Address:       "123 Vo Nguyen Giap",
			Description:   &desc,
			GoogleMapsURL: &mapsURL,
			Rooms: []qsdto.HomestayRoomDto{
				{
					ID:          roomID,
					HomestayID:  1,
					Name:        "Room 101",
					Description: &roomDesc,
					PriceAmount: &price,
					PriceUnit:   &unit,
					SortOrder:   1,
					Images: []qsdto.HomestayImageDto{
						{
							ID:        100,
							RoomID:    roomID,
							ImagePath: "homestays/ocean.jpg",
							SortOrder: 1,
							IsCover:   true,
						},
					},
				},
			},
		},
	}

	router := newFetchHomestayByIDRouter(t, qs)
	response := performGet(t, router, "/homestays/1")

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", response.Code, response.Body.String())
	}

	var res api.FetchHomestayByIDResponse
	if err := json.Unmarshal(response.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}

	if res.Data.Id != 1 || res.Data.Name != "Da Nang Homestay" {
		t.Fatalf("unexpected homestay data: %#v", res.Data)
	}
	if len(res.Data.HomestayRooms) != 1 || len(res.Data.HomestayRooms[0].HomestayImages) != 1 {
		t.Fatalf("unexpected rooms or images: %#v", res.Data.HomestayRooms)
	}
	if res.Data.HomestayRooms[0].HomestayImages[0].ImagePath != "homestays/ocean.jpg" {
		t.Fatalf("unexpected image path: %s", res.Data.HomestayRooms[0].HomestayImages[0].ImagePath)
	}
}

func TestFetchHomestayByID_NotFound(t *testing.T) {
	qs := homestayQueryService{homestay: nil}
	router := newFetchHomestayByIDRouter(t, qs)
	response := performGet(t, router, "/homestays/999")

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestFetchHomestayByID_InvalidID(t *testing.T) {
	qs := homestayQueryService{}
	router := newFetchHomestayByIDRouter(t, qs)
	response := performGet(t, router, "/homestays/not-a-number")

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestFetchHomestayByID_DatabaseError(t *testing.T) {
	qs := homestayQueryService{err: errors.New("db query failed")}
	router := newFetchHomestayByIDRouter(t, qs)
	response := performGet(t, router, "/homestays/1")

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d, body: %s", response.Code, response.Body.String())
	}
}

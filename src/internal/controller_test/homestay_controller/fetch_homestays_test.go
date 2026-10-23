package homestay_controller_test

import (
	"context"
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

type homestayQueryService struct {
	homestays []qsdto.HomestayDto
	err       error
}

func (s homestayQueryService) FetchHomestays(context.Context) ([]qsdto.HomestayDto, error) {
	return s.homestays, s.err
}

var _ homestayqueryservice.HomestayQueryService = homestayQueryService{}

func TestFetchHomestays(t *testing.T) {
	description := "Room description"
	price := 10000.0
	priceUnit := "per_room_per_night"
	roomID := int64(11)
	queryService := homestayQueryService{homestays: []qsdto.HomestayDto{{
		ID: 1, Name: "Homestay", Address: "Da Nang",
		Rooms: []qsdto.HomestayRoomDto{{
			ID: 11, HomestayID: 1, Name: "Room 1", Description: &description, PriceAmount: &price, PriceUnit: &priceUnit, SortOrder: 1,
			Images: []qsdto.HomestayImageDto{{ID: 3, RoomID: roomID, ImagePath: "homestays/image.png", SortOrder: 1, IsCover: true}},
		}},
	}}}
	router := newHomestayRouter(t, queryService)

	request := httptest.NewRequest(http.MethodGet, "/homestays", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}

	var body api.FetchHomestays200Response
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Data) != 1 || len(body.Data[0].HomestayRooms) != 1 || len(body.Data[0].HomestayRooms[0].HomestayImages) != 1 {
		t.Fatalf("response = %#v", body)
	}
	if got := body.Data[0].HomestayRooms[0].HomestayImages[0].ImagePath; got != "homestays/image.png" {
		t.Fatalf("image path = %q", got)
	}
}

func TestFetchHomestaysReturnsError(t *testing.T) {
	router := newHomestayRouter(t, homestayQueryService{err: errors.New("database unavailable")})
	request := httptest.NewRequest(http.MethodGet, "/homestays", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func newHomestayRouter(t *testing.T, queryService homestayqueryservice.HomestayQueryService) http.Handler {
	t.Helper()
	fetchUsecase := homestayusecase.NewFetchHomestaysUseCase(queryService)
	createUsecase := homestayusecase.NewCreateHomestayUsecase(nil)
	handler := controller.NewHomestayController(fetchUsecase, createUsecase)
	router := shared.NewSetupRouter(t)
	router.GET("/homestays", handler.FetchHomestays)
	return router
}

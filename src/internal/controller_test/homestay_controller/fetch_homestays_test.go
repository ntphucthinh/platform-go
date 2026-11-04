package homestay_controller_test

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
	"platform-go/internal/infrastructure/dbutil"
	"platform-go/internal/infrastructure/qsdto"
	"platform-go/internal/interfaces/controller"
	homestayqueryservice "platform-go/internal/usecases/homestays/queryservice"
	homestayusecase "platform-go/internal/usecases/homestays/usecase"
)

type homestayQueryService struct {
	homestays []qsdto.HomestayDto
	homestay  *qsdto.HomestayDto
	err       error
}

func (s homestayQueryService) FetchHomestays(ctx context.Context, page, perPage int) (*dbutil.Paginator[qsdto.HomestayDto], error) {
	if s.err != nil {
		return nil, s.err
	}
	return &dbutil.Paginator[qsdto.HomestayDto]{
		Data: s.homestays,
		PaginatorInfo: api.PaginatorInfo{
			Count:        len(s.homestays),
			CurrentPage:  page,
			HasMorePages: false,
			LastPage:     1,
			PerPage:      perPage,
			Total:        int64(len(s.homestays)),
		},
	}, nil
}

func (s homestayQueryService) FetchHomestayByID(ctx context.Context, id int64) (*qsdto.HomestayDto, error) {
	return s.homestay, s.err
}

var _ homestayqueryservice.HomestayQueryService = homestayQueryService{}

func TestFetchHomestays(t *testing.T) {
	description := "Room description"
	price := 10000.0
	priceUnit := string(api.PriceUnitNightPerRoom)
	roomID := int64(11)
	queryService := homestayQueryService{homestays: []qsdto.HomestayDto{{
		ID: 1, Name: "Homestay", Address: "Da Nang",
		Rooms: []qsdto.HomestayRoomDto{{
			ID: 11, HomestayID: 1, Name: "Room 1", Description: &description, PriceAmount: &price, PriceUnit: &priceUnit, SortOrder: 1,
			Images: []qsdto.HomestayImageDto{{ID: 3, RoomID: roomID, ImagePath: "homestays/image.png", SortOrder: 1, IsCover: true}},
		}},
	}}}
	router := newHomestayRouter(t, queryService)

	request := httptest.NewRequest(http.MethodGet, "/homestays?page=1&perPage=10", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}

	var body api.FetchHomestayResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Data) != 1 || len(body.Data[0].HomestayRooms) != 1 || len(body.Data[0].HomestayRooms[0].HomestayImages) != 1 {
		t.Fatalf("response = %#v", body)
	}
	if got := body.Data[0].HomestayRooms[0].HomestayImages[0].ImagePath; got != "homestays/image.png" {
		t.Fatalf("image path = %q", got)
	}
	if body.PaginatorInfo.CurrentPage != 1 || body.PaginatorInfo.Count != 1 || body.PaginatorInfo.Total != 1 {
		t.Fatalf("paginatorInfo = %#v", body.PaginatorInfo)
	}

	rawJSON := response.Body.String()
	for _, camelField := range []string{
		`"googleMapsUrl"`,
		`"homestayRooms"`,
		`"homestayId"`,
		`"priceAmount"`,
		`"priceUnit"`,
		`"sortOrder"`,
		`"homestayImages"`,
		`"roomId"`,
		`"imagePath"`,
		`"isCover"`,
		`"paginatorInfo"`,
	} {
		if !strings.Contains(rawJSON, camelField) {
			t.Fatalf("expected JSON response to contain camelCase field %s, got: %s", camelField, rawJSON)
		}
	}
}

func TestFetchHomestaysReturnsError(t *testing.T) {
	router := newHomestayRouter(t, homestayQueryService{err: errors.New("database unavailable")})
	request := httptest.NewRequest(http.MethodGet, "/homestays?page=1&perPage=10", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func newHomestayRouter(t *testing.T, queryService homestayqueryservice.HomestayQueryService) http.Handler {
	t.Helper()
	fetchUsecase := homestayusecase.NewFetchHomestaysUseCase(queryService)
	createUsecase := homestayusecase.NewCreateHomestayUsecase(nil, nil)
	fetchByIDUsecase := homestayusecase.NewFetchHomestayByIDUseCase(queryService)
	updateUsecase := homestayusecase.NewUpdateHomestayUsecase(nil, nil)
	deleteUsecase := homestayusecase.NewDeleteHomestayUsecase(nil, nil)
	handler := controller.NewHomestayController(fetchUsecase, createUsecase, fetchByIDUsecase, updateUsecase, deleteUsecase)
	serverController := controller.NewServerController(nil, nil, handler, nil, nil)
	router := shared.NewSetupRouter(t)
	api.RegisterHandlers(router, serverController)
	return router
}

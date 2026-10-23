package homestay_controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/internal/infrastructure/repository"
	"platform-go/internal/interfaces/controller"
	homestayusecase "platform-go/internal/usecases/homestays/usecase"
)

type fakeDB struct {
	txFunc func(ctx context.Context, fn func(tx *gorm.DB) error) error
}

func (f *fakeDB) GetDB() *gorm.DB { return nil }
func (f *fakeDB) GetDBByCtx(ctx context.Context) *gorm.DB { return nil }
func (f *fakeDB) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	if f.txFunc != nil {
		return f.txFunc(ctx, fn)
	}
	return nil
}
func (f *fakeDB) Close() error { return nil }

func newCreateHomestayRouter(t *testing.T, db repository.DBInterface) http.Handler {
	t.Helper()
	fetchUsecase := homestayusecase.NewFetchHomestaysUseCase(nil)
	createUsecase := homestayusecase.NewCreateHomestayUsecase(db)
	handler := controller.NewHomestayController(fetchUsecase, createUsecase)
	router := shared.NewSetupRouter(t)
	router.POST("/homestays", handler.CreateHomestay)
	return router
}

func performPost(t *testing.T, router http.Handler, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/homestays", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func TestCreateHomestay_InvalidBody(t *testing.T) {
	router := newCreateHomestayRouter(t, &fakeDB{})
	response := performPost(t, router, []byte("{invalid-json"))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestCreateHomestay_ValidationErrors(t *testing.T) {
	t.Skip("manual validation is currently commented out in usecase")
}

func TestCreateHomestay_DatabaseError(t *testing.T) {
	db := &fakeDB{
		txFunc: func(ctx context.Context, fn func(tx *gorm.DB) error) error {
			return errors.New("database connection failed")
		},
	}
	router := newCreateHomestayRouter(t, db)

	rooms := []api.CreateHomestayRoomInput{}
	payload := api.CreateHomestayRequest{
		Name:    "Homestay Error",
		Address: "Address",
		Rooms:   &rooms,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	response := performPost(t, router, bodyBytes)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestCreateHomestay_SuccessWithRealDB(t *testing.T) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:root@tcp(platform-db:3306)/platform_db?parseTime=true"
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("skipping integration test: database unavailable")
	}
	sqlDB, err := db.DB()
	if err != nil || sqlDB.Ping() != nil {
		t.Skip("skipping integration test: database unreachable")
	}

	schemaDB := repository.NewSchemaDB(db)
	router := newCreateHomestayRouter(t, schemaDB)

	desc := "A lovely homestay"
	mapsURL := "https://maps.google.com/?q=dalat"
	price := 500000.0
	unit := "per_room_per_night"
	imagePath := "homestays/test-image-1.jpg"

	payload := api.CreateHomestayRequest{
		Name:          "Test Homestay",
		Address:       "456 Tran Phu, Da Lat",
		Description:   &desc,
		GoogleMapsUrl: &mapsURL,
		Rooms: &[]api.CreateHomestayRoomInput{
			{
				Name:        "Deluxe Room",
				Description: &desc,
				PriceAmount: &price,
				PriceUnit:   &unit,
				SortOrder:   1,
				Images: &[]api.HomestayImageInput{
					{
						ImagePath: &imagePath,
						SortOrder: 1,
						IsCover:   true,
					},
				},
			},
		},
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
	if res.Message != homestayusecase.CreateHomestaysUseCaseSuccessMessage {
		t.Fatalf("expected message %q, got %q", homestayusecase.CreateHomestaysUseCaseSuccessMessage, res.Message)
	}
}

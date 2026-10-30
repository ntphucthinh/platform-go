package homestay_controller_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
	"platform-go/internal/interfaces/controller"
	homestayusecase "platform-go/internal/usecases/homestays/usecase"
)

func newDeleteHomestayRouter(t *testing.T, db repository.DBInterface) http.Handler {
	t.Helper()
	fetchUsecase := homestayusecase.NewFetchHomestaysUseCase(nil)
	createUsecase := homestayusecase.NewCreateHomestayUsecase(db)
	fetchByIDUsecase := homestayusecase.NewFetchHomestayByIDUseCase(nil)
	updateUsecase := homestayusecase.NewUpdateHomestayUsecase(db)
	deleteUsecase := homestayusecase.NewDeleteHomestayUsecase(db)
	handler := controller.NewHomestayController(fetchUsecase, createUsecase, fetchByIDUsecase, updateUsecase, deleteUsecase)
	serverController := controller.NewServerController(nil, nil, handler, nil)
	router := shared.NewSetupRouter(t)
	api.RegisterHandlers(router, serverController)
	return router
}

func performHomestayDelete(t *testing.T, router http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodDelete, path, nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func TestDeleteHomestay_InvalidID(t *testing.T) {
	router := newDeleteHomestayRouter(t, &fakeDB{})
	response := performHomestayDelete(t, router, "/homestays/not-a-number")
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestDeleteHomestay_NotFound(t *testing.T) {
	db := &fakeDB{
		txFunc: func(ctx context.Context, fn func(tx *gorm.DB) error) error {
			return homestayusecase.ErrHomestayNotFound
		},
	}
	router := newDeleteHomestayRouter(t, db)
	response := performHomestayDelete(t, router, "/homestays/999")
	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestDeleteHomestay_DatabaseError(t *testing.T) {
	db := &fakeDB{
		txFunc: func(ctx context.Context, fn func(tx *gorm.DB) error) error {
			return errors.New("db error")
		},
	}
	router := newDeleteHomestayRouter(t, db)
	response := performHomestayDelete(t, router, "/homestays/1")
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestDeleteHomestay_SuccessWithRealDB(t *testing.T) {
	db := shared.SetupTestDB(t)

	// 1. Seed homestay with 2 rooms and images
	homestay := dbmodel.Homestay{
		Name:    "Homestay To Delete",
		Address: "123 Delete St",
	}
	if err := db.Create(&homestay).Error; err != nil {
		t.Fatalf("failed to create homestay: %v", err)
	}

	room1 := dbmodel.HomestayRoom{
		HomestayID: homestay.ID,
		Name:       "Room 1",
		SortOrder:  1,
	}
	room2 := dbmodel.HomestayRoom{
		HomestayID: homestay.ID,
		Name:       "Room 2",
		SortOrder:  2,
	}
	if err := db.Create(&room1).Error; err != nil {
		t.Fatalf("failed to create room1: %v", err)
	}
	if err := db.Create(&room2).Error; err != nil {
		t.Fatalf("failed to create room2: %v", err)
	}

	img1 := dbmodel.HomestayImage{
		RoomID:    room1.ID,
		ImagePath: "homestays/img1.jpg",
		SortOrder: 1,
		IsCover:   true,
	}
	img2 := dbmodel.HomestayImage{
		RoomID:    room2.ID,
		ImagePath: "homestays/img2.jpg",
		SortOrder: 1,
		IsCover:   true,
	}
	if err := db.Create(&img1).Error; err != nil {
		t.Fatalf("failed to create img1: %v", err)
	}
	if err := db.Create(&img2).Error; err != nil {
		t.Fatalf("failed to create img2: %v", err)
	}

	schemaDB := repository.NewSchemaDB(db)
	router := newDeleteHomestayRouter(t, schemaDB)

	path := fmt.Sprintf("/homestays/%d", homestay.ID)
	response := performHomestayDelete(t, router, path)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", response.Code, response.Body.String())
	}

	var res api.MessageOnlyResponse
	if err := json.Unmarshal(response.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Message != homestayusecase.DeleteHomestaysUseCaseSuccessMessage {
		t.Fatalf("expected message %q, got %q", homestayusecase.DeleteHomestaysUseCaseSuccessMessage, res.Message)
	}

	// 2. Verify deletion in DB
	var deletedHomestay dbmodel.Homestay
	if err := db.First(&deletedHomestay, homestay.ID).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected homestay to be deleted, got err: %v", err)
	}

	var remainingRooms []dbmodel.HomestayRoom
	if err := db.Where("homestay_id = ?", homestay.ID).Find(&remainingRooms).Error; err != nil {
		t.Fatal(err)
	}
	if len(remainingRooms) != 0 {
		t.Fatalf("expected 0 rooms remaining, found %d", len(remainingRooms))
	}

	var remainingImages []dbmodel.HomestayImage
	if err := db.Where("room_id IN ?", []int64{room1.ID, room2.ID}).Find(&remainingImages).Error; err != nil {
		t.Fatal(err)
	}
	if len(remainingImages) != 0 {
		t.Fatalf("expected 0 images remaining, found %d", len(remainingImages))
	}
}

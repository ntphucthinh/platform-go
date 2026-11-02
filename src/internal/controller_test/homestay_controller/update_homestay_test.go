package homestay_controller_test

import (
	"bytes"
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

func newUpdateHomestayRouter(t *testing.T, db repository.DBInterface) http.Handler {
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

func performPut(t *testing.T, router http.Handler, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPut, path, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func TestUpdateHomestay_InvalidBody(t *testing.T) {
	router := newUpdateHomestayRouter(t, &fakeDB{})
	response := performPut(t, router, "/homestays/1", []byte("{invalid-json"))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestUpdateHomestay_InvalidID(t *testing.T) {
	router := newUpdateHomestayRouter(t, &fakeDB{})
	bodyBytes, _ := json.Marshal(api.UpdateHomestayRequest{
		Name:    "Updated Name",
		Address: "Updated Address",
	})
	response := performPut(t, router, "/homestays/not-a-number", bodyBytes)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestUpdateHomestay_DatabaseError(t *testing.T) {
	db := &fakeDB{
		txFunc: func(ctx context.Context, fn func(tx *gorm.DB) error) error {
			return errors.New("database connection failed")
		},
	}
	router := newUpdateHomestayRouter(t, db)

	bodyBytes, err := json.Marshal(api.UpdateHomestayRequest{
		Name:    "Updated Name",
		Address: "Updated Address",
	})
	if err != nil {
		t.Fatal(err)
	}

	response := performPut(t, router, "/homestays/1", bodyBytes)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestUpdateHomestay_NotFound(t *testing.T) {
	db := &fakeDB{
		txFunc: func(ctx context.Context, fn func(tx *gorm.DB) error) error {
			return homestayusecase.ErrHomestayNotFound
		},
	}
	router := newUpdateHomestayRouter(t, db)

	bodyBytes, err := json.Marshal(api.UpdateHomestayRequest{
		Name:    "Updated Name",
		Address: "Updated Address",
	})
	if err != nil {
		t.Fatal(err)
	}

	response := performPut(t, router, "/homestays/999", bodyBytes)
	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestUpdateHomestay_InvalidPriceUnit(t *testing.T) {
	router := newUpdateHomestayRouter(t, &fakeDB{})

	invalidPayload := map[string]interface{}{
		"name":    "Updated Homestay",
		"address": "456 Tran Phu, Da Lat",
		"rooms": []map[string]interface{}{
			{
				"name":      "Deluxe Room",
				"sortOrder": 1,
				"priceUnit": "invalid_unit",
			},
		},
	}
	bodyBytes, err := json.Marshal(invalidPayload)
	if err != nil {
		t.Fatal(err)
	}

	response := performPut(t, router, "/homestays/1", bodyBytes)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", response.Code, response.Body.String())
	}

	var res map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}

	errorsMap, ok := res["errors"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected errors object in response, got: %#v", res)
	}
	if _, exists := errorsMap["priceUnit"]; !exists {
		t.Fatalf("expected priceUnit validation error, got: %#v", errorsMap)
	}
}

func TestUpdateHomestay_SuccessWithRealDB(t *testing.T) {
	db := shared.SetupTestDB(t)

	// 1. Seed initial homestay with 2 rooms and images
	initialHomestay := dbmodel.Homestay{
		Name:    "Initial Homestay",
		Address: "123 Initial Street",
	}
	if err := db.Create(&initialHomestay).Error; err != nil {
		t.Fatalf("failed to create initial homestay: %v", err)
	}

	room1 := dbmodel.HomestayRoom{
		HomestayID: initialHomestay.ID,
		Name:       "Room 1 To Keep",
		SortOrder:  1,
	}
	room2 := dbmodel.HomestayRoom{
		HomestayID: initialHomestay.ID,
		Name:       "Room 2 To Delete",
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
		RoomID:    room1.ID,
		ImagePath: "homestays/img2-to-delete.jpg",
		SortOrder: 2,
		IsCover:   false,
	}
	img3 := dbmodel.HomestayImage{
		RoomID:    room2.ID,
		ImagePath: "homestays/img3.jpg",
		SortOrder: 1,
		IsCover:   true,
	}
	if err := db.Create(&img1).Error; err != nil {
		t.Fatalf("failed to create img1: %v", err)
	}
	if err := db.Create(&img2).Error; err != nil {
		t.Fatalf("failed to create img2: %v", err)
	}
	if err := db.Create(&img3).Error; err != nil {
		t.Fatalf("failed to create img3: %v", err)
	}

	schemaDB := repository.NewSchemaDB(db)
	router := newUpdateHomestayRouter(t, schemaDB)

	// 2. Perform update:
	// - Homestay updated name and address
	// - Room 1 updated (img1 updated, img2 omitted -> deleted, new img4 added)
	// - Room 2 omitted -> deleted (img3 cascade deleted)
	// - Room 3 added
	newDesc := "Updated description"
	newMapsURL := "https://maps.google.com/?q=updated"
	price1 := 600000.0
	unit1 := api.PriceUnitNightPerRoom
	updatedImg1Path := "homestays/img1-updated.jpg"
	newImg4Path := "homestays/img4-new.jpg"

	price3 := 800000.0
	unit3 := api.PriceUnitNightPerHouse
	newImg5Path := "homestays/img5-new.jpg"

	updatePayload := api.UpdateHomestayRequest{
		Name:          "Updated Homestay Name",
		Address:       "456 Updated Street",
		Description:   &newDesc,
		GoogleMapsUrl: &newMapsURL,
		Rooms: &[]api.UpdateHomestayRoomInput{
			{
				Id:          &room1.ID,
				Name:        "Room 1 Updated Name",
				Description: &newDesc,
				PriceAmount: &price1,
				PriceUnit:   &unit1,
				SortOrder:   10,
				Images: &[]api.UpdateHomestayImageInput{
					{
						Id:        &img1.ID,
						ImagePath: &updatedImg1Path,
						SortOrder: 1,
						IsCover:   true,
					},
					{
						ImagePath: &newImg4Path,
						SortOrder: 2,
						IsCover:   false,
					},
				},
			},
			{
				Name:        "Room 3 Newly Added",
				PriceAmount: &price3,
				PriceUnit:   &unit3,
				SortOrder:   20,
				Images: &[]api.UpdateHomestayImageInput{
					{
						ImagePath: &newImg5Path,
						SortOrder: 1,
						IsCover:   true,
					},
				},
			},
		},
	}

	bodyBytes, err := json.Marshal(updatePayload)
	if err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprintf("/homestays/%d", initialHomestay.ID)
	response := performPut(t, router, path, bodyBytes)
	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", response.Code, response.Body.String())
	}

	var res api.MessageOnlyResponse
	if err := json.Unmarshal(response.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Message != homestayusecase.UpdateHomestaysUseCaseSuccessMessage {
		t.Fatalf("expected message %q, got %q", homestayusecase.UpdateHomestaysUseCaseSuccessMessage, res.Message)
	}

	// 3. Verify database state
	var updatedHomestay dbmodel.Homestay
	if err := db.First(&updatedHomestay, initialHomestay.ID).Error; err != nil {
		t.Fatalf("failed to find updated homestay: %v", err)
	}
	if updatedHomestay.Name != "Updated Homestay Name" || updatedHomestay.Address != "456 Updated Street" {
		t.Fatalf("homestay fields not updated: %#v", updatedHomestay)
	}

	var roomsInDB []dbmodel.HomestayRoom
	if err := db.Where("homestay_id = ?", initialHomestay.ID).Order("sort_order asc").Find(&roomsInDB).Error; err != nil {
		t.Fatalf("failed to fetch rooms: %v", err)
	}
	if len(roomsInDB) != 2 {
		t.Fatalf("expected 2 rooms in DB, got %d", len(roomsInDB))
	}

	// Check room1 updated
	if roomsInDB[0].ID != room1.ID || roomsInDB[0].Name != "Room 1 Updated Name" || roomsInDB[0].SortOrder != 10 {
		t.Fatalf("room 1 was not properly updated: %#v", roomsInDB[0])
	}
	// Check room2 deleted
	var deletedRoom dbmodel.HomestayRoom
	if err := db.First(&deletedRoom, room2.ID).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("room 2 should have been deleted, got err: %v", err)
	}
	// Check room3 created
	if roomsInDB[1].Name != "Room 3 Newly Added" || roomsInDB[1].SortOrder != 20 {
		t.Fatalf("room 3 was not properly created: %#v", roomsInDB[1])
	}

	// Check images for room1
	var room1Images []dbmodel.HomestayImage
	if err := db.Where("room_id = ?", room1.ID).Order("sort_order asc").Find(&room1Images).Error; err != nil {
		t.Fatalf("failed to fetch room1 images: %v", err)
	}
	if len(room1Images) != 2 {
		t.Fatalf("expected 2 images for room1, got %d", len(room1Images))
	}
	if room1Images[0].ID != img1.ID || room1Images[0].ImagePath != updatedImg1Path {
		t.Fatalf("img1 was not updated: %#v", room1Images[0])
	}
	if room1Images[1].ImagePath != newImg4Path {
		t.Fatalf("new image 4 was not created: %#v", room1Images[1])
	}

	// Check img2 deleted
	var deletedImg dbmodel.HomestayImage
	if err := db.First(&deletedImg, img2.ID).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("img 2 should have been deleted, got err: %v", err)
	}
	// Check img3 deleted
	if err := db.First(&deletedImg, img3.ID).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("img 3 should have been deleted, got err: %v", err)
	}
}

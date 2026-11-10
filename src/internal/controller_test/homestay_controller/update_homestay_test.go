package homestay_controller_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/internal/infrastructure/dbmodel"
	homestayusecase "platform-go/internal/usecases/homestays/usecase"
	"platform-go/tests/factory"
)

func TestHomestayController_UpdateHomestay(t *testing.T) {
	t.Run("Normal_UpdateHomestaySuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		initialHomestay := factory.NewHomestayFactory().
			WithName("Initial Homestay").
			WithAddress("123 Initial Street").
			MustInsert(shared.TestDB)

		room1 := factory.NewHomestayRoomFactory().
			WithHomestayID(initialHomestay.ID).
			WithName("Room 1 To Keep").
			WithSortOrder(1).
			MustInsert(shared.TestDB)

		room2 := factory.NewHomestayRoomFactory().
			WithHomestayID(initialHomestay.ID).
			WithName("Room 2 To Delete").
			WithSortOrder(2).
			MustInsert(shared.TestDB)

		img1 := factory.NewHomestayImageFactory().
			WithRoomID(room1.ID).
			WithImagePath("homestays/img1.jpg").
			WithSortOrder(1).
			WithCover(true).
			MustInsert(shared.TestDB)

		img2 := factory.NewHomestayImageFactory().
			WithRoomID(room1.ID).
			WithImagePath("homestays/img2-to-delete.jpg").
			WithSortOrder(2).
			WithCover(false).
			MustInsert(shared.TestDB)

		img3 := factory.NewHomestayImageFactory().
			WithRoomID(room2.ID).
			WithImagePath("homestays/img3.jpg").
			WithSortOrder(1).
			WithCover(true).
			MustInsert(shared.TestDB)

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

		path := fmt.Sprintf("/homestays/%d", initialHomestay.ID)
		req := shared.NewRequest(t, admin, http.MethodPut, path, updatePayload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.MessageOnlyResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, homestayusecase.UpdateHomestaysUseCaseSuccessMessage, res.Message)

		// Verify database state
		var updatedHomestay dbmodel.Homestay
		require.NoError(t, shared.TestDB.First(&updatedHomestay, initialHomestay.ID).Error)
		assert.Equal(t, "Updated Homestay Name", updatedHomestay.Name)
		assert.Equal(t, "456 Updated Street", updatedHomestay.Address)

		var roomsInDB []dbmodel.HomestayRoom
		require.NoError(t, shared.TestDB.Where("homestay_id = ?", initialHomestay.ID).Order("sort_order asc").Find(&roomsInDB).Error)
		require.Len(t, roomsInDB, 2)

		// Room 1 updated
		assert.Equal(t, room1.ID, roomsInDB[0].ID)
		assert.Equal(t, "Room 1 Updated Name", roomsInDB[0].Name)
		assert.Equal(t, int32(10), roomsInDB[0].SortOrder)

		// Room 2 deleted
		var deletedRoom dbmodel.HomestayRoom
		err := shared.TestDB.First(&deletedRoom, room2.ID).Error
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		// Room 3 created
		assert.Equal(t, "Room 3 Newly Added", roomsInDB[1].Name)
		assert.Equal(t, int32(20), roomsInDB[1].SortOrder)

		// Check images for room1
		var room1Images []dbmodel.HomestayImage
		require.NoError(t, shared.TestDB.Where("room_id = ?", room1.ID).Order("sort_order asc").Find(&room1Images).Error)
		require.Len(t, room1Images, 2)
		assert.Equal(t, img1.ID, room1Images[0].ID)
		assert.Equal(t, updatedImg1Path, room1Images[0].ImagePath)
		assert.Equal(t, newImg4Path, room1Images[1].ImagePath)

		// Check img2 and img3 deleted
		var deletedImg dbmodel.HomestayImage
		assert.True(t, errors.Is(shared.TestDB.First(&deletedImg, img2.ID).Error, gorm.ErrRecordNotFound))
		assert.True(t, errors.Is(shared.TestDB.First(&deletedImg, img3.ID).Error, gorm.ErrRecordNotFound))
	})

	t.Run("Error_InvalidPriceUnit", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		homestay := factory.NewHomestayFactory().MustInsert(shared.TestDB)

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

		req := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/homestays/%d", homestay.ID), invalidPayload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())

		var res map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		errorsMap, ok := res["errors"].(map[string]interface{})
		require.True(t, ok)
		assert.Contains(t, errorsMap, "priceUnit")
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		payload := api.UpdateHomestayRequest{
			Name:    "Updated Name",
			Address: "Updated Address",
		}

		req := shared.NewRequest(t, admin, http.MethodPut, "/homestays/999999", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Error_InvalidID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		payload := api.UpdateHomestayRequest{
			Name:    "Updated Name",
			Address: "Updated Address",
		}

		req := shared.NewRequest(t, admin, http.MethodPut, "/homestays/not-a-number", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error_InvalidBody", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodPut, "/homestays/1", strings.NewReader("{invalid-json"))
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		payload := api.UpdateHomestayRequest{
			Name:    "Updated Name",
			Address: "Updated Address",
		}

		req := shared.NewUnAuthRequest(t, http.MethodPut, "/homestays/1", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

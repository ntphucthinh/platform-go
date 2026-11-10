package homestay_controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/internal/infrastructure/dbmodel"
	homestayusecase "platform-go/internal/usecases/homestays/usecase"
)

func TestHomestayController_CreateHomestay(t *testing.T) {
	t.Run("Normal_CreateHomestaySuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		desc := "A lovely homestay"
		mapsURL := "https://maps.google.com/?q=dalat"
		price := 500000.0
		unit := api.PriceUnitNightPerRoom
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

		req := shared.NewRequest(t, admin, http.MethodPost, "/homestays", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.MessageOnlyResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, homestayusecase.CreateHomestaysUseCaseSuccessMessage, res.Message)

		var homestay dbmodel.Homestay
		require.NoError(t, shared.TestDB.Where("name = ?", "Test Homestay").First(&homestay).Error)
		assert.Equal(t, "456 Tran Phu, Da Lat", homestay.Address)

		var rooms []dbmodel.HomestayRoom
		require.NoError(t, shared.TestDB.Where("homestay_id = ?", homestay.ID).Find(&rooms).Error)
		require.Len(t, rooms, 1)
		assert.Equal(t, "Deluxe Room", rooms[0].Name)

		var images []dbmodel.HomestayImage
		require.NoError(t, shared.TestDB.Where("room_id = ?", rooms[0].ID).Find(&images).Error)
		require.Len(t, images, 1)
		assert.Equal(t, "homestays/test-image-1.jpg", images[0].ImagePath)
	})

	t.Run("Error_ValidationErrors", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		invalidPayload := map[string]interface{}{
			"name":    "Test Homestay",
			"address": "456 Tran Phu, Da Lat",
			"rooms": []map[string]interface{}{
				{
					"name":      "Deluxe Room",
					"sortOrder": 1,
					"priceUnit": "invalid_unit",
				},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/homestays", invalidPayload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())

		var res map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		errorsMap, ok := res["errors"].(map[string]interface{})
		require.True(t, ok, "expected errors object in response")
		assert.Contains(t, errorsMap, "priceUnit")
	})

	t.Run("Error_InvalidBody", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodPost, "/homestays", strings.NewReader("{invalid-json"))
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		payload := api.CreateHomestayRequest{
			Name:    "Test Homestay",
			Address: "456 Tran Phu, Da Lat",
		}

		req := shared.NewUnAuthRequest(t, http.MethodPost, "/homestays", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

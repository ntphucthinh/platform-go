package homestay_controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/tests/factory"
)

func TestHomestayController_FetchHomestayByID(t *testing.T) {
	t.Run("Normal_GetHomestayByIDSuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		desc := "A lovely homestay in Da Nang"
		mapsURL := "https://maps.google.com/?q=danang"
		roomDesc := "Deluxe Ocean View"
		price := 1200000.0
		unit := string(api.PriceUnitNightPerRoom)

		homestay := factory.NewHomestayFactory().
			WithName("Da Nang Homestay").
			WithAddress("123 Vo Nguyen Giap").
			WithDescription(desc).
			WithGoogleMapsURL(mapsURL).
			MustInsert(shared.TestDB)

		room := factory.NewHomestayRoomFactory().
			WithHomestayID(homestay.ID).
			WithName("Room 101").
			WithDescription(roomDesc).
			WithPrice(price, unit).
			WithSortOrder(1).
			MustInsert(shared.TestDB)

		factory.NewHomestayImageFactory().
			WithRoomID(room.ID).
			WithImagePath("homestays/ocean.jpg").
			WithSortOrder(1).
			WithCover(true).
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, fmt.Sprintf("/homestays/%d", homestay.ID), nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.FetchHomestayByIDResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))

		assert.Equal(t, homestay.ID, res.Data.Id)
		assert.Equal(t, "Da Nang Homestay", res.Data.Name)
		require.Len(t, res.Data.HomestayRooms, 1)
		assert.Equal(t, "Room 101", res.Data.HomestayRooms[0].Name)
		require.Len(t, res.Data.HomestayRooms[0].HomestayImages, 1)
		assert.Equal(t, "homestays/ocean.jpg", res.Data.HomestayRooms[0].HomestayImages[0].ImagePath)

		rawJSON := w.Body.String()
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
		} {
			assert.True(t, strings.Contains(rawJSON, camelField), "expected camelField %s in response", camelField)
		}
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodGet, "/homestays/999999", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Error_InvalidID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodGet, "/homestays/not-a-number", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})


}

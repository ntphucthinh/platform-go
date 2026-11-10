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
	"platform-go/tests/factory"
)

func TestHomestayController_FetchHomestays(t *testing.T) {
	t.Run("Normal_GetHomestaysList", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		desc := "Room description"
		price := 10000.0
		unit := string(api.PriceUnitNightPerRoom)

		homestay := factory.NewHomestayFactory().
			WithName("Homestay Da Nang").
			WithAddress("Da Nang").
			MustInsert(shared.TestDB)

		room := factory.NewHomestayRoomFactory().
			WithHomestayID(homestay.ID).
			WithName("Room 1").
			WithDescription(desc).
			WithPrice(price, unit).
			WithSortOrder(1).
			MustInsert(shared.TestDB)

		factory.NewHomestayImageFactory().
			WithRoomID(room.ID).
			WithImagePath("homestays/image.png").
			WithSortOrder(1).
			WithCover(true).
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, "/homestays?page=1&perPage=10", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var body api.FetchHomestayResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

		require.Len(t, body.Data, 1)
		assert.Equal(t, homestay.Name, body.Data[0].Name)
		require.Len(t, body.Data[0].HomestayRooms, 1)
		assert.Equal(t, "Room 1", body.Data[0].HomestayRooms[0].Name)
		require.Len(t, body.Data[0].HomestayRooms[0].HomestayImages, 1)
		assert.Equal(t, "homestays/image.png", body.Data[0].HomestayRooms[0].HomestayImages[0].ImagePath)
		assert.Equal(t, 1, body.PaginatorInfo.CurrentPage)
		assert.Equal(t, 1, body.PaginatorInfo.Count)
		assert.Equal(t, int64(1), body.PaginatorInfo.Total)

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
			`"paginatorInfo"`,
		} {
			assert.True(t, strings.Contains(rawJSON, camelField), "expected camelField %s in response", camelField)
		}
	})


}

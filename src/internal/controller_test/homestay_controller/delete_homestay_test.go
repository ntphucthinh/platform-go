package homestay_controller_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestHomestayController_DeleteHomestay(t *testing.T) {
	t.Run("Normal_DeleteHomestaySuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		homestay := factory.NewHomestayFactory().
			WithName("Homestay To Delete").
			WithAddress("123 Delete St").
			MustInsert(shared.TestDB)

		room1 := factory.NewHomestayRoomFactory().
			WithHomestayID(homestay.ID).
			WithName("Room 1").
			WithSortOrder(1).
			MustInsert(shared.TestDB)

		room2 := factory.NewHomestayRoomFactory().
			WithHomestayID(homestay.ID).
			WithName("Room 2").
			WithSortOrder(2).
			MustInsert(shared.TestDB)

		img1 := factory.NewHomestayImageFactory().
			WithRoomID(room1.ID).
			WithImagePath("homestays/img1.jpg").
			MustInsert(shared.TestDB)

		img2 := factory.NewHomestayImageFactory().
			WithRoomID(room1.ID).
			WithImagePath("homestays/img2.jpg").
			MustInsert(shared.TestDB)

		img3 := factory.NewHomestayImageFactory().
			WithRoomID(room2.ID).
			WithImagePath("homestays/img3.jpg").
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodDelete, fmt.Sprintf("/homestays/%d", homestay.ID), nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.MessageOnlyResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, homestayusecase.DeleteHomestaysUseCaseSuccessMessage, res.Message)

		// Verify database state: homestay, rooms, images deleted
		var deletedHomestay dbmodel.Homestay
		err := shared.TestDB.First(&deletedHomestay, homestay.ID).Error
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "homestay should be deleted")

		var deletedRoom dbmodel.HomestayRoom
		assert.True(t, errors.Is(shared.TestDB.First(&deletedRoom, room1.ID).Error, gorm.ErrRecordNotFound))
		assert.True(t, errors.Is(shared.TestDB.First(&deletedRoom, room2.ID).Error, gorm.ErrRecordNotFound))

		var deletedImg dbmodel.HomestayImage
		assert.True(t, errors.Is(shared.TestDB.First(&deletedImg, img1.ID).Error, gorm.ErrRecordNotFound))
		assert.True(t, errors.Is(shared.TestDB.First(&deletedImg, img2.ID).Error, gorm.ErrRecordNotFound))
		assert.True(t, errors.Is(shared.TestDB.First(&deletedImg, img3.ID).Error, gorm.ErrRecordNotFound))
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodDelete, "/homestays/999999", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Error_InvalidID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodDelete, "/homestays/not-a-number", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		req := shared.NewUnAuthRequest(t, http.MethodDelete, "/homestays/1", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

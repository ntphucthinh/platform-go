package selection_controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/tests/factory"
)

func TestSelectionController_FetchSelectionTravelThemes(t *testing.T) {
	t.Run("Normal_GetTravelThemesList", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		theme1 := factory.NewTravelThemeFactory().
			WithName("Thiên nhiên").
			MustInsert(shared.TestDB)

		theme2 := factory.NewTravelThemeFactory().
			WithName("Ẩm thực").
			MustInsert(shared.TestDB)

		req := shared.NewRequest(t, admin, http.MethodGet, "/selections/travel-themes", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.SelectionTravelThemeResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))

		require.Len(t, res.Data, 2)
		assert.Equal(t, theme1.ID, res.Data[0].Id)
		assert.Equal(t, "Thiên nhiên", res.Data[0].Name)
		assert.Equal(t, theme2.ID, res.Data[1].Id)
		assert.Equal(t, "Ẩm thực", res.Data[1].Name)
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		req := shared.NewUnAuthRequest(t, http.MethodGet, "/selections/travel-themes", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

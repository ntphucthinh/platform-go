package auth_controller_test

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

func TestAuthController_Flow(t *testing.T) {
	t.Run("Normal_AuthFlowSuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		user := factory.NewUserFactory().
			WithName("Admin User").
			WithEmail("admin@platform.com").
			WithActive().
			MustInsert(shared.TestDB)

		// 1. Login
		loginPayload := api.LoginRequest{
			Email:    "admin@platform.com",
			Password: "123456",
		}
		loginReq := shared.NewUnAuthRequest(t, http.MethodPost, "/auth/login", loginPayload)
		loginRec := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(loginRec, loginReq)

		require.Equal(t, http.StatusOK, loginRec.Code, loginRec.Body.String())

		var loginBody api.AuthResponse
		require.NoError(t, json.Unmarshal(loginRec.Body.Bytes(), &loginBody))
		assert.NotEmpty(t, loginBody.Data.CsrfToken)

		cookies := loginRec.Result().Cookies()
		jwtCookie := findCookie(t, cookies, "jwt")
		refreshCookie := findCookie(t, cookies, "refresh_token")
		csrfCookie := findCookie(t, cookies, "csrf")

		assert.Equal(t, "/", refreshCookie.Path)
		assert.True(t, jwtCookie.HttpOnly)

		// 2. Fetch current user (/auth/me)
		meReq := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
		meReq.AddCookie(jwtCookie)
		meReq.AddCookie(csrfCookie)
		meReq.Header.Set("X-CSRF-Token", csrfCookie.Value)
		meRec := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(meRec, meReq)

		require.Equal(t, http.StatusOK, meRec.Code, meRec.Body.String())

		var meBody api.LoginUser
		require.NoError(t, json.Unmarshal(meRec.Body.Bytes(), &meBody))
		assert.Equal(t, int64(user.ID), meBody.Id)
		assert.Equal(t, "admin@platform.com", string(meBody.Email))

		// 3. Refresh token (/auth/refresh)
		refreshReq := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
		refreshReq.AddCookie(refreshCookie)
		refreshReq.AddCookie(csrfCookie)
		refreshReq.Header.Set("X-CSRF-Token", csrfCookie.Value)
		refreshRec := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(refreshRec, refreshReq)

		require.Equal(t, http.StatusOK, refreshRec.Code, refreshRec.Body.String())

		// 4. Logout (/auth/logout)
		logoutReq := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
		logoutReq.AddCookie(jwtCookie)
		logoutReq.AddCookie(csrfCookie)
		logoutReq.Header.Set("X-CSRF-Token", csrfCookie.Value)
		logoutRec := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(logoutRec, logoutReq)

		require.Equal(t, http.StatusNoContent, logoutRec.Code)
	})

	t.Run("Error_InvalidCredentials", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		factory.NewUserFactory().
			WithEmail("admin@platform.com").
			MustInsert(shared.TestDB)

		loginPayload := api.LoginRequest{
			Email:    "admin@platform.com",
			Password: "wrong-password",
		}
		loginReq := shared.NewUnAuthRequest(t, http.MethodPost, "/auth/login", loginPayload)
		loginRec := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(loginRec, loginReq)

		assert.Equal(t, http.StatusUnauthorized, loginRec.Code)
	})

	t.Run("Error_RejectsUnauthenticatedMe", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		req := shared.NewUnAuthRequest(t, http.MethodGet, "/auth/me", nil)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Error_InvalidLoginPayload", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		req := shared.NewUnAuthRequest(t, http.MethodPost, "/auth/login", strings.NewReader("{invalid-json"))
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func findCookie(t *testing.T, cookies []*http.Cookie, name string) *http.Cookie {
	t.Helper()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	t.Fatalf("missing %s cookie", name)
	return nil
}

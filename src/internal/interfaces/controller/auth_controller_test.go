package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"platform-go/internal/domain/user"
	"platform-go/internal/infrastructure/security"
	"platform-go/internal/interfaces/controller"
	"platform-go/internal/middleware"
	"platform-go/internal/usecases"
)

type controllerAuthRepository struct{ account user.AuthUser }

func (r controllerAuthRepository) CreateUser(context.Context, string) (user.User, error) {
	return user.User{}, errors.New("not implemented")
}
func (r controllerAuthRepository) FindAuthUserByEmail(context.Context, string) (user.AuthUser, error) {
	return r.account, nil
}
func (r controllerAuthRepository) FindAuthUserByID(context.Context, int64) (user.AuthUser, error) {
	return r.account, nil
}

func TestAuthEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	account := user.AuthUser{ID: 7, Name: "Administrator", Email: "admin@example.com", PasswordHash: string(hash), AccountType: "admin", Active: true}
	jwt, err := security.NewJWTService("01234567890123456789012345678901")
	if err != nil {
		t.Fatal(err)
	}
	usecase := usecases.NewAuthUsecase(controllerAuthRepository{account: account}, jwt)
	handler := controller.NewAuthController(usecase, false)
	router := gin.New()
	router.Use(middleware.JWTAuth(usecase), middleware.CSRFProtection())
	router.POST("/auth/login", handler.Login)
	router.POST("/auth/refresh", handler.Refresh)
	router.POST("/auth/logout", handler.Logout)
	router.GET("/auth/me", handler.Me)

	login := request(router, http.MethodPost, "/auth/login", []byte(`{"email":"ADMIN@example.com","password":"correct-password"}`), nil, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login status = %d, body = %s", login.Code, login.Body.String())
	}
	cookies := login.Result().Cookies()
	jwtCookie := findCookie(t, cookies, "jwt")
	refreshCookie := findCookie(t, cookies, "refresh_token")
	csrfCookie := findCookie(t, cookies, "csrf")
	if refreshCookie.Path != "/" {
		t.Fatalf("refresh cookie path = %q, want /", refreshCookie.Path)
	}
	var response struct {
		Data struct {
			CSRFToken string `json:"csrfToken"`
		} `json:"data"`
	}
	if err := json.Unmarshal(login.Body.Bytes(), &response); err != nil || response.Data.CSRFToken == "" {
		t.Fatalf("login response missing csrf token: %v, %s", err, login.Body.String())
	}

	me := request(router, http.MethodGet, "/auth/me", nil, []*http.Cookie{jwtCookie, csrfCookie}, csrfCookie.Value)
	if me.Code != http.StatusOK {
		t.Fatalf("me status = %d, body = %s", me.Code, me.Body.String())
	}
	if !bytes.Contains(me.Body.Bytes(), []byte(`"email":"admin@example.com"`)) || bytes.Contains(me.Body.Bytes(), []byte("password")) {
		t.Fatalf("unexpected me body: %s", me.Body.String())
	}

	refresh := request(router, http.MethodPost, "/auth/refresh", nil, []*http.Cookie{refreshCookie, csrfCookie}, csrfCookie.Value)
	if refresh.Code != http.StatusOK {
		t.Fatalf("refresh status = %d, body = %s", refresh.Code, refresh.Body.String())
	}

	logout := request(router, http.MethodPost, "/auth/logout", nil, []*http.Cookie{jwtCookie, csrfCookie}, csrfCookie.Value)
	if logout.Code != http.StatusNoContent {
		t.Fatalf("logout status = %d", logout.Code)
	}
}

func TestProtectedAuthEndpointsRejectMissingAccessToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwt, err := security.NewJWTService("01234567890123456789012345678901")
	if err != nil {
		t.Fatal(err)
	}
	usecase := usecases.NewAuthUsecase(controllerAuthRepository{}, jwt)
	handler := controller.NewAuthController(usecase, false)
	router := gin.New()
	router.Use(middleware.JWTAuth(usecase), middleware.CSRFProtection())
	router.GET("/auth/me", handler.Me)

	response := request(router, http.MethodGet, "/auth/me", nil, nil, "")
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("missing access token status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func request(router http.Handler, method, path string, body []byte, cookies []*http.Cookie, csrf string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	if csrf != "" {
		req.Header.Set("X-CSRF-Token", csrf)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	return response
}

func findCookie(t *testing.T, cookies []*http.Cookie, name string) *http.Cookie {
	t.Helper()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	t.Fatalf("cookie %q was not set", name)
	return nil
}

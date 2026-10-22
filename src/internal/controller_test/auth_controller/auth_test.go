package auth_controller_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/internal/domain/user"
	"platform-go/internal/infrastructure/security"
	"platform-go/internal/interfaces/controller"
	"platform-go/internal/middleware"
	"platform-go/internal/usecases"
)

type authRepository struct{ account user.AuthUser }

func (r authRepository) CreateUser(context.Context, string) (user.User, error) {
	return user.User{}, nil
}
func (r authRepository) FindAuthUserByEmail(_ context.Context, email string) (user.AuthUser, error) {
	if email == r.account.Email {
		return r.account, nil
	}
	return user.AuthUser{}, nil
}
func (r authRepository) FindAuthUserByID(_ context.Context, id int64) (user.AuthUser, error) {
	if id == r.account.ID {
		return r.account, nil
	}
	return user.AuthUser{}, nil
}

func newAuthRouter(t *testing.T) *gin.Engine {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	jwt, err := security.NewJWTService("test-secret-with-at-least-thirty-two-characters")
	if err != nil {
		t.Fatal(err)
	}
	usecase := usecases.NewAuthUsecase(authRepository{account: user.AuthUser{
		ID:           7,
		Name:         "Administrator",
		Email:        "admin@example.com",
		PasswordHash: string(hash),
		AccountType:  "admin",
		Active:       true,
	}}, jwt)
	handler := controller.NewAuthController(usecase, false)
	router := shared.NewSetupRouter(t)
	router.Use(middleware.JWTAuth(usecase), middleware.CSRFProtection())
	router.POST("/auth/login", handler.Login)
	router.POST("/auth/refresh", handler.Refresh)
	router.POST("/auth/logout", handler.Logout)
	router.GET("/auth/me", handler.Me)
	return router
}

func TestAuthControllerFlow(t *testing.T) {
	router := newAuthRouter(t)
	login := performJSON(router, http.MethodPost, "/auth/login", `{"email":"ADMIN@example.com","password":"correct-password"}`, nil, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login status = %d, body = %s", login.Code, login.Body.String())
	}

	var loginBody api.AuthResponse
	if err := json.Unmarshal(login.Body.Bytes(), &loginBody); err != nil {
		t.Fatal(err)
	}
	if loginBody.Data.CsrfToken == "" {
		t.Fatal("login response did not include csrfToken")
	}
	cookies := login.Result().Cookies()
	jwtCookie := cookie(t, cookies, "jwt")
	refreshCookie := cookie(t, cookies, "refresh_token")
	csrfCookie := cookie(t, cookies, "csrf")
	if refreshCookie.Path != "/" || jwtCookie.HttpOnly != true {
		t.Fatalf("unexpected auth cookie configuration: jwt=%#v refresh=%#v", jwtCookie, refreshCookie)
	}

	me := performJSON(router, http.MethodGet, "/auth/me", "", []*http.Cookie{jwtCookie, csrfCookie}, csrfCookie.Value)
	if me.Code != http.StatusOK {
		t.Fatalf("me status = %d, body = %s", me.Code, me.Body.String())
	}
	var meBody api.LoginUser
	if err := json.Unmarshal(me.Body.Bytes(), &meBody); err != nil {
		t.Fatal(err)
	}
	if meBody.Id != 7 || meBody.Email != openapi_types.Email("admin@example.com") {
		t.Fatalf("unexpected current user: %#v", meBody)
	}

	refresh := performJSON(router, http.MethodPost, "/auth/refresh", "", []*http.Cookie{refreshCookie, csrfCookie}, csrfCookie.Value)
	if refresh.Code != http.StatusOK {
		t.Fatalf("refresh status = %d, body = %s", refresh.Code, refresh.Body.String())
	}

	logout := performJSON(router, http.MethodPost, "/auth/logout", "", []*http.Cookie{jwtCookie, csrfCookie}, csrfCookie.Value)
	if logout.Code != http.StatusNoContent {
		t.Fatalf("logout status = %d", logout.Code)
	}
}

func TestAuthControllerRejectsUnauthenticatedMe(t *testing.T) {
	response := performJSON(newAuthRouter(t), http.MethodGet, "/auth/me", "", nil, "")
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func performJSON(router http.Handler, method, path, body string, cookies []*http.Cookie, csrf string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	if csrf != "" {
		request.Header.Set("X-CSRF-Token", csrf)
	}
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func cookie(t *testing.T, cookies []*http.Cookie, name string) *http.Cookie {
	t.Helper()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	t.Fatalf("missing %s cookie", name)
	return nil
}

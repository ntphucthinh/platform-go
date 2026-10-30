package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"

	api "platform-go/generated/api"
	"platform-go/internal/common/ctxkeys"
	"platform-go/internal/pkg/types"
	"platform-go/internal/usecases"
)

type AuthController struct {
	usecase       *usecases.AuthUsecase
	secureCookies bool
}

func NewAuthController(usecase *usecases.AuthUsecase, secureCookies bool) *AuthController {
	return &AuthController{usecase: usecase, secureCookies: secureCookies}
}

func (ac *AuthController) Login(c *gin.Context) {
	var body api.LoginJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	email := strings.TrimSpace(strings.ToLower(string(body.Email)))
	if email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}
	tokens, err := ac.usecase.Login(c.Request.Context(), email, body.Password)
	if err != nil {
		ac.writeAuthenticationError(c, err, "invalid email or password")
		return
	}
	ac.setCookies(c, tokens)
	ac.writeAuthResponse(c, http.StatusOK, tokens, "Login successful")
}

func (ac *AuthController) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		ac.clearCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	tokens, err := ac.usecase.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		ac.clearCookies(c)
		ac.writeAuthenticationError(c, err, "invalid refresh token")
		return
	}
	ac.setCookies(c, tokens)
	ac.writeAuthResponse(c, http.StatusOK, tokens, "Token refreshed")
}

// Logout clears the client-held credentials. Tokens remain stateless and are not server-side revoked.
func (ac *AuthController) Logout(c *gin.Context) {
	ac.clearCookies(c)
	c.Status(http.StatusNoContent)
}

func (ac *AuthController) Me(c *gin.Context) {
	value, exists := c.Get(ctxkeys.LoginAdministratorKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	principal, ok := value.(types.LoginAdministrator)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	account, err := ac.usecase.CurrentUser(c.Request.Context(), strconv.FormatUint(uint64(principal.ID), 10))
	if err != nil {
		ac.writeAuthenticationError(c, err, "unauthorized")
		return
	}
	c.JSON(http.StatusOK, api.LoginUser{
		Id:          account.ID,
		Name:        account.Name,
		Email:       openapi_types.Email(account.Email),
		AccountType: account.AccountType,
	})
}

func (ac *AuthController) setCookies(c *gin.Context, tokens usecases.AuthTokens) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("jwt", tokens.AccessToken, 24*60*60, "/", "", true, true)
	c.SetCookie("refresh_token", tokens.RefreshToken, 7*24*60*60, "/", "", true, true)
	// Expire the former path-scoped cookie so a legacy /refresh request cannot prefer it.
	c.SetCookie("refresh_token", "", -1, "/refresh", "", true, true)
	// The client must read this non-HttpOnly value and mirror it in X-CSRF-Token.
	c.SetCookie("csrf", tokens.CSRFToken, 24*60*60, "/", "", true, false)
}

func (ac *AuthController) clearCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("jwt", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/refresh", "", true, true)
	c.SetCookie("csrf", "", -1, "/", "", true, false)
}

func (ac *AuthController) writeAuthResponse(c *gin.Context, status int, tokens usecases.AuthTokens, message string) {
	c.JSON(status, api.AuthResponse{
		Data: api.AuthResponseData{
			CsrfToken: tokens.CSRFToken,
		},
		Message: message,
	})
}

func (ac *AuthController) writeAuthenticationError(c *gin.Context, err error, message string) {
	if errors.Is(err, usecases.ErrInvalidCredentials) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication failed"})
}

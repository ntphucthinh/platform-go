package middleware

import (
	"crypto/subtle"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"platform-go/internal/common/ctxkeys"
	"platform-go/internal/pkg/types"
	"platform-go/internal/usecases"
)

// Endpoints that do not require an access token.
// Add new public endpoints here, grouped by HTTP method.
var accessTokenExemptions = map[string]map[string]struct{}{
	http.MethodPost: {
		"/auth/login":   {},
		"/login":        {},
		"/auth/refresh": {},
		"/refresh":      {},
	},
	http.MethodGet: {
		"/health":    {},
		"/homestays": {},
	},
}

// Endpoints that do not require CSRF validation.
// Refresh endpoints are intentionally excluded because they rotate cookies.
var csrfExemptions = map[string]map[string]struct{}{
	http.MethodPost: {
		"/auth/login": {},
		"/login":      {},
	},
}

// JWTAuth validates the access token and attaches the authenticated account
// to the request context.
func JWTAuth(auth *usecases.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path

		if isAccessTokenExempt(method, path) {
			c.Next()
			return
		}

		// Prefer the JWT cookie; fall back to the Authorization header.
		token, cookieErr := c.Cookie("jwt")

		if token == "" {
			token = bearerToken(c.GetHeader("Authorization"))
		}

		if token == "" {
			log.Printf(
				"[JWT] Missing token: method=%s path=%s cookieErr=%v",
				method,
				path,
				cookieErr,
			)
			unauthorized(c)
			return
		}

		claims, err := auth.ParseAccess(token)
		if err != nil {
			log.Printf(
				"[JWT] ParseAccess failed: method=%s path=%s err=%v",
				method,
				path,
				err,
			)
			unauthorized(c)
			return
		}

		account, err := auth.CurrentUser(
			c.Request.Context(),
			claims.Subject,
		)
		if err != nil {
			log.Printf(
				"[JWT] CurrentUser failed: method=%s path=%s err=%v",
				method,
				path,
				err,
			)
			unauthorized(c)
			return
		}

		c.Set(
			ctxkeys.LoginAdministratorKey,
			types.LoginAdministrator{
				ID:          uint(account.ID),
				Name:        account.Name,
				Email:       account.Email,
				AccountType: account.AccountType,
			},
		)

		c.Next()
	}
}

// CSRFProtection validates the CSRF cookie against the request header.
// Safe HTTP methods do not require CSRF validation.
func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path

		if isCSRFExempt(method, path) {
			c.Next()
			return
		}

		cookie, err := c.Cookie("csrf")
		header := c.GetHeader("X-CSRF-Token")

		if err != nil ||
			cookie == "" ||
			header == "" ||
			subtle.ConstantTimeCompare(
				[]byte(cookie),
				[]byte(header),
			) != 1 {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "CSRF token is invalid"},
			)
			return
		}

		c.Next()
	}
}

// bearerToken extracts a Bearer token from the Authorization header.
func bearerToken(header string) string {
	parts := strings.Fields(header)

	if len(parts) != 2 ||
		!strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return parts[1]
}

// isAccessTokenExempt checks whether the endpoint is public.
func isAccessTokenExempt(method, path string) bool {
	if isPathExempt(accessTokenExemptions, method, path) {
		return true
	}

	// Public homestay detail endpoints.
	return method == http.MethodGet &&
		strings.HasPrefix(path, "/homestays/")
}

// isCSRFExempt checks whether the endpoint can skip CSRF validation.
func isCSRFExempt(method, path string) bool {
	if method == http.MethodGet ||
		method == http.MethodHead ||
		method == http.MethodOptions {
		return true
	}

	return isPathExempt(csrfExemptions, method, path)
}

// isPathExempt checks an exact method and path match.
func isPathExempt(
	exemptions map[string]map[string]struct{},
	method, path string,
) bool {
	paths, exists := exemptions[method]
	if !exists {
		return false
	}

	_, exists = paths[path]
	return exists
}

// unauthorized returns a generic authentication error to the client.
func unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(
		http.StatusUnauthorized,
		gin.H{"error": "unauthorized"},
	)
}

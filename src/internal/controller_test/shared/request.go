package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/security"
	"platform-go/tests/factory"
)

// CreateLoginUser creates and inserts an active admin user in the test database.
func CreateLoginUser(modifiers ...func(*factory.UserFactory)) *dbmodel.User {
	email := fmt.Sprintf("user_%d@example.com", time.Now().UnixNano())
	f := factory.NewUserFactory().
		WithEmail(email).
		WithName("Test User").
		WithAccountTypeAdmin().
		WithActive()

	for _, mod := range modifiers {
		mod(f)
	}

	return f.MustInsert(TestDB)
}

// NewRequest creates an authenticated HTTP request for the given user.
func NewRequest(t *testing.T, user *dbmodel.User, method, path string, body any) *http.Request {
	t.Helper()
	var bodyReader io.Reader
	if body != nil {
		switch v := body.(type) {
		case io.Reader:
			bodyReader = v
		case string:
			bodyReader = strings.NewReader(v)
		case []byte:
			bodyReader = bytes.NewReader(v)
		default:
			data, err := json.Marshal(v)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}
			bodyReader = bytes.NewReader(data)
		}
	}

	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if user != nil {
		email := ""
		if user.Email != nil {
			email = *user.Email
		}
		token, err := testJWT.Create(security.Claims{
			Subject:     strconv.FormatInt(int64(user.ID), 10),
			Email:       email,
			AccountType: user.AccountType,
			TokenType:   security.AccessTokenType,
		}, 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to generate access token: %v", err)
		}

		csrf, err := security.NewCSRFToken()
		if err != nil {
			t.Fatalf("failed to generate csrf token: %v", err)
		}

		req.AddCookie(&http.Cookie{Name: "jwt", Value: token})
		req.AddCookie(&http.Cookie{Name: "csrf", Value: csrf})
		req.Header.Set("X-CSRF-Token", csrf)
	}

	return req
}

// NewUnAuthRequest creates an unauthenticated HTTP request.
func NewUnAuthRequest(t *testing.T, method, path string, body any) *http.Request {
	t.Helper()
	return NewRequest(t, nil, method, path, body)
}


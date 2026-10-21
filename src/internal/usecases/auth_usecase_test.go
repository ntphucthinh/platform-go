package usecases

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"platform-go/internal/domain/user"
	"platform-go/internal/infrastructure/security"
)

type authTestRepository struct {
	byEmail  user.AuthUser
	byID     user.AuthUser
	emailErr error
	idErr    error
}

func (r authTestRepository) CreateUser(context.Context, string) (user.User, error) {
	return user.User{}, errors.New("not implemented")
}
func (r authTestRepository) FindAuthUserByEmail(context.Context, string) (user.AuthUser, error) {
	return r.byEmail, r.emailErr
}
func (r authTestRepository) FindAuthUserByID(context.Context, int64) (user.AuthUser, error) {
	return r.byID, r.idErr
}

func TestAuthUsecaseLoginAndRefresh(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	account := user.AuthUser{ID: 7, Name: "Administrator", Email: "admin@example.com", PasswordHash: string(hash), AccountType: "admin", Active: true}
	jwt, err := security.NewJWTService("01234567890123456789012345678901")
	if err != nil {
		t.Fatal(err)
	}
	usecase := NewAuthUsecase(authTestRepository{byEmail: account, byID: account}, jwt)

	tokens, err := usecase.Login(context.Background(), account.Email, "correct-password")
	if err != nil {
		t.Fatal(err)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" || tokens.CSRFToken == "" {
		t.Fatal("login did not produce a complete token set")
	}
	if _, err := jwt.Parse(tokens.AccessToken, security.AccessTokenType); err != nil {
		t.Fatalf("access token is invalid: %v", err)
	}
	if _, err := jwt.Parse(tokens.RefreshToken, security.RefreshTokenType); err != nil {
		t.Fatalf("refresh token is invalid: %v", err)
	}

	refreshed, err := usecase.Refresh(context.Background(), tokens.RefreshToken)
	if err != nil {
		t.Fatal(err)
	}
	if refreshed.AccessToken == "" || refreshed.RefreshToken == "" || refreshed.CSRFToken == "" {
		t.Fatal("refresh did not produce a complete token set")
	}
	if _, err := usecase.Refresh(context.Background(), tokens.AccessToken); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("access token accepted as refresh token: %v", err)
	}
}

func TestAuthUsecaseRejectsInvalidOrInactiveAccounts(t *testing.T) {
	jwt, err := security.NewJWTService("01234567890123456789012345678901")
	if err != nil {
		t.Fatal(err)
	}
	inactive := user.AuthUser{ID: 7, AccountType: "admin", Active: false}
	usecase := NewAuthUsecase(authTestRepository{byEmail: inactive, byID: inactive}, jwt)

	if _, err := usecase.Login(context.Background(), "missing@example.com", "password"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("inactive login error = %v, want invalid credentials", err)
	}
	missingByEmail := NewAuthUsecase(authTestRepository{emailErr: gorm.ErrRecordNotFound}, jwt)
	if _, err := missingByEmail.Login(context.Background(), "missing@example.com", "password"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("missing login error = %v, want invalid credentials", err)
	}
	activeWithWrongPassword := user.AuthUser{ID: 7, PasswordHash: string([]byte("not-a-bcrypt-hash")), AccountType: "admin", Active: true}
	wrongPassword := NewAuthUsecase(authTestRepository{byEmail: activeWithWrongPassword}, jwt)
	if _, err := wrongPassword.Login(context.Background(), "admin@example.com", "wrong-password"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("wrong password error = %v, want invalid credentials", err)
	}
	refresh, err := jwt.Create(security.Claims{Subject: "7", TokenType: security.RefreshTokenType}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := usecase.Refresh(context.Background(), refresh); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("inactive refresh error = %v, want invalid credentials", err)
	}
	missing := NewAuthUsecase(authTestRepository{idErr: gorm.ErrRecordNotFound}, jwt)
	if _, err := missing.CurrentUser(context.Background(), "7"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("missing current user error = %v, want invalid credentials", err)
	}
}

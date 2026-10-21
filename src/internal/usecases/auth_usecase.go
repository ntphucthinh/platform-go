package usecases

import (
	"context"
	"errors"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"platform-go/internal/domain/user"
	"platform-go/internal/infrastructure/security"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

const (
	accessTokenLifetime  = 24 * time.Hour
	refreshTokenLifetime = 7 * 24 * time.Hour
)

type AuthUsecase struct {
	repo user.UserRepository
	jwt  *security.JWTService
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	CSRFToken    string
}

func NewAuthUsecase(repo user.UserRepository, jwt *security.JWTService) *AuthUsecase {
	return &AuthUsecase{repo: repo, jwt: jwt}
}

func (u *AuthUsecase) Login(ctx context.Context, email, password string) (AuthTokens, error) {
	account, err := u.repo.FindAuthUserByEmail(ctx, email)
	if err != nil || !isActiveAdministrator(account) || bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password)) != nil {
		return AuthTokens{}, ErrInvalidCredentials
	}
	return u.newTokens(account)
}

func (u *AuthUsecase) Refresh(ctx context.Context, refreshToken string) (AuthTokens, error) {
	claims, err := u.jwt.Parse(refreshToken, security.RefreshTokenType)
	if err != nil {
		return AuthTokens{}, ErrInvalidCredentials
	}
	id, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return AuthTokens{}, ErrInvalidCredentials
	}
	account, err := u.repo.FindAuthUserByID(ctx, id)
	if err != nil || !isActiveAdministrator(account) {
		return AuthTokens{}, ErrInvalidCredentials
	}
	return u.newTokens(account)
}

// CurrentUser reloads the authenticated account to enforce its active administrator status.
func (u *AuthUsecase) CurrentUser(ctx context.Context, subject string) (user.AuthUser, error) {
	id, err := strconv.ParseInt(subject, 10, 64)
	if err != nil {
		return user.AuthUser{}, ErrInvalidCredentials
	}
	account, err := u.repo.FindAuthUserByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) || !isActiveAdministrator(account) {
		return user.AuthUser{}, ErrInvalidCredentials
	}
	return account, err
}

func (u *AuthUsecase) ParseAccess(token string) (security.Claims, error) {
	return u.jwt.Parse(token, security.AccessTokenType)
}

func (u *AuthUsecase) newTokens(account user.AuthUser) (AuthTokens, error) {
	access, err := u.createToken(account, security.AccessTokenType, accessTokenLifetime)
	if err != nil {
		return AuthTokens{}, err
	}
	refresh, err := u.createToken(account, security.RefreshTokenType, refreshTokenLifetime)
	if err != nil {
		return AuthTokens{}, err
	}
	csrf, err := security.NewCSRFToken()
	if err != nil {
		return AuthTokens{}, err
	}
	return AuthTokens{AccessToken: access, RefreshToken: refresh, CSRFToken: csrf}, nil
}

func isActiveAdministrator(account user.AuthUser) bool {
	return account.Active && account.AccountType == "admin"
}

func (u *AuthUsecase) createToken(account user.AuthUser, kind string, duration time.Duration) (string, error) {
	return u.jwt.Create(security.Claims{Subject: strconv.FormatInt(account.ID, 10), Email: account.Email, AccountType: account.AccountType, TokenType: kind}, duration)
}

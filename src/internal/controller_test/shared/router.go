package shared

import (
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/di"
	"platform-go/internal/infrastructure/security"
	"platform-go/internal/middleware"
	"platform-go/internal/pkg/cvalidator"
)

var (
	TestDB          *gorm.DB
	TestRouter      *gin.Engine
	TestApplication *di.Application
	testJWT         *security.JWTService
)

// SetupTestEnvironment initializes the test database, DI container, and registers all API handlers onto TestRouter.
func SetupTestEnvironment(t *testing.T) {
	t.Helper()
	TestDB = SetupTestDB(t)

	var err error
	testJWT, err = security.NewJWTService("test-secret-with-at-least-thirty-two-characters")
	if err != nil {
		t.Fatalf("failed to init jwt service: %v", err)
	}

	TestApplication = di.InitializeApplication(TestDB, testJWT, false)

	cvalidator.NewValidate()
	v := cvalidator.Validator()
	cvalidator.CustomTagName(v)
	cvalidator.ValidatorTranslate(v)
	cvalidator.AddCustomValidation(v, TestApplication.DB)

	gin.SetMode(gin.TestMode)
	TestRouter = gin.New()
	TestRouter.Use(gin.Recovery())
	TestRouter.Use(middleware.JWTAuth(TestApplication.AuthUsecase))
	TestRouter.Use(middleware.CSRFProtection())

	api.RegisterHandlers(TestRouter, TestApplication.Controller)
}

// NewSetupRouter creates an isolated Gin router for controller HTTP tests.
// Deprecated: Use SetupTestEnvironment instead.
func NewSetupRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())
	return router
}

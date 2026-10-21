//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"platform-go/internal/infrastructure/queryservice"
	"platform-go/internal/infrastructure/repository"
	"platform-go/internal/infrastructure/security"
	"platform-go/internal/interfaces/controller"
	"platform-go/internal/usecases"
)

// Application contains the HTTP endpoint controller and middleware dependencies.
// It is assembled by the application's composition root.
type Application struct {
	Controller  *controller.BaseController
	AuthUsecase *usecases.AuthUsecase
}

func InitializeApplication(db *gorm.DB, jwtService *security.JWTService, secureCookies bool) *Application {
	wire.Build(
		repository.NewMySQLUserRepository,
		queryservice.NewUserQueryService,
		queryservice.NewHomestayQueryService,
		usecases.NewUserUsecase,
		usecases.NewAuthUsecase,
		controller.NewUserController,
		controller.NewAuthController,
		controller.NewHomestayController,
		controller.NewBaseController,
		wire.Struct(new(Application), "*"),
	)
	return nil
}

//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"platform-go/internal/infrastructure/cloudinary"
	"platform-go/internal/infrastructure/queryservice"
	"platform-go/internal/infrastructure/repository"
	"platform-go/internal/infrastructure/security"
	"platform-go/internal/interfaces/controller"
	"platform-go/internal/usecases"
	homestayqueryservice "platform-go/internal/usecases/homestays/queryservice"
	homestayusecase "platform-go/internal/usecases/homestays/usecase"
	"platform-go/internal/usecases/uploads"
)

// Application contains the HTTP endpoint controller and middleware dependencies.
// It is assembled by the application's composition root.
type Application struct {
	Controller  *controller.ServerController
	AuthUsecase *usecases.AuthUsecase
}

func InitializeApplication(db *gorm.DB, jwtService *security.JWTService, secureCookies bool) *Application {
	wire.Build(
		repository.NewMySQLUserRepository,
		queryservice.NewUserQueryService,
		queryservice.NewHomestayQueryService,
		cloudinary.NewClient,
		wire.Bind(new(uploads.AssetClient), new(*cloudinary.Client)),
		wire.Bind(new(homestayqueryservice.HomestayQueryService), new(*queryservice.HomestayQueryService)),
		usecases.NewUserUsecase,
		usecases.NewAuthUsecase,
		uploads.NewUseCase,
		homestayusecase.NewFetchHomestaysUseCase,
		controller.NewUserController,
		controller.NewAuthController,
		controller.NewHomestayController,
		controller.NewUploadController,
		controller.NewServerController,
		wire.Struct(new(Application), "*"),
	)
	return nil
}

//go:build wireinject
// +build wireinject

package di

import (
	"testing"

	"github.com/google/wire"
	"gorm.io/gorm"

	"platform-go/internal/infrastructure/queryservice"
	"platform-go/internal/infrastructure/repository"
	"platform-go/internal/interfaces/controller"
	"platform-go/internal/usecases"
	homestayqueryservice "platform-go/internal/usecases/homestays/queryservice"
	homestayusecase "platform-go/internal/usecases/homestays/usecase"
)

type TestHandlerSet struct {
	Handler         *controller.UserController
	HomestayHandler *controller.HomestayController
	DB              *gorm.DB
}

func TestInitializeControllers(t *testing.T, db *gorm.DB) (*TestHandlerSet, error) {
	wire.Build(
		repository.NewMySQLUserRepository,
		queryservice.NewUserQueryService,
		queryservice.NewHomestayQueryService,
		wire.Bind(new(homestayqueryservice.HomestayQueryService), new(*queryservice.HomestayQueryService)),
		usecases.NewUserUsecase,
		homestayusecase.NewFetchHomestaysUseCase,
		controller.NewUserController,
		controller.NewHomestayController,
		wire.Struct(new(TestHandlerSet), "*"),
	)
	return nil, nil
}

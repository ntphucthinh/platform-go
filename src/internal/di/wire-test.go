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
	travellocationqueryservice "platform-go/internal/usecases/travel_locations/queryservice"
	travellocationusecase "platform-go/internal/usecases/travel_locations/usecase"
)

type TestHandlerSet struct {
	Handler               *controller.UserController
	HomestayHandler       controller.HomestayController
	TravelLocationHandler controller.TravelLocationController
	DB                    *gorm.DB
}

func TestInitializeControllers(t *testing.T, db *gorm.DB) (*TestHandlerSet, error) {
	wire.Build(
		repository.NewSchemaDB,
		repository.NewMySQLUserRepository,
		repository.NewHomestayRepository,
		queryservice.NewUserQueryService,
		queryservice.NewHomestayQueryService,
		queryservice.NewTravelLocationQueryService,
		wire.Bind(new(homestayqueryservice.HomestayQueryService), new(*queryservice.HomestayQueryService)),
		wire.Bind(new(travellocationqueryservice.TravelLocationQueryService), new(*queryservice.TravelLocationQueryService)),
		usecases.NewUserUsecase,
		homestayusecase.NewFetchHomestaysUseCase,
		homestayusecase.NewCreateHomestayUsecase,
		homestayusecase.NewFetchHomestayByIDUseCase,
		homestayusecase.NewUpdateHomestayUsecase,
		homestayusecase.NewDeleteHomestayUsecase,
		travellocationusecase.NewFetchTravelLocationsUseCase,
		travellocationusecase.NewFetchTravelLocationByIDUseCase,
		controller.NewUserController,
		controller.NewHomestayController,
		controller.NewTravelLocationController,
		wire.Struct(new(TestHandlerSet), "*"),
	)
	return nil, nil
}

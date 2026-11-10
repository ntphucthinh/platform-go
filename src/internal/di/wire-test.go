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
	travelmemoryqueryservice "platform-go/internal/usecases/travel_memories/queryservice"
	travelmemoryusecase "platform-go/internal/usecases/travel_memories/usecase"
)

type TestHandlerSet struct {
	Handler               *controller.UserController
	HomestayHandler       controller.HomestayController
	TravelLocationHandler controller.TravelLocationController
	TravelMemoryHandler   controller.TravelMemoryController
	DB                    *gorm.DB
}

func TestInitializeControllers(t *testing.T, db *gorm.DB) (*TestHandlerSet, error) {
	wire.Build(
		repository.NewSchemaDB,
		repository.NewMySQLUserRepository,
		repository.NewHomestayRepository,
		repository.NewTravelLocationRepository,
		queryservice.NewUserQueryService,
		queryservice.NewHomestayQueryService,
		queryservice.NewTravelLocationQueryService,
		queryservice.NewTravelMemoryQueryService,
		wire.Bind(new(homestayqueryservice.HomestayQueryService), new(*queryservice.HomestayQueryService)),
		wire.Bind(new(travellocationqueryservice.TravelLocationQueryService), new(*queryservice.TravelLocationQueryService)),
		wire.Bind(new(travelmemoryqueryservice.TravelMemoryQueryService), new(*queryservice.TravelMemoryQueryService)),
		usecases.NewUserUsecase,
		homestayusecase.NewFetchHomestaysUseCase,
		homestayusecase.NewCreateHomestayUsecase,
		homestayusecase.NewFetchHomestayByIDUseCase,
		homestayusecase.NewUpdateHomestayUsecase,
		homestayusecase.NewDeleteHomestayUsecase,
		travellocationusecase.NewFetchTravelLocationsUsecase,
		travellocationusecase.NewCreateTravelLocationUsecase,
		travellocationusecase.NewFetchTravelLocationByIDUsecase,
		travellocationusecase.NewUpdateTravelLocationUsecase,
		travellocationusecase.NewDeleteTravelLocationUsecase,
		travelmemoryusecase.NewFetchTravelMemoriesUseCase,
		travelmemoryusecase.NewFetchTravelMemoryByIDUseCase,
		controller.NewUserController,
		controller.NewHomestayController,
		controller.NewTravelLocationController,
		controller.NewTravelMemoryController,
		wire.Struct(new(TestHandlerSet), "*"),
	)
	return nil, nil
}

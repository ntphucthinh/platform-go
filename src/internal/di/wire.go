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
	travellocationqueryservice "platform-go/internal/usecases/travel_locations/queryservice"
	travellocationusecase "platform-go/internal/usecases/travel_locations/usecase"
	travelmemoryqueryservice "platform-go/internal/usecases/travel_memories/queryservice"
	travelmemoryusecase "platform-go/internal/usecases/travel_memories/usecase"
	"platform-go/internal/usecases/uploads"
)

// Application contains the HTTP endpoint controller and middleware dependencies.
// It is assembled by the application's composition root.
type Application struct {
	Controller  *controller.ServerController
	AuthUsecase *usecases.AuthUsecase
	DB          repository.DBInterface
}

func InitializeApplication(db *gorm.DB, jwtService *security.JWTService, secureCookies bool) *Application {
	wire.Build(
		repository.NewSchemaDB,
		repository.NewMySQLUserRepository,
		repository.NewHomestayRepository,
		repository.NewTravelLocationRepository,
		queryservice.NewUserQueryService,
		queryservice.NewHomestayQueryService,
		queryservice.NewTravelLocationQueryService,
		queryservice.NewTravelMemoryQueryService,
		cloudinary.NewClient,
		wire.Bind(new(uploads.AssetClient), new(*cloudinary.Client)),
		wire.Bind(new(homestayqueryservice.HomestayQueryService), new(*queryservice.HomestayQueryService)),
		wire.Bind(new(travellocationqueryservice.TravelLocationQueryService), new(*queryservice.TravelLocationQueryService)),
		wire.Bind(new(travelmemoryqueryservice.TravelMemoryQueryService), new(*queryservice.TravelMemoryQueryService)),
		usecases.NewUserUsecase,
		usecases.NewAuthUsecase,
		uploads.NewUseCase,
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
		controller.NewAuthController,
		controller.NewHomestayController,
		controller.NewTravelLocationController,
		controller.NewTravelMemoryController,
		controller.NewUploadController,
		controller.NewServerController,
		wire.Struct(new(Application), "*"),
	)
	return nil
}

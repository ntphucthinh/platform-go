package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	api "platform-go/generated/api"
	"platform-go/internal/common"
	travelmemory "platform-go/internal/domain/travel_memory"
	travelmemoryday "platform-go/internal/domain/travel_memory_day"
	travelmemoryphoto "platform-go/internal/domain/travel_memory_photo"
	travelmemoryusecase "platform-go/internal/usecases/travel_memories/usecase"
)

type TravelMemoryController interface {
	FetchTravelMemories(c *gin.Context, params api.FetchTravelMemoriesParams)
	CreateTravelMemory(c *gin.Context)
	FetchTravelMemoryByID(c *gin.Context, travelMemoryID int64)
}

type travelMemoryController struct {
	fetchTravelMemoriesUseCase   *travelmemoryusecase.FetchTravelMemoriesUseCase
	createTravelMemoryUseCase    *travelmemoryusecase.CreateTravelMemoryUseCase
	fetchTravelMemoryByIDUseCase *travelmemoryusecase.FetchTravelMemoryByIDUseCase
}

func NewTravelMemoryController(
	fetchTravelMemoriesUseCase *travelmemoryusecase.FetchTravelMemoriesUseCase,
	createTravelMemoryUseCase *travelmemoryusecase.CreateTravelMemoryUseCase,
	fetchTravelMemoryByIDUseCase *travelmemoryusecase.FetchTravelMemoryByIDUseCase,
) TravelMemoryController {
	return &travelMemoryController{
		fetchTravelMemoriesUseCase:   fetchTravelMemoriesUseCase,
		createTravelMemoryUseCase:    createTravelMemoryUseCase,
		fetchTravelMemoryByIDUseCase: fetchTravelMemoryByIDUseCase,
	}
}

// Get travel memory list
// (GET /travel-memories)
func (con *travelMemoryController) FetchTravelMemories(c *gin.Context, params api.FetchTravelMemoriesParams) {
	if params.CountryCode == nil {
		if cc := c.Query("country_code"); cc != "" {
			params.CountryCode = &cc
		}
	}

	res, err := con.fetchTravelMemoriesUseCase.Execute(c.Request.Context(), params.Page, params.CountryCode, params.LocationId)
	if err != nil {
		common.IsErrorWithMessage(c, err, "Failed to fetch travel memories")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Create travel memory
// (POST /travel-memories)
func (con *travelMemoryController) CreateTravelMemory(c *gin.Context) {
	var request api.CreateTravelMemoryRequest
	if err := common.BindJSON(c, &request); err != nil {
		return
	}

	res, err := con.createTravelMemoryUseCase.Execute(c.Request.Context(), request)
	if err != nil {
		if errors.Is(err, travelmemory.ErrInvalidDateRange) ||
			errors.Is(err, travelmemoryday.ErrInvalidDayDate) ||
			errors.Is(err, travelmemoryday.ErrDuplicateDayDate) ||
			errors.Is(err, travelmemory.ErrInvalidImagePath) ||
			errors.Is(err, travelmemoryphoto.ErrInvalidImagePath) ||
			errors.Is(err, travelmemoryphoto.ErrDuplicateImagePath) ||
			errors.Is(err, travelmemoryphoto.ErrInvalidSortOrder) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to create travel memory")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get travel memory detail
// (GET /travel-memories/{travelMemoryID})
func (con *travelMemoryController) FetchTravelMemoryByID(c *gin.Context, travelMemoryID int64) {
	res, err := con.fetchTravelMemoryByIDUseCase.Execute(c.Request.Context(), travelMemoryID)
	if err != nil {
		if errors.Is(err, travelmemoryusecase.ErrTravelMemoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Travel memory not found"})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to fetch travel memory detail")
		return
	}

	c.JSON(http.StatusOK, res)
}

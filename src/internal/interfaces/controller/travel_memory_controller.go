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
	UpdateTravelMemory(c *gin.Context, travelMemoryID int64)
	DeleteTravelMemory(c *gin.Context, travelMemoryID int64)
}

type travelMemoryController struct {
	fetchTravelMemoriesUseCase   *travelmemoryusecase.FetchTravelMemoriesUseCase
	createTravelMemoryUseCase    *travelmemoryusecase.CreateTravelMemoryUseCase
	fetchTravelMemoryByIDUseCase *travelmemoryusecase.FetchTravelMemoryByIDUseCase
	updateTravelMemoryUseCase    *travelmemoryusecase.UpdateTravelMemoryUseCase
	deleteTravelMemoryUseCase    *travelmemoryusecase.DeleteTravelMemoryUseCase
}

func NewTravelMemoryController(
	fetchTravelMemoriesUseCase *travelmemoryusecase.FetchTravelMemoriesUseCase,
	createTravelMemoryUseCase *travelmemoryusecase.CreateTravelMemoryUseCase,
	fetchTravelMemoryByIDUseCase *travelmemoryusecase.FetchTravelMemoryByIDUseCase,
	updateTravelMemoryUseCase *travelmemoryusecase.UpdateTravelMemoryUseCase,
	deleteTravelMemoryUseCase *travelmemoryusecase.DeleteTravelMemoryUseCase,
) TravelMemoryController {
	return &travelMemoryController{
		fetchTravelMemoriesUseCase:   fetchTravelMemoriesUseCase,
		createTravelMemoryUseCase:    createTravelMemoryUseCase,
		fetchTravelMemoryByIDUseCase: fetchTravelMemoryByIDUseCase,
		updateTravelMemoryUseCase:    updateTravelMemoryUseCase,
		deleteTravelMemoryUseCase:    deleteTravelMemoryUseCase,
	}
}

// Get travel memory list
// (GET /travel-memories)
func (con *travelMemoryController) FetchTravelMemories(c *gin.Context, params api.FetchTravelMemoriesParams) {
	if params.Year != nil && (*params.Year < 1000 || *params.Year > 9999) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year parameter"})
		return
	}

	if params.CountryCode == nil {
		if cc := c.Query("country_code"); cc != "" {
			params.CountryCode = &cc
		}
	}

	res, err := con.fetchTravelMemoriesUseCase.Execute(c.Request.Context(), params.Page, params.CountryCode, params.LocationId, params.Year)
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

// Update travel memory
// (PUT /travel-memories/{travelMemoryID})
func (con *travelMemoryController) UpdateTravelMemory(c *gin.Context, travelMemoryID int64) {
	var request api.UpdateTravelMemoryRequest
	if err := common.BindJSON(c, &request); err != nil {
		return
	}

	res, err := con.updateTravelMemoryUseCase.Execute(c.Request.Context(), travelMemoryID, request)
	if err != nil {
		if errors.Is(err, travelmemoryusecase.ErrTravelMemoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Travel memory not found"})
			return
		}
		if errors.Is(err, travelmemoryusecase.ErrUpdateTravelMemoryInvalidID) ||
			errors.Is(err, travelmemoryusecase.ErrUpdateTravelMemoryInvalidDayID) ||
			errors.Is(err, travelmemoryusecase.ErrUpdateTravelMemoryDuplicateDayID) ||
			errors.Is(err, travelmemoryusecase.ErrUpdateTravelMemoryInvalidPhotoID) ||
			errors.Is(err, travelmemoryusecase.ErrUpdateTravelMemoryDuplicatePhotoID) ||
			errors.Is(err, travelmemory.ErrInvalidDateRange) ||
			errors.Is(err, travelmemoryday.ErrInvalidDayDate) ||
			errors.Is(err, travelmemoryday.ErrDuplicateDayDate) ||
			errors.Is(err, travelmemory.ErrInvalidImagePath) ||
			errors.Is(err, travelmemoryphoto.ErrInvalidImagePath) ||
			errors.Is(err, travelmemoryphoto.ErrInvalidSortOrder) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to update travel memory")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Delete travel memory
// (DELETE /travel-memories/{travelMemoryID})
func (con *travelMemoryController) DeleteTravelMemory(c *gin.Context, travelMemoryID int64) {
	if travelMemoryID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Travel memory ID must be greater than zero"})
		return
	}

	res, err := con.deleteTravelMemoryUseCase.Execute(c.Request.Context(), travelMemoryID)
	if err != nil {
		if errors.Is(err, travelmemoryusecase.ErrTravelMemoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Travel memory not found"})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to delete travel memory")
		return
	}

	c.JSON(http.StatusOK, res)
}

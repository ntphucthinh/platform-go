package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	api "platform-go/generated/api"
	"platform-go/internal/common"
	"platform-go/internal/domain/itinerary"
	itineraryactivity "platform-go/internal/domain/itinerary_activity"
	itineraryday "platform-go/internal/domain/itinerary_day"
	itinerarytheme "platform-go/internal/domain/itinerary_theme"
	itineraryusecase "platform-go/internal/usecases/itineraries/usecase"
)

type ItineraryController interface {
	FetchItineraries(c *gin.Context, params api.FetchItinerariesParams)
	FetchItineraryByID(c *gin.Context, itineraryID int64)
	CreateItinerary(c *gin.Context)
	UpdateItinerary(c *gin.Context, itineraryID int64)
	DeleteItinerary(c *gin.Context, itineraryID int64)
}

type itineraryController struct {
	fetchItinerariesUseCase   *itineraryusecase.FetchItinerariesUseCase
	fetchItineraryByIDUseCase *itineraryusecase.FetchItineraryByIDUseCase
	createItineraryUseCase    *itineraryusecase.CreateItineraryUseCase
	updateItineraryUseCase    *itineraryusecase.UpdateItineraryUseCase
	deleteItineraryUseCase    *itineraryusecase.DeleteItineraryUseCase
}

func NewItineraryController(
	fetchItinerariesUseCase *itineraryusecase.FetchItinerariesUseCase,
	fetchItineraryByIDUseCase *itineraryusecase.FetchItineraryByIDUseCase,
	createItineraryUseCase *itineraryusecase.CreateItineraryUseCase,
	updateItineraryUseCase *itineraryusecase.UpdateItineraryUseCase,
	deleteItineraryUseCase *itineraryusecase.DeleteItineraryUseCase,
) ItineraryController {
	return &itineraryController{
		fetchItinerariesUseCase:   fetchItinerariesUseCase,
		fetchItineraryByIDUseCase: fetchItineraryByIDUseCase,
		createItineraryUseCase:    createItineraryUseCase,
		updateItineraryUseCase:    updateItineraryUseCase,
		deleteItineraryUseCase:    deleteItineraryUseCase,
	}
}

// Get itinerary list
// (GET /itineraries)
func (con *itineraryController) FetchItineraries(c *gin.Context, params api.FetchItinerariesParams) {
	res, err := con.fetchItinerariesUseCase.Execute(c.Request.Context(), params.Page, params.PerPage, params.LocationId)
	if err != nil {
		common.IsErrorWithMessage(c, err, "Failed to fetch itineraries")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get itinerary detail
// (GET /itineraries/{itineraryID})
func (con *itineraryController) FetchItineraryByID(c *gin.Context, itineraryID int64) {
	res, err := con.fetchItineraryByIDUseCase.Execute(c.Request.Context(), itineraryID)
	if err != nil {
		if errors.Is(err, itineraryusecase.ErrItineraryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Itinerary not found"})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to fetch itinerary detail")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Create itinerary
// (POST /itineraries)
func (con *itineraryController) CreateItinerary(c *gin.Context) {
	var request api.CreateItineraryRequest
	if err := common.BindJSON(c, &request); err != nil {
		return
	}

	res, err := con.createItineraryUseCase.Execute(c.Request.Context(), request)
	if err != nil {
		if errors.Is(err, itinerary.ErrInvalidImagePath) ||
			errors.Is(err, itinerarytheme.ErrDuplicateThemeID) ||
			errors.Is(err, itinerarytheme.ErrInvalidThemeID) ||
			errors.Is(err, itineraryday.ErrDayNotFound) ||
			errors.Is(err, itineraryday.ErrInvalidDayNumber) ||
			errors.Is(err, itineraryday.ErrDuplicateDayNumber) ||
			errors.Is(err, itineraryday.ErrInvalidImagePath) ||
			errors.Is(err, itineraryactivity.ErrInvalidStartTime) ||
			errors.Is(err, itineraryactivity.ErrInvalidDurationMinutes) ||
			errors.Is(err, itineraryactivity.ErrInvalidImagePath) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to create itinerary")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Update itinerary
// (PUT /itineraries/{itineraryID})
func (con *itineraryController) UpdateItinerary(c *gin.Context, itineraryID int64) {
	var request api.UpdateItineraryRequest
	if err := common.BindJSON(c, &request); err != nil {
		return
	}

	res, err := con.updateItineraryUseCase.Execute(c.Request.Context(), itineraryID, request)
	if err != nil {
		if errors.Is(err, itineraryusecase.ErrItineraryNotFound) ||
			errors.Is(err, itinerary.ErrItineraryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Itinerary not found"})
			return
		}
		if errors.Is(err, itineraryusecase.ErrUpdateItineraryInvalidID) ||
			errors.Is(err, itineraryusecase.ErrUpdateItineraryInvalidDayID) ||
			errors.Is(err, itineraryusecase.ErrUpdateItineraryDuplicateDayID) ||
			errors.Is(err, itineraryusecase.ErrUpdateItineraryInvalidActivityID) ||
			errors.Is(err, itineraryusecase.ErrUpdateItineraryDuplicateActivityID) ||
			errors.Is(err, itinerary.ErrInvalidImagePath) ||
			errors.Is(err, itinerarytheme.ErrDuplicateThemeID) ||
			errors.Is(err, itinerarytheme.ErrInvalidThemeID) ||
			errors.Is(err, itineraryday.ErrDayNotFound) ||
			errors.Is(err, itineraryday.ErrInvalidDayNumber) ||
			errors.Is(err, itineraryday.ErrDuplicateDayNumber) ||
			errors.Is(err, itineraryday.ErrInvalidImagePath) ||
			errors.Is(err, itineraryactivity.ErrActivityNotFound) ||
			errors.Is(err, itineraryactivity.ErrInvalidStartTime) ||
			errors.Is(err, itineraryactivity.ErrInvalidDurationMinutes) ||
			errors.Is(err, itineraryactivity.ErrInvalidImagePath) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to update itinerary")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Delete itinerary
// (DELETE /itineraries/{itineraryID})
func (con *itineraryController) DeleteItinerary(c *gin.Context, itineraryID int64) {
	if itineraryID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Itinerary ID must be greater than zero"})
		return
	}

	res, err := con.deleteItineraryUseCase.Execute(c.Request.Context(), itineraryID)
	if err != nil {
		if errors.Is(err, itineraryusecase.ErrItineraryNotFound) ||
			errors.Is(err, itinerary.ErrItineraryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Itinerary not found"})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to delete itinerary")
		return
	}

	c.JSON(http.StatusOK, res)
}

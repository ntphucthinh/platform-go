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
}

type itineraryController struct {
	fetchItinerariesUseCase   *itineraryusecase.FetchItinerariesUseCase
	fetchItineraryByIDUseCase *itineraryusecase.FetchItineraryByIDUseCase
	createItineraryUseCase    *itineraryusecase.CreateItineraryUseCase
}

func NewItineraryController(
	fetchItinerariesUseCase *itineraryusecase.FetchItinerariesUseCase,
	fetchItineraryByIDUseCase *itineraryusecase.FetchItineraryByIDUseCase,
	createItineraryUseCase *itineraryusecase.CreateItineraryUseCase,
) ItineraryController {
	return &itineraryController{
		fetchItinerariesUseCase:   fetchItinerariesUseCase,
		fetchItineraryByIDUseCase: fetchItineraryByIDUseCase,
		createItineraryUseCase:    createItineraryUseCase,
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

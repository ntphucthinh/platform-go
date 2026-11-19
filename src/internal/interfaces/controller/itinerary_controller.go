package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	api "platform-go/generated/api"
	"platform-go/internal/common"
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
}

func NewItineraryController(
	fetchItinerariesUseCase *itineraryusecase.FetchItinerariesUseCase,
	fetchItineraryByIDUseCase *itineraryusecase.FetchItineraryByIDUseCase,
) ItineraryController {
	return &itineraryController{
		fetchItinerariesUseCase:   fetchItinerariesUseCase,
		fetchItineraryByIDUseCase: fetchItineraryByIDUseCase,
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

	c.JSON(http.StatusOK, api.MessageOnlyResponse{Message: "Success"})
}

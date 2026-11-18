package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "platform-go/generated/api"
	"platform-go/internal/common"
	itineraryusecase "platform-go/internal/usecases/itineraries/usecase"
)

type ItineraryController interface {
	FetchItineraries(c *gin.Context, params api.FetchItinerariesParams)
}

type itineraryController struct {
	fetchItinerariesUseCase *itineraryusecase.FetchItinerariesUseCase
}

func NewItineraryController(fetchItinerariesUseCase *itineraryusecase.FetchItinerariesUseCase) ItineraryController {
	return &itineraryController{
		fetchItinerariesUseCase: fetchItinerariesUseCase,
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

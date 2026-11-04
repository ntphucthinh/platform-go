package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	api "platform-go/generated/api"
	"platform-go/internal/common"
	travellocationusecase "platform-go/internal/usecases/travel_locations/usecase"
)

type TravelLocationController interface {
	FetchTravelLocations(c *gin.Context, params api.FetchTravelLocationsParams)
	FetchTravelLocationByID(c *gin.Context, travelLocationID int64)
}

type travelLocationController struct {
	fetchTravelLocationsUseCase    *travellocationusecase.FetchTravelLocationsUseCase
	fetchTravelLocationByIDUseCase *travellocationusecase.FetchTravelLocationByIDUseCase
}

func NewTravelLocationController(
	fetchTravelLocationsUseCase *travellocationusecase.FetchTravelLocationsUseCase,
	fetchTravelLocationByIDUseCase *travellocationusecase.FetchTravelLocationByIDUseCase,
) TravelLocationController {
	return &travelLocationController{
		fetchTravelLocationsUseCase:    fetchTravelLocationsUseCase,
		fetchTravelLocationByIDUseCase: fetchTravelLocationByIDUseCase,
	}
}

// Get travel location list
// (GET /travel-locations)
func (con *travelLocationController) FetchTravelLocations(c *gin.Context, params api.FetchTravelLocationsParams) {
	if params.Status != nil && *params.Status != api.TravelLocationStatusVisited && *params.Status != api.TravelLocationStatusWishlist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status parameter"})
		return
	}

	res, err := con.fetchTravelLocationsUseCase.Execute(c.Request.Context(), params)
	if err != nil {
		common.IsErrorWithMessage(c, err, "Failed to fetch travel locations")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get travel location detail
// (GET /travel-locations/{travelLocationID})
func (con *travelLocationController) FetchTravelLocationByID(c *gin.Context, travelLocationID int64) {
	res, err := con.fetchTravelLocationByIDUseCase.Execute(c.Request.Context(), travelLocationID)
	if err != nil {
		if errors.Is(err, travellocationusecase.ErrTravelLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Travel location not found"})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to fetch travel location detail")
		return
	}

	c.JSON(http.StatusOK, res)
}


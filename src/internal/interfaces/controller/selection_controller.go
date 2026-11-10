package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "platform-go/generated/api"
	"platform-go/internal/common"
	selectionusecase "platform-go/internal/usecases/selections/usecase"
)

type SelectionController interface {
	FetchSelectionTravelLocations(c *gin.Context, params api.FetchSelectionTravelLocationsParams)
}

type selectionController struct {
	fetchSelectionTravelLocationsUseCase *selectionusecase.FetchSelectionTravelLocationsUseCase
}

func NewSelectionController(
	fetchSelectionTravelLocationsUseCase *selectionusecase.FetchSelectionTravelLocationsUseCase,
) SelectionController {
	return &selectionController{
		fetchSelectionTravelLocationsUseCase: fetchSelectionTravelLocationsUseCase,
	}
}

// Get travel locations for selection
// (GET /selections/travel-locations)
func (con *selectionController) FetchSelectionTravelLocations(c *gin.Context, params api.FetchSelectionTravelLocationsParams) {
	if params.CountryCode == nil {
		if cc := c.Query("country_code"); cc != "" {
			params.CountryCode = &cc
		}
	}

	res, err := con.fetchSelectionTravelLocationsUseCase.Execute(c.Request.Context(), params.CountryCode)
	if err != nil {
		common.IsErrorWithMessage(c, err, "Failed to fetch selection travel locations")
		return
	}

	c.JSON(http.StatusOK, res)
}


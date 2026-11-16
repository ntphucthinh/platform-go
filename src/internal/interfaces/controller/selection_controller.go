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
	FetchSelectionTravelThemes(c *gin.Context)
}

type selectionController struct {
	fetchSelectionTravelLocationsUseCase *selectionusecase.FetchSelectionTravelLocationsUseCase
	fetchSelectionTravelThemesUseCase    *selectionusecase.FetchSelectionTravelThemesUseCase
}

func NewSelectionController(
	fetchSelectionTravelLocationsUseCase *selectionusecase.FetchSelectionTravelLocationsUseCase,
	fetchSelectionTravelThemesUseCase *selectionusecase.FetchSelectionTravelThemesUseCase,
) SelectionController {
	return &selectionController{
		fetchSelectionTravelLocationsUseCase: fetchSelectionTravelLocationsUseCase,
		fetchSelectionTravelThemesUseCase:    fetchSelectionTravelThemesUseCase,
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

// Get travel themes for selection
// (GET /selections/travel-themes)
func (con *selectionController) FetchSelectionTravelThemes(c *gin.Context) {
	res, err := con.fetchSelectionTravelThemesUseCase.Execute(c.Request.Context())
	if err != nil {
		common.IsErrorWithMessage(c, err, "Failed to fetch selection travel themes")
		return
	}

	c.JSON(http.StatusOK, res)
}

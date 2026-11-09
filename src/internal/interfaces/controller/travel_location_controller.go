package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	api "platform-go/generated/api"
	"platform-go/internal/common"
	travellocation "platform-go/internal/domain/travel_location"
	travellocationusecase "platform-go/internal/usecases/travel_locations/usecase"
)

type TravelLocationController interface {
	FetchTravelLocations(c *gin.Context, params api.FetchTravelLocationsParams)
	CreateTravelLocation(c *gin.Context)
	FetchTravelLocationByID(c *gin.Context, travelLocationID int64)
	UpdateTravelLocation(c *gin.Context, travelLocationID int64)
	DeleteTravelLocation(c *gin.Context, travelLocationID int64)
}

type travelLocationController struct {
	fetchTravelLocationsUsecase    *travellocationusecase.FetchTravelLocationsUsecase
	createTravelLocationUsecase    *travellocationusecase.CreateTravelLocationUsecase
	fetchTravelLocationByIDUsecase *travellocationusecase.FetchTravelLocationByIDUsecase
	updateTravelLocationUsecase    *travellocationusecase.UpdateTravelLocationUsecase
	deleteTravelLocationUsecase    *travellocationusecase.DeleteTravelLocationUsecase
}

func NewTravelLocationController(
	fetchTravelLocationsUsecase *travellocationusecase.FetchTravelLocationsUsecase,
	createTravelLocationUsecase *travellocationusecase.CreateTravelLocationUsecase,
	fetchTravelLocationByIDUsecase *travellocationusecase.FetchTravelLocationByIDUsecase,
	updateTravelLocationUsecase *travellocationusecase.UpdateTravelLocationUsecase,
	deleteTravelLocationUsecase *travellocationusecase.DeleteTravelLocationUsecase,
) TravelLocationController {
	return &travelLocationController{
		fetchTravelLocationsUsecase:    fetchTravelLocationsUsecase,
		createTravelLocationUsecase:    createTravelLocationUsecase,
		fetchTravelLocationByIDUsecase: fetchTravelLocationByIDUsecase,
		updateTravelLocationUsecase:    updateTravelLocationUsecase,
		deleteTravelLocationUsecase:    deleteTravelLocationUsecase,
	}
}

// Get travel location list
// (GET /travel-locations)
func (con *travelLocationController) FetchTravelLocations(c *gin.Context, params api.FetchTravelLocationsParams) {
	if params.Status != nil && *params.Status != api.TravelLocationStatusVisited && *params.Status != api.TravelLocationStatusWishlist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status parameter"})
		return
	}

	if params.CountryCode == nil {
		if cc := c.Query("country_code"); cc != "" {
			params.CountryCode = &cc
		}
	}
	if params.LocationName == nil {
		if ln := c.Query("name"); ln != "" {
			params.LocationName = &ln
		}
	}

	res, err := con.fetchTravelLocationsUsecase.Execute(c.Request.Context(), params)
	if err != nil {
		common.IsErrorWithMessage(c, err, "Failed to fetch travel locations")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Create travel location
// (POST /travel-locations)
func (con *travelLocationController) CreateTravelLocation(c *gin.Context) {
	var request api.CreateTravelLocationRequest
	if err := common.BindJSON(c, &request); err != nil {
		return
	}

	res, err := con.createTravelLocationUsecase.Execute(c.Request.Context(), request)
	if err != nil {
		if errors.Is(err, travellocation.ErrInvalidCountryCode) ||
			errors.Is(err, travellocation.ErrInvalidCoordinates) ||
			errors.Is(err, travellocation.ErrInvalidImagePath) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to create travel location")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get travel location detail
// (GET /travel-locations/{travelLocationID})
func (con *travelLocationController) FetchTravelLocationByID(c *gin.Context, travelLocationID int64) {
	res, err := con.fetchTravelLocationByIDUsecase.Execute(c.Request.Context(), travelLocationID)
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

// Update travel location
// (PUT /travel-locations/{travelLocationID})
func (con *travelLocationController) UpdateTravelLocation(c *gin.Context, travelLocationID int64) {
	var request api.UpdateTravelLocationRequest
	if err := common.BindJSON(c, &request); err != nil {
		return
	}

	res, err := con.updateTravelLocationUsecase.Execute(c.Request.Context(), travelLocationID, request)
	if err != nil {
		if errors.Is(err, travellocationusecase.ErrTravelLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Travel location not found"})
			return
		}
		if errors.Is(err, travellocation.ErrInvalidCountryCode) ||
			errors.Is(err, travellocation.ErrInvalidCoordinates) ||
			errors.Is(err, travellocation.ErrInvalidImagePath) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to update travel location")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Delete travel location
// (DELETE /travel-locations/{travelLocationID})
func (con *travelLocationController) DeleteTravelLocation(c *gin.Context, travelLocationID int64) {
	res, err := con.deleteTravelLocationUsecase.Execute(c.Request.Context(), travelLocationID)
	if err != nil {
		if errors.Is(err, travellocationusecase.ErrTravelLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Travel location not found"})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to delete travel location")
		return
	}

	c.JSON(http.StatusOK, res)
}

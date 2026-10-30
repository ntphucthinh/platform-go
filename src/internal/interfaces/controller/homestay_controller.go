package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	api "platform-go/generated/api"
	"platform-go/internal/common"
	homestayusecase "platform-go/internal/usecases/homestays/usecase"
)

type HomestayController interface {
	FetchHomestays(c *gin.Context, params api.FetchHomestaysParams)
	CreateHomestay(c *gin.Context)
	FetchHomestayByID(c *gin.Context, homestayID int64)
	UpdateHomestay(c *gin.Context, homestayID int64)
	DeleteHomestay(c *gin.Context, homestayID int64)
}

type homestayController struct {
	fetchHomestaysUseCase    *homestayusecase.FetchHomestaysUseCase
	createHomestaysUseCase   *homestayusecase.CreateHomestayUsecase
	fetchHomestayByIDUseCase *homestayusecase.FetchHomestayByIDUseCase
	updateHomestayUseCase    *homestayusecase.UpdateHomestayUsecase
	deleteHomestayUseCase    *homestayusecase.DeleteHomestayUsecase
}

func NewHomestayController(
	fetchHomestaysUseCase *homestayusecase.FetchHomestaysUseCase,
	createHomestaysUseCase *homestayusecase.CreateHomestayUsecase,
	fetchHomestayByIDUseCase *homestayusecase.FetchHomestayByIDUseCase,
	updateHomestayUseCase *homestayusecase.UpdateHomestayUsecase,
	deleteHomestayUseCase *homestayusecase.DeleteHomestayUsecase,
) HomestayController {
	return &homestayController{
		fetchHomestaysUseCase:    fetchHomestaysUseCase,
		createHomestaysUseCase:   createHomestaysUseCase,
		fetchHomestayByIDUseCase: fetchHomestayByIDUseCase,
		updateHomestayUseCase:    updateHomestayUseCase,
		deleteHomestayUseCase:    deleteHomestayUseCase,
	}
}

// Get homestay list
// (GET /homestays)
func (con *homestayController) FetchHomestays(c *gin.Context, params api.FetchHomestaysParams) {
	res, err := con.fetchHomestaysUseCase.Execute(c.Request.Context(), params.Page, params.PerPage)
	if err != nil {
		common.IsErrorWithMessage(c, err, "Failed to fetch homestays")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Create homestay
// (POST /homestays)
func (con *homestayController) CreateHomestay(c *gin.Context) {
	var request api.CreateHomestayRequest
	if err := common.BindJSON(c, &request); err != nil {
		return
	}
	res, err := con.createHomestaysUseCase.Execute(c.Request.Context(), request)
	if err != nil {
		if errors.Is(err, homestayusecase.ErrInvalidHomestayInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to create homestay")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get homestay detail
// (GET /homestays/{homestayID})
func (con *homestayController) FetchHomestayByID(c *gin.Context, homestayID int64) {
	res, err := con.fetchHomestayByIDUseCase.Execute(c.Request.Context(), homestayID)
	if err != nil {
		if errors.Is(err, homestayusecase.ErrHomestayNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Homestay not found"})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to fetch homestay detail")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Update homestay
// (PUT /homestays/{homestayID})
func (con *homestayController) UpdateHomestay(c *gin.Context, homestayID int64) {
	var request api.UpdateHomestayRequest
	if err := common.BindJSON(c, &request); err != nil {
		return
	}
	res, err := con.updateHomestayUseCase.Execute(c.Request.Context(), homestayID, request)
	if err != nil {
		if errors.Is(err, homestayusecase.ErrHomestayNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Homestay not found"})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to update homestay")
		return
	}

	c.JSON(http.StatusOK, res)
}

// Delete homestay
// (DELETE /homestays/{homestayID})
func (con *homestayController) DeleteHomestay(c *gin.Context, homestayID int64) {
	res, err := con.deleteHomestayUseCase.Execute(c.Request.Context(), homestayID)
	if err != nil {
		if errors.Is(err, homestayusecase.ErrHomestayNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Homestay not found"})
			return
		}
		common.IsErrorWithMessage(c, err, "Failed to delete homestay")
		return
	}

	c.JSON(http.StatusOK, res)
}

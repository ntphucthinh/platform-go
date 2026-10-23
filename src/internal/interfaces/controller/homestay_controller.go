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
	FetchHomestays(c *gin.Context)
	CreateHomestay(c *gin.Context)
	// FetchHomestayByID(c *gin.Context, homestayID int64)
	// UpdateHomestay(c *gin.Context, homestayID int64)
}

type homestayController struct {
	fetchHomestaysUseCase  *homestayusecase.FetchHomestaysUseCase
	createHomestaysUseCase *homestayusecase.CreateHomestayUsecase
}

func NewHomestayController(
	fetchHomestaysUseCase *homestayusecase.FetchHomestaysUseCase,
	createHomestaysUseCase *homestayusecase.CreateHomestayUsecase,
) HomestayController {
	return &homestayController{
		fetchHomestaysUseCase:  fetchHomestaysUseCase,
		createHomestaysUseCase: createHomestaysUseCase,
	}
}

// Get homestay list
// (GET /homestays)
func (con *homestayController) FetchHomestays(c *gin.Context) {
	items, err := con.fetchHomestaysUseCase.Execute(c.Request.Context())
	if err != nil {
		common.IsErrorWithMessage(c, err, "Failed to fetch homestays")
		return
	}

	c.JSON(http.StatusOK, api.FetchHomestays200Response{
		Data: items,
	})
}

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

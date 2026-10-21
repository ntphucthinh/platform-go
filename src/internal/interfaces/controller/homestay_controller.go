package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "platform-go/generated/api"
	"platform-go/internal/common"
	homestayusecase "platform-go/internal/usecases/homestays/usecase"
)

type HomestayController struct {
	usecase *homestayusecase.FetchHomestaysUseCase
}

func NewHomestayController(usecase *homestayusecase.FetchHomestaysUseCase) *HomestayController {
	return &HomestayController{
		usecase: usecase,
	}
}

// Get homestay list
// (GET /homestays)
func (hc *HomestayController) FetchHomestays(c *gin.Context) {
	items, err := hc.usecase.Execute(c.Request.Context())
	if err != nil {
		common.IsErrorWithMessage(c, err, "Failed to fetch homestays")
		return
	}

	c.JSON(http.StatusOK, api.FetchHomestays200Response{
		Data: items,
	})
}

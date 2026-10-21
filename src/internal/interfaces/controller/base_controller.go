package controller

import (
	"github.com/gin-gonic/gin"

	api "platform-go/generated/api"
)

// BaseController is the HTTP composition root for generated API endpoints.
// It delegates each endpoint to the controller that owns its business workflow.
type BaseController struct {
	users     *UserController
	auth      *AuthController
	homestays *HomestayController
}

func NewBaseController(users *UserController, auth *AuthController, homestays *HomestayController) *BaseController {
	return &BaseController{
		users:     users,
		auth:      auth,
		homestays: homestays,
	}
}

func (bc *BaseController) Login(c *gin.Context) {
	bc.auth.Login(c)
}

func (bc *BaseController) Logout(c *gin.Context) {
	bc.auth.Logout(c)
}

func (bc *BaseController) Me(c *gin.Context) {
	bc.auth.Me(c)
}

func (bc *BaseController) Refresh(c *gin.Context) {
	bc.auth.Refresh(c)
}

func (bc *BaseController) FetchHomestays(c *gin.Context) {
	bc.homestays.FetchHomestays(c)
}

func (bc *BaseController) GetUsers(c *gin.Context) {
	bc.users.GetUsers(c)
}

func (bc *BaseController) CreateUser(c *gin.Context) {
	bc.users.CreateUser(c)
}

var _ api.ServerInterface = (*BaseController)(nil)

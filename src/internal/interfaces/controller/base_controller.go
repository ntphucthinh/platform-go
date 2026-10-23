package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "platform-go/generated/api"
)

var _ api.ServerInterface = (*ServerController)(nil)

// ServerController is the HTTP composition root for generated API endpoints.
// It delegates each endpoint to the controller that owns its business workflow.
type ServerController struct {
	userController     *UserController
	authController     *AuthController
	homestayController HomestayController
	uploadController   *UploadController
}

func NewServerController(userController *UserController, authController *AuthController, homestayController HomestayController, uploadController *UploadController) *ServerController {
	return &ServerController{
		userController:     userController,
		authController:     authController,
		homestayController: homestayController,
		uploadController:   uploadController,
	}
}

func (con *ServerController) UploadFiles(c *gin.Context) { con.uploadController.UploadFiles(c) }

func (con *ServerController) DeleteUploads(c *gin.Context) { con.uploadController.DeleteUploads(c) }

func (con *ServerController) Login(c *gin.Context) {
	con.authController.Login(c)
}

func (con *ServerController) Logout(c *gin.Context) {
	con.authController.Logout(c)
}

func (con *ServerController) Me(c *gin.Context) {
	con.authController.Me(c)
}

func (con *ServerController) Refresh(c *gin.Context) {
	con.authController.Refresh(c)
}

func (con *ServerController) FetchHomestays(c *gin.Context) {
	con.homestayController.FetchHomestays(c)
}

// TODO: delegate to HomestayController once the create use case is implemented.
func (con *ServerController) CreateHomestay(c *gin.Context) {
	con.homestayController.CreateHomestay(c)
}

// TODO: delegate to HomestayController once the detail use case is implemented.
func (con *ServerController) FetchHomestayByID(c *gin.Context, homestayID int64) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "not implemented",
	})
}

// TODO: delegate to HomestayController once the update use case is implemented.
func (con *ServerController) UpdateHomestay(c *gin.Context, homestayID int64) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "not implemented",
	})
}

func (con *ServerController) GetUsers(c *gin.Context) {
	con.userController.GetUsers(c)
}

func (con *ServerController) CreateUser(c *gin.Context) {
	con.userController.CreateUser(c)
}

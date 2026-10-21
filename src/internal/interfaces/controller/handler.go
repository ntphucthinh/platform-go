package controller

import api "platform-go/generated/api"

// Handler combines generated endpoints.
type Handler struct {
	*UserController
	*AuthController
	*HomestayController
}

func NewHandler(users *UserController, auth *AuthController, homestays *HomestayController) *Handler {
	return &Handler{
		UserController:     users,
		AuthController:     auth,
		HomestayController: homestays,
	}
}

var _ api.ServerInterface = (*Handler)(nil)

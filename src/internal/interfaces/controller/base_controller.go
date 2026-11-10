package controller

import (
	"github.com/gin-gonic/gin"

	api "platform-go/generated/api"
)

var _ api.ServerInterface = (*ServerController)(nil)

// ServerController is the HTTP composition root for generated API endpoints.
// It delegates each endpoint to the controller that owns its business workflow.
type ServerController struct {
	userController           *UserController
	authController           *AuthController
	homestayController       HomestayController
	travelLocationController TravelLocationController
	travelMemoryController   TravelMemoryController
	uploadController         *UploadController
	selectionController      SelectionController
}

func NewServerController(
	userController *UserController,
	authController *AuthController,
	homestayController HomestayController,
	travelLocationController TravelLocationController,
	travelMemoryController TravelMemoryController,
	uploadController *UploadController,
	selectionController SelectionController,
) *ServerController {
	return &ServerController{
		userController:           userController,
		authController:           authController,
		homestayController:       homestayController,
		travelLocationController: travelLocationController,
		travelMemoryController:   travelMemoryController,
		uploadController:         uploadController,
		selectionController:      selectionController,
	}
}

// -----------------------------------------------------------------------------
// Uploads APIs (UploadController)
// -----------------------------------------------------------------------------

// UploadFiles handles uploading files to cloud storage.
// (POST /uploads)
func (con *ServerController) UploadFiles(c *gin.Context) {
	con.uploadController.UploadFiles(c)
}

// DeleteUploads handles deleting unused uploaded files from cloud storage.
// (DELETE /uploads)
func (con *ServerController) DeleteUploads(c *gin.Context) {
	con.uploadController.DeleteUploads(c)
}

// -----------------------------------------------------------------------------
// Authentication APIs (AuthController)
// -----------------------------------------------------------------------------

// Login handles user authentication and issues JWT tokens.
// (POST /auth/login)
func (con *ServerController) Login(c *gin.Context) {
	con.authController.Login(c)
}

// Logout handles clearing authentication cookies/session.
// (POST /auth/logout)
func (con *ServerController) Logout(c *gin.Context) {
	con.authController.Logout(c)
}

// Me retrieves the currently authenticated user's profile.
// (GET /auth/me)
func (con *ServerController) Me(c *gin.Context) {
	con.authController.Me(c)
}

// Refresh handles refreshing the JWT access token using the refresh token.
// (POST /auth/refresh)
func (con *ServerController) Refresh(c *gin.Context) {
	con.authController.Refresh(c)
}

// -----------------------------------------------------------------------------
// Homestays APIs (HomestayController)
// -----------------------------------------------------------------------------

// FetchHomestays retrieves a list of homestays with their rooms and images.
// (GET /homestays)
func (con *ServerController) FetchHomestays(c *gin.Context, params api.FetchHomestaysParams) {
	con.homestayController.FetchHomestays(c, params)
}

// CreateHomestay creates a new homestay along with its rooms and images.
// (POST /homestays)
func (con *ServerController) CreateHomestay(c *gin.Context) {
	con.homestayController.CreateHomestay(c)
}

// FetchHomestayByID retrieves details of a single homestay by ID.
// (GET /homestays/{homestayID})
func (con *ServerController) FetchHomestayByID(c *gin.Context, homestayID int64) {
	con.homestayController.FetchHomestayByID(c, homestayID)
}

// UpdateHomestay updates a homestay and synchronizes its rooms and images.
// (PUT /homestays/{homestayID})
func (con *ServerController) UpdateHomestay(c *gin.Context, homestayID int64) {
	con.homestayController.UpdateHomestay(c, homestayID)
}

// DeleteHomestay removes a homestay and all its associated rooms and images.
// (DELETE /homestays/{homestayID})
func (con *ServerController) DeleteHomestay(c *gin.Context, homestayID int64) {
	con.homestayController.DeleteHomestay(c, homestayID)
}

// -----------------------------------------------------------------------------
// Travel Locations APIs (TravelLocationController)
// -----------------------------------------------------------------------------

// FetchTravelLocations retrieves a list of travel locations.
// (GET /travel-locations)
func (con *ServerController) FetchTravelLocations(c *gin.Context, params api.FetchTravelLocationsParams) {
	con.travelLocationController.FetchTravelLocations(c, params)
}

// CreateTravelLocation creates a new travel location.
// (POST /travel-locations)
func (con *ServerController) CreateTravelLocation(c *gin.Context) {
	con.travelLocationController.CreateTravelLocation(c)
}

// FetchTravelLocationByID retrieves details of a single travel location by ID.
// (GET /travel-locations/{travelLocationID})
func (con *ServerController) FetchTravelLocationByID(c *gin.Context, travelLocationID int64) {
	con.travelLocationController.FetchTravelLocationByID(c, travelLocationID)
}

// UpdateTravelLocation updates an existing travel location by ID.
// (PUT /travel-locations/{travelLocationID})
func (con *ServerController) UpdateTravelLocation(c *gin.Context, travelLocationID int64) {
	con.travelLocationController.UpdateTravelLocation(c, travelLocationID)
}

// DeleteTravelLocation soft-deletes a travel location by ID.
// (DELETE /travel-locations/{travelLocationID})
func (con *ServerController) DeleteTravelLocation(c *gin.Context, travelLocationID int64) {
	con.travelLocationController.DeleteTravelLocation(c, travelLocationID)
}

// -----------------------------------------------------------------------------
// Travel Memories APIs (TravelMemoryController)
// -----------------------------------------------------------------------------

// FetchTravelMemories retrieves a paginated list of travel memories.
// (GET /travel-memories)
func (con *ServerController) FetchTravelMemories(c *gin.Context, params api.FetchTravelMemoriesParams) {
	con.travelMemoryController.FetchTravelMemories(c, params)
}

// FetchTravelMemoryByID retrieves details of a single travel memory by ID.
// (GET /travel-memories/{travelMemoryID})
func (con *ServerController) FetchTravelMemoryByID(c *gin.Context, travelMemoryID int64) {
	con.travelMemoryController.FetchTravelMemoryByID(c, travelMemoryID)
}

// -----------------------------------------------------------------------------
// Selections APIs (SelectionController)
// -----------------------------------------------------------------------------

// FetchSelectionTravelLocations retrieves a lightweight list of travel locations for dropdown selection.
// (GET /selections/travel-locations)
func (con *ServerController) FetchSelectionTravelLocations(c *gin.Context, params api.FetchSelectionTravelLocationsParams) {
	con.selectionController.FetchSelectionTravelLocations(c, params)
}

// -----------------------------------------------------------------------------
// Users APIs (UserController)
// -----------------------------------------------------------------------------

// GetUsers retrieves a list of users.
// (GET /users)
func (con *ServerController) GetUsers(c *gin.Context) {
	con.userController.GetUsers(c)
}

// CreateUser creates a new user account.
// (POST /users)
func (con *ServerController) CreateUser(c *gin.Context) {
	con.userController.CreateUser(c)
}

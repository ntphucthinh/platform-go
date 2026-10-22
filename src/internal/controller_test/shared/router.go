package shared

import (
	"testing"

	"github.com/gin-gonic/gin"
)

// NewSetupRouter creates an isolated Gin router for controller HTTP tests.
// Upload tests use a mocked external provider and do not require database fixtures.
func NewSetupRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())
	return router
}

package middleware

import (
	"platform-go/internal/common/ctxkeys"

	"github.com/gin-gonic/gin"
)

func RequestID(c *gin.Context) string {
	v, ok := c.Get(ctxkeys.RequestIDKey)
	if !ok {
		return ""
	}
	requestID, _ := v.(string)
	return requestID
}

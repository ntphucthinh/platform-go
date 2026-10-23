package common

import (
	"context"
	"net/http"
	"platform-go/internal/common/ctxkeys"
	"platform-go/internal/middleware"
	"platform-go/internal/pkg/cvalidator"
	"platform-go/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

func requestID(c *gin.Context) string {
	return middleware.RequestID(c)
}

// BindJSON binds JSON request body and returns 400 Bad Request if validation fails
func BindJSON(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		errMessages := cvalidator.TranslateError(err)
		log.Warn(err.Error(), "request_id", requestID(c), "errors", errMessages)

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": errMessages,
		})
		return err
	}

	v := cvalidator.Validator()
	ctx := context.WithValue(c.Request.Context(), ctxkeys.GinContextKey{}, c)
	if err := v.StructCtx(ctx, req); err != nil {
		errMessages := cvalidator.TranslateError(err)
		log.Warn(err.Error(), "request_id", requestID(c), "errors", errMessages)

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": errMessages,
		})
		return err
	}

	return nil
}

// IsErrorCustom returns an error response with custom HTTP status code and message
func IsErrorCustom(c *gin.Context, err error, status int, message string) {
	c.JSON(status, gin.H{
		"error":   message,
		"details": err.Error(),
	})
}

// IsErrorWithMessage returns a 500 Internal Server Error with custom message
func IsErrorWithMessage(c *gin.Context, err error, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   message,
		"details": err.Error(),
	})
}

// IsValidationError checks whether an error is a validation error
func IsValidationError(err error) bool {
	// In practice, check the error type (e.g. validator.ValidationErrors)
	return err != nil
}

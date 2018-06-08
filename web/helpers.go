package web

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
)

// StatusCodeForError returns an http status code for an error type.
func StatusCodeForError(err interface{}) int {
	switch err.(type) {
	case *services.ValidationError:
		return 400
	default:
		return 500
	}
}

// publicError adds an error to the gin context and sets
// the JSON value of errors.
func publicError(c *gin.Context, statusCode int, err error) {
	c.Error(err).SetType(gin.ErrorTypePublic)
	c.JSON(statusCode, gin.H{"errors": []string{err.Error()}})
}

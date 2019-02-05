package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
)

// StatusCodeForError returns an http status code for an error type.
func StatusCodeForError(err interface{}) int {
	switch err.(type) {
	case *models.ValidationError:
		return 400
	default:
		return 500
	}
}

// publicError adds an error to the gin context and sets
// the JSON value of errors.
func publicError(c *gin.Context, statusCode int, err error) {
	c.Error(err).SetType(gin.ErrorTypePublic)
	switch v := err.(type) {
	case *models.JSONAPIErrors:
		c.JSON(statusCode, v)
	default:
		c.JSON(statusCode, models.NewJSONAPIErrorsWith(err.Error()))
	}
}

func paginatedResponse(
	c *gin.Context,
	name string,
	size int,
	page int,
	resource interface{},
	count int,
	err error,
) {
	if err == orm.ErrorNotFound {
		err = nil
	}

	if err != nil {
		c.AbortWithError(500, fmt.Errorf("error getting paged %s: %+v", name, err))
	} else if buffer, err := NewPaginatedResponse(*c.Request.URL, size, page, count, resource); err != nil {
		c.AbortWithError(500, fmt.Errorf("failed to marshal document: %+v", err))
	} else {
		c.Data(200, MediaType, buffer)
	}
}

func paginatedRequest(action func(*gin.Context, int, int, int)) func(*gin.Context) {
	return func(c *gin.Context) {
		size, page, offset, err := ParsePaginatedRequest(c.Query("size"), c.Query("page"))
		if err != nil {
			publicError(c, 422, err)
			return
		}
		action(c, size, page, offset)
	}
}

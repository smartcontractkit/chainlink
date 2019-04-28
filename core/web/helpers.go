package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

// StatusCodeForError returns an http status code for an error type.
func StatusCodeForError(err interface{}) int {
	switch err.(type) {
	case *models.ValidationError:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
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
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting paged %s: %+v", name, err))
	} else if buffer, err := NewPaginatedResponse(*c.Request.URL, size, page, count, resource); err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to marshal document: %+v", err))
	} else {
		c.Data(http.StatusOK, MediaType, buffer)
	}
}

func paginatedRequest(action func(*gin.Context, int, int, int)) func(*gin.Context) {
	return func(c *gin.Context) {
		size, page, offset, err := ParsePaginatedRequest(c.Query("size"), c.Query("page"))
		if err != nil {
			publicError(c, http.StatusUnprocessableEntity, err)
			return
		}
		action(c, size, page, offset)
	}
}

func jsonAPIResponseWithStatus(c *gin.Context, resource interface{}, name string, status int) {
	json, err := jsonapi.Marshal(resource)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to marshal %s using jsonapi: %+v", name, err))
	} else {
		c.Data(status, MediaType, json)
	}
}

func jsonAPIResponse(c *gin.Context, resource interface{}, name string) {
	jsonAPIResponseWithStatus(c, resource, name, http.StatusOK)
}

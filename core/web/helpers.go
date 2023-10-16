package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// jsonAPIError adds an error to the gin context and sets
// the JSON value of errors.
func jsonAPIError(c *gin.Context, statusCode int, err error) {
	_ = c.Error(err).SetType(gin.ErrorTypePublic)
	var jsonErr *models.JSONAPIErrors
	if errors.As(err, &jsonErr) {
		c.JSON(statusCode, jsonErr)
		return
	}
	c.JSON(statusCode, models.NewJSONAPIErrorsWith(err.Error()))
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
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("error getting paged %s: %+v", name, err))
	} else if buffer, err := NewPaginatedResponse(*c.Request.URL, size, page, count, resource); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("failed to marshal document: %+v", err))
	} else {
		c.Data(http.StatusOK, MediaType, buffer)
	}
}

func paginatedRequest(action func(*gin.Context, int, int, int)) func(*gin.Context) {
	return func(c *gin.Context) {
		size, page, offset, err := ParsePaginatedRequest(c.Query("size"), c.Query("page"))
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, err)
			return
		}
		action(c, size, page, offset)
	}
}

func jsonAPIResponseWithStatus(c *gin.Context, resource interface{}, name string, status int) {
	json, err := jsonapi.Marshal(resource)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("failed to marshal %s using jsonapi: %+v", name, err))
	} else {
		c.Data(status, MediaType, json)
	}
}

func jsonAPIResponse(c *gin.Context, resource interface{}, name string) {
	jsonAPIResponseWithStatus(c, resource, name, http.StatusOK)
}

func Router(t testing.TB, app chainlink.Application, prometheus *ginprom.Prometheus) *gin.Engine {
	r, err := NewRouter(app, prometheus)
	require.NoError(t, err)
	return r
}

package web

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// FeedsManagerController manages the feeds managers
type FeedsManagerController struct {
	App chainlink.Application
}

// CreateFeedsManagerRequest represents a JSONAPI request for registering a
// feeds manager
type CreateFeedsManagerRequest struct {
	Name      string          `json:"name"`
	URI       string          `json:"uri"`
	JobTypes  []string        `json:"jobTypes"`
	PublicKey feeds.PublicKey `json:"publicKey"`
	Network   string          `json:"network"`
}

// Create registers a new feeds manager.
// Example:
// "POST <application>/feeds_managers"
func (fmc *FeedsManagerController) Create(c *gin.Context) {
	request := CreateFeedsManagerRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	ms := &feeds.ManagerService{
		URI:       request.URI,
		Name:      request.Name,
		PublicKey: request.PublicKey,
		JobTypes:  request.JobTypes,
		Network:   request.Network,
	}

	feedsService := fmc.App.GetFeedsService()

	id, err := feedsService.RegisterManagerService(ms)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	ms, err = feedsService.GetManagerService(id)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c,
		presenters.NewFeedsManagerResource(*ms),
		"feeds_managers",
		http.StatusCreated,
	)
}

// List retrieves all the feeds managers
// Example:
// "GET <application>/feeds_managers"
func (fmc *FeedsManagerController) List(c *gin.Context) {
	mss, err := fmc.App.GetFeedsService().ListManagerServices()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewFeedsManagerResources(mss), "feeds_managers")
}

// Show retrieve a feeds manager by id
// Example:
// "GET <application>/feeds_managers/<id>"
func (fmc *FeedsManagerController) Show(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	ms, err := fmc.App.GetFeedsService().GetManagerService(int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			jsonAPIError(c, http.StatusNotFound, errors.New("feeds Manager not found"))
			return
		}

		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewFeedsManagerResource(*ms), "feeds_managers")
}

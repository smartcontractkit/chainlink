package web

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// ExternalInitiatorsController manages external initiators
type ExternalInitiatorsController struct {
	App services.Application
}

// Create builds and saves a new service agreement record.
func (eic *ExternalInitiatorsController) Create(c *gin.Context) {
	if !eic.App.GetStore().Config.Dev() {
		publicError(c, http.StatusMethodNotAllowed, errors.New("External Initiators are currently under development and not yet usable outside of development mode"))
		return
	}

	eia := models.NewExternalInitiatorAuthentication()
	if ea, err := models.NewExternalInitiator(eia); err != nil {
		publicError(c, http.StatusInternalServerError, err)
	} else if err := eic.App.GetStore().CreateExternalInitiator(ea); err != nil {
		publicError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponseWithStatus(c, eia, "external initiator authenticaion", http.StatusCreated)
	}
}

// Destroy deletes an ExternalInitiator
func (eic *ExternalInitiatorsController) Destroy(c *gin.Context) {
	if !eic.App.GetStore().Config.Dev() {
		publicError(c, http.StatusMethodNotAllowed, errors.New("External Initiators are currently under development and not yet usable outside of development mode"))
		return
	}

	id := c.Param("AccessKey")
	if err := eic.App.GetStore().DeleteExternalInitiator(id); err != nil {
		publicError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponseWithStatus(c, nil, "external initiator", http.StatusNoContent)
	}
}

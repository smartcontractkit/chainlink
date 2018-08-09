package web

import (
	"errors"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
)

// ServiceAgreementsController manages service agreements.
type ServiceAgreementsController struct {
	App *services.ChainlinkApplication
}

// Create builds and saves a new service agreement record.
func (sac *ServiceAgreementsController) Create(c *gin.Context) {
	var sar models.ServiceAgreementRequest
	if !sac.App.Store.Config.Dev {
		publicError(c, 500, errors.New("Service Agreements are currently under development and not yet usable outside of development mode"))
	} else if err := c.ShouldBindJSON(&sar); err != nil {
		publicError(c, 400, err)
	} else if sa, err := models.NewServiceAgreementFromRequest(sar); err != nil {
		publicError(c, 400, err)
	} else if err = sac.App.Store.SaveServiceAgreement(&sa); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.JSON(200, sa)
	}
}

// Show returns the details of a ServiceAgreement.
// Example:
//  "<application>/service_agreements/:SAID"
func (sac *ServiceAgreementsController) Show(c *gin.Context) {
	id := c.Param("SAID")
	if sa, err := sac.App.Store.FindServiceAgreement(id); err == storm.ErrNotFound {
		publicError(c, 404, errors.New("ServiceAgreement not found"))
	} else if err != nil {
		c.AbortWithError(500, err)
	} else if doc, err := jsonapi.MarshalToStruct(presenters.ServiceAgreement{ServiceAgreement: sa}, nil); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.JSON(200, doc)
	}
}

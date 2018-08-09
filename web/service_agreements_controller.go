package web

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
)

// ServiceAgreementsController manages service agreements.
type ServiceAgreementsController struct {
	App *services.ChainlinkApplication
}

// Create builds and saves a new service agreement record.
func (sac *ServiceAgreementsController) Create(c *gin.Context) {
	var sar models.ServiceAgreementRequest
	if err := c.ShouldBindJSON(&sar); err != nil {
		publicError(c, 400, err)
	} else if sa, err := models.NewServiceAgreementFromRequest(sar); err != nil {
		publicError(c, 400, err)
	} else if err = sac.App.Store.SaveServiceAgreement(&sa); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.JSON(200, sa)
	}
}

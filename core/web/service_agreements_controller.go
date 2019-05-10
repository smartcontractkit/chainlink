package web

import (
	"errors"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
)

// ServiceAgreementsController manages service agreements.
type ServiceAgreementsController struct {
	App services.Application
}

// Create builds and saves a new service agreement record.
func (sac *ServiceAgreementsController) Create(c *gin.Context) {
	if !sac.App.GetStore().Config.Dev() {
		publicError(c, http.StatusMethodNotAllowed, errors.New("Service Agreements are currently under development and not yet usable outside of development mode"))
		return
	}

	us, err := models.NewUnsignedServiceAgreementFromRequest(c.Request.Body)
	if err != nil {
		publicError(c, http.StatusUnprocessableEntity, err)
		return
	}

	sa, err := sac.App.GetStore().FindServiceAgreement(us.ID.String())
	if err == orm.ErrorNotFound {
		sa, err = models.BuildServiceAgreement(us, sac.App.GetStore().KeyStore)
		if err != nil {
			publicError(c, http.StatusUnprocessableEntity, err)
			return
		} else if err = services.ValidateServiceAgreement(sa, sac.App.GetStore()); err != nil {
			publicError(c, http.StatusUnprocessableEntity, err)
			return
		} else if err = sac.App.AddServiceAgreement(&sa); err != nil {
			publicError(c, http.StatusInternalServerError, err)
			return
		}
	}
	jsonAPIResponse(c, sa, "service agreement")
}

// Show returns the details of a ServiceAgreement.
// Example:
//  "<application>/service_agreements/:SAID"
func (sac *ServiceAgreementsController) Show(c *gin.Context) {
	id := common.HexToHash(c.Param("SAID"))
	if sa, err := sac.App.GetStore().FindServiceAgreement(id.String()); err == orm.ErrorNotFound {
		publicError(c, http.StatusNotFound, errors.New("ServiceAgreement not found"))
	} else if err != nil {
		publicError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponse(c, presenters.ServiceAgreement{ServiceAgreement: sa}, "service agreement")
	}
}

package web

import (
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// ServiceAgreementsController manages service agreements.
type ServiceAgreementsController struct {
	App chainlink.Application
}

// Create builds and saves a new service agreement record.
func (sac *ServiceAgreementsController) Create(c *gin.Context) {
	if !sac.App.GetStore().Config.Dev() {
		jsonAPIError(c, http.StatusMethodNotAllowed, errors.New("Service Agreements are currently under development and not yet usable outside of development mode"))
		return
	}

	us, err := models.NewUnsignedServiceAgreementFromRequest(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	sa, err := sac.App.GetStore().FindServiceAgreement(us.ID.String())
	if errors.Cause(err) == orm.ErrorNotFound {
		sa, err = models.BuildServiceAgreement(us, models.NullSigner{})
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, err)
			return
		}
		if err = services.ValidateServiceAgreement(sa, sac.App.GetStore(), sac.App.GetKeyStore()); err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, err)
			return
		}
		if err = sac.App.AddServiceAgreement(&sa); err != nil {
			jsonAPIError(c, http.StatusInternalServerError, errors.Wrap(err, "#AddServiceAgreement"))
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

	sa, err := sac.App.GetStore().FindServiceAgreement(id.String())
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("ServiceAgreement not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.ServiceAgreement{ServiceAgreement: sa}, "service agreement")
}

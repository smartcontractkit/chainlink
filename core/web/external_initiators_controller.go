package web

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	externalInitiatorNameRegexp = regexp.MustCompile("^[a-zA-Z0-9-_]+$")
)

// ValidateExternalInitiator checks whether External Initiator parameters are
// safe for processing.
func ValidateExternalInitiator(
	exi *models.ExternalInitiatorRequest,
	store *store.Store,
) error {
	fe := models.NewJSONAPIErrors()
	if len([]rune(exi.Name)) == 0 {
		fe.Add("No name specified")
	} else if !externalInitiatorNameRegexp.MatchString(exi.Name) {
		fe.Add("Name must be alphanumeric and may contain '_' or '-'")
	} else if _, err := store.FindExternalInitiatorByName(exi.Name); err == nil {
		fe.Add(fmt.Sprintf("Name %v already exists", exi.Name))
	} else if err != orm.ErrorNotFound {
		return errors.Wrap(err, "validating external initiator")
	}
	return fe.CoerceEmptyToNil()
}

// ExternalInitiatorsController manages external initiators
type ExternalInitiatorsController struct {
	App chainlink.Application
}

func (eic *ExternalInitiatorsController) Index(c *gin.Context, size, page, offset int) {
	eis, count, err := eic.App.GetStore().ExternalInitiatorsSorted(offset, size)
	var resources []presenters.ExternalInitiatorResource
	for _, ei := range eis {
		resources = append(resources, presenters.NewExternalInitiatorResource(ei))
	}

	paginatedResponse(c, "externalInitiators", size, page, resources, count, err)
}

// Create builds and saves a new external initiator
func (eic *ExternalInitiatorsController) Create(c *gin.Context) {
	eir := &models.ExternalInitiatorRequest{}
	if !eic.App.GetConfig().Dev() && !eic.App.GetConfig().FeatureExternalInitiators() {
		err := errors.New("The External Initiator feature is disabled by configuration")
		jsonAPIError(c, http.StatusMethodNotAllowed, err)
		return
	}

	eia := auth.NewToken()
	if err := c.ShouldBindJSON(eir); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	ei, err := models.NewExternalInitiator(eia, eir)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	if err := ValidateExternalInitiator(eir, eic.App.GetStore()); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	if err := eic.App.GetStore().CreateExternalInitiator(ei); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	resp := presenters.NewExternalInitiatorAuthentication(*ei, *eia)
	jsonAPIResponseWithStatus(c, resp, "external initiator authentication", http.StatusCreated)
}

// Destroy deletes an ExternalInitiator
func (eic *ExternalInitiatorsController) Destroy(c *gin.Context) {
	name := c.Param("Name")
	exi, err := eic.App.GetStore().FindExternalInitiatorByName(name)
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("external initiator not found"))
		return
	}
	if err := eic.App.GetStore().DeleteExternalInitiator(exi.Name); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "external initiator", http.StatusNoContent)
}

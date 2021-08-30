package web

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgconn"

	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// BridgeTypesController manages BridgeType requests in the node.
type BridgeTypesController struct {
	App chainlink.Application
}

// Create adds the BridgeType to the given context.
func (btc *BridgeTypesController) Create(c *gin.Context) {
	btr := &models.BridgeTypeRequest{}

	if err := c.ShouldBindJSON(btr); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	bta, bt, err := models.NewBridgeType(btr)
	if err != nil {
		jsonAPIError(c, StatusCodeForError(err), err)
		return
	}
	if e := services.ValidateBridgeType(btr, btc.App.GetStore()); e != nil {
		jsonAPIError(c, http.StatusBadRequest, e)
		return
	}
	if e := services.ValidateBridgeTypeNotExist(btr, btc.App.GetStore()); e != nil {
		jsonAPIError(c, http.StatusBadRequest, e)
		return
	}
	if e := btc.App.GetStore().CreateBridgeType(bt); e != nil {
		jsonAPIError(c, http.StatusInternalServerError, e)
		return
	}
	switch e := err.(type) {
	case *pgconn.PgError:
		var apiErr error
		if e.ConstraintName == "external_initiators_name_key" {
			apiErr = fmt.Errorf("bridge Type %v conflict", bt.Name)
		} else {
			apiErr = err
		}
		jsonAPIError(c, http.StatusConflict, apiErr)
		return
	default:
		resource := presenters.NewBridgeResource(*bt)
		resource.IncomingToken = bta.IncomingToken

		jsonAPIResponse(c, resource, "bridge")
	}
}

// Index lists Bridges, one page at a time.
func (btc *BridgeTypesController) Index(c *gin.Context, size, page, offset int) {
	bridges, count, err := btc.App.GetStore().BridgeTypes(offset, size)

	var resources []presenters.BridgeResource
	for _, bridge := range bridges {
		resources = append(resources, *presenters.NewBridgeResource(bridge))
	}

	paginatedResponse(c, "Bridges", size, page, resources, count, err)
}

// Show returns the details of a specific Bridge.
func (btc *BridgeTypesController) Show(c *gin.Context) {
	name := c.Param("BridgeName")

	taskType, err := models.NewTaskType(name)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	bt, err := btc.App.GetStore().FindBridge(taskType)
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("bridge not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewBridgeResource(bt), "bridge")
}

// Update can change the restricted attributes for a bridge
func (btc *BridgeTypesController) Update(c *gin.Context) {
	name := c.Param("BridgeName")
	btr := &models.BridgeTypeRequest{}

	taskType, err := models.NewTaskType(name)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	bt, err := btc.App.GetStore().FindBridge(taskType)
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("bridge not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	if err := c.ShouldBindJSON(btr); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	if err := services.ValidateBridgeType(btr, btc.App.GetStore()); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	if err := btc.App.GetStore().UpdateBridgeType(&bt, btr); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewBridgeResource(bt), "bridge")
}

// Destroy removes a specific Bridge.
func (btc *BridgeTypesController) Destroy(c *gin.Context) {
	name := c.Param("BridgeName")

	taskType, err := models.NewTaskType(name)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	bt, err := btc.App.GetStore().FindBridge(taskType)
	if errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("bridge not found"))
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("error searching for bridge: %+v", err))
		return
	}
	jobsUsingBridge, err := btc.App.JobORM().FindJobIDsWithBridge(name)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("error searching for associated v2 jobs: %+v", err))
		return
	}
	if len(jobsUsingBridge) > 0 {
		jsonAPIError(c, http.StatusConflict, fmt.Errorf("can't remove the bridge because jobs %v are associated with it", jobsUsingBridge))
		return
	}
	if err = btc.App.GetStore().DeleteBridgeType(&bt); err != nil {
		jsonAPIError(c, StatusCodeForError(err), fmt.Errorf("failed to delete bridge: %+v", err))
		return
	}

	jsonAPIResponse(c, presenters.NewBridgeResource(bt), "bridge")
}

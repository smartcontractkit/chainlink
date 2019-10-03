package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"chainlink/core/services"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"
)

// BridgeTypesController manages BridgeType requests in the node.
type BridgeTypesController struct {
	App services.Application
}

// Create adds the BridgeType to the given context.
func (btc *BridgeTypesController) Create(c *gin.Context) {
	btr := &models.BridgeTypeRequest{}

	if err := c.ShouldBindJSON(btr); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	} else if bta, bt, err := models.NewBridgeType(btr); err != nil {
		jsonAPIError(c, StatusCodeForError(err), err)
	} else if err := services.ValidateBridgeType(btr, btc.App.GetStore()); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
	} else if err := btc.App.GetStore().CreateBridgeType(bt); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponse(c, bta, "bridge")
	}
}

// Index lists Bridges, one page at a time.
func (btc *BridgeTypesController) Index(c *gin.Context, size, page, offset int) {
	bridges, count, err := btc.App.GetStore().BridgeTypes(offset, size)
	paginatedResponse(c, "Bridges", size, page, bridges, count, err)
}

// Show returns the details of a specific Bridge.
func (btc *BridgeTypesController) Show(c *gin.Context) {
	name := c.Param("BridgeName")
	if taskType, err := models.NewTaskType(name); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	} else if bt, err := btc.App.GetStore().FindBridge(taskType); errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("bridge not found"))
	} else if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponse(c, bt, "bridge")
	}
}

// Update can change the restricted attributes for a bridge
func (btc *BridgeTypesController) Update(c *gin.Context) {
	name := c.Param("BridgeName")
	btr := &models.BridgeTypeRequest{}

	if taskType, err := models.NewTaskType(name); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	} else if bt, err := btc.App.GetStore().FindBridge(taskType); errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("bridge not found"))
	} else if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else if err := c.ShouldBindJSON(btr); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	} else if err := btc.App.GetStore().UpdateBridgeType(&bt, btr); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponse(c, bt, "bridge")
	}
}

// Destroy removes a specific Bridge.
func (btc *BridgeTypesController) Destroy(c *gin.Context) {
	name := c.Param("BridgeName")

	if taskType, err := models.NewTaskType(name); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	} else if bt, err := btc.App.GetStore().FindBridge(taskType); errors.Cause(err) == orm.ErrorNotFound {
		jsonAPIError(c, http.StatusNotFound, errors.New("bridge not found"))
	} else if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("Error searching for bridge for BTC Destroy: %+v", err))
	} else if jobFounds, err := btc.App.GetStore().AnyJobWithType(name); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("Error searching for associated jobs for BTC Destroy: %+v", err))
	} else if jobFounds {
		jsonAPIError(c, http.StatusConflict, fmt.Errorf("Can't remove the bridge because there are jobs associated with it: %+v", err))
	} else if err = btc.App.GetStore().DeleteBridgeType(&bt); err != nil {
		jsonAPIError(c, StatusCodeForError(err), fmt.Errorf("failed to initialise BTC Destroy: %+v", err))
	} else {
		jsonAPIResponse(c, bt, "bridge")
	}
}

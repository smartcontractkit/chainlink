package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/forms"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

// BridgeTypesController manages BridgeType requests in the node.
type BridgeTypesController struct {
	App services.Application
}

// Create adds the BridgeType to the given context.
func (btc *BridgeTypesController) Create(c *gin.Context) {
	btr := &models.BridgeTypeRequest{}

	if err := c.ShouldBindJSON(btr); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else if bta, bt, err := models.NewBridgeType(btr); err != nil {
		publicError(c, StatusCodeForError(err), err)
	} else if err := btc.App.GetStore().CreateBridgeType(bt); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else if doc, err := jsonapi.Marshal(bta); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.Data(http.StatusOK, MediaType, doc)
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
	if bt, err := btc.App.GetStore().FindBridge(name); err == orm.ErrorNotFound {
		publicError(c, http.StatusNotFound, errors.New("bridge name not found"))
	} else if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else if doc, err := jsonapi.Marshal(bt); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.Data(http.StatusOK, MediaType, doc)
	}
}

// Update can change the restricted attributes for a bridge
func (btc *BridgeTypesController) Update(c *gin.Context) {
	bn := c.Param("BridgeName")
	form, err := forms.NewUpdateBridgeType(btc.App.GetStore(), bn)

	if err == orm.ErrorNotFound {
		publicError(c, http.StatusNotFound, errors.New("bridge name not found"))
		return
	}

	if err = c.BindJSON(&form); err != nil {
		publicError(c, 400, fmt.Errorf("unable to parse JSON: %v", err))
		return
	}

	err = form.Save()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	doc, err := form.Marshal()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, MediaType, doc)
}

// Destroy removes a specific Bridge.
func (btc *BridgeTypesController) Destroy(c *gin.Context) {
	name := c.Param("BridgeName")
	if bt, err := btc.App.GetStore().FindBridge(name); err == orm.ErrorNotFound {
		publicError(c, http.StatusNotFound, errors.New("bridge name not found"))
	} else if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Error searching for bridge for BTC Destroy: %+v", err))
	} else if jobFounds, err := btc.App.GetStore().AnyJobWithType(name); err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Error searching for associated jobs for BTC Destroy: %+v", err))
	} else if jobFounds {
		c.AbortWithError(http.StatusConflict, fmt.Errorf("Can't remove the bridge because there are jobs associated with it: %+v", err))
	} else if err = btc.App.GetStore().DeleteBridgeType(&bt); err != nil {
		c.AbortWithError(StatusCodeForError(err), fmt.Errorf("failed to initialise BTC Destroy: %+v", err))
	} else {
		c.JSON(http.StatusOK, bt)
	}
}

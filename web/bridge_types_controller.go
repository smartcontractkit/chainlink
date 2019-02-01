package web

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/forms"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
	"github.com/smartcontractkit/chainlink/store/presenters"
)

// BridgeTypesController manages BridgeType requests in the node.
type BridgeTypesController struct {
	App services.Application
}

// Create adds the BridgeType to the given context.
func (btc *BridgeTypesController) Create(c *gin.Context) {
	bt := &models.BridgeType{}

	if err := c.ShouldBindJSON(bt); err != nil {
		c.AbortWithError(500, err)
	} else if err = btc.App.AddAdapter(bt); err != nil {
		publicError(c, StatusCodeForError(err), err)
	} else if doc, err := jsonapi.Marshal(presenters.BridgeType{BridgeType: *bt}); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.Data(200, MediaType, doc)
	}
}

// Index lists Bridges, one page at a time.
func (btc *BridgeTypesController) Index(c *gin.Context, size, page, offset int) {
	bridges, count, err := btc.App.GetStore().BridgeTypes(offset, size)
	pbt := make([]presenters.BridgeType, len(bridges))
	for i, j := range bridges {
		pbt[i] = presenters.BridgeType{BridgeType: j}
	}

	paginatedResponse(c, "Bridges", size, page, pbt, count, err)
}

// Show returns the details of a specific Bridge.
func (btc *BridgeTypesController) Show(c *gin.Context) {
	name := c.Param("BridgeName")
	if bt, err := btc.App.GetStore().FindBridge(name); err == orm.ErrorNotFound {
		publicError(c, 404, errors.New("bridge name not found"))
	} else if err != nil {
		c.AbortWithError(500, err)
	} else if doc, err := jsonapi.Marshal(presenters.BridgeType{BridgeType: bt}); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.Data(200, MediaType, doc)
	}
}

// Update can change the restricted attributes for a bridge
func (btc *BridgeTypesController) Update(c *gin.Context) {
	bn := c.Param("BridgeName")
	form, err := forms.NewUpdateBridgeType(btc.App.GetStore(), bn)

	if err == orm.ErrorNotFound {
		publicError(c, 404, errors.New("bridge name not found"))
		return
	}

	c.BindJSON(&form)
	err = form.Save()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	doc, err := form.Marshal()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Data(200, MediaType, doc)
}

// Destroy removes a specific Bridge.
func (btc *BridgeTypesController) Destroy(c *gin.Context) {
	name := c.Param("BridgeName")
	if bt, err := btc.App.GetStore().FindBridge(name); err == orm.ErrorNotFound {
		publicError(c, 404, errors.New("bridge name not found"))
	} else if err != nil {
		c.AbortWithError(500, fmt.Errorf("Error searching for bridge for BTC Destroy: %+v", err))
	} else if jobFounds, err := btc.App.GetStore().AnyJobWithType(name); err != nil {
		c.AbortWithError(500, fmt.Errorf("Error searching for associated jobs for BTC Destroy: %+v", err))
	} else if jobFounds {
		c.AbortWithError(409, fmt.Errorf("Can't remove the bridge because there are jobs associated with it: %+v", err))
	} else if err = btc.App.RemoveAdapter(&bt); err != nil {
		c.AbortWithError(StatusCodeForError(err), fmt.Errorf("failed to initialise BTC Destroy: %+v", err))
	} else {
		c.JSON(200, presenters.BridgeType{BridgeType: bt})
	}
}

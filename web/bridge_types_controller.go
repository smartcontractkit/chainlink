package web

import (
	"errors"
	"fmt"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
)

// BridgeTypesController manages BridgeType requests in the node.
type BridgeTypesController struct {
	App *services.ChainlinkApplication
}

// Create adds the BridgeType to the given context.
func (btc *BridgeTypesController) Create(c *gin.Context) {
	bt := &models.BridgeType{}

	if err := c.ShouldBindJSON(bt); err != nil {
		c.AbortWithError(500, err)
	} else if err = btc.App.AddAdapter(bt); err != nil {
		publicError(c, StatusCodeForError(err), err)
	} else {
		c.JSON(200, bt)
	}
}

// Index lists Bridges, one page at a time.
func (btc *BridgeTypesController) Index(c *gin.Context) {
	size, page, offset, err := ParsePaginatedRequest(c.Query("size"), c.Query("page"))
	if err != nil {
		publicError(c, 422, err)
		return
	}

	skip := storm.Skip(offset)
	limit := storm.Limit(size)

	var bridges []models.BridgeType

	count, err := btc.App.Store.Count(&models.BridgeType{})
	if err != nil {
		c.AbortWithError(500, fmt.Errorf("error getting count of bridges: %+v", err))
		return
	}
	if err := btc.App.Store.AllByIndex("Name", &bridges, skip, limit); err != nil {
		c.AbortWithError(500, fmt.Errorf("error fetching all bridges: %+v", err))
		return
	}
	pbt := make([]presenters.BridgeType, len(bridges))
	for i, j := range bridges {
		pbt[i] = presenters.BridgeType{BridgeType: j}
	}
	buffer, err := NewPaginatedResponse(*c.Request.URL, size, page, count, pbt)
	if err != nil {
		c.AbortWithError(500, fmt.Errorf("failed to marshal document: %+v", err))
	} else {
		c.Data(200, MediaType, buffer)
	}
}

// Show returns the details of a specific Bridge.
func (btc *BridgeTypesController) Show(c *gin.Context) {
	name := c.Param("BridgeName")
	if bt, err := btc.App.Store.FindBridge(name); err == storm.ErrNotFound {
		publicError(c, 404, errors.New("bridge name not found"))
	} else if err != nil {
		c.AbortWithError(500, err)
	} else {
		c.JSON(200, presenters.BridgeType{BridgeType: bt})
	}
}

// Destroy removes a specific Bridge.
func (btc *BridgeTypesController) Destroy(c *gin.Context) {
	name := c.Param("BridgeName")
	if bt, err := btc.App.Store.FindBridge(name); err == storm.ErrNotFound {
		publicError(c, 404, errors.New("bridge name not found"))
	} else if err != nil {
		c.AbortWithError(500, fmt.Errorf("Error searching for bridge for BTC Destroy: %+v", err))
	} else if err = btc.App.RemoveAdapter(&bt); err != nil {
		c.AbortWithError(StatusCodeForError(err), fmt.Errorf("failed to initialise BTC Destroy: %+v", err))
	} else {
		c.JSON(200, presenters.BridgeType{BridgeType: bt})
	}
}

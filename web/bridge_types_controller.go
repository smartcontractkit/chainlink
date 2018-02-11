package web

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
)

// BridgeTypesController manages BridgeType requests in the node.
type BridgeTypesController struct {
	App *services.ChainlinkApplication
}

// Create adds the BridgeType to the given context.
func (btc *BridgeTypesController) Create(c *gin.Context) {
	bt := &models.BridgeType{}

	if err := c.ShouldBindJSON(bt); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = btc.App.GetStore().Save(bt); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, bt)
	}
}

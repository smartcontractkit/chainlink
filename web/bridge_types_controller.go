package web

import (
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
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = btc.App.AddAdapter(bt); err != nil {
		c.JSON(StatusCodeForError(err), gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, bt)
	}
}

// Index lists Bridges, one page at a time.
func (btc *BridgeTypesController) Index(c *gin.Context) {
	size, page, offset, err := ParsePaginatedRequest(c.Query("size"), c.Query("page"))
	if err != nil {
		c.JSON(422, gin.H{
			"errors": []string{err.Error()},
		})
	}

	skip := storm.Skip(offset)
	limit := storm.Limit(size)

	var bridges []models.BridgeType

	count, err := btc.App.Store.Count(&models.BridgeType{})
	if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{fmt.Errorf("error getting count of Bridges: %+v", err).Error()},
		})
		return
	}
	if err := btc.App.Store.AllByIndex("Name", &bridges, skip, limit); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{fmt.Errorf("erorr fetching All Bridges: %+v", err).Error()},
		})
		return
	}
	pbt := make([]presenters.BridgeType, len(bridges))
	for i, j := range bridges {
		pbt[i] = presenters.BridgeType{BridgeType: j}
	}
	buffer, err := NewPaginatedResponse(*c.Request.URL, size, page, count, pbt)
	if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{fmt.Errorf("failed to marshal document: %+v", err).Error()},
		})
	} else {
		c.Data(200, MediaType, buffer)
	}
}

// Show returns the details of a specific Bridge.
func (btc *BridgeTypesController) Show(c *gin.Context) {
	name := c.Param("BridgeName")
	if bt, err := btc.App.Store.FindBridge(name); err == storm.ErrNotFound {
		c.JSON(404, gin.H{
			"errors": []string{"Bridge Name not found."},
		})
	} else if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, presenters.BridgeType{BridgeType: bt})
	}
}

// RemoveOne removes a specific Bridge.
func (btc *BridgeTypesController) RemoveOne(c *gin.Context) {
	name := c.Param("BridgeName")
	if bt, err := btc.App.Store.FindBridge(name); err == storm.ErrNotFound {
		c.JSON(404, gin.H{
			"errors": []string{"Bridge Name not found."},
		})
	} else if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = btc.App.RemoveAdapter(&bt); err != nil {
		fmt.Println([]string{err.Error()})
		c.JSON(StatusCodeForError(err), gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, presenters.BridgeType{BridgeType: bt})
	}
}

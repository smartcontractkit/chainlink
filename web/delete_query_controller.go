package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
)

// DeleteQueryController manages delete queries in the node.
type DeleteQueryController struct {
	App *services.ChainlinkApplication
}

// DeleteBridges removes several bridges based on the content of the JSON passed with the request.
func (dqc *DeleteQueryController) DeleteBridges(dq *models.DeleteQueryParams) ([]models.BridgeType, error) {
	respBridges, err := dqc.App.Store.AdvancedBridgeSearch(dq.Query)
	if err != nil {
		return nil, fmt.Errorf("Advanced Bridge Search error: %v", err)
	}
	for _, rmBridge := range respBridges {
		if err = dqc.App.RemoveAdapter(&rmBridge); err != nil {
			return nil, fmt.Errorf("Error deleting Bridge: %v", err)
		}
	}
	return respBridges, nil
}

// DeleteJobRuns removes several jobruns based on the content of the JSON passed with the request.
func (dqc *DeleteQueryController) DeleteJobRuns(dq *models.DeleteQueryParams) ([]models.JobRun, error) {
	respJobRuns, err := dqc.App.Store.AdvancedJobRunSearch(dq.Query)
	if err != nil {
		return nil, fmt.Errorf("Advanced JobRun Search error: %v", err)
	}
	for _, rmJobRun := range respJobRuns {
		if err = dqc.App.RemoveJobRun(&rmJobRun); err != nil {
			return nil, fmt.Errorf("Error deleting JobRun: %v", err)
		}
	}
	return respJobRuns, nil
}

// DeleteQuery removes several entries based on the content of the JSON passed with the request.
func (dqc *DeleteQueryController) DeleteQuery(c *gin.Context) {
	dq := &models.DeleteQueryParams{}
	if err := c.ShouldBindJSON(dq); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{fmt.Errorf("Error parsing query for DeleteQuery: %+v", err).Error()},
		})
		fmt.Println(err)
		return
	}

	switch dq.Collection {
	case "bridges":
		resp, err := dqc.DeleteBridges(dq)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"errors": []string{fmt.Errorf("DeleteQuery error: %+v", err).Error()},
			})
			return
		}
		c.JSON(200, resp)

	case "jobruns":
		resp, err := dqc.DeleteJobRuns(dq)
		if err != nil {
			c.JSON(500, gin.H{
				"errors": []string{fmt.Errorf("DeleteQuery error: %+v", err).Error()},
			})
			return
		}
		c.JSON(200, resp)

	default:
		c.JSON(500, gin.H{
			"errors": []string{fmt.Errorf("Collection not found: %+v", dq.Collection).Error()},
		})
		fmt.Println("Collection not found")
		return
	}
}

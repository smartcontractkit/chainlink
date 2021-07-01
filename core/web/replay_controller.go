package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type ReplayController struct {
	App chainlink.Application
}

// ReplayBlocks causes the node to process blocks again from the given block number
// Example:
//  "<application>/replay_blocks"
func (bdc *ReplayController) ReplayBlocks(c *gin.Context) {
	request := &models.ReplayBlocksRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	if err := models.ValidateReplayRequest(request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	if err := bdc.App.ReplayFromBlockNumber(request.BlockNumber); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "nil", http.StatusNoContent)
}

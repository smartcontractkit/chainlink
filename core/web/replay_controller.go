package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

type ReplayController struct {
	App chainlink.Application
}

// ReplayBlocks causes the node to process blocks again from the given block number
// Example:
//  "<application>/replay_blocks"
func (bdc *ReplayController) ReplayBlocks(c *gin.Context) {

	if c.Param("number") == "" {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.New("missing 'number' parameter"))
		return
	}

	blockNumber, err := strconv.ParseInt(c.Param("number"), 10, 0)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	if blockNumber < 0 {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("block number cannot be negative: %v", blockNumber))
		return
	}
	if err := bdc.App.ReplayFromBlockNumber(uint64(blockNumber)); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "nil", http.StatusNoContent)
}

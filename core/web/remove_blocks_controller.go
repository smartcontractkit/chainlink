package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type RemoveBlocksController struct {
	App chainlink.Application
}

// RemoveBlocks causes the LogPoller to remove data starting from the specified block number to the latest block
// Example:
//
//	"<application>/v2/remove_blocks/:from"
func (bdc *RemoveBlocksController) RemoveBlocks(c *gin.Context) {
	if c.Param("start") == "" {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.New("missing 'start' parameter"))
		return
	}

	blockNumber, err := strconv.ParseInt(c.Param("start"), 10, 0)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	if blockNumber < 0 {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("start cannot be negative: %v", blockNumber))
		return
	}

	chain, err := getChain(bdc.App.GetRelayers().LegacyEVMChains(), c.Query("evmChainID"))
	if err != nil {
		if errors.Is(err, ErrInvalidChainID) || errors.Is(err, ErrMultipleChains) || errors.Is(err, ErrMissingChainID) {
			jsonAPIError(c, http.StatusUnprocessableEntity, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	chainID := chain.ID()

	if err := bdc.App.DeleteLogPollerDataAfter(c.Request.Context(), chainID, blockNumber); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

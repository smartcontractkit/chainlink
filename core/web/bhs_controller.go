package web

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

type BHSController struct {
	App chainlink.Application
}

type backwardsModeRequest struct {
	startBlock      int64
	endBlock        int64
	batchBHSAddress common.Address
	batchSize       int64
}

type BackwardsResponse struct {
	Message    string   `json:"message"`
	EVMChainID *big.Int `json:"evmChainID"`
}

// BackwardsMode triggers a task to fill the blockhash store with blockhashes
// of the blocks in the block range specified by the request.
// Example:
// POST <app-url>/v2/bhs/backwards?start_block=<integer>&end_block=<integer>&batch_bhs_address=<address>&batch_size=<integer>
func (c *BHSController) BackwardsMode(cx *gin.Context) {
	r, err := c.validateAndExtractBackwardsParams(cx)
	if err != nil {
		jsonAPIError(cx, http.StatusUnprocessableEntity, err)
		return
	}

	// The caller of this endpoint will have to ensure that there is enough ETH
	// in all sending keys to perform the backwards operation.
	sendingKeys, err := c.App.GetKeyStore().Eth().SendingKeys()
	if err != nil {
		jsonAPIError(cx, http.StatusInternalServerError, err)
		return
	}

	if len(sendingKeys) == 0 {
		jsonAPIError(cx, http.StatusInternalServerError, errors.New("no sending keys configured on chainlink node"))
	}

	chain, err := getChain(c.App.GetChains().EVM, cx.Query("evmChainID"))
	backwardsBHS := blockhashstore.NewBackwardsBHS(
		chain.TxManager(),
		chain.Client(),
		sendingKeys[0].Address.Address(),
		c.App.GetLogger().Named("BackwardsBHS"),
	)

	// Run backwards feeding async, caller will have to check CL node logs
	// to observe progress.
	go func() {
		if err := backwardsBHS.Backwards(r.startBlock, r.endBlock, r.batchBHSAddress, r.batchSize); err != nil {
			c.App.GetLogger().Errorw("error encountered running backwards bhs", "err", err)
		}
	}()

	response := &BackwardsResponse{
		Message:    "backwards feed started",
		EVMChainID: chain.ID(),
	}
	jsonAPIResponse(cx, response, "response")
}

func (c *BHSController) validateAndExtractBackwardsParams(cx *gin.Context) (*backwardsModeRequest, error) {
	req := &backwardsModeRequest{}
	if startBlockS := cx.Param("start_block"); startBlockS == "" {
		return nil, errors.New("missing 'start_block' parameter")
	} else {
		startBlock, err := strconv.ParseInt(startBlockS, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "parse start_block parameter")
		}
		if startBlock <= 0 {
			return nil, fmt.Errorf("start_block must be positive, given %d", startBlock)
		}
		req.startBlock = startBlock
	}

	if endBlockS := cx.Param("end_block"); endBlockS == "" {
		return nil, errors.New("missing 'end_block' parameter")
	} else {
		endBlock, err := strconv.ParseInt(endBlockS, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "parse end_block parameter")
		}
		if endBlock <= 0 {
			return nil, fmt.Errorf("end_block must be positive, given %d", endBlock)
		}
		req.endBlock = endBlock
	}

	if req.startBlock < req.endBlock {
		return nil, fmt.Errorf("start_block %d cannot be less than end_block %d", req.startBlock, req.endBlock)
	}

	if addrS := cx.Param("batch_bhs_address"); addrS == "" {
		return nil, errors.New("missing 'batch_bhs_address' parameter")
	} else {
		addr := common.HexToAddress(addrS)
		if addr.Hex() == (common.Address{}).Hex() {
			return nil, fmt.Errorf("invalid address %s", addrS)
		}
		req.batchBHSAddress = addr
	}

	if batchSizeS := cx.Param("batch_size"); batchSizeS == "" {
		return nil, errors.New("missing 'batch_size' parameter")
	} else {
		batchSize, err := strconv.ParseInt(batchSizeS, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "parse batch_size parameter")
		}
		if batchSize <= 0 {
			return nil, fmt.Errorf("batch_size must be positive, given %d", batchSize)
		}
		req.batchSize = batchSize
	}

	return req, nil
}

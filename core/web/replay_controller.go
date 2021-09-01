package web

import (
	"math/big"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type ReplayController struct {
	App chainlink.Application
}

// ReplayFromBlock causes the node to process blocks again from the given block number
// Example:
//  "<application>/v2/replay_from_block/:number"
func (bdc *ReplayController) ReplayFromBlock(c *gin.Context) {
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
	var chainID *big.Int
	if bdc.App.GetChainSet().ChainCount() > 1 {
		if c.Param("evmChainID") == "" {
			jsonAPIError(c, http.StatusUnprocessableEntity, errors.New("more than one chain available, you must specify evmChainID parameter"))
			return
		}
		var ok bool
		chainID, ok = big.NewInt(0).SetString(c.Param("evmChainID"), 10)
		if !ok {
			jsonAPIError(c, http.StatusUnprocessableEntity, errors.New("invalid evmChainID"))
			return
		}
	} else {
		chain, err := bdc.App.GetChainSet().Default()
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
		chainID = chain.ID()
	}

	if err := bdc.App.ReplayFromBlock(chainID, uint64(blockNumber)); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	response := ReplayResponse{
		Message:    "Replay started",
		EVMChainID: utils.NewBig(chainID),
	}
	jsonAPIResponse(c, &response, "response")
}

type ReplayResponse struct {
	Message    string     `json:"message"`
	EVMChainID *utils.Big `json:"evmChainID"`
}

// GetID returns the jsonapi ID.
func (s ReplayResponse) GetID() string {
	return "replayID"
}

// GetName returns the collection name for jsonapi.
func (ReplayResponse) GetName() string {
	return "replay"
}

// SetID is used to conform to the UnmarshallIdentifier interface for
// deserializing from jsonapi documents.
func (*ReplayResponse) SetID(string) error {
	return nil
}

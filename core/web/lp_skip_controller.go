package web

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

type LPSkipController struct {
	App chainlink.Application
}

type LPSkipToBlockRequest struct {
	BlockNumber int64  `json:"blockNumber"`
	Family      string `json:"family"`
	ChainID     string `json:"chain-id"`
}

// LPSkipToBlock repositions the LogPoller to start processing from the given block number.
// Example:
//
//	"<application>/v2/lp_skip_to_block"
func (c *LPSkipController) LPSkipToBlock(gctx *gin.Context) {
	var request LPSkipToBlockRequest
	if err := gctx.ShouldBindJSON(&request); err != nil {
		jsonAPIError(gctx, http.StatusUnprocessableEntity, err)
		return
	}
	if request.BlockNumber < 2 {
		jsonAPIError(gctx, http.StatusUnprocessableEntity, errors.Errorf("block number must be >= 2: %v", request.BlockNumber))
		return
	}

	if request.Family == "" {
		jsonAPIError(gctx, http.StatusUnprocessableEntity, errors.New("chain family was not provided"))
		return
	}
	if request.Family != relay.NetworkEVM {
		jsonAPIError(gctx, http.StatusUnprocessableEntity, errors.Errorf("unsupported chain family %q, only %s is supported", request.Family, relay.NetworkEVM))
		return
	}

	if strings.TrimSpace(request.ChainID) == "" {
		jsonAPIError(gctx, http.StatusUnprocessableEntity, errors.New("chain-id was not provided"))
		return
	}

	ctx := gctx.Request.Context()
	if err := c.App.LPSkipToBlock(ctx, request.Family, request.ChainID, request.BlockNumber); err != nil {
		if errors.Is(err, chainlink.ErrNoSuchRelayer) {
			jsonAPIError(gctx, http.StatusBadRequest, err)
			return
		}
		jsonAPIError(gctx, http.StatusInternalServerError, err)
		return
	}

	response := LPSkipToBlockResponse{
		Message:     "Log poller will start processing from the new block on next tick",
		ChainID:     request.ChainID,
		BlockNumber: request.BlockNumber,
	}
	jsonAPIResponse(gctx, &response, "response")
}

type LPSkipToBlockResponse struct {
	Message     string `json:"message"`
	ChainID     string `json:"chain-id"`
	BlockNumber int64  `json:"blockNumber"`
}

// GetID returns the jsonapi ID.
func (s LPSkipToBlockResponse) GetID() string {
	return "lpSkipToBlockID"
}

// GetName returns the collection name for jsonapi.
func (LPSkipToBlockResponse) GetName() string {
	return "lp_skip_to_block"
}

// SetID is used to conform to the UnmarshallIdentifier interface for
// deserializing from jsonapi documents.
func (*LPSkipToBlockResponse) SetID(string) error {
	return nil
}

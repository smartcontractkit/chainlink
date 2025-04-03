package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type ReplayController struct {
	App chainlink.Application
}

// ReplayFromBlock causes the node to process blocks again from the given block number
// Example:
//
//	"<application>/v2/replay_from_block/:number"
func (bdc *ReplayController) ReplayFromBlock(c *gin.Context) {
	if c.Param("number") == "" {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.New("missing 'number' parameter"))
		return
	}

	// check if "force" query string parameter provided
	var force bool
	var err error
	if fb := c.Query("force"); fb != "" {
		force, err = strconv.ParseBool(fb)
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, errors.Wrap(err, "boolean value required for 'force' query string param"))
			return
		}
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

	chainFamily := c.Query("family")
	if chainFamily == "" {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.New("chain family was not provoded"))
		return
	}

	chainID := c.Query("ChainID")
	if chainID == "" {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.New("chain-id was not provoded"))
		return
	}

	if chainFamily == "evm" {
		_, err := getChain(bdc.App.GetRelayers().LegacyEVMChains(), c.Query("ChainID"))
		if err != nil {
			if errors.Is(err, ErrInvalidChainID) || errors.Is(err, ErrMultipleChains) || errors.Is(err, ErrMissingChainID) {
				jsonAPIError(c, http.StatusUnprocessableEntity, err)
				return
			}
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
	}

	ctx := c.Request.Context()
	if err := bdc.App.ReplayFromBlock(ctx, chainFamily, chainID, uint64(blockNumber), force); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	response := ReplayResponse{
		Message: "Replay started",
		ChainID: chainID,
	}
	jsonAPIResponse(c, &response, "response")
}

type ReplayResponse struct {
	Message string `json:"message"`
	ChainID string `json:"ChainID"`
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

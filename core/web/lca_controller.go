package web

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type LCAController struct {
	App chainlink.Application
}

// FindLCA compares chain of blocks available in the DB with chain provided by an RPC and returns last common ancestor
// Example:
//
//	"<application>/v2/find_lca"
func (bdc *LCAController) FindLCA(c *gin.Context) {
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

	lca, err := bdc.App.FindLCA(c.Request.Context(), chainID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	if lca == nil {
		jsonAPIError(c, http.StatusNotFound, errors.New("failed to find last common ancestor"))
		return
	}

	response := LCAResponse{
		BlockNumber: lca.BlockNumber,
		Hash:        lca.BlockHash.String(),
		EVMChainID:  big.New(chainID),
	}
	jsonAPIResponse(c, &response, "response")
}

type LCAResponse struct {
	BlockNumber int64    `json:"blockNumber"`
	Hash        string   `json:"hash"`
	EVMChainID  *big.Big `json:"evmChainID"`
}

// GetID returns the jsonapi ID.
func (s LCAResponse) GetID() string {
	return "LCAResponseID"
}

// GetName returns the collection name for jsonapi.
func (LCAResponse) GetName() string {
	return "lca_response"
}

// SetID is used to conform to the UnmarshallIdentifier interface for
// deserializing from jsonapi documents.
func (*LCAResponse) SetID(string) error {
	return nil
}

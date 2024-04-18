package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
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
	println("find_lca get chain")
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

	println("find_lca got chain")
	lca, err := bdc.App.FindLCA(c.Request.Context(), chainID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	println("find_lca found")

	if lca == nil {
		jsonAPIError(c, http.StatusNotFound, fmt.Errorf("failed to find last common ancestor"))
		return
	}

	println("find_lca found not nil")

	response := LCAResponse{
		BlockNumber: lca.BlockNumber,
		Hash:        lca.BlockHash.String(),
		EVMChainID:  big.New(chainID),
	}
	jsonAPIResponse(c, &response, "response")

	println("find_lca rendered: " + lca.BlockHash.String())
}

type LCAResponse struct {
	BlockNumber int64    `json:"block_number"`
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

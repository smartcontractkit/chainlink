package web

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// SolanaChainsController manages Solana chains.
type SolanaChainsController struct {
	App chainlink.Application
}

// Index lists Solana chains.
func (cc *SolanaChainsController) Index(c *gin.Context, size, page, offset int) {
	solanaChains := cc.App.GetChains().Solana
	if solanaChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrSolanaNotEnabled)
		return
	}
	chains, count, err := solanaChains.ORM().Chains(offset, size)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	var resources []presenters.SolanaChainResource
	for _, chain := range chains {
		resources = append(resources, presenters.NewSolanaChainResource(chain))
	}

	paginatedResponse(c, "solana_chain", size, page, resources, count, err)
}

// CreateSolanaChainRequest is a JSONAPI request for creating a Solana chain.
type CreateSolanaChainRequest struct {
	ID     string      `json:"chainID"`
	Config db.ChainCfg `json:"config"`
}

// Show gets a Solana chain by chain id.
func (cc *SolanaChainsController) Show(c *gin.Context) {
	solanaChains := cc.App.GetChains().Solana
	if solanaChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrSolanaNotEnabled)
		return
	}
	chain, err := solanaChains.ORM().Chain(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, presenters.NewSolanaChainResource(chain), "solana_chain")
}

// Create adds a new Solana chain.
func (cc *SolanaChainsController) Create(c *gin.Context) {
	solanaChains := cc.App.GetChains().Solana
	if solanaChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrSolanaNotEnabled)
		return
	}

	request := &CreateSolanaChainRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	chain, err := solanaChains.Add(c.Request.Context(), request.ID, request.Config)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponseWithStatus(c, presenters.NewSolanaChainResource(chain), "solana_chain", http.StatusCreated)
}

// UpdateSolanaChainRequest is a JSONAPI request for updating a Solana chain.
type UpdateSolanaChainRequest struct {
	Enabled bool        `json:"enabled"`
	Config  db.ChainCfg `json:"config"`
}

// Update configures an existing Solana chain.
func (cc *SolanaChainsController) Update(c *gin.Context) {
	solanaChains := cc.App.GetChains().Solana
	if solanaChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrSolanaNotEnabled)
		return
	}

	var request UpdateSolanaChainRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	chain, err := solanaChains.Configure(c.Request.Context(), c.Param("ID"), request.Enabled, request.Config)

	if errors.Is(err, sql.ErrNoRows) {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	} else if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, presenters.NewSolanaChainResource(chain), "solana_chain")
}

// Delete removes a Solana chain.
func (cc *SolanaChainsController) Delete(c *gin.Context) {
	solanaChains := cc.App.GetChains().Solana
	if solanaChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrSolanaNotEnabled)
		return
	}
	err := solanaChains.Remove(c.Param("ID"))

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "solana_chain", http.StatusNoContent)
}

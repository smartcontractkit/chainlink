package web

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// TerraChainsController manages Terra chains.
type TerraChainsController struct {
	App chainlink.Application
}

// Index lists Terra chains.
func (cc *TerraChainsController) Index(c *gin.Context, size, page, offset int) {
	terraChains := cc.App.GetChains().Terra
	if terraChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrTerraNotEnabled)
		return
	}
	chains, count, err := terraChains.ORM().Chains(offset, size)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	var resources []presenters.TerraChainResource
	for _, chain := range chains {
		resources = append(resources, presenters.NewTerraChainResource(chain))
	}

	paginatedResponse(c, "terra_chain", size, page, resources, count, err)
}

// CreateTerraChainRequest is a JSONAPI request for creating a Terra chain.
type CreateTerraChainRequest struct {
	ID     string      `json:"chainID"`
	Config db.ChainCfg `json:"config"`
}

// Show gets a Terra chain by chain id.
func (cc *TerraChainsController) Show(c *gin.Context) {
	terraChains := cc.App.GetChains().Terra
	if terraChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrTerraNotEnabled)
		return
	}
	chain, err := terraChains.ORM().Chain(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, presenters.NewTerraChainResource(chain), "terra_chain")
}

// Create adds a new Terra chain.
func (cc *TerraChainsController) Create(c *gin.Context) {
	terraChains := cc.App.GetChains().Terra
	if terraChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrTerraNotEnabled)
		return
	}

	request := &CreateTerraChainRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	chain, err := terraChains.Add(c.Request.Context(), request.ID, request.Config)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponseWithStatus(c, presenters.NewTerraChainResource(chain), "terra_chain", http.StatusCreated)
}

// UpdateTerraChainRequest is a JSONAPI request for updating a Terra chain.
type UpdateTerraChainRequest struct {
	Enabled bool        `json:"enabled"`
	Config  db.ChainCfg `json:"config"`
}

// Update configures an existing Terra chain.
func (cc *TerraChainsController) Update(c *gin.Context) {
	terraChains := cc.App.GetChains().Terra
	if terraChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrTerraNotEnabled)
		return
	}

	var request UpdateTerraChainRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	chain, err := terraChains.Configure(c.Request.Context(), c.Param("ID"), request.Enabled, request.Config)

	if errors.Is(err, sql.ErrNoRows) {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	} else if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, presenters.NewTerraChainResource(chain), "terra_chain")
}

// Delete removes a Terra chain.
func (cc *TerraChainsController) Delete(c *gin.Context) {
	terraChains := cc.App.GetChains().Terra
	if terraChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrTerraNotEnabled)
		return
	}
	err := terraChains.Remove(c.Param("ID"))

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "terra_chain", http.StatusNoContent)
}

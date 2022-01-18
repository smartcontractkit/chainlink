package web

import (
	"database/sql"
	"net/http"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/gin-gonic/gin"
)

type TerraChainsController struct {
	App chainlink.Application
}

func (cc *TerraChainsController) Index(c *gin.Context, size, page, offset int) {
	chains, count, err := cc.App.TerraORM().Chains(offset, size)

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

type CreateTerraChainRequest struct {
	ID     string         `json:"chainID"`
	Config types.ChainCfg `json:"config"`
}

func (cc *TerraChainsController) Show(c *gin.Context) {
	chain, err := cc.App.TerraORM().Chain(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, presenters.NewTerraChainResource(chain), "terra_chain")
}

func (cc *TerraChainsController) Create(c *gin.Context) {
	request := &CreateTerraChainRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	chain, err := cc.App.GetChains().Terra.Add(request.ID, request.Config)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponseWithStatus(c, presenters.NewTerraChainResource(chain), "terra_chain", http.StatusCreated)
}

type UpdateTerraChainRequest struct {
	Enabled bool           `json:"enabled"`
	Config  types.ChainCfg `json:"config"`
}

func (cc *TerraChainsController) Update(c *gin.Context) {
	var request UpdateTerraChainRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	chain, err := cc.App.GetChains().Terra.Configure(c.Param("ID"), request.Enabled, request.Config)

	if errors.Is(err, sql.ErrNoRows) {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	} else if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, presenters.NewTerraChainResource(chain), "terra_chain")
}

func (cc *TerraChainsController) Delete(c *gin.Context) {
	err := cc.App.GetChains().Terra.Remove(c.Param("ID"))

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "terra_chain", http.StatusNoContent)
}

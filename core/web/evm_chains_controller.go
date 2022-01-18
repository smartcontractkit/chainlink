package web

import (
	"database/sql"
	"net/http"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/gin-gonic/gin"
)

type EVMChainsController struct {
	App chainlink.Application
}

func (cc *EVMChainsController) Index(c *gin.Context, size, page, offset int) {
	chains, count, err := cc.App.EVMORM().Chains(offset, size)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	var resources []presenters.EVMChainResource
	for _, chain := range chains {
		resources = append(resources, presenters.NewEVMChainResource(chain))
	}

	paginatedResponse(c, "chain", size, page, resources, count, err)
}

type CreateEVMChainRequest struct {
	ID     utils.Big      `json:"chainID"`
	Config types.ChainCfg `json:"config"`
}

func (cc *EVMChainsController) Show(c *gin.Context) {
	id := utils.Big{}
	err := id.UnmarshalText([]byte(c.Param("ID")))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	chain, err := cc.App.EVMORM().Chain(id)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, presenters.NewEVMChainResource(chain), "chain")
}

func (cc *EVMChainsController) Create(c *gin.Context) {
	request := &CreateEVMChainRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	chain, err := cc.App.GetChains().EVM.Add(request.ID.ToInt(), request.Config)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponseWithStatus(c, presenters.NewEVMChainResource(chain), "chain", http.StatusCreated)
}

type UpdateEVMChainRequest struct {
	Enabled bool           `json:"enabled"`
	Config  types.ChainCfg `json:"config"`
}

func (cc *EVMChainsController) Update(c *gin.Context) {
	id := utils.Big{}
	err := id.UnmarshalText([]byte(c.Param("ID")))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	var request UpdateEVMChainRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	chain, err := cc.App.GetChains().EVM.Configure(id.ToInt(), request.Enabled, request.Config)

	if errors.Is(err, sql.ErrNoRows) {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	} else if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, presenters.NewEVMChainResource(chain), "chain")
}

func (cc *EVMChainsController) Delete(c *gin.Context) {
	id := utils.Big{}
	err := id.UnmarshalText([]byte(c.Param("ID")))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = cc.App.GetChains().EVM.Remove(id.ToInt())

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "chain", http.StatusNoContent)
}

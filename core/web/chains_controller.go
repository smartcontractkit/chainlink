package web

import (
	"net/http"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/gin-gonic/gin"
)

type ChainsController struct {
	App chainlink.Application
}

func (cc *ChainsController) Index(c *gin.Context, size, page, offset int) {
	chains, count, err := cc.App.EVMORM().Chains(offset, size)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	var resources []presenters.ChainResource
	for _, chain := range chains {
		resources = append(resources, presenters.NewChainResource(chain))
	}

	paginatedResponse(c, "chain", size, page, resources, count, err)
}

type CreateChainRequest struct {
	ID     utils.Big      `json:"chainID"`
	Config types.ChainCfg `json:"config"`
}

func (cc *ChainsController) Create(c *gin.Context) {
	request := &CreateChainRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	chain, err := cc.App.EVMORM().CreateChain(request.ID, request.Config)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, presenters.NewChainResource(chain), "chain")
}

func (cc *ChainsController) Delete(c *gin.Context) {
	id := utils.Big{}
	err := id.UnmarshalText([]byte(c.Param("ID")))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = cc.App.EVMORM().DeleteChain(id)

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "chain", http.StatusNoContent)
}

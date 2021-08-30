package web

import (
	"net/http"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/gin-gonic/gin"
)

type NodesController struct {
	App chainlink.Application
}

func (nc *NodesController) Index(c *gin.Context, size, page, offset int) {
	id := c.Param("ID")

	var nodes []types.Node
	var count int
	var err error

	if id == "" {
		// fetch all nodes
		nodes, count, err = nc.App.EVMORM().Nodes(offset, size)

	} else {
		// fetch nodes for chain ID
		chainID := utils.Big{}
		if err = chainID.UnmarshalText([]byte(id)); err != nil {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
		nodes, count, err = nc.App.EVMORM().NodesForChain(chainID, offset, size)
	}

	var resources []presenters.NodeResource
	for _, node := range nodes {
		resources = append(resources, presenters.NewNodeResource(node))
	}

	paginatedResponse(c, "node", size, page, resources, count, err)
}

func (nc *NodesController) Create(c *gin.Context) {
	var request evm.NewNode

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	node, err := nc.App.EVMORM().CreateNode(request)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, presenters.NewNodeResource(node), "node")
}

func (nc *NodesController) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("ID"), 10, 64)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = nc.App.EVMORM().DeleteNode(id)

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "node", http.StatusNoContent)
}

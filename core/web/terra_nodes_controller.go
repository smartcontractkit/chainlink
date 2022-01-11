package web

import (
	"net/http"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/gin-gonic/gin"
)

// TerraNodesController manages Terra nodes.
type TerraNodesController struct {
	App chainlink.Application
}

// Index lists Terra nodes, and optionally filters by chain id.
func (nc *TerraNodesController) Index(c *gin.Context, size, page, offset int) {
	id := c.Param("ID")

	var nodes []types.Node
	var count int
	var err error

	if id == "" {
		// fetch all nodes
		nodes, count, err = nc.App.TerraORM().Nodes(offset, size)

	} else {
		nodes, count, err = nc.App.TerraORM().NodesForChain(id, offset, size)
	}

	var resources []presenters.TerraNodeResource
	for _, node := range nodes {
		resources = append(resources, presenters.NewTerraNodeResource(node))
	}

	paginatedResponse(c, "node", size, page, resources, count, err)
}

// Create adds a new Terra node.
func (nc *TerraNodesController) Create(c *gin.Context) {
	var request types.NewNode

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	node, err := nc.App.TerraORM().CreateNode(request)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, presenters.NewTerraNodeResource(node), "node")
}

// Delete removes a Terra node.
func (nc *TerraNodesController) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("ID"), 10, 32)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = nc.App.TerraORM().DeleteNode(int32(id))

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "node", http.StatusNoContent)
}

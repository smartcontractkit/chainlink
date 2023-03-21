package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/logger/audit"
)

type NodesController interface {
	// Index lists nodes, and optionally filters by chain id.
	Index(c *gin.Context, size, page, offset int)
}

type nodesController[I chains.ID, N chains.Node, R jsonapi.EntityNamer] struct {
	nodeSet       chains.Nodes[I, N]
	parseChainID  func(string) (I, error)
	errNotEnabled error
	newResource   func(N) R
	auditLogger   audit.AuditLogger
}

func newNodesController[I chains.ID, N chains.Node, R jsonapi.EntityNamer](
	nodeSet chains.Nodes[I, N],
	errNotEnabled error,
	parseChainID func(string) (I, error),
	newResource func(N) R,
	auditLogger audit.AuditLogger,
) NodesController {
	return &nodesController[I, N, R]{
		nodeSet:       nodeSet,
		errNotEnabled: errNotEnabled,
		parseChainID:  parseChainID,
		newResource:   newResource,
		auditLogger:   auditLogger,
	}
}

func (n *nodesController[I, N, R]) Index(c *gin.Context, size, page, offset int) {
	if n.nodeSet == nil {
		jsonAPIError(c, http.StatusBadRequest, n.errNotEnabled)
		return
	}

	id := c.Param("ID")

	var nodes []N
	var count int
	var err error

	if id == "" {
		// fetch all nodes
		nodes, count, err = n.nodeSet.GetNodes(c, offset, size)
	} else {
		// fetch nodes for chain ID
		chainID, err2 := n.parseChainID(id)
		if err2 != nil {
			jsonAPIError(c, http.StatusBadRequest, err2)
			return
		}
		nodes, count, err = n.nodeSet.GetNodesForChain(c, chainID, offset, size)
	}

	var resources []R
	for _, node := range nodes {
		res := n.newResource(node)
		resources = append(resources, res)
	}

	paginatedResponse(c, "node", size, page, resources, count, err)
}

package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
)

type NodesController interface {
	// Index lists nodes, and optionally filters by chain id.
	Index(c *gin.Context, size, page, offset int)
}

type nodesController[R jsonapi.EntityNamer] struct {
	nodeSet       chains.Nodes
	errNotEnabled error
	newResource   func(status types.NodeStatus) R
	auditLogger   audit.AuditLogger
}

func newNodesController[R jsonapi.EntityNamer](
	nodeSet chains.Nodes,
	errNotEnabled error,
	newResource func(status types.NodeStatus) R,
	auditLogger audit.AuditLogger,
) NodesController {
	return &nodesController[R]{
		nodeSet:       nodeSet,
		errNotEnabled: errNotEnabled,
		newResource:   newResource,
		auditLogger:   auditLogger,
	}
}

func (n *nodesController[R]) Index(c *gin.Context, size, page, offset int) {
	if n.nodeSet == nil {
		jsonAPIError(c, http.StatusBadRequest, n.errNotEnabled)
		return
	}

	id := c.Param("ID")

	var nodes []types.NodeStatus
	var count int
	var err error

	if id == "" {
		// fetch all nodes
		nodes, count, err = n.nodeSet.NodeStatuses(c, offset, size)
	} else {
		// fetch nodes for chain ID
		nodes, count, err = n.nodeSet.NodeStatuses(c, offset, size, id)
	}

	var resources []R
	for _, node := range nodes {
		res := n.newResource(node)
		resources = append(resources, res)
	}

	paginatedResponse(c, "node", size, page, resources, count, err)
}

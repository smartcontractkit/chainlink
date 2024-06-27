package web

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type NodesController interface {
	// Index lists nodes, and optionally filters by chain id.
	Index(c *gin.Context, size, page, offset int)
}

type NetworkScopedNodeStatuser struct {
	network  string
	relayers chainlink.RelayerChainInteroperators
}

func NewNetworkScopedNodeStatuser(relayers chainlink.RelayerChainInteroperators, network string) *NetworkScopedNodeStatuser {
	scoped := relayers.List(chainlink.FilterRelayersByType(network))
	return &NetworkScopedNodeStatuser{
		network:  network,
		relayers: scoped,
	}
}

func (n *NetworkScopedNodeStatuser) NodeStatuses(ctx context.Context, offset, limit int, relayIDs ...types.RelayID) (nodes []types.NodeStatus, count int, err error) {
	return n.relayers.NodeStatuses(ctx, offset, limit, relayIDs...)
}

type nodesController[R jsonapi.EntityNamer] struct {
	nodeSet       *NetworkScopedNodeStatuser
	errNotEnabled error
	newResource   func(status types.NodeStatus) R
	auditLogger   audit.AuditLogger
}

func newNodesController[R jsonapi.EntityNamer](
	nodeSet *NetworkScopedNodeStatuser,
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
		// backward compatibility
		var rid types.RelayID
		err = rid.UnmarshalString(id)
		if err != nil {
			rid.ChainID = id
			rid.Network = n.nodeSet.network
		}
		nodes, count, err = n.nodeSet.NodeStatuses(c, offset, size, rid)
	}

	var resources []R
	for _, node := range nodes {
		res := n.newResource(node)
		resources = append(resources, res)
	}

	paginatedResponse(c, "node", size, page, resources, count, err)
}

package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/logger/audit"
)

type NodesController interface {
	// Index lists nodes, and optionally filters by chain id.
	Index(c *gin.Context, size, page, offset int)
	// Create adds a new node.
	Create(*gin.Context)
	// Delete removes a node.
	Delete(*gin.Context)
}

type nodesController[I chains.ID, N chains.Node, R jsonapi.EntityNamer] struct {
	nodeSet       chains.DBNodeSet[I, N]
	parseChainID  func(string) (I, error)
	errNotEnabled error
	newResource   func(N) R
	createNode    func(*gin.Context) (N, error)
	auditLogger   audit.AuditLogger
}

func newNodesController[I chains.ID, N chains.Node, R jsonapi.EntityNamer](
	nodeSet chains.DBNodeSet[I, N],
	errNotEnabled error,
	parseChainID func(string) (I, error),
	newResource func(N) R,
	createNode func(*gin.Context) (N, error),
	lggr logger.Logger,
	auditLogger audit.AuditLogger,
) NodesController {
	return &nodesController[I, N, R]{
		nodeSet:       nodeSet,
		errNotEnabled: errNotEnabled,
		parseChainID:  parseChainID,
		newResource:   newResource,
		createNode:    createNode,
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

func (n *nodesController[I, N, R]) Create(c *gin.Context) {
	if n.nodeSet == nil {
		jsonAPIError(c, http.StatusBadRequest, n.errNotEnabled)
		return
	}

	request, err := n.createNode(c)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	node, err := n.nodeSet.CreateNode(c, request)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	n.auditLogger.Audit(audit.ChainRpcNodeAdded, map[string]interface{}{})

	jsonAPIResponse(c, n.newResource(node), "node")
}

func (n *nodesController[I, N, R]) Delete(c *gin.Context) {
	if n.nodeSet == nil {
		jsonAPIError(c, http.StatusBadRequest, n.errNotEnabled)
		return
	}

	id, err := strconv.ParseInt(c.Param("ID"), 10, 32)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = n.nodeSet.DeleteNode(c, int32(id))

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	n.auditLogger.Audit(audit.ChainDeleted, map[string]interface{}{"id": id})

	jsonAPIResponseWithStatus(c, nil, "node", http.StatusNoContent)
}

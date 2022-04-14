package web

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// ErrSolanaNotEnabled is returned when SOLANA_ENABLED is not true.
var ErrSolanaNotEnabled = errors.New("Solana is disabled. Set SOLANA_ENABLED=true to enable.")

// SolanaNodesController manages Solana nodes.
type SolanaNodesController struct {
	App chainlink.Application
}

// Index lists Solana nodes, and optionally filters by chain id.
func (nc *SolanaNodesController) Index(c *gin.Context, size, page, offset int) {
	solanaChains := nc.App.GetChains().Solana
	if solanaChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrSolanaNotEnabled)
		return
	}
	orm := solanaChains.ORM()

	id := c.Param("ID")

	var nodes []db.Node
	var count int
	var err error

	if id == "" {
		// fetch all nodes
		nodes, count, err = orm.Nodes(offset, size)

	} else {
		nodes, count, err = orm.NodesForChain(id, offset, size)
	}

	var resources []presenters.SolanaNodeResource
	for _, node := range nodes {
		resources = append(resources, presenters.NewSolanaNodeResource(node))
	}

	paginatedResponse(c, "node", size, page, resources, count, err)
}

// Create adds a new Solana node.
func (nc *SolanaNodesController) Create(c *gin.Context) {
	solanaChains := nc.App.GetChains().Solana
	if solanaChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrSolanaNotEnabled)
		return
	}
	orm := solanaChains.ORM()

	var request db.NewNode

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Ensure chain exists.
	if _, err := orm.Chain(request.SolanaChainID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			jsonAPIError(c, http.StatusBadRequest, errors.Errorf("Solana chain %s must be added first", request.SolanaChainID))
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	node, err := orm.CreateNode(request)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, presenters.NewSolanaNodeResource(node), "node")
}

// Delete removes a Solana node.
func (nc *SolanaNodesController) Delete(c *gin.Context) {
	solanaChains := nc.App.GetChains().Solana
	if solanaChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrSolanaNotEnabled)
		return
	}
	orm := solanaChains.ORM()

	id, err := strconv.ParseInt(c.Param("ID"), 10, 32)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = orm.DeleteNode(int32(id))

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "node", http.StatusNoContent)
}

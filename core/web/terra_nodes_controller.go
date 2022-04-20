package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// ErrTerraNotEnabled is returned when TERRA_ENABLED is not true.
var ErrTerraNotEnabled = errors.New("Terra is disabled. Set TERRA_ENABLED=true to enable.")

// TerraNodesController manages Terra nodes.
type TerraNodesController struct {
	App chainlink.Application
}

// Index lists Terra nodes, and optionally filters by chain id.
func (nc *TerraNodesController) Index(c *gin.Context, size, page, offset int) {
	terraChains := nc.App.GetChains().Terra
	if terraChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrTerraNotEnabled)
		return
	}
	orm := terraChains.ORM()

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

	var resources []presenters.TerraNodeResource
	for _, node := range nodes {
		resources = append(resources, presenters.NewTerraNodeResource(node))
	}

	paginatedResponse(c, "node", size, page, resources, count, err)
}

// Create adds a new Terra node.
func (nc *TerraNodesController) Create(c *gin.Context) {
	terraChains := nc.App.GetChains().Terra
	if terraChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrTerraNotEnabled)
		return
	}
	orm := terraChains.ORM()

	var request types.NewNode

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Ensure chain exists.
	if _, err := orm.Chain(request.TerraChainID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("Terra chain %s must be added first", request.TerraChainID))
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	node, err := orm.CreateNode(db.Node{
		Name:          request.Name,
		TerraChainID:  request.TerraChainID,
		TendermintURL: request.TendermintURL,
	})

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, presenters.NewTerraNodeResource(node), "node")
}

// Delete removes a Terra node.
func (nc *TerraNodesController) Delete(c *gin.Context) {
	terraChains := nc.App.GetChains().Terra
	if terraChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrTerraNotEnabled)
		return
	}
	orm := terraChains.ORM()

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

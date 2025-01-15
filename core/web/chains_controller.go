package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type ChainsController interface {
	// Index lists chains.
	Index(c *gin.Context, size, page, offset int)
	// Show gets a chain by id.
	Show(*gin.Context)
}

type errChainDisabled struct {
	name    string
	tomlKey string
}

func (e errChainDisabled) Error() string {
	return fmt.Sprintf("%s is disabled: Set %s=true to enable", e.name, e.tomlKey)
}

type chainsController struct {
	chainStats  chainlink.RelayerChainInteroperators
	newResource func(chainlink.NetworkChainStatus) presenters.ChainResource
	lggr        logger.Logger
	auditLogger audit.AuditLogger
}

func NewChainsController(chainStats chainlink.RelayerChainInteroperators, lggr logger.Logger, auditLogger audit.AuditLogger) *chainsController {
	return &chainsController{
		chainStats:  chainStats,
		newResource: presenters.NewChainResource,
		lggr:        lggr,
		auditLogger: auditLogger,
	}
}

func (cc *chainsController) Index(c *gin.Context, size, page, offset int) {
	chainStats := cc.chainStats
	if network := c.Param("network"); network != "" {
		chainStats = chainStats.List(chainlink.FilterRelayersByType(network))
	}

	chains, count, err := chainStats.ChainStatuses(c.Request.Context(), offset, size)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	resources := []presenters.ChainResource{}
	for _, chain := range chains {
		resources = append(resources, cc.newResource(chain))
	}

	paginatedResponse(c, "chain", size, page, resources, count, err)
}

func (cc *chainsController) Show(c *gin.Context) {
	relayID := types.RelayID{Network: c.Param("network"), ChainID: c.Param("ID")}
	chain, err := cc.chainStats.ChainStatus(c.Request.Context(), relayID)
	status := chainlink.NetworkChainStatus{ChainStatus: chain, Network: relayID.Network}
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, cc.newResource(status), "chain")
}

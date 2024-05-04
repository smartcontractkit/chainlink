package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type ChainsController interface {
	// Index lists chains.
	Index(c *gin.Context, size, page, offset int)
	// Show gets a chain by id.
	Show(*gin.Context)
}

type chainsController[R jsonapi.EntityNamer] struct {
	network       string
	resourceName  string
	chainStats    chainlink.ChainStatuser
	errNotEnabled error
	newResource   func(types.ChainStatus) R
	lggr          logger.Logger
	auditLogger   audit.AuditLogger
}

type errChainDisabled struct {
	name    string
	tomlKey string
}

func (e errChainDisabled) Error() string {
	return fmt.Sprintf("%s is disabled: Set %s=true to enable", e.name, e.tomlKey)
}

func newChainsController[R jsonapi.EntityNamer](network string, chainStats chainlink.ChainsNodesStatuser, errNotEnabled error,
	newResource func(types.ChainStatus) R, lggr logger.Logger, auditLogger audit.AuditLogger) *chainsController[R] {
	return &chainsController[R]{
		network:       network,
		resourceName:  network + "_chain",
		chainStats:    chainStats,
		errNotEnabled: errNotEnabled,
		newResource:   newResource,
		lggr:          lggr,
		auditLogger:   auditLogger,
	}
}

func (cc *chainsController[R]) Index(c *gin.Context, size, page, offset int) {
	if cc.chainStats == nil {
		jsonAPIError(c, http.StatusBadRequest, cc.errNotEnabled)
		return
	}
	chains, count, err := cc.chainStats.ChainStatuses(c, offset, size)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	var resources []R
	for _, chain := range chains {
		resources = append(resources, cc.newResource(chain))
	}

	paginatedResponse(c, cc.resourceName, size, page, resources, count, err)
}

func (cc *chainsController[R]) Show(c *gin.Context) {
	if cc.chainStats == nil {
		jsonAPIError(c, http.StatusBadRequest, cc.errNotEnabled)
		return
	}
	relayID := types.RelayID{Network: cc.network, ChainID: c.Param("ID")}
	chain, err := cc.chainStats.ChainStatus(c, relayID)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, cc.newResource(chain), cc.resourceName)
}

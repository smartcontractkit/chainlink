package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
)

type ChainsController interface {
	// Index lists chains.
	Index(c *gin.Context, size, page, offset int)
	// Show gets a chain by id.
	Show(*gin.Context)
}

type chainsController[I chains.ID, R jsonapi.EntityNamer] struct {
	resourceName  string
	chainSet      chains.Chains[I]
	errNotEnabled error
	parseChainID  func(string) (I, error)
	newResource   func(chains.ChainConfig) R
	lggr          logger.Logger
	auditLogger   audit.AuditLogger
}

type errChainDisabled struct {
	name   string
	envVar string
}

func (e errChainDisabled) Error() string {
	return fmt.Sprintf("%s is disabled: Set %s=true to enable", e.name, e.envVar)
}

func newChainsController[I chains.ID, R jsonapi.EntityNamer](prefix string, chainSet chains.Chains[I], errNotEnabled error,
	parseChainID func(string) (I, error), newResource func(chains.ChainConfig) R, lggr logger.Logger, auditLogger audit.AuditLogger) *chainsController[I, R] {
	return &chainsController[I, R]{
		resourceName:  prefix + "_chain",
		chainSet:      chainSet,
		errNotEnabled: errNotEnabled,
		parseChainID:  parseChainID,
		newResource:   newResource,
		lggr:          lggr,
		auditLogger:   auditLogger,
	}
}

func (cc *chainsController[I, R]) Index(c *gin.Context, size, page, offset int) {
	if cc.chainSet == nil {
		jsonAPIError(c, http.StatusBadRequest, cc.errNotEnabled)
		return
	}
	chains, count, err := cc.chainSet.Index(offset, size)

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

func (cc *chainsController[I, R]) Show(c *gin.Context) {
	if cc.chainSet == nil {
		jsonAPIError(c, http.StatusBadRequest, cc.errNotEnabled)
		return
	}
	id, err := cc.parseChainID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	chain, err := cc.chainSet.Show(id)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, cc.newResource(chain), cc.resourceName)
}

package web

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains"
)

type ChainSet[I, C any] interface {
	Add(ctx context.Context, id I, cfg C) (chains.Chain[I, C], error)
	Show(id I) (chains.Chain[I, C], error)
	Configure(ctx context.Context, id I, enabled bool, cfg C) (chains.Chain[I, C], error)
	Remove(id I) error
	Index(offset, limit int) ([]chains.Chain[I, C], int, error)
}

type ChainsController interface {
	// Index lists chains.
	Index(c *gin.Context, size, page, offset int)
	// Create adds a new chain from a CreateChainRequest.
	Create(*gin.Context)
	// Show gets a chain by id.
	Show(*gin.Context)
	// Update configures an existing chain from an UpdateChainRequest.
	Update(*gin.Context)
	// Delete removes a chain.
	Delete(*gin.Context)
}

type chainsController[I, C, R any] struct {
	resourceName  string
	chainSet      ChainSet[I, C]
	errNotEnabled error
	parseChainID  func(string) (I, error)
	newResource   func(chains.Chain[I, C]) R
}

func NewChainsController[I, C, R any](prefix string, chainSet ChainSet[I, C], errNotEnabled error,
	parseChainID func(string) (I, error), newResource func(chains.Chain[I, C]) R) *chainsController[I, C, R] {
	return &chainsController[I, C, R]{
		resourceName:  prefix + "_chain",
		chainSet:      chainSet,
		errNotEnabled: errNotEnabled,
		parseChainID:  parseChainID,
		newResource:   newResource,
	}
}

func (cc *chainsController[I, C, R]) Index(c *gin.Context, size, page, offset int) {
	if cc.chainSet == nil {
		jsonAPIError(c, http.StatusBadRequest, cc.errNotEnabled)
		return
	}
	chains, count, err := cc.chainSet.Index(offset, size)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	var resources []interface{}
	for _, chain := range chains {
		resources = append(resources, cc.newResource(chain))
	}

	paginatedResponse(c, cc.resourceName, size, page, resources, count, err)
}

type CreateChainRequest[I, C any] struct {
	ID     I `json:"chainID"`
	Config C `json:"config"`
}

func NewCreateChainRequest[I, C any](id I, config C) CreateChainRequest[I, C] {
	return CreateChainRequest[I, C]{ID: id, Config: config}
}

func (cc *chainsController[I, C, R]) Create(c *gin.Context) {
	if cc.chainSet == nil {
		jsonAPIError(c, http.StatusBadRequest, cc.errNotEnabled)
		return
	}

	var request CreateChainRequest[I, C]
	//TODo does this already handle utils.Big....

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	chain, err := cc.chainSet.Add(c.Request.Context(), request.ID, request.Config)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponseWithStatus(c, cc.newResource(chain), cc.resourceName, http.StatusCreated)
}

func (cc *chainsController[I, C, R]) Show(c *gin.Context) {
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

type UpdateChainRequest[C any] struct {
	Enabled bool `json:"enabled"`
	Config  C    `json:"config"`
}

func (cc *chainsController[I, C, R]) Update(c *gin.Context) {
	if cc.chainSet == nil {
		jsonAPIError(c, http.StatusBadRequest, cc.errNotEnabled)
		return
	}

	id, err := cc.parseChainID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	var request UpdateChainRequest[C]
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	chain, err := cc.chainSet.Configure(c.Request.Context(), id, request.Enabled, request.Config)

	if errors.Is(err, sql.ErrNoRows) {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	} else if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	jsonAPIResponse(c, cc.newResource(chain), cc.resourceName)
}

func (cc *chainsController[I, C, R]) Delete(c *gin.Context) {
	if cc.chainSet == nil {
		jsonAPIError(c, http.StatusBadRequest, cc.errNotEnabled)
		return
	}

	id, err := cc.parseChainID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = cc.chainSet.Remove(id)

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, cc.resourceName, http.StatusNoContent)
}

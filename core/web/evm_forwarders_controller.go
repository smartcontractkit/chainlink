package web

import (
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"

	"github.com/gin-gonic/gin"
)

// EVMForwardersController manages EVM forwarders.
type EVMForwardersController struct {
	App chainlink.Application
}

// Index lists EVM forwarders.
func (cc *EVMForwardersController) Index(c *gin.Context, size, page, offset int) {
	orm := forwarders.NewORM(cc.App.GetDB())
	fwds, count, err := orm.FindForwarders(c.Request.Context(), 0, size)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	var resources []presenters.EVMForwarderResource
	for _, fwd := range fwds {
		resources = append(resources, presenters.NewEVMForwarderResource(fwd))
	}

	paginatedResponse(c, "forwarder", size, page, resources, count, err)
}

// TrackEVMForwarderRequest is a JSONAPI request for creating an EVM forwarder.
type TrackEVMForwarderRequest struct {
	EVMChainID *ubig.Big      `json:"evmChainId"`
	Address    common.Address `json:"address"`
}

// Track adds a new EVM forwarder.
func (cc *EVMForwardersController) Track(c *gin.Context) {
	request := &TrackEVMForwarderRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	orm := forwarders.NewORM(cc.App.GetDB())
	fwd, err := orm.CreateForwarder(c.Request.Context(), request.Address, *request.EVMChainID)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	cc.App.GetAuditLogger().Audit(audit.ForwarderCreated, map[string]interface{}{
		"forwarderID":         fwd.ID,
		"forwarderAddress":    fwd.Address,
		"forwarderEVMChainID": fwd.EVMChainID,
	})
	jsonAPIResponseWithStatus(c, presenters.NewEVMForwarderResource(fwd), "forwarder", http.StatusCreated)
}

// Delete removes an EVM Forwarder.
func (cc *EVMForwardersController) Delete(c *gin.Context) {
	id, err := stringutils.ToInt64(c.Param("fwdID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	filterCleanup := func(tx sqlutil.DataSource, evmChainID int64, addr common.Address) error {
		chain, err2 := cc.App.GetRelayers().LegacyEVMChains().Get(big.NewInt(evmChainID).String())
		if err2 != nil {
			// If the chain id doesn't even exist, or logpoller is disabled, then there isn't any filter to clean up.  Returning an error
			// here could be dangerous as it would make it impossible to delete a forwarder with an invalid chain id
			return nil
		}

		if chain.LogPoller() == logpoller.LogPollerDisabled {
			// handle same as non-existent chain id
			return nil
		}
		return chain.LogPoller().UnregisterFilter(c.Request.Context(), forwarders.FilterName(addr))
	}

	orm := forwarders.NewORM(cc.App.GetDB())
	err = orm.DeleteForwarder(c.Request.Context(), id, filterCleanup)

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	cc.App.GetAuditLogger().Audit(audit.ForwarderDeleted, map[string]interface{}{"id": id})
	jsonAPIResponseWithStatus(c, nil, "forwarder", http.StatusNoContent)
}

package web

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/core/logger/audit"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/gin-gonic/gin"
)

// EVMForwardersController manages EVM forwarders.
type EVMForwardersController struct {
	App chainlink.Application
}

// Index lists EVM forwarders.
func (cc *EVMForwardersController) Index(c *gin.Context, size, page, offset int) {
	orm := forwarders.NewORM(cc.App.GetSqlxDB(), cc.App.GetLogger(), cc.App.GetConfig())
	fwds, count, err := orm.FindForwarders(0, size)

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
	EVMChainID *utils.Big     `json:"chainID"`
	Address    common.Address `json:"address"`
}

// Track adds a new EVM forwarder.
func (cc *EVMForwardersController) Track(c *gin.Context) {
	request := &TrackEVMForwarderRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	orm := forwarders.NewORM(cc.App.GetSqlxDB(), cc.App.GetLogger(), cc.App.GetConfig())
	fwd, err := orm.CreateForwarder(request.Address, *request.EVMChainID)

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
	id, err := stringutils.ToInt32(c.Param("fwdID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	orm := forwarders.NewORM(cc.App.GetSqlxDB(), cc.App.GetLogger(), cc.App.GetConfig())
	err = orm.DeleteForwarder(id)

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	cc.App.GetAuditLogger().Audit(audit.ForwarderDeleted, map[string]interface{}{"id": id})
	jsonAPIResponseWithStatus(c, nil, "forwarder", http.StatusNoContent)
}

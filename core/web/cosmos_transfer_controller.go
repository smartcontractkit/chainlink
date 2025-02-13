package web

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"

	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	cosmosmodels "github.com/smartcontractkit/chainlink/v2/core/store/models/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// CosmosTransfersController can send LINK tokens to another address
type CosmosTransfersController struct {
	App chainlink.Application
}

var ErrCosmosNotEnabled = errChainDisabled{name: "Cosmos", tomlKey: "Cosmos.Enabled"}

// Create sends native coins from the Chainlink's account to a specified address.
func (tc *CosmosTransfersController) Create(c *gin.Context) {
	relayers := tc.App.GetRelayers().List(chainlink.FilterRelayersByType(relay.NetworkCosmos))
	if relayers == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrCosmosNotEnabled)
		return
	}

	var tr cosmosmodels.SendRequest
	if err := c.ShouldBindJSON(&tr); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	if tr.CosmosChainID == "" {
		jsonAPIError(c, http.StatusBadRequest, errors.New("missing cosmosChainID"))
		return
	}
	if tr.FromAddress == "" {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("withdrawal source address is missing: %v", tr.FromAddress))
		return
	}

	relayerID := types.RelayID{Network: relay.NetworkCosmos, ChainID: tr.CosmosChainID}
	relayer, err := relayers.Get(relayerID)
	if err != nil {
		if errors.Is(err, chainlink.ErrNoSuchRelayer) {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	cfgs := tc.App.GetConfig().CosmosConfigs()
	i := slices.IndexFunc(cfgs, func(config chainlink.RawConfig) bool { return config.ChainID() == tr.CosmosChainID })
	if i == -1 {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("no config for chain id: %s", tr.CosmosChainID))
		return
	}
	gasToken, ok := cfgs[i]["GasToken"].(string)
	if !ok {
		jsonAPIError(c, http.StatusBadRequest, errors.New("GasToken must be a string"))
		return
	}

	if gasToken != tr.Token {
		jsonAPIError(c, http.StatusBadRequest, errors.Errorf("unable to convert %s to %s: %v", tr.Token, gasToken, err))
		return
	} else if tr.Amount.Sign() != 1 {
		jsonAPIError(c, http.StatusBadRequest, errors.Errorf("amount must be greater than zero: %s", tr.Amount))
		return
	}

	err = relayer.Transact(c, tr.FromAddress, tr.DestinationAddress, tr.Amount, !tr.AllowHigherAmounts)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, errors.Errorf("failed to send transaction: %v", err))
		return
	}

	resource := presenters.NewCosmosMsgResource("cosmos_transfer_"+uuid.New().String(), tr.CosmosChainID, "")
	resource.State = "unstarted"
	tc.App.GetAuditLogger().Audit(audit.CosmosTransactionCreated, map[string]interface{}{
		"cosmosTransactionResource": resource,
	})

	jsonAPIResponse(c, resource, "cosmos_msg")
}

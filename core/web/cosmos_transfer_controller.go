package web

import (
	"fmt"
	"net/http"
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/denom"

	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	cosmosmodels "github.com/smartcontractkit/chainlink/v2/core/store/models/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// CosmosTransfersController can send LINK tokens to another address
type CosmosTransfersController struct {
	App chainlink.Application
}

// Create sends native coins from the Chainlink's account to a specified address.
func (tc *CosmosTransfersController) Create(c *gin.Context) {
	relayers := tc.App.GetRelayers().List(chainlink.FilterRelayersByType(types.NetworkCosmos))
	if relayers == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrSolanaNotEnabled)
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
	if tr.FromAddress.Empty() {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("withdrawal source address is missing: %v", tr.FromAddress))
		return
	}

	relayerID := types.RelayID{Network: types.NetworkCosmos, ChainID: tr.CosmosChainID}
	relayer, err := relayers.Get(relayerID)
	if err != nil {
		if errors.Is(err, chainlink.ErrNoSuchRelayer) {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	var gasToken string
	cfgs := tc.App.GetConfig().CosmosConfigs()
	i := slices.IndexFunc(cfgs, func(config *coscfg.TOMLConfig) bool { return *config.ChainID == tr.CosmosChainID })
	if i == -1 {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("no config for chain id: %s", tr.CosmosChainID))
		return
	}
	gasToken = cfgs[i].GasToken()

	//TODO move this inside?
	coin, err := denom.ConvertDecCoinToDenom(sdk.NewDecCoinFromDec(tr.Token, tr.Amount), gasToken)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, errors.Errorf("unable to convert %s to %s: %v", tr.Token, gasToken, err))
		return
	} else if !coin.Amount.IsPositive() {
		jsonAPIError(c, http.StatusBadRequest, errors.Errorf("amount must be greater than zero: %s", coin.Amount))
		return
	}

	err = relayer.Transact(c, tr.FromAddress.String(), tr.DestinationAddress.String(), coin.Amount.BigInt(), !tr.AllowHigherAmounts)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, errors.Errorf("failed to send transaction: %v", err))
		return
	}

	resource := presenters.NewCosmosMsgResource("cosmos_transfer_"+uuid.New().String(), tr.CosmosChainID, "")
	resource.State = string(db.Unstarted)
	tc.App.GetAuditLogger().Audit(audit.CosmosTransactionCreated, map[string]interface{}{
		"cosmosTransactionResource": resource,
	})

	jsonAPIResponse(c, resource, "cosmos_msg")
}

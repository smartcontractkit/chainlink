package web

import (
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"

	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos/denom"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	cosmosmodels "github.com/smartcontractkit/chainlink/v2/core/store/models/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// maxGasUsedTransfer is an upper bound on how much gas we expect a MsgSend for a single coin to use.
const maxGasUsedTransfer = 100_000

// CosmosTransfersController can send LINK tokens to another address
type CosmosTransfersController struct {
	App chainlink.Application
}

// Create sends Atom and other native coins from the Chainlink's account to a specified address.
func (tc *CosmosTransfersController) Create(c *gin.Context) {
	cosmosChains := tc.App.GetChains().Cosmos
	if cosmosChains == nil {
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
	chain, err := cosmosChains.Chain(c.Request.Context(), tr.CosmosChainID)
	if errors.Is(err, cosmos.ErrChainIDInvalid) || errors.Is(err, cosmos.ErrChainIDEmpty) {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	} else if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	if tr.FromAddress.Empty() {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("withdrawal source address is missing: %v", tr.FromAddress))
		return
	}

	coin, err := denom.DecCoinToUAtom(sdk.NewDecCoinFromDec("atom", tr.Amount))
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, errors.Errorf("unable to convert to uatom: %v", err))
		return
	} else if !coin.Amount.IsPositive() {
		jsonAPIError(c, http.StatusBadRequest, errors.Errorf("amount must be greater than zero: %s", coin.Amount))
		return
	}

	txm := chain.TxManager()

	if !tr.AllowHigherAmounts {
		var reader client.Reader
		reader, err = chain.Reader("")
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, errors.Errorf("chain unreachable: %v", err))
			return
		}
		gasPrice, err2 := txm.GasPrice()
		if err2 != nil {
			jsonAPIError(c, http.StatusInternalServerError, errors.Errorf("gas price unavailable: %v", err2))
			return
		}

		err = cosmosValidateBalance(reader, gasPrice, tr.FromAddress, coin)
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("failed to validate balance: %v", err))
			return
		}
	}

	sendMsg := bank.NewMsgSend(tr.FromAddress, tr.DestinationAddress, sdk.Coins{coin})
	msgID, err := txm.Enqueue("", sendMsg)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, errors.Errorf("transaction failed: %v", err))
		return
	}
	resource := presenters.NewCosmosMsgResource(msgID, tr.CosmosChainID, "")
	msgs, err := txm.GetMsgs(msgID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, errors.Errorf("failed to get message %d: %v", msgID, err))
		return
	}
	if len(msgs) != 1 {
		jsonAPIError(c, http.StatusInternalServerError, errors.Errorf("failed to get message %d: %v", msgID, err))
		return
	}
	msg := msgs[0]
	resource.TxHash = msg.TxHash
	resource.State = string(msg.State)

	tc.App.GetAuditLogger().Audit(audit.CosmosTransactionCreated, map[string]interface{}{
		"cosmosTransactionResource": resource,
	})

	jsonAPIResponse(c, resource, "cosmos_msg")
}

// cosmosValidateBalance validates that fromAddr's balance can cover coin, including fees at gasPrice.
// Note: This is currently limited to uatom only, for both gasPrice and coin.
func cosmosValidateBalance(reader client.Reader, gasPrice sdk.DecCoin, fromAddr sdk.AccAddress, coin sdk.Coin) error {
	const denom = "uatom"
	if gasPrice.Denom != denom {
		return errors.Errorf("unsupported gas price denom: %s", gasPrice.Denom)
	}
	if coin.Denom != denom {
		return errors.Errorf("unsupported coin denom: %s", gasPrice.Denom)
	}

	balance, err := reader.Balance(fromAddr, denom)
	if err != nil {
		return err
	}

	fee := gasPrice.Amount.MulInt64(maxGasUsedTransfer).RoundInt()
	need := coin.Amount.Add(fee)

	if balance.Amount.LT(need) {
		return errors.Errorf("balance %q is too low for this transaction to be executed: need %s total, including %s fee", balance, need, fee)
	}
	return nil
}

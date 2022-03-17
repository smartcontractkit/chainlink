package web

import (
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra/client"

	"github.com/smartcontractkit/chainlink/core/chains/terra"
	"github.com/smartcontractkit/chainlink/core/chains/terra/denom"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	terramodels "github.com/smartcontractkit/chainlink/core/store/models/terra"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// maxGasUsedTransfer is an upper bound on how much gas we expect a MsgSend for a single coin to use.
const maxGasUsedTransfer = 100_000

// TerraTransfersController can send LINK tokens to another address
type TerraTransfersController struct {
	App chainlink.Application
}

// Create sends Luna and other native coins from the Chainlink's account to a specified address.
func (tc *TerraTransfersController) Create(c *gin.Context) {
	terraChains := tc.App.GetChains().Terra
	if terraChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrTerraNotEnabled)
		return
	}

	var tr terramodels.SendRequest
	if err := c.ShouldBindJSON(&tr); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	if tr.TerraChainID == "" {
		jsonAPIError(c, http.StatusBadRequest, errors.New("missing terraChainID"))
		return
	}
	chain, err := terraChains.Chain(c.Request.Context(), tr.TerraChainID)
	switch err {
	case terra.ErrChainIDInvalid, terra.ErrChainIDEmpty:
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	case nil:
		break
	default:
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	if tr.FromAddress.Empty() {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("withdrawal source address is missing: %v", tr.FromAddress))
		return
	}

	coin, err := denom.ConvertToULuna(sdk.NewDecCoinFromDec("luna", tr.Amount))
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, errors.Errorf("unable to convert to uluna: %v", err))
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
		gasPrice, err := txm.GasPrice()
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, errors.Errorf("gas price unavailable: %v", err))
			return
		}

		err = terraValidateBalance(reader, gasPrice, tr.FromAddress, coin)
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
	resource := presenters.NewTerraMsgResource(msgID, tr.TerraChainID, "")
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

	jsonAPIResponse(c, resource, "terra_msg")
}

// terraValidateBalance validates that fromAddr's balance can cover coin, including fees at gasPrice.
// Note: This is currently limited to uluna only, for both gasPrice and coin.
func terraValidateBalance(reader client.Reader, gasPrice sdk.DecCoin, fromAddr sdk.AccAddress, coin sdk.Coin) error {
	const denom = "uluna"
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

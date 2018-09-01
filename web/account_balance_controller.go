package web

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
)

// AccountBalanceController returns information for the active account
type AccountBalanceController struct {
	App *services.ChainlinkApplication
}

// Show returns the address, plus it's ETH & LINK balance
// Example:
//  "<application>/account_balance"
func (abc *AccountBalanceController) Show(c *gin.Context) {
	store := abc.App.Store
	txm := store.TxManager

	if account, err := store.KeyStore.GetAccount(); err != nil {
		publicError(c, 400, err)
	} else if ethBalance, err := txm.GetEthBalance(account.Address); err != nil {
		c.AbortWithError(500, err)
	} else if linkBalance, err := txm.GetLinkBalance(account.Address); err != nil {
		c.AbortWithError(500, err)
	} else {
		ab := presenters.AccountBalance{
			Address:     account.Address.Hex(),
			EthBalance:  ethBalance,
			LinkBalance: linkBalance,
		}
		if json, err := jsonapi.Marshal(ab); err != nil {
			c.AbortWithError(500, fmt.Errorf("failed to marshal account using jsonapi: %+v", err))
		} else {
			c.Data(200, MediaType, json)
		}
	}
}

var naz = assets.NewLink(1)

// Withdraw sends LINK from the configured oracle contract to the given address
func (abc *AccountBalanceController) Withdraw(c *gin.Context) {
	store := abc.App.Store
	txm := store.TxManager
	wr := models.WithdrawalRequest{}
	oca := *store.Config.OracleContractAddress

	if err := c.ShouldBindJSON(&wr); err != nil {
		publicError(c, 400, err)
	} else if wr.Amount.Cmp(naz) < 0 {
		publicError(c, 400, fmt.Errorf("Must withdraw at least %v LINK", naz.String()))
	} else if wr.Address == utils.ZeroAddress { // address is unmarshalled to ZeroAddres if invalid
		publicError(c, 400, errors.New("Invalid withdrawal address"))
	} else if account, err := store.KeyStore.GetAccount(); err != nil {
		c.AbortWithError(500, err)
	} else if linkBalance, err := txm.GetLinkBalance(account.Address); err != nil {
		c.AbortWithError(500, err)
	} else if linkBalance.Cmp(wr.Amount) < 0 {
		publicError(c, 400, fmt.Errorf("Insufficient link balance. Withdrawal Amount: %v Link Balance: %v", wr.Amount.String(), linkBalance.String()))
	} else if tx, err := txm.CreateTx(oca, nil); err != nil {
		c.AbortWithError(500, err)
		// functionSelector := models.HexToFunctionSelector("0x70a08231") // balanceOf(address)
		// data, err := utils.HexToBytes(functionSelector.String(), common.ToHex(common.LeftPadBytes(address.Bytes(), utils.EVMWordByteLen)))
	} else {
		c.JSON(200, tx.Hash)
	}
}

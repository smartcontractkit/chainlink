package web

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/presenters"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// KeysController manages account keys

// KeysController manages account keys
type ETHKeysController struct {
	App chainlink.Application
}

// Index returns the node's Ethereum keys and the account balances of ETH & LINK.
// Example:
//  "<application>/keys/eth"
func (ekc *ETHKeysController) Index(c *gin.Context) {
	store := ekc.App.GetStore()
	keys, err := store.AllKeys()
	if err != nil {
		err = errors.Errorf("error fetching ETH keys from database: %v", err)
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	var pkeys []presenters.ETHKey
	for _, key := range keys {
		ethBalance, err := store.EthClient.BalanceAt(c.Request.Context(), key.Address.Address(), nil)
		if err != nil {
			err = errors.Errorf("error calling getEthBalance on Ethereum node: %v", err)
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}

		linkAddress := common.HexToAddress(store.Config.LinkContractAddress())
		linkBalance, err := store.EthClient.GetLINKBalance(linkAddress, key.Address.Address())
		if err != nil {
			err = errors.Errorf("error calling getLINKBalance on Ethereum node: %v", err)
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}

		key, err := store.ORM.KeyByAddress(key.Address.Address())
		if err != nil {
			err = errors.Errorf("error fetching ETH key from DB: %v", err)
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
		pkeys = append(pkeys, presenters.ETHKey{
			Address:     key.Address.Hex(),
			EthBalance:  (*assets.Eth)(ethBalance),
			LinkBalance: linkBalance,
			NextNonce:   key.NextNonce,
			LastUsed:    key.LastUsed,
			IsFunding:   key.IsFunding,
			CreatedAt:   key.CreatedAt,
			UpdatedAt:   key.UpdatedAt,
			DeletedAt:   key.DeletedAt,
		})
	}
	jsonAPIResponse(c, pkeys, "keys")
}

// Create adds a new account
// Example:
//  "<application>/keys/eth"
func (ekc *ETHKeysController) Create(c *gin.Context) {
	store := ekc.App.GetStore()
	account, err := store.KeyStore.NewAccount()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	if err := store.SyncDiskKeyStoreToDB(); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	key, err := store.KeyByAddress(account.Address)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	ethBalance, err := store.EthClient.BalanceAt(c.Request.Context(), account.Address, nil)
	if err != nil {
		logger.Errorf("error calling getEthBalance on Ethereum node: %v", err)
	}
	linkAddress := common.HexToAddress(store.Config.LinkContractAddress())
	linkBalance, err := store.EthClient.GetLINKBalance(linkAddress, account.Address)
	if err != nil {
		logger.Errorf("error calling getLINKBalance on Ethereum node: %v", err)
	}

	pek := presenters.ETHKey{
		Address:     account.Address.Hex(),
		EthBalance:  (*assets.Eth)(ethBalance),
		LinkBalance: linkBalance,
		NextNonce:   key.NextNonce,
		LastUsed:    key.LastUsed,
		IsFunding:   key.IsFunding,
		CreatedAt:   key.CreatedAt,
		UpdatedAt:   key.UpdatedAt,
		DeletedAt:   key.DeletedAt,
	}
	jsonAPIResponseWithStatus(c, pek, "account", http.StatusCreated)
}

// Delete an ETH key bundle
// Example:
// "DELETE <application>/keys/eth/:keyID"
// "DELETE <application>/keys/eth/:keyID?hard=true"
func (ekc *ETHKeysController) Delete(c *gin.Context) {
	var hardDelete bool
	var err error
	if c.Query("hard") != "" {
		hardDelete, err = strconv.ParseBool(c.Query("hard"))
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, err)
			return
		}
	}

	if !common.IsHexAddress(c.Param("keyID")) {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.New("invalid address"))
		return
	}
	address := common.HexToAddress(c.Param("keyID"))
	if exists, err := ekc.App.GetStore().KeyExists(address); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	} else if !exists {
		jsonAPIError(c, http.StatusNotFound, errors.New("Key does not exist"))
		return
	}

	key, err := ekc.App.GetStore().KeyByAddress(address)
	store := ekc.App.GetStore()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	if hardDelete {
		err = store.DeleteKey(address)
	} else {
		err = store.ArchiveKey(address)
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ethBalance, err := store.EthClient.BalanceAt(c.Request.Context(), address, nil)
	if err != nil {
		logger.Errorf("error calling getEthBalance on Ethereum node: %v", err)
	}
	linkAddress := common.HexToAddress(store.Config.LinkContractAddress())
	linkBalance, err := store.EthClient.GetLINKBalance(linkAddress, address)
	if err != nil {
		logger.Errorf("error calling getLINKBalance on Ethereum node: %v", err)
	}

	pek := presenters.ETHKey{
		Address:     address.Hex(),
		EthBalance:  (*assets.Eth)(ethBalance),
		LinkBalance: linkBalance,
		NextNonce:   key.NextNonce,
		LastUsed:    key.LastUsed,
		IsFunding:   key.IsFunding,
		CreatedAt:   key.CreatedAt,
		UpdatedAt:   key.UpdatedAt,
		DeletedAt:   key.DeletedAt,
	}
	jsonAPIResponse(c, pek, "account")
}

func (ekc *ETHKeysController) Import(c *gin.Context) {
	defer logger.ErrorIfCalling(c.Request.Body.Close)

	store := ekc.App.GetStore()

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")

	acct, err := store.KeyStore.Import(bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	err = store.SyncDiskKeyStoreToDB()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ethBalance, err := store.EthClient.BalanceAt(c.Request.Context(), acct.Address, nil)
	if err != nil {
		err = errors.Errorf("error calling getEthBalance on Ethereum node: %v", err)
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	linkAddress := common.HexToAddress(store.Config.LinkContractAddress())
	linkBalance, err := store.EthClient.GetLINKBalance(linkAddress, acct.Address)
	if err != nil {
		err = errors.Errorf("error calling getLINKBalance on Ethereum node: %v", err)
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	key, err := store.ORM.KeyByAddress(acct.Address)
	if err != nil {
		err = errors.Errorf("error fetching ETH key from DB: %v", err)
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	pek := presenters.ETHKey{
		Address:     key.Address.Hex(),
		EthBalance:  (*assets.Eth)(ethBalance),
		LinkBalance: linkBalance,
		NextNonce:   key.NextNonce,
		LastUsed:    key.LastUsed,
		IsFunding:   key.IsFunding,
		CreatedAt:   key.CreatedAt,
		UpdatedAt:   key.UpdatedAt,
		DeletedAt:   key.DeletedAt,
	}
	jsonAPIResponse(c, pek, "account")
}

func (ekc *ETHKeysController) Export(c *gin.Context) {
	defer logger.ErrorIfCalling(c.Request.Body.Close)

	addressStr := c.Param("address")
	address := common.HexToAddress(addressStr)
	newPassword := c.Query("newpassword")

	bytes, err := ekc.App.GetStore().KeyStore.Export(address, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, MediaType, bytes)
}

package web

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// ETHKeysController manages account keys
type ETHKeysController struct {
	App chainlink.Application
}

// Index returns the node's Ethereum keys and the account balances of ETH & LINK.
// Example:
//  "<application>/keys/eth"
func (ekc *ETHKeysController) Index(c *gin.Context) {
	ethKeyStore := ekc.App.GetKeyStore().Eth()
	var keys []ethkey.KeyV2
	var err error
	if ekc.App.GetStore().Config.Dev() {
		keys, err = ethKeyStore.GetAll()
	} else {
		keys, err = ethKeyStore.SendingKeys()
	}
	if err != nil {
		err = errors.Errorf("error getting unlocked keys: %v", err)
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	states, err := ethKeyStore.GetStatesForKeys(keys)
	if err != nil {
		err = errors.Errorf("error getting key states: %v", err)
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	var resources []presenters.ETHKeyResource
	for _, state := range states {
		key, err := ethKeyStore.Get(state.Address.Hex())
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
		r, err := presenters.NewETHKeyResource(key, state,
			ekc.setEthBalance(c.Request.Context(), key.Address.Address()),
			ekc.setLinkBalance(key.Address.Address()),
		)
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}

		resources = append(resources, *r)
	}

	jsonAPIResponse(c, resources, "keys")
}

// Create adds a new account
// Example:
//  "<application>/keys/eth"
func (ekc *ETHKeysController) Create(c *gin.Context) {
	ethKeyStore := ekc.App.GetKeyStore().Eth()
	key, err := ethKeyStore.Create()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	state, err := ethKeyStore.GetState(key.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	r, err := presenters.NewETHKeyResource(key, state,
		ekc.setEthBalance(c.Request.Context(), key.Address.Address()),
		ekc.setLinkBalance(key.Address.Address()),
	)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, r, "account", http.StatusCreated)
}

// Delete an ETH key bundle
// Example:
// "DELETE <application>/keys/eth/:keyID"
// "DELETE <application>/keys/eth/:keyID?hard=true"
func (ekc *ETHKeysController) Delete(c *gin.Context) {
	ethKeyStore := ekc.App.GetKeyStore().Eth()
	var hardDelete bool
	var err error

	if c.Query("hard") != "" {
		hardDelete, err = strconv.ParseBool(c.Query("hard"))
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, err)
			return
		}
	}

	if !hardDelete {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.New("hard delete only"))
		return
	}

	if !common.IsHexAddress(c.Param("keyID")) {
		jsonAPIError(c, http.StatusBadRequest, errors.New("hard delete only"))
		return
	}
	keyID := c.Param("keyID")
	state, err := ethKeyStore.GetState(keyID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	key, err := ethKeyStore.Delete(keyID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	r, err := presenters.NewETHKeyResource(key, state,
		ekc.setEthBalance(c.Request.Context(), key.Address.Address()),
		ekc.setLinkBalance(key.Address.Address()),
	)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, r, "account")
}

// Import imports a key
func (ekc *ETHKeysController) Import(c *gin.Context) {
	ethKeyStore := ekc.App.GetKeyStore().Eth()
	defer logger.ErrorIfCalling(c.Request.Body.Close)

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")

	key, err := ethKeyStore.Import(bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	state, err := ethKeyStore.GetState(key.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	r, err := presenters.NewETHKeyResource(key, state,
		ekc.setEthBalance(c.Request.Context(), key.Address.Address()),
		ekc.setLinkBalance(key.Address.Address()),
	)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, r, "account")
}

func (ekc *ETHKeysController) Export(c *gin.Context) {
	defer logger.ErrorIfCalling(c.Request.Body.Close)

	address := c.Param("address")
	newPassword := c.Query("newpassword")

	bytes, err := ekc.App.GetKeyStore().Eth().Export(address, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, MediaType, bytes)
}

// setEthBalance is a custom functional option for NewEthKeyResource which
// queries the EthClient for the ETH balance at the address and sets it on the
// resource.
func (ekc *ETHKeysController) setEthBalance(ctx context.Context, accountAddr common.Address) presenters.NewETHKeyOption {
	ethClient := ekc.App.GetEthClient()
	bal, err := ethClient.BalanceAt(ctx, accountAddr, nil)

	return func(r *presenters.ETHKeyResource) error {
		if err != nil {
			return errors.Errorf("error calling getEthBalance on Ethereum node: %v", err)
		}

		r.EthBalance = (*assets.Eth)(bal)

		return nil
	}
}

// setLinkBalance is a custom functional option for NewEthKeyResource which
// queries the EthClient for the LINK balance at the address and sets it on the
// resource.
func (ekc *ETHKeysController) setLinkBalance(accountAddr common.Address) presenters.NewETHKeyOption {
	ethClient := ekc.App.GetEthClient()
	addr := common.HexToAddress(ekc.App.GetEVMConfig().LinkContractAddress())
	bal, err := ethClient.GetLINKBalance(addr, accountAddr)

	return func(r *presenters.ETHKeyResource) error {
		if err != nil {
			return errors.Errorf("error calling getLINKBalance on Ethereum node: %v", err)
		}

		r.LinkBalance = bal

		return nil
	}
}

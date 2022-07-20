package web

import (
	"context"
	"io/ioutil"
	"math/big"
	"net/http"
	"sort"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
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
	if ekc.App.GetConfig().Dev() {
		keys, err = ethKeyStore.GetAll()
	} else {
		keys, err = ethKeyStore.SendingKeys(nil)
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
			ekc.setEthBalance(c.Request.Context(), state),
			ekc.setLinkBalance(state),
			ekc.setKeyMaxGasPriceWei(state, key.Address.Address()),
		)
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}

		resources = append(resources, *r)
	}
	// Put funding keys to the end
	sort.SliceStable(resources, func(i, j int) bool {
		return !resources[i].IsFunding && resources[j].IsFunding
	})

	jsonAPIResponse(c, resources, "keys")
}

// Create adds a new account
// Example:
//  "<application>/keys/eth"
func (ekc *ETHKeysController) Create(c *gin.Context) {
	ethKeyStore := ekc.App.GetKeyStore().Eth()

	chain, err := getChain(ekc.App.GetChains().EVM, c.Query("evmChainID"))
	switch err {
	case ErrInvalidChainID, ErrMultipleChains, ErrMissingChainID:
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	case nil:
		break
	default:
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	var maxGasPriceGWei int64
	if c.Query("maxGasPriceGWei") != "" {
		maxGasPriceGWei, err = strconv.ParseInt(c.Query("maxGasPriceGWei"), 10, 64)
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, err)
			return
		}
	}

	key, err := ethKeyStore.Create(chain.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	if maxGasPriceGWei > 0 {
		maxGasPriceWei := assets.GWei(maxGasPriceGWei)
		updateMaxGasPrice := evm.UpdateKeySpecificMaxGasPrice(key.Address.Address(), maxGasPriceWei)
		if err = ekc.App.GetChains().EVM.UpdateConfig(chain.ID(), updateMaxGasPrice); err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
	}

	state, err := ethKeyStore.GetState(key.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	r, err := presenters.NewETHKeyResource(key, state,
		ekc.setEthBalance(c.Request.Context(), state),
		ekc.setLinkBalance(state),
		ekc.setKeyMaxGasPriceWei(state, key.Address.Address()),
	)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, r, "account", http.StatusCreated)
}

// Update an ETH key's parameters
// Example:
// "PUT <application>/keys/eth/:keyID?maxGasPriceGWei=12345"
func (ekc *ETHKeysController) Update(c *gin.Context) {
	ethKeyStore := ekc.App.GetKeyStore().Eth()

	if c.Query("maxGasPriceGWei") == "" {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.New("no parameters passed to update"))
		return
	}

	maxGasPriceGWei, err := strconv.ParseInt(c.Query("maxGasPriceGWei"), 10, 64)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	keyID := c.Param("keyID")
	state, err := ethKeyStore.GetState(keyID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	key, err := ethKeyStore.Get(keyID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	maxGasPriceWei := assets.GWei(maxGasPriceGWei)
	updateMaxGasPrice := evm.UpdateKeySpecificMaxGasPrice(key.Address.Address(), maxGasPriceWei)
	if err = ekc.App.GetChains().EVM.UpdateConfig((*big.Int)(&state.EVMChainID), updateMaxGasPrice); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	r, err := presenters.NewETHKeyResource(key, state,
		ekc.setEthBalance(c.Request.Context(), state),
		ekc.setLinkBalance(state),
		ekc.setKeyMaxGasPriceWei(state, key.Address.Address()),
	)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, r, "account", http.StatusOK)
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
		ekc.setEthBalance(c.Request.Context(), state),
		ekc.setLinkBalance(state),
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
	defer ekc.App.GetLogger().ErrorIfClosing(c.Request.Body, "Import request body")

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	chain, err := getChain(ekc.App.GetChains().EVM, c.Query("evmChainID"))
	switch err {
	case ErrInvalidChainID, ErrMultipleChains, ErrMissingChainID:
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	case nil:
		break
	default:
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	key, err := ethKeyStore.Import(bytes, oldPassword, chain.ID())
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
		ekc.setEthBalance(c.Request.Context(), state),
		ekc.setLinkBalance(state),
	)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, r, "account")
}

func (ekc *ETHKeysController) Export(c *gin.Context) {
	defer ekc.App.GetLogger().ErrorIfClosing(c.Request.Body, "Export request body")

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
func (ekc *ETHKeysController) setEthBalance(ctx context.Context, state ethkey.State) presenters.NewETHKeyOption {
	var bal *big.Int
	chain, err := ekc.App.GetChains().EVM.Get(state.EVMChainID.ToInt())
	if err == nil {
		ethClient := chain.Client()
		bal, err = ethClient.BalanceAt(ctx, state.Address.Address(), nil)
	}
	return func(r *presenters.ETHKeyResource) error {
		if errors.Is(errors.Cause(err), evm.ErrNoChains) {
			return nil
		}

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
func (ekc *ETHKeysController) setLinkBalance(state ethkey.State) presenters.NewETHKeyOption {
	var bal *assets.Link
	chain, err := ekc.App.GetChains().EVM.Get(state.EVMChainID.ToInt())
	if err == nil {
		ethClient := chain.Client()
		addr := common.HexToAddress(chain.Config().LinkContractAddress())
		bal, err = ethClient.GetLINKBalance(addr, state.Address.Address())
	}

	return func(r *presenters.ETHKeyResource) error {
		if errors.Is(errors.Cause(err), evm.ErrNoChains) {
			return nil
		}
		if err != nil {
			return errors.Errorf("error calling getLINKBalance on Ethereum node: %v", err)
		}

		r.LinkBalance = bal

		return nil
	}
}

// setKeyMaxGasPriceWei is a custom functional option for NewEthKeyResource which
// gets the key specific max gas price from the chain config and sets it on the
// resource.
func (ekc *ETHKeysController) setKeyMaxGasPriceWei(state ethkey.State, keyAddress common.Address) presenters.NewETHKeyOption {
	var price *big.Int
	chain, err := ekc.App.GetChains().EVM.Get(state.EVMChainID.ToInt())
	if err == nil {
		price = chain.Config().KeySpecificMaxGasPriceWei(keyAddress)
	}

	return func(r *presenters.ETHKeyResource) error {
		if errors.Is(errors.Cause(err), evm.ErrNoChains) {
			return nil
		}
		if err != nil {
			return errors.Errorf("error getting EVM Chain: %v", err)
		}

		r.MaxGasPriceWei = *utils.NewBig(price)

		return nil
	}
}

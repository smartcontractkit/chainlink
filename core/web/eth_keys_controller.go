package web

import (
	"context"
	"io"
	"math/big"
	"net/http"
	"sort"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/logger/audit"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

// ETHKeysController manages account keys
type ETHKeysController struct {
	app  chainlink.Application
	lggr logger.Logger
}

func NewETHKeysController(app chainlink.Application) *ETHKeysController {
	return &ETHKeysController{
		app:  app,
		lggr: app.GetLogger().Named("ETHKeysController"),
	}
}

// Index returns the node's Ethereum keys and the account balances of ETH & LINK.
// Example:
//
//	"<application>/keys/eth"
func (ekc *ETHKeysController) Index(c *gin.Context) {
	ethKeyStore := ekc.app.GetKeyStore().Eth()
	var keys []ethkey.KeyV2
	var err error
	keys, err = ethKeyStore.GetAll()
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

		r := presenters.NewETHKeyResource(key, state,
			ekc.setEthBalance(c.Request.Context(), state),
			ekc.setLinkBalance(c.Request.Context(), state),
			ekc.setKeyMaxGasPriceWei(state, key.Address),
		)

		resources = append(resources, *r)
	}
	// Put disabled keys to the end
	sort.SliceStable(resources, func(i, j int) bool {
		return !resources[i].Disabled && resources[j].Disabled
	})

	jsonAPIResponse(c, resources, "keys")
}

// Create adds a new account
// Example:
//
//	"<application>/keys/eth"
func (ekc *ETHKeysController) Create(c *gin.Context) {
	ethKeyStore := ekc.app.GetKeyStore().Eth()

	chain, err := getChain(ekc.app.GetChains().EVM, c.Query("evmChainID"))
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
		updateMaxGasPrice := evm.UpdateKeySpecificMaxGasPrice(key.Address, maxGasPriceWei)
		if err = ekc.app.GetChains().EVM.UpdateConfig(chain.ID(), updateMaxGasPrice); err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
	}

	state, err := ethKeyStore.GetState(key.ID(), chain.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	r := presenters.NewETHKeyResource(key, state,
		ekc.setEthBalance(c.Request.Context(), state),
		ekc.setLinkBalance(c.Request.Context(), state),
		ekc.setKeyMaxGasPriceWei(state, key.Address),
	)

	ekc.app.GetAuditLogger().Audit(audit.KeyCreated, map[string]interface{}{
		"type": "ethereum",
		"id":   key.ID(),
	})

	jsonAPIResponseWithStatus(c, r, "account", http.StatusCreated)
}

// Update an ETH key's parameters
// Example:
// "PUT <application>/keys/eth/:keyID?maxGasPriceGWei=12345"
func (ekc *ETHKeysController) Update(c *gin.Context) {
	ethKeyStore := ekc.app.GetKeyStore().Eth()

	if c.Query("maxGasPriceGWei") == "" {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.New("no parameters passed to update"))
		return
	}

	maxGasPriceGWei, err := strconv.ParseInt(c.Query("maxGasPriceGWei"), 10, 64)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	chain, err := getChain(ekc.app.GetChains().EVM, c.Query("evmChainID"))
	if errors.Is(err, ErrInvalidChainID) || errors.Is(err, ErrMultipleChains) || errors.Is(err, ErrMissingChainID) {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	} else if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	keyID := c.Param("keyID")
	state, err := ethKeyStore.GetState(keyID, chain.ID())
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
	updateMaxGasPrice := evm.UpdateKeySpecificMaxGasPrice(key.Address, maxGasPriceWei)
	if err = ekc.app.GetChains().EVM.UpdateConfig((*big.Int)(&state.EVMChainID), updateMaxGasPrice); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	r := presenters.NewETHKeyResource(key, state,
		ekc.setEthBalance(c.Request.Context(), state),
		ekc.setLinkBalance(c.Request.Context(), state),
		ekc.setKeyMaxGasPriceWei(state, key.Address),
	)

	ekc.app.GetAuditLogger().Audit(audit.KeyUpdated, map[string]interface{}{
		"type": "ethereum",
		"id":   keyID,
	})

	jsonAPIResponseWithStatus(c, r, "account", http.StatusOK)
}

// Delete an ETH key bundle (irreversible!)
// Example:
// "DELETE <application>/keys/eth/:keyID"
func (ekc *ETHKeysController) Delete(c *gin.Context) {
	ethKeyStore := ekc.app.GetKeyStore().Eth()
	var err error

	keyID := c.Param("keyID")
	if !common.IsHexAddress(keyID) {
		jsonAPIError(c, http.StatusBadRequest, errors.Errorf("invalid keyID: %s, must be hex address", keyID))
		return
	}

	_, err = ethKeyStore.Delete(keyID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ekc.app.GetAuditLogger().Audit(audit.KeyDeleted, map[string]interface{}{
		"type": "ethereum",
		"id":   keyID,
	})
	c.Status(http.StatusNoContent)
}

// Import imports a key
func (ekc *ETHKeysController) Import(c *gin.Context) {
	ethKeyStore := ekc.app.GetKeyStore().Eth()
	defer ekc.app.GetLogger().ErrorIfFn(c.Request.Body.Close, "Error closing Import request body")

	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	chain, err := getChain(ekc.app.GetChains().EVM, c.Query("evmChainID"))
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

	state, err := ethKeyStore.GetState(key.ID(), chain.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	r := presenters.NewETHKeyResource(key, state,
		ekc.setEthBalance(c.Request.Context(), state),
		ekc.setLinkBalance(c.Request.Context(), state),
	)

	ekc.app.GetAuditLogger().Audit(audit.KeyImported, map[string]interface{}{
		"type": "ethereum",
		"id":   key.ID(),
	})

	jsonAPIResponse(c, r, "account")
}

func (ekc *ETHKeysController) Export(c *gin.Context) {
	defer ekc.app.GetLogger().ErrorIfFn(c.Request.Body.Close, "Error closing Export request body")

	id := c.Param("address")
	newPassword := c.Query("newpassword")

	bytes, err := ekc.app.GetKeyStore().Eth().Export(id, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ekc.app.GetAuditLogger().Audit(audit.KeyExported, map[string]interface{}{
		"type": "ethereum",
		"id":   id,
	})

	c.Data(http.StatusOK, MediaType, bytes)
}

// Chain updates settings for a given chain for the key
func (ekc *ETHKeysController) Chain(c *gin.Context) {
	kst := ekc.app.GetKeyStore().Eth()
	defer ekc.app.GetLogger().ErrorIfFn(c.Request.Body.Close, "Error closing Import request body")

	addressHex := c.Query("address")
	addressBytes, err := hexutil.Decode(addressHex)
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Wrap(err, "invalid address"))
		return
	}
	address := common.BytesToAddress(addressBytes)

	cid := c.Query("evmChainID")
	chain, err := getChain(ekc.app.GetChains().EVM, cid)
	if errors.Is(err, ErrInvalidChainID) || errors.Is(err, ErrMultipleChains) || errors.Is(err, ErrMissingChainID) {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	} else if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	var nonce int64 = -1
	if nonceStr := c.Query("nextNonce"); nonceStr != "" {
		nonce, err = strconv.ParseInt(nonceStr, 10, 64)
		if err != nil || nonce < 0 {
			jsonAPIError(c, http.StatusUnprocessableEntity, errors.Wrapf(err, "invalid value for nonce: expected 0 or positive int, got: %s", nonceStr))
			return
		}
	}
	abandon := false
	if abandonStr := c.Query("abandon"); abandonStr != "" {
		abandon, err = strconv.ParseBool(abandonStr)
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, errors.Wrapf(err, "invalid value for abandon: expected boolean, got: %s", abandonStr))
			return
		}
	}

	// Reset the chain
	if abandon || nonce >= 0 {
		var resetErr error
		err = chain.TxManager().Reset(func() {
			if nonce >= 0 {
				resetErr = kst.Reset(address, chain.ID(), nonce)
			}
		}, address, abandon)
		err = multierr.Combine(err, resetErr)
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
	}

	enabledStr := c.Query("enabled")
	if enabledStr != "" {
		var enabled bool
		enabled, err = strconv.ParseBool(enabledStr)
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, errors.Wrap(err, "enabled must be bool"))
			return
		}

		if enabled {
			err = kst.Enable(address, chain.ID())
		} else {
			err = kst.Disable(address, chain.ID())
		}
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
	}

	key, err := kst.Get(address.Hex())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	state, err := kst.GetState(key.ID(), chain.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	r := presenters.NewETHKeyResource(key, state,
		ekc.setEthBalance(c.Request.Context(), state),
		ekc.setLinkBalance(c.Request.Context(), state),
	)

	jsonAPIResponse(c, r, "account")
}

// setEthBalance is a custom functional option for NewEthKeyResource which
// queries the EthClient for the ETH balance at the address and sets it on the
// resource.
func (ekc *ETHKeysController) setEthBalance(ctx context.Context, state ethkey.State) presenters.NewETHKeyOption {
	var bal *big.Int
	chainID := state.EVMChainID.ToInt()
	chain, err := ekc.app.GetChains().EVM.Get(chainID)
	if err != nil {
		if !errors.Is(errors.Cause(err), evm.ErrNoChains) {
			ekc.lggr.Errorw("Failed to get EVM Chain", "chainID", chainID, "address", state.Address, "error", err)
		}
	} else {
		ethClient := chain.Client()
		bal, err = ethClient.BalanceAt(ctx, state.Address.Address(), nil)
		if err != nil {
			ekc.lggr.Errorw("Failed to get ETH balance", "chainID", chainID, "address", state.Address, "error", err)
		}
	}
	return presenters.SetETHKeyEthBalance((*assets.Eth)(bal))
}

// setLinkBalance is a custom functional option for NewEthKeyResource which
// queries the EthClient for the LINK balance at the address and sets it on the
// resource.
func (ekc *ETHKeysController) setLinkBalance(ctx context.Context, state ethkey.State) presenters.NewETHKeyOption {
	var bal *assets.Link
	chainID := state.EVMChainID.ToInt()
	chain, err := ekc.app.GetChains().EVM.Get(chainID)
	if err != nil {
		if !errors.Is(errors.Cause(err), evm.ErrNoChains) {
			ekc.lggr.Errorw("Failed to get EVM Chain", "chainID", chainID, "error", err)
		}
	} else {
		ethClient := chain.Client()
		addr := common.HexToAddress(chain.Config().LinkContractAddress())
		bal, err = ethClient.GetLINKBalance(ctx, addr, state.Address.Address())
		if err != nil {
			ekc.lggr.Errorw("Failed to get LINK balance", "chainID", chainID, "address", state.Address, "error", err)
		}
	}
	return presenters.SetETHKeyLinkBalance(bal)
}

// setKeyMaxGasPriceWei is a custom functional option for NewEthKeyResource which
// gets the key specific max gas price from the chain config and sets it on the
// resource.
func (ekc *ETHKeysController) setKeyMaxGasPriceWei(state ethkey.State, keyAddress common.Address) presenters.NewETHKeyOption {
	var price *assets.Wei
	chainID := state.EVMChainID.ToInt()
	chain, err := ekc.app.GetChains().EVM.Get(chainID)
	if err != nil {
		if !errors.Is(errors.Cause(err), evm.ErrNoChains) {
			ekc.lggr.Errorw("Failed to get EVM Chain", "chainID", chainID, "error", err)
		}
	} else {
		price = chain.Config().KeySpecificMaxGasPriceWei(keyAddress)
	}
	return presenters.SetETHKeyMaxGasPriceWei(utils.NewBig(price.ToInt()))
}

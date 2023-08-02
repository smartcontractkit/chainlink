package web

import (
	"context"
	"io"
	"math/big"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"

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

func createETHKeyResource(c *gin.Context, ekc *ETHKeysController, key ethkey.KeyV2, state ethkey.State) *presenters.ETHKeyResource {
	r := presenters.NewETHKeyResource(key, state,
		ekc.setEthBalance(c.Request.Context(), state),
		ekc.setLinkBalance(c.Request.Context(), state),
		ekc.setKeyMaxGasPriceWei(state, key.Address),
	)
	return r
}

func (ekc *ETHKeysController) formatETHKeyResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if the response has not been written yet
		if !c.Writer.Written() {
			// Get the key and state from the Gin context
			key, keyExists := c.Get("key")
			state, stateExists := c.Get("state")

			// If key and state exist, format the response
			if keyExists && stateExists {
				r := createETHKeyResource(c, ekc, key.(ethkey.KeyV2), state.(ethkey.State))
				jsonAPIResponse(c, r, "keys")
			} else {
				err := errors.Errorf("error getting eth key and state: %v", c)
				jsonAPIError(c, http.StatusInternalServerError, err)
			}
		}
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

	jsonAPIResponseWithStatus(c, resources, "keys", http.StatusOK)
}

// Create adds a new account
// Example:
//
//	"<application>/keys/eth"
func (ekc *ETHKeysController) Create(c *gin.Context) {
	ethKeyStore := ekc.app.GetKeyStore().Eth()

	cid := c.Query("evmChainID")
	chain, ok := ekc.getChain(c, ekc.app.GetChains().EVM, cid)
	if !ok {
		return
	}

	if c.Query("maxGasPriceGWei") != "" {
		jsonAPIError(c, http.StatusBadRequest, v2.ErrUnsupported)
		return
	}

	key, err := ethKeyStore.Create(chain.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	state, err := ethKeyStore.GetState(key.ID(), chain.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.Set("key", key)
	c.Set("state", state)

	ekc.app.GetAuditLogger().Audit(audit.KeyCreated, map[string]interface{}{
		"type": "ethereum",
		"id":   key.ID(),
	})
}

// Delete an ETH key bundle (irreversible!)
// Example:
// "DELETE <application>/keys/eth/:keyID"
func (ekc *ETHKeysController) Delete(c *gin.Context) {
	ethKeyStore := ekc.app.GetKeyStore().Eth()

	keyID := c.Param("address")
	if !common.IsHexAddress(keyID) {
		jsonAPIError(c, http.StatusInternalServerError, errors.Errorf("invalid keyID: %s, must be hex address", keyID))
		return
	}

	key, err := ethKeyStore.Get(keyID)
	if err != nil {
		if errors.Is(err, keystore.ErrKeyNotFound) {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	state, err := ethKeyStore.GetStateForKey(key)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	_, err = ethKeyStore.Delete(keyID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.Set("key", key)
	c.Set("state", state)

	ekc.app.GetAuditLogger().Audit(audit.KeyDeleted, map[string]interface{}{
		"type": "ethereum",
		"id":   keyID,
	})
}

// Import imports a key
func (ekc *ETHKeysController) Import(c *gin.Context) {
	ethKeyStore := ekc.app.GetKeyStore().Eth()
	defer ekc.app.GetLogger().ErrorIfFn(c.Request.Body.Close, "Error closing Import request body")

	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	cid := c.Query("evmChainID")
	chain, ok := ekc.getChain(c, ekc.app.GetChains().EVM, cid)
	if !ok {
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

	c.Set("key", key)
	c.Set("state", state)
	c.Status(http.StatusCreated)

	ekc.app.GetAuditLogger().Audit(audit.KeyImported, map[string]interface{}{
		"type": "ethereum",
		"id":   key.ID(),
	})
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
		jsonAPIError(c, http.StatusInternalServerError, errors.Wrap(err, "invalid address"))
		return
	}
	address := common.BytesToAddress(addressBytes)

	cid := c.Query("evmChainID")
	chain, ok := ekc.getChain(c, ekc.app.GetChains().EVM, cid)
	if !ok {
		return
	}

	var nonce int64 = -1
	if nonceStr := c.Query("nextNonce"); nonceStr != "" {
		nonce, err = strconv.ParseInt(nonceStr, 10, 64)
		if err != nil || nonce < 0 {
			jsonAPIError(c, http.StatusInternalServerError, errors.Wrapf(err, "invalid value for nonce: expected 0 or positive int, got: %s", nonceStr))
			return
		}
	}
	abandon := false
	if abandonStr := c.Query("abandon"); abandonStr != "" {
		abandon, err = strconv.ParseBool(abandonStr)
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, errors.Wrapf(err, "invalid value for abandon: expected boolean, got: %s", abandonStr))
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
			if strings.Contains(err.Error(), "key state not found with address") {
				jsonAPIError(c, http.StatusInternalServerError, err)
			}
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
	}

	enabledStr := c.Query("enabled")
	if enabledStr != "" {
		var enabled bool
		enabled, err = strconv.ParseBool(enabledStr)
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, errors.Wrap(err, "enabled must be bool"))
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
		if errors.Is(err, keystore.ErrKeyNotFound) {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	state, err := kst.GetState(key.ID(), chain.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.Set("key", key)
	c.Set("state", state)
	c.Status(http.StatusOK)
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
		addr := common.HexToAddress(chain.Config().EVM().LinkContractAddress())
		bal, err = ethClient.LINKBalance(ctx, state.Address.Address(), addr)
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
		price = chain.Config().EVM().KeySpecificMaxGasPriceWei(keyAddress)
	}
	return presenters.SetETHKeyMaxGasPriceWei(utils.NewBig(price.ToInt()))
}

// getChain is a convenience wrapper to retrieve a chain for a given request
// and call the corresponding API response error function for 400, 404 and 500 results
func (ekc *ETHKeysController) getChain(c *gin.Context, cs evm.ChainSet, chainIDstr string) (chain evm.Chain, ok bool) {
	chain, err := getChain(ekc.app.GetChains().EVM, chainIDstr)
	if err != nil {
		if errors.Is(err, ErrInvalidChainID) {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return nil, false
		} else if errors.Is(err, ErrMultipleChains) {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return nil, false
		} else if errors.Is(err, ErrMissingChainID) {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return nil, false
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return nil, false
	}
	return chain, true
}

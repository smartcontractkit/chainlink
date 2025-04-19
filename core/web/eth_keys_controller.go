package web

import (
	"context"
	"io"
	"math/big"
	"net/http"
	"sort"
	"strconv"
	"strings"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"

	"github.com/ethereum/go-ethereum/common"
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
	ethBalance := ekc.getEthBalance(c.Request.Context(), state)
	linkBalance := ekc.getLinkBalance(c.Request.Context(), state)
	maxGasPrice := ekc.getKeyMaxGasPriceWei(state, key.Address)

	r := presenters.NewETHKeyResource(key, state,
		ekc.setEthBalance(ethBalance),
		ekc.setLinkBalance(linkBalance),
		ekc.setKeyMaxGasPriceWei(maxGasPrice),
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
	keys, err = ethKeyStore.GetAll(c.Request.Context())
	if err != nil {
		err = errors.Errorf("error getting unlocked keys: %v", err)
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	states, err := ethKeyStore.GetStatesForKeys(c.Request.Context(), keys)
	if err != nil {
		err = errors.Errorf("error getting key states: %v", err)
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	var resources []presenters.ETHKeyResource
	for _, state := range states {
		key, err := ethKeyStore.Get(c.Request.Context(), state.Address.Hex())
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}

		r := createETHKeyResource(c, ekc, key, state)

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
	chain, ok := ekc.getChain(c, cid)
	if !ok {
		return
	}

	if c.Query("maxGasPriceGWei") != "" {
		jsonAPIError(c, http.StatusBadRequest, toml.ErrUnsupported)
		return
	}

	key, err := ethKeyStore.Create(c.Request.Context(), chain.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	state, err := ethKeyStore.GetState(c.Request.Context(), key.ID(), chain.ID())
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
		jsonAPIError(c, http.StatusBadRequest, errors.Errorf("invalid keyID: %s, must be hex address", keyID))
		return
	}

	key, err := ethKeyStore.Get(c.Request.Context(), keyID)
	if err != nil {
		if errors.Is(err, keystore.ErrKeyNotFound) {
			jsonAPIError(c, http.StatusNotFound, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	state, err := ethKeyStore.GetStateForKey(c.Request.Context(), key)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	_, err = ethKeyStore.Delete(c.Request.Context(), keyID)
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
	chain, ok := ekc.getChain(c, cid)
	if !ok {
		return
	}

	key, err := ethKeyStore.Import(c.Request.Context(), bytes, oldPassword, chain.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	state, err := ethKeyStore.GetState(c.Request.Context(), key.ID(), chain.ID())
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

	bytes, err := ekc.app.GetKeyStore().Eth().Export(c.Request.Context(), id, newPassword)
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
	var err error
	kst := ekc.app.GetKeyStore().Eth()
	defer ekc.app.GetLogger().ErrorIfFn(c.Request.Body.Close, "Error closing Import request body")

	keyID := c.Query("address")
	if !common.IsHexAddress(keyID) {
		jsonAPIError(c, http.StatusBadRequest, errors.Errorf("invalid address: %s, must be hex address", keyID))
		return
	}
	address := common.HexToAddress(keyID)

	cid := c.Query("evmChainID")
	chain, ok := ekc.getChain(c, cid)
	if !ok {
		return
	}

	abandon := false
	if abandonStr := c.Query("abandon"); abandonStr != "" {
		abandon, err = strconv.ParseBool(abandonStr)
		if err != nil {
			jsonAPIError(c, http.StatusBadRequest, errors.Wrapf(err, "invalid value for abandon: expected boolean, got: %s", abandonStr))
			return
		}
	}

	// Reset the chain
	if abandon {
		var resetErr error
		err = chain.TxManager().Reset(address, abandon)
		err = multierr.Combine(err, resetErr)
		if err != nil {
			if strings.Contains(err.Error(), "key state not found with address") {
				jsonAPIError(c, http.StatusNotFound, err)
				return
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
			jsonAPIError(c, http.StatusBadRequest, errors.Wrap(err, "enabled must be bool"))
			return
		}

		if enabled {
			err = kst.Enable(c.Request.Context(), address, chain.ID())
		} else {
			err = kst.Disable(c.Request.Context(), address, chain.ID())
		}
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
	}

	key, err := kst.Get(c.Request.Context(), keyID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	state, err := kst.GetState(c.Request.Context(), key.ID(), chain.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.Set("key", key)
	c.Set("state", state)
	c.Status(http.StatusOK)
}

func (ekc *ETHKeysController) setEthBalance(bal *big.Int) presenters.NewETHKeyOption {
	return presenters.SetETHKeyEthBalance((*assets.Eth)(bal))
}

// queries the EthClient for the ETH balance at the address associated with state
func (ekc *ETHKeysController) getEthBalance(ctx context.Context, state ethkey.State) *big.Int {
	chainID := state.EVMChainID.ToInt()
	chain, err := ekc.app.GetRelayers().LegacyEVMChains().Get(chainID.String())
	if err != nil {
		if !errors.Is(errors.Cause(err), evmrelay.ErrNoChains) {
			ekc.lggr.Errorw("Failed to get EVM Chain", "chainID", chainID, "address", state.Address, "err", err)
		}
		return nil
	}

	ethClient := chain.Client()
	bal, err := ethClient.BalanceAt(ctx, state.Address.Address(), nil)
	if err != nil {
		ekc.lggr.Errorw("Failed to get ETH balance", "chainID", chainID, "address", state.Address, "err", err)
		return nil
	}

	return bal
}

func (ekc *ETHKeysController) setLinkBalance(bal *commonassets.Link) presenters.NewETHKeyOption {
	return presenters.SetETHKeyLinkBalance(bal)
}

// queries the EthClient for the LINK balance at the address associated with state
func (ekc *ETHKeysController) getLinkBalance(ctx context.Context, state ethkey.State) *commonassets.Link {
	var bal *commonassets.Link
	chainID := state.EVMChainID.ToInt()
	chain, err := ekc.app.GetRelayers().LegacyEVMChains().Get(chainID.String())
	if err != nil {
		if !errors.Is(errors.Cause(err), evmrelay.ErrNoChains) {
			ekc.lggr.Errorw("Failed to get EVM Chain", "chainID", chainID, "err", err)
		}
	} else {
		ethClient := chain.Client()
		addr := common.HexToAddress(chain.Config().EVM().LinkContractAddress())
		bal, err = ethClient.LINKBalance(ctx, state.Address.Address(), addr)
		if err != nil {
			ekc.lggr.Errorw("Failed to get LINK balance", "chainID", chainID, "address", state.Address, "err", err)
		}
	}
	return bal
}

// setKeyMaxGasPriceWei is a custom functional option for NewEthKeyResource which
// gets the key specific max gas price from the chain config and sets it on the
// resource.
func (ekc *ETHKeysController) setKeyMaxGasPriceWei(price *assets.Wei) presenters.NewETHKeyOption {
	return presenters.SetETHKeyMaxGasPriceWei(ubig.New(price.ToInt()))
}

func (ekc *ETHKeysController) getKeyMaxGasPriceWei(state ethkey.State, keyAddress common.Address) *assets.Wei {
	var price *assets.Wei
	chainID := state.EVMChainID.ToInt()
	chain, err := ekc.app.GetRelayers().LegacyEVMChains().Get(chainID.String())
	if err != nil {
		if !errors.Is(errors.Cause(err), evmrelay.ErrNoChains) {
			ekc.lggr.Errorw("Failed to get EVM Chain", "chainID", chainID, "err", err)
		}
	} else {
		price = chain.Config().EVM().GasEstimator().PriceMaxKey(keyAddress)
	}
	return price
}

// getChain is a convenience wrapper to retrieve a chain for a given request
// and call the corresponding API response error function for 400, 404 and 500 results
func (ekc *ETHKeysController) getChain(c *gin.Context, chainIDstr string) (chain legacyevm.Chain, ok bool) {
	chain, err := getChain(ekc.app.GetRelayers().LegacyEVMChains(), chainIDstr)
	if err != nil {
		if errors.Is(err, ErrInvalidChainID) || errors.Is(err, ErrMultipleChains) {
			jsonAPIError(c, http.StatusBadRequest, err)
			return nil, false
		} else if errors.Is(err, ErrMissingChainID) {
			jsonAPIError(c, http.StatusNotFound, err)
			return nil, false
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return nil, false
	}
	return chain, true
}

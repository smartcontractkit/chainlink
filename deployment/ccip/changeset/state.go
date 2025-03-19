package changeset

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"golang.org/x/sync/errgroup"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_1/burn_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/evm_2_evm_onramp"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/latest/log_message_data_receiver"
	price_registry_1_2_0 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_1/burn_with_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_1/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc20"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc677"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/latest/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/latest/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_1/usdc_token_pool"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view"
	viewv1_0 "github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_2"
	viewv1_5 "github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_5"
	viewv1_5_1 "github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_5_1"
	viewv1_6 "github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_6"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	common_v1_0 "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/helpers"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/token_admin_registry"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/fee_quoter"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/latest/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_0_0/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/rmn_remote"
	capabilities_registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/multicall3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/weth9"
)

var (
	// Legacy
	CommitStore   deployment.ContractType = "CommitStore"
	PriceRegistry deployment.ContractType = "PriceRegistry"
	RMN           deployment.ContractType = "RMN"

	// Not legacy
	MockRMN              deployment.ContractType = "MockRMN"
	RMNRemote            deployment.ContractType = "RMNRemote"
	ARMProxy             deployment.ContractType = "ARMProxy"
	WETH9                deployment.ContractType = "WETH9"
	Router               deployment.ContractType = "Router"
	TokenAdminRegistry   deployment.ContractType = "TokenAdminRegistry"
	RegistryModule       deployment.ContractType = "RegistryModuleOwnerCustom"
	NonceManager         deployment.ContractType = "NonceManager"
	FeeQuoter            deployment.ContractType = "FeeQuoter"
	CCIPHome             deployment.ContractType = "CCIPHome"
	RMNHome              deployment.ContractType = "RMNHome"
	OnRamp               deployment.ContractType = "OnRamp"
	OffRamp              deployment.ContractType = "OffRamp"
	CapabilitiesRegistry deployment.ContractType = "CapabilitiesRegistry"
	PriceFeed            deployment.ContractType = "PriceFeed"

	// Test contracts. Note test router maps to a regular router contract.
	TestRouter             deployment.ContractType = "TestRouter"
	Multicall3             deployment.ContractType = "Multicall3"
	CCIPReceiver           deployment.ContractType = "CCIPReceiver"
	LogMessageDataReceiver deployment.ContractType = "LogMessageDataReceiver"
	USDCMockTransmitter    deployment.ContractType = "USDCMockTransmitter"

	// Pools
	BurnMintToken                  deployment.ContractType = "BurnMintToken"
	ERC20Token                     deployment.ContractType = "ERC20Token"
	ERC677Token                    deployment.ContractType = "ERC677Token"
	BurnMintTokenPool              deployment.ContractType = "BurnMintTokenPool"
	BurnWithFromMintTokenPool      deployment.ContractType = "BurnWithFromMintTokenPool"
	BurnFromMintTokenPool          deployment.ContractType = "BurnFromMintTokenPool"
	LockReleaseTokenPool           deployment.ContractType = "LockReleaseTokenPool"
	USDCToken                      deployment.ContractType = "USDCToken"
	USDCTokenMessenger             deployment.ContractType = "USDCTokenMessenger"
	USDCTokenPool                  deployment.ContractType = "USDCTokenPool"
	HybridLockReleaseUSDCTokenPool deployment.ContractType = "HybridLockReleaseUSDCTokenPool"

	// Firedrill
	FiredrillEntrypointType deployment.ContractType = "FiredrillEntrypoint"
)

// CCIPChainState holds a Go binding for all the currently deployed CCIP contracts
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type CCIPChainState struct {
	commoncs.MCMSWithTimelockState
	commoncs.LinkTokenState
	commoncs.StaticLinkTokenState
	ABIByAddress       map[string]string
	OnRamp             onramp.OnRampInterface
	OffRamp            offramp.OffRampInterface
	FeeQuoter          *fee_quoter.FeeQuoter
	RMNProxy           *rmn_proxy_contract.RMNProxy
	NonceManager       *nonce_manager.NonceManager
	TokenAdminRegistry *token_admin_registry.TokenAdminRegistry
	RegistryModule     *registry_module_owner_custom.RegistryModuleOwnerCustom
	Router             *router.Router
	Weth9              *weth9.WETH9
	RMNRemote          *rmn_remote.RMNRemote
	// Map between token Descriptor (e.g. LinkSymbol, WethSymbol)
	// and the respective token / token pool contract(s) (only one of which would be active on the registry).
	// This is more of an illustration of how we'll have tokens, and it might need some work later to work properly.
	ERC20Tokens                map[TokenSymbol]*erc20.ERC20
	ERC677Tokens               map[TokenSymbol]*erc677.ERC677
	BurnMintTokens677          map[TokenSymbol]*burn_mint_erc677.BurnMintERC677
	BurnMintTokenPools         map[TokenSymbol]map[semver.Version]*burn_mint_token_pool.BurnMintTokenPool
	BurnWithFromMintTokenPools map[TokenSymbol]map[semver.Version]*burn_with_from_mint_token_pool.BurnWithFromMintTokenPool
	BurnFromMintTokenPools     map[TokenSymbol]map[semver.Version]*burn_from_mint_token_pool.BurnFromMintTokenPool
	USDCTokenPools             map[semver.Version]*usdc_token_pool.USDCTokenPool
	LockReleaseTokenPools      map[TokenSymbol]map[semver.Version]*lock_release_token_pool.LockReleaseTokenPool
	// Map between token Symbol (e.g. LinkSymbol, WethSymbol)
	// and the respective aggregator USD feed contract
	USDFeeds map[TokenSymbol]*aggregator_v3_interface.AggregatorV3Interface

	// Note we only expect one of these (on the home chain)
	CapabilityRegistry *capabilities_registry.CapabilitiesRegistry
	CCIPHome           *ccip_home.CCIPHome
	RMNHome            *rmn_home.RMNHome

	// Test contracts
	Receiver               maybe_revert_message_receiver.MaybeRevertMessageReceiverInterface
	LogMessageDataReceiver *log_message_data_receiver.LogMessageDataReceiver
	TestRouter             *router.Router
	MockUSDCTransmitter    *mock_usdc_token_transmitter.MockE2EUSDCTransmitter
	MockUSDCTokenMessenger *mock_usdc_token_messenger.MockE2EUSDCTokenMessenger
	Multicall3             *multicall3.Multicall3

	// Legacy contracts
	EVM2EVMOnRamp  map[uint64]*evm_2_evm_onramp.EVM2EVMOnRamp   // mapping of dest chain selector -> EVM2EVMOnRamp
	CommitStore    map[uint64]*commit_store.CommitStore         // mapping of source chain selector -> CommitStore
	EVM2EVMOffRamp map[uint64]*evm_2_evm_offramp.EVM2EVMOffRamp // mapping of source chain selector -> EVM2EVMOffRamp
	MockRMN        *mock_rmn_contract.MockRMNContract
	PriceRegistry  *price_registry_1_2_0.PriceRegistry
	RMN            *rmn_contract.RMNContract
}

func (c CCIPChainState) TokenAddressBySymbol() (map[TokenSymbol]common.Address, error) {
	tokenAddresses := make(map[TokenSymbol]common.Address)
	for symbol, token := range c.ERC20Tokens {
		tokenAddresses[symbol] = token.Address()
	}
	for symbol, token := range c.ERC677Tokens {
		tokenAddresses[symbol] = token.Address()
	}
	for symbol, token := range c.BurnMintTokens677 {
		tokenAddresses[symbol] = token.Address()
	}
	var err error
	tokenAddresses[LinkSymbol], err = c.LinkTokenAddress()
	if err != nil {
		return nil, err
	}
	if c.Weth9 == nil {
		return nil, errors.New("no WETH contract found in the state")
	}
	tokenAddresses[WethSymbol] = c.Weth9.Address()
	return tokenAddresses, nil
}

// TokenDetailsBySymbol get token mapping from the state. It contains only tokens that we have in address book
func (c CCIPChainState) TokenDetailsBySymbol() (map[TokenSymbol]TokenDetails, error) {
	tokenDetails := make(map[TokenSymbol]TokenDetails)
	for symbol, token := range c.ERC20Tokens {
		tokenDetails[symbol] = token
	}
	for symbol, token := range c.ERC677Tokens {
		tokenDetails[symbol] = token
	}
	for symbol, token := range c.BurnMintTokens677 {
		tokenDetails[symbol] = token
	}
	if c.LinkToken != nil {
		tokenDetails[LinkSymbol] = c.LinkToken
	}
	if c.StaticLinkToken != nil {
		tokenDetails[LinkSymbol] = c.StaticLinkToken
	}

	if _, ok := tokenDetails[LinkSymbol]; !ok {
		return nil, errors.New("no LINK contract found in the state")
	}

	if c.Weth9 == nil {
		return nil, errors.New("no WETH contract found in the state")
	}
	tokenDetails[WethSymbol] = c.Weth9
	return tokenDetails, nil
}

func (c CCIPChainState) LinkTokenAddress() (common.Address, error) {
	if c.LinkToken != nil {
		return c.LinkToken.Address(), nil
	}
	if c.StaticLinkToken != nil {
		return c.StaticLinkToken.Address(), nil
	}
	return common.Address{}, errors.New("no link token found in the state")
}

func (c CCIPChainState) GenerateView() (view.ChainView, error) {
	chainView := view.NewChain()
	if c.Router != nil {
		routerView, err := v1_2.GenerateRouterView(c.Router, false)
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate router view for router %s", c.Router.Address().String())
		}
		chainView.Router[c.Router.Address().Hex()] = routerView
	}
	if c.TestRouter != nil {
		testRouterView, err := v1_2.GenerateRouterView(c.TestRouter, true)
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate router view for test router %s", c.TestRouter.Address().String())
		}
		chainView.Router[c.TestRouter.Address().Hex()] = testRouterView
	}
	if c.TokenAdminRegistry != nil {
		taView, err := viewv1_5.GenerateTokenAdminRegistryView(c.TokenAdminRegistry)
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate token admin registry view for token admin registry %s", c.TokenAdminRegistry.Address().String())
		}
		chainView.TokenAdminRegistry[c.TokenAdminRegistry.Address().Hex()] = taView
	}
	tpUpdateGrp := errgroup.Group{}
	for tokenSymbol, versionToPool := range c.BurnMintTokenPools {
		for _, tokenPool := range versionToPool {
			tpUpdateGrp.Go(func() error {
				tokenPoolView, err := viewv1_5_1.GenerateTokenPoolView(tokenPool, c.usdFeedOrDefault(tokenSymbol))
				if err != nil {
					return errors.Wrapf(err, "failed to generate burn mint token pool view for %s", tokenPool.Address().String())
				}
				chainView.UpdateTokenPool(tokenSymbol.String(), tokenPool.Address().Hex(), viewv1_5_1.PoolView{
					TokenPoolView: tokenPoolView,
				})
				return nil
			})
		}
	}
	for tokenSymbol, versionToPool := range c.BurnWithFromMintTokenPools {
		for _, tokenPool := range versionToPool {
			tpUpdateGrp.Go(func() error {
				tokenPoolView, err := viewv1_5_1.GenerateTokenPoolView(tokenPool, c.usdFeedOrDefault(tokenSymbol))
				if err != nil {
					return errors.Wrapf(err, "failed to generate burn mint token pool view for %s", tokenPool.Address().String())
				}
				chainView.UpdateTokenPool(tokenSymbol.String(), tokenPool.Address().Hex(), viewv1_5_1.PoolView{
					TokenPoolView: tokenPoolView,
				})
				return nil
			})
		}
	}
	for tokenSymbol, versionToPool := range c.BurnFromMintTokenPools {
		for _, tokenPool := range versionToPool {
			tpUpdateGrp.Go(func() error {
				tokenPoolView, err := viewv1_5_1.GenerateTokenPoolView(tokenPool, c.usdFeedOrDefault(tokenSymbol))
				if err != nil {
					return errors.Wrapf(err, "failed to generate burn mint token pool view for %s", tokenPool.Address().String())
				}
				chainView.UpdateTokenPool(tokenSymbol.String(), tokenPool.Address().Hex(), viewv1_5_1.PoolView{
					TokenPoolView: tokenPoolView,
				})
				return nil
			})
		}
	}
	for tokenSymbol, versionToPool := range c.LockReleaseTokenPools {
		for _, tokenPool := range versionToPool {
			tpUpdateGrp.Go(func() error {
				tokenPoolView, err := viewv1_5_1.GenerateLockReleaseTokenPoolView(tokenPool, c.usdFeedOrDefault(tokenSymbol))
				if err != nil {
					return errors.Wrapf(err, "failed to generate lock release token pool view for %s", tokenPool.Address().String())
				}
				chainView.UpdateTokenPool(tokenSymbol.String(), tokenPool.Address().Hex(), tokenPoolView)
				return nil
			})
		}
	}
	for _, pool := range c.USDCTokenPools {
		tpUpdateGrp.Go(func() error {
			tokenPoolView, err := viewv1_5_1.GenerateUSDCTokenPoolView(pool)
			if err != nil {
				return errors.Wrapf(err, "failed to generate USDC token pool view for %s", pool.Address().String())
			}
			chainView.UpdateTokenPool(string(USDCSymbol), pool.Address().Hex(), tokenPoolView)
			return nil
		})
	}
	// wait for all pool updates to finish to ensure we are not rate limited by rpc end point by a lot of concurrent calls for other contract queries
	if err := tpUpdateGrp.Wait(); err != nil {
		return chainView, err
	}
	if c.NonceManager != nil {
		nmView, err := viewv1_6.GenerateNonceManagerView(c.NonceManager)
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate nonce manager view for nonce manager %s", c.NonceManager.Address().String())
		}
		chainView.NonceManager[c.NonceManager.Address().Hex()] = nmView
	}
	if c.RMNRemote != nil {
		rmnView, err := viewv1_6.GenerateRMNRemoteView(c.RMNRemote)
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate rmn remote view for rmn remote %s", c.RMNRemote.Address().String())
		}
		chainView.RMNRemote[c.RMNRemote.Address().Hex()] = rmnView
	}

	if c.RMNHome != nil {
		rmnHomeView, err := viewv1_6.GenerateRMNHomeView(c.RMNHome)
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate rmn home view for rmn home %s", c.RMNHome.Address().String())
		}
		chainView.RMNHome[c.RMNHome.Address().Hex()] = rmnHomeView
	}

	if c.FeeQuoter != nil && c.Router != nil && c.TokenAdminRegistry != nil {
		// FeeQuoter knows only about tokens that managed by CCIP (i.e. imported from address book)
		tokenDetails, err := c.TokenDetailsBySymbol()
		if err != nil {
			return chainView, err
		}
		tokens := make([]common.Address, 0, len(tokenDetails))
		for _, tokenDetail := range tokenDetails {
			tokens = append(tokens, tokenDetail.Address())
		}
		fqView, err := viewv1_6.GenerateFeeQuoterView(c.FeeQuoter, c.Router, tokens)
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate fee quoter view for fee quoter %s", c.FeeQuoter.Address().String())
		}
		chainView.FeeQuoter[c.FeeQuoter.Address().Hex()] = fqView
	}

	if c.OnRamp != nil && c.Router != nil && c.TokenAdminRegistry != nil {
		onRampView, err := viewv1_6.GenerateOnRampView(
			c.OnRamp,
			c.Router,
			c.TokenAdminRegistry,
		)
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate on ramp view for on ramp %s", c.OnRamp.Address().String())
		}
		chainView.OnRamp[c.OnRamp.Address().Hex()] = onRampView
	}

	if c.OffRamp != nil && c.Router != nil {
		offRampView, err := viewv1_6.GenerateOffRampView(
			c.OffRamp,
			c.Router,
		)
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate off ramp view for off ramp %s", c.OffRamp.Address().String())
		}
		chainView.OffRamp[c.OffRamp.Address().Hex()] = offRampView
	}

	if c.RMNProxy != nil {
		rmnProxyView, err := viewv1_0.GenerateRMNProxyView(c.RMNProxy)
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate rmn proxy view for rmn proxy %s", c.RMNProxy.Address().String())
		}
		chainView.RMNProxy[c.RMNProxy.Address().Hex()] = rmnProxyView
	}
	if c.CCIPHome != nil && c.CapabilityRegistry != nil {
		chView, err := viewv1_6.GenerateCCIPHomeView(c.CapabilityRegistry, c.CCIPHome)
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate CCIP home view for CCIP home %s", c.CCIPHome.Address())
		}
		chainView.CCIPHome[c.CCIPHome.Address().Hex()] = chView
	}
	if c.CapabilityRegistry != nil {
		capRegView, err := common_v1_0.GenerateCapabilityRegistryView(c.CapabilityRegistry)
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate capability registry view for capability registry %s", c.CapabilityRegistry.Address().String())
		}
		chainView.CapabilityRegistry[c.CapabilityRegistry.Address().Hex()] = capRegView
	}
	if c.MCMSWithTimelockState.Timelock != nil {
		mcmsView, err := c.MCMSWithTimelockState.GenerateMCMSWithTimelockView()
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate MCMS with timelock view for MCMS with timelock %s", c.MCMSWithTimelockState.Timelock.Address().String())
		}
		chainView.MCMSWithTimelock = mcmsView
	}
	if c.LinkToken != nil {
		linkTokenView, err := c.GenerateLinkView()
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate link token view for link token %s", c.LinkToken.Address().String())
		}
		chainView.LinkToken = linkTokenView
	}
	if c.StaticLinkToken != nil {
		staticLinkTokenView, err := c.GenerateStaticLinkView()
		if err != nil {
			return chainView, err
		}
		chainView.StaticLinkToken = staticLinkTokenView
	}
	// Legacy contracts
	if c.CommitStore != nil {
		for source, commitStore := range c.CommitStore {
			commitStoreView, err := viewv1_5.GenerateCommitStoreView(commitStore)
			if err != nil {
				return chainView, errors.Wrapf(err, "failed to generate commit store view for commit store %s for source %d", commitStore.Address().String(), source)
			}
			chainView.CommitStore[commitStore.Address().Hex()] = commitStoreView
		}
	}

	if c.PriceRegistry != nil {
		priceRegistryView, err := v1_2.GeneratePriceRegistryView(c.PriceRegistry)
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate price registry view for price registry %s", c.PriceRegistry.Address().String())
		}
		chainView.PriceRegistry[c.PriceRegistry.Address().String()] = priceRegistryView
	}

	if c.RMN != nil {
		rmnView, err := viewv1_5.GenerateRMNView(c.RMN)
		if err != nil {
			return chainView, errors.Wrapf(err, "failed to generate rmn view for rmn %s", c.RMN.Address().String())
		}
		chainView.RMN[c.RMN.Address().Hex()] = rmnView
	}

	if c.EVM2EVMOffRamp != nil {
		for source, offRamp := range c.EVM2EVMOffRamp {
			offRampView, err := viewv1_5.GenerateOffRampView(offRamp)
			if err != nil {
				return chainView, errors.Wrapf(err, "failed to generate off ramp view for off ramp %s for source %d", offRamp.Address().String(), source)
			}
			chainView.EVM2EVMOffRamp[offRamp.Address().Hex()] = offRampView
		}
	}

	if c.EVM2EVMOnRamp != nil {
		for dest, onRamp := range c.EVM2EVMOnRamp {
			onRampView, err := viewv1_5.GenerateOnRampView(onRamp)
			if err != nil {
				return chainView, errors.Wrapf(err, "failed to generate on ramp view for on ramp %s for dest %d", onRamp.Address().String(), dest)
			}
			chainView.EVM2EVMOnRamp[onRamp.Address().Hex()] = onRampView
		}
	}

	return chainView, nil
}

func (c CCIPChainState) usdFeedOrDefault(symbol TokenSymbol) common.Address {
	if feed, ok := c.USDFeeds[symbol]; ok {
		return feed.Address()
	}
	return common.Address{}
}

// CCIPOnChainState state always derivable from an address book.
// Offchain state always derivable from a list of nodeIds.
// Note can translate this into Go struct needed for MCMS/Docs/UI.
type CCIPOnChainState struct {
	// Populated go bindings for the appropriate version for all contracts.
	// We would hold 2 versions of each contract here. Once we upgrade we can phase out the old one.
	// When generating bindings, make sure the package name corresponds to the version.
	Chains      map[uint64]CCIPChainState
	SolChains   map[uint64]SolCCIPChainState
	AptosChains map[uint64]AptosCCIPChainState
}

func (s CCIPOnChainState) OffRampPermissionLessExecutionThresholdSeconds(ctx context.Context, env deployment.Environment, selector uint64) (uint32, error) {
	family, err := chain_selectors.GetSelectorFamily(selector)
	if err != nil {
		return 0, err
	}
	switch family {
	case chain_selectors.FamilyEVM:
		c, ok := s.Chains[selector]
		if !ok {
			return 0, fmt.Errorf("chain %d not found in the state", selector)
		}
		offRamp := c.OffRamp
		if offRamp == nil {
			return 0, fmt.Errorf("offramp not found in the state for chain %d", selector)
		}
		dCfg, err := offRamp.GetDynamicConfig(&bind.CallOpts{
			Context: ctx,
		})
		if err != nil {
			return dCfg.PermissionLessExecutionThresholdSeconds, fmt.Errorf("fetch dynamic config from offRamp %s for chain %d: %w", offRamp.Address().String(), selector, err)
		}
		return dCfg.PermissionLessExecutionThresholdSeconds, nil
	case chain_selectors.FamilySolana:
		c, ok := s.SolChains[selector]
		if !ok {
			return 0, fmt.Errorf("chain %d not found in the state", selector)
		}
		chain, ok := env.SolChains[selector]
		if !ok {
			return 0, fmt.Errorf("solana chain %d not found in the environment", selector)
		}
		if c.OffRamp.IsZero() {
			return 0, fmt.Errorf("offramp not found in existing state, deploy the offramp first for chain %d", selector)
		}
		var offRampConfig solOffRamp.Config
		offRampConfigPDA, _, _ := solState.FindOfframpConfigPDA(c.OffRamp)
		err := chain.GetAccountDataBorshInto(context.Background(), offRampConfigPDA, &offRampConfig)
		if err != nil {
			return 0, fmt.Errorf("offramp config not found in existing state, initialize the offramp first %d", chain.Selector)
		}
		// #nosec G115
		return uint32(offRampConfig.EnableManualExecutionAfter), nil
	}
	return 0, fmt.Errorf("unsupported chain family %s", family)
}

func (s CCIPOnChainState) Validate() error {
	for sel, chain := range s.Chains {
		// cannot have static link and link together
		if chain.LinkToken != nil && chain.StaticLinkToken != nil {
			return fmt.Errorf("cannot have both link and static link token on the same chain %d", sel)
		}
	}
	return nil
}

func (s CCIPOnChainState) GetAllProposerMCMSForChains(chains []uint64) (map[uint64]*gethwrappers.ManyChainMultiSig, error) {
	multiSigs := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for _, chain := range chains {
		chainState, ok := s.Chains[chain]
		if !ok {
			return nil, fmt.Errorf("chain %d not found", chain)
		}
		if chainState.ProposerMcm == nil {
			return nil, fmt.Errorf("proposer mcm not found for chain %d", chain)
		}
		multiSigs[chain] = chainState.ProposerMcm
	}
	return multiSigs, nil
}

func (s CCIPOnChainState) GetAllTimeLocksForChains(chains []uint64) (map[uint64]common.Address, error) {
	timelocks := make(map[uint64]common.Address)
	for _, chain := range chains {
		chainState, ok := s.Chains[chain]
		if !ok {
			return nil, fmt.Errorf("chain %d not found", chain)
		}
		if chainState.Timelock == nil {
			return nil, fmt.Errorf("timelock not found for chain %d", chain)
		}
		timelocks[chain] = chainState.Timelock.Address()
	}
	return timelocks, nil
}

func (s CCIPOnChainState) SupportedChains() map[uint64]struct{} {
	chains := make(map[uint64]struct{})
	for chain := range s.Chains {
		chains[chain] = struct{}{}
	}
	for chain := range s.SolChains {
		chains[chain] = struct{}{}
	}
	return chains
}

func (s CCIPOnChainState) View(chains []uint64) (map[string]view.ChainView, error) {
	m := make(map[string]view.ChainView)
	mu := sync.Mutex{}
	grp := errgroup.Group{}
	for _, chainSelector := range chains {
		var name string
		var chainView view.ChainView
		chainSelector := chainSelector
		grp.Go(func() error {
			chainInfo, err := deployment.ChainInfo(chainSelector)
			if err != nil {
				return err
			}
			if _, ok := s.Chains[chainSelector]; !ok {
				return fmt.Errorf("chain not supported %d", chainSelector)
			}
			chainState := s.Chains[chainSelector]
			chainView, err = chainState.GenerateView()
			if err != nil {
				return err
			}
			name = chainInfo.ChainName
			if chainInfo.ChainName == "" {
				name = strconv.FormatUint(chainSelector, 10)
			}
			chainView.ChainSelector = chainSelector
			id, err := chain_selectors.GetChainIDFromSelector(chainSelector)
			if err != nil {
				return fmt.Errorf("failed to get chain id from selector %d: %w", chainSelector, err)
			}
			chainView.ChainID = id
			mu.Lock()
			m[name] = chainView
			mu.Unlock()
			return nil
		})
	}
	return m, grp.Wait()
}

func (s CCIPOnChainState) GetOffRampAddressBytes(chainSelector uint64) ([]byte, error) {
	family, err := chain_selectors.GetSelectorFamily(chainSelector)
	if err != nil {
		return nil, err
	}

	var offRampAddress []byte
	switch family {
	case chain_selectors.FamilyEVM:
		offRampAddress = s.Chains[chainSelector].OffRamp.Address().Bytes()
	case chain_selectors.FamilySolana:
		offRampAddress = s.SolChains[chainSelector].OffRamp.Bytes()
	default:
		return nil, fmt.Errorf("unsupported chain family %s", family)
	}

	return offRampAddress, nil
}

func (s CCIPOnChainState) GetOnRampAddressBytes(chainSelector uint64) ([]byte, error) {
	family, err := chain_selectors.GetSelectorFamily(chainSelector)
	if err != nil {
		return nil, err
	}

	var onRampAddressBytes []byte
	switch family {
	case chain_selectors.FamilyEVM:
		if s.Chains[chainSelector].OnRamp == nil {
			return nil, fmt.Errorf("no onramp found in the state for chain %d", chainSelector)
		}
		onRampAddressBytes = s.Chains[chainSelector].OnRamp.Address().Bytes()
	case chain_selectors.FamilySolana:
		if s.SolChains[chainSelector].Router.IsZero() {
			return nil, fmt.Errorf("no router found in the state for chain %d", chainSelector)
		}
		onRampAddressBytes = s.SolChains[chainSelector].Router.Bytes()
	default:
		return nil, fmt.Errorf("unsupported chain family %s", family)
	}

	return onRampAddressBytes, nil
}

func (s CCIPOnChainState) ValidateRamp(chainSelector uint64, rampType deployment.ContractType) error {
	family, err := chain_selectors.GetSelectorFamily(chainSelector)
	if err != nil {
		return err
	}
	switch family {
	case chain_selectors.FamilyEVM:
		chainState, exists := s.Chains[chainSelector]
		if !exists {
			return fmt.Errorf("chain %d does not exist", chainSelector)
		}
		switch rampType {
		case OffRamp:
			if chainState.OffRamp == nil {
				return fmt.Errorf("offramp contract does not exist on evm chain %d", chainSelector)
			}
		case OnRamp:
			if chainState.OnRamp == nil {
				return fmt.Errorf("onramp contract does not exist on evm chain %d", chainSelector)
			}
		default:
			return fmt.Errorf("unknown ramp type %s", rampType)
		}

	case chain_selectors.FamilySolana:
		chainState, exists := s.SolChains[chainSelector]
		if !exists {
			return fmt.Errorf("chain %d does not exist", chainSelector)
		}
		switch rampType {
		case OffRamp:
			if chainState.OffRamp.IsZero() {
				return fmt.Errorf("offramp contract does not exist on solana chain %d", chainSelector)
			}
		case OnRamp:
			if chainState.Router.IsZero() {
				return fmt.Errorf("router contract does not exist on solana chain %d", chainSelector)
			}
		default:
			return fmt.Errorf("unknown ramp type %s", rampType)
		}

	default:
		return fmt.Errorf("unknown chain family %s", family)
	}
	return nil
}

func LoadOnchainState(e deployment.Environment) (CCIPOnChainState, error) {
	solState, err := LoadOnchainStateSolana(e)
	if err != nil {
		return CCIPOnChainState{}, err
	}
	state := CCIPOnChainState{
		Chains:    make(map[uint64]CCIPChainState),
		SolChains: solState.SolChains,
	}
	for chainSelector, chain := range e.Chains {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if errors.Is(err, deployment.ErrChainNotFound) {
				addresses = make(map[string]deployment.TypeAndVersion)
			} else {
				return state, err
			}
		}
		chainState, err := LoadChainState(e.GetContext(), chain, addresses)
		if err != nil {
			return state, err
		}
		state.Chains[chainSelector] = chainState
	}
	return state, state.Validate()
}

// LoadChainState Loads all state for a chain into state
func LoadChainState(ctx context.Context, chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (CCIPChainState, error) {
	var state CCIPChainState
	mcmsWithTimelock, err := commoncs.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
	if err != nil {
		return state, err
	}
	state.MCMSWithTimelockState = *mcmsWithTimelock

	linkState, err := commoncs.MaybeLoadLinkTokenChainState(chain, addresses)
	if err != nil {
		return state, err
	}
	state.LinkTokenState = *linkState
	staticLinkState, err := commoncs.MaybeLoadStaticLinkTokenState(chain, addresses)
	if err != nil {
		return state, err
	}
	state.StaticLinkTokenState = *staticLinkState
	state.ABIByAddress = make(map[string]string)
	for address, tvStr := range addresses {
		switch tvStr.String() {
		case deployment.NewTypeAndVersion(commontypes.RBACTimelock, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = gethwrappers.RBACTimelockABI
		case deployment.NewTypeAndVersion(commontypes.CallProxy, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = gethwrappers.CallProxyABI
		case deployment.NewTypeAndVersion(commontypes.ProposerManyChainMultisig, deployment.Version1_0_0).String(),
			deployment.NewTypeAndVersion(commontypes.CancellerManyChainMultisig, deployment.Version1_0_0).String(),
			deployment.NewTypeAndVersion(commontypes.BypasserManyChainMultisig, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = gethwrappers.ManyChainMultiSigABI
		case deployment.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = link_token.LinkTokenABI
		case deployment.NewTypeAndVersion(commontypes.StaticLinkToken, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = link_token_interface.LinkTokenABI
		case deployment.NewTypeAndVersion(CapabilitiesRegistry, deployment.Version1_0_0).String():
			cr, err := capabilities_registry.NewCapabilitiesRegistry(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.CapabilityRegistry = cr
			state.ABIByAddress[address] = capabilities_registry.CapabilitiesRegistryABI
		case deployment.NewTypeAndVersion(OnRamp, deployment.Version1_6_0).String():
			onRampC, err := onramp.NewOnRamp(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.OnRamp = onRampC
			state.ABIByAddress[address] = onramp.OnRampABI
		case deployment.NewTypeAndVersion(OffRamp, deployment.Version1_6_0).String():
			offRamp, err := offramp.NewOffRamp(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.OffRamp = offRamp
			state.ABIByAddress[address] = offramp.OffRampABI
		case deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_0_0).String():
			armProxy, err := rmn_proxy_contract.NewRMNProxy(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMNProxy = armProxy
			state.ABIByAddress[address] = rmn_proxy_contract.RMNProxyABI
		case deployment.NewTypeAndVersion(RMNRemote, deployment.Version1_6_0).String():
			rmnRemote, err := rmn_remote.NewRMNRemote(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMNRemote = rmnRemote
			state.ABIByAddress[address] = rmn_remote.RMNRemoteABI
		case deployment.NewTypeAndVersion(RMNHome, deployment.Version1_6_0).String():
			rmnHome, err := rmn_home.NewRMNHome(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMNHome = rmnHome
			state.ABIByAddress[address] = rmn_home.RMNHomeABI
		case deployment.NewTypeAndVersion(WETH9, deployment.Version1_0_0).String():
			_weth9, err := weth9.NewWETH9(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Weth9 = _weth9
			state.ABIByAddress[address] = weth9.WETH9ABI
		case deployment.NewTypeAndVersion(NonceManager, deployment.Version1_6_0).String():
			nm, err := nonce_manager.NewNonceManager(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.NonceManager = nm
			state.ABIByAddress[address] = nonce_manager.NonceManagerABI
		case deployment.NewTypeAndVersion(TokenAdminRegistry, deployment.Version1_5_0).String():
			tm, err := token_admin_registry.NewTokenAdminRegistry(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.TokenAdminRegistry = tm
			state.ABIByAddress[address] = token_admin_registry.TokenAdminRegistryABI
		case deployment.NewTypeAndVersion(RegistryModule, deployment.Version1_5_0).String():
			rm, err := registry_module_owner_custom.NewRegistryModuleOwnerCustom(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RegistryModule = rm
			state.ABIByAddress[address] = registry_module_owner_custom.RegistryModuleOwnerCustomABI
		case deployment.NewTypeAndVersion(Router, deployment.Version1_2_0).String():
			r, err := router.NewRouter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Router = r
			state.ABIByAddress[address] = router.RouterABI
		case deployment.NewTypeAndVersion(TestRouter, deployment.Version1_2_0).String():
			r, err := router.NewRouter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.TestRouter = r
			state.ABIByAddress[address] = router.RouterABI
		case deployment.NewTypeAndVersion(FeeQuoter, deployment.Version1_6_0).String():
			fq, err := fee_quoter.NewFeeQuoter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.FeeQuoter = fq
			state.ABIByAddress[address] = fee_quoter.FeeQuoterABI
		case deployment.NewTypeAndVersion(USDCToken, deployment.Version1_0_0).String():
			ut, err := burn_mint_erc677.NewBurnMintERC677(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.BurnMintTokens677 = map[TokenSymbol]*burn_mint_erc677.BurnMintERC677{
				USDCSymbol: ut,
			}
			state.ABIByAddress[address] = burn_mint_erc677.BurnMintERC677ABI
		case deployment.NewTypeAndVersion(USDCTokenPool, deployment.Version1_5_1).String():
			utp, err := usdc_token_pool.NewUSDCTokenPool(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.USDCTokenPools == nil {
				state.USDCTokenPools = make(map[semver.Version]*usdc_token_pool.USDCTokenPool)
			}
			state.USDCTokenPools[deployment.Version1_5_1] = utp
		case deployment.NewTypeAndVersion(HybridLockReleaseUSDCTokenPool, deployment.Version1_5_1).String():
			utp, err := usdc_token_pool.NewUSDCTokenPool(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.USDCTokenPools == nil {
				state.USDCTokenPools = make(map[semver.Version]*usdc_token_pool.USDCTokenPool)
			}
			state.USDCTokenPools[deployment.Version1_5_1] = utp
			state.ABIByAddress[address] = usdc_token_pool.USDCTokenPoolABI
		case deployment.NewTypeAndVersion(USDCMockTransmitter, deployment.Version1_0_0).String():
			umt, err := mock_usdc_token_transmitter.NewMockE2EUSDCTransmitter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.MockUSDCTransmitter = umt
			state.ABIByAddress[address] = mock_usdc_token_transmitter.MockE2EUSDCTransmitterABI
		case deployment.NewTypeAndVersion(USDCTokenMessenger, deployment.Version1_0_0).String():
			utm, err := mock_usdc_token_messenger.NewMockE2EUSDCTokenMessenger(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.MockUSDCTokenMessenger = utm
			state.ABIByAddress[address] = mock_usdc_token_messenger.MockE2EUSDCTokenMessengerABI
		case deployment.NewTypeAndVersion(CCIPHome, deployment.Version1_6_0).String():
			ccipHome, err := ccip_home.NewCCIPHome(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.CCIPHome = ccipHome
			state.ABIByAddress[address] = ccip_home.CCIPHomeABI
		case deployment.NewTypeAndVersion(CCIPReceiver, deployment.Version1_0_0).String():
			mr, err := maybe_revert_message_receiver.NewMaybeRevertMessageReceiver(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Receiver = mr
			state.ABIByAddress[address] = maybe_revert_message_receiver.MaybeRevertMessageReceiverABI
		case deployment.NewTypeAndVersion(LogMessageDataReceiver, deployment.Version1_0_0).String():
			mr, err := log_message_data_receiver.NewLogMessageDataReceiver(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.LogMessageDataReceiver = mr
			state.ABIByAddress[address] = log_message_data_receiver.LogMessageDataReceiverABI
		case deployment.NewTypeAndVersion(Multicall3, deployment.Version1_0_0).String():
			mc, err := multicall3.NewMulticall3(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Multicall3 = mc
			state.ABIByAddress[address] = multicall3.Multicall3ABI
		case deployment.NewTypeAndVersion(PriceFeed, deployment.Version1_0_0).String():
			feed, err := aggregator_v3_interface.NewAggregatorV3Interface(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.USDFeeds == nil {
				state.USDFeeds = make(map[TokenSymbol]*aggregator_v3_interface.AggregatorV3Interface)
			}
			desc, err := feed.Description(&bind.CallOpts{})
			if err != nil {
				return state, err
			}
			key, ok := DescriptionToTokenSymbol[desc]
			if !ok {
				return state, fmt.Errorf("unknown feed description %s", desc)
			}
			state.USDFeeds[key] = feed
			state.ABIByAddress[address] = aggregator_v3_interface.AggregatorV3InterfaceABI
		case deployment.NewTypeAndVersion(BurnMintTokenPool, deployment.Version1_5_1).String():
			ethAddress := common.HexToAddress(address)
			pool, metadata, err := NewTokenPoolWithMetadata(ctx, burn_mint_token_pool.NewBurnMintTokenPool, ethAddress, chain.Client)
			if err != nil {
				return state, fmt.Errorf("failed to connect address %s with token pool bindings and get token symbol: %w", ethAddress, err)
			}
			state.BurnMintTokenPools = helpers.AddValueToNestedMap(state.BurnMintTokenPools, metadata.Symbol, metadata.Version, pool)
			state.ABIByAddress[address] = burn_mint_token_pool.BurnMintTokenPoolABI
		case deployment.NewTypeAndVersion(BurnWithFromMintTokenPool, deployment.Version1_5_1).String():
			ethAddress := common.HexToAddress(address)
			pool, metadata, err := NewTokenPoolWithMetadata(ctx, burn_with_from_mint_token_pool.NewBurnWithFromMintTokenPool, ethAddress, chain.Client)
			if err != nil {
				return state, fmt.Errorf("failed to connect address %s with token pool bindings and get token symbol: %w", ethAddress, err)
			}
			state.BurnWithFromMintTokenPools = helpers.AddValueToNestedMap(state.BurnWithFromMintTokenPools, metadata.Symbol, metadata.Version, pool)
			state.ABIByAddress[address] = burn_with_from_mint_token_pool.BurnWithFromMintTokenPoolABI
		case deployment.NewTypeAndVersion(BurnFromMintTokenPool, deployment.Version1_5_1).String():
			ethAddress := common.HexToAddress(address)
			pool, metadata, err := NewTokenPoolWithMetadata(ctx, burn_from_mint_token_pool.NewBurnFromMintTokenPool, ethAddress, chain.Client)
			if err != nil {
				return state, fmt.Errorf("failed to connect address %s with token pool bindings and get token symbol: %w", ethAddress, err)
			}
			state.BurnFromMintTokenPools = helpers.AddValueToNestedMap(state.BurnFromMintTokenPools, metadata.Symbol, metadata.Version, pool)
			state.ABIByAddress[address] = burn_from_mint_token_pool.BurnFromMintTokenPoolABI
		case deployment.NewTypeAndVersion(LockReleaseTokenPool, deployment.Version1_5_1).String():
			ethAddress := common.HexToAddress(address)
			pool, metadata, err := NewTokenPoolWithMetadata(ctx, lock_release_token_pool.NewLockReleaseTokenPool, ethAddress, chain.Client)
			if err != nil {
				return state, fmt.Errorf("failed to connect address %s with token pool bindings and get token symbol: %w", ethAddress, err)
			}
			state.LockReleaseTokenPools = helpers.AddValueToNestedMap(state.LockReleaseTokenPools, metadata.Symbol, metadata.Version, pool)
			state.ABIByAddress[address] = lock_release_token_pool.LockReleaseTokenPoolABI
		case deployment.NewTypeAndVersion(BurnMintToken, deployment.Version1_0_0).String():
			tok, err := burn_mint_erc677.NewBurnMintERC677(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.BurnMintTokens677 == nil {
				state.BurnMintTokens677 = make(map[TokenSymbol]*burn_mint_erc677.BurnMintERC677)
			}
			symbol, err := tok.Symbol(nil)
			if err != nil {
				return state, fmt.Errorf("failed to get token symbol of token at %s: %w", address, err)
			}
			state.BurnMintTokens677[TokenSymbol(symbol)] = tok
			state.ABIByAddress[address] = burn_mint_erc677.BurnMintERC677ABI
		case deployment.NewTypeAndVersion(ERC20Token, deployment.Version1_0_0).String():
			tok, err := erc20.NewERC20(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.ERC20Tokens == nil {
				state.ERC20Tokens = make(map[TokenSymbol]*erc20.ERC20)
			}
			symbol, err := tok.Symbol(nil)
			if err != nil {
				return state, fmt.Errorf("failed to get token symbol of token at %s: %w", address, err)
			}
			state.ERC20Tokens[TokenSymbol(symbol)] = tok
			state.ABIByAddress[address] = erc20.ERC20ABI
		case deployment.NewTypeAndVersion(ERC677Token, deployment.Version1_0_0).String():
			tok, err := erc677.NewERC677(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.ERC677Tokens == nil {
				state.ERC677Tokens = make(map[TokenSymbol]*erc677.ERC677)
			}
			symbol, err := tok.Symbol(nil)
			if err != nil {
				return state, fmt.Errorf("failed to get token symbol of token at %s: %w", address, err)
			}
			state.ERC677Tokens[TokenSymbol(symbol)] = tok
			state.ABIByAddress[address] = erc677.ERC677ABI
		// legacy addresses below
		case deployment.NewTypeAndVersion(OnRamp, deployment.Version1_5_0).String():
			onRampC, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			sCfg, err := onRampC.GetStaticConfig(nil)
			if err != nil {
				return state, fmt.Errorf("failed to get static config chain %s: %w", chain.String(), err)
			}
			if state.EVM2EVMOnRamp == nil {
				state.EVM2EVMOnRamp = make(map[uint64]*evm_2_evm_onramp.EVM2EVMOnRamp)
			}
			state.EVM2EVMOnRamp[sCfg.DestChainSelector] = onRampC
			state.ABIByAddress[address] = evm_2_evm_onramp.EVM2EVMOnRampABI
		case deployment.NewTypeAndVersion(OffRamp, deployment.Version1_5_0).String():
			offRamp, err := evm_2_evm_offramp.NewEVM2EVMOffRamp(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			sCfg, err := offRamp.GetStaticConfig(nil)
			if err != nil {
				return state, err
			}
			if state.EVM2EVMOffRamp == nil {
				state.EVM2EVMOffRamp = make(map[uint64]*evm_2_evm_offramp.EVM2EVMOffRamp)
			}
			state.EVM2EVMOffRamp[sCfg.SourceChainSelector] = offRamp
			state.ABIByAddress[address] = evm_2_evm_offramp.EVM2EVMOffRampABI
		case deployment.NewTypeAndVersion(CommitStore, deployment.Version1_5_0).String():
			commitStore, err := commit_store.NewCommitStore(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			sCfg, err := commitStore.GetStaticConfig(nil)
			if err != nil {
				return state, err
			}
			if state.CommitStore == nil {
				state.CommitStore = make(map[uint64]*commit_store.CommitStore)
			}
			state.CommitStore[sCfg.SourceChainSelector] = commitStore
			state.ABIByAddress[address] = commit_store.CommitStoreABI
		case deployment.NewTypeAndVersion(PriceRegistry, deployment.Version1_2_0).String():
			pr, err := price_registry_1_2_0.NewPriceRegistry(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.PriceRegistry = pr
			state.ABIByAddress[address] = price_registry_1_2_0.PriceRegistryABI
		case deployment.NewTypeAndVersion(RMN, deployment.Version1_5_0).String():
			rmnC, err := rmn_contract.NewRMNContract(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMN = rmnC
			state.ABIByAddress[address] = rmn_contract.RMNContractABI
		case deployment.NewTypeAndVersion(MockRMN, deployment.Version1_0_0).String():
			mockRMN, err := mock_rmn_contract.NewMockRMNContract(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.MockRMN = mockRMN
			state.ABIByAddress[address] = mock_rmn_contract.MockRMNContractABI
		case deployment.NewTypeAndVersion(FiredrillEntrypointType, deployment.Version1_5_0).String(),
			deployment.NewTypeAndVersion(FiredrillEntrypointType, deployment.Version1_6_0).String():
			// Ignore firedrill contracts
			// Firedrill contracts are unknown to core and their state is being loaded separately
		default:
			// ManyChainMultiSig 1.0.0 can have any of these labels, it can have either 1,2 or 3 of these -
			// bypasser, proposer and canceller
			// if you try to compare tvStr.String() you will have to compare all combinations of labels
			// so we will compare the type and version only
			if tvStr.Type == commontypes.ManyChainMultisig && tvStr.Version == deployment.Version1_0_0 {
				state.ABIByAddress[address] = gethwrappers.ManyChainMultiSigABI
				continue
			}
			return state, fmt.Errorf("unknown contract %s", tvStr)
		}
	}
	return state, nil
}

func ValidateChain(env deployment.Environment, state CCIPOnChainState, chainSel uint64, mcmsCfg *MCMSConfig) error {
	err := deployment.IsValidChainSelector(chainSel)
	if err != nil {
		return fmt.Errorf("is not valid chain selector %d: %w", chainSel, err)
	}
	chain, ok := env.Chains[chainSel]
	if !ok {
		return fmt.Errorf("chain with selector %d does not exist in environment", chainSel)
	}
	chainState, ok := state.Chains[chainSel]
	if !ok {
		return fmt.Errorf("%s does not exist in state", chain)
	}
	if mcmsCfg != nil {
		if chainState.Timelock == nil {
			return fmt.Errorf("missing timelock on %s", chain)
		}
		if mcmsCfg.MCMSAction == timelock.Schedule && chainState.ProposerMcm == nil {
			return fmt.Errorf("missing proposerMcm on %s", chain)
		}
		if mcmsCfg.MCMSAction == timelock.Cancel && chainState.CancellerMcm == nil {
			return fmt.Errorf("missing cancellerMcm on %s", chain)
		}
		if mcmsCfg.MCMSAction == timelock.Bypass && chainState.BypasserMcm == nil {
			return fmt.Errorf("missing bypasserMcm on %s", chain)
		}
	}
	return nil
}

package ccip_integration_tests

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/arm_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_multi_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_multi_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"

	"github.com/stretchr/testify/require"
)

var (
	homeChainID = chainsel.GETH_TESTNET.EvmChainID
)

const (
	CapabilityLabelledName = "ccip"
	CapabilityVersion      = "v1.0.0"
)

func e18Mult(amount uint64) *big.Int {
	return new(big.Int).Mul(uintBigInt(amount), uintBigInt(1e18))
}

func uintBigInt(i uint64) *big.Int {
	return new(big.Int).SetUint64(i)
}

type homeChain struct {
	backend            *backends.SimulatedBackend
	owner              *bind.TransactOpts
	chainID            uint64
	capabilityRegistry *kcr.CapabilitiesRegistry
	ccipConfigContract common.Address
}

type onchainUniverse struct {
	backend            *backends.SimulatedBackend
	owner              *bind.TransactOpts
	chainID            uint64
	linkToken          *link_token.LinkToken
	weth               *weth9.WETH9
	router             *router.Router
	rmnProxy           *arm_proxy_contract.ARMProxyContract
	rmn                *mock_arm_contract.MockARMContract
	onramp             *evm_2_evm_multi_onramp.EVM2EVMMultiOnRamp
	offramp            *evm_2_evm_multi_offramp.EVM2EVMMultiOffRamp
	priceRegistry      *price_registry.PriceRegistry
	tokenAdminRegistry *token_admin_registry.TokenAdminRegistry
	nonceManager       *nonce_manager.NonceManager
}

type chainBase struct {
	backend *backends.SimulatedBackend
	owner   *bind.TransactOpts
}

// createUniverses does the following:
// 1. Creates 1 home chain and `numChains`-1 non-home chains
// 2. Sets up home chain with the capability registry and the CCIP config contract
// 2. Deploys the CCIP contracts to all chains.
// 3. Sets up the initial configurations for the contracts on all chains.
// 4. Wires the chains together.
func createUniverses(
	t *testing.T,
	numUniverses int,
) (homeChainUni homeChain, universes map[uint64]onchainUniverse) {
	chains := createChains(t, numUniverses)

	homeChainBase, ok := chains[homeChainID]
	require.True(t, ok, "home chain backend not available")
	// Set up home chain first
	homeChainUniverse := setupHomeChain(t, homeChainBase.owner, homeChainBase.backend)

	// deploy the ccip contracts on all chains
	universes = make(map[uint64]onchainUniverse)
	for chainID, base := range chains {
		owner := base.owner
		backend := base.backend
		// deploy the CCIP contracts
		linkToken := deployLinkToken(t, owner, backend, chainID)
		rmn := deployMockARMContract(t, owner, backend, chainID)
		rmnProxy := deployARMProxyContract(t, owner, backend, rmn.Address(), chainID)
		weth := deployWETHContract(t, owner, backend, chainID)
		rout := deployRouter(t, owner, backend, weth.Address(), rmnProxy.Address(), chainID)
		priceRegistry := deployPriceRegistry(t, owner, backend, linkToken.Address(), weth.Address(), chainID)
		tokenAdminRegistry := deployTokenAdminRegistry(t, owner, backend, chainID)
		nonceManager := deployNonceManager(t, owner, backend, chainID)

		//======================================================================
		//							OnRamp
		//======================================================================
		onRampAddr, _, _, err := evm_2_evm_multi_onramp.DeployEVM2EVMMultiOnRamp(
			owner,
			backend,
			evm_2_evm_multi_onramp.EVM2EVMMultiOnRampStaticConfig{
				LinkToken:          linkToken.Address(),
				ChainSelector:      chainID,
				RmnProxy:           rmnProxy.Address(),
				MaxFeeJuelsPerMsg:  big.NewInt(1e18),
				NonceManager:       nonceManager.Address(),
				TokenAdminRegistry: tokenAdminRegistry.Address(),
			},
			evm_2_evm_multi_onramp.EVM2EVMMultiOnRampDynamicConfig{
				Router:        rout.Address(),
				PriceRegistry: priceRegistry.Address(),
				//`withdrawFeeTokens` onRamp function is not part of the message flow
				// so we can set this to any address
				FeeAggregator: testutils.NewAddress(),
			},
			// Destination chain configs will be set up later once we have all chains
			[]evm_2_evm_multi_onramp.EVM2EVMMultiOnRampDestChainConfigArgs{},
			// PremiumMultiplier is always needed if the onramp is enabled
			[]evm_2_evm_multi_onramp.EVM2EVMMultiOnRampPremiumMultiplierWeiPerEthArgs{
				{
					PremiumMultiplierWeiPerEth: 9e17, //0.9 ETH
					Token:                      linkToken.Address(),
				},
				{
					PremiumMultiplierWeiPerEth: 1e18,
					Token:                      weth.Address(),
				},
			},
			//TODO: We'll need to have TransferFeeConfigArgs when we start testing with sending tokens
			[]evm_2_evm_multi_onramp.EVM2EVMMultiOnRampTokenTransferFeeConfigArgs{},
		)
		require.NoErrorf(t, err, "failed to deploy onramp on chain id %d", chainID)
		backend.Commit()
		onramp, err := evm_2_evm_multi_onramp.NewEVM2EVMMultiOnRamp(onRampAddr, backend)
		require.NoError(t, err)

		//======================================================================
		//							OffRamp
		//======================================================================
		offrampAddr, _, _, err := evm_2_evm_multi_offramp.DeployEVM2EVMMultiOffRamp(
			owner,
			backend,
			evm_2_evm_multi_offramp.EVM2EVMMultiOffRampStaticConfig{
				ChainSelector:      chainID,
				RmnProxy:           rmnProxy.Address(),
				TokenAdminRegistry: tokenAdminRegistry.Address(),
				NonceManager:       nonceManager.Address(),
			},
			// Source chain configs will be set up later once we have all chains
			[]evm_2_evm_multi_offramp.EVM2EVMMultiOffRampSourceChainConfigArgs{},
		)
		require.NoErrorf(t, err, "failed to deploy offramp on chain id %d", chainID)
		backend.Commit()
		offramp, err := evm_2_evm_multi_offramp.NewEVM2EVMMultiOffRamp(offrampAddr, backend)
		require.NoError(t, err)

		universe := onchainUniverse{
			backend:            backend,
			owner:              owner,
			chainID:            chainID,
			linkToken:          linkToken,
			weth:               weth,
			router:             rout,
			rmnProxy:           rmnProxy,
			rmn:                rmn,
			onramp:             onramp,
			offramp:            offramp,
			priceRegistry:      priceRegistry,
			tokenAdminRegistry: tokenAdminRegistry,
			nonceManager:       nonceManager,
		}
		// Set up the initial configurations for the contracts
		setupUniverseBasics(t, universe)

		universes[chainID] = universe
	}

	// Once we have all chains created and contracts deployed, we can set up the initial configurations and wire chains together
	connectUniverses(t, universes)

	return homeChainUniverse, universes
}

// Creates 1 home chain and `numChains`-1 non-home chains
func createChains(t *testing.T, numChains int) map[uint64]chainBase {
	chains := make(map[uint64]chainBase)

	homeChainOwner := testutils.MustNewSimTransactor(t)
	chains[homeChainID] = chainBase{
		owner: homeChainOwner,
		backend: backends.NewSimulatedBackend(core.GenesisAlloc{
			homeChainOwner.From: core.GenesisAccount{
				Balance: assets.Ether(10_000).ToInt(),
			},
		}, 30e6),
	}

	for chainID := chainsel.TEST_90000001.EvmChainID; len(chains) < numChains && chainID < chainsel.TEST_90000020.EvmChainID; chainID++ {
		owner := testutils.MustNewSimTransactor(t)
		chains[chainID] = chainBase{
			owner: owner,
			backend: backends.NewSimulatedBackend(core.GenesisAlloc{
				owner.From: core.GenesisAccount{
					Balance: assets.Ether(10_000).ToInt(),
				},
			}, 30e6),
		}
	}

	return chains
}

func setupHomeChain(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend) homeChain {
	// deploy the capability registry on the home chain
	crAddress, _, _, err := kcr.DeployCapabilitiesRegistry(owner, backend)
	require.NoError(t, err, "failed to deploy capability registry on home chain")
	backend.Commit()

	capabilityRegistry, err := kcr.NewCapabilitiesRegistry(crAddress, backend)
	require.NoError(t, err)

	ccAddress, _, _, err := ccip_config.DeployCCIPConfig(owner, backend, crAddress)
	require.NoError(t, err)
	backend.Commit()

	capabilityConfig, err := ccip_config.NewCCIPConfig(ccAddress, backend)
	require.NoError(t, err)

	_, err = capabilityRegistry.AddCapabilities(owner, []kcr.CapabilitiesRegistryCapability{
		{
			LabelledName:          CapabilityLabelledName,
			Version:               CapabilityVersion,
			CapabilityType:        2, // consensus. not used (?)
			ResponseType:          0, // report. not used (?)
			ConfigurationContract: ccAddress,
		},
	})
	require.NoError(t, err, "failed to add capabilities to the capability registry")
	backend.Commit()

	return homeChain{
		backend:            backend,
		owner:              owner,
		chainID:            homeChainID,
		capabilityRegistry: capabilityRegistry,
		ccipConfigContract: capabilityConfig.Address(),
	}
}

func connectUniverses(
	t *testing.T,
	universes map[uint64]onchainUniverse,
) {
	for _, uni := range universes {
		wireRouter(t, uni, universes)
		wireOnRamp(t, uni, universes)
		wireOffRamp(t, uni, universes)
		initRemoteChainsGasPrices(t, uni, universes)
	}
}

// setupUniverseBasics sets up the initial configurations for the CCIP contracts on a single chain.
// 1. Mint 1000 LINK to the owner
// 2. Set the price registry with local token prices
// 3. Authorize the onRamp and offRamp on the nonce manager
func setupUniverseBasics(t *testing.T, uni onchainUniverse) {
	//=============================================================================
	//			Universe specific  updates/configs
	//		These updates are specific to each universe and are set up here
	//      These updates don't depend on other chains
	//=============================================================================
	owner := uni.owner
	//=============================================================================
	//							Mint 1000 LINK to owner
	//=============================================================================
	_, err := uni.linkToken.GrantMintRole(owner, owner.From)
	require.NoError(t, err)
	_, err = uni.linkToken.Mint(owner, owner.From, e18Mult(1000))
	require.NoError(t, err)
	uni.backend.Commit()

	//=============================================================================
	//						Price updates for tokens
	//			These are the prices of the fee tokens of local chain in USD
	//=============================================================================
	tokenPriceUpdates := []price_registry.InternalTokenPriceUpdate{
		{
			SourceToken: uni.linkToken.Address(),
			UsdPerToken: e18Mult(20),
		},
		{
			SourceToken: uni.weth.Address(),
			UsdPerToken: e18Mult(4000),
		},
	}
	_, err = uni.priceRegistry.UpdatePrices(owner, price_registry.InternalPriceUpdates{
		TokenPriceUpdates: tokenPriceUpdates,
	})
	require.NoErrorf(t, err, "failed to apply price registry updates on chain id %d", uni.chainID)
	uni.backend.Commit()

	//=============================================================================
	//						Authorize OnRamp & OffRamp on NonceManager
	//	Otherwise the onramp will not be able to call the nonceManager to get next Nonce
	//=============================================================================
	authorizedCallersAuthorizedCallerArgs := nonce_manager.AuthorizedCallersAuthorizedCallerArgs{
		AddedCallers: []common.Address{
			uni.onramp.Address(),
			uni.offramp.Address(),
		},
	}
	_, err = uni.nonceManager.ApplyAuthorizedCallerUpdates(owner, authorizedCallersAuthorizedCallerArgs)
	require.NoError(t, err)
	uni.backend.Commit()
}

// As we can't change router contract. The contract was expecting onRamp and offRamp per lane and not per chain
// In the new architecture we have only one onRamp and one offRamp per chain.
// hence we add the mapping for all remote chains to the onRamp/offRamp contract of the local chain
func wireRouter(t *testing.T, uni onchainUniverse, universes map[uint64]onchainUniverse) {
	owner := uni.owner
	var (
		routerOnrampUpdates  []router.RouterOnRamp
		routerOfframpUpdates []router.RouterOffRamp
	)
	for remoteChainID := range universes {
		if remoteChainID == uni.chainID {
			continue
		}
		routerOnrampUpdates = append(routerOnrampUpdates, router.RouterOnRamp{
			DestChainSelector: remoteChainID,
			OnRamp:            uni.onramp.Address(),
		})
		routerOfframpUpdates = append(routerOfframpUpdates, router.RouterOffRamp{
			SourceChainSelector: remoteChainID,
			OffRamp:             uni.offramp.Address(),
		})
	}
	_, err := uni.router.ApplyRampUpdates(owner, routerOnrampUpdates, []router.RouterOffRamp{}, routerOfframpUpdates)
	require.NoErrorf(t, err, "failed to apply ramp updates on router on chain id %d", uni.chainID)
	uni.backend.Commit()
}

// Setting OnRampDestChainConfigs
func wireOnRamp(t *testing.T, uni onchainUniverse, universes map[uint64]onchainUniverse) {
	owner := uni.owner
	var onrampDestChainConfigArgs []evm_2_evm_multi_onramp.EVM2EVMMultiOnRampDestChainConfigArgs
	for remoteChainID := range universes {
		if remoteChainID == uni.chainID {
			continue
		}
		onrampDestChainConfigArgs = append(onrampDestChainConfigArgs, evm_2_evm_multi_onramp.EVM2EVMMultiOnRampDestChainConfigArgs{
			DestChainSelector: remoteChainID,
			DynamicConfig:     defaultOnRampDynamicConfig(t),
		})
	}
	_, err := uni.onramp.ApplyDestChainConfigUpdates(owner, onrampDestChainConfigArgs)
	require.NoErrorf(t, err, "failed to apply dest chain config updates on onramp on chain id %d", uni.chainID)
	uni.backend.Commit()
}

// Setting OffRampSourceChainConfigs
func wireOffRamp(t *testing.T, uni onchainUniverse, universes map[uint64]onchainUniverse) {
	owner := uni.owner
	var offrampSourceChainConfigArgs []evm_2_evm_multi_offramp.EVM2EVMMultiOffRampSourceChainConfigArgs
	for remoteChainID, remoteUniverse := range universes {
		if remoteChainID == uni.chainID {
			continue
		}
		offrampSourceChainConfigArgs = append(offrampSourceChainConfigArgs, evm_2_evm_multi_offramp.EVM2EVMMultiOffRampSourceChainConfigArgs{
			SourceChainSelector: remoteChainID, // for each destination chain, add a source chain config
			IsEnabled:           true,
			OnRamp:              remoteUniverse.onramp.Address().Bytes(),
		})
	}
	_, err := uni.offramp.ApplySourceChainConfigUpdates(owner, offrampSourceChainConfigArgs)
	require.NoErrorf(t, err, "failed to apply source chain config updates on offramp on chain id %d", uni.chainID)
	uni.backend.Commit()
}

// initRemoteChainsGasPrices sets the gas prices for all chains except the local chain in the local price registry
func initRemoteChainsGasPrices(t *testing.T, uni onchainUniverse, universes map[uint64]onchainUniverse) {
	var gasPriceUpdates []price_registry.InternalGasPriceUpdate
	for remoteChainID := range universes {
		if remoteChainID == uni.chainID {
			continue
		}
		gasPriceUpdates = append(gasPriceUpdates,
			price_registry.InternalGasPriceUpdate{
				DestChainSelector: remoteChainID,
				UsdPerUnitGas:     big.NewInt(2e12),
			},
		)
	}
	_, err := uni.priceRegistry.UpdatePrices(uni.owner, price_registry.InternalPriceUpdates{
		GasPriceUpdates: gasPriceUpdates,
	})
	require.NoError(t, err)
}

func defaultOnRampDynamicConfig(t *testing.T) evm_2_evm_multi_onramp.EVM2EVMMultiOnRampDestChainDynamicConfig {
	// https://github.com/smartcontractkit/ccip/blob/c4856b64bd766f1ddbaf5d13b42d3c4b12efde3a/contracts/src/v0.8/ccip/libraries/Internal.sol#L337-L337
	/*
		```Solidity
			// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
			bytes4 public constant CHAIN_FAMILY_SELECTOR_EVM = 0x2812d52c;
		```
	*/
	evmFamilySelector, err := hex.DecodeString("2812d52c")
	require.NoError(t, err)
	return evm_2_evm_multi_onramp.EVM2EVMMultiOnRampDestChainDynamicConfig{
		IsEnabled:                         true,
		MaxNumberOfTokensPerMsg:           10,
		MaxDataBytes:                      256,
		MaxPerMsgGasLimit:                 3_000_000,
		DestGasOverhead:                   50_000,
		DefaultTokenFeeUSDCents:           1,
		DestGasPerPayloadByte:             10,
		DestDataAvailabilityOverheadGas:   0,
		DestGasPerDataAvailabilityByte:    100,
		DestDataAvailabilityMultiplierBps: 1,
		DefaultTokenDestGasOverhead:       50_000,
		DefaultTokenDestBytesOverhead:     32,
		DefaultTxGasLimit:                 200_000,
		GasMultiplierWeiPerEth:            1,
		NetworkFeeUSDCents:                1,
		ChainFamilySelector:               [4]byte(evmFamilySelector),
	}
}

func deployLinkToken(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend, chainID uint64) *link_token.LinkToken {
	linkAddr, _, _, err := link_token.DeployLinkToken(owner, backend)
	require.NoErrorf(t, err, "failed to deploy link token on chain id %d", chainID)
	backend.Commit()
	linkToken, err := link_token.NewLinkToken(linkAddr, backend)
	require.NoError(t, err)
	return linkToken
}

func deployMockARMContract(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend, chainID uint64) *mock_arm_contract.MockARMContract {
	rmnAddr, _, _, err := mock_arm_contract.DeployMockARMContract(owner, backend)
	require.NoErrorf(t, err, "failed to deploy mock arm on chain id %d", chainID)
	backend.Commit()
	rmn, err := mock_arm_contract.NewMockARMContract(rmnAddr, backend)
	require.NoError(t, err)
	return rmn
}

func deployARMProxyContract(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend, rmnAddr common.Address, chainID uint64) *arm_proxy_contract.ARMProxyContract {
	rmnProxyAddr, _, _, err := arm_proxy_contract.DeployARMProxyContract(owner, backend, rmnAddr)
	require.NoErrorf(t, err, "failed to deploy arm proxy on chain id %d", chainID)
	backend.Commit()
	rmnProxy, err := arm_proxy_contract.NewARMProxyContract(rmnProxyAddr, backend)
	require.NoError(t, err)
	return rmnProxy
}

func deployWETHContract(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend, chainID uint64) *weth9.WETH9 {
	wethAddr, _, _, err := weth9.DeployWETH9(owner, backend)
	require.NoErrorf(t, err, "failed to deploy weth contract on chain id %d", chainID)
	backend.Commit()
	weth, err := weth9.NewWETH9(wethAddr, backend)
	require.NoError(t, err)
	return weth
}

func deployRouter(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend, wethAddr, rmnProxyAddr common.Address, chainID uint64) *router.Router {
	routerAddr, _, _, err := router.DeployRouter(owner, backend, wethAddr, rmnProxyAddr)
	require.NoErrorf(t, err, "failed to deploy router on chain id %d", chainID)
	backend.Commit()
	rout, err := router.NewRouter(routerAddr, backend)
	require.NoError(t, err)
	return rout
}

func deployPriceRegistry(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend, linkAddr, wethAddr common.Address, chainID uint64) *price_registry.PriceRegistry {
	priceRegistryAddr, _, _, err := price_registry.DeployPriceRegistry(owner, backend, []common.Address{}, []common.Address{linkAddr, wethAddr}, 24*60*60, []price_registry.PriceRegistryTokenPriceFeedUpdate{})
	require.NoErrorf(t, err, "failed to deploy price registry on chain id %d", chainID)
	backend.Commit()
	priceRegistry, err := price_registry.NewPriceRegistry(priceRegistryAddr, backend)
	require.NoError(t, err)
	return priceRegistry
}

func deployTokenAdminRegistry(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend, chainID uint64) *token_admin_registry.TokenAdminRegistry {
	tarAddr, _, _, err := token_admin_registry.DeployTokenAdminRegistry(owner, backend)
	require.NoErrorf(t, err, "failed to deploy token admin registry on chain id %d", chainID)
	backend.Commit()
	tokenAdminRegistry, err := token_admin_registry.NewTokenAdminRegistry(tarAddr, backend)
	require.NoError(t, err)
	return tokenAdminRegistry
}

func deployNonceManager(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend, chainID uint64) *nonce_manager.NonceManager {
	nonceManagerAddr, _, _, err := nonce_manager.DeployNonceManager(owner, backend, []common.Address{owner.From})
	require.NoErrorf(t, err, "failed to deploy nonce_manager on chain id %d", chainID)
	backend.Commit()
	nonceManager, err := nonce_manager.NewNonceManager(nonceManagerAddr, backend)
	require.NoError(t, err)
	return nonceManager
}

package crib

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/rs/zerolog"
	xerrgroup "golang.org/x/sync/errgroup"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/ethereum/go-ethereum/common"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// DeployHomeChainContracts deploys the home chain contracts so that the chainlink nodes can use the CR address in Capabilities.ExternalRegistry
// Afterward, we call DeployHomeChainChangeset changeset with nodeinfo ( the peer id and all)
func DeployHomeChainContracts(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel uint64, feedChainSel uint64) (deployment.CapabilityRegistryConfig, cldf.AddressBook, error) {
	e, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	if e == nil {
		return deployment.CapabilityRegistryConfig{}, nil, errors.New("environment is nil")
	}

	evmChains := e.BlockChains.EVMChains()

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("failed to get node info from env: %w", err)
	}
	p2pIds := nodes.NonBootstraps().PeerIDs()
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	for _, chain := range e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM)) {
		mcmsConfig, err := mcmstypes.NewConfig(1, []common.Address{evmChains[chain].DeployerKey.From}, []mcmstypes.Config{})
		if err != nil {
			return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("failed to create mcms config: %w", err)
		}
		cfg[chain] = commontypes.MCMSWithTimelockConfigV2{
			Canceller:        mcmsConfig,
			Bypasser:         mcmsConfig,
			Proposer:         mcmsConfig,
			TimelockMinDelay: big.NewInt(0),
		}
	}
	*e, err = commonchangeset.Apply(nil, *e, nil,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			cfg,
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
				HomeChainSel:             homeChainSel,
				RMNStaticConfig:          testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig:         testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:            testhelpers.NewTestNodeOperator(evmChains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{"NodeOperator": p2pIds},
			},
		),
	)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("changeset sequence execution failed with error: %w", err)
	}
	state, err := stateview.LoadOnchainState(*e)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("failed to load on chain state: %w", err)
	}
	capRegAddr := state.Chains[homeChainSel].CapabilityRegistry.Address()
	if capRegAddr == common.HexToAddress("0x") {
		return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("cap Reg address not found: %w", err)
	}
	capRegConfig := deployment.CapabilityRegistryConfig{
		EVMChainID:  homeChainSel,
		Contract:    state.Chains[homeChainSel].CapabilityRegistry.Address(),
		NetworkType: relay.NetworkEVM,
	}
	return capRegConfig, e.ExistingAddresses, nil
}

// DeployCCIPAndAddLanes is the actual ccip setup once the nodes are initialized.
func DeployCCIPAndAddLanes(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab cldf.AddressBook, rmnEnabled bool) (DeployCCIPOutput, error) {
	e, don, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to initiate new environment: %w", err)
	}
	e.ExistingAddresses = ab

	// ------ Part 1 -----
	// Setup because we only need to deploy the contracts and distribute job specs
	lggr.Infow("setting up chains...")
	*e, err = setupChains(lggr, e, homeChainSel)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply changesets for setting up chain: %w", err)
	}

	state, err := stateview.LoadOnchainState(*e)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	lggr.Infow("setting up lanes...")
	// Add all lanes
	*e, err = setupLanes(e, state)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply changesets for connecting lanes: %w", err)
	}
	// ------ Part 1 -----

	// ----- Part 2 -----
	lggr.Infow("setting up ocr...")
	*e, err = mustOCR(e, homeChainSel, feedChainSel, true, rmnEnabled)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply changesets for setting up OCR: %w", err)
	}

	// distribute funds to transmitters
	// we need to use the nodeinfo from the envConfig here, because multiAddr is not
	// populated in the environment variable
	lggr.Infow("distributing funds...")
	err = distributeTransmitterFunds(lggr, don.PluginNodes(), *e)
	if err != nil {
		return DeployCCIPOutput{}, err
	}

	addresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to convert address book to address book map: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *cldf.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, nil
}

// DeployCCIPChains is a group of changesets used from CRIB to set up new chains
// It sets up CCIP contracts on all chains. We expect that MCMS has already been deployed and set up
func DeployCCIPChains(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab cldf.AddressBook) (DeployCCIPOutput, error) {
	e, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to initiate new environment: %w", err)
	}
	e.ExistingAddresses = ab

	// Setup because we only need to deploy the contracts and distribute job specs
	lggr.Infow("setting up chains...")
	*e, err = setupChains(lggr, e, homeChainSel)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply changesets for setting up chain: %w", err)
	}
	addresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get convert address book to address book map: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *cldf.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, nil
}

// ConnectCCIPLanes is a group of changesets used from CRIB to set up new lanes
// It creates a fully connected mesh where all chains are connected to all chains
func ConnectCCIPLanes(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab cldf.AddressBook) (DeployCCIPOutput, error) {
	e, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to initiate new environment: %w", err)
	}
	e.ExistingAddresses = ab

	state, err := stateview.LoadOnchainState(*e)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	lggr.Infow("setting up lanes...")
	// Add all lanes
	*e, err = setupLanes(e, state)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply changesets for connecting lanes: %w", err)
	}

	addresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get convert address book to address book map: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *cldf.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, nil
}

// ConfigureCCIPOCR is a group of changesets used from CRIB to redeploy the chainlink don on an existing setup
func ConfigureCCIPOCR(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab cldf.AddressBook, rmnEnabled bool) (DeployCCIPOutput, error) {
	e, don, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to initiate new environment: %w", err)
	}
	e.ExistingAddresses = ab

	lggr.Infow("resetting ocr...")
	*e, err = mustOCR(e, homeChainSel, feedChainSel, false, rmnEnabled)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply changesets for setting up OCR: %w", err)
	}
	err = distributeTransmitterFunds(lggr, don.PluginNodes(), *e)
	if err != nil {
		return DeployCCIPOutput{}, err
	}

	addresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get convert address book to address book map: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *cldf.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, nil
}

// FundCCIPTransmitters is used from CRIB to provide funds to the node transmitters
// This function sends funds from the deployer key to the chainlink node transmitters
func FundCCIPTransmitters(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, ab cldf.AddressBook) (DeployCCIPOutput, error) {
	e, don, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to initiate new environment: %w", err)
	}
	e.ExistingAddresses = ab

	// distribute funds to transmitters
	// we need to use the nodeinfo from the envConfig here, because multiAddr is not
	// populated in the environment variable
	lggr.Infow("distributing funds...")
	err = distributeTransmitterFunds(lggr, don.PluginNodes(), *e)
	if err != nil {
		return DeployCCIPOutput{}, err
	}

	addresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get convert address book to address book map: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *cldf.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, nil
}

func setupChains(lggr logger.Logger, e *cldf.Environment, homeChainSel uint64) (cldf.Environment, error) {
	chainSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
	chainConfigs := make(map[uint64]v1_6.ChainConfig)
	nodeInfo, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return *e, fmt.Errorf("failed to get node info from env: %w", err)
	}
	prereqCfgs := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
	contractParams := make(map[uint64]ccipseq.ChainContractParams)

	for _, chain := range chainSelectors {
		prereqCfgs = append(prereqCfgs, changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
		chainConfigs[chain] = v1_6.ChainConfig{
			Readers: nodeInfo.NonBootstraps().PeerIDs(),
			// Number of nodes is 3f+1
			//nolint:gosec // this should always be less than max uint8
			FChain: uint8(len(nodeInfo.NonBootstraps().PeerIDs()) / 3),
			EncodableChainConfig: chainconfig.ChainConfig{
				GasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(1000)},
				DAGasPriceDeviationPPB:  cciptypes.BigInt{Int: big.NewInt(1_000_000)},
				OptimisticConfirmations: 1,
			},
		}
		contractParams[chain] = ccipseq.ChainContractParams{
			FeeQuoterParams: ccipops.DefaultFeeQuoterParams(),
			OffRampParams:   ccipops.DefaultOffRampParams(),
		}
	}
	env, err := commonchangeset.Apply(nil, *e, nil,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateChainConfigChangeset),
			v1_6.UpdateChainConfigConfig{
				HomeChainSelector: homeChainSel,
				RemoteChainAdds:   chainConfigs,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
			chainSelectors,
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
			changeset.DeployPrerequisiteConfig{
				Configs: prereqCfgs,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployChainContractsChangeset),
			ccipseq.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ContractParamsPerChain: contractParams,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.SetRMNRemoteOnRMNProxyChangeset),
			v1_6.SetRMNRemoteOnRMNProxyConfig{
				ChainSelectors: chainSelectors,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.CCIPCapabilityJobspecChangeset),
			nil, // ChangeSet does not use a config.
		),
	)
	if err != nil {
		return *e, fmt.Errorf("failed to apply changesets: %w", err)
	}
	lggr.Infow("setup Link pools")
	return setupLinkPools(&env)
}

func setupLinkPools(e *cldf.Environment) (cldf.Environment, error) {
	evmChains := e.BlockChains.EVMChains()
	state, err := stateview.LoadOnchainState(*e)
	if err != nil {
		return *e, fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
	poolInput := make(map[uint64]v1_5_1.DeployTokenPoolInput)
	pools := make(map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo)

	for _, chain := range chainSelectors {
		poolInput[chain] = v1_5_1.DeployTokenPoolInput{
			Type:               shared.BurnMintTokenPool,
			LocalTokenDecimals: 18,
			AllowList:          []common.Address{},
			TokenAddress:       state.Chains[chain].LinkToken.Address(),
		}
		pools[chain] = map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
			shared.LinkSymbol: {
				Type:          shared.BurnMintTokenPool,
				Version:       deployment.Version1_5_1,
				ExternalAdmin: evmChains[chain].DeployerKey.From,
			},
		}
	}
	env, err := commonchangeset.Apply(nil, *e, nil,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_5_1.DeployTokenPoolContractsChangeset),
			v1_5_1.DeployTokenPoolContractsConfig{
				TokenSymbol: shared.LinkSymbol,
				NewPools:    poolInput,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_5_1.ProposeAdminRoleChangeset),
			v1_5_1.TokenAdminRegistryChangesetConfig{
				Pools: pools,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_5_1.AcceptAdminRoleChangeset),
			v1_5_1.TokenAdminRegistryChangesetConfig{
				Pools: pools,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_5_1.SetPoolChangeset),
			v1_5_1.TokenAdminRegistryChangesetConfig{
				Pools: pools,
			},
		),
	)

	if err != nil {
		return *e, fmt.Errorf("failed to apply changesets: %w", err)
	}

	state, err = stateview.LoadOnchainState(env)
	if err != nil {
		return *e, fmt.Errorf("failed to load onchain state: %w", err)
	}

	for _, chain := range chainSelectors {
		linkPool := state.Chains[chain].BurnMintTokenPools[shared.LinkSymbol][deployment.Version1_5_1]
		linkToken := state.Chains[chain].LinkToken
		tx, err := linkToken.GrantMintAndBurnRoles(evmChains[chain].DeployerKey, linkPool.Address())
		_, err = cldf.ConfirmIfNoError(evmChains[chain], tx, err)
		if err != nil {
			return *e, fmt.Errorf("failed to grant mint and burn roles for link pool: %w", err)
		}
	}
	return env, err
}

func setupLanes(e *cldf.Environment, state stateview.CCIPOnChainState) (cldf.Environment, error) {
	eg := xerrgroup.Group{}
	poolUpdates := make(map[uint64]v1_5_1.TokenPoolConfig)
	rateLimitPerChain := make(v1_5_1.RateLimiterPerChain)
	mu := sync.Mutex{}
	evmChains := e.BlockChains.EVMChains()
	for src := range evmChains {
		src := src
		eg.Go(func() error {
			onRampUpdatesByChain := make(map[uint64]map[uint64]v1_6.OnRampDestinationUpdate)
			pricesByChain := make(map[uint64]v1_6.FeeQuoterPriceUpdatePerSource)
			feeQuoterDestsUpdatesByChain := make(map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig)
			updateOffRampSources := make(map[uint64]map[uint64]v1_6.OffRampSourceUpdate)
			updateRouterChanges := make(map[uint64]v1_6.RouterUpdates)
			onRampUpdatesByChain[src] = make(map[uint64]v1_6.OnRampDestinationUpdate)
			pricesByChain[src] = v1_6.FeeQuoterPriceUpdatePerSource{
				TokenPrices: map[common.Address]*big.Int{
					state.Chains[src].LinkToken.Address(): testhelpers.DefaultLinkPrice,
					state.Chains[src].Weth9.Address():     testhelpers.DefaultWethPrice,
				},
				GasPrices: make(map[uint64]*big.Int),
			}
			feeQuoterDestsUpdatesByChain[src] = make(map[uint64]fee_quoter.FeeQuoterDestChainConfig)
			updateOffRampSources[src] = make(map[uint64]v1_6.OffRampSourceUpdate)
			updateRouterChanges[src] = v1_6.RouterUpdates{
				OffRampUpdates: make(map[uint64]bool),
				OnRampUpdates:  make(map[uint64]bool),
			}

			for dst := range evmChains {
				if src != dst {
					onRampUpdatesByChain[src][dst] = v1_6.OnRampDestinationUpdate{
						IsEnabled:        true,
						AllowListEnabled: false,
					}
					pricesByChain[src].GasPrices[dst] = testhelpers.DefaultGasPrice
					feeQuoterDestsUpdatesByChain[src][dst] = v1_6.DefaultFeeQuoterDestChainConfig(true)

					updateOffRampSources[src][dst] = v1_6.OffRampSourceUpdate{
						IsEnabled:                 true,
						IsRMNVerificationDisabled: true,
					}

					updateRouterChanges[src].OffRampUpdates[dst] = true
					updateRouterChanges[src].OnRampUpdates[dst] = true
					mu.Lock()
					rateLimitPerChain[dst] = v1_5_1.RateLimiterConfig{
						Inbound: token_pool.RateLimiterConfig{
							IsEnabled: false,
							Capacity:  big.NewInt(0),
							Rate:      big.NewInt(0),
						},
						Outbound: token_pool.RateLimiterConfig{
							IsEnabled: false,
							Capacity:  big.NewInt(0),
							Rate:      big.NewInt(0),
						},
					}
					mu.Unlock()
				}
			}
			mu.Lock()
			poolUpdates[src] = v1_5_1.TokenPoolConfig{
				Type:         shared.BurnMintTokenPool,
				Version:      deployment.Version1_5_1,
				ChainUpdates: rateLimitPerChain,
			}
			mu.Unlock()

			_, err := commonchangeset.Apply(nil, *e, nil,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateOnRampsDestsChangeset),
					v1_6.UpdateOnRampDestsConfig{
						UpdatesByChain: onRampUpdatesByChain,
					},
				),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterPricesChangeset),
					v1_6.UpdateFeeQuoterPricesConfig{
						PricesByChain: pricesByChain,
					},
				),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterDestsChangeset),
					v1_6.UpdateFeeQuoterDestsConfig{
						UpdatesByChain: feeQuoterDestsUpdatesByChain,
					},
				),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateOffRampSourcesChangeset),
					v1_6.UpdateOffRampSourcesConfig{
						UpdatesByChain: updateOffRampSources,
					},
				),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateRouterRampsChangeset),
					v1_6.UpdateRouterRampsConfig{
						UpdatesByChain: updateRouterChanges,
					},
				),
			)
			return err
		})
	}

	err := eg.Wait()
	if err != nil {
		return *e, err
	}

	_, err = commonchangeset.Apply(nil, *e, nil, commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
		v1_5_1.ConfigureTokenPoolContractsConfig{
			TokenSymbol: shared.LinkSymbol,
			PoolUpdates: poolUpdates,
		},
	))

	return *e, err
}

func mustOCR(e *cldf.Environment, homeChainSel uint64, feedChainSel uint64, newDons bool, rmnEnabled bool) (cldf.Environment, error) {
	chainSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
	var commitOCRConfigPerSelector = make(map[uint64]v1_6.CCIPOCRParams)
	var execOCRConfigPerSelector = make(map[uint64]v1_6.CCIPOCRParams)
	// Should be configured in the future based on the load test scenario
	chainType := v1_6.Default

	overrides := func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams { return params }
	if rmnEnabled {
		overrides = func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
			params.CommitOffChainConfig.RMNEnabled = true
			return params
		}
	}

	for selector := range e.BlockChains.EVMChains() {
		commitOCRConfigPerSelector[selector] = v1_6.DeriveOCRParamsForCommit(chainType, feedChainSel, nil, overrides)
		execOCRConfigPerSelector[selector] = v1_6.DeriveOCRParamsForExec(chainType, nil, nil)
	}

	var commitChangeset commonchangeset.ConfiguredChangeSet
	if newDons {
		commitChangeset = commonchangeset.Configure(
			// Add the DONs and candidate commit OCR instances for the chain
			cldf.CreateLegacyChangeSet(v1_6.AddDonAndSetCandidateChangeset),
			v1_6.AddDonAndSetCandidateChangesetConfig{
				SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
					HomeChainSelector: homeChainSel,
					FeedChainSelector: feedChainSel,
				},
				PluginInfo: v1_6.SetCandidatePluginInfo{
					OCRConfigPerRemoteChainSelector: commitOCRConfigPerSelector,
					PluginType:                      types.PluginTypeCCIPCommit,
				},
			},
		)
	} else {
		commitChangeset = commonchangeset.Configure(
			// Update commit OCR instances for existing chains
			cldf.CreateLegacyChangeSet(v1_6.SetCandidateChangeset),
			v1_6.SetCandidateChangesetConfig{
				SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
					HomeChainSelector: homeChainSel,
					FeedChainSelector: feedChainSel,
				},
				PluginInfo: []v1_6.SetCandidatePluginInfo{
					{
						OCRConfigPerRemoteChainSelector: commitOCRConfigPerSelector,
						PluginType:                      types.PluginTypeCCIPCommit,
					},
				},
			},
		)
	}

	return commonchangeset.Apply(nil, *e, nil,
		commitChangeset,
		commonchangeset.Configure(
			// Add the exec OCR instances for the new chains
			cldf.CreateLegacyChangeSet(v1_6.SetCandidateChangeset),
			v1_6.SetCandidateChangesetConfig{
				SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
					HomeChainSelector: homeChainSel,
					FeedChainSelector: feedChainSel,
				},
				PluginInfo: []v1_6.SetCandidatePluginInfo{
					{
						OCRConfigPerRemoteChainSelector: execOCRConfigPerSelector,
						PluginType:                      types.PluginTypeCCIPExec,
					},
				},
			},
		),
		commonchangeset.Configure(
			// Promote everything
			cldf.CreateLegacyChangeSet(v1_6.PromoteCandidateChangeset),
			v1_6.PromoteCandidateChangesetConfig{
				HomeChainSelector: homeChainSel,
				PluginInfo: []v1_6.PromoteCandidatePluginInfo{
					{
						RemoteChainSelectors: chainSelectors,
						PluginType:           types.PluginTypeCCIPCommit,
					},
					{
						RemoteChainSelectors: chainSelectors,
						PluginType:           types.PluginTypeCCIPExec,
					},
				},
			},
		),
		commonchangeset.Configure(
			// Enable the OCR config on the remote chains
			cldf.CreateLegacyChangeSet(v1_6.SetOCR3OffRampChangeset),
			v1_6.SetOCR3OffRampConfig{
				HomeChainSel:       homeChainSel,
				RemoteChainSels:    chainSelectors,
				CCIPHomeConfigType: globals.ConfigTypeActive,
			},
		),
	)
}

type RMNNodeConfig struct {
	v1_6.RMNNopConfig
	RageProxyKeystore string
	RMNKeystore       string
	Passphrase        string
}

func SetupRMNNodeOnAllChains(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab cldf.AddressBook, nodes []RMNNodeConfig) (DeployCCIPOutput, error) {
	e, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to create environment: %w", err)
	}

	e.ExistingAddresses = ab

	allChains := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
	allUpdates := make(map[uint64]map[uint64]v1_6.OffRampSourceUpdate)
	for _, chainIdx := range allChains {
		updates := make(map[uint64]v1_6.OffRampSourceUpdate)

		for _, subChainID := range allChains {
			if subChainID == chainIdx {
				continue
			}
			updates[subChainID] = v1_6.OffRampSourceUpdate{
				IsRMNVerificationDisabled: false,
				IsEnabled:                 true,
			}
		}

		allUpdates[chainIdx] = updates
	}

	_, err = v1_6.UpdateOffRampSourcesChangeset(*e, v1_6.UpdateOffRampSourcesConfig{
		UpdatesByChain: allUpdates,
	})

	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to update dynamic off ramp config: %w", err)
	}

	rmnNodes := make([]rmn_home.RMNHomeNode, len(nodes))
	bitmap := new(big.Int)
	for i, node := range nodes {
		rmnNodes[i] = rmn_home.RMNHomeNode{
			PeerId:            node.PeerId,
			OffchainPublicKey: node.OffchainPublicKey,
		}
		bitmap.SetBit(bitmap, i, 1)
	}

	sourceChains := make([]rmn_home.RMNHomeSourceChain, len(allChains))
	for i, chain := range allChains {
		sourceChains[i] = rmn_home.RMNHomeSourceChain{
			ChainSelector:       chain,
			FObserve:            1,
			ObserverNodesBitmap: bitmap,
		}
	}

	env, err := commonchangeset.Apply(nil, *e, nil,
		commonchangeset.Configure(
			// Enable the OCR config on the remote chains
			cldf.CreateLegacyChangeSet(v1_6.SetRMNHomeCandidateConfigChangeset),
			v1_6.SetRMNHomeCandidateConfig{
				HomeChainSelector: homeChainSel,
				RMNStaticConfig: rmn_home.RMNHomeStaticConfig{
					Nodes:          rmnNodes,
					OffchainConfig: []byte{},
				},
				RMNDynamicConfig: rmn_home.RMNHomeDynamicConfig{
					OffchainConfig: []byte{},
					SourceChains:   sourceChains,
				},
			},
		),
	)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to set rmn node candidate: %w", err)
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to load chain state: %w", err)
	}

	configDigest, err := state.Chains[homeChainSel].RMNHome.GetCandidateDigest(nil)

	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get rmn home candidate digest: %w", err)
	}

	env, err = commonchangeset.Apply(nil, *e, nil,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.PromoteRMNHomeCandidateConfigChangeset),
			v1_6.PromoteRMNHomeCandidateConfig{
				HomeChainSelector: homeChainSel,
				DigestToPromote:   configDigest,
			},
		),
	)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to promote rmn node candidate: %w", err)
	}

	signers := make([]rmn_remote.RMNRemoteSigner, len(nodes))
	for i, node := range nodes {
		signers[i] = node.ToRMNRemoteSigner()
	}

	g, ctx := xerrgroup.WithContext(context.Background())
	for _, chain := range allChains {
		g.Go(func() error {
			rmnRemoteConfig := map[uint64]ccipops.RMNRemoteConfig{
				chain: {
					Signers: signers,
					F:       1,
				},
			}

			_, err := v1_6.SetRMNRemoteConfigChangeset(*e, ccipseq.SetRMNRemoteConfig{
				RMNRemoteConfigs: rmnRemoteConfig,
			})
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to set rmn remote config: %w", err)
	}

	addresses, err := env.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *cldf.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, nil
}

func GenerateRMNNodeIdentities(rmnNodeCount uint, rageProxyImageURI, rageProxyImageTag, afn2proxyImageURI,
	afn2proxyImageTag string, imagePlatform string) ([]RMNNodeConfig, error) {
	lggr := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	rmnNodeConfigs := make([]RMNNodeConfig, rmnNodeCount)

	for i := uint(0); i < rmnNodeCount; i++ {
		peerID, rawKeystore, _, err := devenv.GeneratePeerID(zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}), rageProxyImageURI, rageProxyImageTag, imagePlatform)
		if err != nil {
			return nil, err
		}

		keys, rawRMNKeystore, afnPassphrase, err := devenv.GenerateRMNKeyStore(lggr, afn2proxyImageURI, afn2proxyImageTag, imagePlatform)
		if err != nil {
			return nil, err
		}

		newPeerID, err := p2pkey.MakePeerID(peerID.String())
		if err != nil {
			return nil, err
		}

		rmnNodeConfigs[i] = RMNNodeConfig{
			RMNNopConfig: v1_6.RMNNopConfig{
				NodeIndex:           uint64(i),
				OffchainPublicKey:   [32]byte(keys.OffchainPublicKey),
				EVMOnChainPublicKey: keys.EVMOnchainPublicKey,
				PeerId:              newPeerID,
			},
			RageProxyKeystore: rawKeystore,
			RMNKeystore:       rawRMNKeystore,
			Passphrase:        afnPassphrase,
		}
	}
	return rmnNodeConfigs, nil
}

package crib

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_1/token_pool"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// DeployHomeChainContracts deploys the home chain contracts so that the chainlink nodes can use the CR address in Capabilities.ExternalRegistry
// Afterward, we call DeployHomeChainChangeset changeset with nodeinfo ( the peer id and all)
func DeployHomeChainContracts(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel uint64, feedChainSel uint64) (deployment.CapabilityRegistryConfig, deployment.AddressBook, error) {
	e, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	if e == nil {
		return deployment.CapabilityRegistryConfig{}, nil, errors.New("environment is nil")
	}

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("failed to get node info from env: %w", err)
	}
	p2pIds := nodes.NonBootstraps().PeerIDs()
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	for _, chain := range e.AllChainSelectors() {
		mcmsConfig, err := config.NewConfig(1, []common.Address{e.Chains[chain].DeployerKey.From}, []config.Config{})
		if err != nil {
			return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("failed to create mcms config: %w", err)
		}
		cfg[chain] = commontypes.MCMSWithTimelockConfig{
			Canceller:        *mcmsConfig,
			Bypasser:         *mcmsConfig,
			Proposer:         *mcmsConfig,
			TimelockMinDelay: big.NewInt(0),
		}
	}
	*e, err = commonchangeset.Apply(nil, *e, nil,
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelock),
			cfg,
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.DeployHomeChainChangeset),
			changeset.DeployHomeChainConfig{
				HomeChainSel:             homeChainSel,
				RMNStaticConfig:          testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig:         testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:            testhelpers.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{"NodeOperator": p2pIds},
			},
		),
	)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("changeset sequence execution failed with error: %w", err)
	}
	state, err := changeset.LoadOnchainState(*e)
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
func DeployCCIPAndAddLanes(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab deployment.AddressBook) (DeployCCIPOutput, error) {
	e, don, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to initiate new environment: %w", err)
	}
	e.ExistingAddresses = ab

	// ------ Part 1 -----
	// Setup because we only need to deploy the contracts and distribute job specs
	fmt.Println("setting up chains...")
	*e, err = setupChains(lggr, e, homeChainSel)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply changesets for setting up chain: %w", err)
	}

	state, err := changeset.LoadOnchainState(*e)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	fmt.Println("setting up lanes...")
	// Add all lanes
	*e, err = setupLanes(e, state)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply changesets for connecting lanes: %w", err)
	}
	// ------ Part 1 -----

	// ----- Part 2 -----
	fmt.Println("setting up ocr...")
	*e, err = setupOCR(e, homeChainSel, feedChainSel)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply changesets for setting up OCR: %w", err)
	}

	// distribute funds to transmitters
	// we need to use the nodeinfo from the envConfig here, because multiAddr is not
	// populated in the environment variable
	fmt.Println("distributing funds...")
	err = distributeTransmitterFunds(lggr, don.PluginNodes(), *e)
	if err != nil {
		return DeployCCIPOutput{}, err
	}

	addresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to convert address book to address book map: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *deployment.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, nil
}

// DeployCCIPChains is a group of changesets used from CRIB to set up new chains
// It sets up CCIP contracts on all chains. We expect that MCMS has already been deployed and set up
func DeployCCIPChains(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab deployment.AddressBook) (DeployCCIPOutput, error) {
	e, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to initiate new environment: %w", err)
	}
	e.ExistingAddresses = ab

	// Setup because we only need to deploy the contracts and distribute job specs
	fmt.Println("setting up chains...")
	*e, err = setupChains(lggr, e, homeChainSel)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply changesets for setting up chain: %w", err)
	}
	addresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get convert address book to address book map: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *deployment.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, nil
}

// ConnectCCIPLanes is a group of changesets used from CRIB to set up new lanes
// It creates a fully connected mesh where all chains are connected to all chains
func ConnectCCIPLanes(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab deployment.AddressBook) (DeployCCIPOutput, error) {
	e, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to initiate new environment: %w", err)
	}
	e.ExistingAddresses = ab

	state, err := changeset.LoadOnchainState(*e)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	fmt.Println("setting up lanes...")
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
		AddressBook: *deployment.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, nil
}

// ConfigureCCIPOCR is a group of changesets used from CRIB to configure OCR on a new setup
// This sets up OCR on all chains in the envConfig by configuring the CCIP home chain
func ConfigureCCIPOCR(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab deployment.AddressBook) (DeployCCIPOutput, error) {
	e, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to initiate new environment: %w", err)
	}
	e.ExistingAddresses = ab

	fmt.Println("setting up ocr...")
	*e, err = setupOCR(e, homeChainSel, feedChainSel)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply changesets for setting up OCR: %w", err)
	}

	addresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get convert address book to address book map: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *deployment.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, nil
}

// FundCCIPTransmitters is used from CRIB to provide funds to the node transmitters
// This function sends funds from the deployer key to the chainlink node transmitters
func FundCCIPTransmitters(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, ab deployment.AddressBook) (DeployCCIPOutput, error) {
	e, don, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to initiate new environment: %w", err)
	}
	e.ExistingAddresses = ab

	// distribute funds to transmitters
	// we need to use the nodeinfo from the envConfig here, because multiAddr is not
	// populated in the environment variable
	fmt.Println("distributing funds...")
	err = distributeTransmitterFunds(lggr, don.PluginNodes(), *e)
	if err != nil {
		return DeployCCIPOutput{}, err
	}

	addresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get convert address book to address book map: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *deployment.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, nil
}

func setupChains(lggr logger.Logger, e *deployment.Environment, homeChainSel uint64) (deployment.Environment, error) {
	chainSelectors := e.AllChainSelectors()
	chainConfigs := make(map[uint64]changeset.ChainConfig)
	nodeInfo, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return *e, fmt.Errorf("failed to get node info from env: %w", err)
	}
	prereqCfgs := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
	contractParams := make(map[uint64]changeset.ChainContractParams)

	for _, chain := range chainSelectors {
		prereqCfgs = append(prereqCfgs, changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
		chainConfigs[chain] = changeset.ChainConfig{
			Readers: nodeInfo.NonBootstraps().PeerIDs(),
			FChain:  1,
			EncodableChainConfig: chainconfig.ChainConfig{
				GasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(1000)},
				DAGasPriceDeviationPPB:  cciptypes.BigInt{Int: big.NewInt(1_000_000)},
				OptimisticConfirmations: 1,
			},
		}
		contractParams[chain] = changeset.ChainContractParams{
			FeeQuoterParams: changeset.DefaultFeeQuoterParams(),
			OffRampParams:   changeset.DefaultOffRampParams(),
		}
	}
	env, err := commonchangeset.Apply(nil, *e, nil,
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.UpdateChainConfigChangeset),
			changeset.UpdateChainConfigConfig{
				HomeChainSelector: homeChainSel,
				RemoteChainAdds:   chainConfigs,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
			chainSelectors,
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
			changeset.DeployPrerequisiteConfig{
				Configs: prereqCfgs,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.DeployChainContractsChangeset),
			changeset.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ContractParamsPerChain: contractParams,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.SetRMNRemoteOnRMNProxyChangeset),
			changeset.SetRMNRemoteOnRMNProxyConfig{
				ChainSelectors: chainSelectors,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.CCIPCapabilityJobspecChangeset),
			nil, // ChangeSet does not use a config.
		),
	)
	if err != nil {
		return *e, fmt.Errorf("failed to apply changesets: %w", err)
	}
	lggr.Infow("setup Link pools")
	return setupLinkPools(&env)
}

func setupLinkPools(e *deployment.Environment) (deployment.Environment, error) {
	state, err := changeset.LoadOnchainState(*e)
	if err != nil {
		return *e, fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainSelectors := e.AllChainSelectors()
	poolInput := make(map[uint64]changeset.DeployTokenPoolInput)
	pools := make(map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo)
	for _, chain := range chainSelectors {
		poolInput[chain] = changeset.DeployTokenPoolInput{
			Type:               changeset.BurnMintTokenPool,
			LocalTokenDecimals: 18,
			AllowList:          []common.Address{},
			TokenAddress:       state.Chains[chain].LinkToken.Address(),
		}
		pools[chain] = map[changeset.TokenSymbol]changeset.TokenPoolInfo{
			changeset.LinkSymbol: {
				Type:          changeset.BurnMintTokenPool,
				Version:       deployment.Version1_5_1,
				ExternalAdmin: e.Chains[chain].DeployerKey.From,
			},
		}
	}
	env, err := commonchangeset.Apply(nil, *e, nil,
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.DeployTokenPoolContractsChangeset),
			changeset.DeployTokenPoolContractsConfig{
				TokenSymbol: changeset.LinkSymbol,
				NewPools:    poolInput,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.ProposeAdminRoleChangeset),
			changeset.TokenAdminRegistryChangesetConfig{
				Pools: pools,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.AcceptAdminRoleChangeset),
			changeset.TokenAdminRegistryChangesetConfig{
				Pools: pools,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.SetPoolChangeset),
			changeset.TokenAdminRegistryChangesetConfig{
				Pools: pools,
			},
		),
	)

	if err != nil {
		return *e, fmt.Errorf("failed to apply changesets: %w", err)
	}

	state, err = changeset.LoadOnchainState(env)
	if err != nil {
		return *e, fmt.Errorf("failed to load onchain state: %w", err)
	}

	for _, chain := range chainSelectors {
		linkPool := state.Chains[chain].BurnMintTokenPools[changeset.LinkSymbol][deployment.Version1_5_1]
		linkToken := state.Chains[chain].LinkToken
		tx, err := linkToken.GrantMintAndBurnRoles(e.Chains[chain].DeployerKey, linkPool.Address())
		_, err = deployment.ConfirmIfNoError(e.Chains[chain], tx, err)
		if err != nil {
			return *e, fmt.Errorf("failed to grant mint and burn roles for link pool: %w", err)
		}
	}
	return env, err
}

func setupLanes(e *deployment.Environment, state changeset.CCIPOnChainState) (deployment.Environment, error) {
	onRampUpdatesByChain := make(map[uint64]map[uint64]changeset.OnRampDestinationUpdate)
	pricesByChain := make(map[uint64]changeset.FeeQuoterPriceUpdatePerSource)
	feeQuoterDestsUpdatesByChain := make(map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig)
	updateOffRampSources := make(map[uint64]map[uint64]changeset.OffRampSourceUpdate)
	updateRouterChanges := make(map[uint64]changeset.RouterUpdates)
	poolUpdates := make(map[uint64]changeset.TokenPoolConfig)
	for src := range e.Chains {
		onRampUpdatesByChain[src] = make(map[uint64]changeset.OnRampDestinationUpdate)
		pricesByChain[src] = changeset.FeeQuoterPriceUpdatePerSource{
			TokenPrices: map[common.Address]*big.Int{
				state.Chains[src].LinkToken.Address(): testhelpers.DefaultLinkPrice,
				state.Chains[src].Weth9.Address():     testhelpers.DefaultWethPrice,
			},
			GasPrices: make(map[uint64]*big.Int),
		}
		feeQuoterDestsUpdatesByChain[src] = make(map[uint64]fee_quoter.FeeQuoterDestChainConfig)
		updateOffRampSources[src] = make(map[uint64]changeset.OffRampSourceUpdate)
		updateRouterChanges[src] = changeset.RouterUpdates{
			OffRampUpdates: make(map[uint64]bool),
			OnRampUpdates:  make(map[uint64]bool),
		}
		rateLimitPerChain := make(changeset.RateLimiterPerChain)

		for dst := range e.Chains {
			if src != dst {
				onRampUpdatesByChain[src][dst] = changeset.OnRampDestinationUpdate{
					IsEnabled:        true,
					AllowListEnabled: false,
				}
				pricesByChain[src].GasPrices[dst] = testhelpers.DefaultGasPrice
				feeQuoterDestsUpdatesByChain[src][dst] = changeset.DefaultFeeQuoterDestChainConfig(true)

				updateOffRampSources[src][dst] = changeset.OffRampSourceUpdate{
					IsEnabled:                 true,
					IsRMNVerificationDisabled: true,
				}

				updateRouterChanges[src].OffRampUpdates[dst] = true
				updateRouterChanges[src].OnRampUpdates[dst] = true
				rateLimitPerChain[dst] = changeset.RateLimiterConfig{
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
			}
		}

		poolUpdates[src] = changeset.TokenPoolConfig{
			Type:         changeset.BurnMintTokenPool,
			Version:      deployment.Version1_5_1,
			ChainUpdates: rateLimitPerChain,
		}
	}

	return commonchangeset.Apply(nil, *e, nil,
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.ConfigureTokenPoolContractsChangeset),
			changeset.ConfigureTokenPoolContractsConfig{
				TokenSymbol: changeset.LinkSymbol,
				PoolUpdates: poolUpdates,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.UpdateOnRampsDestsChangeset),
			changeset.UpdateOnRampDestsConfig{
				UpdatesByChain: onRampUpdatesByChain,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.UpdateFeeQuoterPricesChangeset),
			changeset.UpdateFeeQuoterPricesConfig{
				PricesByChain: pricesByChain,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.UpdateFeeQuoterDestsChangeset),
			changeset.UpdateFeeQuoterDestsConfig{
				UpdatesByChain: feeQuoterDestsUpdatesByChain,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.UpdateOffRampSourcesChangeset),
			changeset.UpdateOffRampSourcesConfig{
				UpdatesByChain: updateOffRampSources,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.UpdateRouterRampsChangeset),
			changeset.UpdateRouterRampsConfig{
				UpdatesByChain: updateRouterChanges,
			},
		),
	)
}

func setupOCR(e *deployment.Environment, homeChainSel uint64, feedChainSel uint64) (deployment.Environment, error) {
	chainSelectors := e.AllChainSelectors()
	var ocrConfigPerSelector = make(map[uint64]changeset.CCIPOCRParams)
	for selector := range e.Chains {
		ocrConfigPerSelector[selector] = changeset.DeriveCCIPOCRParams(changeset.WithDefaultCommitOffChainConfig(feedChainSel, nil),
			changeset.WithDefaultExecuteOffChainConfig(nil),
		)
	}
	return commonchangeset.Apply(nil, *e, nil,
		commonchangeset.Configure(
			// Add the DONs and candidate commit OCR instances for the chain
			deployment.CreateLegacyChangeSet(changeset.AddDonAndSetCandidateChangeset),
			changeset.AddDonAndSetCandidateChangesetConfig{
				SetCandidateConfigBase: changeset.SetCandidateConfigBase{
					HomeChainSelector: homeChainSel,
					FeedChainSelector: feedChainSel,
				},
				PluginInfo: changeset.SetCandidatePluginInfo{
					OCRConfigPerRemoteChainSelector: ocrConfigPerSelector,
					PluginType:                      types.PluginTypeCCIPCommit,
				},
			},
		),
		commonchangeset.Configure(
			// Add the exec OCR instances for the new chains
			deployment.CreateLegacyChangeSet(changeset.SetCandidateChangeset),
			changeset.SetCandidateChangesetConfig{
				SetCandidateConfigBase: changeset.SetCandidateConfigBase{
					HomeChainSelector: homeChainSel,
					FeedChainSelector: feedChainSel,
				},
				PluginInfo: []changeset.SetCandidatePluginInfo{
					{
						OCRConfigPerRemoteChainSelector: ocrConfigPerSelector,
						PluginType:                      types.PluginTypeCCIPExec,
					},
				},
			},
		),
		commonchangeset.Configure(
			// Promote everything
			deployment.CreateLegacyChangeSet(changeset.PromoteCandidateChangeset),
			changeset.PromoteCandidateChangesetConfig{
				HomeChainSelector: homeChainSel,
				PluginInfo: []changeset.PromoteCandidatePluginInfo{
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
			deployment.CreateLegacyChangeSet(changeset.SetOCR3OffRampChangeset),
			changeset.SetOCR3OffRampConfig{
				HomeChainSel:       homeChainSel,
				RemoteChainSels:    chainSelectors,
				CCIPHomeConfigType: globals.ConfigTypeActive,
			},
		),
	)
}

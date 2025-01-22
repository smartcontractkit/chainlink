package crib

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// DeployHomeChainContracts deploys the home chain contracts so that the chainlink nodes can use the CR address in Capabilities.ExternalRegistry
// Afterwards, we call DeployHomeChainChangeset changeset with nodeinfo ( the peer id and all)
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
	*e, err = commonchangeset.ApplyChangesets(nil, *e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployHomeChainChangeset),
			Config: changeset.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					"NodeOperator": p2pIds,
				},
			},
		},
	})

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
	e, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to initiate new environment: %w", err)
	}
	e.ExistingAddresses = ab
	chainSelectors := e.AllChainSelectors()
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	var prereqCfgs []changeset.DeployPrerequisiteConfigPerChain
	for _, chain := range e.AllChainSelectors() {
		mcmsConfig, err := config.NewConfig(1, []common.Address{e.Chains[chain].DeployerKey.From}, []config.Config{})
		if err != nil {
			return DeployCCIPOutput{}, fmt.Errorf("failed to create mcms config: %w", err)
		}
		cfg[chain] = commontypes.MCMSWithTimelockConfig{
			Canceller:        *mcmsConfig,
			Bypasser:         *mcmsConfig,
			Proposer:         *mcmsConfig,
			TimelockMinDelay: big.NewInt(0),
		}
		prereqCfgs = append(prereqCfgs, changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
	}

	// set up chains
	chainConfigs := make(map[uint64]changeset.ChainConfig)
	nodeInfo, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get node info from env: %w", err)
	}
	for _, chain := range chainSelectors {
		chainConfigs[chain] = changeset.ChainConfig{
			Readers: nodeInfo.NonBootstraps().PeerIDs(),
			FChain:  1,
			EncodableChainConfig: chainconfig.ChainConfig{
				GasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(1000)},
				DAGasPriceDeviationPPB:  cciptypes.BigInt{Int: big.NewInt(1_000_000)},
				OptimisticConfirmations: 1,
			},
		}
	}

	// Setup because we only need to deploy the contracts and distribute job specs
	*e, err = commonchangeset.ApplyChangesets(nil, *e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.UpdateChainConfigChangeset),
			Config: changeset.UpdateChainConfigConfig{
				HomeChainSelector: homeChainSel,
				RemoteChainAdds:   chainConfigs,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    chainSelectors,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployPrerequisitesChangeset),
			Config: changeset.DeployPrerequisiteConfig{
				Configs: prereqCfgs,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config:    cfg,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployChainContractsChangeset),
			Config: changeset.DeployChainContractsConfig{
				ChainSelectors:    chainSelectors,
				HomeChainSelector: homeChainSel,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.CCIPCapabilityJobspecChangeset),
			Config:    struct{}{},
		},
	})
	state, err := changeset.LoadOnchainState(*e)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	// Add all lanes
	for from := range e.Chains {
		for to := range e.Chains {
			if from != to {
				stateChain1 := state.Chains[from]
				newEnv, err := commonchangeset.ApplyChangesets(nil, *e, nil, []commonchangeset.ChangesetApplication{
					{
						Changeset: commonchangeset.WrapChangeSet(changeset.UpdateOnRampsDestsChangeset),
						Config: changeset.UpdateOnRampDestsConfig{
							UpdatesByChain: map[uint64]map[uint64]changeset.OnRampDestinationUpdate{
								from: {
									to: {
										IsEnabled:        true,
										TestRouter:       false,
										AllowListEnabled: false,
									},
								},
							},
						},
					},
					{
						Changeset: commonchangeset.WrapChangeSet(changeset.UpdateFeeQuoterPricesChangeset),
						Config: changeset.UpdateFeeQuoterPricesConfig{
							PricesByChain: map[uint64]changeset.FeeQuoterPriceUpdatePerSource{
								from: {
									TokenPrices: map[common.Address]*big.Int{
										stateChain1.LinkToken.Address(): testhelpers.DefaultLinkPrice,
										stateChain1.Weth9.Address():     testhelpers.DefaultWethPrice,
									},
									GasPrices: map[uint64]*big.Int{
										to: testhelpers.DefaultGasPrice,
									},
								},
							},
						},
					},
					{
						Changeset: commonchangeset.WrapChangeSet(changeset.UpdateFeeQuoterDestsChangeset),
						Config: changeset.UpdateFeeQuoterDestsConfig{
							UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
								from: {
									to: changeset.DefaultFeeQuoterDestChainConfig(),
								},
							},
						},
					},
					{
						Changeset: commonchangeset.WrapChangeSet(changeset.UpdateOffRampSourcesChangeset),
						Config: changeset.UpdateOffRampSourcesConfig{
							UpdatesByChain: map[uint64]map[uint64]changeset.OffRampSourceUpdate{
								to: {
									from: {
										IsEnabled:  true,
										TestRouter: true,
									},
								},
							},
						},
					},
					{
						Changeset: commonchangeset.WrapChangeSet(changeset.UpdateRouterRampsChangeset),
						Config: changeset.UpdateRouterRampsConfig{
							TestRouter: true,
							UpdatesByChain: map[uint64]changeset.RouterUpdates{
								// onRamp update on source chain
								from: {
									OnRampUpdates: map[uint64]bool{
										to: true,
									},
								},
								// off
								from: {
									OffRampUpdates: map[uint64]bool{
										to: true,
									},
								},
							},
						},
					},
				})
				if err != nil {
					return DeployCCIPOutput{}, fmt.Errorf("failed to apply changesets: %w", err)
				}
				e = &newEnv
			}
		}
	}

	addresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get convert address book to address book map: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *deployment.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, err
}

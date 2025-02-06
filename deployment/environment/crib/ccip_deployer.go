package crib

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/rs/zerolog"

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
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
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
	*e, err = commonchangeset.ApplyChangesets(nil, *e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config:    cfg,
		},
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
	e, don, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to initiate new environment: %w", err)
	}
	e.ExistingAddresses = ab
	chainSelectors := e.AllChainSelectors()
	var prereqCfgs []changeset.DeployPrerequisiteConfigPerChain
	for _, chain := range e.AllChainSelectors() {
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
	contractParams := make(map[uint64]changeset.ChainContractParams)
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
		contractParams[chain] = changeset.ChainContractParams{
			FeeQuoterParams: changeset.DefaultFeeQuoterParams(),
			OffRampParams:   changeset.DefaultOffRampParams(),
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
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployChainContractsChangeset),
			Config: changeset.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ContractParamsPerChain: contractParams,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.SetRMNRemoteOnRMNProxyChangeset),
			Config: changeset.SetRMNRemoteOnRMNProxyConfig{
				ChainSelectors: chainSelectors,
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

	var ocrConfigPerSelector = make(map[uint64]changeset.CCIPOCRParams)
	for selector := range e.Chains {
		ocrConfigPerSelector[selector] = changeset.DeriveCCIPOCRParams(changeset.WithDefaultCommitOffChainConfig(feedChainSel, nil),
			changeset.WithDefaultExecuteOffChainConfig(nil),
		)
	}

	*e, err = commonchangeset.ApplyChangesets(nil, *e, nil, []commonchangeset.ChangesetApplication{
		{
			// Add the DONs and candidate commit OCR instances for the chain.
			Changeset: commonchangeset.WrapChangeSet(changeset.AddDonAndSetCandidateChangeset),
			Config: changeset.AddDonAndSetCandidateChangesetConfig{
				SetCandidateConfigBase: changeset.SetCandidateConfigBase{
					HomeChainSelector: homeChainSel,
					FeedChainSelector: feedChainSel,
				},
				PluginInfo: changeset.SetCandidatePluginInfo{
					OCRConfigPerRemoteChainSelector: ocrConfigPerSelector,
					PluginType:                      types.PluginTypeCCIPCommit,
				},
			},
		},
		{
			// Add the exec OCR instances for the new chains.
			Changeset: commonchangeset.WrapChangeSet(changeset.SetCandidateChangeset),
			Config: changeset.SetCandidateChangesetConfig{
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
		},
		{
			// Promote everything
			Changeset: commonchangeset.WrapChangeSet(changeset.PromoteCandidateChangeset),
			Config: changeset.PromoteCandidateChangesetConfig{
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
		},
		{
			// Enable the OCR config on the remote chains.
			Changeset: commonchangeset.WrapChangeSet(changeset.SetOCR3OffRampChangeset),
			Config: changeset.SetOCR3OffRampConfig{
				HomeChainSel:    homeChainSel,
				RemoteChainSels: chainSelectors,
			},
		},
	})
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply changesets: %w", err)
	}

	// Add all lanes
	for src := range e.Chains {
		for dst := range e.Chains {
			if src != dst {
				stateChain1 := state.Chains[src]
				newEnv, err := commonchangeset.ApplyChangesets(nil, *e, nil, []commonchangeset.ChangesetApplication{
					{
						Changeset: commonchangeset.WrapChangeSet(changeset.UpdateOnRampsDestsChangeset),
						Config: changeset.UpdateOnRampDestsConfig{
							UpdatesByChain: map[uint64]map[uint64]changeset.OnRampDestinationUpdate{
								src: {
									dst: {
										IsEnabled:        true,
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
								src: {
									TokenPrices: map[common.Address]*big.Int{
										stateChain1.LinkToken.Address(): testhelpers.DefaultLinkPrice,
										stateChain1.Weth9.Address():     testhelpers.DefaultWethPrice,
									},
									GasPrices: map[uint64]*big.Int{
										dst: testhelpers.DefaultGasPrice,
									},
								},
							},
						},
					},
					{
						Changeset: commonchangeset.WrapChangeSet(changeset.UpdateFeeQuoterDestsChangeset),
						Config: changeset.UpdateFeeQuoterDestsConfig{
							UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
								src: {
									dst: changeset.DefaultFeeQuoterDestChainConfig(true),
								},
							},
						},
					},
					{
						Changeset: commonchangeset.WrapChangeSet(changeset.UpdateOffRampSourcesChangeset),
						Config: changeset.UpdateOffRampSourcesConfig{
							UpdatesByChain: map[uint64]map[uint64]changeset.OffRampSourceUpdate{
								dst: {
									src: {
										IsEnabled: true,
									},
								},
							},
						},
					},
					{
						Changeset: commonchangeset.WrapChangeSet(changeset.UpdateRouterRampsChangeset),
						Config: changeset.UpdateRouterRampsConfig{
							UpdatesByChain: map[uint64]changeset.RouterUpdates{
								src: {
									OffRampUpdates: map[uint64]bool{
										dst: true,
									},
									OnRampUpdates: map[uint64]bool{
										dst: true,
									},
								},
								dst: {
									OffRampUpdates: map[uint64]bool{
										src: true,
									},
									OnRampUpdates: map[uint64]bool{
										src: true,
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

	// distribute funds to transmitters
	// we need to use the nodeinfo from the envConfig here, because multiAddr is not
	// populated in the environment variable
	distributeFunds(lggr, don.PluginNodes(), *e)

	addresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get convert address book to address book map: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *deployment.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, err
}

type RMNNodeConfig struct {
	changeset.RMNNopConfig
	rageproxyKeystore string
	rmnKeystore       string
	passphrase        string
}

func SetupRMNNodeOnAllChains(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab deployment.AddressBook, nodes []RMNNodeConfig) error {
	e, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return err
	}

	rmnNodes := make([]rmn_home.RMNHomeNode, len(nodes))
	for i, node := range nodes {
		rmnNodes[i] = rmn_home.RMNHomeNode{
			PeerId:            node.PeerId,
			OffchainPublicKey: node.OffchainPublicKey,
		}
	}

	_, err = changeset.SetRMNHomeCandidateConfigChangeset(*e, changeset.SetRMNHomeCandidateConfig{
		HomeChainSelector: homeChainSel,
		RMNStaticConfig: rmn_home.RMNHomeStaticConfig{
			Nodes:          rmnNodes,
			OffchainConfig: []byte{},
		},
	})
	if err != nil {
		return err
	}

	state, err := changeset.LoadOnchainState(*e)
	if err != nil {
		return err
	}

	configDigest, err := state.Chains[homeChainSel].RMNHome.GetCandidateDigest(nil)

	if err != nil {
		return err
	}

	_, err = changeset.PromoteRMNHomeCandidateConfigChangeset(*e, changeset.PromoteRMNHomeCandidateConfig{
		HomeChainSelector: homeChainSel,
		DigestToPromote:   configDigest,
	})

	if err != nil {
		return err
	}

	signers := make([]rmn_remote.RMNRemoteSigner, len(nodes))
	for i, node := range nodes {
		signers[i] = node.ToRMNRemoteSigner()
	}

	allChains := e.AllChainSelectors()
	rmnRemoteConfig := make(map[uint64]changeset.RMNRemoteConfig)
	for _, chain := range allChains {
		rmnRemoteConfig[chain] = changeset.RMNRemoteConfig{
			Signers: signers,
			F:       1,
		}
	}

	_, err = changeset.SetRMNRemoteConfigChangeset(*e, changeset.SetRMNRemoteConfig{
		HomeChainSelector: homeChainSel,
		RMNRemoteConfigs:  rmnRemoteConfig,
	})
	if err != nil {
		return err
	}

	updates := make(map[uint64]changeset.OffRampParams)
	for _, chainIdx := range allChains {
		updates[chainIdx] = changeset.OffRampParams{
			IsRMNVerificationDisabled:               false,
			PermissionLessExecutionThresholdSeconds: 86400,
		}
	}

	_, err = changeset.UpdateDynamicConfigOffRampChangeset(*e, changeset.UpdateDynamicConfigOffRampConfig{
		Updates: updates,
	})

	if err != nil {
		return err
	}

	return nil
}

func GenerateRMNNodeIdentities(rmnNodeCount uint) ([]RMNNodeConfig, error) {
	lggr := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	rmnNodeConfigs := make([]RMNNodeConfig, rmnNodeCount)

	for i := uint(0); i < rmnNodeCount; i++ {
		peerID, rawKeystore, _, err := devenv.GeneratePeerID(zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}), "rageproxy", "latest")
		if err != nil {
			return nil, err
		}

		keys, rawRMNKeystore, afnPassphrase, err := devenv.GenerateRMNKeyStore(lggr, "afn2proxy", "latest")
		if err != nil {
			return nil, err
		}

		newPeerID, err := p2pkey.MakePeerID(peerID.String())
		if err != nil {
			return nil, err
		}

		rmnNodeConfigs[i] = RMNNodeConfig{
			RMNNopConfig: changeset.RMNNopConfig{
				NodeIndex:           uint64(i),
				OffchainPublicKey:   [32]byte(keys.OffchainPublicKey),
				EVMOnChainPublicKey: keys.EVMOnchainPublicKey,
				PeerId:              newPeerID,
			},
			rageproxyKeystore: rawKeystore,
			rmnKeystore:       rawRMNKeystore,
			passphrase:        afnPassphrase,
		}
	}
	return rmnNodeConfigs, nil
}

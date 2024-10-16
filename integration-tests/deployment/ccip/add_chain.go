package ccipdeployment

import (
	"fmt"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
)

// NewChainInboundProposal generates a proposal
// to connect the new chain to the existing chains.
func NewChainInboundProposal(
	e deployment.Environment,
	state CCIPOnChainState,
	homeChainSel uint64,
	feedChainSel uint64,
	newChainSel uint64,
	sources []uint64,
	tokenConfig TokenConfig,
	rmnHomeAddress []byte,
) (*timelock.MCMSWithTimelockProposal, error) {
	// Generate proposal which enables new destination (from test router) on all source chains.
	var batches []timelock.BatchChainOperation
	metaDataPerChain := make(map[mcms.ChainIdentifier]mcms.ChainMetadata)
	timelockAddresses := make(map[mcms.ChainIdentifier]common.Address)
	for _, source := range sources {
		chain, _ := chainsel.ChainBySelector(source)
		enableOnRampDest, err := state.Chains[source].OnRamp.ApplyDestChainConfigUpdates(deployment.SimTransactOpts(), []onramp.OnRampDestChainConfigArgs{
			{
				DestChainSelector: newChainSel,
				Router:            state.Chains[source].TestRouter.Address(),
			},
		})
		if err != nil {
			return nil, err
		}
		enableFeeQuoterDest, err := state.Chains[source].FeeQuoter.ApplyDestChainConfigUpdates(
			deployment.SimTransactOpts(),
			[]fee_quoter.FeeQuoterDestChainConfigArgs{
				{
					DestChainSelector: newChainSel,
					DestChainConfig:   defaultFeeQuoterDestChainConfig(),
				},
			})
		if err != nil {
			return nil, err
		}
		initialPrices, err := state.Chains[source].FeeQuoter.UpdatePrices(
			deployment.SimTransactOpts(),
			fee_quoter.InternalPriceUpdates{
				TokenPriceUpdates: []fee_quoter.InternalTokenPriceUpdate{},
				GasPriceUpdates: []fee_quoter.InternalGasPriceUpdate{
					{
						DestChainSelector: newChainSel,
						// TODO: parameterize
						UsdPerUnitGas: big.NewInt(2e12),
					},
				}})
		if err != nil {
			return nil, err
		}
		batches = append(batches, timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(chain.Selector),
			Batch: []mcms.Operation{
				{
					// Enable the source in on ramp
					To:    state.Chains[source].OnRamp.Address(),
					Data:  enableOnRampDest.Data(),
					Value: big.NewInt(0),
				},
				{
					// Set initial dest prices to unblock testing.
					To:    state.Chains[source].FeeQuoter.Address(),
					Data:  initialPrices.Data(),
					Value: big.NewInt(0),
				},
				{
					To:    state.Chains[source].FeeQuoter.Address(),
					Data:  enableFeeQuoterDest.Data(),
					Value: big.NewInt(0),
				},
			},
		})
		opCount, err := state.Chains[source].ProposerMcm.GetOpCount(nil)
		if err != nil {
			return nil, err
		}
		metaDataPerChain[mcms.ChainIdentifier(chain.Selector)] = mcms.ChainMetadata{
			StartingOpCount: opCount.Uint64(),
			MCMAddress:      state.Chains[source].ProposerMcm.Address(),
		}
		timelockAddresses[mcms.ChainIdentifier(chain.Selector)] = state.Chains[source].Timelock.Address()
	}

	// Home chain new don.
	// - Add new DONs for destination to home chain
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return nil, err
	}
	encodedExtraChainConfig, err := chainconfig.EncodeChainConfig(chainconfig.ChainConfig{
		GasPriceDeviationPPB:    ccipocr3.NewBigIntFromInt64(1000),
		DAGasPriceDeviationPPB:  ccipocr3.NewBigIntFromInt64(0),
		OptimisticConfirmations: 1,
	})
	if err != nil {
		return nil, err
	}
	chainConfig := SetupConfigInfo(newChainSel, nodes.NonBootstraps().PeerIDs(),
		nodes.DefaultF(), encodedExtraChainConfig)
	addChain, err := state.Chains[homeChainSel].CCIPHome.ApplyChainConfigUpdates(
		deployment.SimTransactOpts(), nil, []ccip_home.CCIPHomeChainConfigArgs{
			chainConfig,
		})
	if err != nil {
		return nil, err
	}

	newDONArgs, err := BuildAddDONArgs(
		e.Logger,
		state.Chains[newChainSel].OffRamp,
		e.Chains[newChainSel],
		feedChainSel,
		tokenConfig.GetTokenInfo(e.Logger, state.Chains[newChainSel].LinkToken),
		nodes.NonBootstraps(),
		rmnHomeAddress,
	)
	if err != nil {
		return nil, err
	}
	mcmsOps, err := CreateDON(
		e.Logger,
		state.Chains[homeChainSel].CapabilityRegistry,
		state.Chains[homeChainSel].CCIPHome,
		newDONArgs,
		e.Chains[homeChainSel],
		nodes,
	)
	//addDON, err := state.Chains[homeChainSel].CapabilityRegistry.AddDON(SimTransactOpts(),
	//	nodes.NonBootstraps().PeerIDs(), []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
	//		{
	//			CapabilityId: CCIPCapabilityID,
	//			Config:       newDONArgs,
	//		},
	//	}, false, false, nodes.NonBootstraps().DefaultF())
	if err != nil {
		return nil, err
	}
	homeChain, _ := chainsel.ChainBySelector(homeChainSel)
	opCount, err := state.Chains[homeChainSel].ProposerMcm.GetOpCount(nil)
	if err != nil {
		return nil, err
	}
	metaDataPerChain[mcms.ChainIdentifier(homeChain.Selector)] = mcms.ChainMetadata{
		StartingOpCount: opCount.Uint64(),
		MCMAddress:      state.Chains[homeChainSel].ProposerMcm.Address(),
	}
	timelockAddresses[mcms.ChainIdentifier(homeChain.Selector)] = state.Chains[homeChainSel].Timelock.Address()
	batches = append(batches, timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(homeChain.Selector),
		Batch: append([]mcms.Operation{
			{
				// Add the chain first, don needs it to be there.
				To:    state.Chains[homeChainSel].CCIPHome.Address(),
				Data:  addChain.Data(),
				Value: big.NewInt(0),
			},
		}, mcmsOps...),
	})
	return timelock.NewMCMSWithTimelockProposal(
		"1",
		2004259681, // TODO: should be parameterized and based on current block timestamp.
		[]mcms.Signature{},
		false,
		metaDataPerChain,
		timelockAddresses,
		"blah", // TODO
		batches,
		timelock.Schedule,
		"0s", // TODO: Should be parameterized.
	)
}

func NewChainCandidateProposal(
	e deployment.Environment,
	state CCIPOnChainState,
	homeChainSel uint64,
	feedChainSel uint64,
	newChainSel uint64,
	sources []uint64,
	tokenConfig TokenConfig,
	rmnHomeAddress []byte,
) (map[cctypes.PluginType]ccip_home.CCIPHomeOCR3Config, deployment.Nodes, uint32, *timelock.MCMSWithTimelockProposal, error) {
	// Generate proposal which enables new destination (from test router) on all source chains.
	var batches []timelock.BatchChainOperation
	metaDataPerChain := make(map[mcms.ChainIdentifier]mcms.ChainMetadata)
	timelockAddresses := make(map[mcms.ChainIdentifier]common.Address)
	for _, source := range sources {
		chain, _ := chainsel.ChainBySelector(source)
		enableOnRampDest, err := state.Chains[source].OnRamp.ApplyDestChainConfigUpdates(deployment.SimTransactOpts(), []onramp.OnRampDestChainConfigArgs{
			{
				DestChainSelector: newChainSel,
				Router:            state.Chains[source].TestRouter.Address(),
			},
		})
		if err != nil {
			return nil, nil, 0, nil, err
		}
		enableFeeQuoterDest, err := state.Chains[source].FeeQuoter.ApplyDestChainConfigUpdates(
			deployment.SimTransactOpts(),
			[]fee_quoter.FeeQuoterDestChainConfigArgs{
				{
					DestChainSelector: newChainSel,
					DestChainConfig:   defaultFeeQuoterDestChainConfig(),
				},
			})
		if err != nil {
			return nil, nil, 0, nil, err
		}
		initialPrices, err := state.Chains[source].FeeQuoter.UpdatePrices(
			deployment.SimTransactOpts(),
			fee_quoter.InternalPriceUpdates{
				TokenPriceUpdates: []fee_quoter.InternalTokenPriceUpdate{},
				GasPriceUpdates: []fee_quoter.InternalGasPriceUpdate{
					{
						DestChainSelector: newChainSel,
						// TODO: parameterize
						UsdPerUnitGas: big.NewInt(2e12),
					},
				}})
		if err != nil {
			return nil, nil, 0, nil, err
		}
		batches = append(batches, timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(chain.Selector),
			Batch: []mcms.Operation{
				{
					// Enable the source in on ramp
					To:    state.Chains[source].OnRamp.Address(),
					Data:  enableOnRampDest.Data(),
					Value: big.NewInt(0),
				},
				{
					// Set initial dest prices to unblock testing.
					To:    state.Chains[source].FeeQuoter.Address(),
					Data:  initialPrices.Data(),
					Value: big.NewInt(0),
				},
				{
					To:    state.Chains[source].FeeQuoter.Address(),
					Data:  enableFeeQuoterDest.Data(),
					Value: big.NewInt(0),
				},
			},
		})
		opCount, err := state.Chains[source].ProposerMcm.GetOpCount(nil)
		if err != nil {
			return nil, nil, 0, nil, err
		}
		metaDataPerChain[mcms.ChainIdentifier(chain.Selector)] = mcms.ChainMetadata{
			StartingOpCount: opCount.Uint64(),
			MCMAddress:      state.Chains[source].ProposerMcm.Address(),
		}
		timelockAddresses[mcms.ChainIdentifier(chain.Selector)] = state.Chains[source].Timelock.Address()
	}

	// Home chain new don.
	// - Add new DONs for destination to home chain
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return nil, nil, 0, nil, err
	}
	encodedExtraChainConfig, err := chainconfig.EncodeChainConfig(chainconfig.ChainConfig{
		GasPriceDeviationPPB:    ccipocr3.NewBigIntFromInt64(1000),
		DAGasPriceDeviationPPB:  ccipocr3.NewBigIntFromInt64(0),
		OptimisticConfirmations: 1,
	})
	if err != nil {
		return nil, nil, 0, nil, err
	}
	chainConfig := SetupConfigInfo(newChainSel, nodes.NonBootstraps().PeerIDs(),
		nodes.DefaultF(), encodedExtraChainConfig)
	addChain, err := state.Chains[homeChainSel].CCIPHome.ApplyChainConfigUpdates(
		deployment.SimTransactOpts(), nil, []ccip_home.CCIPHomeChainConfigArgs{
			chainConfig,
		})
	if err != nil {
		return nil, nil, 0, nil, err
	}

	newDONArgs, err := BuildAddDONArgs(
		e.Logger,
		state.Chains[newChainSel].OffRamp,
		e.Chains[newChainSel],
		feedChainSel,
		tokenConfig.GetTokenInfo(e.Logger, state.Chains[newChainSel].LinkToken),
		nodes.NonBootstraps(),
		rmnHomeAddress,
	)
	if err != nil {
		return nil, nil, 0, nil, err
	}

	newDonWithCandidateOps, donId, err := createDONwithCandidates_owned(
		state.Chains[homeChainSel].CapabilityRegistry,
		newDONArgs,
		nodes,
	)
	if err != nil {
		return nil, nil, 0, nil, err
	}

	opCount, err := state.Chains[homeChainSel].ProposerMcm.GetOpCount(nil)
	if err != nil {
		return nil, nil, 0, nil, err
	}
	metaDataPerChain[mcms.ChainIdentifier(homeChainSel)] = mcms.ChainMetadata{
		StartingOpCount: opCount.Uint64(),
		MCMAddress:      state.Chains[homeChainSel].ProposerMcm.Address(),
	}
	timelockAddresses[mcms.ChainIdentifier(homeChainSel)] = state.Chains[homeChainSel].Timelock.Address()
	batches = append(batches, timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(homeChainSel),
		Batch: append([]mcms.Operation{
			{
				// Add the chain first, don needs it to be there.
				To:    state.Chains[homeChainSel].CCIPHome.Address(),
				Data:  addChain.Data(),
				Value: big.NewInt(0),
			},
		}, newDonWithCandidateOps...),
	})

	newProp, err := timelock.NewMCMSWithTimelockProposal(
		"1",
		2004259681, // TODO: should be parameterized and based on current block timestamp.
		[]mcms.Signature{},
		false,
		metaDataPerChain,
		timelockAddresses,
		"initial new chain setup",
		batches,
		timelock.Schedule,
		"0s", // TODO: Should be parameterized.
	)
	return newDONArgs, nodes, donId, newProp, err
}

func NewChainPromoteProposal(
	ocr3Configs map[cctypes.PluginType]ccip_home.CCIPHomeOCR3Config,
	state CCIPOnChainState,
	homeSelector uint64,
	nodes deployment.Nodes,
	donId uint32,
) (*timelock.MCMSWithTimelockProposal, error) {
	ccipHome := state.Chains[homeSelector].CCIPHome
	capReg := state.Chains[homeSelector].CapabilityRegistry
	tl := state.Chains[homeSelector].Timelock
	mcm := state.Chains[homeSelector].ProposerMcm
	ops, err := promoteDONCandidates_owned(ccipHome, capReg, ocr3Configs, nodes, donId)
	if err != nil {
		return nil, fmt.Errorf("promoteDONCandidates_owned: %w", err)
	}
	return CreateSingleChainMCMSOps(ops, homeSelector, tl, mcm)
}

func createDONwithCandidates_owned(
	capReg *capabilities_registry.CapabilitiesRegistry,
	ocr3Configs map[cctypes.PluginType]ccip_home.CCIPHomeOCR3Config,
	nodes deployment.Nodes,
) ([]mcms.Operation, uint32, error) {

	commitConfig, ok := ocr3Configs[cctypes.PluginTypeCCIPCommit]
	if !ok {
		return nil, 0, fmt.Errorf("missing commit plugin in ocr3Configs")
	}

	execConfig, ok := ocr3Configs[cctypes.PluginTypeCCIPExec]
	if !ok {
		return nil, 0, fmt.Errorf("missing exec plugin in ocr3Configs")
	}

	tabi, err := ccip_home.CCIPHomeMetaData.GetAbi()
	if err != nil {
		return nil, 0, err
	}
	latestDon, err := LatestCCIPDON(capReg)
	if err != nil {
		return nil, 0, err
	}

	donID := latestDon.Id + 1
	mcmsOps := []mcms.Operation{}

	donOp, err := addNewDonWithoutCapabilites_owned(capReg, nodes)
	if err != nil {
		return nil, 0, err
	}
	mcmsOps = append(mcmsOps, donOp...)

	proposeCommitPluginOp, err := proposePlugin_owned(tabi, donID, commitConfig, capReg, nodes)
	if err != nil {
		return nil, 0, err
	}
	mcmsOps = append(mcmsOps, proposeCommitPluginOp...)

	proposeExecPluginOp, err := proposePlugin_owned(tabi, donID, execConfig, capReg, nodes)
	if err != nil {
		return nil, 0, err
	}
	mcmsOps = append(mcmsOps, proposeExecPluginOp...)

	return mcmsOps, donID, nil
}

func promoteDONCandidates_owned(
	ccipHome *ccip_home.CCIPHome,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ocr3Configs map[cctypes.PluginType]ccip_home.CCIPHomeOCR3Config,
	nodes deployment.Nodes,
	donId uint32) ([]mcms.Operation, error) {

	commitConfig, ok := ocr3Configs[cctypes.PluginTypeCCIPCommit]
	if !ok {
		return nil, fmt.Errorf("missing commit plugin in ocr3Configs")
	}

	execConfig, ok := ocr3Configs[cctypes.PluginTypeCCIPExec]
	if !ok {
		return nil, fmt.Errorf("missing exec plugin in ocr3Configs")
	}

	tabi, err := ccip_home.CCIPHomeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	promoteCommitPluginOp, err := promotePlugin_owned(tabi, donId, commitConfig, nodes, capReg, ccipHome)
	if err != nil {
		return nil, err
	}

	promoteExecPluginOp, err := promotePlugin_owned(tabi, donId, execConfig, nodes, capReg, ccipHome)
	if err != nil {
		return nil, err
	}

	return append(promoteCommitPluginOp, promoteExecPluginOp...), nil
}

func addNewDonWithoutCapabilites_owned(
	capReg *capabilities_registry.CapabilitiesRegistry,
	nodes deployment.Nodes) ([]mcms.Operation, error) {
	addDonTx, err := capReg.AddDON(deployment.SimTransactOpts(), nodes.PeerIDs(), []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{}, false, false, nodes.DefaultF())
	if err != nil {
		return nil, fmt.Errorf("Could not generate AddDon Tx %w", err)
	}

	return []mcms.Operation{{
		To:    capReg.Address(),
		Data:  addDonTx.Data(),
		Value: big.NewInt(0),
	}}, nil
}

func proposePlugin_owned(
	tabi *abi.ABI,
	donID uint32,
	pluginConfig ccip_home.CCIPHomeOCR3Config,
	capReg *capabilities_registry.CapabilitiesRegistry,
	nodes deployment.Nodes) ([]mcms.Operation, error) {

	encodedSetCandidateCall, err := tabi.Pack(
		"setCandidate",
		donID,
		pluginConfig.PluginType,
		pluginConfig,
		[32]byte{},
	)
	if err != nil {
		return nil, fmt.Errorf("pack set candidate call: %w", err)
	}

	// set candidate call
	updateDonCall, err := capReg.UpdateDON(
		deployment.SimTransactOpts(),
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: CCIPCapabilityID,
				Config:       encodedSetCandidateCall,
			},
		},
		false,
		nodes.DefaultF(),
	)
	if err != nil {
		return nil, fmt.Errorf("update don w/ plugin config: %w", err)
	}

	return []mcms.Operation{{
		To:    capReg.Address(),
		Data:  updateDonCall.Data(),
		Value: big.NewInt(0),
	}}, nil
}

func promotePlugin_owned(
	tabi *abi.ABI,
	donID uint32,
	pluginConfig ccip_home.CCIPHomeOCR3Config,
	nodes deployment.Nodes,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome) ([]mcms.Operation, error) {

	candidateDigest, err := ccipHome.GetCandidateDigest(nil, donID, pluginConfig.PluginType)
	if err != nil {
		return nil, fmt.Errorf("get candidate digest: %w", err)
	}
	if candidateDigest == [32]byte{} {
		return nil, fmt.Errorf("candidate digest is empty, expected nonempty")
	}

	activeDigest, err := ccipHome.GetActiveDigest(nil, donID, pluginConfig.PluginType)
	if err != nil {
		return nil, fmt.Errorf("active digest for %d: %w", pluginConfig.PluginType, err)
	}

	// promote candidate call
	encodedPromotionCall, err := tabi.Pack(
		"promoteCandidateAndRevokeActive",
		donID,
		pluginConfig.PluginType,
		candidateDigest,
		activeDigest,
	)
	if err != nil {
		return nil, fmt.Errorf("pack promotion call: %w", err)
	}

	// set candidate call
	updateDonCall, err := capReg.UpdateDON(
		deployment.SimTransactOpts(),
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: CCIPCapabilityID,
				Config:       encodedPromotionCall,
			},
		},
		false,
		nodes.DefaultF(),
	)
	if err != nil {
		return nil, fmt.Errorf("update don w/ plugin config for promotion: %w", err)
	}

	return []mcms.Operation{{
		To:    capReg.Address(),
		Data:  updateDonCall.Data(),
		Value: big.NewInt(0),
	}}, nil
}

package v1_6_test

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	fqv2ops "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/operations/fee_quoter"
	fqv2seq "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/sequences"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type laneDefinition struct {
	Source v1_6.ChainDefinition
	Dest   v1_6.ChainDefinition
}

func getAllPossibleLanes(chains []v1_6.ChainDefinition, disable bool) []v1_6.BidirectionalLaneDefinition {
	lanes := make([]v1_6.BidirectionalLaneDefinition, 0)
	paired := make(map[uint64]map[uint64]bool)

	for i, chainA := range chains {
		for j, chainB := range chains {
			if i == j {
				continue
			}
			if paired[chainA.Selector] != nil && paired[chainA.Selector][chainB.Selector] {
				continue
			}
			if paired[chainB.Selector] != nil && paired[chainB.Selector][chainA.Selector] {
				continue
			}

			lanes = append(lanes, v1_6.BidirectionalLaneDefinition{
				Chains:     [2]v1_6.ChainDefinition{chainA, chainB},
				IsDisabled: disable,
			})
			if paired[chainA.Selector] == nil {
				paired[chainA.Selector] = make(map[uint64]bool)
			}
			paired[chainA.Selector][chainB.Selector] = true
			if paired[chainB.Selector] == nil {
				paired[chainB.Selector] = make(map[uint64]bool)
			}
			paired[chainB.Selector][chainA.Selector] = true
		}
	}

	return lanes
}

func getRemoteChains(chains []v1_6.ChainDefinition, currentIndex int) []v1_6.ChainDefinition {
	remoteChains := make([]v1_6.ChainDefinition, 0, len(chains)-1)
	for i, chain := range chains {
		if i == currentIndex {
			continue
		}
		remoteChains = append(remoteChains, chain)
	}
	return remoteChains
}

func checkBidirectionalLaneConnectivity(
	t *testing.T,
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	chainOne v1_6.ChainDefinition,
	chainTwo v1_6.ChainDefinition,
	testRouter bool,
	disable bool,
) {
	lanes := []laneDefinition{
		{
			Source: chainOne,
			Dest:   chainTwo,
		},
		{
			Source: chainTwo,
			Dest:   chainOne,
		},
	}
	for _, lane := range lanes {
		onRamp := state.Chains[lane.Source.Selector].OnRamp
		offRamp := state.Chains[lane.Dest.Selector].OffRamp
		feeQuoterOnSrc := state.Chains[lane.Source.Selector].FeeQuoter
		routerOnSrc := state.Chains[lane.Source.Selector].Router
		routerOnDest := state.Chains[lane.Dest.Selector].Router
		if testRouter {
			routerOnSrc = state.Chains[lane.Source.Selector].TestRouter
			routerOnDest = state.Chains[lane.Dest.Selector].TestRouter
		}

		destChainConfig, err := onRamp.GetDestChainConfig(nil, lane.Dest.Selector)
		require.NoError(t, err, "must get dest chain config from onRamp")
		routerAddr := routerOnSrc.Address().Hex()
		if disable {
			routerAddr = common.HexToAddress("0x0").Hex()
		}
		require.Equal(t, routerAddr, destChainConfig.Router.Hex(), "router must equal expected")
		require.Equal(t, lane.Dest.AllowListEnabled, destChainConfig.AllowlistEnabled, "allowListEnabled must equal expected")

		srcChainConfig, err := offRamp.GetSourceChainConfig(nil, lane.Source.Selector)
		require.NoError(t, err, "must get src chain config from offRamp")
		require.Equal(t, !disable, srcChainConfig.IsEnabled, "isEnabled must be expected")
		require.Equal(t, lane.Source.RMNVerificationDisabled, srcChainConfig.IsRMNVerificationDisabled, "rmnVerificationDisabled must equal expected")
		require.Equal(t, common.LeftPadBytes(state.Chains[lane.Source.Selector].OnRamp.Address().Bytes(), 32), srcChainConfig.OnRamp, "remote onRamp must be set on offRamp")
		require.Equal(t, routerOnDest.Address().Hex(), srcChainConfig.Router.Hex(), "router must equal expected")

		isOffRamp, err := routerOnSrc.IsOffRamp(nil, lane.Dest.Selector, state.Chains[lane.Source.Selector].OffRamp.Address())
		require.NoError(t, err, "must check if router has offRamp")
		require.Equal(t, !disable, isOffRamp, "isOffRamp result must equal expected")
		onRampOnRouter, err := routerOnSrc.GetOnRamp(nil, lane.Dest.Selector)
		require.NoError(t, err, "must get onRamp from router")
		onRampAddr := state.Chains[lane.Source.Selector].OnRamp.Address().Hex()
		if disable {
			onRampAddr = common.HexToAddress("0x0").Hex()
		}
		require.Equal(t, onRampAddr, onRampOnRouter.Hex(), "onRamp must equal expected")

		isOffRamp, err = routerOnDest.IsOffRamp(nil, lane.Source.Selector, state.Chains[lane.Dest.Selector].OffRamp.Address())
		require.NoError(t, err, "must check if router has offRamp")
		require.Equal(t, !disable, isOffRamp, "isOffRamp result must equal expected")
		onRampOnRouter, err = routerOnDest.GetOnRamp(nil, lane.Source.Selector)
		require.NoError(t, err, "must get onRamp from router")
		onRampAddr = state.Chains[lane.Dest.Selector].OnRamp.Address().Hex()
		if disable {
			onRampAddr = common.HexToAddress("0x0").Hex()
		}
		require.Equal(t, onRampAddr, onRampOnRouter.Hex(), "onRamp must equal expected")

		feeQuoterDestConfig, err := feeQuoterOnSrc.GetDestChainConfig(nil, lane.Dest.Selector)
		require.NoError(t, err, "must get dest chain config from feeQuoter")
		require.Equal(t, lane.Dest.FeeQuoterDestChainConfig, feeQuoterDestConfig, "feeQuoter dest chain config must equal expected")

		price, err := feeQuoterOnSrc.GetDestinationChainGasPrice(nil, lane.Dest.Selector)
		require.NoError(t, err, "must get price from feeQuoter")
		require.Equal(t, lane.Dest.GasPrice, price.Value, "price must equal expected")
	}
}

func TestBuildConfigs(t *testing.T) {
	selectors := []uint64{1, 2}

	chains := make([]v1_6.ChainDefinition, len(selectors))
	for i, selector := range selectors {
		chains[i] = v1_6.ChainDefinition{
			ConnectionConfig: v1_6.ConnectionConfig{
				RMNVerificationDisabled: true,
				AllowListEnabled:        false,
			},
			Selector:                 selector,
			GasPrice:                 big.NewInt(1e17),
			FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
		}
	}

	cfg := v1_6.UpdateBidirectionalLanesConfig{
		TestRouter: false,
		MCMSConfig: &proposalutils.TimelockConfig{
			MinDelay:   0 * time.Second,
			MCMSAction: types.TimelockActionSchedule,
		},
		Lanes: getAllPossibleLanes(chains, false),
	}

	configs := cfg.BuildConfigs()

	require.Equal(t, v1_6.UpdateFeeQuoterDestsConfig{
		UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
			1: {
				2: v1_6.DefaultFeeQuoterDestChainConfig(true),
			},
			2: {
				1: v1_6.DefaultFeeQuoterDestChainConfig(true),
			},
		},
		MCMS: cfg.MCMSConfig,
	}, configs.UpdateFeeQuoterDestsConfig)
	require.Equal(t, v1_6.UpdateFeeQuoterPricesConfig{
		PricesByChain: map[uint64]v1_6.FeeQuoterPriceUpdatePerSource{
			1: {
				GasPrices: map[uint64]*big.Int{
					2: big.NewInt(1e17),
				},
			},
			2: {
				GasPrices: map[uint64]*big.Int{
					1: big.NewInt(1e17),
				},
			},
		},
		MCMS: cfg.MCMSConfig,
	}, configs.UpdateFeeQuoterPricesConfig)
	require.Equal(t, v1_6.UpdateOffRampSourcesConfig{
		UpdatesByChain: map[uint64]map[uint64]v1_6.OffRampSourceUpdate{
			1: {
				2: {
					IsEnabled:                 true,
					TestRouter:                false,
					IsRMNVerificationDisabled: true,
				},
			},
			2: {
				1: {
					IsEnabled:                 true,
					TestRouter:                false,
					IsRMNVerificationDisabled: true,
				},
			},
		},
		MCMS: cfg.MCMSConfig,
	}, configs.UpdateOffRampSourcesConfig)
	require.Equal(t, v1_6.UpdateOnRampDestsConfig{
		UpdatesByChain: map[uint64]map[uint64]v1_6.OnRampDestinationUpdate{
			1: {
				2: {
					IsEnabled:        true,
					TestRouter:       false,
					AllowListEnabled: false,
				},
			},
			2: {
				1: {
					IsEnabled:        true,
					TestRouter:       false,
					AllowListEnabled: false,
				},
			},
		},
		MCMS: cfg.MCMSConfig,
	}, configs.UpdateOnRampDestsConfig)
	require.Equal(t, v1_6.UpdateRouterRampsConfig{
		UpdatesByChain: map[uint64]v1_6.RouterUpdates{
			1: {
				OnRampUpdates: map[uint64]bool{
					2: true,
				},
				OffRampUpdates: map[uint64]bool{
					2: true,
				},
			},
			2: {
				OnRampUpdates: map[uint64]bool{
					1: true,
				},
				OffRampUpdates: map[uint64]bool{
					1: true,
				},
			},
		},
		MCMS: cfg.MCMSConfig,
	}, configs.UpdateRouterRampsConfig)
}

func TestUpdateBidirectionalLanesChangeset(t *testing.T) {
	t.Parallel()

	type test struct {
		Msg        string
		TestRouter bool
		MCMS       *proposalutils.TimelockConfig
		Disable    bool
	}

	mcmsConfig := &proposalutils.TimelockConfig{
		MinDelay:   0 * time.Second,
		MCMSAction: types.TimelockActionSchedule,
	}

	tests := []test{
		{
			Msg:        "Use production router (with MCMS) & disable afterwards",
			TestRouter: false,
			MCMS:       mcmsConfig,
			Disable:    true,
		},
		{
			Msg:        "Use production router (with MCMS)",
			TestRouter: false,
			MCMS:       mcmsConfig,
		},
		{
			Msg:        "Use test router (without MCMS)",
			TestRouter: true,
			MCMS:       nil,
		},
		{
			Msg:        "Use test router (with MCMS for other contracts)",
			TestRouter: true,
			MCMS:       mcmsConfig,
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(t, func(testCfg *testhelpers.TestConfigs) {
				testCfg.Chains = 3
			})
			e := deployedEnvironment.Env

			state, err := stateview.LoadOnchainState(e)
			require.NoError(t, err, "must load onchain state")

			selectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))

			if test.MCMS != nil {
				contractsToTransfer := make(map[uint64][]common.Address, len(selectors))
				for _, selector := range selectors {
					contractsToTransfer[selector] = []common.Address{
						state.Chains[selector].OnRamp.Address(),
						state.Chains[selector].OffRamp.Address(),
						state.Chains[selector].Router.Address(),
						state.Chains[selector].FeeQuoter.Address(),
						state.Chains[selector].TokenAdminRegistry.Address(),
						state.Chains[selector].RMNRemote.Address(),
						state.Chains[selector].RMNProxy.Address(),
						state.Chains[selector].NonceManager.Address(),
					}
				}
				e, err = commonchangeset.Apply(t, e,
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
						commonchangeset.TransferToMCMSWithTimelockConfig{
							ContractsByChain: contractsToTransfer,
							MCMSConfig: proposalutils.TimelockConfig{
								MinDelay:   0 * time.Second,
								MCMSAction: types.TimelockActionSchedule,
							},
						},
					),
				)
				require.NoError(t, err, "must apply TransferToMCMSWithTimelock")
			}

			chains := make([]v1_6.ChainDefinition, len(selectors))
			for i, selector := range selectors {
				chains[i] = v1_6.ChainDefinition{
					ConnectionConfig: v1_6.ConnectionConfig{
						RMNVerificationDisabled: true,
						AllowListEnabled:        false,
					},
					Selector:                 selector,
					GasPrice:                 big.NewInt(1e17),
					FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
				}
			}

			e, err = commonchangeset.Apply(t, e,
				commonchangeset.Configure(
					v1_6.UpdateBidirectionalLanesChangeset,
					v1_6.UpdateBidirectionalLanesConfig{
						TestRouter: test.TestRouter,
						MCMSConfig: test.MCMS,
						Lanes:      getAllPossibleLanes(chains, false),
					},
				),
			)
			require.NoError(t, err, "must apply AddBidirectionalLanesChangeset")

			for i, chain := range chains {
				remoteChains := getRemoteChains(chains, i)
				for _, remoteChain := range remoteChains {
					checkBidirectionalLaneConnectivity(t, e, state, chain, remoteChain, test.TestRouter, false)
				}
			}

			if test.Disable {
				e, err = commonchangeset.Apply(t, e,
					commonchangeset.Configure(
						v1_6.UpdateBidirectionalLanesChangeset,
						v1_6.UpdateBidirectionalLanesConfig{
							TestRouter: test.TestRouter,
							MCMSConfig: test.MCMS,
							Lanes:      getAllPossibleLanes(chains, true),
						},
					),
				)
				require.NoError(t, err, "must apply AddBidirectionalLanesChangeset")

				for i, chain := range chains {
					remoteChains := getRemoteChains(chains, i)
					for _, remoteChain := range remoteChains {
						checkBidirectionalLaneConnectivity(t, e, state, chain, remoteChain, test.TestRouter, true)
					}
				}
			}
		})
	}
}

func TestUpdateBidirectionalLanesChangesetWithV2FeeQuoter(t *testing.T) {
	t.Parallel()

	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(t, func(testCfg *testhelpers.TestConfigs) {
		testCfg.Chains = 3
	})
	e := deployedEnvironment.Env

	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err, "must load onchain state")

	selectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	require.Len(t, selectors, 3, "must have 3 chains")

	v2FQChainSel := selectors[0]
	evmChain := e.BlockChains.EVMChains()[v2FQChainSel]

	// Deploy a v2 FeeQuoter on the first chain
	parsedABI, err := abi.JSON(strings.NewReader(fqv2ops.FeeQuoterABI))
	require.NoError(t, err, "must parse v2 FeeQuoter ABI")

	fqV2Addr, tx, _, err := bind.DeployContract(
		evmChain.DeployerKey,
		parsedABI,
		common.FromHex(fqv2ops.FeeQuoterBin),
		evmChain.Client,
		fqv2ops.StaticConfig{
			MaxFeeJuelsPerMsg: big.NewInt(1e18),
			LinkToken:         state.Chains[v2FQChainSel].LinkToken.Address(),
		},
		[]common.Address{evmChain.DeployerKey.From},
		[]fqv2ops.TokenTransferFeeConfigArgs{},
		[]fqv2ops.DestChainConfigArgs{},
	)
	require.NoError(t, err, "must deploy v2 FeeQuoter")

	_, err = evmChain.Confirm(tx)
	require.NoError(t, err, "must confirm v2 FeeQuoter deployment")

	// AddressBook does NOT contain FeeQuoter v2
	// Add FeeQuoter v2 only to the DataStore
	ds, err := shared.PopulateDataStore(e.ExistingAddresses)
	require.NoError(t, err, "must populate datastore from existing addresses")

	err = ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: v2FQChainSel,
		Address:       fqV2Addr.Hex(),
		Type:          datastore.ContractType(fqv2ops.ContractType),
		Version:       fqv2ops.Version,
		Qualifier:     fmt.Sprintf("%s-%s", fqV2Addr.Hex(), fqv2ops.ContractType),
	})
	require.NoError(t, err, "must add v2 FeeQuoter to datastore")

	e.DataStore = ds.Seal()

	chains := make([]v1_6.ChainDefinition, len(selectors))
	for i, selector := range selectors {
		chains[i] = v1_6.ChainDefinition{
			ConnectionConfig: v1_6.ConnectionConfig{
				RMNVerificationDisabled: true,
				AllowListEnabled:        false,
			},
			Selector:                 selector,
			GasPrice:                 big.NewInt(1e17),
			FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
		}
	}

	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_6.UpdateBidirectionalLanesChangeset,
			v1_6.UpdateBidirectionalLanesConfig{
				TestRouter: true,
				MCMSConfig: nil,
				Lanes:      getAllPossibleLanes(chains, false),
			},
		),
	)
	require.NoError(t, err, "must apply UpdateBidirectionalLanesChangeset")

	// v1.6 source chains must have dest configs for all remote chains including the chain with v2 FQ
	for _, srcSel := range selectors {
		if srcSel == v2FQChainSel {
			continue
		}
		fq := state.Chains[srcSel].FeeQuoter
		for _, destSel := range selectors {
			if destSel == srcSel {
				continue
			}
			destCfg, err := fq.GetDestChainConfig(nil, destSel)
			require.NoError(t, err, "must get dest chain config from feeQuoter")
			require.Equal(t, v1_6.DefaultFeeQuoterDestChainConfig(true), destCfg, "feeQuoter dest chain config must equal expected")

			price, err := fq.GetDestinationChainGasPrice(nil, destSel)
			require.NoError(t, err, "must get price from feeQuoter")
			require.Equal(t, big.NewInt(1e17), price.Value, "price must equal expected")
		}
	}

	// v2 FeeQuoter on the chain with v2 FQ must have dest configs and prices for each remote chain
	fqV2, err := fqv2ops.NewFeeQuoterContract(fqV2Addr, evmChain.Client)
	require.NoError(t, err, "must bind v2 FeeQuoter")

	for _, destSel := range selectors {
		if destSel == v2FQChainSel {
			continue
		}
		destCfg, err := fqV2.GetDestChainConfig(nil, destSel)
		require.NoError(t, err, "must get v2 FQ dest chain config")
		require.True(t, destCfg.IsEnabled, "v2 dest chain config must be enabled")
		require.Equal(t, uint16(10), destCfg.NetworkFeeUSDCents, "NetworkFeeUSDCents must be converted to uint16")
		require.Equal(t, fqv2seq.LinkFeeMultiplierPercent, destCfg.LinkFeeMultiplierPercent, "LinkFeeMultiplierPercent must be set")

		price, err := fqV2.GetDestinationChainGasPrice(nil, destSel)
		require.NoError(t, err, "must get v2 FQ gas price")
		require.Equal(t, big.NewInt(1e17), price.Value, "price must equal expected")
	}

	// Active v1.6 FeeQuoter on the chain with v2 FQ must also receive updates because
	// it is still the FeeQuoter referenced by the OnRamp during migration.
	activeV1FQ := state.Chains[v2FQChainSel].FeeQuoter
	for _, destSel := range selectors {
		if destSel == v2FQChainSel {
			continue
		}
		destCfg, err := activeV1FQ.GetDestChainConfig(nil, destSel)
		require.NoError(t, err, "must get active v1.6 FQ dest config")
		require.Equal(t, v1_6.DefaultFeeQuoterDestChainConfig(true), destCfg, "active v1.6 FQ must have dest config")

		price, err := activeV1FQ.GetDestinationChainGasPrice(nil, destSel)
		require.NoError(t, err, "must get active v1.6 FQ gas price")
		require.Equal(t, big.NewInt(1e17), price.Value, "active v1.6 FQ gas price must equal expected")
	}
}

// TestUpdateBidirectionalLanesIdempotentWithV2FeeQuoter verifies that running the changeset
// twice on a chain with both v1.6 and v2 FQs is idempotent — the second run must filter out
// all already-enabled dests from BOTH FQs using the correct addresses.
func TestUpdateBidirectionalLanesIdempotentWithV2FeeQuoter(t *testing.T) {
	t.Parallel()

	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(t, func(testCfg *testhelpers.TestConfigs) {
		testCfg.Chains = 3
	})
	e := deployedEnvironment.Env

	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err, "must load onchain state")

	selectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	require.Len(t, selectors, 3, "must have 3 chains")

	v2FQChainSel := selectors[0]
	evmChain := e.BlockChains.EVMChains()[v2FQChainSel]

	// Deploy a v2 FeeQuoter on the first chain
	parsedABI, err := abi.JSON(strings.NewReader(fqv2ops.FeeQuoterABI))
	require.NoError(t, err, "must parse v2 FeeQuoter ABI")

	fqV2Addr, tx, _, err := bind.DeployContract(
		evmChain.DeployerKey,
		parsedABI,
		common.FromHex(fqv2ops.FeeQuoterBin),
		evmChain.Client,
		fqv2ops.StaticConfig{
			MaxFeeJuelsPerMsg: big.NewInt(1e18),
			LinkToken:         state.Chains[v2FQChainSel].LinkToken.Address(),
		},
		[]common.Address{evmChain.DeployerKey.From},
		[]fqv2ops.TokenTransferFeeConfigArgs{},
		[]fqv2ops.DestChainConfigArgs{},
	)
	require.NoError(t, err, "must deploy v2 FeeQuoter")

	_, err = evmChain.Confirm(tx)
	require.NoError(t, err, "must confirm v2 FeeQuoter deployment")

	ds, err := shared.PopulateDataStore(e.ExistingAddresses)
	require.NoError(t, err, "must populate datastore from existing addresses")

	err = ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: v2FQChainSel,
		Address:       fqV2Addr.Hex(),
		Type:          datastore.ContractType(fqv2ops.ContractType),
		Version:       fqv2ops.Version,
		Qualifier:     fmt.Sprintf("%s-%s", fqV2Addr.Hex(), fqv2ops.ContractType),
	})
	require.NoError(t, err, "must add v2 FeeQuoter to datastore")

	e.DataStore = ds.Seal()

	chains := make([]v1_6.ChainDefinition, len(selectors))
	for i, selector := range selectors {
		chains[i] = v1_6.ChainDefinition{
			ConnectionConfig: v1_6.ConnectionConfig{
				RMNVerificationDisabled: true,
				AllowListEnabled:        false,
			},
			Selector:                 selector,
			GasPrice:                 big.NewInt(1e17),
			FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
		}
	}

	cfg := v1_6.UpdateBidirectionalLanesConfig{
		TestRouter: true,
		MCMSConfig: nil,
		Lanes:      getAllPossibleLanes(chains, false),
	}

	// First run — populates both v1 and v2 FQs
	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(v1_6.UpdateBidirectionalLanesChangeset, cfg),
	)
	require.NoError(t, err, "first run must succeed")

	// Second run — must succeed without errors; all dests already configured should be filtered
	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(v1_6.UpdateBidirectionalLanesChangeset, cfg),
	)
	require.NoError(t, err, "second (idempotent) run must succeed")

	// Verify state is unchanged — v2 FQ still has correct configs
	fqV2, err := fqv2ops.NewFeeQuoterContract(fqV2Addr, evmChain.Client)
	require.NoError(t, err, "must bind v2 FeeQuoter")

	for _, destSel := range selectors {
		if destSel == v2FQChainSel {
			continue
		}
		destCfg, err := fqV2.GetDestChainConfig(nil, destSel)
		require.NoError(t, err, "must get v2 FQ dest chain config after idempotent run")
		require.True(t, destCfg.IsEnabled, "v2 dest chain config must still be enabled")
	}

	// Verify active v1.6 FQ still has correct configs
	activeV1FQ := state.Chains[v2FQChainSel].FeeQuoter
	for _, destSel := range selectors {
		if destSel == v2FQChainSel {
			continue
		}
		destCfg, err := activeV1FQ.GetDestChainConfig(nil, destSel)
		require.NoError(t, err, "must get active v1.6 FQ dest config after idempotent run")
		require.Equal(t, v1_6.DefaultFeeQuoterDestChainConfig(true), destCfg,
			"active v1.6 FQ dest config must be unchanged after idempotent run")
	}

	// Verify v1-only chains are also unchanged
	for _, srcSel := range selectors {
		if srcSel == v2FQChainSel {
			continue
		}
		fq := state.Chains[srcSel].FeeQuoter
		for _, destSel := range selectors {
			if destSel == srcSel {
				continue
			}
			destCfg, err := fq.GetDestChainConfig(nil, destSel)
			require.NoError(t, err, "must get v1 FQ dest config after idempotent run")
			require.Equal(t, v1_6.DefaultFeeQuoterDestChainConfig(true), destCfg,
				"v1 FQ dest config must be unchanged after idempotent run")
		}
	}
}

func TestFilterOutExistingDestChainConfigs(t *testing.T) {
	t.Parallel()

	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(t, func(testCfg *testhelpers.TestConfigs) {
		testCfg.Chains = 2
	})
	e := deployedEnvironment.Env

	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err, "must load onchain state")

	selectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	require.Len(t, selectors, 2, "must have 2 chains")

	chainSel := selectors[0]
	otherChainSel := selectors[1]
	evmChain := e.BlockChains.EVMChains()[chainSel]

	parsedABI, err := abi.JSON(strings.NewReader(fqv2ops.FeeQuoterABI))
	require.NoError(t, err, "must parse v2 FeeQuoter ABI")

	// Deploy a v2 FeeQuoter, then configure one destination (otherChainSel) after deployment.
	// Constructor-time dest config with sparse fields reverts, so we apply it post-deploy.
	fqV2Addr, tx, _, err := bind.DeployContract(
		evmChain.DeployerKey,
		parsedABI,
		common.FromHex(fqv2ops.FeeQuoterBin),
		evmChain.Client,
		fqv2ops.StaticConfig{
			MaxFeeJuelsPerMsg: big.NewInt(1e18),
			LinkToken:         state.Chains[chainSel].LinkToken.Address(),
		},
		[]common.Address{evmChain.DeployerKey.From},
		[]fqv2ops.TokenTransferFeeConfigArgs{},
		[]fqv2ops.DestChainConfigArgs{},
	)
	require.NoError(t, err, "must deploy v2 FeeQuoter")

	_, err = evmChain.Confirm(tx)
	require.NoError(t, err, "must confirm v2 FeeQuoter deployment")

	fqV2, err := fqv2ops.NewFeeQuoterContract(fqV2Addr, evmChain.Client)
	require.NoError(t, err, "must bind v2 FeeQuoter")

	// Convert a default v1.6 config to v2 format and apply it so otherChainSel is already enabled
	v2Cfgs, err := v1_6.ConvertV16FeeQuoterDestUpdatesToV2([]fee_quoter.FeeQuoterDestChainConfigArgs{
		{
			DestChainSelector: otherChainSel,
			DestChainConfig:   v1_6.DefaultFeeQuoterDestChainConfig(true),
		},
	})
	require.NoError(t, err, "must convert v1.6 dest config to v2")

	applyTx, err := fqV2.ApplyDestChainConfigUpdates(evmChain.DeployerKey, v2Cfgs)
	require.NoError(t, err, "must apply dest chain config")
	_, err = evmChain.Confirm(applyTx)
	require.NoError(t, err, "must confirm dest chain config tx")

	unconfiguredChainSel := uint64(999)

	// Call FilterOutExistingDestChainConfigs with both destinations
	input := []fqv2ops.DestChainConfigArgs{
		{
			DestChainSelector: otherChainSel,
			DestChainConfig:   fqv2ops.DestChainConfig{IsEnabled: true, MaxDataBytes: 50_000},
		},
		{
			DestChainSelector: unconfiguredChainSel,
			DestChainConfig:   fqv2ops.DestChainConfig{IsEnabled: true, MaxDataBytes: 60_000},
		},
	}

	filtered, err := v1_6.FilterOutExistingDestChainConfigs(e, fqV2Addr, chainSel, input)
	require.NoError(t, err, "FilterOutExistingDestChainConfigs must not error")

	// Only the unconfigured chain should remain
	require.Len(t, filtered, 1, "must filter out the already-enabled destination")
	assert.Equal(t, unconfiguredChainSel, filtered[0].DestChainSelector,
		"remaining entry must be the unconfigured chain")
	assert.Equal(t, uint32(60_000), filtered[0].DestChainConfig.MaxDataBytes,
		"remaining entry must preserve original config")

	// Configure otherChainSel on the v1.6 FeeQuoter so it's already enabled
	fqV16 := state.Chains[chainSel].FeeQuoter
	applyTxV16, err := fqV16.ApplyDestChainConfigUpdates(evmChain.DeployerKey, []fee_quoter.FeeQuoterDestChainConfigArgs{
		{
			DestChainSelector: otherChainSel,
			DestChainConfig:   v1_6.DefaultFeeQuoterDestChainConfig(true),
		},
	})
	require.NoError(t, err, "must apply v1.6 dest chain config")
	_, err = evmChain.Confirm(applyTxV16)
	require.NoError(t, err, "must confirm v1.6 dest chain config tx")

	// Call FilterOutExistingDestChainConfigs with v1.6 types
	v16Input := []fee_quoter.FeeQuoterDestChainConfigArgs{
		{
			DestChainSelector: otherChainSel,
			DestChainConfig:   v1_6.DefaultFeeQuoterDestChainConfig(true),
		},
		{
			DestChainSelector: unconfiguredChainSel,
			DestChainConfig:   v1_6.DefaultFeeQuoterDestChainConfig(true),
		},
	}

	filteredV16, err := v1_6.FilterOutExistingDestChainConfigs(e, fqV16.Address(), chainSel, v16Input)
	require.NoError(t, err, "FilterOutExistingDestChainConfigs must not error")

	require.Len(t, filteredV16, 1, "must filter out the already-enabled destination")
	assert.Equal(t, unconfiguredChainSel, filteredV16[0].DestChainSelector,
		"remaining entry must be the unconfigured chain")
}

func TestUpdateBidirectionalLanesChangesetWithV2FeeQuoterWithMCMS(t *testing.T) {
	t.Parallel()

	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(t, func(testCfg *testhelpers.TestConfigs) {
		testCfg.Chains = 3
	})
	e := deployedEnvironment.Env

	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err, "must load onchain state")

	selectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	require.Len(t, selectors, 3, "must have 3 chains")

	v2FQChainSel := selectors[0]
	evmChain := e.BlockChains.EVMChains()[v2FQChainSel]

	// Deploy a v2 FeeQuoter on the first chain
	parsedABI, err := abi.JSON(strings.NewReader(fqv2ops.FeeQuoterABI))
	require.NoError(t, err, "must parse v2 FeeQuoter ABI")

	fqV2Addr, tx, _, err := bind.DeployContract(
		evmChain.DeployerKey,
		parsedABI,
		common.FromHex(fqv2ops.FeeQuoterBin),
		evmChain.Client,
		fqv2ops.StaticConfig{
			MaxFeeJuelsPerMsg: big.NewInt(1e18),
			LinkToken:         state.Chains[v2FQChainSel].LinkToken.Address(),
		},
		// Include timelock as authorized caller
		[]common.Address{evmChain.DeployerKey.From, state.Chains[v2FQChainSel].Timelock.Address()},
		[]fqv2ops.TokenTransferFeeConfigArgs{},
		[]fqv2ops.DestChainConfigArgs{},
	)
	require.NoError(t, err, "must deploy v2 FeeQuoter")

	_, err = evmChain.Confirm(tx)
	require.NoError(t, err, "must confirm v2 FeeQuoter deployment")

	// AddressBook does NOT contain FeeQuoter v2
	// Add FeeQuoter v2 only to the DataStore
	ds, err := shared.PopulateDataStore(e.ExistingAddresses)
	require.NoError(t, err, "must populate datastore from existing addresses")

	err = ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: v2FQChainSel,
		Address:       fqV2Addr.Hex(),
		Type:          datastore.ContractType(fqv2ops.ContractType),
		Version:       fqv2ops.Version,
		Qualifier:     fmt.Sprintf("%s-%s", fqV2Addr.Hex(), fqv2ops.ContractType),
	})
	require.NoError(t, err, "must add v2 FeeQuoter to datastore")

	e.DataStore = ds.Seal()

	// Transfer contracts to MCMS timelock
	mcmsConfig := &proposalutils.TimelockConfig{
		MinDelay:   0 * time.Second,
		MCMSAction: types.TimelockActionSchedule,
	}

	contractsToTransfer := make(map[uint64][]common.Address, len(selectors))
	for _, selector := range selectors {
		contractsToTransfer[selector] = []common.Address{
			state.Chains[selector].OnRamp.Address(),
			state.Chains[selector].OffRamp.Address(),
			state.Chains[selector].Router.Address(),
			state.Chains[selector].FeeQuoter.Address(),
			state.Chains[selector].TokenAdminRegistry.Address(),
			state.Chains[selector].RMNRemote.Address(),
			state.Chains[selector].RMNProxy.Address(),
			state.Chains[selector].NonceManager.Address(),
		}
	}
	// Also transfer the v2 FeeQuoter to MCMS timelock
	contractsToTransfer[v2FQChainSel] = append(contractsToTransfer[v2FQChainSel], fqV2Addr)

	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
			commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: contractsToTransfer,
				MCMSConfig: proposalutils.TimelockConfig{
					MinDelay:   0 * time.Second,
					MCMSAction: types.TimelockActionSchedule,
				},
			},
		),
	)
	require.NoError(t, err, "must apply TransferToMCMSWithTimelock")

	chains := make([]v1_6.ChainDefinition, len(selectors))
	for i, selector := range selectors {
		chains[i] = v1_6.ChainDefinition{
			ConnectionConfig: v1_6.ConnectionConfig{
				RMNVerificationDisabled: true,
				AllowListEnabled:        false,
			},
			Selector:                 selector,
			GasPrice:                 big.NewInt(1e17),
			FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
		}
	}

	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			v1_6.UpdateBidirectionalLanesChangeset,
			v1_6.UpdateBidirectionalLanesConfig{
				TestRouter: false,
				MCMSConfig: mcmsConfig,
				Lanes:      getAllPossibleLanes(chains, false),
			},
		),
	)
	require.NoError(t, err, "must apply UpdateBidirectionalLanesChangeset with MCMS")

	// Save v1.6 FQ reference before reload; LoadOnchainState picks the highest-version
	// FeeQuoter as chainState.FeeQuoter, which will be the v2 FQ after registration.
	activeV1FQ := state.Chains[v2FQChainSel].FeeQuoter

	// Reload state after MCMS proposal execution
	state, err = stateview.LoadOnchainState(e)
	require.NoError(t, err, "must reload onchain state")

	// v1.6 source chains must have dest configs for all remote chains
	for _, srcSel := range selectors {
		if srcSel == v2FQChainSel {
			continue
		}
		fq := state.Chains[srcSel].FeeQuoter
		for _, destSel := range selectors {
			if destSel == srcSel {
				continue
			}
			destCfg, err := fq.GetDestChainConfig(nil, destSel)
			require.NoError(t, err, "must get dest chain config from feeQuoter")
			require.Equal(t, v1_6.DefaultFeeQuoterDestChainConfig(true), destCfg, "feeQuoter dest chain config must equal expected")

			price, err := fq.GetDestinationChainGasPrice(nil, destSel)
			require.NoError(t, err, "must get price from feeQuoter")
			require.Equal(t, big.NewInt(1e17), price.Value, "price must equal expected")
		}
	}

	// v2 FeeQuoter must have correct state
	fqV2, err := fqv2ops.NewFeeQuoterContract(fqV2Addr, evmChain.Client)
	require.NoError(t, err, "must bind v2 FeeQuoter")

	for _, destSel := range selectors {
		if destSel == v2FQChainSel {
			continue
		}
		destCfg, err := fqV2.GetDestChainConfig(nil, destSel)
		require.NoError(t, err, "must get v2 FQ dest chain config")
		require.True(t, destCfg.IsEnabled, "v2 dest chain config must be enabled")
		require.Equal(t, uint16(10), destCfg.NetworkFeeUSDCents, "NetworkFeeUSDCents must be converted to uint16")
		require.Equal(t, fqv2seq.LinkFeeMultiplierPercent, destCfg.LinkFeeMultiplierPercent, "LinkFeeMultiplierPercent must be set")

		price, err := fqV2.GetDestinationChainGasPrice(nil, destSel)
		require.NoError(t, err, "must get v2 FQ gas price")
		require.Equal(t, big.NewInt(1e17), price.Value, "price must equal expected")
	}

	// Active v1.6 FeeQuoter on the chain with v2 FQ must also receive updates because
	// it is still the FeeQuoter referenced by the OnRamp during migration.
	for _, destSel := range selectors {
		if destSel == v2FQChainSel {
			continue
		}
		destCfg, err := activeV1FQ.GetDestChainConfig(nil, destSel)
		require.NoError(t, err, "must get active v1.6 FQ dest config")
		require.Equal(t, v1_6.DefaultFeeQuoterDestChainConfig(true), destCfg, "active v1.6 FQ must have dest config")

		price, err := activeV1FQ.GetDestinationChainGasPrice(nil, destSel)
		require.NoError(t, err, "must get active v1.6 FQ gas price")
		require.Equal(t, big.NewInt(1e17), price.Value, "active v1.6 FQ gas price must equal expected")
	}
}

func TestConvertV16FeeQuoterDestUpdatesToV2(t *testing.T) {
	t.Parallel()

	t.Run("converts valid dest chain configs", func(t *testing.T) {
		in := []fee_quoter.FeeQuoterDestChainConfigArgs{
			{
				DestChainSelector: 100,
				DestChainConfig: fee_quoter.FeeQuoterDestChainConfig{
					IsEnabled:                   true,
					MaxDataBytes:                30000,
					MaxPerMsgGasLimit:           3000000,
					DestGasOverhead:             50000,
					DestGasPerPayloadByteBase:   16,
					ChainFamilySelector:         [4]byte{1, 2, 3, 4},
					DefaultTokenFeeUSDCents:     25,
					DefaultTokenDestGasOverhead: 90000,
					DefaultTxGasLimit:           200000,
					NetworkFeeUSDCents:          10,
				},
			},
		}

		out, err := v1_6.ConvertV16FeeQuoterDestUpdatesToV2(in)
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, uint64(100), out[0].DestChainSelector)
		assert.True(t, out[0].DestChainConfig.IsEnabled)
		assert.Equal(t, uint16(10), out[0].DestChainConfig.NetworkFeeUSDCents)
		assert.Equal(t, fqv2seq.LinkFeeMultiplierPercent, out[0].DestChainConfig.LinkFeeMultiplierPercent)
		assert.Equal(t, uint32(30000), out[0].DestChainConfig.MaxDataBytes)
	})

	t.Run("rejects NetworkFeeUSDCents exceeding uint16 max", func(t *testing.T) {
		in := []fee_quoter.FeeQuoterDestChainConfigArgs{
			{
				DestChainSelector: 200,
				DestChainConfig: fee_quoter.FeeQuoterDestChainConfig{
					NetworkFeeUSDCents: math.MaxUint16 + 1,
				},
			},
		}

		_, err := v1_6.ConvertV16FeeQuoterDestUpdatesToV2(in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "exceeds uint16 max")
	})
}

func TestConvertV16FeeQuoterPriceUpdatesToV2(t *testing.T) {
	t.Parallel()

	in := fee_quoter.InternalPriceUpdates{
		TokenPriceUpdates: []fee_quoter.InternalTokenPriceUpdate{
			{
				SourceToken: common.HexToAddress("0x1234"),
				UsdPerToken: big.NewInt(1e18),
			},
		},
		GasPriceUpdates: []fee_quoter.InternalGasPriceUpdate{
			{
				DestChainSelector: 100,
				UsdPerUnitGas:     big.NewInt(1e17),
			},
		},
	}

	out := v1_6.ConvertV16FeeQuoterPriceUpdatesToV2(in)
	require.Len(t, out.TokenPriceUpdates, 1)
	assert.Equal(t, common.HexToAddress("0x1234"), out.TokenPriceUpdates[0].SourceToken)
	assert.Equal(t, big.NewInt(1e18), out.TokenPriceUpdates[0].UsdPerToken)

	require.Len(t, out.GasPriceUpdates, 1)
	assert.Equal(t, uint64(100), out.GasPriceUpdates[0].DestChainSelector)
	assert.Equal(t, big.NewInt(1e17), out.GasPriceUpdates[0].UsdPerUnitGas)
}

func TestUpdateBidirectionalLanesChangeset_NonceManagerAutoDetect(t *testing.T) {
	t.Parallel()

	t.Run("no_v15_contracts_no_nonce_update", func(t *testing.T) {
		// Standard environment with no v1.5 contracts - NonceManager should NOT be updated
		tenv, _ := testhelpers.NewMemoryEnvironment(t)

		allChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
		require.Len(t, allChains, 2)
		chainA, chainB := allChains[0], allChains[1]

		// Apply UpdateBidirectionalLanesChangeset
		chains := []v1_6.ChainDefinition{
			{Selector: chainA, GasPrice: big.NewInt(1e9), FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true)},
			{Selector: chainB, GasPrice: big.NewInt(1e9), FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true)},
		}
		cfg := v1_6.UpdateBidirectionalLanesConfig{
			Lanes: getAllPossibleLanes(chains, false),
		}
		_, err := commonchangeset.Apply(t, tenv.Env,
			commonchangeset.Configure(v1_6.UpdateBidirectionalLanesChangeset, cfg))
		require.NoError(t, err)

		// Verify NonceManager does NOT have previous ramps set
		state, err := stateview.LoadOnchainState(tenv.Env)
		require.NoError(t, err)

		prevRamps, err := state.Chains[chainA].NonceManager.GetPreviousRamps(nil, chainB)
		require.NoError(t, err)
		require.Equal(t, common.Address{}, prevRamps.PrevOnRamp, "should NOT set PrevOnRamp when no v1.5 exists")
		require.Equal(t, common.Address{}, prevRamps.PrevOffRamp, "should NOT set PrevOffRamp when no v1.5 exists")
	})

	t.Run("skip_nonce_manager_updates_flag", func(t *testing.T) {
		// Even with setup that could trigger updates, the skip flag should prevent NonceManager changes
		tenv, _ := testhelpers.NewMemoryEnvironment(t)

		allChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
		chainA, chainB := allChains[0], allChains[1]

		chains := []v1_6.ChainDefinition{
			{Selector: chainA, GasPrice: big.NewInt(1e9), FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true)},
			{Selector: chainB, GasPrice: big.NewInt(1e9), FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true)},
		}
		cfg := v1_6.UpdateBidirectionalLanesConfig{
			Lanes:                   getAllPossibleLanes(chains, false),
			SkipNonceManagerUpdates: true, // OPT-OUT
		}
		_, err := commonchangeset.Apply(t, tenv.Env,
			commonchangeset.Configure(v1_6.UpdateBidirectionalLanesChangeset, cfg))
		require.NoError(t, err)

		// Verify NonceManager is NOT updated
		state, err := stateview.LoadOnchainState(tenv.Env)
		require.NoError(t, err)

		prevRamps, err := state.Chains[chainA].NonceManager.GetPreviousRamps(nil, chainB)
		require.NoError(t, err)
		require.Equal(t, common.Address{}, prevRamps.PrevOnRamp, "should NOT set when SkipNonceManagerUpdates=true")
		require.Equal(t, common.Address{}, prevRamps.PrevOffRamp, "should NOT set when SkipNonceManagerUpdates=true")
	})

	t.Run("idempotent_run_twice_no_error", func(t *testing.T) {
		// Running the changeset twice should succeed without error (idempotent)
		tenv, _ := testhelpers.NewMemoryEnvironment(t)

		allChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
		chainA, chainB := allChains[0], allChains[1]

		chains := []v1_6.ChainDefinition{
			{Selector: chainA, GasPrice: big.NewInt(1e9), FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true)},
			{Selector: chainB, GasPrice: big.NewInt(1e9), FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true)},
		}
		cfg := v1_6.UpdateBidirectionalLanesConfig{
			Lanes: getAllPossibleLanes(chains, false),
		}

		// First run
		var err error
		tenv.Env, err = commonchangeset.Apply(t, tenv.Env,
			commonchangeset.Configure(v1_6.UpdateBidirectionalLanesChangeset, cfg))
		require.NoError(t, err, "first run should succeed")

		// Second run - should also succeed (idempotent)
		_, err = commonchangeset.Apply(t, tenv.Env,
			commonchangeset.Configure(v1_6.UpdateBidirectionalLanesChangeset, cfg))
		require.NoError(t, err, "second run should succeed (idempotent)")
	})

	t.Run("v15_contracts_auto_detected_and_registered", func(t *testing.T) {
		// Environment with v1.5 contracts deployed - NonceManager SHOULD be updated with correct addresses
		e, tenv := testhelpers.NewMemoryEnvironment(
			t,
			testhelpers.WithPrerequisiteDeploymentOnly(&changeset.V1_5DeploymentConfig{
				PriceRegStalenessThreshold: 60 * 60 * 24 * 14,
				RMNConfig: &rmn_contract.RMNConfig{
					BlessWeightThreshold: 2,
					CurseWeightThreshold: 2,
					Voters: []rmn_contract.RMNVoter{
						{
							BlessWeight:   2,
							CurseWeight:   2,
							BlessVoteAddr: utils.RandomAddress(),
							CurseVoteAddr: utils.RandomAddress(),
						},
					},
				},
			}),
			testhelpers.WithEVMChainsBySelectors([]uint64{
				chain_selectors.GETH_TESTNET.Selector,
				chain_selectors.TEST_90000001.Selector,
			}),
		)

		allChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
		require.Len(t, allChains, 2)
		chainA, chainB := allChains[0], allChains[1]

		// Add CCIP v1.6 contracts
		e = testhelpers.AddCCIPContractsToEnvironment(t, allChains, tenv, false)

		// Load state to get v1.5 contracts
		state, err := stateview.LoadOnchainState(e.Env, stateview.WithLoadLegacyContracts(true))
		require.NoError(t, err)

		// Add v1.5 lanes
		pairs := []testhelpers.SourceDestPair{
			{SourceChainSelector: chainA, DestChainSelector: chainB},
			{SourceChainSelector: chainB, DestChainSelector: chainA},
		}
		e.Env = v1_5.AddLanes(t, e.Env, state, pairs)
		state, err = stateview.LoadOnchainState(e.Env, stateview.WithLoadLegacyContracts(true))
		require.NoError(t, err)

		// Verify v1.5 contracts exist before running the changeset
		require.NotNil(t, state.Chains[chainA].EVM2EVMOnRamp, "v1.5 OnRamps should exist")
		require.NotNil(t, state.Chains[chainA].EVM2EVMOnRamp[chainB], "v1.5 OnRamp A->B should exist")
		v15OnRampAtoB := state.Chains[chainA].EVM2EVMOnRamp[chainB].Address()
		require.NotEqual(t, common.Address{}, v15OnRampAtoB, "v1.5 OnRamp address should be non-zero")

		// Apply UpdateBidirectionalLanesChangeset
		chains := []v1_6.ChainDefinition{
			{Selector: chainA, GasPrice: big.NewInt(1e9), FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true)},
			{Selector: chainB, GasPrice: big.NewInt(1e9), FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true)},
		}
		cfg := v1_6.UpdateBidirectionalLanesConfig{
			Lanes:                   getAllPossibleLanes(chains, false),
			SkipNonceManagerUpdates: false, // Explicitly enable
		}
		_, err = commonchangeset.Apply(t, e.Env,
			commonchangeset.Configure(v1_6.UpdateBidirectionalLanesChangeset, cfg))
		require.NoError(t, err)

		// Verify NonceManager was updated with the correct v1.5 addresses
		state, err = stateview.LoadOnchainState(e.Env, stateview.WithLoadLegacyContracts(true))
		require.NoError(t, err)

		// Check NonceManager on chainA has previous ramps for chainB
		prevRampsA, err := state.Chains[chainA].NonceManager.GetPreviousRamps(nil, chainB)
		require.NoError(t, err)
		require.Equal(t, state.Chains[chainA].EVM2EVMOnRamp[chainB].Address(), prevRampsA.PrevOnRamp,
			"NonceManager on chainA should have correct v1.5 OnRamp for chainB")

		// Check NonceManager on chainB has previous ramps for chainA
		prevRampsB, err := state.Chains[chainB].NonceManager.GetPreviousRamps(nil, chainA)
		require.NoError(t, err)
		require.Equal(t, state.Chains[chainB].EVM2EVMOnRamp[chainA].Address(), prevRampsB.PrevOnRamp,
			"NonceManager on chainB should have correct v1.5 OnRamp for chainA")
	})
}

package solana_test

import (
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	solBinary "github.com/gagliardetto/binary"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	commonState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestDeployChainContractsChangesetSolana(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     1,
		SolChains:  1,
		Nodes:      4,
	})
	evmSelectors := e.AllChainSelectors()
	homeChainSel := evmSelectors[0]
	solChainSelectors := e.AllChainSelectorsSolana()
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	contractParams := make(map[uint64]v1_6.ChainContractParams)
	for _, chain := range e.AllChainSelectors() {
		cfg[chain] = proposalutils.SingleGroupTimelockConfigV2(t)
		contractParams[chain] = v1_6.ChainContractParams{
			FeeQuoterParams: v1_6.DefaultFeeQuoterParams(),
			OffRampParams:   v1_6.DefaultOffRampParams(),
		}
	}
	prereqCfg := make([]ccipChangeset.DeployPrerequisiteConfigPerChain, 0)
	for _, chain := range e.AllChainSelectors() {
		prereqCfg = append(prereqCfg, ccipChangeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
	}

	feeAggregatorPrivKey, _ := solana.NewRandomPrivateKey()
	feeAggregatorPubKey := feeAggregatorPrivKey.PublicKey()
	ci := os.Getenv("CI") == "true"
	if ci {
		testhelpers.SavePreloadedSolAddresses(t, e, solChainSelectors[0])
	}

	e, err = commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					testhelpers.TestNodeOperator: nodes.NonBootstraps().PeerIDs(),
				},
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
			e.AllChainSelectors(),
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
			e.AllChainSelectorsSolana(),
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			cfg,
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangeset.DeployPrerequisitesChangeset),
			ccipChangeset.DeployPrerequisiteConfig{
				Configs: prereqCfg,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(v1_6.DeployChainContractsChangeset),
			v1_6.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ContractParamsPerChain: contractParams,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			map[uint64]commontypes.MCMSWithTimelockConfigV2{
				solChainSelectors[0]: {
					Canceller:        proposalutils.SingleGroupMCMSV2(t),
					Proposer:         proposalutils.SingleGroupMCMSV2(t),
					Bypasser:         proposalutils.SingleGroupMCMSV2(t),
					TimelockMinDelay: big.NewInt(0),
				},
			},
		),
	)
	require.NoError(t, err)
	addresses, err := e.ExistingAddresses.AddressesForChain(solChainSelectors[0])
	require.NoError(t, err)
	mcmState, err := commonState.MaybeLoadMCMSWithTimelockChainStateSolana(e.SolChains[solChainSelectors[0]], addresses)
	require.NoError(t, err)

	// Fund signer PDAs for timelock and mcm
	// If we don't fund, execute() calls will fail with "no funds" errors.
	timelockSignerPDA := commonState.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	mcmSignerPDA := commonState.GetMCMSignerPDA(mcmState.McmProgram, mcmState.ProposerMcmSeed)
	memory.FundSolanaAccounts(e.GetContext(), t, []solana.PublicKey{timelockSignerPDA, mcmSignerPDA},
		100, e.SolChains[solChainSelectors[0]].Client)
	t.Logf("funded timelock signer PDA: %s", timelockSignerPDA.String())
	t.Logf("funded mcm signer PDA: %s", mcmSignerPDA.String())
	upgradeAuthority := timelockSignerPDA

	// we can't upgrade in place locally so we have to change where we build
	buildCs := commonchangeset.Configure(
		deployment.CreateLegacyChangeSet(ccipChangesetSolana.BuildSolanaChangeset),
		ccipChangesetSolana.BuildSolanaConfig{
			ChainSelector:       solChainSelectors[0],
			GitCommitSha:        "0863d8fed5fbada9f352f33c405e1753cbb7d72c",
			DestinationDir:      e.SolChains[solChainSelectors[0]].ProgramsPath,
			CleanDestinationDir: true,
		},
	)
	deployCs := commonchangeset.Configure(
		deployment.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
		ccipChangesetSolana.DeployChainContractsConfig{
			HomeChainSelector: homeChainSel,
			ContractParamsPerChain: map[uint64]ccipChangesetSolana.ChainContractParams{
				solChainSelectors[0]: {
					FeeQuoterParams: ccipChangesetSolana.FeeQuoterParams{
						DefaultMaxFeeJuelsPerMsg: solBinary.Uint128{Lo: 300000000, Hi: 0, Endianness: nil},
					},
					OffRampParams: ccipChangesetSolana.OffRampParams{
						EnableExecutionAfter: int64(globals.PermissionLessExecutionThreshold.Seconds()),
					},
				},
			},
		},
	)
	// set the fee aggregator address
	feeAggregatorCs := commonchangeset.Configure(
		deployment.CreateLegacyChangeSet(ccipChangesetSolana.SetFeeAggregator),
		ccipChangesetSolana.SetFeeAggregatorConfig{
			ChainSelector: solChainSelectors[0],
			FeeAggregator: feeAggregatorPubKey.String(),
		},
	)
	transferOwnershipCs := commonchangeset.Configure(
		deployment.CreateLegacyChangeSet(ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolana),
		ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolanaConfig{
			MinDelay: 1 * time.Second,
			ContractsByChain: map[uint64]ccipChangesetSolana.CCIPContractsToTransfer{
				solChainSelectors[0]: {
					Router:    true,
					FeeQuoter: true,
					OffRamp:   true,
				},
			},
		},
	)
	// make sure idempotency works and setting the upgrade authority
	upgradeAuthorityCs := commonchangeset.Configure(
		deployment.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
		ccipChangesetSolana.DeployChainContractsConfig{
			HomeChainSelector: homeChainSel,
			ContractParamsPerChain: map[uint64]ccipChangesetSolana.ChainContractParams{
				solChainSelectors[0]: {
					FeeQuoterParams: ccipChangesetSolana.FeeQuoterParams{
						DefaultMaxFeeJuelsPerMsg: solBinary.Uint128{Lo: 300000000, Hi: 0, Endianness: nil},
					},
					OffRampParams: ccipChangesetSolana.OffRampParams{
						EnableExecutionAfter: int64(globals.PermissionLessExecutionThreshold.Seconds()),
					},
				},
			},
			NewUpgradeAuthority: &upgradeAuthority,
		},
	)
	upgradeCs := commonchangeset.Configure(
		deployment.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
		ccipChangesetSolana.DeployChainContractsConfig{
			HomeChainSelector: homeChainSel,
			ContractParamsPerChain: map[uint64]ccipChangesetSolana.ChainContractParams{
				solChainSelectors[0]: {
					FeeQuoterParams: ccipChangesetSolana.FeeQuoterParams{
						DefaultMaxFeeJuelsPerMsg: solBinary.Uint128{Lo: 300000000, Hi: 0, Endianness: nil},
					},
					OffRampParams: ccipChangesetSolana.OffRampParams{
						EnableExecutionAfter: int64(globals.PermissionLessExecutionThreshold.Seconds()),
					},
				},
			},
			UpgradeConfig: ccipChangesetSolana.UpgradeConfig{
				NewFeeQuoterVersion: &deployment.Version1_1_0,
				NewRouterVersion:    &deployment.Version1_1_0,
				UpgradeAuthority:    upgradeAuthority,
				SpillAddress:        upgradeAuthority,
				MCMS: &ccipChangeset.MCMSConfig{
					MinDelay: 1 * time.Second,
				},
			},
		},
	)
	// because we cannot upgrade in place locally, we can't redeploy offramp
	offRampCs := commonchangeset.Configure(
		deployment.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
		ccipChangesetSolana.DeployChainContractsConfig{
			HomeChainSelector: homeChainSel,
			ContractParamsPerChain: map[uint64]ccipChangesetSolana.ChainContractParams{
				solChainSelectors[0]: {
					FeeQuoterParams: ccipChangesetSolana.FeeQuoterParams{
						DefaultMaxFeeJuelsPerMsg: solBinary.Uint128{Lo: 300000000, Hi: 0, Endianness: nil},
					},
					OffRampParams: ccipChangesetSolana.OffRampParams{
						EnableExecutionAfter: int64(globals.PermissionLessExecutionThreshold.Seconds()),
					},
				},
			},
			UpgradeConfig: ccipChangesetSolana.UpgradeConfig{
				NewOffRampVersion: &deployment.Version1_1_0,
			},
		},
	)
	if ci {
		e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
			deployCs,
			feeAggregatorCs,
			upgradeAuthorityCs,
			transferOwnershipCs,
		})
		require.NoError(t, err)
		state, err := ccipChangeset.LoadOnchainStateSolana(e)
		require.NoError(t, err)
		oldOffRampAddress := state.SolChains[solChainSelectors[0]].OffRamp
		// add a second offramp address
		e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
			buildCs,
			upgradeCs,
			offRampCs,
		})
		require.NoError(t, err)
		// verify the offramp address is different
		state, err = ccipChangeset.LoadOnchainStateSolana(e)
		require.NoError(t, err)
		newOffRampAddress := state.SolChains[solChainSelectors[0]].OffRamp
		require.NotEqual(t, oldOffRampAddress, newOffRampAddress)
	} else {
		e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
			buildCs,
			deployCs,
			feeAggregatorCs,
			upgradeAuthorityCs,
			upgradeCs,
		})
	}
	require.NoError(t, err)
	// Verify router and fee quoter upgraded in place
	// and offramp had 2nd address added
	addresses, err = e.ExistingAddresses.AddressesForChain(solChainSelectors[0])
	require.NoError(t, err)
	numRouters := 0
	numFeeQuoters := 0
	numOffRamps := 0
	for _, address := range addresses {
		if address.Type == ccipChangeset.Router {
			numRouters++
		}
		if address.Type == ccipChangeset.FeeQuoter {
			numFeeQuoters++
		}
		if address.Type == ccipChangeset.OffRamp {
			numOffRamps++
		}
	}
	require.Equal(t, 1, numRouters)
	require.Equal(t, 1, numFeeQuoters)
	if ci {
		require.Equal(t, 2, numOffRamps)
	} else {
		require.Equal(t, 1, numOffRamps)
	}
	require.NoError(t, err)
	// solana verification
	testhelpers.ValidateSolanaState(t, e, solChainSelectors)
}

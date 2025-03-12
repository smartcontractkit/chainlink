package crib

import (
	"context"
	"errors"

	solBinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func DeploySolanaHomeChainContracts(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig) (deployment.CapabilityRegistryConfig, deployment.AddressBook, error) {
	e, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	if e == nil {
		return deployment.CapabilityRegistryConfig{}, nil, errors.New("environment is nil")
	}

	evmSelectors := e.AllChainSelectors()
	homeChainSel := evmSelectors[0]
	solChainSelectors := e.AllChainSelectorsSolana()
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	contractParams := make(map[uint64]v1_6.ChainContractParams)
	for _, chain := range e.AllChainSelectors() {
		cfg[chain] = proposalutils.SingleGroupTimelockConfigV2(nil)
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

	solCfg := ccipChangesetSolana.BuildSolanaConfig{
		ChainSelector:       solChainSelectors[0],
		GitCommitSha:        "0863d8fed5fbada9f352f33c405e1753cbb7d72c",
		DestinationDir:      e.SolChains[solChainSelectors[0]].ProgramsPath,
		CleanDestinationDir: true,
		CleanGitDir:         true,
	}

	homeChainCfg := v1_6.DeployHomeChainConfig{
		HomeChainSel:     homeChainSel,
		RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
		RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
		NodeOperators:    testhelpers.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
		NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
			testhelpers.TestNodeOperator: nodes.NonBootstraps().PeerIDs(),
		},
	}

	homeChainChangeset, err := v1_6.DeployHomeChainChangeset(*e, homeChainCfg)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}

	linkTokenChangesetEvm, err := commonchangeset.DeployLinkToken(*e, evmSelectors)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}

	linkTokenChangesetSol, err := commonchangeset.DeployLinkToken(*e, solChainSelectors)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}

	timelockChangeset, err := commonchangeset.DeployMCMSWithTimelockV2(*e, cfg)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}

	prereqChangeset, err := ccipChangeset.DeployPrerequisitesChangeset(*e, ccipChangeset.DeployPrerequisiteConfig{
		Configs: prereqCfg,
	})
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}

	evmChainContractsChangeset, err := v1_6.DeployChainContractsChangeset(*e, v1_6.DeployChainContractsConfig{
		HomeChainSelector:      homeChainSel,
		ContractParamsPerChain: contractParams,
	})
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}

	solChainContractsChangeset, err := ccipChangesetSolana.DeployChainContractsChangeset(*e, ccipChangesetSolana.DeployChainContractsConfig{
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
	})
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}

	feeAggregatorChangeset, err := ccipChangesetSolana.SetFeeAggregator(*e, ccipChangesetSolana.SetFeeAggregatorConfig{
		ChainSelector: solChainSelectors[0],
		FeeAggregator: feeAggregatorPubKey.String(),
	},
	)

	state, err := changeset.LoadOnchainStateSolana(*e)

	_, err = ccipChangesetSolana.BuildSolanaChangeset(*e, solCfg)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	return deployment.CapabilityRegistryConfig{}, nil, nil
}

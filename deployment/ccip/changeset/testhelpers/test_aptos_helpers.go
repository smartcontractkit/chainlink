package testhelpers

import (
	"math/big"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func DeployChainContractsToAptosCS(t *testing.T, e DeployedEnv, chainSelector uint64) commonchangeset.ConfiguredChangeSet {
	// Set mock link token address on Address book (to skip deploying)
	err := e.Env.ExistingAddresses.Save(chainSelector, aptoscs.MockLinkAddress, cldf.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_6_0))
	require.NoError(t, err)

	linkTokenAddress := aptos.AccountAddress{}
	// TODO APT - replace with LINK once available
	linkTokenAddress.ParseStringRelaxed("0xa")
	ccipConfig := config.DeployAptosChainConfig{
		ContractParamsPerChain: map[uint64]config.ChainContractParams{
			chainSelector: {
				FeeQuoterParams: config.FeeQuoterParams{
					MaxFeeJuelsPerMsg:            1000000000,
					TokenPriceStalenessThreshold: 24 * 60 * 60,
					FeeTokens:                    []aptos.AccountAddress{linkTokenAddress},
				},
				OffRampParams: config.OffRampParams{
					ChainSelector:                    chainSelector,
					PermissionlessExecutionThreshold: uint32(globals.PermissionLessExecutionThreshold.Seconds()),
					IsRMNVerificationDisabled:        nil,
					SourceChainSelectors:             nil,
					SourceChainIsEnabled:             nil,
					SourceChainsOnRamp:               nil,
				},
				OnRampParams: config.OnRampParams{
					ChainSelector:  chainSelector,
					AllowlistAdmin: e.Env.BlockChains.AptosChains()[chainSelector].DeployerSigner.AccountAddress(),
					FeeAggregator:  e.Env.BlockChains.AptosChains()[chainSelector].DeployerSigner.AccountAddress(),
				},
			},
		},
		MCMSDeployConfigPerChain: map[uint64]commontypes.MCMSWithTimelockConfigV2{
			chainSelector: {
				Canceller:        proposalutils.SingleGroupMCMSV2(t),
				Proposer:         proposalutils.SingleGroupMCMSV2(t),
				Bypasser:         proposalutils.SingleGroupMCMSV2(t),
				TimelockMinDelay: big.NewInt(1),
			},
		},
	}

	return commonchangeset.Configure(aptoscs.DeployAptosChain{}, ccipConfig)
}

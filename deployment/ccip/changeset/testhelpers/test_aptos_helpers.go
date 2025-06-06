package testhelpers

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func DeployChainContractsToAptosCS(t *testing.T, e DeployedEnv, chainSelector uint64) commonchangeset.ConfiguredChangeSet {
	// Set mock link token address on Address book (to skip deploying)
	err := e.Env.ExistingAddresses.Save(chainSelector, aptoscs.MockLinkAddress, cldf.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_6_0))
	require.NoError(t, err)

	//  Deploy contracts
	mockCCIPParams := aptoscs.GetMockChainContractParams(t, chainSelector)
	ccipConfig := config.DeployAptosChainConfig{
		ContractParamsPerChain: map[uint64]config.ChainContractParams{
			chainSelector: mockCCIPParams,
		},
		MCMSDeployConfigPerChain: map[uint64]commontypes.MCMSWithTimelockConfigV2{
			chainSelector: {
				Canceller:        proposalutils.SingleGroupMCMSV2(t),
				Proposer:         proposalutils.SingleGroupMCMSV2(t),
				Bypasser:         proposalutils.SingleGroupMCMSV2(t),
				TimelockMinDelay: big.NewInt(0),
			},
		},
	}

	return commonchangeset.Configure(aptoscs.DeployAptosChain{}, ccipConfig)
}

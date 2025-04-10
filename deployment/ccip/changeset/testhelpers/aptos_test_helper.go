package testhelpers

import (
	"testing"

	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	aptoshelper "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

func DeployChainContractsToAptosCS(t *testing.T, e DeployedEnv, chainSelector uint64) commonchangeset.ConfiguredChangeSet {
	mockCCIPParams := aptoshelper.GetMockChainContractParams(t, chainSelector)
	ccipConfig := config.DeployAptosChainConfig{
		ContractParamsPerChain: map[uint64]config.ChainContractParams{
			chainSelector: mockCCIPParams,
		},
		MCMSConfigPerChain: map[uint64]mcmstypes.Config{
			chainSelector: proposalutils.SingleGroupMCMSV2(t),
		},
	}

	return commonchangeset.Configure(aptoscs.DeployAptosChain{}, ccipConfig)
}

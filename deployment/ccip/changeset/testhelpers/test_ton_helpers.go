package testhelpers

import (
	"testing"

	"github.com/xssnick/tonutils-go/address"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"

	toncs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/ton"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func DeployChainContractsToTonCS(t *testing.T, e DeployedEnv, chainSelector uint64) commoncs.ConfiguredChangeSet {
	// TODO(ton): Implement this function to deploy chain contracts to Ton chain, https://smartcontract-it.atlassian.net/browse/NONEVM-1938
	ccipConfig := toncs.DeployTONChainConfig{}
	return commoncs.Configure(toncs.DeployTONChain{}, ccipConfig)
}

func AddLaneTONChangesets(e *DeployedEnv, from, to uint64, fromFamily, toFamily string) commoncs.ConfiguredChangeSet {
	laneConfig := toncs.AddLaneTONChainConfig{
		FromChainSelector: from,
		ToChainSelector:   to,
		FromFamily:        fromFamily,
		ToFamily:          toFamily,
	}
	return commoncs.Configure(toncs.AddLaneTONChain{}, laneConfig)
}

// TODO: move to ton package
func ConfirmCommitWithExpectedSeqNumRangeTon(
	t *testing.T,
	srcSelector uint64,
	dest cldf_ton.Chain,
	offrampAddr address.Address,
	startBlock uint64,
	expectedSeqNumRange ccipocr3.SeqNumRange,
	enforceSingleCommit bool,
) (any, error) {
	t.Logf("DEBUG: ConfirmCommitWithExpectedSeqNumRangeTon srcSelector: %d, startBlock: %+v, expectedSeqNumRange: %+v, enforceSingleCommit: %+v\n", srcSelector, startBlock, expectedSeqNumRange, enforceSingleCommit)
	// TODO once offramp contracts are supported, we can add the logic to confirm commit with expected sequence number range
	return true, nil
}

func ConfirmExecWithSeqNrsTon(
	t *testing.T,
	sourceChain uint64,
	dest cldf_ton.Chain,
	offRampAddress address.Address,
	startBlock *uint64,
	expectedSeqNrs []uint64,
) (executionStates map[uint64]int, err error) {
	t.Logf("DEBUG: ConfirmExecWithSeqNrsTon srcSelector: %d, dest: %s, startBlock: %+v, expectedSeqNrs: %+v\n", sourceChain, startBlock, expectedSeqNrs)
	// TODO once offramp contracts are supported, we can add the logic to confirm execution with sequence numbers
	t.Logf("DEBUG: TODO(ton): ConfirmExecWithSeqNrsTon\n")
	seqNrsToWatch := make(map[uint64]int)
	for _, seqNr := range expectedSeqNrs {
		seqNrsToWatch[seqNr] = EXECUTION_STATE_SUCCESS
	}
	return seqNrsToWatch, nil
}

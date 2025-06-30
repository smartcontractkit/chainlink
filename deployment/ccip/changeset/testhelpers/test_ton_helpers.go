package testhelpers

import (
	"testing"

	toncs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/ton"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func DeployChainContractsToTonCS(t *testing.T, e DeployedEnv, chainSelector uint64) commoncs.ConfiguredChangeSet {
	// TODO(ton): Implement this function to deploy chain contracts to Ton chain, https://smartcontract-it.atlassian.net/browse/NONEVM-1938
	ccipConfig := toncs.DeployCCIPContractsCfg{
		TonChainSelector: chainSelector,
	}
	return commoncs.Configure(toncs.DeployCCIPContracts{}, ccipConfig)
}

func AddLaneTONChangesets(e *DeployedEnv, from, to uint64, fromFamily, toFamily string) commoncs.ConfiguredChangeSet {
	laneConfig := toncs.AddLaneCfg{
		FromChainSelector: from,
		ToChainSelector:   to,
		FromFamily:        fromFamily,
		ToFamily:          toFamily,
	}
	return commoncs.Configure(toncs.AddLane{}, laneConfig)
}

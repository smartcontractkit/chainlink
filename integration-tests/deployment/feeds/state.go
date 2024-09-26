package feeds

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

// Can serialize to RDD style data
type FeedsOnChainState struct {
	Feeds map[common.Address]*offchainaggregator.OffchainAggregator
}

func LoadFeedsChainState(chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (FeedsOnChainState, error) {
	feeds := make(map[common.Address]*offchainaggregator.OffchainAggregator)
	for address, tvStr := range addresses {
		switch tvStr.String() {
		case deployment.NewTypeAndVersion("OffchainAggregator", deployment.Version1_0_0).String():
			addr := common.HexToAddress(address)
			aggr, _ := offchainaggregator.NewOffchainAggregator(common.HexToAddress(address), chain.Client)
			feeds[addr] = aggr
		default:
			panic("unknown contract type and version")
		}
	}
	return FeedsOnChainState{
		Feeds: feeds,
	}, nil
}

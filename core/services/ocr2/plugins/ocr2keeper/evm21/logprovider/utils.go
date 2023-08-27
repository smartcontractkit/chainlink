package logprovider

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// FilterLogsByContent filters logs based on the filter selector and topics 1/2/3
func FilterLogsByContent(filter upkeepFilter, logs []logpoller.Log, lggr logger.Logger) []logpoller.Log {
	if filter.selector == 0 {
		return logs
	}
	var filteredLogs []logpoller.Log
	// filter logs based on filter selector (topic1, topic2, topic3) if it's configured
	checkTopic1 := filter.selector%2 == 1
	checkTopic2 := (filter.selector>>1)%2 == 1
	checkTopic3 := (filter.selector>>2)%2 == 1
	for _, l := range logs {
		if checkTopic1 && common.BytesToHash(l.Topics[1]).Hex() != filter.topics[1].Hex() {
			lggr.Debugf("upkeep Id %s topics[1] %s does not match log topics[1] %s", filter.upkeepID, filter.topics[1].Hex(), common.BytesToHash(l.Topics[1]).Hex())
			continue
		}
		if checkTopic2 && common.BytesToHash(l.Topics[2]).Hex() != filter.topics[2].Hex() {
			lggr.Debugf("upkeep Id %s topics[2] %s does not match log topics[2] %s", filter.upkeepID, filter.topics[2].Hex(), common.BytesToHash(l.Topics[2]).Hex())
			continue
		}
		if checkTopic3 && common.BytesToHash(l.Topics[3]).Hex() != filter.topics[3].Hex() {
			lggr.Debugf("upkeep Id %s topics[3] %s does not match log topics[3] %s", filter.upkeepID, filter.topics[3].Hex(), common.BytesToHash(l.Topics[3]).Hex())
			continue
		}
		lggr.Infof("upkeep Id %s topics matches log %s", filter.upkeepID, l.TxHash.Hex())
		filteredLogs = append(filteredLogs, l)
	}
	return filteredLogs
}

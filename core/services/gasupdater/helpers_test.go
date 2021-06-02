package gasupdater

import "github.com/smartcontractkit/chainlink/core/services/headtracker"

func GasUpdaterToStruct(gu GasUpdater) *gasUpdater {
	return gu.(*gasUpdater)
}
func SetRollingBlockHistory(gu GasUpdater, blocks []headtracker.Block) {
	gu.(*gasUpdater).rollingBlockHistory = blocks
}

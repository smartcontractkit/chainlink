package gasupdater

import "github.com/smartcontractkit/chainlink/core/store/models"

func GasUpdaterToStruct(gu GasUpdater) *gasUpdater {
	return gu.(*gasUpdater)
}
func SetRollingBlockHistory(gu GasUpdater, blocks []models.Block) {
	gu.(*gasUpdater).rollingBlockHistory = blocks
}

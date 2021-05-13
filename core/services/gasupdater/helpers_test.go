package gasupdater

func GasUpdaterToStruct(gu GasUpdater) *gasUpdater {
	return gu.(*gasUpdater)
}
func SetRollingBlockHistory(gu GasUpdater, blocks []Block) {
	gu.(*gasUpdater).rollingBlockHistory = blocks
}

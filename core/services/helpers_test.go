package services

func (ht *HeadTracker) ExportedDone() chan struct{} {
	return ht.done
}

func GasUpdaterToStruct(gu GasUpdater) *gasUpdater {
	return gu.(*gasUpdater)
}

func SetRollingBlockHistory(gu GasUpdater, blocks []Block) {
	gu.(*gasUpdater).rollingBlockHistory = blocks
}

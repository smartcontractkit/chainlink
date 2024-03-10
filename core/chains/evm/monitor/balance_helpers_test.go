package monitor

func (bm *balanceMonitor) WorkDone() <-chan struct{} {
	return bm.sleeperTask.WorkDone()
}

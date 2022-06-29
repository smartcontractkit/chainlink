package utils

func (once *StartStopOnce) LoadState() StartStopOnceState {
	return StartStopOnceState(once.state.Load())
}

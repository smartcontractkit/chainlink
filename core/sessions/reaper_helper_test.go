package sessions

func (sr *sessionReaper) RunSignal() <-chan struct{} {
	return sr.runSignal
}

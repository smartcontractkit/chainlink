package telemetry

type NoopAgent struct {
}

// SendLog sends a telemetry log to the explorer
func (t *NoopAgent) SendLog(log []byte) {
}

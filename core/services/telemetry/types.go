package telemetry

// FIXME: remove this

// MonitoringEndpoint is where the OCR protocol sends monitoring output
type MonitoringEndpoint interface {
	SendLog(log []byte)
}

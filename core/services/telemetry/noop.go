package telemetry

import (
	"os"

	ocrtypes "github.com/smartcontractkit/libocr/commontypes"
)

var _ MonitoringEndpointGenerator = &NoopAgent{}

type NoopAgent struct {
}

// SendLog sends a telemetry log to the explorer
func (t *NoopAgent) SendLog(log []byte) {
	if os.Getenv("TELEMETRY_INGRESS_LOG_FILE") != "" {
		f, err := os.Create("telemetry.log")
		if err != nil {
			return
		}
		defer f.Close()
		_, err = f.Write(log)
		if err != nil {
			return
		}
	}
}

// GenMonitoringEndpoint creates a monitoring endpoint for telemetry
func (t *NoopAgent) GenMonitoringEndpoint(contractID string) ocrtypes.MonitoringEndpoint {
	return t
}

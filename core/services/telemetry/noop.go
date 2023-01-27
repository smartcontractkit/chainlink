package telemetry

import (
	"encoding/base64"
	"os"

	ocrtypes "github.com/smartcontractkit/libocr/commontypes"
)

var _ MonitoringEndpointGenerator = &NoopAgent{}

type NoopAgent struct {
}

// SendLog sends a telemetry log to the explorer
func (t *NoopAgent) SendLog(log []byte) {
	if os.Getenv("TELEMETRY_INGRESS_LOG_FILE") != "" {
		f, err := os.OpenFile("telemetry.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
		if err != nil {
			return
		}
		defer f.Close()
		_, err = f.WriteString(base64.StdEncoding.EncodeToString(log) + "\n")
		if err != nil {
			return
		}
	}
}

// GenMonitoringEndpoint creates a monitoring endpoint for telemetry
func (t *NoopAgent) GenMonitoringEndpoint(contractID string) ocrtypes.MonitoringEndpoint {
	return t
}

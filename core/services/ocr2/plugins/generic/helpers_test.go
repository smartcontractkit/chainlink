package generic

import "github.com/smartcontractkit/libocr/commontypes"

func (t *TelemetryAdapter) Endpoints() map[[4]string]commontypes.MonitoringEndpoint {
	return t.endpoints
}

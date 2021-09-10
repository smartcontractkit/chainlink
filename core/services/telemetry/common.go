package telemetry

import (
	"github.com/ethereum/go-ethereum/common"
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
)

type MonitoringEndpointGenerator interface {
	GenMonitoringEndpoint(addr common.Address) ocrcommontypes.MonitoringEndpoint
}

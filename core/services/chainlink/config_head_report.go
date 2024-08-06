package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type headReport struct {
	h toml.HeadReport
}

func (h headReport) TelemetryEnabled() bool {
	return *h.h.TelemetryEnabled
}

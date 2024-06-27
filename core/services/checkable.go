package services

import (
	"github.com/smartcontractkit/chainlink-relay/pkg/services"
)

// Deprecated: use services.HealthReporter
type Checkable = services.HealthReporter

// Deprecated: use services.CopyHealth
func CopyHealth(dest, src map[string]error) { services.CopyHealth(dest, src) }

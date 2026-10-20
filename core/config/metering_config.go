package config

import "github.com/smartcontractkit/chainlink-common/pkg/resourcemanager"

// Metering exposes durable resource-metering configuration: the emission
// toggles and the coarse deployment/node identity dimensions stamped on emitted
// MeterRecords and MeterSnapshots. These are passed via loop.EnvConfig to every LOOP
// plugin.
type Metering interface {
	MeterRecordsEnabled() bool
	MeterSnapshotsEnabled() bool
	Product() string
	Tenant() string
	NumericTenantID() string
	Environment() string
	Zone() string
	NodeID() string
	// FaultInjection returns the usage record fault injection settings (off by default).
	FaultInjection() resourcemanager.UsageFaultInjection
}

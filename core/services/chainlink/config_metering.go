package chainlink

import (
	"github.com/smartcontractkit/chainlink-common/pkg/resourcemanager"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type meteringConfig struct {
	s toml.Metering
}

func (b *meteringConfig) MeterRecordsEnabled() bool {
	if b.s.MeterRecordsEnabled == nil {
		return false
	}
	return *b.s.MeterRecordsEnabled
}

func (b *meteringConfig) MeterSnapshotsEnabled() bool {
	if b.s.MeterSnapshotsEnabled == nil {
		return false
	}
	return *b.s.MeterSnapshotsEnabled
}

// Product returns the deployment product identity dimension. The parsed config
// defaults it to "cre" via docs.CoreDefaults so metering is never enabled with
// an empty product dimension; a zero-value toml.Metering that has not been run
// through setDefaults returns UnsetProduct (the nil-pointer fallback below).
func (b *meteringConfig) Product() string {
	if b.s.Product == nil {
		return resourcemanager.UnsetProduct
	}
	return *b.s.Product
}

func (b *meteringConfig) Tenant() string {
	if b.s.Tenant == nil {
		return ""
	}
	return *b.s.Tenant
}

func (b *meteringConfig) NumericTenantID() string {
	if b.s.NumericTenantID == nil {
		return ""
	}
	return *b.s.NumericTenantID
}

func (b *meteringConfig) Environment() string {
	if b.s.Environment == nil {
		return ""
	}
	return *b.s.Environment
}

func (b *meteringConfig) Zone() string {
	if b.s.Zone == nil {
		return ""
	}
	return *b.s.Zone
}

func (b *meteringConfig) NodeID() string {
	if b.s.NodeID == nil {
		return ""
	}
	return *b.s.NodeID
}

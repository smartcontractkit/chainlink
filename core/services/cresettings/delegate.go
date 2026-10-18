// cresettings jobs are used to distribute updates for CRE settings overrides.
// See: https://pkg.go.dev/github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings
// Only one Job of type CRESettings may run at a time per config type. Attempts to create a second job of the same config type will fail.
package cresettings

import (
	"context"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

const (
	ConfigTypeSettings        = "settings"
	ConfigTypeShardAssignment = "shard_assignment"
	ConfigTypeCapRegistry     = "capabilities_registry"

	// offchainConfigKey is the settings TOML key holding the capabilities_registry proto-JSON payload.
	offchainConfigKey = "offchain_config"
)

func NewDelegate(lggr logger.Logger, atomicSettings *loop.AtomicSettings, shardAssignmentSettings *loop.AtomicSettings, capRegistry *globalconfig.GlobalConfig) *delegate {
	return &delegate{
		lggr:                    lggr,
		atomicSettings:          atomicSettings,
		shardAssignmentSettings: shardAssignmentSettings,
		capRegistry:             capRegistry,
		capRegMetrics:           globalconfig.DefaultMetrics(),
	}
}

var _ job.Delegate = (*delegate)(nil)

type delegate struct {
	lggr                    logger.Logger
	atomicSettings          *loop.AtomicSettings
	shardAssignmentSettings *loop.AtomicSettings
	capRegistry             *globalconfig.GlobalConfig
	capRegMetrics           *globalconfig.Metrics

	// activeJobIDs maps config type to the active job ID.
	activeJobIDs sync.Map
}

func (d *delegate) JobType() job.Type {
	return job.CRESettings
}

func (d *delegate) BeforeJobCreated(j job.Job) {}

func (d *delegate) ServicesForSpec(ctx context.Context, j job.Job) ([]job.ServiceCtx, error) {
	spec := j.CRESettingsSpec
	configType := d.configType(spec)
	switch configType {
	case ConfigTypeSettings, ConfigTypeShardAssignment, ConfigTypeCapRegistry:
	default:
		return nil, fmt.Errorf("unknown config_type %q", configType)
	}

	if activeJobID, loaded := d.activeJobIDs.LoadOrStore(configType, j.ID); loaded {
		return nil, fmt.Errorf("another %s job with config_type %q is already active: %d", job.CRESettings, configType, activeJobID.(int32))
	}

	switch configType {
	case ConfigTypeCapRegistry:
		var version uint64
		payload, err := extractOffchainConfig(spec.Settings)
		if err == nil {
			version, err = d.capRegistry.Store(payload, spec.Hash)
		}
		domain, env := globalconfig.PayloadLabels(payload)
		if err != nil {
			d.capRegMetrics.RecordValidationError(ctx, domain, env)
			d.activeJobIDs.CompareAndDelete(configType, j.ID)
			return nil, fmt.Errorf("failed to store capabilities registry config: %w", err)
		}
		d.capRegMetrics.RecordAppliedVersion(ctx, domain, env, version)
		d.lggr.Infow("Updated capabilities registry config", "hash", spec.Hash, "offchainVersion", version)

	case ConfigTypeShardAssignment:
		if err := d.shardAssignmentSettings.Store(core.SettingsUpdate{
			Settings: spec.Settings,
			Hash:     spec.Hash,
		}); err != nil {
			return nil, fmt.Errorf("failed to store shard assignment settings: %w", err)
		}
		d.lggr.Infow("Updated shard assignment config", "hash", spec.Hash)

	case ConfigTypeSettings:
		if err := d.atomicSettings.Store(core.SettingsUpdate{
			Settings: spec.Settings,
			Hash:     spec.Hash,
		}); err != nil {
			return nil, fmt.Errorf("failed to update settings: %w", err)
		}
		d.lggr.Infow("Updated settings", "hash", spec.Hash, "settings", spec.Settings)
	}

	return nil, nil
}

func (d *delegate) AfterJobCreated(j job.Job) {}

func (d *delegate) BeforeJobDeleted(j job.Job) {}

func (d *delegate) OnDeleteJob(ctx context.Context, jb job.Job) error {
	configType := d.configType(jb.CRESettingsSpec)
	// A deleted capabilities_registry job keeps its payload applied: JD replaces a job by deleting
	// and re-creating it, and clearing in between would briefly launch capabilities without config.
	if !d.activeJobIDs.CompareAndDelete(configType, jb.ID) {
		d.lggr.Errorf("job %d was not the active %s job for config_type %q", jb.ID, job.CRESettings, configType)
	}
	return nil
}

func (d *delegate) configType(spec *job.CRESettingsSpec) string {
	if spec == nil || spec.Settings == "" {
		return ConfigTypeSettings
	}
	if ct, ok := extractConfigType(spec.Settings); ok {
		return ct
	}
	d.lggr.Infow("No config_type specified, defaulting to settings")
	return ConfigTypeSettings
}

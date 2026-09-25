// cresettings jobs are used to distribute updates for CRE settings overrides.
// See: https://pkg.go.dev/github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings
// At most one CRESettings job per config_type may run at a time; a second job of the same
// config_type will fail. Different config_types (settings, shard_assignment,
// capabilities_registry) coexist.
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

func NewDelegate(lggr logger.Logger, atomicSettings *loop.AtomicSettings, shardAssignmentSettings *loop.AtomicSettings, globalConfig *globalconfig.GlobalConfig) *delegate {
	return &delegate{
		lggr:                    lggr,
		atomicSettings:          atomicSettings,
		shardAssignmentSettings: shardAssignmentSettings,
		globalConfig:            globalConfig,
	}
}

var _ job.Delegate = (*delegate)(nil)

type delegate struct {
	lggr                    logger.Logger
	atomicSettings          *loop.AtomicSettings
	shardAssignmentSettings *loop.AtomicSettings
	globalConfig            *globalconfig.GlobalConfig

	// activeJobIDs tracks the active job ID per config_type, so one job of each config_type
	// can run concurrently while rejecting a second job of the same config_type.
	activeJobIDs sync.Map // configType(string) -> jobID(int32)
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
	case ConfigTypeShardAssignment:
		if err := d.shardAssignmentSettings.Store(core.SettingsUpdate{
			Settings: spec.Settings,
			Hash:     spec.Hash,
		}); err != nil {
			return nil, fmt.Errorf("failed to store shard assignment settings: %w", err)
		}
		d.lggr.Infow("Updated shard assignment config", "hash", spec.Hash)

	case ConfigTypeCapRegistry:
		if d.globalConfig == nil {
			return nil, fmt.Errorf("no global config store configured for config_type %q", configType)
		}
		if err := d.globalConfig.Store(globalconfig.Update{
			Raw:  spec.OffchainConfig,
			Hash: spec.Hash,
		}); err != nil {
			return nil, fmt.Errorf("failed to store offchain capabilities registry config: %w", err)
		}
		d.lggr.Infow("Updated offchain capabilities registry config", "hash", spec.Hash)

	case ConfigTypeSettings:
		if err := d.atomicSettings.Store(core.SettingsUpdate{
			Settings: spec.Settings,
			Hash:     spec.Hash,
		}); err != nil {
			return nil, fmt.Errorf("failed to update settings: %w", err)
		}
		d.lggr.Infow("Updated settings", "hash", spec.Hash, "settings", spec.Settings)
	}

	d.activeJobIDs.Store(configType, j.ID)
	return nil, nil
}

func (d *delegate) AfterJobCreated(j job.Job) {}

func (d *delegate) BeforeJobDeleted(j job.Job) {}

func (d *delegate) OnDeleteJob(ctx context.Context, jb job.Job) error {
	configType := d.configType(jb.CRESettingsSpec)
	if !d.activeJobIDs.CompareAndDelete(configType, jb.ID) {
		d.lggr.Errorf("job %d was not the active %s job for config_type %q", jb.ID, job.CRESettings, configType)
	}
	return nil
}

func (d *delegate) configType(spec *job.CRESettingsSpec) string {
	if spec == nil {
		return ConfigTypeSettings
	}
	// Must match validate.go's resolveConfigType: the top-level ConfigType field wins (used by
	// capabilities_registry, whose payload lives in OffchainConfig with empty Settings), then a
	// config_type key embedded in Settings (shard_assignment), else the default settings.
	return resolveConfigType(*spec)
}

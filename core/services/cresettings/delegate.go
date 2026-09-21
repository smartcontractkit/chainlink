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

	// activeJobs tracks the active job ID per config_type, so one job of each config_type
	// can run concurrently while rejecting a second job of the same config_type.
	activeJobs sync.Map // configType(string) -> jobID(int32)
}

func (d *delegate) JobType() job.Type {
	return job.CRESettings
}

func (d *delegate) BeforeJobCreated(j job.Job) {}

func (d *delegate) ServicesForSpec(ctx context.Context, j job.Job) ([]job.ServiceCtx, error) {
	spec := j.CRESettingsSpec
	configType := resolveConfigType(*spec)

	if existing, ok := d.activeJobs.Load(configType); ok && existing.(int32) != j.ID {
		return nil, fmt.Errorf("another %s job with config_type %q is already active: %d", job.CRESettings, configType, existing)
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

	default:
		return nil, fmt.Errorf("unknown config_type %q", configType)
	}

	d.activeJobs.Store(configType, j.ID)
	return nil, nil
}

func (d *delegate) AfterJobCreated(j job.Job) {}

func (d *delegate) BeforeJobDeleted(j job.Job) {}

func (d *delegate) OnDeleteJob(ctx context.Context, jb job.Job) error {
	d.activeJobs.Range(func(k, v any) bool {
		if v.(int32) == jb.ID {
			d.activeJobs.Delete(k)
			return false
		}
		return true
	})
	return nil
}

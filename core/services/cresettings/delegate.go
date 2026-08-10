package cresettings

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func NewDelegate(lggr logger.Logger, atomicSettings *loop.AtomicSettings, shardAssignment *atomic.Pointer[ShardAssignmentConfig]) *delegate {
	return &delegate{
		lggr:            lggr,
		atomicSettings:  atomicSettings,
		shardAssignment: shardAssignment,
	}
}

var _ job.Delegate = (*delegate)(nil)

type delegate struct {
	lggr            logger.Logger
	atomicSettings  *loop.AtomicSettings
	shardAssignment *atomic.Pointer[ShardAssignmentConfig]

	activeJobID atomic.Pointer[int32]
}

func (d *delegate) JobType() job.Type {
	return job.CRESettings
}

func (d *delegate) BeforeJobCreated(j job.Job) {}

func (d *delegate) ServicesForSpec(ctx context.Context, j job.Job) ([]job.ServiceCtx, error) {
	if activeJobID := d.activeJobID.Load(); activeJobID != nil {
		return nil, fmt.Errorf("another %s job is already active: %d", job.CRESettings, activeJobID)
	}

	spec := j.CRESettingsSpec
	configType := spec.ConfigType
	if configType == "" {
		configType = ConfigTypeSettings
	}

	switch configType {
	case ConfigTypeShardAssignment:
		saCfg, err := ParseShardAssignmentConfig(spec.ShardAssignment)
		if err != nil {
			return nil, fmt.Errorf("failed to parse shard_assignment config: %w", err)
		}
		d.shardAssignment.Store(saCfg)
		d.lggr.Infow("Updated shard assignment config", "hash", spec.Hash)

	case ConfigTypeSettings:
		fallthrough
	default:
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
	if !d.activeJobID.CompareAndSwap(&jb.ID, nil) {
		d.lggr.Errorf("job %d was not active", jb.ID)
	}
	return nil
}

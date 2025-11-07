// cresettings jobs are used to distribute updates for CRE settings overrides.
// See: https://pkg.go.dev/github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings
// Only one Job of type CRESettings may run at a time. Attempts to create a second job will fail.
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

func NewDelegate(lggr logger.Logger, atomicSettings *loop.AtomicSettings) *delegate {
	return &delegate{lggr: lggr, atomicSettings: atomicSettings}
}

var _ job.Delegate = (*delegate)(nil)

// delegate manages job.CRESettings jobs.
// It ensures that only one job is created at a time, and broadcasts changes via loop.AtomicSettings.
type delegate struct {
	lggr           logger.Logger
	atomicSettings *loop.AtomicSettings

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

	if err := d.atomicSettings.Store(core.SettingsUpdate{
		Settings: j.CRESettingsSpec.Settings,
		Hash:     j.CRESettingsSpec.Hash,
	}); err != nil {
		return nil, fmt.Errorf("failed to update settings: %w", err)
	}
	d.lggr.Info("Updated settings", "hash", j.CRESettingsSpec.Hash)

	return nil, nil // no active services
}

func (d *delegate) AfterJobCreated(j job.Job) {}

func (d *delegate) BeforeJobDeleted(j job.Job) {}

func (d *delegate) OnDeleteJob(ctx context.Context, jb job.Job) error {
	if !d.activeJobID.CompareAndSwap(&jb.ID, nil) {
		d.lggr.Errorf("job %d was not active", jb.ID)
	}
	return nil
}

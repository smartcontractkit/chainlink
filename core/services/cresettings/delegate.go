// cresettings jobs are used to distribute updates for CRE settings overrides.
// See: https://pkg.go.dev/github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings
// At most one CRESettings job per config_type may run at a time; a second job of the same
// config_type will fail. Different config_types (settings, shard_assignment,
// capabilities_registry) coexist. For capabilities_registry, uniqueness and version
// monotonicity are enforced by the database and the runtime config is projected from committed
// state (see CapRegistryProjector).
package cresettings

import (
	"context"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func NewDelegate(lggr logger.Logger, atomicSettings *loop.AtomicSettings, shardAssignmentSettings *loop.AtomicSettings, capRegistry *CapRegistryProjector) *delegate {
	return &delegate{
		lggr:                    lggr,
		atomicSettings:          atomicSettings,
		shardAssignmentSettings: shardAssignmentSettings,
		capRegistry:             capRegistry,
	}
}

var _ job.Delegate = (*delegate)(nil)

type delegate struct {
	lggr                    logger.Logger
	atomicSettings          *loop.AtomicSettings
	shardAssignmentSettings *loop.AtomicSettings
	capRegistry             *CapRegistryProjector

	// activeJobIDs tracks the active job ID per settings-based config_type, so one job of each
	// can run concurrently while rejecting a second job of the same config_type.
	// capabilities_registry does not use it: the database enforces a single job.
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

	if configType == ConfigTypeCapRegistry {
		return nil, d.syncCapRegistry(ctx)
	}

	if activeJobID, loaded := d.activeJobIDs.LoadOrStore(configType, j.ID); loaded {
		return nil, fmt.Errorf("another %s job with config_type %q is already active: %d", job.CRESettings, configType, activeJobID.(int32))
	}

	if err := d.apply(configType, spec); err != nil {
		// Release the slot claimed above: a job whose payload was rejected is not active, and
		// must not block a later valid job of the same config_type. The feeds manager deletes
		// the previous job before creating its replacement, so a reserved slot here would wedge
		// the config_type until restart.
		d.activeJobIDs.CompareAndDelete(configType, j.ID)
		return nil, err
	}
	return nil, nil
}

// apply stores the spec's payload into the store for its config_type.
func (d *delegate) apply(configType string, spec *job.CRESettingsSpec) error {
	switch configType {
	case ConfigTypeShardAssignment:
		if err := d.shardAssignmentSettings.Store(core.SettingsUpdate{
			Settings: spec.Settings,
			Hash:     spec.Hash,
		}); err != nil {
			return fmt.Errorf("failed to store shard assignment settings: %w", err)
		}
		d.lggr.Infow("Updated shard assignment config", "hash", spec.Hash)

	case ConfigTypeSettings:
		if err := d.atomicSettings.Store(core.SettingsUpdate{
			Settings: spec.Settings,
			Hash:     spec.Hash,
		}); err != nil {
			return fmt.Errorf("failed to update settings: %w", err)
		}
		d.lggr.Infow("Updated settings", "hash", spec.Hash, "settings", spec.Settings)
	}
	return nil
}

// syncCapRegistry brings the runtime offchain capabilities registry in line with committed
// state. It never applies the spec it was called with directly: this may run inside the
// (uncommitted) transaction creating the job. On boot the job is already committed, so the
// synchronous refresh applies it before other jobs start; inside a transaction it observes the
// previous committed state, and the hint makes the projector pick up the commit promptly.
func (d *delegate) syncCapRegistry(ctx context.Context) error {
	if d.capRegistry == nil {
		return fmt.Errorf("no offchain capabilities registry configured for config_type %q", ConfigTypeCapRegistry)
	}
	if err := d.capRegistry.Refresh(ctx); err != nil {
		// Not fatal for the job: the projector keeps retrying in the background.
		d.lggr.Warnw("Failed to refresh offchain capabilities registry config", "err", err)
	}
	d.capRegistry.Trigger()
	return nil
}

func (d *delegate) AfterJobCreated(j job.Job) {}

func (d *delegate) BeforeJobDeleted(j job.Job) {}

func (d *delegate) OnDeleteJob(ctx context.Context, jb job.Job) error {
	configType := d.configType(jb.CRESettingsSpec)
	if configType == ConfigTypeCapRegistry {
		// Runs inside the (uncommitted) delete transaction, so nothing is cleared here. Once the
		// delete commits, the projector clears the payload and capabilities fall back to
		// on-chain/TOML config; if a replacement commits in the same transaction, it moves
		// straight to the new payload; if the transaction rolls back, nothing changes.
		if d.capRegistry != nil {
			d.capRegistry.Trigger()
		}
		return nil
	}
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

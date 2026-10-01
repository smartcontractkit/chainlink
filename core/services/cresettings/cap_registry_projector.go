package cresettings

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

const (
	// defaultCapRegistryIdleInterval bounds how long a committed change (or a change missed for
	// any reason) can take to be applied when no hint was received.
	defaultCapRegistryIdleInterval = 10 * time.Second
	// After a hint (job created/deleted), poll quickly for a while: the job may be part of a
	// larger transaction (e.g. a feeds manager approval) that commits shortly after.
	defaultCapRegistryFastInterval = 200 * time.Millisecond
	defaultCapRegistryFastWindow   = 30 * time.Second
)

// CapRegistryProjector keeps the runtime offchain capabilities registry (GlobalConfig) equal to
// the committed capabilities_registry job in the database.
//
// The cresettings delegate is invoked inside the job create/delete transaction, before it
// commits (the feeds manager deletes the old job and creates its replacement in one
// transaction). Applying the payload from there would expose uncommitted state: a rolled-back
// create or delete would still change the runtime config, and a replacement would briefly
// clear it between the delete and the create. Instead the delegate only hints, and the
// projector reads committed state on a connection outside that transaction:
//   - a rolled-back transaction is never observed, so the runtime config is unchanged;
//   - a replacement commits atomically, so the runtime config moves from the old payload
//     straight to the new one;
//   - a committed delete clears the payload (capabilities fall back to on-chain/TOML).
//
// Validity, uniqueness and version monotonicity are enforced durably by the database when the
// spec is inserted (see job.ErrCRESettingsCapRegistryStale), so whatever is committed is
// applied.
type CapRegistryProjector struct {
	services.StateMachine
	lggr    logger.Logger
	load    func(ctx context.Context) (*job.CRESettingsSpec, error)
	gc      *globalconfig.GlobalConfig
	metrics *globalconfig.Metrics
	// lastDomain/lastEnv label the applied-version gauge; kept so it can be reset to 0 for the
	// same series when the payload is withdrawn. Guarded by mu.
	lastDomain, lastEnv string

	idleInterval, fastInterval, fastWindow time.Duration

	mu      sync.Mutex // serializes refreshes
	trigger chan struct{}
	stopCh  services.StopChan
	wg      sync.WaitGroup
}

// NewCapRegistryProjector returns a projector reading committed state from ds, which must not
// be a transaction.
func NewCapRegistryProjector(lggr logger.Logger, ds sqlutil.DataSource, gc *globalconfig.GlobalConfig) *CapRegistryProjector {
	return newCapRegistryProjector(lggr, func(ctx context.Context) (*job.CRESettingsSpec, error) {
		return job.LoadCapabilitiesRegistrySpec(ctx, ds)
	}, gc)
}

func newCapRegistryProjector(lggr logger.Logger, load func(context.Context) (*job.CRESettingsSpec, error), gc *globalconfig.GlobalConfig) *CapRegistryProjector {
	return &CapRegistryProjector{
		lggr:         logger.Named(lggr, "CapRegistryProjector"),
		load:         load,
		gc:           gc,
		metrics:      globalconfig.DefaultMetrics(),
		idleInterval: defaultCapRegistryIdleInterval,
		fastInterval: defaultCapRegistryFastInterval,
		fastWindow:   defaultCapRegistryFastWindow,
		trigger:      make(chan struct{}, 1),
		stopCh:       make(services.StopChan),
	}
}

// GlobalConfig returns the runtime registry this projector maintains.
func (p *CapRegistryProjector) GlobalConfig() *globalconfig.GlobalConfig { return p.gc }

func (p *CapRegistryProjector) Name() string { return p.lggr.Name() }

func (p *CapRegistryProjector) HealthReport() map[string]error {
	return map[string]error{p.Name(): p.Healthy()}
}

func (p *CapRegistryProjector) Start(context.Context) error {
	return p.StartOnce("CapRegistryProjector", func() error {
		p.wg.Go(p.run)
		return nil
	})
}

func (p *CapRegistryProjector) Close() error {
	return p.StopOnce("CapRegistryProjector", func() error {
		close(p.stopCh)
		p.wg.Wait()
		return nil
	})
}

// Trigger requests a prompt refresh, followed by fast polling for a while so a change that is
// committed shortly afterwards is picked up quickly. It never blocks.
func (p *CapRegistryProjector) Trigger() {
	select {
	case p.trigger <- struct{}{}:
	default:
	}
}

// Refresh applies the committed capabilities_registry payload (or clears the runtime config
// when there is none). Failures are counted in platform_cap_config_apply_errors_total.
func (p *CapRegistryProjector) Refresh(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	spec, err := p.load(ctx)
	if err != nil {
		p.metrics.RecordApplyError(ctx, p.lastDomain, p.lastEnv)
		return err
	}
	_, _, before := p.gc.Info()
	if spec == nil {
		p.gc.Clear()
		if before != 0 {
			p.lggr.Infow("Cleared offchain capabilities registry config (no committed capabilities_registry job)", "previousVersion", before)
		}
		p.recordApplied(ctx)
		return nil
	}
	if err = p.gc.Store(globalconfig.Update{Raw: spec.OffchainConfig, Hash: spec.Hash}); err != nil {
		domain, env := globalconfig.PayloadLabels(spec.OffchainConfig)
		p.metrics.RecordApplyError(ctx, domain, env)
		return fmt.Errorf("failed to apply committed capabilities_registry config: %w", err)
	}
	if _, _, after := p.gc.Info(); after != before {
		p.lggr.Infow("Applied offchain capabilities registry config", "version", after, "previousVersion", before, "hash", spec.Hash)
	}
	p.recordApplied(ctx)
	return nil
}

// recordApplied sets platform_cap_config_applied_version for the applied payload's
// (domain, env), and resets the previous series to 0 when the payload is withdrawn or its
// labels change. Callers must hold mu.
func (p *CapRegistryProjector) recordApplied(ctx context.Context) {
	domain, env, version := p.gc.Info()
	if version == 0 {
		domain, env = p.lastDomain, p.lastEnv
	} else if domain != p.lastDomain || env != p.lastEnv {
		p.metrics.RecordAppliedVersion(ctx, p.lastDomain, p.lastEnv, 0)
	}
	p.metrics.RecordAppliedVersion(ctx, domain, env, version)
	p.lastDomain, p.lastEnv = domain, env
}

func (p *CapRegistryProjector) run() {
	ctx, cancel := p.stopCh.NewCtx()
	defer cancel()

	var fastUntil time.Time
	refresh := func() {
		if err := p.Refresh(ctx); err != nil && ctx.Err() == nil {
			p.lggr.Errorw("Failed to refresh offchain capabilities registry config", "err", err)
		}
	}
	refresh()

	timer := time.NewTimer(p.idleInterval)
	defer timer.Stop()
	for {
		select {
		case <-p.stopCh:
			return
		case <-p.trigger:
			fastUntil = time.Now().Add(p.fastWindow)
			refresh()
		case <-timer.C:
			refresh()
		}
		next := p.idleInterval
		if time.Now().Before(fastUntil) {
			next = p.fastInterval
		}
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(next)
	}
}

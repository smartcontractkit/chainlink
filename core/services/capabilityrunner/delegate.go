// Package capabilityrunner manages job.CapabilityRunner jobs: capability
// binaries that the node launches via the empty LOOP (go-plugin liveness only,
// no RPC surface) and notifies of CRE settings updates over a local HTTP reload
// endpoint.
package capabilityrunner

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/capabilities/libs/standalone/capability"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

// ReloaderHealthName is the health check name of the per-job settings reloader.
const ReloaderHealthName = "SettingsReloader"

// SettingsPath is the file settings updates are dumped to. The path and the
// reload route are the contract between the node and the runner, so both come
// from the capabilities repo rather than being spelled out on either side.
func SettingsPath() string { return capability.SettingsPath() }

// Delegate manages job.CapabilityRunner jobs.
type Delegate struct {
	lggr            logger.Logger
	registrarConfig plugins.RegistrarConfig
	settings        core.SettingsBroadcaster
}

var _ job.Delegate = (*Delegate)(nil)

// NewDelegate takes the settings broadcaster the reloaders subscribe to. Pass a
// FileBackedSettings so that updates are dumped to SettingsPath() before any
// reloader is handed them.
func NewDelegate(lggr logger.Logger, registrarConfig plugins.RegistrarConfig, settings core.SettingsBroadcaster) *Delegate {
	return &Delegate{lggr: lggr, registrarConfig: registrarConfig, settings: settings}
}

func (d *Delegate) JobType() job.Type { return job.CapabilityRunner }

func (d *Delegate) BeforeJobCreated(job.Job) {}
func (d *Delegate) AfterJobCreated(job.Job)  {}
func (d *Delegate) BeforeJobDeleted(job.Job) {}

func (d *Delegate) OnDeleteJob(context.Context, job.Job) error { return nil }

func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) ([]job.ServiceCtx, error) {
	s := spec.CapabilityRunnerSpec
	if s == nil {
		return nil, errors.New("CapabilityRunnerSpec is required")
	}

	httpPort, err := HTTPPortFromArgs(s.Args)
	if err != nil {
		return nil, err
	}

	lggr := d.lggr.Named("CapabilityRunner").Named(spec.Name.ValueOrZero())

	// Resolve and register the binary the same way standard capability LOOPs
	// do: the spec's command, plus the shared capabilities env file.
	envVars, err := plugins.ParseEnvFile(env.CapabilitiesPlugin.Env.Get())
	if err != nil {
		return nil, fmt.Errorf("failed to parse capabilities env file: %w", err)
	}
	cmdFn, grpcOpts, err := d.registrarConfig.RegisterLOOP(plugins.CmdConfig{
		ID:   lggr.Name(),
		Cmd:  s.Command,
		Args: s.Args,
		Env:  envVars,
	})
	if err != nil {
		return nil, fmt.Errorf("error registering loop: %w", err)
	}

	// The empty LOOP launches and supervises the binary (liveness only).
	runner := loop.NewEmptyService(lggr, grpcOpts, cmdFn)
	reloader := newSettingsReloader(lggr, d.settings, runner, httpPort)

	return []job.ServiceCtx{runner, reloader}, nil
}

// settingsReloader keeps the runner's view of CRE settings current.
//
// Once the runner reports healthy it reloads once unconditionally, provided the
// settings file exists - the file outlives the process, so a runner restarted
// without any new settings update still needs to be pointed at what is already
// on disk. After that it reloads on every settings update.
//
// A reload is a GET to localhost:{http_port}/reload/<file>; 2xx means the runner
// reloaded successfully, anything else is a failure. The file itself is written
// by the fileBackedSettings wrapper before the update reaches this subscriber,
// so by the time one arrives here the on-disk state already matches it.
type settingsReloader struct {
	services.Service
	eng *services.Engine

	settings core.SettingsBroadcaster
	runner   services.HealthReporter
	path     string
	url      string
	client   *http.Client
}

func newSettingsReloader(lggr logger.Logger, settings core.SettingsBroadcaster, runner services.HealthReporter, httpPort int) *settingsReloader {
	r := &settingsReloader{
		settings: settings,
		runner:   runner,
		path:     SettingsPath(),
		url:      fmt.Sprintf("http://localhost:%d%s", httpPort, capability.ReloadPath()),
		client:   &http.Client{Timeout: 10 * time.Second},
	}
	r.Service, r.eng = services.Config{
		Name:  ReloaderHealthName,
		Start: r.start,
	}.NewServiceEngine(lggr)
	return r
}

func (r *settingsReloader) start(context.Context) error {
	r.eng.Go(r.run)
	return nil
}

func (r *settingsReloader) run(ctx context.Context) {
	r.initialReload(ctx)
	for ctx.Err() == nil {
		if !r.subscribeOnce(ctx) {
			return
		}
	}
}

// initialReload waits for the runner to become healthy, then reloads once if a
// settings file is already on disk.
func (r *settingsReloader) initialReload(ctx context.Context) {
	if !r.waitHealthy(ctx) {
		return
	}
	if _, err := os.Stat(r.path); err != nil {
		if !os.IsNotExist(err) {
			r.eng.Errorw("failed to stat settings file; skipping initial reload", "err", err, "path", r.path)
		}
		return
	}
	if err := r.reload(ctx); err != nil {
		r.eng.Errorw("initial reload failed", "err", err)
		return
	}
	r.eng.Infow("performed initial reload", "path", r.path)
}

// waitHealthy blocks until the runner reports healthy, returning false if ctx is
// done first.
func (r *settingsReloader) waitHealthy(ctx context.Context) bool {
	for {
		err := r.runner.Ready()
		if err == nil {
			return true
		}
		r.eng.Debugw("waiting for runner to become healthy", "err", err)
		select {
		case <-ctx.Done():
			return false
		case <-time.After(time.Second):
		}
	}
}

// subscribeOnce subscribes and consumes updates until the channel closes
// (returns true, to resubscribe) or ctx is done (returns false).
func (r *settingsReloader) subscribeOnce(ctx context.Context) bool {
	ch, err := r.settings.Subscribe(ctx)
	if err != nil {
		r.eng.Errorw("failed to subscribe to settings updates", "err", err)
		select {
		case <-ctx.Done():
			return false
		case <-time.After(5 * time.Second):
			return true
		}
	}
	if local, ok := r.settings.(core.SettingsBroadcasterLocal); ok {
		defer local.Unsubscribe(ch)
	}
	r.eng.Info("subscribed to settings updates")
	for {
		select {
		case <-ctx.Done():
			return false
		case update, ok := <-ch:
			if !ok {
				return true // resubscribe
			}
			if err := r.reload(ctx); err != nil {
				r.eng.Errorw("failed to reload settings", "err", err, "hash", update.Hash)
				continue
			}
			r.eng.Infow("reloaded settings", "hash", update.Hash)
		}
	}
}

func (r *settingsReloader) reload(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.url, nil)
	if err != nil {
		return fmt.Errorf("failed to build reload request: %w", err)
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("reload request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("reload endpoint returned status %d", resp.StatusCode)
	}
	return nil
}

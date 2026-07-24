// Package capabilityrunner manages job.CapabilityRunner jobs: capability
// binaries that the node launches via the empty LOOP (go-plugin liveness only,
// no RPC surface) and notifies of CRE settings (limits) updates over a local
// HTTP reload endpoint.
package capabilityrunner

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

const (
	// limitsDirName names the directory (under os.TempDir()) where settings
	// updates are dumped for runner binaries. The capabilityrunner dependency in
	// the capabilities repo reads the same conventional location when its
	// /reload endpoint is hit; both processes share the container so
	// os.TempDir() resolves identically.
	limitsDirName = "cre_limits"
	// LimitsFileName is the dumped settings file name and the path suffix of
	// the runner's reload endpoint: /reload/<LimitsFileName>.
	LimitsFileName = "cre_limits.txt"
)

// LimitsDir is the directory settings updates are dumped to.
func LimitsDir() string { return filepath.Join(os.TempDir(), limitsDirName) }

// Delegate manages job.CapabilityRunner jobs.
type Delegate struct {
	lggr            logger.Logger
	registrarConfig plugins.RegistrarConfig
	settings        core.SettingsBroadcaster
}

var _ job.Delegate = (*Delegate)(nil)

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
	reloader := newLimitsReloader(lggr, d.settings, httpPort)

	return []job.ServiceCtx{runner, reloader}, nil
}

// limitsReloader subscribes to CRE settings updates; on each update it dumps
// the settings to LimitsDir()/LimitsFileName and hits the runner's
// localhost:{http_port}/reload/<LimitsFileName> endpoint. A 2xx response means
// the runner reloaded successfully; anything else is a failure.
type limitsReloader struct {
	services.Service
	eng *services.Engine

	settings core.SettingsBroadcaster
	url      string
	client   *http.Client
}

func newLimitsReloader(lggr logger.Logger, settings core.SettingsBroadcaster, httpPort int) *limitsReloader {
	r := &limitsReloader{
		settings: settings,
		url:      fmt.Sprintf("http://localhost:%d/reload/%s", httpPort, LimitsFileName),
		client:   &http.Client{Timeout: 10 * time.Second},
	}
	r.Service, r.eng = services.Config{
		Name:  "LimitsReloader",
		Start: r.start,
	}.NewServiceEngine(lggr)
	return r
}

func (r *limitsReloader) start(context.Context) error {
	r.eng.Go(r.run)
	return nil
}

func (r *limitsReloader) run(ctx context.Context) {
	for ctx.Err() == nil {
		if !r.subscribeOnce(ctx) {
			return
		}
	}
}

// subscribeOnce subscribes and consumes updates until the channel closes
// (returns true, to resubscribe) or ctx is done (returns false).
func (r *limitsReloader) subscribeOnce(ctx context.Context) bool {
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
			if err := r.apply(ctx, update); err != nil {
				r.eng.Errorw("failed to reload limits", "err", err, "hash", update.Hash)
				continue
			}
			r.eng.Infow("reloaded limits", "hash", update.Hash)
		}
	}
}

func (r *limitsReloader) apply(ctx context.Context, update core.SettingsUpdate) error {
	dir := LimitsDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create limits dir: %w", err)
	}
	// Write to a temp file and rename so the runner never reads a torn file.
	tmp, err := os.CreateTemp(dir, LimitsFileName+".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp limits file: %w", err)
	}
	if _, err = tmp.WriteString(update.Settings); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
		return fmt.Errorf("failed to write limits file: %w", err)
	}
	if err = tmp.Close(); err != nil {
		_ = os.Remove(tmp.Name())
		return fmt.Errorf("failed to close limits file: %w", err)
	}
	if err = os.Rename(tmp.Name(), filepath.Join(dir, LimitsFileName)); err != nil {
		_ = os.Remove(tmp.Name())
		return fmt.Errorf("failed to move limits file into place: %w", err)
	}

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

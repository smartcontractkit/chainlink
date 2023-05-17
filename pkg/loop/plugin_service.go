package loop

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/hashicorp/go-plugin"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
)

const keepAliveTickDuration = 5 * time.Second //TODO from config

type grpcPlugin interface {
	plugin.Plugin
	plugin.GRPCPlugin
	ClientConfig() *plugin.ClientConfig
}

// pluginService is a [types.Service] wrapper that maintains an internal [types.Service] created from a [grpcPlugin]
// client instance by launching and re-launching as necessary.
type pluginService[P grpcPlugin, S types.Service] struct {
	utils.StartStopOnce

	pluginName string

	lggr logger.Logger
	cmd  func() *exec.Cmd

	wg     sync.WaitGroup
	stopCh chan struct{}

	grpcPlug P

	client         *plugin.Client
	clientProtocol plugin.ClientProtocol

	newService func(context.Context, any) (S, error)

	serviceCh chan struct{} // closed when service is available
	service   S

	testInterrupt chan func(*pluginService[P, S]) // tests only (via TestHook) to enable access to internals without racing
}

func newPluginService[P grpcPlugin, S types.Service](pluginName string, p P, newService func(context.Context, any) (S, error), lggr logger.Logger, cmd func() *exec.Cmd, stopCh chan struct{}) *pluginService[P, S] {
	return &pluginService[P, S]{
		pluginName: pluginName,
		lggr:       lggr,
		cmd:        cmd,
		stopCh:     stopCh,
		grpcPlug:   p,
		newService: newService,
		serviceCh:  make(chan struct{}),
	}
}

func (r *pluginService[P, S]) keepAlive() {
	defer r.wg.Done()

	r.lggr.Debugw("Staring keepAlive", "tick", keepAliveTickDuration)

	t := time.NewTicker(keepAliveTickDuration)
	defer t.Stop()
	for {
		select {
		case <-r.stopCh:
			return
		case <-t.C:
			c := r.client
			cp := r.clientProtocol
			if c != nil && !c.Exited() && cp != nil {
				// launched
				err := cp.Ping()
				if err == nil {
					continue // healthy
				}
				r.lggr.Errorw("Relaunching unhealthy plugin", "err", err)
			}
			if err := r.tryLaunch(cp); err != nil {
				r.lggr.Errorw("Failed to launch plugin", "err", err)
			}
		case fn := <-r.testInterrupt:
			fn(r)
		}
	}
}

func (r *pluginService[P, S]) tryLaunch(old plugin.ClientProtocol) (err error) {
	if old != nil && r.clientProtocol != old {
		// already replaced by another routine
		return nil
	}
	r.client, r.clientProtocol, err = r.launch()
	return
}

func (r *pluginService[P, S]) launch() (*plugin.Client, plugin.ClientProtocol, error) {
	ctx, cancelFn := utils.ContextFromChan(r.stopCh)
	defer cancelFn()

	r.lggr.Debug("Launching")

	cc := r.grpcPlug.ClientConfig()
	cc.Cmd = r.cmd()
	client := plugin.NewClient(cc)
	cp, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, nil, fmt.Errorf("failed to create ClientProtocol: %w", err)
	}
	abort := func() {
		if cerr := cp.Close(); cerr != nil {
			r.lggr.Errorw("Error closing ClientProtocol", "err", cerr)
		}
		client.Kill()
	}
	i, err := cp.Dispense(r.pluginName)
	if err != nil {
		abort()
		return nil, nil, fmt.Errorf("failed to Dispense %q plugin: %w", r.pluginName, err)
	}

	select {
	case <-r.serviceCh:
		// r.service already set
	default:
		r.service, err = r.newService(ctx, i)
		if err != nil {
			abort()
			return nil, nil, fmt.Errorf("failed to create service: %w", err)
		}
		defer close(r.serviceCh)
	}
	return client, cp, nil
}

func (r *pluginService[P, S]) Start(context.Context) error {
	return r.StartOnce("PluginService", func() error {
		r.wg.Add(1)
		go r.keepAlive()
		return nil
	})
}

func (r *pluginService[P, S]) Ready() error {
	select {
	case <-r.serviceCh:
		return r.service.Ready()
	default:
		return ErrPluginUnavailable
	}
}

func (r *pluginService[P, S]) Name() string { return r.lggr.Name() }

func (r *pluginService[P, S]) HealthReport() map[string]error {
	select {
	case <-r.serviceCh:
		hr := map[string]error{r.Name(): r.Healthy()}
		maps.Copy(hr, r.service.HealthReport())
		return hr
	default:
		return map[string]error{r.Name(): ErrPluginUnavailable}
	}
}

func (r *pluginService[P, S]) Close() error {
	return r.StopOnce("PluginService", func() (err error) {
		close(r.stopCh)
		r.wg.Wait()

		select {
		case <-r.serviceCh:
			if cerr := r.service.Close(); !errors.Is(cerr, context.Canceled) && status.Code(cerr) != codes.Canceled {
				err = errors.Join(err, cerr)
			}
		default:
		}
		if r.clientProtocol != nil {
			if cerr := r.clientProtocol.Close(); !errors.Is(cerr, context.Canceled) {
				err = errors.Join(err, cerr)
			}
		}
		if r.client != nil {
			r.client.Kill()
		}
		return
	})
}

func (r *pluginService[P, S]) wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case <-r.serviceCh:
		return nil
	}
}

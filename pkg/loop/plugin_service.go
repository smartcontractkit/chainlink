package loop

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal"
	"github.com/smartcontractkit/chainlink-relay/pkg/services"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
)

const keepAliveTickDuration = 5 * time.Second //TODO from config

type BrokerConfig = internal.BrokerConfig

type grpcPlugin interface {
	plugin.Plugin
	plugin.GRPCPlugin
	ClientConfig() *plugin.ClientConfig
}

// pluginService is a [types.Service] wrapper that maintains an internal [types.Service] created from a [grpcPlugin]
// client instance by launching and re-launching as necessary.
type pluginService[P grpcPlugin, S services.Service] struct {
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

func (s *pluginService[P, S]) init(pluginName string, p P, newService func(context.Context, any) (S, error), lggr logger.Logger, cmd func() *exec.Cmd, stopCh chan struct{}) {
	s.pluginName = pluginName
	s.lggr = lggr
	s.cmd = cmd
	s.stopCh = stopCh
	s.grpcPlug = p
	s.newService = newService
	s.serviceCh = make(chan struct{})
}

func (s *pluginService[P, S]) keepAlive() {
	defer s.wg.Done()

	s.lggr.Debugw("Starting keepAlive", "tick", keepAliveTickDuration)

	check := func() {
		c := s.client
		cp := s.clientProtocol
		if c != nil && !c.Exited() && cp != nil {
			// launched
			err := cp.Ping()
			if err == nil {
				return // healthy
			}
			s.lggr.Errorw("Relaunching unhealthy plugin", "err", err)
		}
		if err := s.tryLaunch(cp); err != nil {
			s.lggr.Errorw("Failed to launch plugin", "err", err)
		}
	}

	check() // no delay

	t := time.NewTicker(keepAliveTickDuration)
	defer t.Stop()
	for {
		select {
		case <-s.stopCh:
			return
		case <-t.C:
			check()
		case fn := <-s.testInterrupt:
			fn(s)
		}
	}
}

func (s *pluginService[P, S]) tryLaunch(old plugin.ClientProtocol) (err error) {
	if old != nil && s.clientProtocol != old {
		// already replaced by another routine
		return nil
	}
	if cerr := s.closeClient(); cerr != nil {
		s.lggr.Errorw("Error closing old client", "err", cerr)
	}
	s.client, s.clientProtocol, err = s.launch()
	return
}

func (s *pluginService[P, S]) launch() (*plugin.Client, plugin.ClientProtocol, error) {
	ctx, cancelFn := utils.ContextFromChan(s.stopCh)
	defer cancelFn()

	s.lggr.Debug("Launching")

	cc := s.grpcPlug.ClientConfig()
	cc.Cmd = s.cmd()
	client := plugin.NewClient(cc)
	cp, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, nil, fmt.Errorf("failed to create ClientProtocol: %w", err)
	}
	abort := func() {
		if cerr := cp.Close(); cerr != nil {
			s.lggr.Errorw("Error closing ClientProtocol", "err", cerr)
		}
		client.Kill()
	}
	i, err := cp.Dispense(s.pluginName)
	if err != nil {
		abort()
		return nil, nil, fmt.Errorf("failed to Dispense %q plugin: %w", s.pluginName, err)
	}

	select {
	case <-s.serviceCh:
		// s.service already set
	default:
		s.service, err = s.newService(ctx, i)
		if err != nil {
			abort()
			return nil, nil, fmt.Errorf("failed to create service: %w", err)
		}
		defer close(s.serviceCh)
	}
	return client, cp, nil
}

func (s *pluginService[P, S]) Start(context.Context) error {
	return s.StartOnce("PluginService", func() error {
		s.wg.Add(1)
		go s.keepAlive()
		return nil
	})
}

func (s *pluginService[P, S]) Ready() error {
	select {
	case <-s.serviceCh:
		return s.service.Ready()
	default:
		return ErrPluginUnavailable
	}
}

func (s *pluginService[P, S]) Name() string { return s.lggr.Name() }

func (s *pluginService[P, S]) HealthReport() map[string]error {
	select {
	case <-s.serviceCh:
		hr := map[string]error{s.Name(): s.Healthy()}
		services.CopyHealth(hr, s.service.HealthReport())
		return hr
	default:
		return map[string]error{s.Name(): ErrPluginUnavailable}
	}
}

func (s *pluginService[P, S]) Close() error {
	return s.StopOnce("PluginService", func() (err error) {
		close(s.stopCh)
		s.wg.Wait()

		select {
		case <-s.serviceCh:
			if cerr := s.service.Close(); !errors.Is(cerr, context.Canceled) && status.Code(cerr) != codes.Canceled {
				err = errors.Join(err, cerr)
			}
		default:
		}
		err = errors.Join(err, s.closeClient())
		return
	})
}

func (s *pluginService[P, S]) closeClient() (err error) {
	if s.clientProtocol != nil {
		if cerr := s.clientProtocol.Close(); !errors.Is(cerr, context.Canceled) {
			err = cerr
		}
	}
	if s.client != nil {
		s.client.Kill()
	}
	return
}

func (s *pluginService[P, S]) wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case <-s.serviceCh:
		return nil
	}
}

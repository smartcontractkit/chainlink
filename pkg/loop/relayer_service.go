package loop

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/hashicorp/go-plugin"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
)

var ErrPluginUnavailable = errors.New("plugin unavailable")

const keepAliveTickDuration = 5 * time.Second //TODO from config

var _ types.Service = (*RelayerService)(nil)

// RelayerService is a [types.Service] that maintains an internal [Relayer] from a [PluginRelayer] client
// instance by launching and re-launching as necessary.
type RelayerService struct {
	utils.StartStopOnce

	lggr     logger.Logger
	cmd      func() *exec.Cmd
	config   string
	keystore Keystore

	wg     sync.WaitGroup
	stopCh chan struct{}

	mu             sync.RWMutex
	client         *plugin.Client
	clientProtocol plugin.ClientProtocol
	plug           PluginRelayer
	relayer        Relayer
}

// NewRelayerService returns a new [*RelayerService].
// cmd must return a new exec.Cmd each time it is called.
func NewRelayerService(lggr logger.Logger, cmd func() *exec.Cmd, config string, keystore Keystore) *RelayerService {
	return &RelayerService{lggr: lggr, cmd: cmd, config: config, keystore: keystore, stopCh: make(chan struct{})}
}

func (p *RelayerService) launch() (*plugin.Client, plugin.ClientProtocol, PluginRelayer, Relayer, error) {
	ctx, cancelFn := utils.ContextFromChan(p.stopCh)
	defer cancelFn()
	cc := PluginRelayerClientConfig(p.lggr)
	cc.Cmd = p.cmd()
	client := plugin.NewClient(cc)
	cp, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, nil, nil, nil, fmt.Errorf("failed to create plugin Client: %w", err)
	}
	abort := func() {
		if cerr := cp.Close(); cerr != nil {
			p.lggr.Errorw("Error closing ClientProtocol", "err", cerr)
		}
		client.Kill()
	}
	i, err := cp.Dispense(PluginRelayerName)
	if err != nil {
		abort()
		return nil, nil, nil, nil, fmt.Errorf("failed to Dispense %q plugin: %w", PluginRelayerName, err)
	}
	plug, ok := i.(PluginRelayer)
	if !ok {
		abort()
		return nil, nil, nil, nil, fmt.Errorf("expected PluginRelayer but got %T", i)
	}
	relayer, err := plug.NewRelayer(ctx, p.config, p.keystore)
	if err != nil {
		abort()
		return nil, nil, nil, nil, fmt.Errorf("failed to create Relayer: %w", err)
	}
	err = relayer.Start(ctx)
	if err != nil {
		abort()
		return nil, nil, nil, nil, fmt.Errorf("failed to start Relayer: %w", err)
	}
	return client, cp, plug, relayer, nil
}

func (p *RelayerService) tryLaunch(old plugin.ClientProtocol) (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if old != nil && p.clientProtocol != old {
		// already replaced by another routine
		return nil
	}
	p.client, p.clientProtocol, p.plug, p.relayer, err = p.launch()
	return
}

func (p *RelayerService) Start(context.Context) error {
	return p.StartOnce("RelayerService", func() error {
		if err := p.tryLaunch(nil); err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				return fmt.Errorf("missing plugin executable: %w", err)
			}
			p.lggr.Error("Failed to launch plugin", "err", err)
		}

		p.wg.Add(1)
		go p.keepAlive()
		return nil
	})
}

func (p *RelayerService) Name() string { return p.lggr.Name() }

func (p *RelayerService) HealthReport() map[string]error {
	hr := map[string]error{
		p.Name(): p.Healthy(),
	}
	p.mu.RLock()
	relayer := p.relayer
	p.mu.RUnlock()
	if relayer != nil {
		for n, e := range relayer.HealthReport() {
			hr[n] = e
		}
	}
	return hr
}

func (p *RelayerService) ping() error {
	p.mu.RLock()
	cp := p.clientProtocol
	p.mu.RUnlock()
	if cp == nil {
		return ErrPluginUnavailable
	}
	return cp.Ping()
}

func (p *RelayerService) Ready() error { return p.ping() }

func (p *RelayerService) keepAlive() {
	defer p.wg.Done()

	t := time.NewTicker(keepAliveTickDuration)
	defer t.Stop()
	for {
		select {
		case <-p.stopCh:
			return
		case <-t.C:
			p.mu.RLock()
			c := p.client
			cp := p.clientProtocol
			p.mu.RUnlock()
			if c != nil && !c.Exited() && cp != nil {
				// launched
				err := cp.Ping()
				if err == nil {
					continue // healthy
				}
				p.lggr.Errorw("Relaunching unhealthy plugin", "err", err)
			}
			if err := p.tryLaunch(cp); err != nil {
				p.lggr.Errorw("Failed to launch plugin", "err", err)
			}
		}
	}
}

func (p *RelayerService) Close() error {
	return p.StopOnce("RelayerService", func() (err error) {
		close(p.stopCh)
		p.wg.Wait()

		p.mu.RLock()
		defer p.mu.RUnlock()
		if p.relayer != nil {
			err = errors.Join(err, p.relayer.Close())
		}
		if p.clientProtocol != nil {
			err = errors.Join(err, p.clientProtocol.Close())
		}
		if p.client != nil {
			p.client.Kill()
		}
		return
	})
}

func (p *RelayerService) Relayer() (Relayer, error) {
	p.mu.RLock()
	relayer := p.relayer
	p.mu.RUnlock()
	if relayer == nil {
		return nil, ErrPluginUnavailable
	}
	return relayer, nil
}

package internal

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
)

// Broker is a subset of the methods exported by *plugin.GRPCBroker.
type Broker interface {
	Accept(id uint32) (net.Listener, error)
	DialWithOptions(id uint32, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error)
	NextId() uint32
}

var _ Broker = (*atomicBroker)(nil)

// An atomicBroker implements [Broker] and is backed by a swappable [*plugin.GRPCBroker]
type atomicBroker struct {
	broker atomic.Pointer[Broker]
}

func (a *atomicBroker) store(b Broker) { a.broker.Store(&b) }
func (a *atomicBroker) load() Broker   { return *a.broker.Load() }

func (a *atomicBroker) Accept(id uint32) (net.Listener, error) {
	return a.load().Accept(id)
}

func (a *atomicBroker) DialWithOptions(id uint32, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	return a.load().DialWithOptions(id, opts...)
}

func (a *atomicBroker) NextId() uint32 { //nolint:revive
	return a.load().NextId()
}

// GRPCOpts has GRPC client and server options.
type GRPCOpts struct {
	// Optionally include additional options when dialing a client.
	// Normally aligned with [plugin.ClientConfig.GRPCDialOptions].
	DialOpts []grpc.DialOption
	// Optionally override the default *grpc.Server constructor.
	// Normally aligned with [plugin.ServeConfig.GRPCServer].
	NewServer func([]grpc.ServerOption) *grpc.Server
}

// BrokerConfig holds Broker configuration fields.
type BrokerConfig struct {
	StopCh <-chan struct{}
	Logger logger.Logger

	GRPCOpts // optional
}

// brokerExt extends a Broker with various helper methods.
type brokerExt struct {
	broker Broker
	BrokerConfig
}

// withName returns a new [*brokerExt] with name added to the logger.
func (b *brokerExt) withName(name string) *brokerExt {
	bn := *b
	bn.Logger = logger.Named(b.Logger, name)
	return &bn
}

// newClientConn return a new *clientConn backed by this *brokerExt.
func (b *brokerExt) newClientConn(name string, newClient newClientFn) *clientConn {
	return &clientConn{
		brokerExt: b.withName(name),
		newClient: newClient,
		name:      name,
	}
}

func (b *brokerExt) stopCtx() (context.Context, context.CancelFunc) {
	return utils.ContextFromChan(b.StopCh)
}

func (b *brokerExt) dial(id uint32) (conn *grpc.ClientConn, err error) {
	return b.broker.DialWithOptions(id, b.DialOpts...)
}

func (b *brokerExt) serveNew(name string, register func(*grpc.Server), deps ...resource) (uint32, resource, error) {
	var server *grpc.Server
	if b.NewServer == nil {
		server = grpc.NewServer()
	} else {
		server = b.NewServer(nil)
	}
	register(server)
	return b.serve(name, server, deps...)
}

func (b *brokerExt) serve(name string, server *grpc.Server, deps ...resource) (uint32, resource, error) {
	id := b.broker.NextId()
	b.Logger.Debugf("Serving %s on connection %d", name, id)
	lis, err := b.broker.Accept(id)
	if err != nil {
		b.closeAll(deps...)
		return 0, resource{}, ErrConnAccept{Name: name, ID: id, Err: err}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer b.closeAll(deps...)
		if err := server.Serve(lis); err != nil {
			b.Logger.Errorw(fmt.Sprintf("Failed to serve %s on connection %d", name, id), "err", err)
		}
	}()

	done := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-b.StopCh:
			server.Stop()
		case <-done:
		}
	}()

	return id, resource{fnCloser(func() {
		server.Stop()
		close(done)
		wg.Wait()
	}), name}, nil
}

func (b *brokerExt) closeAll(deps ...resource) {
	for _, d := range deps {
		if err := d.Close(); err != nil {
			b.Logger.Error(fmt.Sprintf("Error closing %s", d.name), "err", err)
		}
	}
}

type resource struct {
	io.Closer
	name string
}

type resources []resource

func (rs *resources) Add(r resource) {
	*rs = append(*rs, r)
}

func (rs *resources) Stop(s interface{ Stop() }, name string) {
	rs.Add(resource{fnCloser(s.Stop), name})
}

func (rs *resources) Close(c io.Closer, name string) {
	rs.Add(resource{c, name})
}

// fnCloser implements io.Closer with a func().
type fnCloser func()

func (s fnCloser) Close() error {
	s()
	return nil
}

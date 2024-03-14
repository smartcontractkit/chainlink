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

// An atomicBroker implements [Broker] and is backed by a swappable [*plugin.GRPCBroker].
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

// BrokerExt extends a Broker with various helper methods.
type BrokerExt struct {
	Broker Broker
	BrokerConfig
}

// WithName returns a new [*BrokerExt] with Name added to the logger.
func (b *BrokerExt) WithName(name string) *BrokerExt {
	bn := *b
	bn.Logger = logger.Named(b.Logger, name)
	return &bn
}

// NewClientConn return a new *clientConn backed by this *BrokerExt.
func (b *BrokerExt) NewClientConn(name string, newClient newClientFn) *clientConn {
	return &clientConn{
		BrokerExt: b.WithName(name),
		newClient: newClient,
		name:      name,
	}
}

func (b *BrokerExt) StopCtx() (context.Context, context.CancelFunc) {
	return utils.ContextFromChan(b.StopCh)
}

func (b *BrokerExt) Dial(id uint32) (conn *grpc.ClientConn, err error) {
	return b.Broker.DialWithOptions(id, b.DialOpts...)
}

func (b *BrokerExt) ServeNew(name string, register func(*grpc.Server), deps ...Resource) (uint32, Resource, error) {
	var server *grpc.Server
	if b.NewServer == nil {
		server = grpc.NewServer()
	} else {
		server = b.NewServer(nil)
	}
	register(server)
	return b.Serve(name, server, deps...)
}

func (b *BrokerExt) Serve(name string, server *grpc.Server, deps ...Resource) (uint32, Resource, error) {
	id := b.Broker.NextId()
	b.Logger.Debugf("Serving %s on connection %d", name, id)
	lis, err := b.Broker.Accept(id)
	if err != nil {
		b.CloseAll(deps...)
		return 0, Resource{}, ErrConnAccept{Name: name, ID: id, Err: err}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer b.CloseAll(deps...)
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

	return id, Resource{fnCloser(func() {
		server.Stop()
		close(done)
		wg.Wait()
	}), name}, nil
}

func (b *BrokerExt) CloseAll(deps ...Resource) {
	for _, d := range deps {
		if err := d.Close(); err != nil {
			b.Logger.Error(fmt.Sprintf("Error closing %s", d.Name), "err", err)
		}
	}
}

type Resource struct {
	io.Closer
	Name string
}

type Resources []Resource

func (rs *Resources) Add(r Resource) {
	*rs = append(*rs, r)
}

func (rs *Resources) Stop(s interface{ Stop() }, name string) {
	rs.Add(Resource{fnCloser(s.Stop), name})
}

func (rs *Resources) Close(c io.Closer, name string) {
	rs.Add(Resource{c, name})
}

// fnCloser implements io.Closer with a func().
type fnCloser func()

func (s fnCloser) Close() error {
	s()
	return nil
}

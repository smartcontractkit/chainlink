package loop

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
)

// broker is a subset of the methods exported by *plugin.GRPCBroker.
type broker interface {
	Accept(id uint32) (net.Listener, error)
	Dial(id uint32) (conn *grpc.ClientConn, err error)
	NextId() uint32
}

var _ broker = (*atomicBroker)(nil)

// An atomicBroker implements [broker] and is backed by a swappable [*plugin.GRPCBroker]
type atomicBroker struct {
	broker atomic.Pointer[plugin.GRPCBroker]
}

func (a *atomicBroker) store(b *plugin.GRPCBroker) { a.broker.Store(b) }

func (a *atomicBroker) Accept(id uint32) (net.Listener, error) {
	return a.broker.Load().Accept(id)
}

func (a *atomicBroker) Dial(id uint32) (conn *grpc.ClientConn, err error) {
	return a.broker.Load().Dial(id)
}

func (a *atomicBroker) NextId() uint32 {
	return a.broker.Load().NextId()
}

// brokerExt extends a broker with various helper methods.
type brokerExt struct {
	stopCh <-chan struct{}
	lggr   logger.Logger
	broker broker
}

// named returns a new [*brokerExt] with name added to the logger.
func (b *brokerExt) named(name string) *brokerExt {
	return &brokerExt{
		stopCh: b.stopCh,
		lggr:   logger.Named(b.lggr, name),
		broker: b.broker,
	}
}

// newClientConn return a new *clientConn backed by this *brokerExt.
func (b *brokerExt) newClientConn(name string, newClient newClientFn) *clientConn {
	return &clientConn{
		brokerExt: b.named(name),
		newClient: newClient,
		name:      name,
	}
}

func (b *brokerExt) ctx() (context.Context, context.CancelFunc) {
	return utils.ContextFromChan(b.stopCh)
}

func (b *brokerExt) serve(name string, register func(*grpc.Server), deps ...resource) (uint32, resource, error) {
	id := b.broker.NextId()
	b.lggr.Debugf("Serving %s on connection %d", name, id)
	lis, err := b.broker.Accept(id)
	if err != nil {
		b.closeAll(deps...)
		return 0, resource{}, ErrConnAccept{Name: name, ID: id, Err: err}
	}

	server := grpc.NewServer()
	register(server)
	go func() {
		defer b.closeAll(deps...)
		if err := server.Serve(lis); err != nil {
			b.lggr.Errorw(fmt.Sprintf("Failed to serve %s on connection %d", name, id), "err", err)
		}
	}()

	done := make(chan struct{})
	go func() {
		select {
		case <-b.stopCh:
			server.Stop()
		case <-done:
		}
	}()

	return id, resource{fnCloser(func() {
		server.Stop()
		close(done)
	}), name}, nil
}

func (b *brokerExt) closeAll(deps ...resource) {
	for _, d := range deps {
		if err := d.Close(); err != nil {
			b.lggr.Error(fmt.Sprintf("Error closing %s", d.name), "err", err)
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

package internal

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jpillora/backoff"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var _ grpc.ClientConnInterface = (*atomicClient)(nil)

// An atomicClient implements [grpc.ClientConnInterface] and is backed by a swappable [*grpc.ClientConn]
type atomicClient struct {
	cc atomic.Pointer[grpc.ClientConn]
}

func (a *atomicClient) store(cc *grpc.ClientConn) { a.cc.Store(cc) }

func (a *atomicClient) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	return a.cc.Load().Invoke(ctx, method, args, reply, opts...)
}

func (a *atomicClient) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return a.cc.Load().NewStream(ctx, desc, method, opts...)
}

var _ grpc.ClientConnInterface = (*clientConn)(nil)

// newClientFn returns a new client connection id to dial, and a set of resource dependencies to close.
type newClientFn func(context.Context) (id uint32, deps resources, err error)

// clientConn is a [grpc.ClientConnInterface] backed by a [*grpc.ClientConn] which can be recreated and swapped out
// via the provided [newClientFn].
// New instances should be created via brokerExt.newClientConn.
type clientConn struct {
	*brokerExt
	newClient newClientFn
	name      string

	mu   sync.RWMutex
	deps resources
	cc   *grpc.ClientConn
}

func (c *clientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	c.mu.RLock()
	cc := c.cc
	c.mu.RUnlock()

	if cc == nil {
		cc = c.refresh(ctx, nil)
	}
	for cc != nil {
		err := cc.Invoke(ctx, method, args, reply, opts...)
		if isErrTerminal(err) {
			c.Logger.Warnw("clientConn: Invoke: terminal error, refreshing connection", "err", err)
			cc = c.refresh(ctx, cc)
			continue
		}
		return err
	}
	return context.Cause(ctx)
}

func (c *clientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	c.mu.RLock()
	cc := c.cc
	c.mu.RUnlock()

	if cc == nil {
		cc = c.refresh(ctx, nil)
	}
	for cc != nil {
		s, err := cc.NewStream(ctx, desc, method, opts...)
		if isErrTerminal(err) {
			c.Logger.Warnw("clientConn: NewStream: terminal error, refreshing connection", "err", err)
			cc = c.refresh(ctx, cc)
			continue
		}
		return s, err
	}
	return nil, context.Cause(ctx)
}

// refresh replaces c.cc with a new (different from orig) *grpc.ClientConn, and returns it as well.
// It will block until a new connection is successfully dialed, or return nil if the context expires.
func (c *clientConn) refresh(ctx context.Context, orig *grpc.ClientConn) *grpc.ClientConn {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cc != orig {
		return c.cc
	}
	if c.cc != nil {
		if err := c.cc.Close(); err != nil {
			c.Logger.Errorw("Client close failed", "err", err)
		}
		c.closeAll(c.deps...)
	}

	try := func() bool {
		c.Logger.Debug("Client refresh")
		id, deps, err := c.newClient(ctx)
		if err != nil {
			c.Logger.Errorw("Client refresh attempt failed", "err", err)
			c.closeAll(deps...)
			return false
		}
		c.deps = deps

		lggr := logger.With(c.Logger, "id", id)
		lggr.Debug("Client dial")
		c.cc, err = c.dial(id)
		if err != nil {
			if ctx.Err() != nil {
				lggr.Errorw("Client dial failed", "err", ErrConnDial{Name: c.name, ID: id, Err: err})
			}
			c.closeAll(c.deps...)
			return false
		}
		return true
	}

	b := backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    5 * time.Second,
		Factor: 2,
	}
	for !try() {
		if ctx.Err() != nil {
			c.Logger.Errorw("Client refresh failed: aborting refresh due to context error", "err", ctx.Err())
			return nil
		}
		wait := b.Duration()
		c.Logger.Infow("Waiting to refresh", "wait", wait)
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(wait):
		}
	}

	return c.cc
}

// isErrTerminal returns true if the grpc [status] [codes.Code] indicates that the plugin connection has terminated and
// must be refreshed.
func isErrTerminal(err error) bool {
	switch status.Code(err) {
	case codes.Unavailable, codes.Canceled:
		return true
	case codes.OK, codes.Unknown, codes.InvalidArgument, codes.DeadlineExceeded, codes.NotFound, codes.AlreadyExists,
		codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted, codes.OutOfRange,
		codes.Unimplemented, codes.Internal, codes.DataLoss, codes.Unauthenticated:
		return false
	}
	return false
}

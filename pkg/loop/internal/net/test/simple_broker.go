package test

import (
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/hashicorp/consul/sdk/freeport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	loopnet "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
)

// Broker is a test implementation of loopnet.Broker.
type Broker struct {
	T  *testing.T
	mu sync.Mutex

	// The next ID to be assigned.
	nextID uint32

	// The listeners that have been created.
	// use lazy initialization
	listeners map[uint32]net.Listener

	once sync.Once
}

var _ loopnet.Broker = (*Broker)(nil)

func (v *Broker) init() {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.listeners == nil {
		v.listeners = make(map[uint32]net.Listener)
	}
	v.T.Cleanup(func() { v.close() })
}

// Accept implements net.Broker.
func (v *Broker) Accept(id uint32) (net.Listener, error) {
	v.once.Do(v.init)

	v.mu.Lock()
	defer v.mu.Unlock()

	if l, exists := v.listeners[id]; exists {
		return l, nil
	}

	port := freeport.GetOne(v.T)
	l, err := net.Listen("tcp", "localhost:"+fmt.Sprint(port))
	if err != nil {
		return nil, err
	}

	v.listeners[id] = l
	return l, nil
}

// DialWithOptions implements net.Broker.
func (v *Broker) DialWithOptions(id uint32, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	v.once.Do(v.init)
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.listeners == nil {
		return nil, fmt.Errorf("listener with id %d does not exist", id)
	}

	l, exists := v.listeners[id]
	if !exists {
		return nil, fmt.Errorf("listener with id %d does not exist", id)
	}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return grpc.Dial(l.Addr().String(), opts...)
}

// NextId implements net.Broker.
// nolint:revive
func (v *Broker) NextId() uint32 {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.nextID++
	return v.nextID
}

func (v *Broker) close() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	for _, l := range v.listeners {
		l.Close()
	}
	return nil
}

// Package creregistry wires the core node to the plain-gRPC
// CapabilitiesRegistry served by the p2p proxy process (crecore).
//
// The client itself lives in chainlink-common
// (pkg/capabilities/registry/client); this package holds only what is specific
// to core: the migration shim for capabilities that are still handed over as
// in-process values, and the selector that decides whether to use the remote
// registry at all.
package creregistry

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	registryclient "github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry/client"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

// LOOPShim adapts capabilities that are still delivered as in-process values to
// the address-addressed registry.
//
// The registry registers by address, so a capability handed over as a value has to
// be served somewhere first. Core receives them one at a time, so this shim gives
// each its own listener and calls AddAt with that address.
//
// Migration story: a LOOP calls Add on core's go-plugin registry and hands over a
// capability value backed by a broker connection into the LOOP's own process. The
// remote registry cannot use that value — it dials, it does not accept broker IDs
// — so the shim serves the capability locally and registers that address instead.
// Once capabilities call Add on the registry directly with their own address
// (registryclient.New plus registryclient.RegisterCapability), this shim has no
// users and can be deleted.
//
// The listener binds to loopback by default: it exists so a co-located registry
// process can reach capabilities in this process, not to publish them.
type LOOPShim struct {
	lggr   logger.Logger
	client core.AddressableCapabilitiesRegistry

	// listenAddr is what each capability's listener binds to. Host-only values
	// (e.g. "127.0.0.1:0") pick a free port per capability.
	listenAddr string

	// listen creates a capability's listener. It exists so tests can supply an
	// in-memory listener instead of binding a port; production always uses
	// net.Listen.
	listen func(addr string) (net.Listener, error)

	mu      sync.Mutex
	served  map[string]*servedCapability
	stopped bool
}

// servedCapability is one capability's listener and gRPC server.
//
// One server per capability, not one shared server: grpc.Server's service
// registry is immutable once Serve starts, so a shared server could only ever
// host the capabilities known before the first Add. Per-capability servers also
// make the address model literal — one capability, one address.
type servedCapability struct {
	srv  *grpc.Server
	addr string
	// cap is retained so a re-registration can tell whether the previous one is
	// still live, the same question baseRegistry answers before replacing.
	cap capabilities.BaseCapability
}

var _ core.CapabilitiesRegistry = (*LOOPShim)(nil)

// NewLOOPShim wraps an address-based registry so Add accepts in-process capability
// values, serving each at an address of its own and registering that.
func NewLOOPShim(lggr logger.Logger, client core.AddressableCapabilitiesRegistry, listenAddr string) *LOOPShim {
	if listenAddr == "" {
		listenAddr = "127.0.0.1:0"
	}
	return &LOOPShim{
		lggr:       logger.Named(lggr, "CRERegistryLOOPShim"),
		client:     client,
		listenAddr: listenAddr,
		listen:     func(addr string) (net.Listener, error) { return net.Listen("tcp", addr) },
		served:     map[string]*servedCapability{},
	}
}

// Add serves cap on its own local listener and registers that address.
func (s *LOOPShim) Add(ctx context.Context, cap capabilities.BaseCapability) error {
	info, err := cap.Info(ctx)
	if err != nil {
		return fmt.Errorf("could not read capability info: %w", err)
	}

	// Both paths share one ID namespace, so the same duplicate rule applies to
	// each: a re-registration is how a restarted host re-announces itself and
	// replaces a dead entry, while a live entry means a second claimant and is
	// refused. This mirrors baseRegistry, which core uses when the registry is
	// local.
	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		return errors.New("shim is closed")
	}
	prev, replacing := s.served[info.ID]
	if replacing && !isDeadCapability(prev.cap) {
		s.mu.Unlock()
		return fmt.Errorf("%w: id %s is already served by this shim", registry.ErrCapabilityAlreadyExists, info.ID)
	}
	s.mu.Unlock()

	srv := grpc.NewServer()
	if err = registryclient.RegisterCapability(s.lggr, srv, cap, info.CapabilityType); err != nil {
		srv.Stop()
		return fmt.Errorf("could not serve capability %s: %w", info.ID, err)
	}

	lis, err := s.listen(s.listenAddr)
	if err != nil {
		srv.Stop()
		return fmt.Errorf("failed to listen on %s: %w", s.listenAddr, err)
	}
	addr := lis.Addr().String()

	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		srv.Stop()
		_ = lis.Close()
		return errors.New("shim is closed")
	}
	s.served[info.ID] = &servedCapability{srv: srv, addr: addr, cap: cap}
	s.mu.Unlock()

	s.retire(info.ID, prev, replacing, addr)

	go func() {
		s.lggr.Infow("serving local capability for the CRE registry", "capabilityID", info.ID, "address", addr)
		if err := srv.Serve(lis); err != nil {
			s.lggr.Errorw("capability shim gRPC server stopped", "capabilityID", info.ID, "err", err)
		}
	}()

	if err = s.client.AddAt(ctx, info.ID, info.CapabilityType, addr); err != nil {
		// Registration failed, so nothing will ever dial this listener.
		s.mu.Lock()
		delete(s.served, info.ID)
		s.mu.Unlock()
		srv.Stop()
		return err
	}
	return nil
}

// Remove deregisters the capability and stops serving it.
func (s *LOOPShim) Remove(ctx context.Context, id string) error {
	err := s.client.Remove(ctx, id)

	s.mu.Lock()
	served, ok := s.served[id]
	delete(s.served, id)
	s.mu.Unlock()

	// Stop serving even if the deregistration call failed: leaving the listener up
	// would keep answering for a capability the caller asked to drop.
	if ok {
		served.srv.GracefulStop()
	}
	return err
}

// Close stops every capability server this shim started.
func (s *LOOPShim) Close() error {
	s.mu.Lock()
	s.stopped = true
	served := s.served
	s.served = map[string]*servedCapability{}
	s.mu.Unlock()

	for _, c := range served {
		c.srv.GracefulStop()
	}
	return nil
}

// retire stops what the new registration superseded, once the replacement is
// already serving.
func (s *LOOPShim) retire(id string, prev *servedCapability, replacing bool, addr string) {
	if !replacing {
		return
	}
	s.lggr.Infow("replaced a dead capability registration",
		"capabilityID", id, "previousAddress", prev.addr, "address", addr)
	prev.srv.GracefulStop()
}

// isDeadCapability reports whether a previously registered capability is known to
// be unreachable.
//
// The states mirror baseRegistry's replace condition. A capability that cannot
// report a state is treated as live, so an unrelated second claimant on the same
// ID is still refused.
func isDeadCapability(c capabilities.BaseCapability) bool {
	sg, ok := c.(registry.StateGetter)
	if !ok {
		return false
	}
	switch sg.GetState() {
	case connectivity.Shutdown, connectivity.TransientFailure, connectivity.Idle:
		return true
	default:
		return false
	}
}

// --- pass-through of everything except Add/Remove ---

func (s *LOOPShim) Get(ctx context.Context, id string) (capabilities.BaseCapability, error) {
	return s.client.Get(ctx, id)
}

func (s *LOOPShim) GetTrigger(ctx context.Context, id string) (capabilities.TriggerCapability, error) {
	return s.client.GetTrigger(ctx, id)
}

func (s *LOOPShim) GetExecutable(ctx context.Context, id string) (capabilities.ExecutableCapability, error) {
	return s.client.GetExecutable(ctx, id)
}

func (s *LOOPShim) List(ctx context.Context) ([]capabilities.BaseCapability, error) {
	return s.client.List(ctx)
}

func (s *LOOPShim) LocalNode(ctx context.Context) (capabilities.Node, error) {
	return s.client.LocalNode(ctx)
}

func (s *LOOPShim) NodeByPeerID(ctx context.Context, peerID ragetypes.PeerID) (capabilities.Node, error) {
	return s.client.NodeByPeerID(ctx, peerID)
}

func (s *LOOPShim) ConfigForCapability(ctx context.Context, capabilityID string, donID uint32) (capabilities.CapabilityConfiguration, error) {
	return s.client.ConfigForCapability(ctx, capabilityID, donID)
}

func (s *LOOPShim) DONsForCapability(ctx context.Context, capabilityID string) ([]capabilities.DONWithNodes, error) {
	return s.client.DONsForCapability(ctx, capabilityID)
}

func (s *LOOPShim) DONByID(ctx context.Context, donID uint32) (capabilities.DON, error) {
	return s.client.DONByID(ctx, donID)
}

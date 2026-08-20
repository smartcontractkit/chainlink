package creregistry

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
)

// ProxyConfig is the subset of config.CapabilitiesProxy this package needs.
// Declared locally so the selector can be used (and tested) without pulling in
// the whole config interface.
type ProxyConfig interface {
	Enabled() bool
	GRPCPort() uint16
}

// Select returns the CapabilitiesRegistry the node should use.
//
// When the rage proxy is enabled, the node has already delegated its P2P to
// crecore; the registry comes from the same process, so remote capability shims,
// the launcher, and the on-chain read all live on that side of the boundary and
// the node stops running its own registrysyncer. When it is disabled, this
// returns core's in-process registry and nothing changes.
//
// The returned io.Closer is non-nil only in the proxy case; callers must close
// it to release capability connections and the local capability listeners.
func Select(lggr logger.Logger, cfg ProxyConfig) (*capabilities.Registry, func() error, error) {
	if cfg == nil || !cfg.Enabled() {
		return capabilities.NewRegistry(lggr), nil, nil
	}

	// crecore is launched as a LOOP on loopback, so the registry lives at the
	// same address the OCR proxy client already uses.
	addr := fmt.Sprintf("127.0.0.1:%d", cfg.GRPCPort())

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CRE registry client for %s: %w", addr, err)
	}

	// Capabilities registered by this node are served on loopback by the shim, and
	// the proxy itself is a co-located process, so insecure credentials are stated
	// explicitly here rather than defaulted in the client. A deployment that moves
	// capabilities off-host has to change this line, which is the point.
	// addresses is where this process serves the capabilities it registers, filled in by the shim as
	// it opens a listener for each: the registry on the other end holds addresses rather than
	// values, so registering with it means naming where the value can be reached.
	addresses := map[string]string{}
	reg := registry.Local(lggr).WithRemote(conn, addresses,
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Capabilities in this process are still handed over as values, so wrap the
	// registry in the migration shim that serves them at an address. Once every
	// capability registers with its own address, drop the shim and return reg.
	shim := NewLOOPShim(lggr, reg, addresses, "")

	closeFn := func() error {
		shimErr := shim.Close()
		clientErr := reg.Close()
		connErr := conn.Close()
		switch {
		case shimErr != nil:
			return shimErr
		case clientErr != nil:
			return clientErr
		default:
			return connErr
		}
	}

	lggr.Infow("using the CRE registry served by the p2p proxy", "address", addr)
	return capabilities.NewDelegatingRegistry(lggr, shim), closeFn, nil
}

package adapters

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/tunnel"
)

type TunnelAdapter[T any] interface {
	DescribeEndpoints(componentID string, output *T) ([]tunnel.EndpointRef, error)
	RewriteWithBindings(output *T, bindings []tunnel.TunnelBinding) error
}

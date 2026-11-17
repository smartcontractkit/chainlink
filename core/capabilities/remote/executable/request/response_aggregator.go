package request

import (
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// ResponseAggregator decides the outcome of a client request from peer replies.
// All methods return (nil, nil) until decided, (resp, nil) to resolve the request
// with resp, or (nil, err) to resolve the request with err. An observationError
// is returned to the OnMessage caller without deciding the request.
type ResponseAggregator interface {
	// OnUnmarshalable: the peer returned an OK response whose payload could not be unmarshalled.
	OnUnmarshalable(peer p2ptypes.PeerID, err error) (*commoncap.CapabilityResponse, error)
	// OnError: the peer returned a non-OK response.
	OnError(peer p2ptypes.PeerID, transportErr types.Error, errMsg string) (*commoncap.CapabilityResponse, error)
	// OnResponse: the peer returned an OK response with a valid payload.
	OnResponse(peer p2ptypes.PeerID, resp commoncap.CapabilityResponse) (*commoncap.CapabilityResponse, error)
}

// AggregatorFactory returns nil to keep the default quorum aggregator.
type AggregatorFactory func(req commoncap.CapabilityRequest, remoteDON commoncap.DON) ResponseAggregator

// observationError tells ClientRequest to surface err to the OnMessage caller
// without resolving the request.
type observationError struct{ err error }

var _ error = observationError{}

func (e observationError) Error() string { return e.err.Error() }
func (e observationError) Unwrap() error { return e.err }

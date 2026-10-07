package request

import (
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// ResponseAggregator replaces the identical-response quorum for requests whose
// peers return different payloads. Both methods return (nil, nil) until decided.
type ResponseAggregator interface {
	OnResponse(peer p2ptypes.PeerID, resp commoncap.CapabilityResponse) (*commoncap.CapabilityResponse, error)
	OnError(peer p2ptypes.PeerID, errMsg string) (*commoncap.CapabilityResponse, error)
}

// AggregatorFactory returns nil to keep the default quorum.
type AggregatorFactory func(req commoncap.CapabilityRequest, remoteDON commoncap.DON) ResponseAggregator

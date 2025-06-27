package utils

import (
	"context"

	httpserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/http/server"
	consensusserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/consensus/server"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/fakes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// NewFakeCapabilities builds faked capabilities, then registers them with the capability registry.
func NewFakeComputeCapabilities(ctx context.Context, lggr logger.Logger, registry *capabilities.Registry) ([]services.Service, error) {
	caps := make([]services.Service, 0)

	// Consensus
	fakeConsensusNoDAG := fakes.NewFakeConsensusNoDAG(lggr)
	fakeConsensusServer := consensusserver.NewConsensusServer(fakeConsensusNoDAG)
	if err := registry.Add(ctx, fakeConsensusServer); err != nil {
		return nil, err
	}
	caps = append(caps, fakeConsensusServer)

	// HTTP Action
	httpAction := fakes.NewDirectHTTPAction(lggr)
	httpActionServer := httpserver.NewClientServer(httpAction)
	if err := registry.Add(ctx, httpActionServer); err != nil {
		return nil, err
	}
	caps = append(caps, httpActionServer)

	return caps, nil
}

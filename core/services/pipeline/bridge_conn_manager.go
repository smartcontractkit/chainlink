package pipeline

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
)

type BridgeConnManager interface {
	GetObservation(ctx context.Context, lggr logger.Logger, bridge bridges.BridgeType, requestData MapParam) ([]byte, error)
}

type dummyBridgeConnManager struct{}

func NewBridgeConnManager() BridgeConnManager {
	return &dummyBridgeConnManager{}
}

func (m *dummyBridgeConnManager) GetObservation(ctx context.Context, lggr logger.Logger, bridge bridges.BridgeType, requestData MapParam) ([]byte, error) {
	return nil, errors.Errorf("bridge connection manager transport is not implemented for bridge %q", bridge.Name.String())
}

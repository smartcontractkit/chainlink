package gateway

import (
	"context"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// ConnectionManager holds all connections between Gateway and Nodes.
type ConnectionManager interface {
	job.ServiceCtx
}

type DONConnectionManager interface {
	SetHandler(handler Handler)
	SendToNode(ctx context.Context, nodeAddress string, msg *Message) error
}

type connectionManager struct {
	dons map[string]DONConnectionManager
	lggr logger.Logger
}

type donConnectionManager struct {
	donConfig *DONConfig
	handler   Handler
	codec     Codec
	mu        sync.Mutex
}

func NewConnectionManager(dons map[string]DONConnectionManager, lggr logger.Logger) ConnectionManager {
	return &connectionManager{
		dons: dons,
		lggr: lggr,
	}
}

func (m *connectionManager) Start(context.Context) error {
	m.lggr.Info("starting connection manager")
	return nil
}

func (m *connectionManager) Close() error {
	m.lggr.Info("closing connection manager")
	return nil
}

func NewDONConnectionManager(donConfig *DONConfig, codec Codec) DONConnectionManager {
	return &donConnectionManager{donConfig: donConfig, codec: codec}
}

func (m *donConnectionManager) SetHandler(handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handler = handler
}

func (m *donConnectionManager) SendToNode(ctx context.Context, nodeAddress string, msg *Message) error {
	return nil
}

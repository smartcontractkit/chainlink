package gateway

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	gw_net "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// ConnectionManager holds all connections between Gateway and Nodes.
type ConnectionManager interface {
	job.ServiceCtx
	network.ConnectionAcceptor

	DONConnectionManager(donId string) DONConnectionManager
}

type DONConnectionManager interface {
	SetHandler(handler Handler)
	SendToNode(ctx context.Context, nodeAddress string, msg *Message) error
}

type connectionManager struct {
	utils.StartStopOnce

	dons     map[string]*donConnectionManager
	wsServer gw_net.WebSocketServer
	lggr     logger.Logger
}

type donConnectionManager struct {
	donConfig *DONConfig
	handler   Handler
	codec     Codec
	mu        sync.Mutex
}

func NewConnectionManager(config *GatewayConfig, codec Codec, lggr logger.Logger) (ConnectionManager, error) {
	dons := make(map[string]*donConnectionManager)
	for _, donConfig := range config.Dons {
		donConfig := donConfig
		if donConfig.DonId == "" {
			return nil, errors.New("empty DON ID")
		}
		_, ok := dons[donConfig.DonId]
		if ok {
			return nil, fmt.Errorf("duplicate DON ID %s", donConfig.DonId)
		}
		dons[donConfig.DonId] = &donConnectionManager{donConfig: &donConfig, codec: codec}
	}
	connMgr := &connectionManager{
		dons: dons,
		lggr: lggr.Named("ConnectionManager"),
	}
	wsServer := gw_net.NewWebSocketServer(&config.NodeServerConfig, connMgr, lggr)
	connMgr.wsServer = wsServer
	return connMgr, nil
}

func (m *connectionManager) DONConnectionManager(donId string) DONConnectionManager {
	return m.dons[donId]
}

func (m *connectionManager) Start(ctx context.Context) error {
	return m.StartOnce("ConnectionManager", func() error {
		m.lggr.Info("starting connection manager")
		return m.wsServer.Start(ctx)
	})
}

func (m *connectionManager) Close() error {
	return m.StopOnce("ConnectionManager", func() (err error) {
		m.lggr.Info("closing connection manager")
		return m.wsServer.Close()
	})
}

func (m *connectionManager) StartHandshake(authHeader []byte) (attemptId string, challenge []byte, err error) {
	m.lggr.Debug("StartHandshake")
	// TODO (FUN-469): implement the handshake (extract node address, validate signature, create a new attempt, generate challenge)
	return "", []byte{}, nil
}

func (m *connectionManager) FinalizeHandshake(attemptId string, response []byte, conn *websocket.Conn) error {
	m.lggr.Debug("FinalizeHandshake", attemptId, response)
	// TODO (FUN-469): implement the handshake (validate signature, add a new connection)
	return nil
}

func (m *connectionManager) AbortHandshake(attemptId string) {
	m.lggr.Debug("AbortHandshake", attemptId)
	// TODO (FUN-469): implement the handshake (clear cached attempt)
}

func (m *donConnectionManager) SetHandler(handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handler = handler
}

func (m *donConnectionManager) SendToNode(ctx context.Context, nodeAddress string, msg *Message) error {
	return nil
}

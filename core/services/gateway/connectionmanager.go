package gateway

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

var promKeepalivesSent = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "gateway_keepalives_sent",
	Help: "Metric to track the number of successful keepalive ping messages per DON",
}, []string{"don_id"})

// ConnectionManager holds all connections between Gateway and Nodes.
type ConnectionManager interface {
	job.ServiceCtx
	network.ConnectionAcceptor

	DONConnectionManager(donId string) *donConnectionManager
	GetPort() int
}

type connectionManager struct {
	services.StateMachine

	config             *config.ConnectionManagerConfig
	dons               map[string]*donConnectionManager
	wsServer           network.WebSocketServer
	clock              clockwork.Clock
	connAttempts       map[string]*connAttempt
	connAttemptCounter uint64
	connAttemptsMu     sync.Mutex
	lggr               logger.Logger
}

func (m *connectionManager) HealthReport() map[string]error {
	hr := map[string]error{m.Name(): m.Healthy()}
	for _, d := range m.dons {
		for _, n := range d.nodes {
			services.CopyHealth(hr, n.conn.HealthReport())
		}
	}
	return hr
}

func (m *connectionManager) Name() string { return m.lggr.Name() }

type donConnectionManager struct {
	donConfig  *config.DONConfig
	nodes      map[string]*nodeState
	handler    handlers.Handler
	codec      api.Codec
	closeWait  sync.WaitGroup
	shutdownCh services.StopChan
	lggr       logger.Logger
}

type nodeState struct {
	name string
	conn network.WSConnectionWrapper
}

// immutable
type connAttempt struct {
	nodeState   *nodeState
	nodeAddress string
	challenge   network.ChallengeElems
	timestamp   uint32
}

func NewConnectionManager(gwConfig *config.GatewayConfig, clock clockwork.Clock, lggr logger.Logger) (ConnectionManager, error) {
	codec := &api.JsonRPCCodec{}
	dons := make(map[string]*donConnectionManager)
	for _, donConfig := range gwConfig.Dons {
		donConfig := donConfig
		if donConfig.DonId == "" {
			return nil, errors.New("empty DON ID")
		}
		_, ok := dons[donConfig.DonId]
		if ok {
			return nil, fmt.Errorf("duplicate DON ID %s", donConfig.DonId)
		}
		nodes := make(map[string]*nodeState)
		for _, nodeConfig := range donConfig.Members {
			nodeAddress := strings.ToLower(nodeConfig.Address)
			_, ok := nodes[nodeAddress]
			if ok {
				return nil, fmt.Errorf("duplicate node address %s in DON %s", nodeAddress, donConfig.DonId)
			}
			connWrapper := network.NewWSConnectionWrapper(lggr)
			if connWrapper == nil {
				return nil, fmt.Errorf("error creating WSConnectionWrapper for node %s", nodeAddress)
			}
			nodes[nodeAddress] = &nodeState{
				name: nodeConfig.Name,
				conn: connWrapper,
			}
		}
		dons[donConfig.DonId] = &donConnectionManager{
			donConfig:  &donConfig,
			codec:      codec,
			nodes:      nodes,
			shutdownCh: make(chan struct{}),
			lggr:       lggr.Named("DONConnectionManager." + donConfig.DonId),
		}
	}
	connMgr := &connectionManager{
		config:       &gwConfig.ConnectionManagerConfig,
		dons:         dons,
		connAttempts: make(map[string]*connAttempt),
		clock:        clock,
		lggr:         lggr.Named("ConnectionManager"),
	}
	wsServer := network.NewWebSocketServer(&gwConfig.NodeServerConfig, connMgr, lggr)
	connMgr.wsServer = wsServer
	return connMgr, nil
}

func (m *connectionManager) DONConnectionManager(donId string) *donConnectionManager {
	return m.dons[donId]
}

func (m *connectionManager) Start(ctx context.Context) error {
	return m.StartOnce("ConnectionManager", func() error {
		m.lggr.Info("starting connection manager")
		for _, donConnMgr := range m.dons {
			donConnMgr.closeWait.Add(len(donConnMgr.nodes))
			for nodeAddress, nodeState := range donConnMgr.nodes {
				if err := nodeState.conn.Start(ctx); err != nil {
					return err
				}
				go donConnMgr.readLoop(nodeAddress, nodeState)
			}
			donConnMgr.closeWait.Add(1)
			go donConnMgr.keepaliveLoop(m.config.HeartbeatIntervalSec)
		}
		return m.wsServer.Start(ctx)
	})
}

func (m *connectionManager) Close() error {
	return m.StopOnce("ConnectionManager", func() (err error) {
		m.lggr.Info("closing connection manager")
		err = multierr.Combine(err, m.wsServer.Close())
		for _, donConnMgr := range m.dons {
			close(donConnMgr.shutdownCh)
			for _, nodeState := range donConnMgr.nodes {
				nodeState.conn.Close()
			}
		}
		for _, donConnMgr := range m.dons {
			donConnMgr.closeWait.Wait()
		}
		return
	})
}

func (m *connectionManager) StartHandshake(authHeader []byte) (attemptId string, challenge []byte, err error) {
	m.lggr.Debug("StartHandshake")
	authHeaderElems, signer, err := network.UnpackSignedAuthHeader(authHeader)
	if err != nil {
		return "", nil, multierr.Append(network.ErrAuthHeaderParse, err)
	}
	nodeAddress := "0x" + hex.EncodeToString(signer)
	donConnMgr, ok := m.dons[authHeaderElems.DonId]
	if !ok {
		return "", nil, network.ErrAuthInvalidDonId
	}
	nodeState, ok := donConnMgr.nodes[nodeAddress]
	if !ok {
		return "", nil, network.ErrAuthInvalidNode
	}
	if authHeaderElems.GatewayId != m.config.AuthGatewayId {
		return "", nil, network.ErrAuthInvalidGateway
	}
	nowTs := uint32(m.clock.Now().Unix())
	ts := authHeaderElems.Timestamp
	if ts < nowTs-m.config.AuthTimestampToleranceSec || nowTs+m.config.AuthTimestampToleranceSec < ts {
		return "", nil, network.ErrAuthInvalidTimestamp
	}
	attemptId, challenge, err = m.newAttempt(nodeState, nodeAddress, ts)
	if err != nil {
		return "", nil, err
	}
	return attemptId, challenge, nil
}

func (m *connectionManager) newAttempt(nodeSt *nodeState, nodeAddress string, timestamp uint32) (string, []byte, error) {
	challengeBytes := make([]byte, m.config.AuthChallengeLen)
	_, err := rand.Read(challengeBytes)
	if err != nil {
		return "", nil, err
	}
	challenge := network.ChallengeElems{Timestamp: timestamp, GatewayId: m.config.AuthGatewayId, ChallengeBytes: challengeBytes}
	m.connAttemptsMu.Lock()
	defer m.connAttemptsMu.Unlock()
	m.connAttemptCounter++
	newId := fmt.Sprintf("%s_%d", nodeAddress, m.connAttemptCounter)
	m.connAttempts[newId] = &connAttempt{nodeState: nodeSt, nodeAddress: nodeAddress, challenge: challenge, timestamp: timestamp}
	return newId, network.PackChallenge(&challenge), nil
}

func (m *connectionManager) FinalizeHandshake(attemptId string, response []byte, conn *websocket.Conn) error {
	m.lggr.Debugw("FinalizeHandshake", "attemptId", attemptId)
	m.connAttemptsMu.Lock()
	attempt, ok := m.connAttempts[attemptId]
	delete(m.connAttempts, attemptId)
	m.connAttemptsMu.Unlock()
	if !ok {
		return network.ErrChallengeAttemptNotFound
	}
	signer, err := common.ExtractSigner(response, network.PackChallenge(&attempt.challenge))
	if err != nil || attempt.nodeAddress != "0x"+hex.EncodeToString(signer) {
		return network.ErrChallengeInvalidSignature
	}
	if conn != nil {
		conn.SetPongHandler(func(data string) error {
			m.lggr.Debugw("received keepalive pong from node", "nodeAddress", attempt.nodeAddress)
			return nil
		})
	}
	attempt.nodeState.conn.Reset(conn)
	m.lggr.Infof("node %s connected", attempt.nodeAddress)
	return nil
}

func (m *connectionManager) AbortHandshake(attemptId string) {
	m.lggr.Debugw("AbortHandshake", "attemptId", attemptId)
	m.connAttemptsMu.Lock()
	defer m.connAttemptsMu.Unlock()
	delete(m.connAttempts, attemptId)
}

func (m *connectionManager) GetPort() int {
	return m.wsServer.GetPort()
}

func (m *donConnectionManager) SetHandler(handler handlers.Handler) {
	m.handler = handler
}

func (m *donConnectionManager) SendToNode(ctx context.Context, nodeAddress string, msg *api.Message) error {
	if msg == nil {
		return errors.New("nil message")
	}
	data, err := m.codec.EncodeRequest(msg)
	if err != nil {
		return fmt.Errorf("error encoding request for node %s: %w", nodeAddress, err)
	}
	nodeState := m.nodes[nodeAddress]
	if nodeState == nil {
		return fmt.Errorf("node %s not found", nodeAddress)
	}
	return nodeState.conn.Write(ctx, websocket.BinaryMessage, data)
}

func (m *donConnectionManager) readLoop(nodeAddress string, nodeState *nodeState) {
	ctx, _ := m.shutdownCh.NewCtx()
	for {
		select {
		case <-m.shutdownCh:
			m.closeWait.Done()
			return
		case item := <-nodeState.conn.ReadChannel():
			msg, err := m.codec.DecodeResponse(item.Data)
			if err != nil {
				m.lggr.Errorw("parse error when reading from node", "nodeAddress", nodeAddress, "err", err)
				break
			}
			if err = msg.Validate(); err != nil {
				m.lggr.Errorw("message validation error when reading from node", "nodeAddress", nodeAddress, "err", err)
				break
			}
			if msg.Body.Sender != nodeAddress {
				m.lggr.Errorw("message sender mismatch when reading from node", "nodeAddress", nodeAddress, "sender", msg.Body.Sender)
				break
			}
			err = m.handler.HandleNodeMessage(ctx, msg, nodeAddress)
			if err != nil {
				m.lggr.Error("error when calling HandleNodeMessage ", err)
			}
		}
	}
}

func (m *donConnectionManager) keepaliveLoop(intervalSec uint32) {
	ctx, _ := m.shutdownCh.NewCtx()
	defer m.closeWait.Done()

	if intervalSec == 0 {
		m.lggr.Errorw("keepalive interval is 0, keepalive disabled", "donID", m.donConfig.DonId)
		return
	}
	m.lggr.Infow("starting keepalive loop", "donID", m.donConfig.DonId)

	keepaliveTicker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	defer keepaliveTicker.Stop()

	for {
		select {
		case <-m.shutdownCh:
			return
		case <-keepaliveTicker.C:
			errorCount := 0
			for nodeAddress, nodeState := range m.nodes {
				err := nodeState.conn.Write(ctx, websocket.PingMessage, []byte{})
				if err != nil {
					m.lggr.Debugw("unable to send keepalive ping to node", "nodeAddress", nodeAddress, "name", nodeState.name, "donID", m.donConfig.DonId, "err", err)
					errorCount++
				}
			}
			promKeepalivesSent.WithLabelValues(m.donConfig.DonId).Set(float64(len(m.nodes) - errorCount))
			m.lggr.Infow("sent keepalive pings to nodes", "donID", m.donConfig.DonId, "errCount", errorCount)
		}
	}
}

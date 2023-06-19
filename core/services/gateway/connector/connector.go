package connector

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// GatewayConnector is a component run by Nodes to connect to a set of Gateways.
type GatewayConnector interface {
	job.ServiceCtx
	network.ConnectionInitiator

	SendToGateway(ctx context.Context, gatewayId string, msg *api.Message) error
}

// Signer implementation needs to be provided by a GatewayConnector user (node)
// in order to sign handshake messages with node's private key.

//go:generate mockery --quiet --name Signer --output ./mocks/ --case=underscore
type Signer interface {
	// Sign keccak256 hash of data.
	Sign(data ...[]byte) ([]byte, error)
}

// GatewayConnector user (node) implements application logic in the Handler interface.

//go:generate mockery --quiet --name GatewayConnectorHandler --output ./mocks/ --case=underscore
type GatewayConnectorHandler interface {
	job.ServiceCtx

	HandleGatewayMessage(gatewayId string, msg *api.Message)
}

type gatewayConnector struct {
	utils.StartStopOnce

	config      *ConnectorConfig
	codec       api.Codec
	clock       utils.Clock
	nodeAddress []byte
	signer      Signer
	handler     GatewayConnectorHandler
	gateways    map[string]*gatewayState
	closeWait   sync.WaitGroup
	shutdownCh  chan struct{}
	lggr        logger.Logger
}

type gatewayState struct {
	conn     network.WSConnectionWrapper
	config   *ConnectorGatewayConfig
	url      *url.URL
	wsClient network.WebSocketClient
}

func NewGatewayConnector(config *ConnectorConfig, signer Signer, handler GatewayConnectorHandler, clock utils.Clock, lggr logger.Logger) (GatewayConnector, error) {
	if signer == nil || handler == nil || clock == nil {
		return nil, errors.New("nil dependency")
	}
	if len(config.DonId) == 0 || len(config.DonId) > int(network.HandshakeDonIdLen) {
		return nil, errors.New("invalid DON ID")
	}
	addressBytes, err := utils.TryParseHex(config.NodeAddress)
	if err != nil {
		return nil, err
	}
	connector := &gatewayConnector{
		config:      config,
		codec:       &api.JsonRPCCodec{},
		clock:       clock,
		nodeAddress: addressBytes,
		signer:      signer,
		handler:     handler,
		shutdownCh:  make(chan struct{}),
		lggr:        lggr,
	}
	gateways := make(map[string]*gatewayState)
	for _, gw := range config.Gateways {
		gw := gw
		_, ok := gateways[gw.Id]
		if ok {
			return nil, fmt.Errorf("duplicate Gateway ID %s", gw.Id)
		}
		parsedURL, err := url.Parse(gw.URL)
		if err != nil {
			return nil, err
		}
		gateway := &gatewayState{
			config:   &gw,
			url:      parsedURL,
			wsClient: network.NewWebSocketClient(config.WsClientConfig, connector, lggr),
		}
		gateways[gw.Id] = gateway
	}
	connector.gateways = gateways
	return connector, nil
}

func (c *gatewayConnector) SendToGateway(ctx context.Context, gatewayId string, msg *api.Message) error {
	data, err := c.codec.EncodeResponse(msg)
	if err != nil {
		return fmt.Errorf("error encoding response for gateway %s: %v", gatewayId, err)
	}
	gateway, ok := c.gateways[gatewayId]
	if !ok {
		return fmt.Errorf("invalid Gateway ID %s", gatewayId)
	}
	if gateway.conn == nil {
		return fmt.Errorf("connector not started")
	}
	return gateway.conn.Write(ctx, websocket.BinaryMessage, data)
}

func (c *gatewayConnector) readLoop(gatewayState *gatewayState) {
	for {
		select {
		case <-c.shutdownCh:
			c.closeWait.Done()
			return
		case item := <-gatewayState.conn.ReadChannel():
			msg, err := c.codec.DecodeRequest(item.Data)
			if err != nil {
				c.lggr.Errorw("parse error when reading from Gateway", "id", gatewayState.config.Id, "err", err)
				break
			}
			c.handler.HandleGatewayMessage(gatewayState.config.Id, msg)
		}
	}
}

func (c *gatewayConnector) reconnectLoop(gatewayState *gatewayState) {
	redialBackoff := utils.NewRedialBackoff()
	ctx, _ := utils.StopChan(c.shutdownCh).NewCtx()
	for {
		conn, err := gatewayState.wsClient.Connect(ctx, gatewayState.url)
		if err != nil {
			c.lggr.Error("connection error")
		} else {
			closeCh := gatewayState.conn.Restart(conn)
			<-closeCh
			c.lggr.Info("connection closed")
			// reset backoff
			redialBackoff = utils.NewRedialBackoff()
		}
		select {
		case <-c.shutdownCh:
			c.closeWait.Done()
			return
		case <-time.After(redialBackoff.Duration()):
			c.lggr.Info("reconnecting ...")
		}
	}
}

func (c *gatewayConnector) Start(ctx context.Context) error {
	return c.StartOnce("GatewayConnector", func() error {
		c.lggr.Info("starting gateway connector")
		if err := c.handler.Start(ctx); err != nil {
			return err
		}
		c.closeWait.Add(2 * len(c.gateways))
		for _, gatewayState := range c.gateways {
			gatewayState := gatewayState
			gatewayState.conn = network.NewWSConnectionWrapper()
			go c.readLoop(gatewayState)
			go c.reconnectLoop(gatewayState)
		}
		return nil
	})
}

func (c *gatewayConnector) Close() error {
	return c.StopOnce("GatewayConnector", func() (err error) {
		c.lggr.Info("closing gateway connector")
		close(c.shutdownCh)
		for _, gatewayState := range c.gateways {
			gatewayState.conn.Close()
		}
		c.closeWait.Wait()
		return c.handler.Close()
	})
}

func (c *gatewayConnector) NewAuthHeader(url *url.URL) ([]byte, error) {
	authHeaderElems := &network.AuthHeaderElems{
		Timestamp:  uint32(c.clock.Now().Unix()),
		DonId:      c.config.DonId,
		GatewayURL: url.String(),
	}
	packedElems := network.Pack(authHeaderElems)
	signature, err := c.signer.Sign(packedElems)
	if err != nil {
		return nil, err
	}
	return append(packedElems, signature...), nil
}

func (c *gatewayConnector) ChallengeResponse(challenge []byte) ([]byte, error) {
	if len(challenge) < c.config.MinHandshakeChallengeLen {
		return nil, errors.New("handshake challenge too short")
	}
	return c.signer.Sign(challenge)
}

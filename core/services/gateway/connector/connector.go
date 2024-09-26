package connector

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/hex"

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

	AddHandler(methods []string, handler GatewayConnectorHandler) error
	SendToGateway(ctx context.Context, gatewayId string, msg *api.Message) error
}

// Signer implementation needs to be provided by a GatewayConnector user (node)
// in order to sign handshake messages with node's private key.
type Signer interface {
	// Sign keccak256 hash of data.
	Sign(data ...[]byte) ([]byte, error)
}

// GatewayConnector user (node) implements application logic in the Handler interface.
type GatewayConnectorHandler interface {
	job.ServiceCtx

	HandleGatewayMessage(ctx context.Context, gatewayId string, msg *api.Message)
}

type gatewayConnector struct {
	services.StateMachine

	config      *ConnectorConfig
	codec       api.Codec
	clock       clockwork.Clock
	nodeAddress []byte
	signer      Signer
	handlers    map[string]GatewayConnectorHandler
	gateways    map[string]*gatewayState
	urlToId     map[string]string
	closeWait   sync.WaitGroup
	shutdownCh  services.StopChan
	lggr        logger.Logger
}

func (c *gatewayConnector) HealthReport() map[string]error {
	m := map[string]error{c.Name(): c.Healthy()}
	for _, g := range c.gateways {
		services.CopyHealth(m, g.conn.HealthReport())
	}
	return m
}

func (c *gatewayConnector) Name() string { return c.lggr.Name() }

type gatewayState struct {
	conn     network.WSConnectionWrapper
	config   ConnectorGatewayConfig
	url      *url.URL
	wsClient network.WebSocketClient
}

func NewGatewayConnector(config *ConnectorConfig, signer Signer, clock clockwork.Clock, lggr logger.Logger) (GatewayConnector, error) {
	if config == nil || signer == nil || clock == nil || lggr == nil {
		return nil, errors.New("nil dependency")
	}
	if len(config.DonId) == 0 || len(config.DonId) > network.HandshakeDonIdLen {
		return nil, errors.New("invalid DON ID")
	}
	addressBytes, err := hex.DecodeString(config.NodeAddress)
	if err != nil {
		return nil, err
	}
	connector := &gatewayConnector{
		config:      config,
		codec:       &api.JsonRPCCodec{},
		clock:       clock,
		nodeAddress: addressBytes,
		signer:      signer,
		handlers:    make(map[string]GatewayConnectorHandler),
		shutdownCh:  make(chan struct{}),
		lggr:        lggr.Named("GatewayConnector"),
	}
	gateways := make(map[string]*gatewayState)
	urlToId := make(map[string]string)
	for _, gw := range config.Gateways {
		gw := gw
		if _, exists := gateways[gw.Id]; exists {
			return nil, fmt.Errorf("duplicate Gateway ID %s", gw.Id)
		}
		if _, exists := urlToId[gw.URL]; exists {
			return nil, fmt.Errorf("duplicate Gateway URL %s", gw.URL)
		}
		parsedURL, err := url.Parse(gw.URL)
		if err != nil {
			return nil, err
		}
		gateway := &gatewayState{
			conn:     network.NewWSConnectionWrapper(lggr),
			config:   gw,
			url:      parsedURL,
			wsClient: network.NewWebSocketClient(config.WsClientConfig, connector, lggr),
		}
		gateways[gw.Id] = gateway
		urlToId[gw.URL] = gw.Id
	}
	connector.gateways = gateways
	connector.urlToId = urlToId
	return connector, nil
}

func (c *gatewayConnector) AddHandler(methods []string, handler GatewayConnectorHandler) error {
	if handler == nil {
		return errors.New("cannot add a nil handler")
	}
	for _, method := range methods {
		if _, exists := c.handlers[method]; exists {
			return fmt.Errorf("handler for method %s already exists", method)
		}
	}
	// add all or nothing
	for _, method := range methods {
		c.handlers[method] = handler
	}
	return nil
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
	ctx, cancel := c.shutdownCh.NewCtx()
	defer cancel()

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
			if err = msg.Validate(); err != nil {
				c.lggr.Errorw("failed to validate message signature", "id", gatewayState.config.Id, "err", err)
				break
			}
			handler, exists := c.handlers[msg.Body.Method]
			if !exists {
				c.lggr.Errorw("no handler for method", "id", gatewayState.config.Id, "method", msg.Body.Method)
				break
			}
			handler.HandleGatewayMessage(ctx, gatewayState.config.Id, msg)
		}
	}
}

func (c *gatewayConnector) reconnectLoop(gatewayState *gatewayState) {
	redialBackoff := utils.NewRedialBackoff()
	ctx, cancel := c.shutdownCh.NewCtx()
	defer cancel()

	for {
		conn, err := gatewayState.wsClient.Connect(ctx, gatewayState.url)
		if err != nil {
			c.lggr.Errorw("connection error", "url", gatewayState.url, "err", err)
		} else {
			c.lggr.Infow("connected successfully", "url", gatewayState.url)
			closeCh := gatewayState.conn.Reset(conn)
			<-closeCh
			c.lggr.Infow("connection closed", "url", gatewayState.url)
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
		for _, gatewayState := range c.gateways {
			gatewayState := gatewayState
			if err := gatewayState.conn.Start(ctx); err != nil {
				return err
			}
			c.closeWait.Add(2)
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
		var errs error
		for _, gatewayState := range c.gateways {
			errs = errors.Join(errs, gatewayState.conn.Close())
		}
		c.closeWait.Wait()
		return errs
	})
}

func (c *gatewayConnector) NewAuthHeader(url *url.URL) ([]byte, error) {
	gatewayId, found := c.urlToId[url.String()]
	if !found {
		return nil, network.ErrAuthInvalidGateway
	}
	authHeaderElems := &network.AuthHeaderElems{
		Timestamp: uint32(c.clock.Now().Unix()),
		DonId:     c.config.DonId,
		GatewayId: gatewayId,
	}
	packedElems := network.PackAuthHeader(authHeaderElems)
	signature, err := c.signer.Sign(packedElems)
	if err != nil {
		return nil, err
	}
	return append(packedElems, signature...), nil
}

func (c *gatewayConnector) ChallengeResponse(url *url.URL, challenge []byte) ([]byte, error) {
	challengeElems, err := network.UnpackChallenge(challenge)
	if err != nil {
		return nil, err
	}
	if len(challengeElems.ChallengeBytes) < c.config.AuthMinChallengeLen {
		return nil, network.ErrChallengeTooShort
	}
	gatewayId, found := c.urlToId[url.String()]
	if !found || challengeElems.GatewayId != gatewayId {
		return nil, network.ErrAuthInvalidGateway
	}
	nowTs := uint32(c.clock.Now().Unix())
	ts := challengeElems.Timestamp
	if ts < nowTs-c.config.AuthTimestampToleranceSec || nowTs+c.config.AuthTimestampToleranceSec < ts {
		return nil, network.ErrAuthInvalidTimestamp
	}
	return c.signer.Sign(challenge)
}

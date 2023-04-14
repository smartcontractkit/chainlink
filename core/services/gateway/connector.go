package gateway

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type GatewayConnectorConfig struct {
	DONID            string
	GatewayAddresses []string
	SignerAddress    string
}

type GatewayConnector interface {
	job.ServiceCtx
}

type gatewayConnector struct {
	config    *GatewayConnectorConfig
	logger    logger.Logger
	handler   Handler
	conn      *websocket.Conn
	closeChan chan bool
}

type connectorCallback struct {
	conn *websocket.Conn
}

func (c *connectorCallback) SendResponse(msg *Message) {
	b, _ := Encode(msg)
	c.conn.WriteMessage(websocket.TextMessage, b)
}

func NewGatewayConnector(config *GatewayConnectorConfig, handler Handler, logger logger.Logger) *gatewayConnector {
	return &gatewayConnector{
		config:    config,
		logger:    logger,
		handler:   handler,
		closeChan: make(chan bool),
	}
}

func (g *gatewayConnector) run(context.Context) error {
	// Connect and maintain persistent connections to all Gateways in the config.
	if len(g.config.GatewayAddresses) != 1 {
		return fmt.Errorf("only one gateway is supported now")
	}
	addr := g.config.GatewayAddresses[0]

	u := url.URL{Scheme: "ws", Host: addr, Path: "/node"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		g.logger.Error("GatewayConnector: couldn't connect to ", addr)
		return err
	}
	g.logger.Info("GatewayConnector: connected to ", addr)
	g.conn = c
	_, helloMsg, err := c.ReadMessage()
	if err != nil {
		return err
	}
	msg, _ := Decode(helloMsg)
	if msg.Method != "hello" {
		return fmt.Errorf("wrong method name %s", msg.Method)
	}
	g.logger.Info("GatewayConnector: received Hello from gateway ", addr)

	var helloResponse Message
	helloResponse.Method = "hello"
	helloResponse.SenderAddress = g.config.SignerAddress
	helloResponse.DonId = g.config.DONID
	respBytes, _ := Encode(&helloResponse)
	c.WriteMessage(websocket.TextMessage, respBytes)
	g.logger.Info("GatewayConnector: sent Hello back")
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				g.logger.Info("GatewayConnector: connection terminated")
				return
			}
			msg, _ := Decode(message)
			if msg.Method == "heartbeat" {
				g.logger.Trace("GatewayConnector: heartbeat")
			} else {
				g.handler.HandleUserMessage(msg, &connectorCallback{conn: c})
			}
		}
	}()
	return nil
}

func (g *gatewayConnector) Start(ctx context.Context) error {
	g.run(ctx)
	return nil
}

func (g *gatewayConnector) Close() error {
	g.conn.Close()
	//g.closeChan <- true
	return nil
}

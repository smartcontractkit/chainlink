package gateway

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type GatewayConfig struct {
	UserEndpointPort uint16
	NodeEndpointPort uint16
	Dons             []GatewayDONConfig
}

type GatewayDONConfig struct {
	DonId       string
	HandlerName string
	Members     []DONMember
}

type DONMember struct {
	Name          string
	SignerAddress string
}

type Gateway interface {
	job.ServiceCtx
}

type gateway struct {
	lggr     logger.SugaredLogger
	server   *http.Server
	connMgrs map[string]ConnectionManager
	handlers map[string]Handler
}

type gatewayCallback struct {
	w    http.ResponseWriter
	done chan bool
}

func (c *gatewayCallback) SendResponse(msg *Message) {
	b, _ := Encode(msg)
	c.w.Write(b)
	c.done <- true
}

var upgrader = websocket.Upgrader{} // use default options

func (g *gateway) node(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	g.lggr.Info("Gateway Incoming Node Connection")
	c.WriteMessage(websocket.TextMessage, []byte(`{"method": "hello", "payload": {"rand": "rand_bytes"}}`))

	_, message, err := c.ReadMessage()
	if err != nil {
		g.lggr.Info("Gateway Read error")
		c.Close()
		return
	}
	msg, err := Decode(message)
	if err != nil {
		g.lggr.Info("Gateway Read error")
		c.Close()
		return
	}
	val, ok := g.connMgrs[msg.DonId]
	if !ok {
		g.lggr.Warn("Gateway Incoming unknown DON ", msg.DonId)
		c.Close()
		return
	}
	if msg.Method == "hello" {
		g.lggr.Info("Gateway Received Hello from node ", msg.SenderAddress, " - accepting!")
		// TODO validate
		val.AddConnection(msg.SenderAddress, c)
	}
}

func (g *gateway) user(w http.ResponseWriter, r *http.Request) {
	g.lggr.Info("Gateway User Message")
	b, err := io.ReadAll(r.Body)
	if err != nil {
		g.lggr.Error("Gateway error reading user message", err)
		// TODO responses here and below
		return
	}
	msg, err := Decode(b)
	if err != nil {
		g.lggr.Error("Gateway error parsing user message", err)
		return
	}
	val, ok := g.handlers[msg.DonId]
	if !ok {
		g.lggr.Error("Gateway unknown DON ", msg.DonId)
		return
	}
	done := make(chan bool)
	cb := &gatewayCallback{w: w, done: done}
	val.HandleUserMessage(msg, cb)
	<-done
}

func (g *gateway) Start(context.Context) error {
	g.lggr.Info("Starting Gateway!")
	http.HandleFunc("/node", g.node)
	http.HandleFunc("/user", g.user)
	g.server = &http.Server{Addr: "localhost:8040"}
	go func() {
		err := g.server.ListenAndServe()
		if err != nil {
			g.lggr.Error("Gateway Server closed with error", err)
		}
	}()
	return nil
}

func (g *gateway) Close() error {
	g.server.Shutdown(context.Background())
	return nil
}

func (g *gateway) Name() string {
	return "Gateway"
}

func NewGateway(lggr logger.SugaredLogger, config *GatewayConfig) Gateway {
	gateway := &gateway{lggr: lggr, connMgrs: make(map[string]ConnectionManager), handlers: make(map[string]Handler)}
	for _, don := range config.Dons {
		lggr.Info("Gateway adding DON: ", don.DonId, " with ", len(don.Members), " members")
		// TODO if don.HandlerName == "functions"
		handler := NewFunctionsHandler(lggr)
		connMgr := NewConnectionManager(&don, handler)
		handler.Init(connMgr, &don)
		gateway.handlers[don.DonId] = handler
		gateway.connMgrs[don.DonId] = connMgr
	}
	return gateway
}

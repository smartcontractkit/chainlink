package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jonboulle/clockwork"
	"github.com/pelletier/go-toml/v2"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

// Script to run Connector outside of the core node.
//
// Usage (without TLS):
//
//	go run run_connector.go --config sample_config.toml
type client struct {
	privateKey *ecdsa.PrivateKey
	connector  connector.GatewayConnector
	lggr       logger.Logger
}

func (h *client) HandleGatewayMessage(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) error {
	msg, err := hc.ValidatedMessageFromReq(req)
	if err != nil {
		return err
	}
	resp, err := hc.ValidatedResponseFromMessage(msg)
	if err != nil {
		return err
	}
	h.lggr.Infof("received message from gateway %s. Echoing back.", gatewayID)
	err = h.connector.SendToGateway(ctx, gatewayID, resp)
	if err != nil {
		h.lggr.Errorw("failed to send to gateway", "id", gatewayID, "err", err)
	}
	return nil
}

func (h *client) Sign(ctx context.Context, data ...[]byte) ([]byte, error) {
	return common.SignData(h.privateKey, data...)
}

func (h *client) Start(ctx context.Context) error {
	return nil
}

func (h *client) ID(ctx context.Context) (string, error) {
	return "test_client", nil
}

func (h *client) Close() error {
	return nil
}

func main() {
	configFile := flag.String("config", "", "Path to TOML config file")
	flag.Parse()

	rawConfig, err := os.ReadFile(*configFile)
	if err != nil {
		fmt.Println("error reading config:", err)
		return
	}

	var cfg connector.ConnectorConfig
	err = toml.Unmarshal(rawConfig, &cfg)
	if err != nil {
		fmt.Println("error parsing config:", err)
		return
	}

	sampleKey, _ := crypto.HexToECDSA("cd47d3fafdbd652dd2b66c6104fa79b372c13cb01f4a4fbfc36107cce913ac1d")
	lggr, err := logger.New()
	if err != nil {
		fmt.Println("error creating logger:", err)
		return
	}
	client := &client{privateKey: sampleKey, lggr: lggr}
	// client acts as a signer here
	connector, _ := connector.NewGatewayConnector(&cfg, client, clockwork.NewRealClock(), lggr)
	err = connector.AddHandler(context.Background(), []string{"test_method"}, client)
	if err != nil {
		fmt.Println("error adding handler:", err)
		return
	}
	client.connector = connector

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	err = connector.Start(ctx)
	if err != nil {
		fmt.Println("error staring connector:", err)
		return
	}

	<-ctx.Done()
	err = connector.Close()
	if err != nil {
		fmt.Println("error closing connector:", err)
		return
	}
}

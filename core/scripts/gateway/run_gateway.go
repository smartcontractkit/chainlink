package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
)

// Script to run Gateway outside of the core node. It works only with simple handlers.
// Any handlers that depend on core services will fail in their factory methods.
//
// Usage without TLS:
//
//	go run run_gateway.go --config sample_config.toml
//
//	curl -X POST -d  '{"jsonrpc":"2.0","method":"test","id":"abcd","params":{"body":{"don_id":"wrong"}}}' http://localhost:8080/user
//
// Usage with TLS:
//
//	openssl req -newkey rsa:2048 -nodes -keyout key.pem -x509 -days 365 -out certificate.pem
//	go run run_gateway.go --config sample_config_tls.toml
//
//	curl -X POST -d  '{"jsonrpc":"2.0","method":"test","id":"abcd","params":{"body":{"don_id":"wrong"}}}' https://localhost:8088/user -k
func main() {
	configFile := flag.String("config", "", "Path to TOML config file")
	flag.Parse()

	rawConfig, err := os.ReadFile(*configFile)
	if err != nil {
		fmt.Println("error reading config:", err)
		return
	}

	var cfg gateway.GatewayConfig
	err = toml.Unmarshal(rawConfig, &cfg)
	if err != nil {
		fmt.Println("error parsing config:", err)
		return
	}

	lggr, _ := logger.NewLogger()

	gw, err := gateway.NewGatewayFromConfig(&cfg, lggr)
	if err != nil {
		fmt.Println("error creating Gateway object:", err)
		return
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	gw.Start(ctx)

	<-ctx.Done()
	gw.Close()
}

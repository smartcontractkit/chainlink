package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

// Script to run Connector outside of the core node.
//
// Usage without TLS:
//
//	go run run_connector.go --url ws://localhost:8081/node
//
// Usage with TLS:
//
//	go run run_connector.go --url wss://localhost:8089/node
type initiator struct {
}

func (*initiator) NewAuthHeader(url *url.URL) []byte {
	fmt.Println("generating new auth header:", url.String())
	return []byte{}
}

func (*initiator) ChallengeResponse(challenge []byte) ([]byte, error) {
	fmt.Println("generating challenge response:", string(challenge))
	return []byte{}, nil
}

func main() {
	urlStr := flag.String("url", "", "Gateway URL")
	flag.Parse()
	url, err := url.Parse(*urlStr)
	if err != nil {
		fmt.Println("error parsing url:", err)
	}

	lggr, _ := logger.NewLogger()

	config := network.WebSocketClientConfig{
		HandshakeTimeoutMillis: 1000,
	}
	client := network.NewWebSocketClient(config, &initiator{}, lggr)
	conn, err := client.Connect(context.Background(), url)
	if err != nil || conn == nil {
		fmt.Println("connection error:", err)
		return
	}

	fmt.Println("connected successfully!")
	if err = conn.Close(); err != nil {
		fmt.Println("error closing connection", err)
	}
}

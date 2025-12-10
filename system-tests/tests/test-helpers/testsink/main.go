package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	// "github.com/smartcontractkit/chainlink/pkg/chiptestsink"
)

func main() {
	var (
		grpcListen = flag.String("grpc-listen", ":9090", "ChipIngress gRPC listen address")
		httpListen = flag.String("http-listen", ":8080", "HTTP API listen address")
		upstream   = flag.String("upstream", "", "optional upstream ChipIngress gRPC endpoint for pass-through")
	)
	flag.Parse()

	cfg := Config{
		GRPCListen:       *grpcListen,
		HTTPListen:       *httpListen,
		UpstreamEndpoint: *upstream,
	}

	srv, err := NewServer(cfg)
	if err != nil {
		log.Fatalf("create chip test sink: %v", err)
	}

	// Run server
	go func() {
		if err := srv.Run(); err != nil {
			log.Printf("[chip-testsink] stopped: %v", err)
		}
	}()

	// Wait for SIGINT/SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[chip-testsink] shutdown error: %v", err)
	}
}

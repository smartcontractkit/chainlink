package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		grpcListen = flag.String("grpc-listen", ":9090", "ChipIngress gRPC listen address")
		httpListen = flag.String("http-listen", ":8080", "HTTP API listen address")
		upstream   = flag.String("upstream", "", "optional upstream ChipIngress gRPC endpoint for pass-through")
		cacheSize  = flag.Int("cache-size", 2000, "maximum number of events to cache")
		daemon     = flag.Bool("daemon", false, "run as daemon without signal handling (for CI environments)")
	)
	flag.Parse()

	cfg := Config{
		GRPCListen:       *grpcListen,
		HTTPListen:       *httpListen,
		UpstreamEndpoint: *upstream,
		CacheSize:        *cacheSize,
	}

	srv, err := NewServer(cfg)
	if err != nil {
		log.Fatalf("create chip test sink: %v", err)
	}

	// In daemon mode, just run forever without signal handling
	if *daemon {
		log.Printf("[chip-testsink] starting in daemon mode on gRPC=%s, HTTP=%s", *grpcListen, *httpListen)
		if err := srv.Run(); err != nil {
			log.Fatalf("[chip-testsink] daemon stopped: %v", err)
		}
		return
	}

	// Normal mode: run server with graceful shutdown
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

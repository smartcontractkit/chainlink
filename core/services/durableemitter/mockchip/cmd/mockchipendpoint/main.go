// mockchipendpoint runs the mock Chip Ingress gRPC endpoint locally so that
// a DurableEmitter under test can connect to it (directly or through an
// `ngrok tcp` tunnel) and the resulting events can be inspected over a small
// HTTP control plane.
//
// Typical usage:
//
//	go run ./core/services/durableemitter/mockchip/cmd/mockchipendpoint \
//	    -grpc :9095 -http :9096
//	# in a separate terminal
//	ngrok tcp 9095        # paste the forwarding address into your DurableEmitter config
//	curl localhost:9096/stats
//	curl -X POST localhost:9096/outage/on
//	curl -X POST localhost:9096/outage/off
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/durableemitter/mockchip"
)

func main() {
	grpcAddr := flag.String("grpc", ":9095", "gRPC listen address for the mock ChipIngress endpoint")
	httpAddr := flag.String("http", ":9096", "HTTP listen address for the inspection / control plane")
	startOutage := flag.Bool("start-outage", false, "Begin in outage mode (every publish fails until /outage/off is called)")
	flag.Parse()

	srv := mockchip.NewServer()
	if *startOutage {
		srv.SetOutage(true)
	}

	grpcResolved, err := srv.Start(*grpcAddr)
	if err != nil {
		log.Fatalf("mockchip: start gRPC: %v", err)
	}

	ctrl := mockchip.NewHTTPController(srv)
	httpResolved, err := ctrl.Start(*httpAddr)
	if err != nil {
		srv.Stop()
		log.Fatalf("mockchip: start HTTP: %v", err)
	}

	fmt.Printf("mockchip: gRPC listening on %s\n", grpcResolved)
	fmt.Printf("mockchip: HTTP control plane on http://%s\n", httpResolved)
	fmt.Printf("mockchip: outage_active=%v\n", srv.OutageActive())
	fmt.Println("mockchip: expose to the internet with `ngrok tcp " + portOf(grpcResolved) + "`")
	fmt.Println("mockchip: send SIGINT/SIGTERM to exit")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = ctrl.Stop(shutdownCtx)
	srv.Stop()
	fmt.Println("mockchip: shutdown complete")
}

func portOf(addr string) string {
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			return addr[i+1:]
		}
	}
	return addr
}

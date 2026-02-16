package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/agent"
	blockchainsets "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/sets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:18080", "agent listen address")
	flag.Parse()

	lggr := zerolog.New(os.Stderr).With().Timestamp().Logger()
	provider := &infra.Provider{Type: infra.Docker}
	server := agent.NewServer(lggr, blockchainsets.NewDeployerSet(lggr, provider))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	lggr.Info().Msgf("starting local CRE agent on %s", *addr)
	if err := agent.Run(ctx, *addr, server); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "agent failed: %v\n", err)
		os.Exit(1)
	}
}

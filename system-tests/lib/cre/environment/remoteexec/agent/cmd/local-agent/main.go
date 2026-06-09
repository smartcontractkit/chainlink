package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	blockchainsets "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/sets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/agent"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func main() {
	os.Exit(run())
}

func run() int {
	defaultAddr := "127.0.0.1:18080"
	if runtimecfg.IsDirectMode() {
		defaultAddr = "0.0.0.0:18080"
	}
	addr := flag.String("addr", defaultAddr, "agent listen address")
	flag.Parse()

	lggr := zerolog.New(os.Stderr).With().Timestamp().Logger()
	commonLogger := cldlogger.NewSingleFileLogger(nil)
	provider := &infra.Provider{Type: infra.Docker}
	server := agent.NewServer(lggr, commonLogger, blockchainsets.NewDeployerSet(lggr, provider))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	lggr.Info().Msgf("starting local CRE agent on %s", *addr)
	if err := agent.Run(ctx, *addr, server); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "agent failed: %v\n", err)
		return 1
	}
	return 0
}

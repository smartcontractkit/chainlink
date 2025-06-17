package main

import (
	"io"
	defaultLog "log"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	tcLog "github.com/testcontainers/testcontainers-go/log"

	"github.com/smartcontractkit/chainlink/integration-tests/docker/cmd/internal"
)

var rootCmd = &cobra.Command{
	Use:   "coreqa",
	Short: "Core QA test tool",
}

func init() {
	rootCmd.AddCommand(internal.StartNodesCmd)

	// Set default log level for non-testcontainer code
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Discard testcontainers logs
	tcLog.SetDefault(defaultLog.New(io.Discard, "", defaultLog.LstdFlags))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("Error")
		os.Exit(1)
	}
}

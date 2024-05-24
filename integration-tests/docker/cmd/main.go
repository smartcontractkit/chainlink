package main

import (
	"io"
	defaultlog "log"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/testcontainers/testcontainers-go"

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
	testcontainers.Logger = defaultlog.New(io.Discard, "", defaultlog.LstdFlags)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("Error")
		os.Exit(1)
	}
}

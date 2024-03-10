package main

import (
	"io"
	defaultlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "test_env",
		Short: "CL cluster docker test env management tool",
	}

	startEnvCmd := &cobra.Command{
		Use:   "start-env",
		Short: "Start new docker test env",
	}
	rootCmd.AddCommand(startEnvCmd)

	startFullEnvCmd := &cobra.Command{
		Use:   "cl-cluster",
		Short: "Basic CL cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Logger = logging.GetLogger(nil, "CORE_DOCKER_ENV_LOG_LEVEL")
			log.Info().Msg("Starting CL cluster test environment..")

			_, err := test_env.NewCLTestEnvBuilder().
				WithGeth().
				WithMockAdapter().
				WithCLNodes(6).
				Build()
			if err != nil {
				return err
			}

			log.Info().Msg("Cl cluster is ready")

			handleExitSignal()

			return nil
		},
	}

	startEnvCmd.AddCommand(startFullEnvCmd)

	// Set default log level for non-testcontainer code
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Discard testcontainers logs
	testcontainers.Logger = defaultlog.New(io.Discard, "", defaultlog.LstdFlags)

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("Error")
		os.Exit(1)
	}
}

func handleExitSignal() {
	// Create a channel to receive exit signals
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)

	log.Info().Msg("Press Ctrl+C to destroy the test environment")

	// Block until an exit signal is received
	<-exitChan
}

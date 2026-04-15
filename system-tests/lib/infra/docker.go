package infra

import (
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func PrintFailedContainerLogs(logger zerolog.Logger, logLinesCount uint64) {
	if err := framework.PrintFailedContainerLogs(logLinesCount); err != nil {
		logger.Error().Err(err).Msg("failed to print failed Docker container logs")
	}
}

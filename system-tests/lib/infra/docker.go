package infra

import (
	"encoding/binary"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	text "github.com/smartcontractkit/chainlink/system-tests/lib/format"
)

func PrintFailedContainerLogs(logger zerolog.Logger, logLinesCount uint64) {
	logStream, lErr := framework.StreamContainerLogs(framework.ExitedCtfContainersListOpts, container.LogsOptions{
		ShowStderr: true,
		Tail:       strconv.FormatUint(logLinesCount, 10),
	})

	if lErr != nil {
		logger.Error().Err(lErr).Msg("failed to stream Docker container logs")
		return
	}

	logger.Error().Msgf("Containers that failed to start: %s", strings.Join(slices.Collect(maps.Keys(logStream)), ", "))
	for cName, ioReader := range logStream {
		content := ""
		header := make([]byte, 8) // Docker stream header is 8 bytes
		for {
			_, err := io.ReadFull(ioReader, header)
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Error().Err(err).Str("Container", cName).Msg("failed to read log stream header")
				break
			}

			// Extract log message size
			msgSize := binary.BigEndian.Uint32(header[4:8])

			// Read the log message
			msg := make([]byte, msgSize)
			_, err = io.ReadFull(ioReader, msg)
			if err != nil {
				logger.Error().Err(err).Str("Container", cName).Msg("failed to read log message")
				break
			}

			content += string(msg)
		}

		content = strings.TrimSpace(content)
		if len(content) > 0 {
			logger.Info().Str("Container", cName).Msgf("Last 100 lines of logs")
			fmt.Println(text.RedText("%s\n", content))
		}
		_ = ioReader.Close() // can't do much about the error here
	}
}

package infra

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/docker/docker/api/types/container"
	dfilter "github.com/docker/docker/api/types/filters"
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

	logger.Error().Msgf("Containers that exited with non-zero codes: %s", strings.Join(slices.Collect(maps.Keys(logStream)), ", "))
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

// CheckContainersForPanics scans Docker container logs for panic-related patterns.
// When a panic is detected, it displays the panic line and up to maxLinesAfterPanic lines following it.
//
// This is useful for debugging test failures where a Docker container may have panicked.
// The function searches for common panic patterns including:
//   - "panic:" - Standard Go panic
//   - "runtime error:" - Runtime errors (nil pointer, index out of bounds, etc.)
//   - "fatal error:" - Fatal errors
//   - "goroutine N [running]" - Stack trace indicators
//
// The function scans all containers in parallel and stops as soon as the first panic is found.
//
// Parameters:
//   - logger: zerolog.Logger for output
//   - maxLinesAfterPanic: Maximum number of lines to show after the panic line (including stack trace).
//     Recommended: 50-200 depending on how much context you want.
//
// Returns:
//   - true if any panics were found in any container, false otherwise
func CheckContainersForPanics(logger zerolog.Logger, maxLinesAfterPanic int) bool {
	// Panic patterns to search for
	panicPatterns := map[string]*regexp.Regexp{
		"panic":         regexp.MustCompile(`(?i)panic:`),                // Go panic
		"runtime error": regexp.MustCompile(`(?i)runtime error:`),        // Runtime errors
		"fatal error":   regexp.MustCompile(`(?i)fatal error:`),          // Fatal errors
		"stack trace":   regexp.MustCompile(`goroutine \d+ \[running\]`), // Stack trace indicator
	}

	// List only Chainlink node containers
	listOpts := container.ListOptions{
		All: true,
		Filters: dfilter.NewArgs(
			dfilter.KeyValuePair{Key: "label", Value: "framework=ctf"},
		),
	}

	logOpts := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true, // Include timestamps for better debugging
	}

	logStream, err := framework.StreamContainerLogs(listOpts, logOpts)
	if err != nil {
		logger.Error().Err(err).Msg("failed to stream container logs for panic detection")
		return false
	}

	// Create context for early cancellation when first panic is found
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to receive panic results
	panicFoundChan := make(chan bool, len(logStream))
	var wg sync.WaitGroup

	// Scan all containers in parallel
	for containerName, reader := range logStream {
		wg.Add(1)
		go func(name string, r io.ReadCloser) {
			defer wg.Done()
			defer r.Close()

			panicFound := scanContainerForPanics(ctx, logger, name, r, panicPatterns, maxLinesAfterPanic)
			if panicFound {
				panicFoundChan <- true
				cancel() // Signal other goroutines to stop
			}
		}(containerName, reader)
	}

	// Wait for all goroutines to finish in a separate goroutine
	go func() {
		wg.Wait()
		close(panicFoundChan)
	}()

	// Check if any panic was found
	for range panicFoundChan {
		return true // Return as soon as first panic is found
	}

	return false
}

// scanContainerForPanics scans a single container's log stream for panic patterns
// It checks the context for cancellation to enable early termination when another goroutine finds a panic
func scanContainerForPanics(ctx context.Context, logger zerolog.Logger, containerName string, reader io.Reader, patterns map[string]*regexp.Regexp, maxLinesAfter int) bool {
	scanner := bufio.NewScanner(reader)
	// Increase buffer size to handle large log lines
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	var logLines []string
	lineNum := 0
	panicLineNum := -1
	patternNameFound := ""

	// Read all lines and detect panic
	for scanner.Scan() {
		// Check if context is cancelled (another goroutine found a panic)
		select {
		case <-ctx.Done():
			return false // Stop scanning, another container already found a panic
		default:
		}

		line := scanner.Text()
		lineNum++

		// If we found a panic, collect remaining lines up to maxLinesAfter
		if panicLineNum >= 0 {
			logLines = append(logLines, line)
			// Stop reading once we have enough context after the panic
			if lineNum >= panicLineNum+maxLinesAfter+1 {
				break
			}
			continue
		}

		// Still searching for panic - store all lines
		logLines = append(logLines, line)

		// Check if this line contains a panic pattern
		for patternName, pattern := range patterns {
			if pattern.MatchString(line) {
				patternNameFound = patternName
				panicLineNum = lineNum - 1 // Store index (0-based)
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error().Err(err).Str("Container", containerName).Msg("error reading container logs")
		return false
	}

	// If panic was found, display it with context
	if panicLineNum >= 0 {
		logger.Error().
			Str("Container", containerName).
			Int("PanicLineNumber", panicLineNum+1).
			Msgf("🔥 %s DETECTED in container logs", strings.ToUpper(patternNameFound))

		// Calculate range to display
		startLine := panicLineNum
		endLine := min(len(logLines), panicLineNum+maxLinesAfter+1)

		// Build the output
		var output strings.Builder
		output.WriteString(fmt.Sprintf("\n%s\n", strings.Repeat("=", 80)))
		output.WriteString(fmt.Sprintf("%s FOUND IN CONTAINER: %s (showing %d lines from panic)\n", strings.ToUpper(patternNameFound), containerName, endLine-startLine))
		output.WriteString(strings.Repeat("=", 80) + "\n")

		for i := startLine; i < endLine; i++ {
			output.WriteString(logLines[i] + "\n")
		}

		output.WriteString(strings.Repeat("=", 80))

		fmt.Println(text.RedText("%s\n", output.String()))
		return true
	}

	return false
}

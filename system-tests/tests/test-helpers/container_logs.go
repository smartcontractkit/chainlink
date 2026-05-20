package helpers

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"testing"

	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

func clNodeContainerNames(t *testing.T, testEnv *ttypes.TestEnvironment) []string {
	t.Helper()

	names := make(map[string]struct{})
	for _, nodeSet := range testEnv.Config.NodeSets {
		if nodeSet.Out == nil {
			continue
		}
		for _, clNode := range nodeSet.Out.CLNodes {
			if name := clNode.Node.ContainerName; name != "" {
				names[name] = struct{}{}
			}
		}
	}
	require.NotEmpty(t, names, "no container names found in test environment")
	out := make([]string, 0, len(names))
	for name := range names {
		out = append(out, name)
	}
	slices.Sort(out)
	return out
}

func AssertNodeLogs(t *testing.T, testEnv *ttypes.TestEnvironment, needle string) {
	t.Helper()

	targetNames := make(map[string]struct{})
	for _, name := range clNodeContainerNames(t, testEnv) {
		targetNames[name] = struct{}{}
	}

	logStreams, err := framework.StreamContainerLogs(
		client.ContainerListOptions{All: true},
		client.ContainerLogsOptions{ShowStdout: true, ShowStderr: true},
	)
	require.NoError(t, err)

	found := false
	for containerName, reader := range logStreams {
		if _, ok := targetNames[containerName]; !ok {
			_ = reader.Close()
			continue
		}
		content, readErr := readContainerLogs(reader)
		if readErr != nil {
			framework.L.Warn().Str("container", containerName).Err(readErr).Msg("could not read container logs")
			continue
		}
		if strings.Contains(content, needle) {
			found = true
			framework.L.Info().Str("container", containerName).Msg("confirmed: " + needle)
		}
	}
	assert.True(t, found, "expected at least one node container log to contain %q", needle)
}

// readContainerLogs decodes a Docker multiplexed log stream into plain text.
// framework.StreamContainerLogs returns this format; the framework decoder is not exported.
func readContainerLogs(r io.ReadCloser) (string, error) {
	defer func() { _ = r.Close() }()

	var buf strings.Builder
	if err := decodeDockerLogStream(&buf, r); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func decodeDockerLogStream(dst io.Writer, r io.Reader) error {
	header := make([]byte, 8)
	for {
		_, err := io.ReadFull(r, header)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read log stream header: %w", err)
		}

		msgSize := binary.BigEndian.Uint32(header[4:8])
		msg := make([]byte, msgSize)
		if _, err = io.ReadFull(r, msg); err != nil {
			return fmt.Errorf("read log message: %w", err)
		}
		if _, err = dst.Write(msg); err != nil {
			return fmt.Errorf("write log message: %w", err)
		}
	}
}

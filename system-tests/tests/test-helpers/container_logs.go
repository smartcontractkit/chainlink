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

// clNodeContainerNames returns sorted Docker container names for every Chainlink node in testEnv.
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

// nodesetContainerNames returns sorted Docker container names for the nodeset named nodesetName.
func nodesetContainerNames(t *testing.T, testEnv *ttypes.TestEnvironment, nodesetName string) []string {
	t.Helper()

	names := make([]string, 0)
	for _, nodeSet := range testEnv.Config.NodeSets {
		if nodeSet.Name != nodesetName || nodeSet.Out == nil {
			continue
		}
		for _, clNode := range nodeSet.Out.CLNodes {
			if name := clNode.Node.ContainerName; name != "" {
				names = append(names, name)
			}
		}
	}
	require.NotEmptyf(t, names, "no container names found for nodeset %q", nodesetName)
	slices.Sort(names)
	return names
}

// assertContainerLogs scans stdout/stderr of containerNames and checks whether needle appears.
func assertContainerLogs(t *testing.T, containerNames []string, needle string, wantFound bool) {
	t.Helper()

	targetNames := make(map[string]struct{}, len(containerNames))
	for _, name := range containerNames {
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
			framework.L.Info().Str("container", containerName).Str("needle", needle).Bool("want_found", wantFound).Msg("container log match")
		}
	}

	if wantFound {
		assert.True(t, found, "expected at least one of %v to contain %q", containerNames, needle)
		return
	}
	assert.False(t, found, "expected none of %v to contain %q", containerNames, needle)
}

// AssertNodeLogs requires needle to appear in at least one Chainlink node container log.
func AssertNodeLogs(t *testing.T, testEnv *ttypes.TestEnvironment, needle string) {
	t.Helper()
	assertContainerLogs(t, clNodeContainerNames(t, testEnv), needle, true)
}

// AssertContainerLogsForNodeset requires needle in at least one container log for nodesetName.
func AssertContainerLogsForNodeset(t *testing.T, testEnv *ttypes.TestEnvironment, nodesetName, needle string) {
	t.Helper()
	assertContainerLogs(t, nodesetContainerNames(t, testEnv, nodesetName), needle, true)
}

// AssertContainerLogsAbsentForNodeset requires needle in no container logs for nodesetName.
func AssertContainerLogsAbsentForNodeset(t *testing.T, testEnv *ttypes.TestEnvironment, nodesetName, needle string) {
	t.Helper()
	assertContainerLogs(t, nodesetContainerNames(t, testEnv, nodesetName), needle, false)
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

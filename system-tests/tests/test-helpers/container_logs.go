package helpers

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"testing"
	"time"

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

// containerLogLines returns every log line containing needle from the current logs
// of containerNames. It is the shared scanning core: the needle assertions check for
// a non-empty result, and callers that need the lines themselves filter them further.
func containerLogLines(t *testing.T, containerNames []string, needle string) []string {
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

	var matches []string
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
		for line := range strings.Lines(content) {
			if strings.Contains(line, needle) {
				matches = append(matches, line)
				framework.L.Info().Str("container", containerName).Str("needle", needle).Msg("container log match")
			}
		}
	}
	return matches
}

// containerLogsContain reports whether needle appears in the logs of any of containerNames,
// as they stand at the moment of the call.
func containerLogsContain(t *testing.T, containerNames []string, needle string) bool {
	t.Helper()

	return len(containerLogLines(t, containerNames, needle)) > 0
}

// assertContainerLogs scans stdout/stderr of containerNames and checks whether needle appears.
func assertContainerLogs(t *testing.T, containerNames []string, needle string, wantFound bool) {
	t.Helper()

	found := containerLogsContain(t, containerNames, needle)
	if wantFound {
		assert.True(t, found, "expected at least one of %v to contain %q", containerNames, needle)
		return
	}
	assert.False(t, found, "expected none of %v to contain %q", containerNames, needle)
}

// ContainerLogLinesForNodeset returns every log line containing needle from the
// current logs of the nodeset's containers. Unlike the needle assertions, it returns
// the lines so callers can filter further (e.g. by a workflow ID whose position in the
// line is not stable).
func ContainerLogLinesForNodeset(t *testing.T, testEnv *ttypes.TestEnvironment, nodesetName, needle string) []string {
	t.Helper()

	return containerLogLines(t, nodesetContainerNames(t, testEnv, nodesetName), needle)
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

// RequireContainerLogsForNodesetEventually waits for needle to appear in at least one container
// log of nodesetName, rescanning until timeout. Use it instead of AssertContainerLogsForNodeset
// when the line is written asynchronously and may still be in flight when the test reaches it.
func RequireContainerLogsForNodesetEventually(t *testing.T, testEnv *ttypes.TestEnvironment, nodesetName, needle string, timeout, interval time.Duration) {
	t.Helper()

	containerNames := nodesetContainerNames(t, testEnv, nodesetName)
	require.Eventually(t, func() bool {
		return containerLogsContain(t, containerNames, needle)
	}, timeout, interval, "expected at least one of %v to contain %q within %s", containerNames, needle, timeout)
}

// NodesetContainerNames returns sorted Docker container names for the nodeset named nodesetName.
func NodesetContainerNames(t *testing.T, testEnv *ttypes.TestEnvironment, nodesetName string) []string {
	t.Helper()
	return nodesetContainerNames(t, testEnv, nodesetName)
}

// ContainerLogLineMatchesForNodeset reports whether any single log line of nodesetName's
// containers satisfies match. Use it when several values must appear on the same line.
func ContainerLogLineMatchesForNodeset(t *testing.T, testEnv *ttypes.TestEnvironment, nodesetName string, match func(line string) bool) bool {
	t.Helper()

	targetNames := make(map[string]struct{})
	for _, name := range nodesetContainerNames(t, testEnv, nodesetName) {
		targetNames[name] = struct{}{}
	}

	logStreams, err := framework.StreamContainerLogs(
		client.ContainerListOptions{All: true},
		client.ContainerLogsOptions{ShowStdout: true, ShowStderr: true},
	)
	require.NoError(t, err)

	found := false
	for containerName, reader := range logStreams {
		if _, ok := targetNames[containerName]; !ok || found {
			_ = reader.Close()
			continue
		}
		content, readErr := readContainerLogs(reader)
		if readErr != nil {
			framework.L.Warn().Str("container", containerName).Err(readErr).Msg("could not read container logs")
			continue
		}
		for line := range strings.SplitSeq(content, "\n") {
			if match(line) {
				found = true
				break
			}
		}
	}
	return found
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

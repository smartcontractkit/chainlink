package helpers

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

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

// waitForStandardCapabilityStartup polls the given container logs until every container has
// emitted "Started standard capabilities" with capabilityLabelledName in the same log line. Container logs are
// read in parallel for speed. The caller is responsible for selecting which containers to check —
// different capabilities run on different node sets (e.g. evm-worker-1337 on workflow nodes,
// evm-worker-2337 on capabilities nodes), so the caller must resolve the correct set beforehand.
//
// "Started standard capabilities" is logged by standardcapabilities/standard_capabilities.go after
// Initialise() and Infos() both succeed. It is not specific to any capability type.
//
// Background: on startup a race between CRE.RegistrySyncerORM persisting the local DON registry
// and the LOOP plugin calling GetLocalNode() can leave 3 of 4 nodes with a permanently unregistered
// capability. Waiting for this log on every expected container ensures the capability is serving
// before any dependent workflow is deployed.
func waitForStandardCapabilityStartup(t *testing.T, containerNames []string, capabilityLabelledName string, timeout time.Duration) {
	t.Helper()
	if len(containerNames) == 0 {
		framework.L.Warn().Msgf("WaitForStandardCapabilityStartup: no containers to check for capability %s, skipping", capabilityLabelledName)
		return
	}

	expected := make(map[string]struct{}, len(containerNames))
	for _, name := range containerNames {
		expected[name] = struct{}{}
	}

	deadline := time.Now().Add(timeout)
	for {
		logStreams, err := framework.StreamContainerLogs(
			client.ContainerListOptions{All: true},
			client.ContainerLogsOptions{ShowStdout: true, ShowStderr: true},
		)
		if err != nil {
			framework.L.Warn().Err(err).Msg("WaitForStandardCapabilityStartup: error streaming container logs, retrying")
			time.Sleep(5 * time.Second)
			continue
		}

		// Read all relevant container logs in parallel.
		var mu sync.Mutex
		ready := make(map[string]struct{})
		var wg sync.WaitGroup
		for containerName, reader := range logStreams {
			if _, ok := expected[containerName]; !ok {
				_ = reader.Close()
				continue
			}
			wg.Add(1)
			go func(name string, r io.ReadCloser) {
				defer wg.Done()
				content, readErr := readContainerLogs(r)
				if readErr != nil {
					framework.L.Warn().Str("container", name).Err(readErr).Msg("WaitForStandardCapabilityStartup: could not read logs")
					return
				}

				for line := range strings.SplitSeq(content, "\n") {
					if strings.Contains(line, "Started standard capabilities") && strings.Contains(line, capabilityLabelledName) {
						mu.Lock()
						ready[name] = struct{}{}
						mu.Unlock()
						return
					}
				}
			}(containerName, reader)
		}
		wg.Wait()

		if len(ready) >= len(expected) {
			framework.L.Info().Msgf("WaitForStandardCapabilityStartup: %s ready on all %d expected containers", capabilityLabelledName, len(expected))
			return
		}

		framework.L.Info().Msgf("WaitForStandardCapabilityStartup: %s ready on %d/%d containers", capabilityLabelledName, len(ready), len(expected))

		if time.Now().After(deadline) {
			missing := make([]string, 0, len(expected)-len(ready))
			for name := range expected {
				if _, ok := ready[name]; !ok {
					missing = append(missing, name)
				}
			}
			require.FailNowf(t, "capability startup timeout",
				"standard capability %s did not start within %s on containers: %v", capabilityLabelledName, timeout, missing)
		}
		time.Sleep(5 * time.Second)
	}
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

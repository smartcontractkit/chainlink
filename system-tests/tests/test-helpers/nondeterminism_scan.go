package helpers

import (
	"strings"

	"github.com/moby/moby/client"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// NonDeterminismNeedles are the node-log substrings that indicate two nodes in the
// same DON computed different results — the signal the mixed-env topology surfaces.
// They are silent when every node runs identical code. Single source of truth shared
// by the smoke-suite TestMain and the CI gate (cmd/mixed-env-nondeterminism-check).
// See core/scripts/cre/environment/docs/mixed-env.md.
var NonDeterminismNeedles = []string{
	// libocr OCR3 consensus: a peer's report/commit signature failed to verify
	// against this node's own computed report (report-attestation + commit phases).
	"This is commonly caused by non-determinism",
	// DON2DON remote-capability request aggregation (server side).
	"received messages with the same id and different payloads",
	// DON2DON remote-capability response aggregation (client side).
	"received multiple unique responses for the same request",
	"response quorum unreachable",
}

// NeedleHit records a forbidden log substring found in a container's logs.
type NeedleHit struct {
	Container string
	Needle    string
}

// ScanContainersForNeedles streams stdout/stderr of every Docker container and
// returns each (container, needle) pair whose logs contain the needle substring.
//
// Unlike assertContainerLogs, it takes no *testing.T, so it can be called from a
// package TestMain after m.Run() — while the shared CRE containers are still
// alive — to fail the run on a non-determinism marker even if all tests passed.
//
// It is best-effort: containers whose logs cannot be read are logged and skipped
// rather than failing the whole scan.
func ScanContainersForNeedles(needles []string) ([]NeedleHit, error) {
	logStreams, err := framework.StreamContainerLogs(
		client.ContainerListOptions{All: true},
		client.ContainerLogsOptions{ShowStdout: true, ShowStderr: true},
	)
	if err != nil {
		return nil, err
	}

	var hits []NeedleHit
	for containerName, reader := range logStreams {
		content, readErr := readContainerLogs(reader) // closes reader
		if readErr != nil {
			framework.L.Warn().Str("container", containerName).Err(readErr).Msg("could not read container logs")
			continue
		}
		for _, needle := range needles {
			if strings.Contains(content, needle) {
				hits = append(hits, NeedleHit{Container: containerName, Needle: needle})
			}
		}
	}
	return hits, nil
}

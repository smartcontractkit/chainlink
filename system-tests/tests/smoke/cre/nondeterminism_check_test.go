package cre

import (
	"fmt"
	"os"
	"strings"
	"testing"

	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

// nonDeterminismFailurePrefix is emitted (and shows up in CI logs / JUnit output)
// when the check trips, so the failure reason is unambiguous.
const nonDeterminismFailurePrefix = "Non-Determinism introduced"

// TestMain wraps the CRE smoke suite. For the mixed-env topology it scans every
// node's logs for non-determinism markers after all tests have run (the shared
// containers are still alive at this point) and fails the run — even if every test
// passed — when a marker is found. For all other topologies it is a no-op passthrough.
func TestMain(m *testing.M) {
	code := m.Run()

	if nonDeterminismCheckEnabled() {
		if found := reportNonDeterminism(); found && code == 0 {
			code = 1
		}
	}

	os.Exit(code)
}

// nonDeterminismCheckEnabled restricts the scan to the mixed-env topology (or an
// explicit opt-in). Single-image runs never emit these lines, so there is no reason
// to scan them.
func nonDeterminismCheckEnabled() bool {
	if strings.EqualFold(os.Getenv("CRE_NONDETERMINISM_CHECK"), "true") {
		return true
	}
	return strings.Contains(strings.ToLower(os.Getenv("TOPOLOGY_NAME")), "mixed-env")
}

// reportNonDeterminism scans all container logs for the markers and prints any
// hits. It returns true if at least one marker was found. A scan error is treated
// as "not found" (best-effort) so we never fail a run just because logs were
// unreadable.
func reportNonDeterminism() bool {
	hits, err := t_helpers.ScanContainersForNeedles(t_helpers.NonDeterminismNeedles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s check skipped: could not read container logs: %v\n", nonDeterminismFailurePrefix, err)
		return false
	}
	if len(hits) == 0 {
		fmt.Println("mixed-env non-determinism check: no markers found")
		return false
	}

	fmt.Printf("FAIL: %s — mixed-env detected disagreement between PR and develop nodes:\n", nonDeterminismFailurePrefix)
	for _, h := range hits {
		fmt.Printf("  - container=%s marker=%q\n", h.Container, h.Needle)
	}
	return true
}

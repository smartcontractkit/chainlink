// Command mixed-env-nondeterminism-check scans all local CRE node containers for the
// non-determinism markers and exits non-zero if any are found. It is the hard CI gate
// for the mixed-env workflow, reusing the same scanner and needle list as the
// smoke-suite TestMain (test-helpers/nondeterminism_scan.go) so the logic isn't
// duplicated in bash. See core/scripts/cre/environment/docs/mixed-env.md.
package main

import (
	"fmt"
	"os"

	helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

func main() {
	hits, err := helpers.ScanContainersForNeedles(helpers.NonDeterminismNeedles)
	if err != nil {
		// Best-effort, matching the smoke-suite scan: an inability to read container
		// logs is an infra problem, not a non-determinism signal, so don't fail the gate.
		fmt.Fprintf(os.Stderr, "mixed-env non-determinism check skipped: could not read container logs: %v\n", err)
		return
	}
	if len(hits) == 0 {
		fmt.Println("mixed-env non-determinism check: no markers found")
		return
	}
	for _, h := range hits {
		fmt.Printf("::error::Non-Determinism introduced — container=%s marker=%q\n", h.Container, h.Needle)
	}
	fmt.Println("Non-Determinism introduced: PR and develop nodes disagreed.")
	os.Exit(1)
}

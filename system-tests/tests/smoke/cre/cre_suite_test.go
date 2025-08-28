package cre

import (
	"testing"
)

/*
To execute tests locally start the local CRE first:
Inside `core/scripts/cre/environment` directory
 1. Ensure the necessary capabilities (i.e. readcontract, http-trigger, http-action) are listed in the environment configuration
 2. Run: `go run . env start && ctf obs up && ctf bs up` to start env + observability + blockscout.
 3. Execute the tests in `system-tests/tests/smoke/cre`: `go test -timeout 15m -run ^Test_CRE_Suite$`.
*/
func Test_CRE_Suite(t *testing.T) {
	testEnv := SetupTestEnvironment(t)

	// WARNING: currently we can't run these tests in parallel, because each test rebuilds environment structs and that includes
	// logging into CL node with GraphQL API, which allows only 1 session per user at a time.
	t.Run("[v1] CRE Suite", func(t *testing.T) {
		// requires `readcontract`, `cron`
		t.Run("[v1] CRE Proof of Reserve (PoR) Test", func(t *testing.T) {
			ExecutePoRTest(t, testEnv)
		})
	})

	t.Run("[v2] CRE Suite", func(t *testing.T) {
		t.Run("[v2] vault DON test", func(t *testing.T) {
			ExecuteVaultTest(t, testEnv)
		})

		t.Run("[v2] HTTP trigger and action test", func(t *testing.T) {
			// requires `http_trigger`, `http_action`
			ExecuteHTTPTriggerActionTest(t, testEnv)
		})

		t.Run("[v2] DON Time test", func(t *testing.T) {
			const skipReason = "Implement smoke test - https://smartcontract-it.atlassian.net/browse/CAPPL-1028"
			t.Skipf("Skipping test for the following reason: %s", skipReason)
		})

		t.Run("[v2] Beholder test", func(t *testing.T) {
			ExecuteBeholderTest(t, testEnv)
		})

		t.Run("[v2] Consensus test", func(t *testing.T) {
			executeConsensusTest(t, testEnv)
		})
	})
}

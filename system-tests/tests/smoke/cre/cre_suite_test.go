package cre

import (
	"testing"
)

/*
To execute tests locally start the local CRE first:
Inside `core/scripts/cre/environment` directory
 1. Ensure the necessary capabilities (i.e. readcontract, http-trigger, http-action) are listed in the environment configuration
 2. Identify the appropriate topology that you want to test
 3. Stop and clear any existing environment: `go run . env stop -a`
 4. Run: `go run . env start -t <topology> && ./bin/ctf obs up` to start env + observability
 5. Optionally run blockscout `./bin/ctf bs up`
 6. Execute the tests in `system-tests/tests/smoke/cre` with CTF_CONFIG set to the corresponding topology file:
    `export  CTF_CONFIGS=../../../../core/scripts/cre/environment/configs/<topology>.toml; go test -timeout 15m -run ^Test_CRE_Suite$`.
*/
func Test_CRE_Suite(t *testing.T) {
	testEnv := SetupTestEnvironment(t)
	priceProvider, porWfCfg := beforePoRTest(t, testEnv)

	// WARNING: currently we can't run these tests in parallel, because each test rebuilds environment structs and that includes
	// logging into CL node with GraphQL API, which allows only 1 session per user at a time.
	t.Run("[v1] CRE Suite", func(t *testing.T) {
		// requires `readcontract`, `cron`
		t.Run("[v1] CRE Proof of Reserve (PoR) Test", func(t *testing.T) {
			porWfCfg.WorkflowFileLocation = "../../../../core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/main.go"
			porWfCfg.WorkflowName = "por-workflow"
			ExecutePoRTest(t, testEnv, priceProvider, porWfCfg)
		})
	})

	t.Run("[v2] CRE Suite", func(t *testing.T) {
		t.Run("[v2] vault DON test", func(t *testing.T) {
			// Skip till we figure out and fix the issues with environment startup on this test
			t.Skip("Skipping test for the following reason: Skip till the errors with topology TopologyWorkflowGatewayCapabilities are fixed: https://smartcontract-it.atlassian.net/browse/PRIV-160")
			ExecuteVaultTest(t, testEnv)
		})

		t.Run("[v2] HTTP trigger and action test", func(t *testing.T) {
			// requires `http_trigger`, `http_action`
			ExecuteHTTPTriggerActionTest(t, testEnv)
		})

		t.Run("[v2] DON Time test", func(t *testing.T) {
			t.Skipf("Skipping test for the following reason: Implement smoke test - https://smartcontract-it.atlassian.net/browse/CAPPL-1028")
		})

		t.Run("[v2] Beholder test", func(t *testing.T) {
			ExecuteBeholderTest(t, testEnv)
		})

		t.Run("[v2] Consensus test", func(t *testing.T) {
			executeConsensusTest(t, testEnv)
		})
	})
}

func Test_withV2Registries(t *testing.T) {
	t.Run("[v1] CRE Proof of Reserve (PoR) Test", func(t *testing.T) {
		const skipReason = "Integrate v2 registry contracts in local CRE/test setup - https://smartcontract-it.atlassian.net/browse/CRE-635"
		t.Skipf("Skipping test for the following reason: %s", skipReason)
		flags := []string{"--with-contracts-version", "v2"}

		testEnv := SetupTestEnvironment(t, flags...)
		priceProvider, wfConfig := beforePoRTest(t, testEnv)
		wfConfig.WorkflowFileLocation = "../../../../core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/main.go"
		wfConfig.WorkflowName = "por-workflow"
		ExecutePoRTest(t, testEnv, priceProvider, wfConfig)
	})
}

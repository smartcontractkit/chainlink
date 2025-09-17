package cre

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
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
	testEnv := SetupTestEnvironmentWithConfig(t, getDefaultTestConfig(t))
	// WARNING: currently we can't run these tests in parallel, because each test rebuilds environment structs and that includes
	// logging into CL node with GraphQL API, which allows only 1 session per user at a time.
	t.Run("[v1] CRE Suite", func(t *testing.T) {
		// requires `readcontract`, `cron`
		t.Run("[v1] CRE Proof of Reserve (PoR) Test", func(t *testing.T) {
			priceProvider, porWfCfg := beforePoRTest(t, testEnv, "por-workflowV1", PoRWFV1Location)
			ExecutePoRTest(t, testEnv, priceProvider, porWfCfg)
		})
	})

	t.Run("[v2] CRE Suite", func(t *testing.T) {
		t.Run("[v2] vault DON test", func(t *testing.T) {
			ExecuteVaultTest(t, testEnv)
		})

		t.Run("[v2] Cron (Beholder) happy path", func(t *testing.T) {
			ExecuteCronBeholderTest(t, testEnv)
		})

		// negative tests for cron
		// TODO: move to a separate package
		for _, tCase := range cronInvalidSchedulesTests {
			testName := fmt.Sprintf("[v2] Cron (Beholder) fails when schedule is %s (%s)", tCase.name, tCase.invalidSchedule)
			t.Run(testName, func(t *testing.T) {
				CronBeholderFailWithInvalidScheduleTest(t, testEnv, tCase.invalidSchedule)
			})
		}

		t.Run("[v2] HTTP trigger and action test", func(t *testing.T) {
			t.Skip("Skipping flaky test https://chainlink-core.slack.com/archives/C07GQNPVBB5/p1757085817724369")
			// requires `http_trigger`, `http_action`
			ExecuteHTTPTriggerActionTest(t, testEnv)
		})

		t.Run("[v2] DON Time test", func(t *testing.T) {
			ExecuteDonTimeTest(t, testEnv)
		})

		t.Run("[v2] Billing test", func(t *testing.T) {
			ExecuteBillingTest(t, testEnv)
		})

		t.Run("[v2] Consensus test", func(t *testing.T) {
			executeConsensusTest(t, testEnv)
		})
	})
}

func Test_CRE_Suite_EVM(t *testing.T) {
	testEnv := SetupTestEnvironmentWithConfig(t, getDefaultTestConfig(t))

	// TODO remove this when OCR works properly with multiple chains in Local CRE
	testEnv.WrappedBlockchainOutputs = []*cre.WrappedBlockchainOutput{testEnv.WrappedBlockchainOutputs[0]}
	t.Run("[v2] EVM Write Test", func(t *testing.T) {
		priceProvider, porWfCfg := beforePoRTest(t, testEnv, "por-workflowV2", PoRWFV2Location)
		porWfCfg.FeedIDs = []string{porWfCfg.FeedIDs[0]}
		ExecutePoRTest(t, testEnv, priceProvider, porWfCfg)
	})

	t.Run("[v2] EVM Read test", func(t *testing.T) {
		executeEVMReadTest(t, testEnv)
	})
}

func Test_withV2Registries(t *testing.T) {
	t.Run("[v1] CRE Proof of Reserve (PoR) Test", func(t *testing.T) {
		flags := []string{"--with-contracts-version", "v2"}
		testEnv := SetupTestEnvironmentWithConfig(t, getDefaultTestConfig(t), flags...)
		priceProvider, wfConfig := beforePoRTest(t, testEnv, "por-workflow", PoRWFV1Location)
		ExecutePoRTest(t, testEnv, priceProvider, wfConfig)
	})
}

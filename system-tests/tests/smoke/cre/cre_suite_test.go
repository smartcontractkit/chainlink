package cre

import (
	"testing"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

//////////// SMOKE TESTS /////////////
// target happy path and sanity checks
// all other tests (e.g. edge cases, negative conditions)
// should go to a `regression` package.
/////////////////////////////////////

/*
To execute tests locally start the local CRE first:
Inside `core/scripts/cre/environment` directory
 1. Ensure the necessary capabilities (i.e. readcontract, http-trigger, http-action) are listed in the environment configuration
 2. Identify the appropriate topology that you want to test
 3. Stop and clear any existing environment: `go run . env stop -a`
 4. Run: `CTF_CONFIGS=<path-to-your-topology-config> go run . env start && ./bin/ctf obs up` to start env + observability
 5. Optionally run blockscout `./bin/ctf bs up`
 6. Execute the tests in `system-tests/tests/smoke/cre` with CTF_CONFIG set to the corresponding topology file:
    `export  CTF_CONFIGS=../../../../core/scripts/cre/environment/configs/<topology>.toml; go test -timeout 15m -run ^Test_CRE_Suite$`.
*/
func Test_CRE_Suite_V1(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))
	// WARNING: currently we can't run these tests in parallel, because each test rebuilds environment structs and that includes
	// logging into CL node with GraphQL API, which allows only 1 session per user at a time.
	t.Run("[v1] CRE Suite", func(t *testing.T) {
		// requires `readcontract`, `cron`
		t.Run("[v1] Proof of Reserve (PoR) Test", func(t *testing.T) {
			priceProvider, porWfCfg := beforePoRTest(t, testEnv, "por-workflowV1", PoRWFV1Location)
			ExecutePoRTest(t, testEnv, priceProvider, porWfCfg)
		})
	})
}

func Test_CRE_Suite_V1_Tron(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-don-tron.toml"))

	t.Run("[v1] Tron Write Test with PoR", func(t *testing.T) {
		priceProvider, porWfCfg := beforePoRTest(t, testEnv, "por-workflowV1", PoRWFV1Location)
		ExecutePoRTest(t, testEnv, priceProvider, porWfCfg)
	})
}

func Test_CRE_Suite_V1_SecureMint(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-solana-don.toml"))

	t.Run("[v1] SecureMint Test with PoR", func(t *testing.T) {
		ExecuteSecureMintTest(t, testEnv)
	})
}

//////////// V2 TESTS /////////////
/*
To execute tests with v2 contracts start the local CRE first:
 1. Inside `core/scripts/cre/environment` directory: `go run . env restart --with-beholder --with-contracts-version v2`
 2. Execute the tests in `system-tests/tests/smoke/cre` with CTF_CONFIG set to the corresponding topology file:
    `export  CTF_CONFIGS=../../../../core/scripts/cre/environment/configs/<topology>.toml; go test -timeout 15m -run ^Test_CRE_Suite$`.
*/
func Test_CRE_Suite_V2(t *testing.T) {
	flags := []string{"--with-contracts-version", "v2"}
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), flags...)

	t.Run("[v2] CRE Proof of Reserve (PoR) Test", func(t *testing.T) {
		priceProvider, wfConfig := beforePoRTest(t, testEnv, "por-workflow", PoRWFV1Location)
		ExecutePoRTest(t, testEnv, priceProvider, wfConfig)
	})

	t.Run("[v2] Vault DON test", func(t *testing.T) {
		ExecuteVaultTest(t, testEnv)
	})

	t.Run("[v2] Cron (Beholder) happy path", func(t *testing.T) {
		ExecuteCronBeholderTest(t, testEnv)
	})

	t.Run("[v2] HTTP trigger and action test", func(t *testing.T) {
		ExecuteHTTPTriggerActionTest(t, testEnv)
	})

	t.Run("[v2] DON Time test", func(t *testing.T) {
		ExecuteDonTimeTest(t, testEnv)
	})

	t.Run("[v2] Billing test", func(t *testing.T) {
		ExecuteBillingTest(t, testEnv)
	})

	t.Run("[v2] Consensus test", func(t *testing.T) {
		ExecuteConsensusTest(t, testEnv)
	})
}

func Test_CRE_Suite_V2_EVM(t *testing.T) {
	flags := []string{"--with-contracts-version", "v2"}
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), flags...)

	// TODO: remove this when OCR works properly with multiple chains in Local CRE
	testEnv.WrappedBlockchainOutputs = []*cre.WrappedBlockchainOutput{testEnv.WrappedBlockchainOutputs[0]}
	t.Run("[v2] EVM Write happy path test", func(t *testing.T) {
		priceProvider, porWfCfg := beforePoRTest(t, testEnv, "por-workflowV2", PoRWFV2Location)
		porWfCfg.FeedIDs = []string{porWfCfg.FeedIDs[0]}
		ExecutePoRTest(t, testEnv, priceProvider, porWfCfg)
	})

	t.Run("[v2] EVM Read happy path test", func(t *testing.T) {
		ExecuteEVMReadTest(t, testEnv)
	})
}

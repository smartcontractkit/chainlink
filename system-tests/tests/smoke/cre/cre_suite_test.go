package cre

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/quarantine"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

//////////// SMOKE TESTS /////////////
// target happy path and sanity checks
// all other tests (e.g. edge cases, negative conditions)
// should go to a `regression` package
/////////////////////////////////////

var v2RegistriesFlags = []string{"--with-contracts-version", "v2"}

/*
To execute tests locally start the local CRE first:
Inside `core/scripts/cre/environment` directory
 1. Ensure the necessary capabilities (i.e. readcontract, http-trigger, http-action) are listed in the environment configuration
 2. Identify the appropriate topology that you want to test
 3. Stop and clear any existing environment: `go run . env stop -a`
 4. Run: `CTF_CONFIGS=<path-to-your-topology-config> go run . env start && ./bin/ctf obs up` to start env + observability
 5. Optionally run blockscout `./bin/ctf bs up`
 6. Execute the tests in `system-tests/tests/smoke/cre`: `go test -timeout 15m -run "^Test_CRE_V2"`.
*/
func Test_CRE_V1_Proof_Of_Reserve(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))
	// WARNING: currently we can't run these tests in parallel, because each test rebuilds environment structs and that includes
	// logging into CL node with GraphQL API, which allows only 1 session per user at a time.

	// requires `readcontract`, `cron`
	priceProvider, porWfCfg := beforePoRTest(t, testEnv, "por-workflowV1", PoRWFV1Location)
	ExecutePoRTest(t, testEnv, priceProvider, porWfCfg, false)
}

func Test_CRE_V1_Tron(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-don-tron.toml"))

	priceProvider, porWfCfg := beforePoRTest(t, testEnv, "por-workflowV1", PoRWFV1Location)
	ExecutePoRTest(t, testEnv, priceProvider, porWfCfg, false)
}

func Test_CRE_V1_SecureMint(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-don-solana.toml"))

	ExecuteSecureMintTest(t, testEnv)
}

// TODO: Move Billing tests to v2 Registries
func Test_CRE_V1_Billing_EVM_Write(t *testing.T) {
	quarantine.Flaky(t, "DX-1911")
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))

	// TODO remove this when OCR works properly with multiple chains in Local CRE
	testEnv.WrappedBlockchainOutputs = []*cre.WrappedBlockchainOutput{testEnv.WrappedBlockchainOutputs[0]}

	require.NoError(
		t,
		startBillingStackIfIsNotRunning(t, testEnv.TestConfig.RelativePathToRepoRoot, testEnv.TestConfig.EnvironmentDirPath, testEnv),
		"failed to start Billing stack",
	)

	priceProvider, porWfCfg := beforePoRTest(t, testEnv, "por-workflowV2-billing", PoRWFV2Location)
	porWfCfg.FeedIDs = []string{porWfCfg.FeedIDs[0]}
	ExecutePoRTest(t, testEnv, priceProvider, porWfCfg, true)
}

func Test_CRE_V1_Billing_Cron_Beholder(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))

	// TODO remove this when OCR works properly with multiple chains in Local CRE
	testEnv.WrappedBlockchainOutputs = []*cre.WrappedBlockchainOutput{testEnv.WrappedBlockchainOutputs[0]}

	require.NoError(
		t,
		startBillingStackIfIsNotRunning(t, testEnv.TestConfig.RelativePathToRepoRoot, testEnv.TestConfig.EnvironmentDirPath, testEnv),
		"failed to start Billing stack",
	)

	ExecuteBillingTest(t, testEnv)
}

//////////// V2 TESTS /////////////
/*
To execute tests with v2 contracts start the local CRE first:
 1. Inside `core/scripts/cre/environment` directory: `go run . env restart --with-beholder --with-contracts-version v2`
 2. Execute the tests in `system-tests/tests/smoke/cre`: `go test -timeout 15m -run "^Test_CRE_V2"`.
*/
func Test_CRE_V2_Suite(t *testing.T) {
	quarantine.Flaky(t, "DX-2002")
	topology := os.Getenv("TOPOLOGY_NAME")

	t.Run("[v2] Proof Of Reserve - "+topology, func(t *testing.T) {
		// TODO: Review why this test cannot run with two chains? (CRE-983)
		// How to configure evm for both chains and capabilities DON (DON<>DON topology)?
		testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

		// TODO: remove this when OCR works properly with multiple chains in Local CRE
		testEnv.WrappedBlockchainOutputs = []*cre.WrappedBlockchainOutput{testEnv.WrappedBlockchainOutputs[0]}
		priceProvider, wfConfig := beforePoRTest(t, testEnv, "por-workflow-v2", PoRWFV2Location)
		wfConfig.FeedIDs = []string{wfConfig.FeedIDs[0]}
		ExecutePoRTest(t, testEnv, priceProvider, wfConfig, false)
	})

	t.Run("[v2] Vault DON - "+topology, func(t *testing.T) {
		testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

		ExecuteVaultTest(t, testEnv)
	})

	t.Run("[v2] Cron Beholder - "+topology, func(t *testing.T) {
		testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

		ExecuteCronBeholderTest(t, testEnv)
	})

	t.Run("[v2] HTTP Trigger Action - "+topology, func(t *testing.T) {
		testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

		ExecuteHTTPTriggerActionTest(t, testEnv)
	})

	t.Run("[v2] DON Time - "+topology, func(t *testing.T) {
		testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

		ExecuteDonTimeTest(t, testEnv)
	})

	t.Run("[v2] Consensus - "+topology, func(t *testing.T) {
		t.Skip("Quarantined - CRE-1064")
		testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

		ExecuteConsensusTest(t, testEnv)
	})
}

func Test_CRE_V2_EVM_Suite(t *testing.T) {
	topology := os.Getenv("TOPOLOGY_NAME")
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)
	// TODO: remove this when OCR works properly with multiple chains in Local CRE
	testEnv.WrappedBlockchainOutputs = []*cre.WrappedBlockchainOutput{testEnv.WrappedBlockchainOutputs[0]}

	t.Run("[v2] EVM Write - "+topology, func(t *testing.T) {
		priceProvider, porWfCfg := beforePoRTest(t, testEnv, "por-workflowV2", PoRWFV2Location)
		porWfCfg.FeedIDs = []string{porWfCfg.FeedIDs[0]}
		ExecutePoRTest(t, testEnv, priceProvider, porWfCfg, false)
	})

	t.Run("[v2] EVM Read - "+topology, func(t *testing.T) {
		ExecuteEVMReadTest(t, testEnv)
	})
}

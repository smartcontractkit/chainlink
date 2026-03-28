package cre

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	evm_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread/config"
	v2suite_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/v2suite/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

//////////// SMOKE TESTS /////////////
// target happy path and sanity checks
// all other tests (e.g. edge cases, negative conditions)
// should go to a `regression` package
/////////////////////////////////////

var v1RegistriesFlags = []string{"--with-contracts-version", "v1"}

var (
	parallelEnabled = t_helpers.ParallelEnabled()
	fanoutEnabled   = t_helpers.ChipSinkFanoutEnabled()
	// topology is used in test names
	topology = os.Getenv("TOPOLOGY_NAME")
)

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
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v1RegistriesFlags...)
	// WARNING: currently we can't run these tests in parallel, because each test rebuilds environment structs and that includes
	// logging into CL node with GraphQL API, which allows only 1 session per user at a time.

	// requires `readcontract`, `cron`
	priceProvider, porWfCfg := BeforePoRTest(t, testEnv, "por-workflowV1", PoRWFV1Location)
	ExecutePoRTest(t, testEnv, priceProvider, porWfCfg, false)
}

func Test_CRE_V1_Tron(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-don-tron.toml"), v1RegistriesFlags...)

	priceProvider, porWfCfg := BeforePoRTest(t, testEnv, "por-workflowV1", PoRWFV1Location)
	ExecutePoRTest(t, testEnv, priceProvider, porWfCfg, false)
}

/*
// TODO: Move Billing tests to v2 Registries
func Test_CRE_V1_Billing_EVM_Write(t *testing.T) {
	quarantine.Flaky(t, "DX-1911")
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))

	require.NoError(
		t,
		startBillingStackIfIsNotRunning(t, testEnv.TestConfig.RelativePathToRepoRoot, testEnv.TestConfig.EnvironmentDirPath, testEnv),
		"failed to start Billing stack",
	)

	priceProvider, porWfCfg := BeforePoRTest(t, testEnv, "por-workflowV2-billing", PoRWFV2Location)
	porWfCfg.FeedIDs = []string{porWfCfg.FeedIDs[0]}
	ExecutePoRTest(t, testEnv, priceProvider, porWfCfg, true)
}
*/

func Test_CRE_V1_Billing_Cron_Beholder(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v1RegistriesFlags...)

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
 1. Inside `core/scripts/cre/environment` directory: `go run . env restart --with-beholder`
 2. Execute the tests in `system-tests/tests/smoke/cre`: `go test -timeout 15m -run "^Test_CRE_V2"`.
*/
func Test_CRE_V2_Suite_Bucket_A(t *testing.T) {
	runV2SuiteBucket(t, v2suite_config.SuiteBucketA)
}

func Test_CRE_V2_Suite_Bucket_B(t *testing.T) {
	runV2SuiteBucket(t, v2suite_config.SuiteBucketB)
}

func Test_CRE_V2_Suite_Bucket_C(t *testing.T) {
	runV2SuiteBucket(t, v2suite_config.SuiteBucketC)
}

func runV2SuiteBucket(t *testing.T, bucket v2suite_config.SuiteBucket) {
	require.NoError(t, v2suite_config.ValidateSuiteBucketRegistry(), "invalid V2 suite bucket registry")

	scenarios, err := v2suite_config.ScenariosForSuiteBucket(bucket)
	require.NoErrorf(t, err, "failed to load V2 suite bucket %q", bucket)

	executeV2SuiteScenarios(t, topology, scenarios)
}

func executeV2SuiteScenarios(t *testing.T, topology string, scenarios []v2suite_config.SuiteScenario) {
	require.NotEmpty(t, scenarios, "no V2 suite scenarios selected")

	seen := make(map[v2suite_config.SuiteScenario]struct{}, len(scenarios))
	for _, scenario := range scenarios {
		require.GreaterOrEqualf(t, scenario, v2suite_config.SuiteScenario(0), "invalid scenario %d", scenario)
		require.Lessf(t, scenario, v2suite_config.SuiteScenarioLen, "invalid scenario %d", scenario)
		if _, alreadySeen := seen[scenario]; alreadySeen {
			require.Failf(t, "duplicate scenario", "scenario %q selected more than once", scenario.String())
		}
		seen[scenario] = struct{}{}
	}

	for _, scenario := range scenarios {
		runV2SuiteScenario(t, topology, scenario)
	}
}

func runV2SuiteScenario(t *testing.T, topology string, scenario v2suite_config.SuiteScenario) {
	switch scenario {
	case v2suite_config.SuiteScenarioProofOfReserve:
		t.Run("[v2] Proof Of Reserve - "+topology, func(t *testing.T) {
			if parallelEnabled {
				t.Parallel()
			}
			testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
			priceProvider, wfConfig := BeforePoRTest(t, testEnv, "por-workflow-v2", PoRWFV2Location)
			ExecutePoRTest(t, testEnv, priceProvider, wfConfig, false)
		})
	case v2suite_config.SuiteScenarioVaultDON:
		t.Run("[v2] Vault DON - "+topology, func(t *testing.T) {
			if parallelEnabled {
				t.Parallel()
			}
			testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
			ExecuteVaultTest(t, testEnv)
		})
	case v2suite_config.SuiteScenarioCronBeholder:
		// NOTE: this test is not easily parallelisable, because it uses "real" ChIP Ingress stack
		// we don't want to plug it into ChIP fanout, at least not yet
		t.Run("[v2] Cron Beholder - "+topology, func(t *testing.T) {
			testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))
			ExecuteCronBeholderTest(t, testEnv)
		})
	case v2suite_config.SuiteScenarioHTTPTriggerAction:
		t.Run("[v2] HTTP Trigger Action - "+topology, func(t *testing.T) {
			if parallelEnabled {
				t.Parallel()
			}
			testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
			ExecuteHTTPTriggerActionTest(t, testEnv)
		})
	case v2suite_config.SuiteScenarioHTTPActionCRUD:
		t.Run("[v2] HTTP Action CRUD - "+topology, func(t *testing.T) {
			if parallelEnabled && fanoutEnabled {
				t.Parallel()
			}
			testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
			ExecuteHTTPActionCRUDSuccessTest(t, testEnv)
		})
	case v2suite_config.SuiteScenarioDONTime:
		t.Run("[v2] DON Time - "+topology, func(t *testing.T) {
			if parallelEnabled && fanoutEnabled {
				t.Parallel()
			}
			testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
			ExecuteDonTimeTest(t, testEnv)
		})
	case v2suite_config.SuiteScenarioConsensus:
		t.Run("[v2] Consensus - "+topology, func(t *testing.T) {
			if parallelEnabled && fanoutEnabled {
				t.Parallel()
			}
			testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
			ExecuteConsensusTest(t, testEnv)
		})
	default:
		require.Failf(t, "unsupported V2 suite scenario", "scenario %q is not supported by the runner", scenario.String())
	}
}

func Test_CRE_V2_EVM_Write_LogTrigger(t *testing.T) {
	t.Run("[v2] EVM Write - "+topology, func(t *testing.T) {
		if parallelEnabled && fanoutEnabled {
			t.Parallel()
		}
		testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
		priceProvider, porWfCfg := BeforePoRTest(t, testEnv, "por-workflowV2", PoRWFV2Location)
		ExecutePoRTest(t, testEnv, priceProvider, porWfCfg, false)
	})

	t.Run("[v2] EVM LogTrigger - "+topology, func(t *testing.T) {
		if parallelEnabled && fanoutEnabled {
			t.Parallel()
		}
		testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
		ExecuteEVMLogTriggerTest(t, testEnv)
	})
}

func Test_CRE_V2_EVM_Read_HeavyCalls(t *testing.T) {
	runV2EVMReadBucket(t, evm_config.ReadBucketHeavyCalls)
}

func Test_CRE_V2_EVM_Read_StateQueries(t *testing.T) {
	runV2EVMReadBucket(t, evm_config.ReadBucketStateQueries)
}

func Test_CRE_V2_EVM_Read_TxArtifacts(t *testing.T) {
	runV2EVMReadBucket(t, evm_config.ReadBucketTxArtifacts)
}

func runV2EVMReadBucket(t *testing.T, bucket evm_config.ReadBucket) {
	testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
	require.NoError(t, evm_config.ValidateReadBucketRegistry(), "invalid EVM read bucket registry")

	testCases, err := evm_config.CasesForReadBucket(bucket)
	require.NoErrorf(t, err, "failed to load EVM read bucket %q", bucket)

	t.Run(fmt.Sprintf("[v2] EVM Read (%s) - %s", bucket, topology), func(t *testing.T) {
		ExecuteEVMReadTestForCases(t, testEnv, testCases)
	})
}

func Test_CRE_V2_Solana_Suite(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-don-solana.toml"))
	t.Run("[v2] Solana Write", func(t *testing.T) {
		ExecuteSolanaWriteTest(t, testEnv)
	})
}

func Test_CRE_V2_Aptos_Suite(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-gateway-don-aptos.toml"))
	t.Run("[v2] Aptos", func(t *testing.T) {
		ExecuteAptosTest(t, testEnv)
	})
}
func Test_CRE_V2_HTTP_Action_Regression_Suite(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))

	ExecuteHTTPActionRegressionTest(t, testEnv)
}

func Test_CRE_V2_Beholder_Suite(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), "--with-dashboards")

	ExecuteLogStreamingTest(t, testEnv)
}

func Test_CRE_V2_Sharding(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(
		t,
		t_helpers.GetTestConfig(t, "/configs/workflow-gateway-sharded-don.toml"),
	)

	ExecuteShardingTest(t, testEnv)
}

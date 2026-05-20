package cre

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	suite_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/config"
	evm_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

//////////// SMOKE TESTS /////////////
// target happy path and sanity checks
// all other tests (e.g. edge cases, negative conditions)
// should go to a `regression` package
/////////////////////////////////////

var (
	parallelEnabled = t_helpers.ParallelEnabled()
	// topology is used in test names
	topology = os.Getenv("TOPOLOGY_NAME")
)

//////////// CRE TESTS /////////////
/*
To execute tests start the local CRE first:
 1. Inside `core/scripts/cre/environment` directory: `go run . env restart --with-beholder`
 2. Execute the tests in `system-tests/tests/smoke/cre`: `go test -timeout 15m -run "^Test_CRE_"`.
*/
func Test_CRE_V2_Suite_Bucket_A(t *testing.T) {
	runSuiteBucket(t, suite_config.SuiteBucketA)
}

func Test_CRE_V2_Suite_Bucket_B(t *testing.T) {
	runSuiteBucket(t, suite_config.SuiteBucketB)
}

func Test_CRE_V2_Suite_Bucket_C(t *testing.T) {
	runSuiteBucket(t, suite_config.SuiteBucketC)
}

func runSuiteBucket(t *testing.T, bucket suite_config.SuiteBucket) {
	require.NoError(t, suite_config.ValidateSuiteBucketRegistry(), "invalid suite bucket registry")

	scenarios, err := suite_config.ScenariosForSuiteBucket(bucket)
	require.NoErrorf(t, err, "failed to load suite bucket %q", bucket)

	executeSuiteScenarios(t, topology, scenarios)
}

func executeSuiteScenarios(t *testing.T, topology string, scenarios []suite_config.SuiteScenario) {
	require.NotEmpty(t, scenarios, "no suite scenarios selected")

	seen := make(map[suite_config.SuiteScenario]struct{}, len(scenarios))
	for _, scenario := range scenarios {
		require.GreaterOrEqualf(t, scenario, suite_config.SuiteScenario(0), "invalid scenario %d", scenario)
		require.Lessf(t, scenario, suite_config.SuiteScenarioLen, "invalid scenario %d", scenario)
		if _, alreadySeen := seen[scenario]; alreadySeen {
			require.Failf(t, "duplicate scenario", "scenario %q selected more than once", scenario.String())
		}
		seen[scenario] = struct{}{}
	}

	for _, scenario := range scenarios {
		runSuiteScenario(t, topology, scenario)
	}
}

func runSuiteScenario(t *testing.T, topology string, scenario suite_config.SuiteScenario) {
	switch scenario {
	case suite_config.SuiteScenarioProofOfReserve:
		t.Run("Proof Of Reserve - "+topology, func(t *testing.T) {
			if parallelEnabled {
				t.Parallel()
			}
			testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
			priceProvider, wfConfig := BeforePoRTest(t, testEnv, "por-workflow-v2", PoRWFLocation)
			ExecutePoRTest(t, testEnv, priceProvider, wfConfig, false)
		})
	case suite_config.SuiteScenarioVaultDON:
		t.Run("Vault DON - "+topology, func(t *testing.T) {
			if parallelEnabled {
				t.Parallel()
			}
			allowlistSubtestName := "allowlist_auth_when_jwt_auth_disabled"
			jwtSubtestName := "jwt_auth_rejected_when_jwt_auth_disabled"
			vaultConfig := getVaultDefaultTestConfig(t)
			if isVaultJWTAuthEnabledTopology(topology) {
				vaultConfig = getVaultJWTAuthEnabledTestConfig(t)
				allowlistSubtestName = "allowlist_auth_when_jwt_auth_enabled"
				jwtSubtestName = "jwt_auth_when_jwt_auth_enabled"
			}
			fixture := setupVaultSharedScenarioFixture(t, vaultConfig)
			allowlistEnv := fixture.TestEnv
			jwtEnv := fixture.TestEnv
			if parallelEnabled && isVaultJWTAuthEnabledTopology(topology) {
				allowlistEnv = t_helpers.SetupTestEnvironmentWithPerTestKeys(t, fixture.TestEnv.TestConfig)
				jwtEnv = t_helpers.SetupTestEnvironmentWithPerTestKeys(t, fixture.TestEnv.TestConfig)
			}

			t.Run(allowlistSubtestName, func(t *testing.T) {
				if parallelEnabled {
					t.Parallel()
				}
				ExecuteVaultAllowListBasedTests(t, fixture, allowlistEnv)
			})
			if isVaultJWTAuthEnabledTopology(topology) {
				t.Run(jwtSubtestName, func(t *testing.T) {
					if parallelEnabled {
						t.Parallel()
					}
					ExecuteVaultMixedAuthTest(t, fixture, jwtEnv)
				})
				return
			}
			t.Run(jwtSubtestName, func(t *testing.T) {
				if parallelEnabled {
					t.Parallel()
				}
				ExecuteVaultJWTDisabledTest(t, fixture)
			})
		})
	case suite_config.SuiteScenarioCronBeholder:
		t.Run("Cron Beholder - "+topology, func(t *testing.T) {
			if parallelEnabled {
				t.Parallel()
			}
			testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))
			ExecuteCronBeholderTest(t, testEnv)
		})
	case suite_config.SuiteScenarioHTTPTriggerAction:
		t.Run("HTTP Trigger Action - "+topology, func(t *testing.T) {
			if parallelEnabled {
				t.Parallel()
			}
			testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
			ExecuteHTTPTriggerActionTest(t, testEnv)
		})
	case suite_config.SuiteScenarioHTTPActionCRUD:
		t.Run("HTTP Action CRUD - "+topology, func(t *testing.T) {
			if parallelEnabled {
				t.Parallel()
			}
			testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
			ExecuteHTTPActionCRUDSuccessTest(t, testEnv)
		})
	case suite_config.SuiteScenarioDONTime:
		t.Run("DON Time - "+topology, func(t *testing.T) {
			if parallelEnabled {
				t.Parallel()
			}
			testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
			ExecuteDonTimeTest(t, testEnv)
		})
	case suite_config.SuiteScenarioConsensus:
		t.Run("Consensus - "+topology, func(t *testing.T) {
			if parallelEnabled {
				t.Parallel()
			}
			testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
			ExecuteConsensusTest(t, testEnv)
		})
	default:
		require.Failf(t, "unsupported suite scenario", "scenario %q is not supported by the runner", scenario.String())
	}
}

func Test_CRE_V2_EVM_Write_LogTrigger(t *testing.T) {
	t.Run("EVM Write - "+topology, func(t *testing.T) {
		if parallelEnabled {
			t.Parallel()
		}
		testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
		priceProvider, porWfCfg := BeforePoRTest(t, testEnv, "por-workflow", PoRWFLocation)
		ExecutePoRTest(t, testEnv, priceProvider, porWfCfg, false)
	})

	t.Run("EVM LogTrigger - "+topology, func(t *testing.T) {
		if parallelEnabled {
			t.Parallel()
		}
		testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
		ExecuteEVMLogTriggerTest(t, testEnv)
	})
}

func Test_CRE_V2_EVM_Read_HeavyCalls(t *testing.T) {
	runEVMReadBucket(t, evm_config.ReadBucketHeavyCalls)
}

func Test_CRE_V2_EVM_Read_StateQueries(t *testing.T) {
	runEVMReadBucket(t, evm_config.ReadBucketStateQueries)
}

func Test_CRE_V2_EVM_Read_TxArtifacts(t *testing.T) {
	runEVMReadBucket(t, evm_config.ReadBucketTxArtifacts)
}

func runEVMReadBucket(t *testing.T, bucket evm_config.ReadBucket) {
	testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
	require.NoError(t, evm_config.ValidateReadBucketRegistry(), "invalid EVM read bucket registry")

	testCases, err := evm_config.CasesForReadBucket(bucket)
	require.NoErrorf(t, err, "failed to load EVM read bucket %q", bucket)

	t.Run(fmt.Sprintf("EVM Read (%s) - %s", bucket, topology), func(t *testing.T) {
		ExecuteEVMReadTestForCases(t, testEnv, testCases)
	})
}

func Test_CRE_V2_Solana_Suite(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-don-solana.toml"))
	t.Run("Solana Write", func(t *testing.T) {
		ExecuteSolanaWriteTest(t, testEnv)
	})
}

func Test_CRE_V2_Aptos_Suite(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-gateway-don-aptos.toml"))
	t.Run("Aptos", func(t *testing.T) {
		ExecuteAptosTest(t, testEnv)
	})
}

func Test_CRE_V2_Module_Cache(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-gateway-don-cache-test.toml"))

	ExecuteModuleCacheTest(t, testEnv)
}

func Test_CRE_V2_HTTP_Action_Regression_Suite(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))

	ExecuteHTTPActionRegressionTest(t, testEnv)
}

func Test_CRE_V2_Beholder_Suite(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), "--with-dashboards")

	ExecuteLogStreamingTest(t, testEnv)
}

func Test_CRE_V2_DurableEmitter(t *testing.T) {
	t.Skip("CRE-4315 fix CRE_V2_DurableEmitter test on CI")
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))
	ExecuteDurableEmitterTest(t, testEnv)
}

func Test_CRE_V2_Sharding(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(
		t,
		t_helpers.GetTestConfig(t, "/configs/workflow-gateway-sharded-don.toml"),
	)

	ExecuteShardingTest(t, testEnv)
}

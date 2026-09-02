package cre

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	suite_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/config"
	evm_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread/config"
	solana_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/solread/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
 1. Inside `core/scripts/cre/environment` directory: `go run . env restart --with-chip-ingress-stack` (deprecated: `--with-beholder`)
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
			allowlistSubtestName := "allowlist_auth"
			vaultConfig := getVaultDefaultTestConfig(t)
			if isVaultStallPurgeTopology(topology) {
				vaultConfig = getVaultStallPurgeTestConfig(t)
				allowlistSubtestName = "pending_queue_stall_purge"
			} else if isVaultWorkflowDONBindingEnabledTopology(topology) {
				vaultConfig = getVaultWorkflowDONBindingEnabledTestConfig(t)
				allowlistSubtestName = "allowlist_auth_when_workflow_don_binding_enabled"
			}
			fixture := setupVaultSharedScenarioFixture(t, vaultConfig)

			t.Run(allowlistSubtestName, func(t *testing.T) {
				if parallelEnabled {
					t.Parallel()
				}
				allowlistEnv := fixture.TestEnv
				if parallelEnabled {
					allowlistEnv = t_helpers.SetupTestEnvironmentWithPerTestKeys(t, fixture.TestEnv.TestConfig)
				}
				if isVaultStallPurgeTopology(topology) {
					ExecuteVaultPendingQueueStallPurgeSmokeTest(t, fixture, allowlistEnv)
					return
				}
				ExecuteVaultAllowListBasedTests(t, fixture, allowlistEnv)
			})
			if isVaultStallPurgeTopology(topology) {
				return
			}
			t.Run("jwt_auth", func(t *testing.T) {
				if parallelEnabled {
					t.Parallel()
				}
				jwtEnv := fixture.TestEnv
				if parallelEnabled {
					jwtEnv = t_helpers.SetupTestEnvironmentWithPerTestKeys(t, fixture.TestEnv.TestConfig)
				}
				ExecuteVaultMixedAuthTest(t, fixture, jwtEnv)
			})
		})
	case suite_config.SuiteScenarioCronChipIngressStack:
		t.Run("Cron Beholder - "+topology, func(t *testing.T) {
			if parallelEnabled {
				t.Parallel()
			}
			testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))
			ExecuteCronChipIngressStackTest(t, testEnv)
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
	case suite_config.SuiteScenarioHTTPActionMultiGateway:
		t.Run("HTTP Action Multi Gateway - "+topology, func(t *testing.T) {
			if !isMultiGatewayTopology(topology) {
				t.Skipf("skipping multi-gateway HTTP action test on topology %q", topology)
			}
			if parallelEnabled {
				t.Parallel()
			}
			testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, getMultiGatewayTestConfig(t))
			ExecuteHTTPActionMultiGatewayRoutingTest(t, testEnv)
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

const solanaConfigPath = "/configs/workflow-don-solana.toml"

//nolint:paralleltest // isolate local cre env run
func Test_CRE_V2_Solana_Write(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, solanaConfigPath))
	t.Run("Solana Write", func(t *testing.T) {
		ExecuteSolanaWriteTest(t, testEnv)
	})
}

func Test_CRE_V2_Solana_LogTrigger(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, solanaConfigPath))
	t.Run("Solana LogTrigger", func(t *testing.T) {
		ExecuteSolanaLogTriggerTest(t, testEnv)
	})
	t.Run("Solana LogTrigger CPI", func(t *testing.T) {
		ExecuteSolanaLogTriggerCPITest(t, testEnv)
	})
}

//nolint:paralleltest // single test
func Test_CRE_V2_Solana_Read_Accounts(t *testing.T) {
	runSolanaReadBucket(t, solana_config.ReadBucketAccountCalls)
}

//nolint:paralleltest // single test
func Test_CRE_V2_Solana_Read_Block(t *testing.T) {
	runSolanaReadBucket(t, solana_config.ReadBucketBlockCalls)
}

//nolint:paralleltest // single test
func Test_CRE_V2_Solana_Read_Tx(t *testing.T) {
	runSolanaReadBucket(t, solana_config.ReadBucketTxCalls)
}

func runSolanaReadBucket(t *testing.T, bucket solana_config.ReadBucket) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, solanaConfigPath))
	require.NoError(t, solana_config.ValidateReadBucketRegistry(), "invalid Solana read bucket registry")

	testCases, err := solana_config.CasesForReadBucket(bucket)
	require.NoErrorf(t, err, "failed to load Solana read bucket %q", bucket)

	t.Run(fmt.Sprintf("Solana Read (%s) - %s", bucket, topology), func(t *testing.T) {
		ExecuteSolanaReadTestForCases(t, testEnv, testCases)
	})
}

func Test_CRE_V2_Aptos_Suite(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-gateway-don-aptos.toml"))
	t.Run("Aptos", func(t *testing.T) {
		ExecuteAptosTest(t, testEnv)
	})
}

//nolint:paralleltest // isolate local cre env run
func Test_CRE_V2_Stellar_Suite(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-gateway-don-stellar.toml"))

	t.Run("Stellar GetLatestLedger", func(t *testing.T) {
		t.Parallel()
		env, chain, userLogsCh, baseMessageCh := setupStellarScenario(t, testEnv)
		executeStellarReadLatestLedgerTest(t, env, chain, userLogsCh, baseMessageCh)
	})

	t.Run("Stellar ReadContract", func(t *testing.T) {
		t.Parallel()
		env, chain, userLogsCh, baseMessageCh := setupStellarScenario(t, testEnv)
		executeStellarReadContractSmokeTest(t, env, chain, userLogsCh, baseMessageCh)
	})

	t.Run("StellarWrite", func(t *testing.T) {
		t.Parallel()
		env, chain, userLogsCh, baseMessageCh := setupStellarScenario(t, testEnv)
		executeStellarWriteTest(t, env, chain, userLogsCh, baseMessageCh)
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
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))
	ExecuteDurableEmitterTest(t, testEnv)
}

//nolint:paralleltest // subtests share the same sharding config
func Test_CRE_V2_Sharding(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(
		t,
		t_helpers.GetTestConfig(t, "/configs/workflow-gateway-sharded-don.toml"),
	)
	t.Run("ExecuteShardingTestWithCronTrigger", func(t *testing.T) {
		ExecuteShardingTestWithCronTrigger(t, testEnv)
	})
	t.Run("ExecuteShardingTestWithEVMLogTrigger", func(t *testing.T) {
		// Reinitialize OperationsBundle so that it can reexecute shard config updates instead of caching them.
		testEnv.CreEnvironment.CldfEnvironment.OperationsBundle = operations.NewBundle(t.Context, logger.TestLogger(t), operations.NewMemoryReporter())
		ExecuteShardingTestWithEVMLogTrigger(t, testEnv)
	})
}

//nolint:paralleltest // subtests share the same sharding config
func Test_CRE_V2_ShardingWithHttpTrigger(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(
		t,
		t_helpers.GetTestConfig(t, "/configs/workflow-gateway-sharded-don.toml"),
	)
	t.Run("ExecuteShardingTestWithHTTPTrigger", func(t *testing.T) {
		ExecuteShardingTestWithHTTPTrigger(t, testEnv)
	})
}

//nolint:paralleltest // subtests share the same sharding config
func Test_CRE_V2_ShardManualAssignment(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(
		t,
		t_helpers.GetTestConfig(t, "/configs/workflow-gateway-sharded-manual.toml"),
	)
	ExecuteManualShardAssignmentTest(t, testEnv)
}

//nolint:paralleltest // subtests share the same sharding config
func Test_CRE_V2_ShardRingOCROverrides(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(
		t,
		t_helpers.GetTestConfig(t, "/configs/workflow-gateway-sharded-ringocr-overrides.toml"),
	)
	ExecuteRingOCROverridesTest(t, testEnv)
}

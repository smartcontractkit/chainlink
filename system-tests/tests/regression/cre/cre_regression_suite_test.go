package cre

import (
	"fmt"
	"strings"
	"testing"

	"github.com/smartcontractkit/quarantine"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// REGRESSION TESTS target edge cases, negative conditions, etc., all happy path and sanity checks should go to a `smoke` package.

var v2RegistriesFlags = []string{"--with-contracts-version", "v2"}

/*
To execute tests locally start the local CRE first:
Inside `core/scripts/cre/environment` directory
 1. Ensure the necessary capabilities (i.e. readcontract, http-trigger, http-action) are listed in the environment configuration
 2. Identify the appropriate topology that you want to test
 3. Stop and clear any existing environment: `go run . env stop -a`
 4. Run: `CTF_CONFIGS=<path-to-your-topology-config> go run . env start && ./bin/ctf obs up` to start env + observability
 5. Optionally run the Blockscout (chain explorer) `./bin/ctf bs up`
 6. Execute the tests in `system-tests/tests/regression/cre`: `go test -timeout 15m -run "^Test_CRE_V2"`
*/
func Test_CRE_V2_Consensus_Regression(t *testing.T) {
	// a template for Consensus negative tests names to avoid duplication
	const consensusTestNameTemplate = "[v2] Consensus.%s fails with %s" // e.g. "[v2] Consensus.<Function> fails with <invalid input>"

	for _, tCase := range consensusNegativeTestsGenerateReport {
		testName := fmt.Sprintf(consensusTestNameTemplate, tCase.caseToTrigger, tCase.name)
		t.Run(testName, func(t *testing.T) {
			testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)
			ConsensusFailsTest(t, testEnv, tCase)
		})
	}
}

func Test_CRE_V2_Cron_Regression(t *testing.T) {
	for _, tCase := range cronInvalidSchedulesTests {
		testName := "[v2] Cron (Beholder) fails when schedule is " + tCase.name
		t.Run(testName, func(t *testing.T) {
			testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

			CronBeholderFailsWithInvalidScheduleTest(t, testEnv, tCase.invalidSchedule)
		})
	}
}

func Test_CRE_V2_HTTP_Regression(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

	for _, tCase := range httpNegativeTests {
		testName := "[v2] HTTP Trigger fails with " + tCase.name
		t.Run(testName, func(t *testing.T) {
			HTTPTriggerFailsTest(t, testEnv, tCase)
		})
	}
}

// runEVMNegativeTestSuite runs a suite of EVM negative tests with the given test cases
func runEVMNegativeTestSuite(t *testing.T, testCases []evmNegativeTest) {
	// a template for EVM negative tests names to avoid duplication
	const evmTestNameTemplate = "[v2] EVM.%s fails with %s" // e.g. "[v2] EVM.<Function> fails with <invalid input>"

	for _, tCase := range testCases {
		testName := fmt.Sprintf(evmTestNameTemplate, tCase.functionToTest, tCase.name)
		t.Run(testName, func(t *testing.T) {
			testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)
			// TODO remove this when OCR works properly with multiple chains in Local CRE
			testEnv.WrappedBlockchainOutputs = []*cre.WrappedBlockchainOutput{testEnv.WrappedBlockchainOutputs[0]}

			// Check if test name contains "write" to determine which test function to run
			if strings.Contains(strings.ToLower(testName), "writereport") {
				framework.L.Info().Msg("Running EVM Write Regression test")
				EVMWriteFailsTest(t, testEnv, tCase)
			} else {
				framework.L.Info().Msg("Running EVM Read Regression test")
				EVMReadFailsTest(t, testEnv, tCase)
			}
		})
	}
}

func Test_CRE_V2_EVM_BalanceAt_Invalid_Address_Regression(t *testing.T) {
	runEVMNegativeTestSuite(t, evmNegativeTestsBalanceAtInvalidAddress)
}

func Test_CRE_V2_EVM_CallContract_Invalid_Addr_To_Read_Regression(t *testing.T) {
	runEVMNegativeTestSuite(t, evmNegativeTestsCallContractInvalidAddressToRead)
}

func Test_CRE_V2_EVM_CallContract_Invalid_Balance_Reader_Contract_Regression(t *testing.T) {
	runEVMNegativeTestSuite(t, evmNegativeTestsCallContractInvalidBalanceReaderContract)
}

func Test_CRE_V2_EVM_EstimateGas_Invalid_To_Address_Regression(t *testing.T) {
	runEVMNegativeTestSuite(t, evmNegativeTestsEstimateGasInvalidToAddress)
}

func Test_CRE_V2_EVM_FilterLogs_Invalid_Addresses_Regression(t *testing.T) {
	runEVMNegativeTestSuite(t, evmNegativeTestsFilterLogsWithInvalidAddress)
}

func Test_CRE_V2_EVM_FilterLogs_Invalid_FromBlock_Regression(t *testing.T) {
	runEVMNegativeTestSuite(t, evmNegativeTestsFilterLogsWithInvalidFromBlock)
}

func Test_CRE_V2_EVM_FilterLogs_Invalid_ToBlock_Regression(t *testing.T) {
	runEVMNegativeTestSuite(t, evmNegativeTestsFilterLogsWithInvalidToBlock)
}

func Test_CRE_V2_EVM_GetTransactionByHash_Invalid_Hash_Regression(t *testing.T) {
	runEVMNegativeTestSuite(t, evmNegativeTestsGetTransactionByHashInvalidHash)
}

func Test_CRE_V2_EVM_GetTransactionReceipt_Invalid_Hash_Regression(t *testing.T) {
	runEVMNegativeTestSuite(t, evmNegativeTestsGetTransactionReceiptInvalidHash)
}

func Test_CRE_V2_EVM_HeaderByNumber_Invalid_Block_Regression(t *testing.T) {
	runEVMNegativeTestSuite(t, evmNegativeTestsHeaderByNumberInvalidBlock)
}

func Test_CRE_V2_EVM_WriteReport_Invalid_Receiver_Regression(t *testing.T) {
	runEVMNegativeTestSuite(t, evmNegativeTestsWriteReportInvalidReceiver)
}

func Test_CRE_V2_EVM_WriteReport_Corrupt_Receiver_Address_Regression(t *testing.T) {
	quarantine.Flaky(t, "DX-2049")
	runEVMNegativeTestSuite(t, evmNegativeTestsWriteReportCorruptReceiverAddress)
}

func Test_CRE_V2_EVM_WriteReport_Invalid_Gas_Regression(t *testing.T) {
	runEVMNegativeTestSuite(t, evmNegativeTestsWriteReportInvalidGas)
}

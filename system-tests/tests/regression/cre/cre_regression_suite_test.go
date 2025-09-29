package cre

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
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
 5. Optionally run blockscout `./bin/ctf bs up`
 6. Execute the tests in `system-tests/tests/smoke/cre` with CTF_CONFIG set to the corresponding topology file:
    `export  CTF_CONFIGS=../../../../core/scripts/cre/environment/configs/<topology>.toml; go test -timeout 15m -run ^Test_CRE_Suite$`.
*/
func Test_CRE_V2_Cron_Regression(t *testing.T) {
	for _, tCase := range cronInvalidSchedulesTests {
		testName := "[v2] Cron (Beholder) fails when schedule is " + tCase.name
		t.Run(testName, func(t *testing.T) {
			testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), v2RegistriesFlags...)

			CronBeholderFailsWithInvalidScheduleTest(t, testEnv, tCase.invalidSchedule)
		})
	}
}

func Test_CRE_Suite_V2_HTTP_Regression(t *testing.T) {
	flags := []string{"--with-contracts-version", "v2"}
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t), flags...)

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

			EVMReadFailsTest(t, testEnv, tCase)
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
	t.Skip("Skipping because evm.FilterLogs does not return an error or nil for invalid/not existing addresses")
	runEVMNegativeTestSuite(t, evmNegativeTestsFilterLogsWithInvalidAddress)
}

func Test_CRE_V2_EVM_FilterLogs_Invalid_FromBlock_Regression(t *testing.T) {
	runEVMNegativeTestSuite(t, evmNegativeTestsFilterLogsWithInvalidFromBlock)
}

func Test_CRE_V2_EVM_FilterLogs_Invalid_ToBlock_Regression(t *testing.T) {
	runEVMNegativeTestSuite(t, evmNegativeTestsFilterLogsWithInvalidToBlock)
}

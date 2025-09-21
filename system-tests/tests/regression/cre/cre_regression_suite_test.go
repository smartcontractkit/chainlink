package cre

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	tcommon "github.com/smartcontractkit/chainlink/system-tests/tests/common"
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
func Test_CRE_Suite_Regression(t *testing.T) {
	testEnv := tcommon.SetupTestEnvironmentWithConfig(t, tcommon.GetDefaultTestConfig(t))
	t.Run("[v2] CRE Regression Suite", func(t *testing.T) {
		// negative tests for cron
		// TODO: move to a separate package
		for _, tCase := range cronInvalidSchedulesTests {
			testName := fmt.Sprintf("[v2] Cron (Beholder) fails when schedule is %s", tCase.name)
			t.Run(testName, func(t *testing.T) {
				CronBeholderFailsWithInvalidScheduleTest(t, testEnv, tCase.invalidSchedule)
			})
		}
	})
}

func Test_CRE_Suite_EVM(t *testing.T) {
	testEnv := tcommon.SetupTestEnvironmentWithConfig(t, tcommon.GetDefaultTestConfig(t))
	// TODO remove this when OCR works properly with multiple chains in Local CRE
	testEnv.WrappedBlockchainOutputs = []*cre.WrappedBlockchainOutput{testEnv.WrappedBlockchainOutputs[0]}

	// negative tests for evm read
	// TODO: move to a separate package
	for _, tCase := range evmNegativeTests {
		testName := fmt.Sprintf("[v2] EVM.%s fails with %s", tCase.functionToTest, tCase.name)
		t.Run(testName, func(t *testing.T) {
			EVMReadFailsTest(t, testEnv, tCase)
		})
	}
}

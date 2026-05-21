package cre

import (
	"testing"

	suite_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/config"
)

/*
	Add upgrade tests below. These tests are run during the release process to verify that the upgrade procedure works as expected.
*/

func Test_Upgrade_Suite(t *testing.T) {
	executeSuiteScenarios(t, "workflow-gateway-don", []suite_config.SuiteScenario{suite_config.SuiteScenarioProofOfReserve})
}

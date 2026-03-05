package cre

import (
	"testing"

	v2suite_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/v2suite/config"
)

/*
	Add upgrade tests below. These tests are run during the release process to verify that the upgrade procedure works as expected.
*/

func Test_Upgrade_Suite(t *testing.T) {
	executeV2SuiteScenarios(t, "workflow-gateway-don", []v2suite_config.SuiteScenario{v2suite_config.SuiteScenarioProofOfReserve})
}

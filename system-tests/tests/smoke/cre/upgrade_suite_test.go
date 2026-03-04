package cre

import (
	"testing"

	v2suite_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/v2suite/config"
)

/*
Add tests below which will be run, when run upgrade tests during release process to ensure that the upgrade process is working as expected.
*/

func Test_Upgrade_Suite(t *testing.T) {
	executeV2SuiteScenarios(t, "workflow-gateway-don", []v2suite_config.SuiteScenario{v2suite_config.SuiteScenarioProofOfReserve})
}

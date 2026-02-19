//go:build upgrade
// +build upgrade

package cre

import (
	"testing"

	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

func Test_CRE_V2_Upgrade_DONTime(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))
	ExecuteDonTimeTest(t, testEnv)
}

func Test_CRE_V2_Upgrade_Consensus(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))
	ExecuteConsensusTest(t, testEnv)
}

func Test_CRE_V2_Upgrade_HTTPTriggerAction(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetDefaultTestConfig(t))
	ExecuteHTTPTriggerActionTest(t, testEnv)
}

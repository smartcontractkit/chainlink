package smoke

import (
	"testing"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestAutomationNodeUpgrade(t *testing.T) {
	config, err := tc.GetConfig(tc.Smoke, tc.Automation)
	if err != nil {
		t.Fatal(err)
	}
	SetupAutomationBasic(t, true, &config)
}

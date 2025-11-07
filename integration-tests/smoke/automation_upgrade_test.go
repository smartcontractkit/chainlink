package smoke

import (
	"testing"

	"github.com/smartcontractkit/quarantine"
)

func TestAutomationNodeUpgrade(t *testing.T) {
	quarantine.Flaky(t, "DX-2060")
	SetupAutomationBasic(t, true)
}

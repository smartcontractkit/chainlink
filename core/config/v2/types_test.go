package v2

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/chainlink/cfgtest"
)

func TestCoreDefaults_notNil(t *testing.T) {
	cfgtest.AssertFieldsNotNil(t, &defaults)
}

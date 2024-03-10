package docs

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink/cfgtest"
)

func TestCoreDefaults_notNil(t *testing.T) {
	cfgtest.AssertFieldsNotNil(t, CoreDefaults())
}

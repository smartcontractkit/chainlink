package docs

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/config/configtest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestCoreDefaults_notNil(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	configtest.AssertFieldsNotNil(t, CoreDefaults())
}

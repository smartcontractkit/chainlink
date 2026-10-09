package docs

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/config/configtest"
)

func TestCoreDefaults_notNil(t *testing.T) {
	t.Parallel()
	configtest.AssertFieldsNotNil(t, CoreDefaults())
}

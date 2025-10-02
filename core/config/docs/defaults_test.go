package docs

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/config/configtest"
)

func TestCoreDefaults_notNil(t *testing.T) {
	configtest.AssertFieldsNotNil(t, CoreDefaults())
}

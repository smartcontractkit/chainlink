package docs

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/config/configtest"
)

func TestUnit_CoreDefaults_notNil(t *testing.T) {
	configtest.AssertFieldsNotNil(t, CoreDefaults())
}

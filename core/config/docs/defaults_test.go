package docs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/config/configtest"
)

func TestCoreDefaults_notNil(t *testing.T) {
	configtest.AssertFieldsNotNil(t, CoreDefaults())
}

func TestCoreDefaults_distinctInstances(t *testing.T) {
	t.Parallel()
	a := CoreDefaults()
	b := CoreDefaults()
	*a.Insecure.DevWebServer = true
	*a.AuditLogger.Enabled = true
	assert.False(t, *b.Insecure.DevWebServer)
	assert.False(t, *b.AuditLogger.Enabled)
}

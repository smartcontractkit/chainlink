package capabilities

import (
	"testing"

	"github.com/stretchr/testify/assert"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func TestExecutableClientOpts(t *testing.T) {
	t.Parallel()
	assert.Len(t, executableClientOpts(vaultcommon.CapabilityID, vaulttypes.MethodSecretsGet), 1)
	assert.Empty(t, executableClientOpts(vaultcommon.CapabilityID, "other"))
	assert.Empty(t, executableClientOpts("other@1.0.0", vaulttypes.MethodSecretsGet))
}

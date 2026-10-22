package capabilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func TestExecutableClientOpts(t *testing.T) {
	t.Parallel()
	shareSizeLimit, err := limits.MakeUpperBoundLimiter(limits.Factory{Settings: cresettings.DefaultGetter}, cresettings.Default.VaultShareSizeLimit)
	require.NoError(t, err)
	w := &launcher{vaultShareSizeLimit: shareSizeLimit}

	assert.NotNil(t, w.executableClientOpts(vaultcommon.CapabilityID, vaulttypes.MethodSecretsGet))
	assert.Empty(t, w.executableClientOpts(vaultcommon.CapabilityID, "other"))
	assert.Empty(t, w.executableClientOpts("other@1.0.0", vaulttypes.MethodSecretsGet))
}

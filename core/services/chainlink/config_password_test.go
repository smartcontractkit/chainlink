package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestPasswordConfig(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	p := cfg.Password()
	assert.Empty(t, p.VRF())
	assert.Empty(t, p.Keystore())
}

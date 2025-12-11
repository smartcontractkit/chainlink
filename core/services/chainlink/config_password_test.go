package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_PasswordConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	p := cfg.Password()
	assert.Empty(t, p.VRF())
	assert.Empty(t, p.Keystore())
}

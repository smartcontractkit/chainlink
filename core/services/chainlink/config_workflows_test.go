package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowsConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	w := cfg.Workflows()
	assert.Equal(t, int32(200), w.Limits().Global())
	assert.Equal(t, int32(200), w.Limits().PerOwner())
	assert.Equal(t,
		map[string]int32{
			"0xae4E781a6218A8031764928E88d457937A954fC3": 5,
			"0x538aAaB4ea120b2bC2fe5D296852D948F07D849e": 10,
		},
		w.Limits().PerOwnerOverrides())
}

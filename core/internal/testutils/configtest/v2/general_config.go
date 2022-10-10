package v2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

// TODO doc
func NewTestGeneralConfig(t *testing.T) config.GeneralConfig { return TOML{}.NewGeneralConfig(t) }

// TODO doc
type TOML struct {
	Config  string
	Secrets string
}

// TODO doc
func (c TOML) NewGeneralConfig(t *testing.T) config.GeneralConfig {
	g, err := chainlink.NewTOMLGeneralConfig(logger.TestLogger(t), c.Config, c.Secrets, nil, nil)
	require.NoError(t, err)
	return g
}

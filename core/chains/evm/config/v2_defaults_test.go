package config_test

import (
	"fmt"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/pelletier/go-toml/v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/config"
	v2 "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_v2Config_SetDefaults(t *testing.T) {
	var fallbackTOML []byte
	t.Run("fallback", func(t *testing.T) {
		got, name := v2.DefaultsNamed(nil)
		assert.Empty(t, name)
		exp := config.FallbackDefaultsAsV2()

		assertChainsEqual(t, exp, got)

		var err error
		fallbackTOML, err = toml.Marshal(got)
		require.NoError(t, err)
	})
	for id, exp := range config.ChainSpecificConfigDefaultsAsV2() {
		got, name := v2.DefaultsNamed(utils.NewBigI(id))
		t.Run(fmt.Sprintf("%d:%s", id, name), func(t *testing.T) {
			assertChainsEqual(t, exp, got)
		})
	}
	t.Run("fallback-unchanged", func(t *testing.T) {
		got, _ := v2.DefaultsNamed(nil)
		gotTOML, err := toml.Marshal(got)
		require.NoError(t, err)
		assert.Equal(t, fallbackTOML, gotTOML)
	})
}

func assertChainsEqual(t *testing.T, exp, got v2.Chain) {
	t.Helper()
	if !assert.Equal(t, exp, got) {
		eb, err := toml.Marshal(exp)
		require.NoError(t, err)
		gb, err := toml.Marshal(got)
		require.NoError(t, err)
		t.Log("exp:", string(eb))
		t.Log("got:", string(gb))
		t.Log("diff:", diff.Diff(string(eb), string(gb)))
	}
}

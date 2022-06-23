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
	t.Run("fallback", func(t *testing.T) {
		var got v2.Chain
		got.SetDefaults(nil)
		exp := config.V2Fallback()

		assertChainsEqual(t, exp, got)
	})
	for id, exp := range config.V2Defaults() {
		t.Run(fmt.Sprintf("%d", id), func(t *testing.T) {
			var got v2.Chain
			id := utils.NewBigI(id)
			got.SetDefaults(id)

			assertChainsEqual(t, exp, got)
		})
	}
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

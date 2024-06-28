package evm

import (
	"bytes"
	_ "embed"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_resolveup(t *testing.T) {
	t.Run("resolve up migrations", func(t *testing.T) {
		type test struct {
			name   string
			upTmpl string
		}
		cases := []test{
			{
				name:   "forwarders",
				upTmpl: forwardersUpTmpl,
			},
		}
		for _, tt := range cases {
			t.Run("do nothing for evm schema", func(t *testing.T) {
				out := &bytes.Buffer{}
				err := resolve(out, tt.upTmpl, Cfg{Schema: "evm", ChainID: big.NewI(int64(3266))})
				require.NoError(t, err)
				assert.Equal(t, "-- Do nothing for `evm` schema for backward compatibility\n", out.String())
			})

			t.Run("err no chain id", func(t *testing.T) {
				out := &bytes.Buffer{}
				err := resolve(out, tt.upTmpl, Cfg{Schema: "evm_213"})
				require.Error(t, err)
				assert.Empty(t, out.String())
			})

			t.Run("err no schema", func(t *testing.T) {
				out := &bytes.Buffer{}
				err := resolve(out, tt.upTmpl, Cfg{ChainID: big.NewI(int64(3266))})
				require.Error(t, err)
				assert.Empty(t, out.String())
			})

			t.Run("ok", func(t *testing.T) {
				out := &bytes.Buffer{}
				err := resolve(out, tt.upTmpl, Cfg{Schema: "evm_3266", ChainID: big.NewI(int64(3266))})
				require.NoError(t, err)
				lines := strings.Split(out.String(), "\n")
				assert.Greater(t, len(lines), 2)
				assert.Contains(t, out.String(), "CREATE TABLE evm_3266")
			})
		}
	})
}

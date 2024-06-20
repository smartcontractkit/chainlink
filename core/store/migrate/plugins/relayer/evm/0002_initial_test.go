package evm

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_resolveup(t *testing.T) {
	type args struct {
		val Cfg
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "evm template",
			args: args{
				val: Cfg{
					Schema:  "evm",
					ChainID: big.NewI(int64(3266)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := resolveUp(out, tt.args.val)
			require.NoError(t, err)
			assert.NotEmpty(t, out.String())
		})
	}
}

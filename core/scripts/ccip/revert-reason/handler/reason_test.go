package handler

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/revert-reason/config"
)

func Test_RevertReasonFromTx(t *testing.T) {
	type fields struct {
		cfg *config.Config
	}
	type args struct {
		txHash string
	}
	var tests []struct {
		name     string
		fields   fields
		args     args
		expected string
	} // TODO: Add test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &BaseHandler{
				cfg: tt.fields.cfg,
			}
			got, err := h.RevertReasonFromTx(tt.args.txHash)
			require.NoError(t, err)
			require.Equal(t, tt.expected, got)
		})
	}
}

func Test_RevertReasonFromErrorCodeString(t *testing.T) {
	type fields struct {
		cfg *config.Config
	}
	type args struct {
		errorCodeString string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected string
	}{
		{
			name:   "decode error string",
			fields: fields{cfg: &config.Config{}},
			args: args{
				errorCodeString: "0x4e487b710000000000000000000000000000000000000000000000000000000000000032",
			},
			expected: "If you access an array, bytesN or an array slice at an out-of-bounds or negative index (i.e. x[i] where i >= x.length or i < 0).",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &BaseHandler{
				cfg: tt.fields.cfg,
			}
			got, err := h.RevertReasonFromErrorCodeString(tt.args.errorCodeString)
			require.NoError(t, err)
			require.Equal(t, tt.expected, got)
		})
	}
}

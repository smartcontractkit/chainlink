package evm

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_resolve(t *testing.T) {
	type args struct {
		in  string
		val Cfg
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "evm template",
			args: args{
				val: Cfg{
					Schema:  "evm",
					ChainID: 3266,
				},
				in: "schema={{.Schema}}, chainID={{.ChainID}}",
			},
			wantOut: "schema=evm, chainID=3266",
		},
		{
			name: "unknown template",
			args: args{
				val: Cfg{
					Schema:  "evm",
					ChainID: 3266,
				},
				in: "schema={{.WrongField}}, chainID={{.ChainID}}",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := resolve(out, tt.args.in, tt.args.val)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantOut, out.String())
			}
		})
	}
}

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
					ChainID: 3266,
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

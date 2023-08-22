package relay

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifier_UnmarshalString(t *testing.T) {
	type fields struct {
		Network Network
		ChainID ChainID
	}
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		want    fields
		args    args
		wantErr bool
	}{
		{name: "evm",
			args:    args{s: "evm.1"},
			wantErr: false,
			want:    fields{Network: EVM, ChainID: "1"},
		},
		{name: "bad network",
			args:    args{s: "notANetwork.1"},
			wantErr: true,
		},
		{name: "bad pattern",
			args:    args{s: "evm_1"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ID{}
			err := i.UnmarshalString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Identifier.UnmarshalString() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want.Network, i.Network)
			assert.Equal(t, tt.want.ChainID, i.ChainID)
		})
	}
}

func TestNewID(t *testing.T) {
	type args struct {
		n Network
		c ChainID
	}
	tests := []struct {
		name    string
		args    args
		want    ID
		wantErr bool
	}{
		{name: "good evm",
			args: args{n: EVM, c: "1"},
			want: ID{Network: EVM, ChainID: "1"},
		},
		{name: "bad evm",
			args:    args{n: EVM, c: "not a number"},
			want:    ID{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewID(tt.args.n, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got, "got id %v", got)
		})
	}
}

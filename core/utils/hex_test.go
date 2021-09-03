package utils

import (
	"encoding/hex"
	"testing"
)

func TestIsHexBytes(t *testing.T) {
	type args struct {
		arr []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "hex string with prefix",
			args: args{
				arr: []byte("0x" + hex.EncodeToString([]byte(`test`))),
			},
			want: true,
		},
		{
			name: "hex string without prefix",
			args: args{
				arr: []byte(hex.EncodeToString([]byte(`test`))),
			},
			want: true,
		},
		{
			name: "not a hex string",
			args: args{
				arr: []byte(`123 not hex`),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHexBytes(tt.args.arr); got != tt.want {
				t.Errorf("IsHexBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

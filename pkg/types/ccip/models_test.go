package ccip

import "testing"

func TestHash_String(t *testing.T) {
	tests := []struct {
		name string
		h    Hash
		want string
	}{
		{
			name: "empty",
			h:    Hash{},
			want: "0x0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name: "1..",
			h:    Hash{1},
			want: "0x0100000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name: "1..000..1",
			h:    [32]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			want: "0x0100000000000000000000000000000000000000000000000000000000000001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

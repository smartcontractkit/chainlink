package ccip

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_byte32Slice(t *testing.T) {
	tooLong := make([]byte, 33)
	tooLong[32] = 32
	tooShort := make([]byte, 31)
	tooShort[30] = 30
	type args struct {
		pbVal [][]byte
	}
	tests := []struct {
		name     string
		args     args
		ifaceVal [][32]byte
		wantErr  bool
	}{
		{name: "empty", args: args{pbVal: [][]byte{}}, ifaceVal: [][32]byte{}, wantErr: false},
		{name: "non-empty",
			args: args{
				pbVal: [][]byte{
					{0: 1, 31: 2},
					{0: 3, 31: 4},
				},
			},
			ifaceVal: [][32]byte{
				{0: 1, 31: 2},
				{0: 3, 31: 4},
			},
			wantErr: false},
		{name: "too long", args: args{pbVal: [][]byte{tooLong}}, ifaceVal: nil, wantErr: true},
		{name: "too short", args: args{pbVal: [][]byte{tooShort}}, ifaceVal: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("pb-to-iface %s", tt.name), func(t *testing.T) {
			t.Parallel()
			got, err := byte32Slice(tt.args.pbVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("byte32Slice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.ifaceVal) {
				t.Errorf("byte32Slice() = %v, want %v", got, tt.ifaceVal)
			}
		})

		t.Run(fmt.Sprintf("iface-to-pb %s", tt.name), func(t *testing.T) {
			t.Parallel()
			// there are no errors in this direction so skip tests that expect errors
			if tt.wantErr {
				return
			}
			got := byte32SliceToPB(tt.ifaceVal)

			if !reflect.DeepEqual(got, tt.args.pbVal) {
				t.Errorf("byte32SlicePB() = %v, want %v", got, tt.args.pbVal)
			}
		})
	}

	// special case for nil
	t.Run("nil pb-to-iface", func(t *testing.T) {
		t.Parallel()
		got, err := byte32Slice(nil)
		if err != nil {
			t.Errorf("byte32Slice() error = %v, wantErr %v", err, false)
			return
		}
		expected := [][32]byte{}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("byte32Slice() = %v, want %v", got, expected)
		}
	})

	t.Run("nil iface-to-pb", func(t *testing.T) {
		t.Parallel()
		// there are no errors in this direction so skip tests that expect errors
		got := byte32SliceToPB(nil)
		expected := [][]byte{}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("byte32SlicePB() = %v, want %v", got, expected)
		}
	})
}

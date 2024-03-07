package cciptypes

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	json2 "github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

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

func TestAddress_JSON(t *testing.T) {
	addrLower := "0xe8bade28e08b469b4eeec35b9e48b2ce49fb3fc9"
	addrEIP55 := "0xE8BAde28E08B469B4EeeC35b9E48B2Ce49FB3FC9"

	t.Run("arrays", func(t *testing.T) {
		addrArr := []Address{Address(addrLower), Address(addrEIP55)}
		b, err := json2.Marshal(addrArr)
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(`["%s","%s"]`, addrLower, addrLower), string(b))

		evmAddrArr := []common.Address{common.HexToAddress(addrLower), common.HexToAddress(addrEIP55)}
		bEvm, err := json2.Marshal(evmAddrArr)
		assert.NoError(t, err)
		assert.Equal(t, b, bEvm)
	})

	t.Run("maps", func(t *testing.T) {
		m := map[Address]int{Address(addrEIP55): 14}
		b, err := json2.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(`{"%s":14}`, addrLower), string(b), "should be lower when marshalled")

		m2 := map[Address]int{}
		err = json2.Unmarshal(b, &m2)
		assert.NoError(t, err)
		assert.Equal(t, m, m2, "should be eip55 when unmarshalled")
	})
}

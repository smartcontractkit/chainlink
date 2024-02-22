package cciptypes

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
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
	addr1 := utils.RandomAddress()

	addrArr := []Address{Address(addr1.String())}
	evmAddrArr := []common.Address{addr1}

	b, err := json.Marshal(addrArr)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf(`["%s"]`, strings.ToLower(addr1.String())), string(b))

	b2, err := json.Marshal(evmAddrArr)
	assert.NoError(t, err)
	assert.Equal(t, string(b), string(b2), "marshal should produce the same result for common.Address and cciptypes.Address")

	var unmarshalledAddr []Address
	err = json.Unmarshal(b, &unmarshalledAddr)
	assert.NoError(t, err)
	assert.Equal(t, addrArr[0], unmarshalledAddr[0])

	var unmarshalledEvmAddr []common.Address
	err = json.Unmarshal(b, &unmarshalledEvmAddr)
	assert.NoError(t, err)
	assert.Equal(t, evmAddrArr[0], unmarshalledEvmAddr[0])
}

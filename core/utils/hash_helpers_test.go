package utils_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestNewHash(t *testing.T) {
	t.Parallel()

	h1 := utils.NewHash()
	h2 := utils.NewHash()
	assert.NotEqual(t, h1, h2)
	assert.NotEqual(t, h1, common.HexToHash("0x0"))
	assert.NotEqual(t, h2, common.HexToHash("0x0"))
}

func TestPadByteToHash(t *testing.T) {
	t.Parallel()

	h := utils.PadByteToHash(1)
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000001", h.String())
}

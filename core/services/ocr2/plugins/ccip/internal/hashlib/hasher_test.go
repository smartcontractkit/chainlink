package hashlib

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestHashInternal(t *testing.T) {
	k := NewKeccakCtx()
	h1 := common.HexToHash("0x1")
	h2 := common.HexToHash("0x2")
	h12 := k.HashInternal(h1, h2)
	h21 := k.HashInternal(h2, h1)
	assert.Equal(t, h12, h21)
}

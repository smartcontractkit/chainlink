package hashlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashInternal(t *testing.T) {
	k := NewKeccakCtx()
	h1 := [32]byte{1}
	h2 := [32]byte{2}
	h12 := k.HashInternal(h1, h2)
	h21 := k.HashInternal(h2, h1)
	assert.Equal(t, h12, h21)
}

package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_groupID(t *testing.T) {
	g := &counter{}
	assert.Equal(t, [32]byte{}, g.Bytes())
	g.Inc()
	assert.Equal(t, [32]byte{31: 1}, g.Bytes())
	g.Inc()
	assert.Equal(t, [32]byte{31: 2}, g.Bytes())
}

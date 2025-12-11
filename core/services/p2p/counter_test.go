package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func Test_groupID(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	g := &counter{}
	assert.Equal(t, [32]byte{}, g.Bytes())
	g.Inc()
	assert.Equal(t, [32]byte{31: 1}, g.Bytes())
	g.Inc()
	assert.Equal(t, [32]byte{31: 2}, g.Bytes())
}

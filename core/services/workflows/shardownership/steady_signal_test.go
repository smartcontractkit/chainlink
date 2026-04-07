package shardownership

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSteadySignal_SkipCommittedOwnerCheck(t *testing.T) {
	var nilSig *SteadySignal
	assert.False(t, nilSig.SkipCommittedOwnerCheck())

	s := NewSteadySignal()
	assert.False(t, s.SkipCommittedOwnerCheck())

	s.ObserveRoutingSteady(true)
	assert.True(t, s.SkipCommittedOwnerCheck())

	s.ObserveRoutingSteady(false)
	assert.False(t, s.SkipCommittedOwnerCheck())

	s.ObserveRoutingSteady(true)
	assert.True(t, s.SkipCommittedOwnerCheck())
	s.Invalidate()
	assert.False(t, s.SkipCommittedOwnerCheck())
}

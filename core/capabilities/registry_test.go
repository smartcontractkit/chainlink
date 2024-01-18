package capabilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCapability struct {
	Validatable
	CapabilityInfo
}

func TestRegistry(t *testing.T) {
	r := NewRegistry()

	id := Stringer("capability-1")
	ci, err := NewCapabilityInfo(
		id,
		CapabilityTypeAction,
		"capability-1-description",
		"v1.0.0",
	)
	require.NoError(t, err)

	c := &mockCapability{CapabilityInfo: ci}
	err = r.Add(c)
	require.NoError(t, err)

	gc, err := r.Get(id)
	require.NoError(t, err)

	assert.Equal(t, c, gc)

	cs := r.List()
	assert.Len(t, cs, 1)
	assert.Equal(t, c, cs[0])
}

func TestRegistry_NoDuplicateIDs(t *testing.T) {
	r := NewRegistry()

	id := Stringer("capability-1")
	ci, err := NewCapabilityInfo(
		id,
		CapabilityTypeAction,
		"capability-1-description",
		"v1.0.0",
	)
	require.NoError(t, err)

	c := &mockCapability{CapabilityInfo: ci}
	err = r.Add(c)
	require.NoError(t, err)

	ci, err = NewCapabilityInfo(
		id,
		CapabilityTypeReport,
		"capability-2-description",
		"v1.0.0",
	)
	require.NoError(t, err)
	c2 := &mockCapability{CapabilityInfo: ci}

	err = r.Add(c2)
	assert.ErrorContains(t, err, "capability with id: capability-1 already exists")
}

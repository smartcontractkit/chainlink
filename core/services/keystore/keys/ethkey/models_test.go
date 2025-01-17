package ethkey

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestState_AddUsage(t *testing.T) {
	state := State{
		tags: make(map[string]int),
	}

	err := state.Tag("process1")
	require.NoError(t, err)
	require.Equal(t, 1, state.tags["process1"])

	err = state.Tag("process1")
	require.NoError(t, err)
	require.Equal(t, 2, state.tags["process1"])

	err = state.Tag("")
	require.Error(t, err)
	require.Equal(t, "cannot add usage: label string is empty", err.Error())
}

func TestState_RemoveUsage(t *testing.T) {
	state := State{
		tags: map[string]int{
			"process1": 2,
			"process2": 1,
		},
	}

	err := state.Untag("process1")
	require.NoError(t, err)
	require.Equal(t, 1, state.tags["process1"])

	err = state.Untag("process1")
	require.NoError(t, err)
	_, exists := state.tags["process1"]
	require.False(t, exists)

	err = state.Untag("process2")
	require.NoError(t, err)
	_, exists = state.tags["process2"]
	require.False(t, exists)

	err = state.Untag("")
	require.Error(t, err)
	require.Equal(t, "cannot remove usage: label string is empty", err.Error())
}

func TestState_HasTag(t *testing.T) {
	state := State{
		tags: map[string]int{
			"process1": 2,
			"process2": 1,
		},
	}

	hasTag, err := state.HasTag("process1")
	require.NoError(t, err)
	require.True(t, hasTag)

	hasTag, err = state.HasTag("process2")
	require.NoError(t, err)
	require.True(t, hasTag)

	hasTag, err = state.HasTag("process3")
	require.NoError(t, err)
	require.False(t, hasTag)

	hasTag, err = state.HasTag("")
	require.Error(t, err)
	require.False(t, hasTag)
	require.Equal(t, "cannot check usage: label string is empty", err.Error())
}

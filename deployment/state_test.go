package deployment

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNode_GetState_UsageExample(t *testing.T) {
	mgr := NewInMemoryStateManager(FooState{"blah"})
	env := Environment{state: mgr}

	_, err := GetState(env, StateKey[BarState]{})
	require.ErrorIs(t, err, ErrStateNotFound)
}

func TestNode_GetState(t *testing.T) {
	env := Environment{state: NewInMemoryStateManager(FooState{"foo"}, BarState{"bar"})}

	fooState, err := GetState(env, StateKey[FooState]{})
	require.NoError(t, err)
	require.Equal(t, "foo", fooState.foo)
	barState, err := GetState(env, StateKey[BarState]{})
	require.NoError(t, err)
	require.Equal(t, "bar", barState.bar)
}

func TestNode_GetState_NotFound(t *testing.T) {
	env := Environment{state: NewInMemoryStateManager(FooState{"blah"})}

	_, err := GetState(env, StateKey[BarState]{})
	require.ErrorIs(t, err, ErrStateNotFound)
}

// Set up state structs.

var _ StateRepresentation = FooState{}

type FooState struct {
	foo string
}

func (f FooState) MarshalJSON() ([]byte, error) {
	// Alias to avoid recursive calls
	type Alias FooState
	return json.MarshalIndent(&struct{ Alias }{Alias: Alias(f)}, "", " ")
}

var _ StateRepresentation = BarState{}

type BarState struct {
	bar string
}

func (b BarState) MarshalJSON() ([]byte, error) {
	// Alias to avoid recursive calls
	type Alias BarState
	return json.MarshalIndent(&struct{ Alias }{Alias: Alias(b)}, "", " ")
}

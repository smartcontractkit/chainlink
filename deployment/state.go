package deployment

import (
	"errors"
	"reflect"
)

var (
	ErrInvalidStateType = errors.New("invalid state type")
	ErrStateNotFound    = errors.New("no state with this type can be found")

	ErrStateManagerNotFound = errors.New("the environment did not include a state manager")
)

// StateKey is a carrier of type information, to allow a static type to be used as a key, conveniently, without
// having to mess with the reflection package in client code. Its basic usage is:
//
//	key := StateKey[MyStruct]{}
//
// The type parameter information is carried on the inner field, and infrastructure code can obtain that type for use
// in underlying code.
type StateKey[T StateRepresentation] struct {
	_ T // Expose the type to generics.
}

// StateManager is responsible for looking up and supplying state to the typesafe GetState wrapper function. Different
// implementations will obtain the state (and populate the given struct) from different sources, whether in-memory,
// from json files, via RPCs or any other mechanism. The simple InMemoryStateManager can be used for in-memory
// Environments. More complex StateManager implementations can be found where those Environment variations are
// implemented.
type StateManager interface {
	// getState returns a struct of the type given in the "key" parameter, populated with any state data obtained
	// using that same key. If no such state exists, it returns ErrStateNotFound. If it finds any errors in finding
	// and marshalling the data in question, it returns those errors.
	getState(key reflect.Type) (StateRepresentation, error)
}

// InMemoryStateManager is a non-thread-safe StateManager implementation, which simply maintains an internal
// map keyed by reflect.Type. The struct value of the map is required to be of the type represented by the key.
//
// This can be used in a test quite simply to setup state for use by the GetState function, outside the ChangeSet
// under test, injected into the environment, and then the ChangeSet can access that state via the GetState function.
type InMemoryStateManager struct {
	// Maintains the map of type to state struct. Public for testing. This map is unavailable within a ChangeSet
	State map[reflect.Type]StateRepresentation
}

func (m InMemoryStateManager) getState(key reflect.Type) (StateRepresentation, error) {
	state, ok := m.State[key]
	if !ok {
		return nil, ErrStateNotFound
	}
	return state, nil
}

func NewInMemoryStateManager(state ...StateRepresentation) InMemoryStateManager {
	mgr := InMemoryStateManager{map[reflect.Type]StateRepresentation{}}
	for _, s := range state {
		stateType := reflect.TypeOf(s)
		mgr.State[stateType] = s
	}
	return mgr
}

// GetState obtains a domain-specific state struct (as defined by the product/domain team) from the environment.
// It relies on a StateManager inside the Environment to manage the actual fetching and marshalling of such state. The
// basic usage pattern is as follows:
//
// state, err := GetState(env, StateKey[MyCustomStruct]{})
//
// If no such state can be found with that key (e.g., with that struct), then the method will return ErrStateNotFound.
// If there are underlying errors in fetching, marshalling, parsing, or any other preparation, these will be surfaced.
// In tests, state is obtained from an in-memory state store, which simply maintains a map of such state. In other
// contexts, state may derive from json files or from a back-end data store.
//
// This function is intended to be used from within a ChangeSet function, and has undefined behavior outside of that
// context.
func GetState[T StateRepresentation](env Environment, key StateKey[T]) (*T, error) {
	mgr := env.state
	if mgr == nil {
		return nil, ErrStateManagerNotFound
	}
	keyType := reflect.TypeOf(key)
	typeField := keyType.Field(0)
	stateType := typeField.Type
	rawState, err := mgr.getState(stateType)
	if err != nil {
		return nil, err
	}
	stronglyTypeState, ok := rawState.(T)
	if !ok {
		// Given the guarantees of StateManager.getState, this represents a bug in state manager, which should not
		// return an instance of anything except the provided struct type.
		return nil, ErrInvalidStateType
	}
	return &stronglyTypeState, nil
}

package pipeline

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrKeypathNotFound = errors.New("keypath not found")
	ErrVarsRoot        = errors.New("cannot get/set the root of a pipeline.Vars")
	ErrVarsSetNested   = errors.New("cannot set a nested key of a pipeline.Vars")

	variableRegexp = regexp.MustCompile(`\$\(\s*([a-zA-Z0-9_\.]+)\s*\)`)
)

type Vars struct {
	vars map[string]interface{}
}

// NewVarsFrom creates new Vars from the given map.
// If the map is nil, a new map instance will be created.
func NewVarsFrom(m map[string]interface{}) Vars {
	if m == nil {
		m = make(map[string]interface{})
	}
	return Vars{vars: m}
}

// Get returns the value for the given keypath or error.
// The keypath can consist of one or two parts, e.g. "foo" or "foo.bar".
// The second part of the keypath can be an index of a slice.
func (vars Vars) Get(keypathStr string) (interface{}, error) {
	keypathStr = strings.TrimSpace(keypathStr)
	keypath, err := NewKeypathFromString(keypathStr)
	if err != nil {
		return nil, err
	}
	if keypath.NumParts == 0 {
		return nil, ErrVarsRoot
	}

	var val interface{}
	var exists bool

	if keypath.NumParts >= 1 {
		val, exists = vars.vars[keypath.Part0]
		if !exists {
			return nil, errors.Wrapf(ErrKeypathNotFound, "key %v / keypath %v", keypath.Part0, keypathStr)
		}
	}

	if keypath.NumParts == 2 {
		switch v := val.(type) {
		case map[string]interface{}:
			val, exists = v[keypath.Part1]
			if !exists {
				return nil, errors.Wrapf(ErrKeypathNotFound, "key %v / keypath %v", keypath.Part1, keypathStr)
			}
		case []interface{}:
			idx, err := strconv.ParseInt(keypath.Part1, 10, 64)
			if err != nil {
				return nil, errors.Wrapf(ErrKeypathNotFound, "could not parse key as integer: %v", err)
			} else if idx < 0 || idx > int64(len(v)-1) {
				return nil, errors.Wrapf(ErrIndexOutOfRange, "index %v out of range (length %v / keypath %v)", idx, len(v), keypathStr)
			}
			val = v[idx]
		default:
			return nil, errors.Wrapf(ErrKeypathNotFound, "value at key '%v' is a %T, not a map or slice", keypath.Part0, val)
		}
	}

	return val, nil
}

// Set sets a top-level variable specified by dotID.
// Returns error if either dotID is empty or it is a compound keypath.
func (vars Vars) Set(dotID string, value interface{}) error {
	dotID = strings.TrimSpace(dotID)
	if len(dotID) == 0 {
		return ErrVarsRoot
	} else if strings.Contains(dotID, KeypathSeparator) {
		return errors.Wrapf(ErrVarsSetNested, "%s", dotID)
	}

	vars.vars[dotID] = value

	return nil
}

// Copy makes a copy of Vars by copying the underlying map.
// Used by scheduler for new tasks to avoid data races.
func (vars Vars) Copy() Vars {
	newVars := make(map[string]interface{})
	// No need to copy recursively, because only the top-level map is mutable (see Set()).
	for k, v := range vars.vars {
		newVars[k] = v
	}
	return NewVarsFrom(newVars)
}

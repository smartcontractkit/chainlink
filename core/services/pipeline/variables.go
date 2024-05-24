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
// The keypath can consist of one or more parts, e.g. "foo" or "foo.6.a.b".
// Every part except for the first one can be an index of a slice.
func (vars Vars) Get(keypathStr string) (interface{}, error) {
	keypathStr = strings.TrimSpace(keypathStr)
	keypath, err := NewKeypathFromString(keypathStr)
	if err != nil {
		return nil, err
	}
	if len(keypath.Parts) == 0 {
		return nil, ErrVarsRoot
	}

	var exists bool
	var currVal interface{} = vars.vars
	for i, part := range keypath.Parts {
		switch v := currVal.(type) {
		case map[string]interface{}:
			currVal, exists = v[part]
			if !exists {
				return nil, errors.Wrapf(ErrKeypathNotFound, "key %v (segment %v in keypath %v)", part, i, keypathStr)
			}
		case []interface{}:
			idx, err := strconv.ParseInt(part, 10, 64)
			if err != nil {
				return nil, errors.Wrapf(ErrKeypathNotFound, "could not parse key as integer: %v", err)
			} else if idx < 0 || idx > int64(len(v)-1) {
				return nil, errors.Wrapf(ErrIndexOutOfRange, "index %v out of range (segment %v of length %v in keypath %v)", idx, i, len(v), keypathStr)
			}
			currVal = v[idx]
		default:
			return nil, errors.Wrapf(ErrKeypathNotFound, "value at key '%v' is a %T, not a map or slice", part, currVal)
		}
	}

	return currVal, nil
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

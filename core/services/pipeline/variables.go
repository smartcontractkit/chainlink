package pipeline

import (
	"maps"
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
	vars map[string]any
	// timestamps records the source timestamp (unix ms) of each task's result,
	// keyed by dotID. Zero means unknown. Written by the scheduler as results
	// are reported, and read by aggregation tasks (see the median staleness
	// gate). Like vars, it is only mutated in the scheduler's Run loop, and
	// copies handed to tasks are snapshots.
	timestamps map[string]uint64
}

// NewVarsFrom creates new Vars from the given map.
// If the map is nil, a new map instance will be created.
func NewVarsFrom(m map[string]any) Vars {
	if m == nil {
		m = make(map[string]any)
	}
	return Vars{vars: m}
}

// Get returns the value for the given keypath or error.
// The keypath can consist of one or more parts, e.g. "foo" or "foo.6.a.b".
// Every part except for the first one can be an index of a slice.
func (vars Vars) Get(keypathStr string) (any, error) {
	keypathStr = strings.TrimSpace(keypathStr)
	keypath, err := NewKeypathFromString(keypathStr)
	if err != nil {
		return nil, err
	}
	if len(keypath.Parts) == 0 {
		return nil, ErrVarsRoot
	}

	var exists bool
	var currVal any = vars.vars
	for i, part := range keypath.Parts {
		switch v := currVal.(type) {
		case map[string]any:
			currVal, exists = v[part]
			if !exists {
				return nil, errors.Wrapf(ErrKeypathNotFound, "key %v (segment %v in keypath %v)", part, i, keypathStr)
			}
		case []any:
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
func (vars Vars) Set(dotID string, value any) error {
	dotID = strings.TrimSpace(dotID)
	if len(dotID) == 0 {
		return ErrVarsRoot
	} else if strings.Contains(dotID, KeypathSeparator) {
		return errors.Wrapf(ErrVarsSetNested, "%s", dotID)
	}

	vars.vars[dotID] = value

	return nil
}

// SetTimestamp records the source timestamp (unix ms) of a task's result,
// keyed by dotID. Zero means unknown.
func (vars *Vars) SetTimestamp(dotID string, timestamp uint64) {
	if timestamp == 0 {
		// Unknown timestamps are never consulted, so there is nothing to record.
		return
	}
	if vars.timestamps == nil {
		vars.timestamps = make(map[string]uint64)
	}
	vars.timestamps[dotID] = timestamp
}

// GetTimestamp returns the recorded source timestamp for a task's dotID.
// Zero means unknown (never recorded or unknown).
func (vars Vars) GetTimestamp(dotID string) uint64 {
	if vars.timestamps == nil {
		return 0
	}
	return vars.timestamps[dotID]
}

// Copy makes a copy of Vars by copying the underlying map.
// Used by scheduler for new tasks to avoid data races.
func (vars Vars) Copy() Vars {
	newVars := make(map[string]any)
	// No need to copy recursively, because only the top-level map is mutable (see Set()).
	maps.Copy(newVars, vars.vars)
	var newTimestamps map[string]uint64
	if vars.timestamps != nil {
		newTimestamps = make(map[string]uint64, len(vars.timestamps))
		maps.Copy(newTimestamps, vars.timestamps)
	}
	return Vars{vars: newVars, timestamps: newTimestamps}
}

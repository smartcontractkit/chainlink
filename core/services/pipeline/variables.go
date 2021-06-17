package pipeline

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrKeypathNotFound = errors.New("keypath not found")
	ErrKeypathTooDeep  = errors.New("keypath too deep (maximum 2 keys)")
	ErrVarsRoot        = errors.New("cannot get/set the root of a pipeline.Vars")

	variableRegexp = regexp.MustCompile(`\$\(\s*([a-zA-Z0-9_\.]+)\s*\)`)
)

type Vars struct {
	vars map[string]interface{}
}

func NewVarsFrom(m map[string]interface{}) Vars {
	if m == nil {
		m = make(map[string]interface{})
	}
	return Vars{vars: m}
}

func (vars Vars) Copy() Vars {
	m := make(map[string]interface{})
	for k, v := range vars.vars {
		m[k] = v
	}
	return Vars{vars: m}
}

func (vars Vars) Get(keypathStr string) (interface{}, error) {
	keypath, err := newKeypathFromString(keypathStr)
	if err != nil {
		return nil, err
	}

	numParts := keypath.NumParts()

	if numParts == 0 {
		return nil, ErrVarsRoot
	}

	var val interface{}
	var exists bool

	if numParts >= 1 {
		val, exists = vars.vars[string(keypath[0])]
		if !exists {
			return nil, errors.Wrapf(ErrKeypathNotFound, "key %v / keypath %v", string(keypath[0]), keypath.String())
		}
	}

	if numParts == 2 {
		switch v := val.(type) {
		case map[string]interface{}:
			val, exists = v[string(keypath[1])]
			if !exists {
				return nil, errors.Wrapf(ErrKeypathNotFound, "key %v / keypath %v", string(keypath[1]), keypath.String())
			}
		case []interface{}:
			idx, err := strconv.ParseInt(string(keypath[1]), 10, 64)
			if err != nil {
				return nil, errors.Wrapf(ErrKeypathNotFound, "could not parse key as integer: %v", err)
			} else if idx > int64(len(v)-1) {
				return nil, errors.Wrapf(ErrKeypathNotFound, "index %v out of range (length %v / keypath %v)", idx, len(v), keypath.String())
			}
			val = v[idx]
		default:
			return nil, errors.Wrapf(ErrKeypathNotFound, "value at key '%v' is a %T, not a map or slice", string(keypath[0]), val)
		}
	}

	return val, nil
}

func (vars Vars) Set(dotID string, value interface{}) {
	dotID = strings.TrimSpace(dotID)
	if len(dotID) == 0 {
		panic(ErrVarsRoot)
	} else if strings.IndexByte(dotID, keypathSeparator[0]) >= 0 {
		panic("cannot set a nested key of a pipeline.Vars")
	}
	vars.vars[dotID] = value
}

type Keypath [2][]byte

var keypathSeparator = []byte(".")

func newKeypathFromString(keypathStr string) (Keypath, error) {
	if len(keypathStr) == 0 {
		return Keypath{}, nil
	}
	// The bytes package uses platform-dependent hardware optimizations and
	// avoids the extra allocations that are required to work with strings.
	// Keypaths have to be parsed quite a bit, so let's do it well.
	kp := []byte(keypathStr)

	n := 1 + bytes.Count(kp, keypathSeparator)
	if n > 2 {
		return Keypath{}, errors.Wrapf(ErrKeypathTooDeep, "while parsing keypath '%v'", keypathStr)
	}
	idx := bytes.IndexByte(kp, keypathSeparator[0])
	if idx == -1 || idx == len(kp)-1 {
		return Keypath{kp, nil}, nil
	}
	return Keypath{kp[:idx], kp[idx+1:]}, nil
}

func (keypath Keypath) NumParts() int {
	switch {
	case keypath[0] == nil && keypath[1] == nil:
		return 0
	case keypath[0] != nil && keypath[1] == nil:
		return 1
	case keypath[0] == nil && keypath[1] != nil:
		panic("invariant violation: keypath part 1 is non-nil but part 0 is nil")
	default:
		return 2
	}
}

func (keypath Keypath) String() string {
	switch keypath.NumParts() {
	case 0:
		return "(empty)"
	case 1:
		return string(keypath[0])
	case 2:
		return string(keypath[0]) + string(keypathSeparator) + string(keypath[1])
	default:
		panic("invariant violation: keypath must have 0, 1, or 2 parts")
	}
}

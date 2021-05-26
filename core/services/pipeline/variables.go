package pipeline

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

var (
	ErrKeypathNotFound = errors.New("keypath not found")
	ErrKeypathTooDeep  = errors.New("keypath too deep (maximum 2 keys)")

	variableRegexp = regexp.MustCompile(`\$\(([a-zA-Z0-9_\.]+)\)`)
)

func ExpandVars(v string, vars Vars) (string, error) {
	var err error
	resolved := variableRegexp.ReplaceAllFunc([]byte(v), func(keypath []byte) []byte {
		val, err2 := vars.Get(string(keypath[2 : len(keypath)-1]))
		if err2 != nil {
			err = multierr.Append(err, err2)
			return nil
		} else if asErr, isErr := val.(error); isErr {
			err = multierr.Append(err, asErr)
			return nil
		}

		bs, err2 := json.Marshal(val)
		if err2 != nil {
			err = multierr.Append(err, err2)
			return nil
		}
		return bs
	})
	if err != nil {
		return "", err
	}
	return string(resolved), nil
}

type Vars map[string]interface{}

func NewVars() Vars {
	return make(Vars)
}

func (vars Vars) Get(keypath string) (interface{}, error) {
	keypathParts, err := keypathParts(keypath)
	if err != nil {
		return nil, err
	}

	if len(keypathParts[0]) == 0 && len(keypathParts[1]) == 0 {
		return (map[string]interface{})(vars), nil
	}

	var val interface{}
	var exists bool

	if len(keypathParts[0]) > 0 {
		val, exists = vars[string(keypathParts[0])]
		if !exists {
			return nil, errors.Wrapf(ErrKeypathNotFound, "key %v / keypath %v", string(keypathParts[0]), keypathPartsToStrings(keypathParts))
		}
	}

	if len(keypathParts[1]) > 0 {
		switch v := val.(type) {
		case map[string]interface{}:
			val, exists = v[string(keypathParts[1])]
			if !exists {
				return nil, errors.Wrapf(ErrKeypathNotFound, "key %v / keypath %v", string(keypathParts[1]), keypathPartsToStrings(keypathParts))
			}
		case []interface{}:
			idx, err := strconv.ParseInt(string(keypathParts[1]), 10, 64)
			if err != nil {
				return nil, err
			} else if idx > int64(len(v)-1) {
				return nil, errors.Wrapf(ErrKeypathNotFound, "index %v out of range (length %v / keypath %v)", idx, len(v), keypathPartsToStrings(keypathParts))
			}
			val = v[idx]
		}
	}
	return val, nil
}

var keypathSeparator = []byte(".")

func keypathParts(keypath string) ([2][]byte, error) {
	if len(keypath) == 0 {
		return [2][]byte{}, nil
	}
	// The bytes package uses platform-dependent hardware optimizations and
	// avoids the extra allocations that are required to work with strings.
	// Keypaths have to be parsed quite a bit, so let's do it well.
	kp := []byte(keypath)

	n := 1 + bytes.Count(kp, keypathSeparator)
	if n > 2 {
		return [2][]byte{}, ErrKeypathTooDeep
	}
	idx := bytes.IndexByte(kp, keypathSeparator[0])
	if idx == -1 || idx == len(kp)-1 {
		return [2][]byte{kp, nil}, nil
	}
	return [2][]byte{kp[:idx], kp[idx+1:]}, nil
}

func keypathPartsToStrings(bs [2][]byte) []string {
	var s []string
	for _, b := range bs {
		s = append(s, string(b))
	}
	return s
}

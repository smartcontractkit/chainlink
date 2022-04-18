package pipeline

import (
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrWrongKeypath = errors.New("wrong keypath format")
)

const KeypathSeparator = "."

// Keypath contains keypath parsed by NewKeypathFromString.
type Keypath struct {
	NumParts int    // can be 0, 1 or 2
	Part0    string // can be empty string if NumParts is 0
	Part1    string // can be empty string if NumParts is 0 or 1
}

// NewKeypathFromString creates a new Keypath from the given string.
// Returns error if it fails to parse the given keypath string.
func NewKeypathFromString(keypathStr string) (Keypath, error) {
	if len(keypathStr) == 0 {
		return Keypath{}, nil
	}

	parts := strings.Split(keypathStr, KeypathSeparator)

	switch len(parts) {
	case 0:
		return Keypath{}, errors.Wrapf(ErrWrongKeypath, "empty keypath")
	case 1:
		if len(parts[0]) > 0 {
			return Keypath{1, parts[0], ""}, nil
		}
	case 2:
		if len(parts[0]) > 0 && len(parts[1]) > 0 {
			return Keypath{2, parts[0], parts[1]}, nil
		}
	}

	return Keypath{}, errors.Wrapf(ErrWrongKeypath, "while parsing keypath '%v'", keypathStr)
}

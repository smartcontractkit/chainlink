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
	Parts []string
}

// NewKeypathFromString creates a new Keypath from the given string.
// Returns error if it fails to parse the given keypath string.
func NewKeypathFromString(keypathStr string) (Keypath, error) {
	if len(keypathStr) == 0 {
		return Keypath{}, nil
	}

	parts := strings.Split(keypathStr, KeypathSeparator)
	if len(parts) == 0 {
		return Keypath{}, errors.Wrapf(ErrWrongKeypath, "empty keypath")
	}
	for i, part := range parts {
		if len(part) == 0 {
			return Keypath{}, errors.Wrapf(ErrWrongKeypath, "empty keypath segment at index %d", i)
		}
	}

	return Keypath{parts}, nil
}

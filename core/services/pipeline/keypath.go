package pipeline

import (
	"bytes"

	"github.com/pkg/errors"
)

var (
	ErrWrongKeypath  = errors.New("wrong keypath format")
	KeypathSeparator = []byte(".")
)

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

	parts := bytes.Split([]byte(keypathStr), KeypathSeparator)

	switch len(parts) {
	case 0:
		return Keypath{}, errors.Wrapf(ErrWrongKeypath, "empty keypath")
	case 1:
		if len(parts[0]) > 0 {
			part0 := string(parts[0])
			return Keypath{1, part0, ""}, nil
		}
	case 2:
		if len(parts[0]) > 0 && len(parts[1]) > 0 {
			part0 := string(parts[0])
			part1 := string(parts[1])
			return Keypath{2, part0, part1}, nil
		}
	}

	return Keypath{}, errors.Wrapf(ErrWrongKeypath, "while parsing keypath '%v'", keypathStr)
}

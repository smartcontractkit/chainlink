package config

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
)

type InvalidError = config.ErrInvalid

// NewErrDuplicate returns an InvalidError with a standard duplicate message.
func NewErrDuplicate(name string, value any) InvalidError {
	return config.NewErrDuplicate(name, value)
}

type MissingError = config.ErrMissing

type EmptyError = config.ErrEmpty

// UniqueStrings is a helper for tracking unique values in string form.
type UniqueStrings = config.UniqueStrings

type OverrideError struct {
	Name string
}

func (e OverrideError) Error() string {
	return e.Name + ": overrides (duplicate keys or list elements) are not allowed for multiple secrets files"
}

type DeprecatedError struct {
	Name    string
	Version semver.Version
}

func (e DeprecatedError) Error() string {
	when := "a future version"
	if e.Version != (semver.Version{}) {
		when = fmt.Sprintf("version %s", e.Version)
	}
	return fmt.Sprintf("%s: is deprecated and will be removed in %s", e.Name, when)
}

// ErrInvalid is retained for backward compatibility.
type ErrInvalid = InvalidError //nolint:errname // backward compatibility

// ErrMissing is retained for backward compatibility.
type ErrMissing = MissingError //nolint:errname // backward compatibility

// ErrEmpty is retained for backward compatibility.
type ErrEmpty = EmptyError //nolint:errname // backward compatibility

// ErrOverride is retained for backward compatibility.
type ErrOverride = OverrideError //nolint:errname // backward compatibility

// ErrDeprecated is retained for backward compatibility.
type ErrDeprecated = DeprecatedError //nolint:errname // backward compatibility

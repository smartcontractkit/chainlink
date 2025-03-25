package config

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
)

type ErrInvalid = config.ErrInvalid

// NewErrDuplicate returns an ErrInvalid with a standard duplicate message.
func NewErrDuplicate(name string, value any) ErrInvalid {
	return config.NewErrDuplicate(name, value)
}

type ErrMissing = config.ErrMissing

type ErrEmpty = config.ErrEmpty

// UniqueStrings is a helper for tracking unique values in string form.
type UniqueStrings = config.UniqueStrings

type ErrOverride struct {
	Name string
}

func (e ErrOverride) Error() string {
	return e.Name + ": overrides (duplicate keys or list elements) are not allowed for multiple secrets files"
}

type ErrDeprecated struct {
	Name    string
	Version semver.Version
}

func (e ErrDeprecated) Error() string {
	when := "a future version"
	if e.Version != (semver.Version{}) {
		when = fmt.Sprintf("version %s", e.Version)
	}
	return fmt.Sprintf("%s: is deprecated and will be removed in %s", e.Name, when)
}

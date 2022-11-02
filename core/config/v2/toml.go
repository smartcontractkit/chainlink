package v2

import (
	"io"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
)

// DecodeTOML decodes toml from r in to v.
// Requires strict field matches and returns full toml.StrictMissingError details.
func DecodeTOML(r io.Reader, v any) error {
	d := toml.NewDecoder(r).DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		var strict *toml.StrictMissingError
		if errors.As(err, &strict) {
			return errors.New(strict.String())
		}
		return err
	}
	return nil
}

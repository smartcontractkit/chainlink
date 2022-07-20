package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

// Note that adding an entropy requirement wouldn't add much, since a 16
// character password already has an entropy score of 75 even if it's all
// lowercase characters
const PasswordComplexityRequirements = `
Must have a length of 16-50 characters
Must not comprise:
	Leading or trailing whitespace (note that a trailing newline in the password file, if present, will be ignored)
	A user's API email
`

const MinRequiredLen = 16

var LeadingWhitespace = regexp.MustCompile(`^\s+`)
var TrailingWhitespace = regexp.MustCompile(`\s+$`)

var (
	ErrPasswordMinLength = errors.Errorf("must be longer than %d characters", MinRequiredLen)
	ErrWhitespace        = errors.New("must not contain leading or trailing whitespace characters")
)

func VerifyPasswordComplexity(password string, disallowedStrings ...string) (merr error) {
	if LeadingWhitespace.MatchString(password) || TrailingWhitespace.MatchString(password) {
		merr = multierr.Append(merr, ErrWhitespace)
	}

	if len(password) < MinRequiredLen {
		merr = multierr.Append(merr, ErrPasswordMinLength)
	}

	for _, s := range disallowedStrings {
		if strings.Contains(strings.ToLower(password), strings.ToLower(s)) {
			merr = multierr.Append(merr, errors.Errorf("password may not contain: %q", s))
		}
	}

	if merr != nil {
		merr = fmt.Errorf("password does not meet the requirements: %s", merr.Error())
	}

	return
}

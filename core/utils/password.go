package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// PasswordComplexityRequirements defines the complexity requirements message
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
	ErrMsgHeader = fmt.Sprintf(`
Expected password complexity:
Must be longer than %d characters
Must not comprise:
	Leading or trailing whitespace
	A user's API email

Faults:
`, MinRequiredLen)
	ErrWhitespace = errors.New("Password contains a leading or trailing whitespace")
)

func VerifyPasswordComplexity(password string, disallowedStrings ...string) (merr error) {
	errMsg := ErrMsgHeader
	var stringErrs []string

	if LeadingWhitespace.MatchString(password) || TrailingWhitespace.MatchString(password) {
		stringErrs = append(stringErrs, ErrWhitespace.Error())
	}

	if len(password) < MinRequiredLen {
		stringErrs = append(stringErrs, fmt.Sprintf("Password is %d characters long", len(password)))
	}

	for _, s := range disallowedStrings {
		if strings.Contains(strings.ToLower(password), strings.ToLower(s)) {
			stringErrs = append(stringErrs, fmt.Sprintf("Password may not contain: %q", s))
		}
	}

	if len(stringErrs) > 0 {
		for _, stringErr := range stringErrs {
			errMsg = fmt.Sprintf("%s	%s\n", errMsg, stringErr)
		}
		merr = errors.New(errMsg)
	}

	return
}

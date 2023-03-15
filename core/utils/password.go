package utils

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	ErrPasswordWhitespace  = errors.New("leading/trailing whitespace detected in password")
	ErrEmptyPasswordInFile = errors.New("detected empty password in password file")
)

// PasswordComplexityRequirements defines the complexity requirements message
// Note that adding an entropy requirement wouldn't add much, since a 16
// character password already has an entropy score of 75 even if it's all
// lowercase characters
const PasswordComplexityRequirements = `
Must have a length of 16-50 characters
Must not comprise:
	Leading or trailing whitespace (note that a trailing newline in the password file, if present, will be ignored)
`

const MinRequiredLen = 16

var LeadingWhitespace = regexp.MustCompile(`^\s+`)
var TrailingWhitespace = regexp.MustCompile(`\s+$`)

var (
	ErrMsgHeader = fmt.Sprintf(`
Expected password complexity:
Must be at least %d characters long
Must not comprise:
	Leading or trailing whitespace
	A user's API email

Faults:
`, MinRequiredLen)
	ErrWhitespace = errors.New("password contains a leading or trailing whitespace")
)

func VerifyPasswordComplexity(password string, disallowedStrings ...string) (merr error) {
	errMsg := ErrMsgHeader
	var stringErrs []string

	if LeadingWhitespace.MatchString(password) || TrailingWhitespace.MatchString(password) {
		stringErrs = append(stringErrs, ErrWhitespace.Error())
	}

	if len(password) < MinRequiredLen {
		stringErrs = append(stringErrs, fmt.Sprintf("password is less than %d characters long", MinRequiredLen))
	}

	for _, s := range disallowedStrings {
		if strings.Contains(strings.ToLower(password), strings.ToLower(s)) {
			stringErrs = append(stringErrs, fmt.Sprintf("password may not contain: %q", s))
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

func PasswordFromFile(pwdFile string) (string, error) {
	if len(pwdFile) == 0 {
		return "", nil
	}
	dat, err := os.ReadFile(pwdFile)
	// handle POSIX case, when text files may have a trailing \n
	pwd := strings.TrimSuffix(string(dat), "\n")

	if err != nil {
		return pwd, err
	}
	if len(pwd) == 0 {
		return pwd, ErrEmptyPasswordInFile
	}
	if strings.TrimSpace(pwd) != pwd {
		return pwd, ErrPasswordWhitespace
	}
	return pwd, err
}

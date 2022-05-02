package utils

import (
	"fmt"
	"regexp"

	"go.uber.org/multierr"
)

const PasswordComplexityRequirements = `
Must be longer than 12 characters
Must comprise at least 3 of:
	lowercase characters
	uppercase characters
	numbers
	symbols
Must not comprise:
	A user's API email
	More than three identical consecutive characters
`

var (
	lowercase = regexp.MustCompile("[a-z]")
	uppercase = regexp.MustCompile("[A-Z]")
	numbers   = regexp.MustCompile("[0-9]")
	symbols   = regexp.MustCompile(`[!@#$%^&*()-=_+\[\]\\|;:'",<.>/?~` + "`]")
)

func VerifyPasswordComplexity(password string) (merr error) {
	if len(password) <= 12 {
		merr = multierr.Append(merr, fmt.Errorf("must be longer than 12 characters"))
	}
	if len(lowercase.FindAllString(password, -1)) < 3 {
		merr = multierr.Append(merr, fmt.Errorf("must contain at least 3 lowercase characters"))
	}
	if len(uppercase.FindAllString(password, -1)) < 3 {
		merr = multierr.Append(merr, fmt.Errorf("must contain at least 3 uppercase characters"))
	}
	if len(numbers.FindAllString(password, -1)) < 3 {
		merr = multierr.Append(merr, fmt.Errorf("must contain at least 3 numbers"))
	}
	if len(symbols.FindAllString(password, -1)) < 3 {
		merr = multierr.Append(merr, fmt.Errorf("must contain at least 3 symbols"))
	}
	var c byte
	var instances int
	for i := 0; i < len(password); i++ {
		if password[i] == c {
			instances++
		} else {
			instances = 1
		}
		if instances > 3 {
			merr = multierr.Append(merr, fmt.Errorf("must not contain more than 3 identical consecutive characters"))
			break
		}
		c = password[i]
	}

	if merr != nil {
		merr = fmt.Errorf("password does not meet the requirements.\n%+v", merr)
	}

	return
}

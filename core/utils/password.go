package utils

import (
	"errors"
	"fmt"
	"regexp"
	"unicode"

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

	ErrPasswordMinLength     = errors.New("must be longer than 12 characters")
	ErrPasswordMinLowercase  = errors.New("must contain at least 3 lowercase characters")
	ErrPasswordMinUppercase  = errors.New("must contain at least 3 uppercase characters")
	ErrPasswordMinNumbers    = errors.New("must contain at least 3 numbers")
	ErrPasswordMinSymbols    = errors.New("must contain at least 3 symbols")
	ErrPasswordRepeatedChars = errors.New("must not contain more than 3 identical consecutive characters")
)

func countSymbols(password string) (count int) {
	for _, r := range password {
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			count++
		}
	}
	return
}

func VerifyPasswordComplexity(password string) (merr error) {
	if len(password) <= 12 {
		merr = multierr.Append(merr, ErrPasswordMinLength)
	}
	if len(lowercase.FindAllString(password, -1)) < 3 {
		merr = multierr.Append(merr, ErrPasswordMinLowercase)
	}
	if len(uppercase.FindAllString(password, -1)) < 3 {
		merr = multierr.Append(merr, ErrPasswordMinUppercase)
	}
	if len(numbers.FindAllString(password, -1)) < 3 {
		merr = multierr.Append(merr, ErrPasswordMinNumbers)
	}
	if countSymbols(password) < 3 {
		merr = multierr.Append(merr, ErrPasswordMinSymbols)
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
			merr = multierr.Append(merr, ErrPasswordRepeatedChars)
			break
		}
		c = password[i]
	}

	if merr != nil {
		merr = fmt.Errorf("password does not meet the requirements.\n%+v", merr)
	}

	return
}

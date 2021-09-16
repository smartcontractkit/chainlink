package cmd

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
)

// TerminalKeyStoreAuthenticator contains fields for prompting the user and an
// exit code.
type TerminalKeyStoreAuthenticator struct {
	Prompter Prompter
}

func (auth TerminalKeyStoreAuthenticator) authenticate(c *clipkg.Context, keyStore keystore.Master) error {
	isEmpty, err := keyStore.IsEmpty()
	if err != nil {
		return errors.Wrap(err, "error determining if keystore is empty")
	}
	password, err := passwordFromFile(c.String("password"))
	if err != nil {
		return errors.Wrap(err, "error reading password from file")
	}
	passwordProvided := len(password) != 0
	if passwordProvided {
		return keyStore.Unlock(password)
	}
	interactive := auth.Prompter.IsTerminal()
	if !interactive {
		return errors.New("no password provided")
	} else if !isEmpty {
		password = auth.promptExistingPassword()
	} else {
		password, err = auth.promptNewPassword()
	}
	if err != nil {
		return err
	}
	return keyStore.Unlock(password)
}

func (auth TerminalKeyStoreAuthenticator) validatePasswordStrength(password string) error {
	// Password policy:
	//
	// Must be longer than 12 characters
	// Must comprise at least 3 of:
	//     lowercase characters
	//     uppercase characters
	//     numbers
	//     symbols
	// Must not comprise:
	//     A user's API email
	//     More than three identical consecutive characters

	var (
		lowercase = regexp.MustCompile("[a-z]")
		uppercase = regexp.MustCompile("[A-Z]")
		numbers   = regexp.MustCompile("[0-9]")
		symbols   = regexp.MustCompile(`[!@#$%^&*()-=_+\[\]\\|;:'",<.>/?~` + "`]")
	)

	var merr error
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
	return merr
}

func (auth TerminalKeyStoreAuthenticator) promptExistingPassword() string {
	password := auth.Prompter.PasswordPrompt("Enter key store password:")
	return password
}

func (auth TerminalKeyStoreAuthenticator) promptNewPassword() (string, error) {
	for {
		password := auth.Prompter.PasswordPrompt("New key store password: ")
		err := auth.validatePasswordStrength(password)
		if err != nil {
			return password, err
		}
		clearLine()
		passwordConfirmation := auth.Prompter.PasswordPrompt("Confirm password: ")
		clearLine()
		if password != passwordConfirmation {
			fmt.Printf("Passwords don't match. Please try again... ")
			continue
		}
		return password, nil
	}
}

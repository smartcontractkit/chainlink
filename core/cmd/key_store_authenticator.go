package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	clipkg "github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"
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
	return utils.VerifyPasswordComplexity(password)
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

package cmd

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// TerminalKeyStoreAuthenticator contains fields for prompting the user and an
// exit code.
type TerminalKeyStoreAuthenticator struct {
	Prompter Prompter
}

func (auth TerminalKeyStoreAuthenticator) authenticate(keyStore keystore.Master, cfg config.BasicConfig) error {
	isEmpty, err := keyStore.IsEmpty()
	if err != nil {
		return errors.Wrap(err, "error determining if keystore is empty")
	}
	password := cfg.KeystorePassword()

	if len(password) != 0 {
		// Because we changed password requirements to increase complexity, to
		// not break backward compatibility we enforce this only for empty key
		// stores.
		if err = auth.validatePasswordStrength(password); err != nil && isEmpty {
			return err
		}
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
		if err := auth.validatePasswordStrength(password); err != nil {
			return "", err
		}
		if strings.TrimSpace(password) != password {
			return "", utils.ErrPasswordWhitespace
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

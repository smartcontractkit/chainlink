package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// TerminalKeyStoreAuthenticator contains fields for prompting the user and an
// exit code.
type TerminalKeyStoreAuthenticator struct {
	Prompter Prompter
}

type keystorePassword interface {
	Keystore() string
}

func (auth TerminalKeyStoreAuthenticator) authenticate(ctx context.Context, keyStore keystore.Master, password keystorePassword) error {
	isEmpty, err := keyStore.IsEmpty(ctx)
	if err != nil {
		return errors.Wrap(err, "error determining if keystore is empty")
	}
	pw := password.Keystore()

	if len(pw) != 0 {
		// Because we changed password requirements to increase complexity, to
		// not break backward compatibility we enforce this only for empty key
		// stores.
		if err = auth.validatePasswordStrength(pw); err != nil && isEmpty {
			return err
		}
		return keyStore.Unlock(ctx, pw)
	}
	interactive := auth.Prompter.IsTerminal()
	if !interactive {
		return errors.New("no password provided")
	} else if !isEmpty {
		pw = auth.promptExistingPassword()
	} else {
		pw, err = auth.promptNewPassword()
	}
	if err != nil {
		return err
	}
	return keyStore.Unlock(ctx, pw)
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

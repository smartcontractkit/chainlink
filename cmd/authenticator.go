package cmd

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/store"
)

// Authenticator implements the Authenticate method for the store and
// a password string.
type Authenticator interface {
	Authenticate(*store.Store, string) error
}

// TerminalAuthenticator contains fields for prompting the user and an
// exit code.
type TerminalAuthenticator struct {
	Prompter Prompter
}

// Authenticate checks to see if there are accounts present in
// the KeyStore, and if there are none, a new account will be created
// by prompting for a password. If there are accounts present, the
// account which is unlocked by the given password will be used.
func (auth TerminalAuthenticator) Authenticate(store *store.Store, pwd string) error {
	if len(pwd) != 0 {
		return auth.authenticateWithPwd(store, pwd)
	} else if auth.Prompter.IsTerminal() {
		return auth.authenticationPrompt(store)
	} else {
		return errors.New("No password provided")
	}
}

func (auth TerminalAuthenticator) authenticationPrompt(store *store.Store) error {
	if store.KeyStore.HasAccounts() {
		return auth.promptAndCheckPasswordLoop(store)
	}
	return auth.promptAndCreateAccount(store)
}

func (auth TerminalAuthenticator) authenticateWithPwd(store *store.Store, pwd string) error {
	if !store.KeyStore.HasAccounts() {
		fmt.Println("There are no accounts, creating a new account with the specified password")
		return createAccount(store, pwd)
	}
	return checkPassword(store, pwd)
}

func checkPassword(store *store.Store, phrase string) error {
	if err := store.KeyStore.Unlock(phrase); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (auth TerminalAuthenticator) promptAndCheckPasswordLoop(store *store.Store) error {
	for {
		phrase := auth.Prompter.PasswordPrompt("Enter Password:")
		if checkPassword(store, phrase) == nil {
			break
		}
	}

	return nil
}

func (auth TerminalAuthenticator) promptAndCreateAccount(store *store.Store) error {
	for {
		phrase := auth.Prompter.PasswordPrompt("New Password: ")
		clearLine()
		phraseConfirmation := auth.Prompter.PasswordPrompt("Confirm Password: ")
		clearLine()
		if phrase == phraseConfirmation {
			return createAccount(store, phrase)
		}
		fmt.Printf("Passwords don't match. Please try again... ")
	}
}

func createAccount(store *store.Store, password string) error {
	_, err := store.KeyStore.NewAccount(password)
	return err
}

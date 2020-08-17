package cmd

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// KeyStoreAuthenticator implements the Authenticate method for the store and
// a password string.
type KeyStoreAuthenticator interface {
	Authenticate(*store.Store, string) (string, error)
	AuthenticateVRFKey(*store.Store, string) error
}

// TerminalKeyStoreAuthenticator contains fields for prompting the user and an
// exit code.
type TerminalKeyStoreAuthenticator struct {
	Prompter Prompter
}

// Authenticate checks to see if there are accounts present in
// the KeyStore, and if there are none, a new account will be created
// by prompting for a password. If there are accounts present, the
// account which is unlocked by the given password will be used.
func (auth TerminalKeyStoreAuthenticator) Authenticate(store *store.Store, pwd string) (string, error) {
	if len(pwd) != 0 {
		return auth.authenticateWithPwd(store, pwd)
	} else if auth.Prompter.IsTerminal() {
		return auth.authenticationPrompt(store)
	} else {
		return "", errors.New("No password provided")
	}
}

func (auth TerminalKeyStoreAuthenticator) authenticationPrompt(store *store.Store) (string, error) {
	if store.KeyStore.HasAccounts() {
		return auth.promptAndCheckPasswordLoop(store), nil
	}
	return auth.promptAndCreateAccount(store)
}

func (auth TerminalKeyStoreAuthenticator) authenticateWithPwd(store *store.Store, pwd string) (string, error) {
	if !store.KeyStore.HasAccounts() {
		fmt.Println("There are no accounts, creating a new account with the specified password")
		return pwd, createAccount(store, pwd)
	}
	return pwd, checkPassword(store, pwd)
}

func checkPassword(store *store.Store, phrase string) error {
	return store.KeyStore.Unlock(phrase)
}

func (auth TerminalKeyStoreAuthenticator) promptAndCheckPasswordLoop(store *store.Store) string {
	for {
		phrase := auth.Prompter.PasswordPrompt("Enter Password:")
		if checkPassword(store, phrase) == nil {
			return phrase
		}
	}
}

func (auth TerminalKeyStoreAuthenticator) promptAndCreateAccount(store *store.Store) (string, error) {
	for {
		phrase := auth.Prompter.PasswordPrompt("New Password: ")
		clearLine()
		phraseConfirmation := auth.Prompter.PasswordPrompt("Confirm Password: ")
		clearLine()
		if phrase == phraseConfirmation {
			return phrase, createAccount(store, phrase)
		}
		fmt.Printf("Passwords don't match. Please try again... ")
	}
}

func createAccount(store *store.Store, password string) error {
	_, err := store.KeyStore.NewAccount(password)
	if err != nil {
		return errors.Wrapf(err, "while creating ethereum keys")
	}
	return checkPassword(store, password)
}

// AuthenticateVRFKey creates an encrypted VRF key protected by password in
// store's db if db store has no extant keys. It unlocks at least one VRF key
// with given password, or returns an error. password must be non-trivial, as an
// empty password signifies that the VRF oracle functionality is disabled.
func (auth TerminalKeyStoreAuthenticator) AuthenticateVRFKey(store *store.Store,
	password string) error {
	if password == "" {
		return fmt.Errorf("VRF password must be non-trivial")
	}
	keys, err := store.VRFKeyStore.Get()
	if err != nil {
		return errors.Wrapf(err, "while checking for extant VRF keys")
	}
	if len(keys) == 0 {
		fmt.Println(
			"There are no VRF keys; creating a new key encrypted with given password")
		if _, err := store.VRFKeyStore.CreateKey(password); err != nil {
			return errors.Wrapf(err, "while creating a new encrypted VRF key")
		}
	}
	return errors.Wrapf(utils.JustError(store.VRFKeyStore.Unlock(password)),
		"there are VRF keys in the DB, but that password did not unlock any of "+
			"them... please check the password in the file specified by vrfpassword"+
			". You can add and delete VRF keys in the DB using the "+
			"`chainlink local vrf` subcommands")
}

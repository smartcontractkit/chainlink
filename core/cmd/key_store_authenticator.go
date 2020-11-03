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
	AuthenticateOCRKey(store *store.Store, password string) error
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
func (auth TerminalKeyStoreAuthenticator) Authenticate(store *store.Store, password string) (string, error) {
	passwordProvided := len(password) != 0
	interactive := auth.Prompter.IsTerminal()
	hasAccounts := store.KeyStore.HasAccounts()

	if passwordProvided && hasAccounts {
		return auth.unlockExistingWithPassword(store, password)
	} else if passwordProvided && !hasAccounts {
		return auth.unlockNewWithPassword(store, password)
	} else if !passwordProvided && interactive && hasAccounts {
		return auth.promptExistingPassword(store)
	} else if !passwordProvided && interactive && !hasAccounts {
		return auth.promptNewPassword(store)
	} else {
		return "", errors.New("No password provided")
	}
}

func (auth TerminalKeyStoreAuthenticator) promptExistingPassword(store *store.Store) (string, error) {
	for {
		password := auth.Prompter.PasswordPrompt("Enter Password:")
		if store.KeyStore.Unlock(password) == nil {
			return password, nil
		}
	}
}

func (auth TerminalKeyStoreAuthenticator) promptNewPassword(store *store.Store) (string, error) {
	for {
		password := auth.Prompter.PasswordPrompt("New Password: ")
		clearLine()
		passwordConfirmation := auth.Prompter.PasswordPrompt("Confirm Password: ")
		clearLine()
		if password == passwordConfirmation {
			fmt.Printf("Passwords don't match. Please try again... ")
			continue
		}
		_, err := store.KeyStore.NewAccount()
		return password, errors.Wrapf(err, "while creating ethereum keys")
	}
}

func (auth TerminalKeyStoreAuthenticator) unlockNewWithPassword(store *store.Store, password string) (string, error) {
	err := store.KeyStore.Unlock(password)
	if err != nil {
		return "", errors.Wrap(err, "Error unlocking key store")
	}
	fmt.Println("There are no accounts, creating a new account with the specified password")
	_, err = store.KeyStore.NewAccount()
	return password, errors.Wrapf(err, "while creating ethereum keys")
}

func (auth TerminalKeyStoreAuthenticator) unlockExistingWithPassword(store *store.Store, password string) (string, error) {
	err := store.KeyStore.Unlock(password)
	return password, err
}

// AuthenticateVRFKey creates an encrypted VRF key protected by password in
// store's db if db store has no extant keys. It unlocks at least one VRF key
// with given password, or returns an error. password must be non-trivial, as an
// empty password signifies that the VRF oracle functionality is disabled.
func (auth TerminalKeyStoreAuthenticator) AuthenticateVRFKey(store *store.Store, password string) error {
	if password == "" {
		return fmt.Errorf("VRF password must be non-trivial")
	}
	keys, err := store.VRFKeyStore.Get()
	if err != nil {
		return errors.Wrapf(err, "while checking for extant VRF keys")
	}
	if len(keys) == 0 {
		fmt.Println("There are no VRF keys; creating a new key encrypted with given password")
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

func (auth TerminalKeyStoreAuthenticator) AuthenticateOCRKey(store *store.Store, password string) error {
	if password == "" {
		return fmt.Errorf("OCR password must be non-trivial")
	}

	err := store.OCRKeyStore.Unlock(password)
	if err != nil {
		return errors.Wrapf(err,
			"there are OCR/P2P keys in the DB, but there were errors unlocking "+
				"them... please check the password in the file specified by --password"+
				". You can add and delete OCR/P2P keys in the DB using the "+
				"`chainlink node ocr` and `chainlink node p2p` subcommands")
	}

	p2pkeys, err := store.OCRKeyStore.FindEncryptedP2PKeys()
	if err != nil {
		return errors.Wrap(err, "could not fetch encrypted P2P keys from database")
	}
	if len(p2pkeys) == 0 {
		fmt.Println("There are no P2P keys; creating a new key encrypted with given password")
		_, _, err = store.OCRKeyStore.GenerateEncryptedP2PKey()
		if err != nil {
			return errors.Wrapf(err, "while creating a new encrypted P2P key")
		}
	}

	ocrkeys, err := store.OCRKeyStore.FindEncryptedOCRKeyBundles()
	if err != nil {
		return errors.Wrap(err, "could not fetch encrypted OCR keys from database")
	}
	if len(ocrkeys) == 0 {
		fmt.Println("There are no OCR keys; creating a new key encrypted with given password")
		_, _, err := store.OCRKeyStore.GenerateEncryptedOCRKeyBundle()
		if err != nil {
			return errors.Wrapf(err, "while creating a new encrypted OCR key")
		}
	}
	return nil
}

package cmd

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

// KeyStoreAuthenticator implements the Authenticate method for the store and
// a password string.
type KeyStoreAuthenticator interface {
	AuthenticateEthKey(*keystore.Eth, string) (string, error)
	AuthenticateCSAKey(*keystore.CSA, string) error
	AuthenticateVRFKey(*keystore.VRF, string) error
	AuthenticateOCRKey(*keystore.OCR, *orm.Config, string) error
}

// TerminalKeyStoreAuthenticator contains fields for prompting the user and an
// exit code.
type TerminalKeyStoreAuthenticator struct {
	Prompter Prompter
}

// Authenticate checks to see if there are accounts present in
// the KeyStore, and if there are none, a new account will be created
// by prompting for a password. If there are accounts present, all accounts
// will be unlocked.
func (auth TerminalKeyStoreAuthenticator) AuthenticateEthKey(ethKeyStore *keystore.Eth, password string) (string, error) {
	passwordProvided := len(password) != 0
	interactive := auth.Prompter.IsTerminal()
	hasSendingKeys, err := ethKeyStore.HasDBSendingKeys()
	if err != nil {
		return "", errors.Wrap(err, "failed to query DB for send keys")
	}

	if passwordProvided && hasSendingKeys {
		return auth.unlockExistingWithPassword(ethKeyStore, password)
	} else if passwordProvided && !hasSendingKeys {
		return auth.unlockNewWithPassword(ethKeyStore, password)
	} else if !passwordProvided && interactive && hasSendingKeys {
		return auth.promptExistingPassword(ethKeyStore)
	} else if !passwordProvided && interactive && !hasSendingKeys {
		return auth.promptNewPassword(ethKeyStore)
	} else {
		return "", errors.New("No password provided")
	}
}

func (auth TerminalKeyStoreAuthenticator) validatePasswordStrength(ethKeyStore *keystore.Eth, password string) error {
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

func (auth TerminalKeyStoreAuthenticator) promptExistingPassword(ethKeyStore *keystore.Eth) (string, error) {
	for {
		password := auth.Prompter.PasswordPrompt("Enter key store password:")
		if ethKeyStore.Unlock(password) == nil {
			return password, nil
		}
	}
}

func (auth TerminalKeyStoreAuthenticator) promptNewPassword(ethKeyStore *keystore.Eth) (string, error) {
	for {
		password := auth.Prompter.PasswordPrompt("New key store password: ")
		err := auth.validatePasswordStrength(ethKeyStore, password)
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
		err = ethKeyStore.Unlock(password)
		if err != nil {
			return password, errors.Wrap(err, "unexpectedly failed to unlock KeyStore")
		}
		_, err = ethKeyStore.CreateNewKey()
		return password, errors.Wrap(err, "failed to create new ETH key")
	}
}

func (auth TerminalKeyStoreAuthenticator) unlockNewWithPassword(ethKeyStore *keystore.Eth, password string) (string, error) {
	err := auth.validatePasswordStrength(ethKeyStore, password)
	if err != nil {
		return password, err
	}
	err = ethKeyStore.Unlock(password)
	if err != nil {
		return "", errors.Wrap(err, "Error unlocking key store")
	}
	fmt.Println("There are no accounts, creating a new account with the specified password")
	_, err = ethKeyStore.CreateNewKey()
	return password, errors.Wrap(err, "failed to create new ETH key")
}

func (auth TerminalKeyStoreAuthenticator) unlockExistingWithPassword(ethKeyStore *keystore.Eth, password string) (string, error) {
	err := ethKeyStore.Unlock(password)
	return password, err
}

// AuthenticateVRFKey creates an encrypted VRF key protected by password in
// store's db if db store has no extant keys. It unlocks at least one VRF key
// with given password, or returns an error. password must be non-trivial, as an
// empty password signifies that the VRF oracle functionality is disabled.
func (auth TerminalKeyStoreAuthenticator) AuthenticateVRFKey(vrfKeyStore *keystore.VRF, password string) error {
	keys, err := vrfKeyStore.Unlock(password) // Saves the password
	if err != nil {
		return errors.Wrapf(err,
			"there are VRF keys in the DB, but there were errors unlocking "+
				"them... please check the password in the file specified by --vrfpassword"+
				". You can add and delete VRF keys in the DB using the "+
				"`chainlink keys vrf` subcommand")
	}
	if len(keys) == 0 {
		fmt.Println("There are no VRF keys; creating a new key encrypted with given password")
		if _, err2 := vrfKeyStore.CreateKey(); err2 != nil {
			return errors.Wrapf(err2, "while creating a new encrypted VRF key")
		}
	}
	// Double check we can read back the encrypted key
	_, err = vrfKeyStore.Get()
	if err != nil {
		return errors.Wrap(err, "could not fetch encrypted VRF keys from database")
	}
	return nil
}

func (auth TerminalKeyStoreAuthenticator) AuthenticateCSAKey(csaKeyStore *keystore.CSA, password string) error {
	return errors.Wrapf(csaKeyStore.Unlock(password),
		"there are CSA keys in the DB, but that password did not unlock any of "+
			"them... please check the password in the file")
}

// AuthenticateOCRKey authenticates OCR keypairs
func (auth TerminalKeyStoreAuthenticator) AuthenticateOCRKey(ocrKeyStore *keystore.OCR, config *orm.Config, password string) error {
	err := ocrKeyStore.Unlock(password)
	if err != nil {
		return errors.Wrapf(err,
			"there are OCR/P2P keys in the DB, but there were errors unlocking "+
				"them... please check the password in the file specified by --password"+
				". You can add and delete OCR/P2P keys in the DB using the "+
				"`chainlink node ocr` and `chainlink node p2p` subcommands")
	}

	p2pkeys, err := ocrKeyStore.FindEncryptedP2PKeys()
	if err != nil {
		return errors.Wrap(err, "could not fetch encrypted P2P keys from database")
	}
	if len(p2pkeys) == 0 {
		fmt.Println("There are no P2P keys; creating a new key encrypted with given password")
		var k p2pkey.EncryptedP2PKey
		_, k, err = ocrKeyStore.GenerateEncryptedP2PKey()
		if err != nil {
			return errors.Wrapf(err, "while creating a new encrypted P2P key")
		}
		if !config.P2PPeerIDIsSet() {
			config.Set("P2P_PEER_ID", k.PeerID)
		}
	}

	ocrkeys, err := ocrKeyStore.FindEncryptedOCRKeyBundles()
	if err != nil {
		return errors.Wrap(err, "could not fetch encrypted OCR keys from database")
	}
	if len(ocrkeys) == 0 {
		fmt.Println("There are no OCR keys; creating a new key encrypted with given password")
		_, _, err := ocrKeyStore.GenerateEncryptedOCRKeyBundle()
		if err != nil {
			return errors.Wrapf(err, "while creating a new encrypted OCR key")
		}
	}
	return nil
}

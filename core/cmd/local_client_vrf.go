package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	clipkg "github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models/vrfkey"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func vRFKeyStore(cli *Client) *store.VRFKeyStore {
	return cli.AppFactory.NewApplication(cli.Config).GetStore().VRFKeyStore
}

// CreateVRFKey creates a key in the VRF keystore, protected by the password in
// the password file
func (cli *Client) CreateVRFKey(c *clipkg.Context) error {
	cli.Config.Dialect = orm.DialectPostgresWithoutLock
	password, err := getPassword(c)
	if err != nil {
		return err
	}
	key, err := vRFKeyStore(cli).CreateKey(string(password))
	if err != nil {
		return errors.Wrapf(err, "while creating new account")
	}
	uncompressedKey, err := key.StringUncompressed()
	if err != nil {
		return errors.Wrapf(err, "while creating new account")
	}
	hash, err := key.Hash()
	hashStr := hash.Hex()
	if err != nil {
		hashStr = "error while computing hash of public key: " + err.Error()
	}
	fmt.Printf(`Created keypair.

Compressed public key (use this for interactions with the chainlink node):
  %s
Uncompressed public key (use this to register key with the VRFCoordinator):
  %s
Hash of public key (use this to request randomness from your consuming contract):
  %s

The following command will export the encrypted secret key from the db to <save_path>:

chainlink local vrf export -f <save_path> -pk %s
`, key, uncompressedKey, hashStr, key)
	return nil
}

// CreateAndExportWeakVRFKey creates a key in the VRF keystore, protected by the
// password in the password file, but with weak key-derivation-function
// parameters, which makes it cheaper for testing, but also more vulnerable to
// bruteforcing of the encrypted key material. For testing purposes only!
//
// The key is only stored at the specified file location, not stored in the DB.
func (cli *Client) CreateAndExportWeakVRFKey(c *clipkg.Context) error {
	password, err := getPassword(c)
	if err != nil {
		return err
	}
	key, err := vRFKeyStore(cli).CreateWeakInMemoryEncryptedKeyXXXTestingOnly(
		string(password))
	if err != nil {
		return errors.Wrapf(err, "while creating testing key")
	}
	if !c.IsSet("file") || !noFileToOverwrite(c.String("file")) {
		errmsg := "must specify path to key file which does not already exist"
		fmt.Println(errmsg)
		return fmt.Errorf(errmsg)
	}
	fmt.Println("Don't use this key for anything sensitive!")
	return key.WriteToDisk(c.String("file"))
}

// getPassword retrieves the password from the file specified on the CL, or errors
func getPassword(c *clipkg.Context) ([]byte, error) {
	if !c.IsSet("password") {
		return nil, fmt.Errorf("must specify password file")
	}
	rawPassword, err := passwordFromFile(c.String("password"))
	if err != nil {
		return nil, errors.Wrapf(err, "could not read password from file %s",
			c.String("password"))
	}
	return []byte(rawPassword), nil
}

// getPasswordAndKeyFile retrieves the password and key json from the files
// specified on the CL, or errors
func getPasswordAndKeyFile(c *clipkg.Context) (password []byte, keyjson []byte, err error) {
	password, err = getPassword(c)
	if err != nil {
		return nil, nil, err
	}
	if !c.IsSet("file") {
		return nil, nil, fmt.Errorf("must specify key file")
	}
	keypath := c.String("file")
	keyjson, err = ioutil.ReadFile(keypath)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to read file %s", keypath)
	}
	return password, keyjson, nil
}

// ImportVRFKey reads a file into an EncryptedSecretKey in the db
func (cli *Client) ImportVRFKey(c *clipkg.Context) error {
	password, keyjson, err := getPasswordAndKeyFile(c)
	if err != nil {
		return err
	}
	if err := vRFKeyStore(cli).Import(keyjson, string(password)); err != nil {
		if err == store.MatchingVRFKeyError {
			fmt.Println(`The database already has an entry for that public key.`)
			var key struct{ PublicKey string }
			if e := json.Unmarshal(keyjson, &key); e != nil {
				fmt.Println("could not extract public key from json input")
				return errors.Wrapf(e, "while extracting public key from %s", keyjson)
			}
			fmt.Printf(`If you want to import the new key anyway, delete the old key with the command

    %s

(but maybe back it up first, with %s.)
`,
				fmt.Sprintf("chainlink local delete -pk %s", key.PublicKey),
				fmt.Sprintf("`chainlink local export -f <backup_path> -pk %s`",
					key.PublicKey))
			return errors.Wrap(err, "while attempting to import key from CL")
		}
		return err
	}
	return nil
}

// ExportVRFKey saves encrypted copy of VRF key with given public key to
// requested file path.
func (cli *Client) ExportVRFKey(c *clipkg.Context) error {
	encryptedKey, err := getKeys(cli, c)
	if err != nil {
		return err
	}
	if !c.IsSet("file") {
		return fmt.Errorf("must specify file to export to") // Or could default to stdout?
	}
	keypath := c.String("file")
	_, err = os.Stat(keypath)
	if err == nil {
		return fmt.Errorf(
			"refusing to overwrite existing file %s. Please move it or change the save path",
			keypath)
	}
	if !os.IsNotExist(err) {
		return errors.Wrapf(err, "while checking whether file %s exists", keypath)
	}
	if err := encryptedKey.WriteToDisk(keypath); err != nil {
		return errors.Wrapf(err, "could not save %#+v to %s", encryptedKey, keypath)
	}
	return nil
}

// getKeys retrieves the keys for an ExportVRFKey request
func getKeys(cli *Client, c *clipkg.Context) (*vrfkey.EncryptedSecretKey, error) {
	publicKey, err := getPublicKey(c)
	if err != nil {
		return nil, err
	}
	enckey, err := vRFKeyStore(cli).GetSpecificKey(publicKey)
	if err != nil {
		return nil, errors.Wrapf(err,
			"while retrieving keys with matching public key %s", publicKey.String())
	}
	return enckey, nil
}

// DeleteVRFKey deletes the VRF key with given public key from the db
//
// Since this runs in an independent process from any chainlink node, it cannot
// cause running nodes to forget the key, if they already have it unlocked.
func (cli *Client) DeleteVRFKey(c *clipkg.Context) error {
	publicKey, err := getPublicKey(c)
	if err != nil {
		return err
	}
	if err := vRFKeyStore(cli).Delete(publicKey); err != nil {
		if err == store.AttemptToDeleteNonExistentKeyFromDB {
			fmt.Printf("There is already no entry in the DB for %s\n", publicKey)
		}
		return err
	}
	return nil
}

func getPublicKey(c *clipkg.Context) (vrfkey.PublicKey, error) {
	if !c.IsSet("publicKey") {
		return vrfkey.PublicKey{}, fmt.Errorf("must specify public key")
	}
	publicKey, err := vrfkey.NewPublicKeyFromHex(c.String("publicKey"))
	if err != nil {
		return vrfkey.PublicKey{}, errors.Wrap(err, "failed to parse public key")
	}
	return publicKey, nil
}

// ListKeys Lists the keys in the db
func (cli *Client) ListKeys(c *clipkg.Context) error {
	keys, err := vRFKeyStore(cli).ListKeys()
	if err != nil {
		return err
	}
	// TODO(alx) Figure out how to make a nice box out of this, like the other
	// commands do.
	fmt.Println(
		`*********************************************************************************
Public keys of encrypted keys in database (compressed, uncompressed, hash)
*********************************************************************************`)
	for keyidx, key := range keys {
		fmt.Println("compressed  ", key)
		uncompressed, err := key.StringUncompressed()
		if err != nil {
			logger.Infow("keys",
				fmt.Sprintf("while computing uncompressed representation of %+v: %s",
					key, err))
			uncompressed = "error while computing uncompressed representation: " +
				err.Error()
		}
		fmt.Println("uncompressed", uncompressed)
		hash, err := key.Hash()
		if err != nil {
			logger.Infow("keys", "while computing hash of %+v: %s", key, hash)
			fmt.Println("hash        ", "error while computing hash of %+v: "+err.Error())
		} else {
			fmt.Println("hash        ", hash.Hex())
		}
		if keyidx != len(keys)-1 {
			fmt.Println(
				"---------------------------------------------------------------------------------")
		}
	}
	fmt.Println(
		"*********************************************************************************")
	logger.Infow("keys", "keys", keys)
	return nil
}

func noFileToOverwrite(path string) bool {
	return os.IsNotExist(utils.JustError(os.Stat(path)))
}

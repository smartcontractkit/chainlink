package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	clipkg "github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/store/models/ocrkey"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

// CreateOCRKey creates a key and inserts it into encrypted_ocr_keys,
// protected by the password in the password file
func (cli *Client) CreateOCRKey(c *clipkg.Context) error {
	return cli.errorOut(cli.createOCRKey(c))
}

// DeleteOCRKey creates a key and inserts it into encrypted_ocr_keys,
// protected by the password in the password file
func (cli *Client) DeleteOCRKey(c *clipkg.Context) error {
	return cli.errorOut(cli.deleteOCRKey(c))

}

// ListOCRKeys creates a key and inserts it into encrypted_ocr_keys,
// protected by the password in the password file
func (cli *Client) ListOCRKeys(c *clipkg.Context) error {
	return cli.errorOut(cli.listOCRKeys(c))
}

const createOCRKeyMsg = `Created OCR keypair.
Key Set ID:
  %s
On-chain Public Address:
  0x%s
Off-chain Public Key:
  %s
`

func (cli *Client) createOCRKey(c *clipkg.Context) error {
	cli.Config.Dialect = orm.DialectPostgresWithoutLock
	store := cli.AppFactory.NewApplication(cli.Config).GetStore()
	password, err := getPassword(c)
	if err != nil {
		return err
	}
	key, err := ocrkey.NewOCRPrivateKey()
	if err != nil {
		return errors.Wrapf(err, "while generating the new OCR key")
	}
	encryptedKey, err := key.Encrypt(string(password), ocrkey.DefaultScryptParams)
	if err != nil {
		return errors.Wrapf(err, "while encrypting the new OCR key")
	}
	err = store.CreateEncryptedOCRKey(encryptedKey)
	if err != nil {
		return errors.Wrapf(err, "while persisting the new encrypted OCR key")
	}
	addressOnChain := key.PublicKeyAddressOnChain()
	fmt.Printf(
		createOCRKeyMsg,
		key.ID,
		hex.EncodeToString(addressOnChain[:]),
		hex.EncodeToString(key.PublicKeyOffChain()),
	)
	return nil
}

func (cli *Client) listOCRKeys(c *clipkg.Context) error {
	cli.Config.Dialect = orm.DialectPostgresWithoutLock
	store := cli.AppFactory.NewApplication(cli.Config).GetStore()
	keys, err := store.FindEncryptedOCRKeys()
	if err != nil {
		return errors.Wrapf(err, "while fetching encrypted OCR keys")
	}

	fmt.Println(
		`***********************************************************************************
Encrypted Off-Chain Reporting Keys in DB
***********************************************************************************`)
	for keyidx, key := range keys {
		fmt.Println("ID                ", key.ID)
		fmt.Println("On-chain Address  ", "0x"+hex.EncodeToString(key.OnChainSigningAddress[:]))
		fmt.Println("Off-chain PubKey  ", hex.EncodeToString(key.OffChainPublicKey))
		if keyidx != len(keys)-1 {
			fmt.Println(
				"-----------------------------------------------------------------------------------")
		}
	}
	if len(keys) == 0 {
		fmt.Println("None")
	}
	fmt.Println(
		"***********************************************************************************")

	return nil
}

func (cli *Client) deleteOCRKey(c *clipkg.Context) error {
	if !c.Args().Present() {
		return errors.New("Must pass the ID of the OCR Key bundle to delete")
	}
	id := c.Args().First()

	cli.Config.Dialect = orm.DialectPostgresWithoutLock
	store := cli.AppFactory.NewApplication(cli.Config).GetStore()

	key, err := store.FindEncryptedOCRKeyByID(id)
	if gorm.IsRecordNotFoundError(err) {
		return errors.New("Unable to find the key bundle with the provided ID")
	} else if err != nil {
		return errors.Wrapf(err, "while fetching the OCR key")
	}

	err = store.DeleteEncryptedOCRKey(key)
	if err != nil {
		return errors.Wrapf(err, "while deleting the OCR key")
	}

	fmt.Printf("Successfully deleted OCRKeyBundle %s", key.ID)
	return nil
}

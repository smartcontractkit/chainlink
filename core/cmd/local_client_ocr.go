package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	clipkg "github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

// CreateOCRKeyBundle creates a key and inserts it into encrypted_ocr_keys,
// protected by the password in the password file
func (cli *Client) CreateOCRKeyBundle(c *clipkg.Context) error {
	return cli.errorOut(cli.createOCRKeyBundle(c))
}

// DeleteOCRKeyBundle creates a key and inserts it into encrypted_ocr_keys,
// protected by the password in the password file
func (cli *Client) DeleteOCRKeyBundle(c *clipkg.Context) error {
	return cli.errorOut(cli.deleteOCRKeyBundle(c))

}

// ListOCRKeyBundles lists the available OCR Key Bundles
func (cli *Client) ListOCRKeyBundles(c *clipkg.Context) error {
	return cli.errorOut(cli.listOCRKeyBundles(c))
}

const listKey = `
Key Bundle ID:
  %s
On-chain Public Address:
  0x%s
Off-chain Public Key:
  %s
Config Public Key:
  %s
`

func (cli *Client) createOCRKeyBundle(c *clipkg.Context) error {
	cli.Config.Dialect = orm.DialectPostgresWithoutLock
	store := cli.AppFactory.NewApplication(cli.Config).GetStore()
	password, err := getPassword(c)
	if err != nil {
		return err
	}
	key, _, err := store.OCRKeyStore.GenerateEncryptedOCRKeyBundle(string(password))
	if err != nil {
		return err
	}
	addressOnChain := key.PublicKeyAddressOnChain()
	configPublicKey := key.PublicKeyConfig()
	fmt.Println("Created OCR key bundle:")
	fmt.Printf(
		listKey,
		key.ID,
		hex.EncodeToString(addressOnChain[:]),
		hex.EncodeToString(key.PublicKeyOffChain()),
		hex.EncodeToString(configPublicKey[:]),
	)
	return nil
}

func (cli *Client) listOCRKeyBundles(c *clipkg.Context) error {
	cli.Config.Dialect = orm.DialectPostgresWithoutLock
	store := cli.AppFactory.NewApplication(cli.Config).GetStore()
	keys, err := store.OCRKeyStore.FindEncryptedOCRKeyBundles()
	if err != nil {
		return errors.Wrapf(err, "while fetching encrypted OCR key bundles")
	}

	fmt.Println(
		`***********************************************************************************
Encrypted Off-Chain Reporting Key Bundles in DB
***********************************************************************************`)
	for keyidx, key := range keys {
		fmt.Printf(
			listKey,
			key.ID,
			hex.EncodeToString(key.OnChainSigningAddress[:]),
			hex.EncodeToString(key.OffChainPublicKey[:]),
			hex.EncodeToString(key.ConfigPublicKey[:]),
		)
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

func (cli *Client) deleteOCRKeyBundle(c *clipkg.Context) error {
	if !c.Args().Present() {
		return errors.New("Must pass the ID of the OCR key bundle to delete")
	}
	idStr := c.Args().First()

	cli.Config.Dialect = orm.DialectPostgresWithoutLock
	store := cli.AppFactory.NewApplication(cli.Config).GetStore()

	id, err := models.Sha256HashFromHex(idStr)
	if err != nil {
		return errors.Wrap(err, "while decoding key ID")
	}

	key, err := store.OCRKeyStore.FindEncryptedOCRKeyBundleByID(id)
	if gorm.IsRecordNotFoundError(err) {
		return errors.New("Unable to find the OCR key bundle with the provided ID")
	} else if err != nil {
		return errors.Wrapf(err, "while fetching the OCR key bundle")
	}

	err = store.OCRKeyStore.DeleteEncryptedOCRKeyBundle(&key)
	if err != nil {
		return errors.Wrapf(err, "while deleting the OCR key bundle")
	}

	fmt.Printf("Successfully deleted OCR key bundle %s\n", key.ID)
	return nil
}

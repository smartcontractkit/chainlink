package cmd

import (
	"encoding/hex"
	"fmt"

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
		// uncompressed, err := key.StringUncompressed()
		// if err != nil {
		// 	logger.Infow("keys",
		// 		fmt.Sprintf("while computing uncompressed representation of %+v: %s",
		// 			key, err))
		// 	uncompressed = "error while computing uncompressed representation: " +
		// 		err.Error()
		// }
		// fmt.Println("uncompressed", uncompressed)
		// hash, err := key.Hash()
		// if err != nil {
		// 	logger.Infow("keys", "while computing hash of %+v: %s", key, hash)
		// 	fmt.Println("hash        ", "error while computing hash of %+v: "+err.Error())
		// } else {
		// 	fmt.Println("hash        ", hash.Hex())
		// }
		if keyidx != len(keys)-1 {
			fmt.Println(
				"-----------------------------------------------------------------------------------")
		}
	}
	fmt.Println(
		"***********************************************************************************")

	return nil
}

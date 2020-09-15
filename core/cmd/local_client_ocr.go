package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	clipkg "github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/store/orm"
)

// CreateOCRKey creates a key and inserts it into encrypted_ocr_keys,
// protected by the password in the password file
func (cli *Client) CreateOCRKey(c *clipkg.Context) error {
	return cli.errorOut(cli.CreateOCRKey(c))
}

const createOCRKeyMsg = `Created OCR keypair.
Key ID
  %v
On-chain Public Address:
  0x%s
Off-chain Public Key:
  %s
`

func (cli *Client) CreateOCRKey(c *clipkg.Context) error {
	cli.Config.Dialect = orm.DialectPostgresWithoutLock
	store := cli.AppFactory.NewApplication(cli.Config).GetStore()

	password, err := getPassword(c)
	if err != nil {
		return err
	}
	key, err := store.OCRKeyStore.CreateKey(string(password))
	if err != nil {
		return errors.Wrapf(err, "while generating new OCR key")
	}
	fmt.Printf(createOCRKeyMsg, key.ID, key.PublicKeyAddressOnChain(), key.PublicKeyOffChain())
	return nil
}

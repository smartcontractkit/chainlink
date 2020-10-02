package cmd

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	clipkg "github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/store/orm"
)

// CreateP2PKey creates a key and inserts it into encrypted_p2p_keys,
// protected by the password in the password file
func (cli *Client) CreateP2PKey(c *clipkg.Context) error {
	return cli.errorOut(cli.createP2PKey(c))
}

// DeleteP2PKey creates a P2P key protected by the password in the password file
func (cli *Client) DeleteP2PKey(c *clipkg.Context) error {
	return cli.errorOut(cli.deleteP2PKey(c))

}

// ListP2PKeys lists the available P2P keys
func (cli *Client) ListP2PKeys(c *clipkg.Context) error {
	return cli.errorOut(cli.listP2PKeys(c))
}

const createKeyMsg = `Created P2P keypair.
Key ID
  %v
Public key:
  0x%x
Peer ID:
  %s
`

func (cli *Client) createP2PKey(c *clipkg.Context) error {
	cli.Config.Dialect = orm.DialectPostgresWithoutLock
	store := cli.AppFactory.NewApplication(cli.Config).GetStore()

	password, err := getPassword(c)
	if err != nil {
		return err
	}
	_, enc, err := store.OCRKeyStore.GenerateEncryptedP2PKey(string(password))
	if err != nil {
		return errors.Wrap(err, "while generating encrypted p2p key")
	}
	fmt.Printf(createKeyMsg, enc.ID, enc.PubKey, peer.ID(enc.PeerID).Pretty())
	return nil
}

func (cli *Client) listP2PKeys(c *clipkg.Context) error {
	cli.Config.Dialect = orm.DialectPostgresWithoutLock
	store := cli.AppFactory.NewApplication(cli.Config).GetStore()
	keys, err := store.FindEncryptedP2PKeys()
	if err != nil {
		return errors.Wrapf(err, "while fetching encrypted OCR key bundles")
	}

	fmt.Println(
		`***********************************************************************************
Encrypted P2P keys in DB
***********************************************************************************`)
	for keyidx, key := range keys {
		fmt.Println("ID                ", key.ID)
		fmt.Println("PeerID            ", key.PeerID)
		fmt.Println("Public Key        ", hex.EncodeToString(key.PubKey))
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

func (cli *Client) deleteP2PKey(c *clipkg.Context) error {
	if !c.Args().Present() {
		return errors.New("Must pass the ID of the P2P key to delete")
	}
	strID := c.Args().First()

	id, err := strconv.ParseInt(strID, 10, 32)
	if err != nil {
		return errors.New("Unable to convert provided P2P ID into integer")
	}

	cli.Config.Dialect = orm.DialectPostgresWithoutLock
	store := cli.AppFactory.NewApplication(cli.Config).GetStore()

	key, err := store.FindEncryptedP2PKeyByID(int32(id))
	if gorm.IsRecordNotFoundError(err) {
		return errors.New("Unable to find the P2P key with the provided ID")
	} else if err != nil {
		return errors.Wrapf(err, "while fetching the P2P key")
	}

	err = store.DeleteEncryptedP2PKey(key)
	if err != nil {
		return errors.Wrapf(err, "while deleting the P2P key")
	}

	fmt.Printf("Successfully deleted P2P key %d\n", key.ID)
	return nil
}

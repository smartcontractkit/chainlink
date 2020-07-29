package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	clipkg "github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/store/models/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

// CreateP2PKey creates a key and inserts it into encrypted_p2p_keys,
// protected by the password in the password file
func (cli *Client) CreateP2PKey(c *clipkg.Context) error {
	return cli.errorOut(cli.createP2PKey(c))
}

func (cli *Client) createP2PKey(c *clipkg.Context) error {
	cli.Config.Dialect = orm.DialectPostgresWithoutLock
	store := cli.AppFactory.NewApplication(cli.Config).GetStore()

	password, err := getPassword(c)
	if err != nil {
		return err
	}
	k, err := p2pkey.CreateKey()
	if err != nil {
		return errors.Wrapf(err, "while generating new p2p key")
	}
	enc, err := k.Encrypt(string(password))
	if err != nil {
		return errors.Wrapf(err, "while encrypting p2p key")
	}
	err = store.UpsertEncryptedP2PKey(&enc)
	if err != nil {
		return errors.Wrapf(err, "while inserting p2p key")
	}
	peerID, err := k.GetPeerID()
	if err != nil {
		return errors.Wrapf(err, "while getting peer ID")
	}
	fmt.Printf(`Created P2P keypair.
Key ID
  %v
Public key:
  0x%x
Peer ID:
  %s
`, enc.ID, enc.PubKey, peerID.Pretty())
	return nil
}

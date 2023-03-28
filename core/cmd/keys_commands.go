package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// KeysClient is a generic client interface for any type of key.
type KeysClient interface {
	CreateKey(*cli.Context) error
	ImportKey(*cli.Context) error
	ExportKey(*cli.Context) error
	DeleteKey(*cli.Context) error
	ListKeys(*cli.Context) error
}

// keysCommand returns a cli.Command with subcommands for the given KeysClient.
func keysCommand(typ string, c KeysClient) cli.Command {
	lower := strings.ToLower(typ)
	return cli.Command{
		Name:  lower,
		Usage: fmt.Sprintf("Remote commands for administering the node's %s keys", typ),
		Subcommands: cli.Commands{
			{
				Name:   "create",
				Usage:  fmt.Sprintf("Create a %s key", typ),
				Action: c.CreateKey,
			},
			{
				Name:  "import",
				Usage: fmt.Sprintf("Import %s key from keyfile", typ),
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "old-password, oldpassword, p",
						Usage: "`FILE` containing the password used to encrypt the key in the JSON file",
					},
				},
				Action: c.ImportKey,
			},
			{
				Name:  "export",
				Usage: fmt.Sprintf("Export %s key to keyfile", typ),
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "new-password, newpassword, p",
						Usage: "`FILE` containing the password to encrypt the key (required)",
					},
					cli.StringFlag{
						Name:  "output, o",
						Usage: "`FILE` where the JSON file will be saved (required)",
					},
				},
				Action: c.ExportKey,
			},
			{
				Name:  "delete",
				Usage: fmt.Sprintf("Delete %s key if present", typ),
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "skip the confirmation prompt",
					},
					cli.BoolFlag{
						Name:  "hard",
						Usage: "hard-delete the key instead of archiving (irreversible!)",
					},
				},
				Action: c.DeleteKey,
			},
			{
				Name: "list", Usage: fmt.Sprintf("List the %s keys", typ),
				Action: c.ListKeys,
			},
		},
	}
}

type keysClient[K keystore.Key, P TableRenderer, P2 ~[]P] struct {
	*Client
	typ  string
	path string
}

// newKeysClient returns a new KeysClient for a particular type of keystore.Key.
// P is a TableRenderer corresponding to K, and P2 is the slice variant.
func newKeysClient[K keystore.Key, P TableRenderer, P2 ~[]P](typ string, c *Client) KeysClient {
	lower := strings.ToLower(typ)
	return &keysClient[K, P, P2]{
		Client: c,
		typ:    typ,
		path:   "/v2/keys/" + lower,
	}
}

// ListKeys retrieves a list of all keys
func (cli *keysClient[K, P, P2]) ListKeys(c *cli.Context) (err error) {
	resp, err := cli.HTTP.Get(cli.path, nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var p2 P2
	return cli.renderAPIResponse(resp, &p2)
}

// CreateKey creates a new key
func (cli *keysClient[K, P, P2]) CreateKey(c *cli.Context) (err error) {
	resp, err := cli.HTTP.Post(cli.path, nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var p P
	return cli.renderAPIResponse(resp, &p, fmt.Sprintf("Created %s keypair", cli.typ))
}

// DeleteKey deletes a key,
// key ID must be passed
func (cli *keysClient[K, P, P2]) DeleteKey(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the key ID to be deleted"))
	}
	id := c.Args().Get(0)

	if !confirmAction(c) {
		return nil
	}

	var queryStr string
	if c.Bool("hard") {
		queryStr = "?hard=true"
	}

	resp, err := cli.HTTP.Delete(fmt.Sprintf(cli.path+"/%s%s", id, queryStr))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var p P
	return cli.renderAPIResponse(resp, &p, " key deleted")
}

// ImportKey imports and stores a key,
// path to key must be passed
func (cli *keysClient[K, P, P2]) ImportKey(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the filepath of the key to be imported"))
	}

	oldPasswordFile := c.String("old-password")
	if len(oldPasswordFile) == 0 {
		return cli.errorOut(errors.New("Must specify --old-password/-p flag"))
	}
	oldPassword, err := os.ReadFile(oldPasswordFile)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.Args().Get(0)
	keyJSON, err := os.ReadFile(filepath)
	if err != nil {
		return cli.errorOut(err)
	}

	normalizedPassword := normalizePassword(string(oldPassword))
	resp, err := cli.HTTP.Post(cli.path+"/import?oldpassword="+normalizedPassword, bytes.NewReader(keyJSON))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var p P
	return cli.renderAPIResponse(resp, &p, fmt.Sprintf("🔑 Imported %s key", cli.typ))
}

// ExportKey exports a key,
// key ID must be passed
func (cli *keysClient[K, P, P2]) ExportKey(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the ID of the key to export"))
	}

	newPasswordFile := c.String("new-password")
	if len(newPasswordFile) == 0 {
		return cli.errorOut(errors.New("Must specify --new-password/-p flag"))
	}
	newPassword, err := os.ReadFile(newPasswordFile)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.String("output")
	if len(filepath) == 0 {
		return cli.errorOut(errors.New("Must specify --output/-o flag"))
	}

	ID := c.Args().Get(0)

	normalizedPassword := normalizePassword(string(newPassword))
	resp, err := cli.HTTP.Post(cli.path+"/export/"+ID+"?newpassword="+normalizedPassword, nil)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not make HTTP request"))
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return cli.errorOut(fmt.Errorf("error exporting: %w", httpError(resp)))
	}

	keyJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read response body"))
	}

	err = utils.WriteFileWithMaxPerms(filepath, keyJSON, 0600)
	if err != nil {
		return cli.errorOut(errors.Wrapf(err, "Could not write %v", filepath))
	}

	_, err = os.Stderr.WriteString(fmt.Sprintf("🔑 Exported %s key %s to %s\n", cli.typ, ID, filepath))
	if err != nil {
		return cli.errorOut(err)
	}

	return nil
}

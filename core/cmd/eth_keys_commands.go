package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/urfave/cli"
	"go.uber.org/multierr"
)

// ListETHKeys renders the active account address with its ETH & LINK balance
func (cli *Client) ListETHKeys(c *cli.Context) (err error) {
	resp, err := cli.HTTP.Get("/v2/keys/eth")
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var keys []presenters.ETHKeyResource
	return cli.renderAPIResponse(resp, &keys, "🔑 ETH keys")
}

// CreateETHKey creates a new ethereum key with the same password
// as the one used to unlock the existing key.
func (cli *Client) CreateETHKey(c *cli.Context) (err error) {
	resp, err := cli.HTTP.Post("/v2/keys/eth", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var keys presenters.ETHKeyResource
	return cli.renderAPIResponse(resp, &keys, "ETH key created.\n\n🔑 New key")
}

// DeleteETHKey deletes an Ethereum key,
// address of key must be passed
func (cli *Client) DeleteETHKey(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the address of the key to be deleted"))
	}

	if c.Bool("hard") && !confirmAction(c) {
		return nil
	}

	var queryStr string
	var confirmationMsg string
	if c.Bool("hard") {
		queryStr = "?hard=true"
		confirmationMsg = "Deleted ETH key"
	} else {
		confirmationMsg = "Archived ETH key"
	}

	address := c.Args().Get(0)
	resp, err := cli.HTTP.Delete(fmt.Sprintf("/v2/keys/eth/%s%s", address, queryStr))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var key presenters.ETHKeyResource
	return cli.renderAPIResponse(resp, &key, fmt.Sprintf("🔑 %s", confirmationMsg))
}

// ImportETHKey imports an Ethereum key,
// file path must be passed
func (cli *Client) ImportETHKey(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the filepath of the key to be imported"))
	}

	oldPasswordFile := c.String("oldpassword")
	if len(oldPasswordFile) == 0 {
		return cli.errorOut(errors.New("Must specify --oldpassword/-p flag"))
	}
	oldPassword, err := ioutil.ReadFile(oldPasswordFile)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.Args().Get(0)
	keyJSON, err := ioutil.ReadFile(filepath)
	if err != nil {
		return cli.errorOut(err)
	}

	normalizedPassword := normalizePassword(string(oldPassword))
	resp, err := cli.HTTP.Post("/v2/keys/eth/import?oldpassword="+normalizedPassword, bytes.NewReader(keyJSON))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var key presenters.ETHKeyResource
	return cli.renderAPIResponse(resp, &key, "🔑 Imported ETH key")
}

// ExportETHKey exports an ETH key,
// address must be passed
func (cli *Client) ExportETHKey(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the address of the key to export"))
	}

	newPasswordFile := c.String("newpassword")
	if len(newPasswordFile) == 0 {
		return cli.errorOut(errors.New("Must specify --newpassword/-p flag"))
	}
	newPassword, err := ioutil.ReadFile(newPasswordFile)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.String("output")
	if len(newPassword) == 0 {
		return cli.errorOut(errors.New("Must specify --output/-o flag"))
	}

	address := c.Args().Get(0)

	normalizedPassword := normalizePassword(string(newPassword))
	resp, err := cli.HTTP.Post("/v2/keys/eth/export/"+address+"?newpassword="+normalizedPassword, nil)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not make HTTP request"))
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return cli.errorOut(errors.New("Error exporting"))
	}

	keyJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read response body"))
	}

	err = utils.WriteFileWithMaxPerms(filepath, keyJSON, 0600)
	if err != nil {
		return cli.errorOut(errors.Wrapf(err, "Could not write %v", filepath))
	}

	_, err = os.Stderr.WriteString("🔑 Exported ETH key " + address + " to " + filepath + "\n")
	if err != nil {
		return cli.errorOut(err)
	}

	return nil
}

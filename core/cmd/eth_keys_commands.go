package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/urfave/cli"
	"go.uber.org/multierr"
)

type EthKeyPresenter struct {
	presenters.ETHKeyResource
}

func (p *EthKeyPresenter) ToRow() []string {
	return []string{
		p.Address,
		p.EVMChainID.String(),
		p.EthBalance.String(),
		p.LinkBalance.String(),
		fmt.Sprintf("%v", p.IsFunding),
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
		p.MaxGasPriceWei.String(),
	}
}

var ethKeysTableHeaders = []string{"Address", "EVM Chain ID", "ETH", "LINK", "Is funding", "Created", "Updated", "Max Gas Price Wei"}

// RenderTable implements TableRenderer
func (p *EthKeyPresenter) RenderTable(rt RendererTable) error {
	rows := [][]string{p.ToRow()}

	renderList(ethKeysTableHeaders, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

type EthKeyPresenters []EthKeyPresenter

// RenderTable implements TableRenderer
func (ps EthKeyPresenters) RenderTable(rt RendererTable) error {
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(ethKeysTableHeaders, rows, rt.Writer)

	return nil
}

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

	return cli.renderAPIResponse(resp, &EthKeyPresenters{}, "🔑 ETH keys")
}

// CreateETHKey creates a new ethereum key with the same password
// as the one used to unlock the existing key.
func (cli *Client) CreateETHKey(c *cli.Context) (err error) {
	createUrl := url.URL{
		Path: "/v2/keys/eth",
	}
	query := createUrl.Query()

	if c.IsSet("evmChainID") {
		query.Set("evmChainID", c.String("evmChainID"))
	}
	if c.IsSet("maxGasPriceGWei") {
		query.Set("maxGasPriceGWei", c.String("maxGasPriceGWei"))
	}

	createUrl.RawQuery = query.Encode()
	resp, err := cli.HTTP.Post(createUrl.String(), nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &EthKeyPresenter{}, "ETH key created.\n\n🔑 New key")
}

// UpdateETHKey updates an Ethereum key's parameters,
// address of key must be passed as well as at least one parameter to update
func (cli *Client) UpdateETHKey(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the address of the key to be updated"))
	}
	address := c.Args().Get(0)
	updateUrl := url.URL{
		Path: "/v2/keys/eth/" + address,
	}

	query := updateUrl.Query()
	if c.IsSet("maxGasPriceGWei") {
		query.Set("maxGasPriceGWei", c.String("maxGasPriceGWei"))
	} else {
		return cli.errorOut(errors.New("Must pass at least one parameter to update"))
	}

	updateUrl.RawQuery = query.Encode()
	resp, err := cli.HTTP.Put(updateUrl.String(), nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &EthKeyPresenter{}, "ETH key updated.\n\n🔑 Updated key")
}

// DeleteETHKey deletes an Ethereum key,
// address of key must be passed
func (cli *Client) DeleteETHKey(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the address of the key to be deleted"))
	}
	address := c.Args().Get(0)
	deleteUrl := url.URL{
		Path: "/v2/keys/eth/" + address,
	}
	query := deleteUrl.Query()

	if c.Bool("hard") && !confirmAction(c) {
		return nil
	}

	var confirmationMsg string
	if c.Bool("hard") {
		query.Set("hard", "true")
		confirmationMsg = "Deleted ETH key"
	} else {
		confirmationMsg = "Archived ETH key"
	}

	deleteUrl.RawQuery = query.Encode()
	resp, err := cli.HTTP.Delete(deleteUrl.String())
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &EthKeyPresenter{}, fmt.Sprintf("🔑 %s", confirmationMsg))
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

	importUrl := url.URL{
		Path: "/v2/keys/eth/import",
	}
	query := importUrl.Query()

	query.Set("oldpassword", strings.TrimSpace(string(oldPassword)))

	if c.IsSet("evmChainID") {
		query.Set("evmChainID", c.String("evmChainID"))
	}

	importUrl.RawQuery = query.Encode()
	resp, err := cli.HTTP.Post(importUrl.String(), bytes.NewReader(keyJSON))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &EthKeyPresenter{}, "🔑 Imported ETH key")
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
	exportUrl := url.URL{
		Path: "/v2/keys/eth/export/" + address,
	}
	query := exportUrl.Query()
	query.Set("newpassword", strings.TrimSpace(string(newPassword)))

	exportUrl.RawQuery = query.Encode()
	resp, err := cli.HTTP.Post(exportUrl.String(), nil)
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

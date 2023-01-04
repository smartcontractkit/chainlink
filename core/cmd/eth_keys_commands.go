package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

type EthKeyPresenter struct {
	presenters.ETHKeyResource
}

func (p *EthKeyPresenter) ToRow() []string {
	return []string{
		p.Address,
		p.EVMChainID.String(),
		fmt.Sprintf("%d", p.NextNonce),
		p.EthBalance.String(),
		p.LinkBalance.String(),
		fmt.Sprintf("%v", p.Disabled),
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
		p.MaxGasPriceWei.String(),
	}
}

var ethKeysTableHeaders = []string{"Address", "EVM Chain ID", "Next Nonce", "ETH", "LINK", "Disabled", "Created", "Updated", "Max Gas Price Wei"}

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
	resp, err := cli.HTTP.Get("/v2/keys/evm")
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
		Path: "/v2/keys/evm",
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
		Path: "/v2/keys/evm/" + address,
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

// DeleteETHKey hard deletes an Ethereum key,
// address of key must be passed
func (cli *Client) DeleteETHKey(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the address of the key to be deleted"))
	}
	address := c.Args().Get(0)

	if !confirmAction(c) {
		return nil
	}

	resp, err := cli.HTTP.Delete("/v2/keys/evm/" + address)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode != http.StatusNoContent {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return cli.errorOut(errors.Wrap(err, "Failed to read request response"))
		}
		var result *models.JSONAPIErrors
		err = json.Unmarshal(body, &result)
		if err != nil {
			return cli.errorOut(errors.Wrapf(err, "Unable to unmarshal json from body '%s'", string(body)))
		}
		return cli.errorOut(errors.Errorf("Delete ETH key failed: %s", result.Error()))
	}
	fmt.Println(fmt.Sprintf("Deleted ETH key: %s", address))
	return nil
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
	oldPassword, err := os.ReadFile(oldPasswordFile)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.Args().Get(0)
	keyJSON, err := os.ReadFile(filepath)
	if err != nil {
		return cli.errorOut(err)
	}

	importUrl := url.URL{
		Path: "/v2/keys/evm/import",
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
	newPassword, err := os.ReadFile(newPasswordFile)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.String("output")
	if len(newPassword) == 0 {
		return cli.errorOut(errors.New("Must specify --output/-o flag"))
	}

	address := c.Args().Get(0)
	exportUrl := url.URL{
		Path: "/v2/keys/evm/export/" + address,
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

	keyJSON, err := io.ReadAll(resp.Body)
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

// UpdateChainEVMKey updates settings for the given key on the given chain
func (cli *Client) UpdateChainEVMKey(c *cli.Context) (err error) {
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	addr := c.String("address")
	query.Set("address", addr)
	cid := c.String("evmChainID")
	query.Set("evmChainID", cid)
	abandon := c.String("abandon")
	query.Set("abandon", abandon)

	if c.IsSet("setNextNonce") {
		query.Set("nextNonce", c.String("setNextNonce"))
	}
	if c.IsSet("enable") && c.IsSet("disable") {
		return cli.errorOut(errors.New("cannot set both --enable and --disable simultaneously"))
	} else if c.Bool("enable") {
		query.Set("enabled", "true")
	} else if c.Bool("disable") {
		query.Set("enabled", "false")
	}

	chainURL.RawQuery = query.Encode()
	resp, err := cli.HTTP.Post(chainURL.String(), nil)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not make HTTP request"))
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		resp, err := io.ReadAll(resp.Body)
		if err != nil {
			return cli.errorOut(errors.Errorf("Error resetting key: %s", err.Error()))
		}
		return cli.errorOut(errors.Errorf("Error resetting key: %s", resp))
	}

	return cli.renderAPIResponse(resp, &EthKeyPresenter{}, "🔑 Updated ETH key")
}

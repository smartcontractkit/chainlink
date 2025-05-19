package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	cutils "github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func initEthKeysSubCmd(s *Shell) cli.Command {
	return cli.Command{
		Name:  "eth",
		Usage: "Remote commands for administering the node's Ethereum keys",
		Subcommands: cli.Commands{
			{
				Name:   "create",
				Usage:  "Create a key in the node's keystore alongside the existing key; to create an original key, just run the node",
				Action: s.CreateETHKey,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "evm-chain-id, evmChainID",
						Usage: "Chain ID for the key. If left blank, default chain will be used.",
					},
					cli.Uint64Flag{
						Name:  "max-gas-price-gwei, maxGasPriceGWei",
						Usage: "Optional maximum gas price (GWei) for the creating key.",
					},
				},
			},
			{
				Name:   "list",
				Usage:  "List available Ethereum accounts with their ETH & LINK balances and other metadata",
				Action: s.ListETHKeys,
			},
			{
				Name:  "delete",
				Usage: format(`Delete the ETH key by address (irreversible!)`),
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "skip the confirmation prompt",
					},
				},
				Action: s.DeleteETHKey,
			},
			{
				Name:  "import",
				Usage: format(`Import an ETH key from a JSON file`),
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "old-password, oldpassword, p",
						Usage: "`FILE` containing the password used to encrypt the key in the JSON file",
					},
					cli.StringFlag{
						Name:  "evm-chain-id, evmChainID",
						Usage: "Chain ID for the key. If left blank, default chain will be used.",
					},
				},
				Action: s.ImportETHKey,
			},
			{
				Name:  "export",
				Usage: format(`Exports an ETH key to a JSON file`),
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "new-password, newpassword, p",
						Usage: "`FILE` containing the password to encrypt the key (required)",
					},
					cli.StringFlag{
						Name:  "output, o",
						Usage: "Path where the JSON file will be saved (required)",
					},
				},
				Action: s.ExportETHKey,
			},
			{
				Name:   "chain",
				Usage:  "Update an EVM key for the given chain",
				Action: s.UpdateChainEVMKey,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     "address",
						Usage:    "address of the key",
						Required: true,
					},
					cli.StringFlag{
						Name:     "evm-chain-id, evmChainID",
						Usage:    "chain ID of the key",
						Required: true,
					},
					cli.BoolFlag{
						Name:  "enable",
						Usage: "enable the key for the given chain",
					},
					cli.BoolFlag{
						Name:  "disable",
						Usage: "disable the key for the given chain",
					},
					cli.BoolFlag{
						Name:  "abandon",
						Usage: "if set, will abandon all pending and unconfirmed transactions and mark them as fatally errored. Use with caution, this can result in nonce gaps or 'stuck' transactions",
					},
				},
			},
		},
	}
}

type EthKeyPresenter struct {
	presenters.ETHKeyResource
}

func (p *EthKeyPresenter) ToRow() []string {
	eth := "Unknown"
	if p.EthBalance != nil {
		eth = p.EthBalance.String()
	}
	link := "Unknown"
	if p.LinkBalance != nil {
		link = p.LinkBalance.String()
	}
	gas := "None"
	if p.MaxGasPriceWei != nil {
		gas = p.MaxGasPriceWei.String()
	}
	return []string{
		p.Address,
		p.EVMChainID.String(),
		eth,
		link,
		strconv.FormatBool(p.Disabled),
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
		gas,
	}
}

var ethKeysTableHeaders = []string{"Address", "EVM Chain ID", "ETH", "LINK", "Disabled", "Created", "Updated", "Max Gas Price Wei"}

// RenderTable implements TableRenderer
func (p *EthKeyPresenter) RenderTable(rt RendererTable) error {
	rows := [][]string{p.ToRow()}

	renderList(ethKeysTableHeaders, rows, rt.Writer)

	return cutils.JustError(rt.Write([]byte("\n")))
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
func (s *Shell) ListETHKeys(_ *cli.Context) (err error) {
	resp, err := s.HTTP.Get(s.ctx(), "/v2/keys/evm")

	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return s.renderAPIResponse(resp, &EthKeyPresenters{}, "🔑 ETH keys")
}

// CreateETHKey creates a new ethereum key with the same password
// as the one used to unlock the existing key.
func (s *Shell) CreateETHKey(c *cli.Context) (err error) {
	createUrl := url.URL{
		Path: "/v2/keys/evm",
	}
	query := createUrl.Query()

	if c.IsSet("evm-chain-id") {
		query.Set("evmChainID", c.String("evm-chain-id"))
	}
	if c.IsSet("max-gas-price-gwei") {
		query.Set("maxGasPriceGWei", c.String("max-gas-price-gwei"))
	}

	createUrl.RawQuery = query.Encode()
	resp, err := s.HTTP.Post(s.ctx(), createUrl.String(), nil)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return s.renderAPIResponse(resp, &EthKeyPresenter{}, "ETH key created.\n\n🔑 New key")
}

// DeleteETHKey hard deletes an Ethereum key,
// address of key must be passed
func (s *Shell) DeleteETHKey(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return s.errorOut(errors.New("Must pass the address of the key to be deleted"))
	}
	address := c.Args().Get(0)

	if !confirmAction(c) {
		return nil
	}

	resp, err := s.HTTP.Delete(s.ctx(), "/v2/keys/evm/"+address)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return s.errorOut(errors.Wrap(err, "Failed to read request response"))
		}
		var result *models.JSONAPIErrors
		err = json.Unmarshal(body, &result)
		if err != nil {
			return s.errorOut(errors.Wrapf(err, "Unable to unmarshal json from body '%s'", string(body)))
		}
		return s.errorOut(errors.Errorf("Delete ETH key failed: %s", result.Error()))
	}
	return s.renderAPIResponse(resp, &EthKeyPresenter{}, fmt.Sprintf("🔑 Deleted ETH key: %s\n", address))
}

// ImportETHKey imports an Ethereum key,
// file path must be passed
func (s *Shell) ImportETHKey(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return s.errorOut(errors.New("Must pass the filepath of the key to be imported"))
	}

	oldPasswordFile := c.String("old-password")
	if len(oldPasswordFile) == 0 {
		return s.errorOut(errors.New("Must specify --old-password/-p flag"))
	}
	oldPassword, err := os.ReadFile(oldPasswordFile)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.Args().Get(0)
	keyJSON, err := os.ReadFile(filepath)
	if err != nil {
		return s.errorOut(err)
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
	resp, err := s.HTTP.Post(s.ctx(), importUrl.String(), bytes.NewReader(keyJSON))
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return s.renderAPIResponse(resp, &EthKeyPresenter{}, "🔑 Imported ETH key")
}

// ExportETHKey exports an ETH key,
// address must be passed
func (s *Shell) ExportETHKey(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return s.errorOut(errors.New("Must pass the address of the key to export"))
	}

	newPasswordFile := c.String("new-password")
	if len(newPasswordFile) == 0 {
		return s.errorOut(errors.New("Must specify --new-password/-p flag"))
	}
	newPassword, err := os.ReadFile(newPasswordFile)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.String("output")
	if len(newPassword) == 0 {
		return s.errorOut(errors.New("Must specify --output/-o flag"))
	}

	address := c.Args().Get(0)
	exportUrl := url.URL{
		Path: "/v2/keys/evm/export/" + address,
	}
	query := exportUrl.Query()
	query.Set("newpassword", strings.TrimSpace(string(newPassword)))

	exportUrl.RawQuery = query.Encode()
	resp, err := s.HTTP.Post(s.ctx(), exportUrl.String(), nil)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "Could not make HTTP request"))
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return s.errorOut(fmt.Errorf("error exporting: %w", httpError(resp)))
	}

	keyJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "Could not read response body"))
	}

	err = utils.WriteFileWithMaxPerms(filepath, keyJSON, 0o600)
	if err != nil {
		return s.errorOut(errors.Wrapf(err, "Could not write %v", filepath))
	}

	_, err = os.Stderr.WriteString("🔑 Exported ETH key " + address + " to " + filepath + "\n")
	if err != nil {
		return s.errorOut(err)
	}

	return nil
}

// UpdateChainEVMKey updates settings for the given key on the given chain
func (s *Shell) UpdateChainEVMKey(c *cli.Context) (err error) {
	chainURL := url.URL{Path: "/v2/keys/evm/chain"}
	query := chainURL.Query()

	addr := c.String("address")
	query.Set("address", addr)
	cid := c.String("evmChainID")
	query.Set("evmChainID", cid)
	abandon := c.String("abandon")
	query.Set("abandon", abandon)

	if c.IsSet("enable") && c.IsSet("disable") {
		return s.errorOut(errors.New("cannot set both --enable and --disable simultaneously"))
	} else if c.Bool("enable") {
		query.Set("enabled", "true")
	} else if c.Bool("disable") {
		query.Set("enabled", "false")
	}

	chainURL.RawQuery = query.Encode()
	resp, err := s.HTTP.Post(s.ctx(), chainURL.String(), nil)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "Could not make HTTP request"))
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return s.errorOut(fmt.Errorf("error resetting key: %w", httpError(resp)))
	}

	return s.renderAPIResponse(resp, &EthKeyPresenter{}, "🔑 Updated ETH key")
}

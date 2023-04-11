package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func initOCR2KeysSubCmd(client *Client) cli.Command {
	return cli.Command{
		Name:  "ocr2",
		Usage: "Remote commands for administering the node's off chain reporting keys",
		Subcommands: cli.Commands{
			{
				Name:   "create",
				Usage:  format(`Create an OCR2 key bundle, encrypted with password from the password file, and store it in the database`),
				Action: client.CreateOCR2KeyBundle,
			},
			{
				Name:  "delete",
				Usage: format(`Deletes the encrypted OCR2 key bundle matching the given ID`),
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
				Action: client.DeleteOCR2KeyBundle,
			},
			{
				Name:   "list",
				Usage:  format(`List available OCR2 key bundles`),
				Action: client.ListOCR2KeyBundles,
			},
			{
				Name:  "import",
				Usage: format(`Imports an OCR2 key bundle from a JSON file`),
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "old-password, oldpassword, p",
						Usage: "`FILE` containing the password used to encrypt the key in the JSON file",
					},
				},
				Action: client.ImportOCR2Key,
			},
			{
				Name:  "export",
				Usage: format(`Exports an OCR2 key bundle to a JSON file`),
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
				Action: client.ExportOCR2Key,
			},
		},
	}
}

type OCR2KeyBundlePresenter struct {
	JAID // Include this to overwrite the presenter JAID so it can correctly render the ID in JSON
	presenters.OCR2KeysBundleResource
}

// RenderTable implements TableRenderer
func (p *OCR2KeyBundlePresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Type", "On-chain pubkey", "Off-chain pubkey", "Config pubkey"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 OCR Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p *OCR2KeyBundlePresenter) ToRow() []string {
	return []string{
		p.ID,
		p.ChainType,
		p.OnchainPublicKey,
		p.OffChainPublicKey,
		p.ConfigPublicKey,
	}
}

type OCR2KeyBundlePresenters []OCR2KeyBundlePresenter

// RenderTable implements TableRenderer
func (ps OCR2KeyBundlePresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Type", "On-chain pubkey", "Off-chain pubkey", "Config pubkey"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 OCR Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

// ListOCR2KeyBundles lists the available OCR2 Key Bundles
func (cli *Client) ListOCR2KeyBundles(c *cli.Context) error {
	resp, err := cli.HTTP.Get("/v2/keys/ocr2", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenters OCR2KeyBundlePresenters
	return cli.renderAPIResponse(resp, &presenters)
}

// CreateOCR2KeyBundle creates an OCR2 key bundle and saves it to the keystore
func (cli *Client) CreateOCR2KeyBundle(c *cli.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(
			errors.Errorf(`must pass the type to create, options are: %s`, chaintype.SupportedChainTypes.String()),
		)
	}
	chainType := c.Args().Get(0)
	resp, err := cli.HTTP.Post(fmt.Sprintf("/v2/keys/ocr2/%s", chainType), nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenter OCR2KeyBundlePresenter
	return cli.renderAPIResponse(resp, &presenter, "Created OCR key bundle")
}

// DeleteOCR2KeyBundle deletes an OCR2 key bundle
func (cli *Client) DeleteOCR2KeyBundle(c *cli.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the key ID to be deleted"))
	}
	id, err := models.Sha256HashFromHex(c.Args().Get(0))
	if err != nil {
		return cli.errorOut(err)
	}

	if !confirmAction(c) {
		return nil
	}

	var queryStr string
	if c.Bool("hard") {
		queryStr = "?hard=true"
	}

	resp, err := cli.HTTP.Delete(fmt.Sprintf("/v2/keys/ocr2/%s%s", id, queryStr))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenter OCR2KeyBundlePresenter
	return cli.renderAPIResponse(resp, &presenter, "OCR key bundle deleted")
}

// ImportOCR2Key imports OCR2 key bundle
func (cli *Client) ImportOCR2Key(c *cli.Context) (err error) {
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
	resp, err := cli.HTTP.Post("/v2/keys/ocr2/import?oldpassword="+normalizedPassword, bytes.NewReader(keyJSON))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenter OCR2KeyBundlePresenter
	return cli.renderAPIResponse(resp, &presenter, "Imported OCR key bundle")
}

// ExportOCR2Key exports an OCR2 key bundle by ID
func (cli *Client) ExportOCR2Key(c *cli.Context) (err error) {
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
	resp, err := cli.HTTP.Post("/v2/keys/ocr2/export/"+ID+"?newpassword="+normalizedPassword, nil)
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

	err = utils.WriteFileWithMaxPerms(filepath, keyJSON, 0o600)
	if err != nil {
		return cli.errorOut(errors.Wrapf(err, "Could not write %v", filepath))
	}

	_, err = os.Stderr.WriteString(fmt.Sprintf("Exported OCR key bundle %s to %s", ID, filepath))
	if err != nil {
		return cli.errorOut(err)
	}

	return nil
}

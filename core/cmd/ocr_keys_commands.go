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

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func initOCRKeysSubCmd(client *Client) cli.Command {
	return cli.Command{
		Name:  "ocr",
		Usage: "Remote commands for administering the node's legacy off chain reporting keys",
		Subcommands: cli.Commands{
			{
				Name:   "create",
				Usage:  format(`Create an OCR key bundle, encrypted with password from the password file, and store it in the database`),
				Action: client.CreateOCRKeyBundle,
			},
			{
				Name:  "delete",
				Usage: format(`Deletes the encrypted OCR key bundle matching the given ID`),
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
				Action: client.DeleteOCRKeyBundle,
			},
			{
				Name:   "list",
				Usage:  format(`List available OCR key bundles`),
				Action: client.ListOCRKeyBundles,
			},
			{
				Name:  "import",
				Usage: format(`Imports an OCR key bundle from a JSON file`),
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "old-password, oldpassword, p",
						Usage: "`FILE` containing the password used to encrypt the key in the JSON file",
					},
				},
				Action: client.ImportOCRKey,
			},
			{
				Name:  "export",
				Usage: format(`Exports an OCR key bundle to a JSON file`),
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
				Action: client.ExportOCRKey,
			},
		},
	}
}

type OCRKeyBundlePresenter struct {
	JAID // Include this to overwrite the presenter JAID so it can correctly render the ID in JSON
	presenters.OCRKeysBundleResource
}

// RenderTable implements TableRenderer
func (p *OCRKeyBundlePresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "On-chain signing addr", "Off-chain pubkey", "Config pubkey"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 Legacy OCR Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p *OCRKeyBundlePresenter) ToRow() []string {
	return []string{
		p.ID,
		p.OnChainSigningAddress.String(),
		p.OffChainPublicKey.String(),
		p.ConfigPublicKey.String(),
	}
}

type OCRKeyBundlePresenters []OCRKeyBundlePresenter

// RenderTable implements TableRenderer
func (ps OCRKeyBundlePresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "On-chain signing addr", "Off-chain pubkey", "Config pubkey"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 Legacy OCR Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

// ListOCRKeyBundles lists the available OCR Key Bundles
func (cli *Client) ListOCRKeyBundles(c *cli.Context) error {
	resp, err := cli.HTTP.Get("/v2/keys/ocr", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenters OCRKeyBundlePresenters
	return cli.renderAPIResponse(resp, &presenters)
}

// CreateOCR2KeyBundle creates an OCR key bundle and saves it to the keystore
func (cli *Client) CreateOCRKeyBundle(c *cli.Context) error {
	resp, err := cli.HTTP.Post("/v2/keys/ocr", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenter OCRKeyBundlePresenter
	return cli.renderAPIResponse(resp, &presenter, "Created OCR key bundle")
}

// DeleteOCR2KeyBundle deletes an OCR key bundle
func (cli *Client) DeleteOCRKeyBundle(c *cli.Context) error {
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

	resp, err := cli.HTTP.Delete(fmt.Sprintf("/v2/keys/ocr/%s%s", id, queryStr))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenter OCRKeyBundlePresenter
	return cli.renderAPIResponse(resp, &presenter, "OCR key bundle deleted")
}

// ImportOCR2Key imports OCR key bundle
func (cli *Client) ImportOCRKey(c *cli.Context) (err error) {
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
	resp, err := cli.HTTP.Post("/v2/keys/ocr/import?oldpassword="+normalizedPassword, bytes.NewReader(keyJSON))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var presenter OCRKeyBundlePresenter
	return cli.renderAPIResponse(resp, &presenter, "Imported OCR key bundle")
}

// ExportOCR2Key exports an OCR key bundle by ID
func (cli *Client) ExportOCRKey(c *cli.Context) (err error) {
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
	resp, err := cli.HTTP.Post("/v2/keys/ocr/export/"+ID+"?newpassword="+normalizedPassword, nil)
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

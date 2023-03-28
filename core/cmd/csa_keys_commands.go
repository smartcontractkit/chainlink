package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func initCSAKeysSubCmd(client *Client) cli.Command {
	return cli.Command{
		Name:  "csa",
		Usage: "Remote commands for administering the node's CSA keys",
		Subcommands: cli.Commands{
			{
				Name:   "create",
				Usage:  format(`Create a CSA key, encrypted with password from the password file, and store it in the database.`),
				Action: client.CreateCSAKey,
			},
			{
				Name:   "list",
				Usage:  format(`List available CSA keys`),
				Action: client.ListCSAKeys,
			},
			{
				Name:  "import",
				Usage: format(`Imports a CSA key from a JSON file.`),
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "old-password, oldpassword, p",
						Usage: "`FILE` containing the password used to encrypt the key in the JSON file",
					},
				},
				Action: client.ImportCSAKey,
			},
			{
				Name:  "export",
				Usage: format(`Exports an existing CSA key by its ID.`),
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
				Action: client.ExportCSAKey,
			},
		},
	}
}

type CSAKeyPresenter struct {
	JAID
	presenters.CSAKeyResource
}

// RenderTable implements TableRenderer
func (p *CSAKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"Public key"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 CSA Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return nil
}

func (p *CSAKeyPresenter) ToRow() []string {
	row := []string{
		p.PubKey,
	}

	return row
}

type CSAKeyPresenters []CSAKeyPresenter

// RenderTable implements TableRenderer
func (ps CSAKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"Public key"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 CSA Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)
	return utils.JustError(rt.Write([]byte("\n")))
}

// ListCSAKeys retrieves a list of all CSA keys
func (cli *Client) ListCSAKeys(c *cli.Context) (err error) {
	resp, err := cli.HTTP.Get("/v2/keys/csa", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &CSAKeyPresenters{})
}

// CreateCSAKey creates a new CSA key
func (cli *Client) CreateCSAKey(c *cli.Context) (err error) {
	resp, err := cli.HTTP.Post("/v2/keys/csa", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &CSAKeyPresenter{}, "Created CSA key")
}

// ImportCSAKey imports and stores a CSA key. Path to key must be passed.
func (cli *Client) ImportCSAKey(c *cli.Context) (err error) {
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

	exportUrl := url.URL{
		Path: "/v2/keys/csa/import",
	}

	query := exportUrl.Query()
	query.Set("oldpassword", normalizePassword(string(oldPassword)))

	exportUrl.RawQuery = query.Encode()
	resp, err := cli.HTTP.Post(exportUrl.String(), bytes.NewReader(keyJSON))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &CSAKeyPresenter{}, "🔑 Imported CSA key")
}

// ExportCSAKey exports a CSA key. Key ID must be passed.
func (cli *Client) ExportCSAKey(c *cli.Context) (err error) {
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
	exportUrl := url.URL{
		Path: "/v2/keys/csa/export/" + ID,
	}

	query := exportUrl.Query()
	query.Set("newpassword", normalizePassword(string(newPassword)))

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

	_, err = os.Stderr.WriteString(fmt.Sprintf("🔑 Exported P2P key %s to %s\n", ID, filepath))
	if err != nil {
		return cli.errorOut(err)
	}

	return nil
}

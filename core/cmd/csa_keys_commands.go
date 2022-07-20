package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

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

	newPasswordFile := c.String("newpassword")
	if len(newPasswordFile) == 0 {
		return cli.errorOut(errors.New("Must specify --newpassword/-p flag"))
	}

	newPassword, err := ioutil.ReadFile(newPasswordFile)
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

	_, err = os.Stderr.WriteString(fmt.Sprintf("🔑 Exported P2P key %s to %s\n", ID, filepath))
	if err != nil {
		return cli.errorOut(err)
	}

	return nil
}

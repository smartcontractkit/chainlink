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

type SolanaKeyPresenter struct {
	JAID
	presenters.SolanaKeyResource
}

// RenderTable implements TableRenderer
func (p *SolanaKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 Solana Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p *SolanaKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.PubKey,
	}

	return row
}

type SolanaKeyPresenters []SolanaKeyPresenter

// RenderTable implements TableRenderer
func (ps SolanaKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 Solana Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

// ListSolanaKeys retrieves a list of all Solana keys
func (cli *Client) ListSolanaKeys(c *cli.Context) (err error) {
	resp, err := cli.HTTP.Get("/v2/keys/solana", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &SolanaKeyPresenters{})
}

// CreateSolanaKey creates a new Solana key
func (cli *Client) CreateSolanaKey(c *cli.Context) (err error) {
	resp, err := cli.HTTP.Post("/v2/keys/solana", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &SolanaKeyPresenter{}, "Created Solana keypair")
}

// DeleteSolanaKey deletes a Solana key,
// key ID must be passed
func (cli *Client) DeleteSolanaKey(c *cli.Context) (err error) {
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

	resp, err := cli.HTTP.Delete(fmt.Sprintf("/v2/keys/solana/%s%s", id, queryStr))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &SolanaKeyPresenter{}, "Solana key deleted")
}

// ImportSolanaKey imports and stores a Solana key,
// path to key must be passed
func (cli *Client) ImportSolanaKey(c *cli.Context) (err error) {
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
	resp, err := cli.HTTP.Post("/v2/keys/solana/import?oldpassword="+normalizedPassword, bytes.NewReader(keyJSON))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &SolanaKeyPresenter{}, "🔑 Imported Solana key")
}

// ExportSolanaKey exports a Solana key,
// key ID must be passed
func (cli *Client) ExportSolanaKey(c *cli.Context) (err error) {
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

	normalizedPassword := normalizePassword(string(newPassword))
	resp, err := cli.HTTP.Post("/v2/keys/solana/export/"+ID+"?newpassword="+normalizedPassword, nil)
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

	_, err = os.Stderr.WriteString(fmt.Sprintf("🔑 Exported Solana key %s to %s\n", ID, filepath))
	if err != nil {
		return cli.errorOut(err)
	}

	return nil
}

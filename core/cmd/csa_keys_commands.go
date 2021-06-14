package cmd

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/urfave/cli"
	"go.uber.org/multierr"
)

type CSAKeyPresenter struct {
	JAID
	presenters.CSAKeyResource
}

// RenderTable implements TableRenderer
func (p *CSAKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key", "Created", "Updated"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 CSA Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return nil
}

func (p *CSAKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.PubKey,
		fmt.Sprintf("%v", p.CreatedAt),
		fmt.Sprintf("%v", p.UpdatedAt),
	}

	return row
}

type CSAKeyPresenters []CSAKeyPresenter

// RenderTable implements TableRenderer
func (ps CSAKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key", "Created", "Updated"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 CSA Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return nil
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

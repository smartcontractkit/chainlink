package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// SolanaNodePresenter implements TableRenderer for a SolanaNodeResource.
type SolanaNodePresenter struct {
	presenters.SolanaNodeResource
}

// ToRow presents the SolanaNodeResource as a slice of strings.
func (p *SolanaNodePresenter) ToRow() []string {
	row := []string{
		p.GetID(),
		p.Name,
		p.SolanaChainID,
		p.SolanaURL,
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
	}
	return row
}

var solanaNodeHeaders = []string{"ID", "Name", "Chain ID", "URL", "Created", "Updated"}

// RenderTable implements TableRenderer
func (p SolanaNodePresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(solanaNodeHeaders, rows, rt.Writer)

	return nil
}

// SolanaNodePresenters implements TableRenderer for a slice of SolanaNodePresenter.
type SolanaNodePresenters []SolanaNodePresenter

// RenderTable implements TableRenderer
func (ps SolanaNodePresenters) RenderTable(rt RendererTable) error {
	var rows [][]string

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(solanaNodeHeaders, rows, rt.Writer)

	return nil
}

// IndexSolanaNodes returns all Solana nodes.
func (cli *Client) IndexSolanaNodes(c *cli.Context) (err error) {
	return cli.getPage("/v2/nodes/solana", c.Int("page"), &SolanaNodePresenters{})
}

// CreateSolanaNode adds a new node to the nodelink node
func (cli *Client) CreateSolanaNode(c *cli.Context) (err error) {
	name := c.String("name")
	chainID := c.String("chain-id")
	urlStr := c.String("url")

	if name == "" {
		return cli.errorOut(errors.New("missing --name"))
	}
	if chainID == "" {
		return cli.errorOut(errors.New("missing --chain-id"))
	}

	if _, err2 := url.Parse(urlStr); err2 != nil {
		return cli.errorOut(fmt.Errorf("invalid url: %v", err2))
	}

	params := db.NewNode{
		Name:          name,
		SolanaChainID: chainID,
		SolanaURL:     urlStr,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return cli.errorOut(err)
	}

	resp, err := cli.HTTP.Post("/v2/nodes/solana", bytes.NewBuffer(body))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &SolanaNodePresenter{})
}

// RemoveSolanaNode removes a specific Solana Node by name.
func (cli *Client) RemoveSolanaNode(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass the id of the node to be removed"))
	}
	nodeID := c.Args().First()
	resp, err := cli.HTTP.Delete("/v2/nodes/solana/" + nodeID)
	if err != nil {
		return cli.errorOut(err)
	}
	_, err = cli.parseResponse(resp)
	if err != nil {
		return cli.errorOut(err)
	}

	fmt.Printf("Node %v deleted\n", c.Args().First())
	return nil
}

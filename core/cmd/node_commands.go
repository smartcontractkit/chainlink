package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/urfave/cli"
	"go.uber.org/multierr"
	null "gopkg.in/guregu/null.v4"
)

type NodePresenter struct {
	presenters.NodeResource
}

func (p *NodePresenter) ToRow() []string {
	row := []string{
		p.GetID(),
		p.Name,
		p.EVMChainID.ToInt().String(),
		p.WSURL.ValueOrZero(),
		p.HTTPURL.ValueOrZero(),
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
	}
	return row
}

var nodeHeaders = []string{"ID", "Name", "Chain ID", "Websocket URL", "HTTP URL", "Created", "Updated"}

// RenderTable implements TableRenderer
func (p NodePresenter) RenderTable(rt RendererTable) error {
	rows := [][]string{}
	rows = append(rows, p.ToRow())
	renderList(nodeHeaders, rows, rt.Writer)

	return nil
}

type NodePresenters []NodePresenter

// RenderTable implements TableRenderer
func (ps NodePresenters) RenderTable(rt RendererTable) error {
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(nodeHeaders, rows, rt.Writer)

	return nil
}

// IndexNodes returns all nodes.
func (cli *Client) IndexNodes(c *cli.Context) (err error) {
	return cli.getPage("/v2/nodes", c.Int("page"), &NodePresenters{})
}

// CreateNode adds a new node to the nodelink node
func (cli *Client) CreateNode(c *cli.Context) (err error) {
	name := c.String("name")
	t := c.String("type")
	ws := c.String("ws-url")
	httpURLStr := c.String("http-url")
	chainID := c.Int64("chain-id")

	if name == "" {
		return cli.errorOut(errors.New("missing --name"))
	}
	if chainID == 0 {
		return cli.errorOut(errors.New("missing --chain-id"))
	}
	if t != "primary" && t != "sendonly" {
		return cli.errorOut(errors.New("invalid or unspecified --type, must be either primary or sendonly"))
	}
	if t == "primary" && ws == "" {
		return cli.errorOut(errors.New("missing --ws-url"))
	}
	var httpURL = null.NewString(httpURLStr, true)
	if httpURLStr == "" {
		httpURL = null.NewString(httpURLStr, false)
	}

	var wsURL null.String
	if ws != "" {
		wsURL = null.StringFrom(ws)
	}

	params := evmtypes.NewNode{
		Name:       name,
		EVMChainID: *utils.NewBigI(chainID),
		WSURL:      wsURL,
		HTTPURL:    httpURL,
		SendOnly:   t == "sendonly",
	}

	body, err := json.Marshal(params)
	if err != nil {
		return cli.errorOut(err)
	}

	resp, err := cli.HTTP.Post("/v2/nodes", bytes.NewBuffer(body))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &NodePresenter{})
}

// RemoveNode removes a specific Node by name.
func (cli *Client) RemoveNode(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass the id of the node to be removed"))
	}
	nodeID := c.Args().First()
	resp, err := cli.HTTP.Delete("/v2/nodes/" + nodeID)
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

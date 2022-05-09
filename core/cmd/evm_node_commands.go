package cmd

import (
	"errors"

	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// EVMNodePresenter implements TableRenderer for an EVMNodeResource.
type EVMNodePresenter struct {
	presenters.EVMNodeResource
}

// ToRow presents the EVMNodeResource as a slice of strings.
func (p *EVMNodePresenter) ToRow() []string {
	row := []string{
		p.GetID(),
		p.Name,
		p.EVMChainID.ToInt().String(),
		p.WSURL.ValueOrZero(),
		p.HTTPURL.ValueOrZero(),
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
		p.State,
	}
	return row
}

var evmNodeHeaders = []string{"ID", "Name", "Chain ID", "Websocket URL", "HTTP URL", "Created", "Updated", "State"}

// RenderTable implements TableRenderer
func (p EVMNodePresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(evmNodeHeaders, rows, rt.Writer)

	return nil
}

// EVMNodePresenters implements TableRenderer for a slice of EVMNodePresenter.
type EVMNodePresenters []EVMNodePresenter

// RenderTable implements TableRenderer
func (ps EVMNodePresenters) RenderTable(rt RendererTable) error {
	var rows [][]string

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(evmNodeHeaders, rows, rt.Writer)

	return nil
}

func NewEVMNodeClient(c *Client) NodeClient[evmtypes.Node, presenters.EVMNodeResource, EVMNodePresenter, EVMNodePresenters] {
	createNode := func(c *cli.Context) (node evmtypes.Node, err error) {
		name := c.String("name")
		t := c.String("type")
		ws := c.String("ws-url")
		httpURLStr := c.String("http-url")
		chainID := c.Int64("chain-id")

		if name == "" {
			err = errors.New("missing --name")
			return
		}
		if chainID == 0 {
			err = errors.New("missing --chain-id")
			return
		}
		if t != "primary" && t != "sendonly" {
			err = errors.New("invalid or unspecified --type, must be either primary or sendonly")
			return
		}
		if t == "primary" && ws == "" {
			err = errors.New("missing --ws-url")
			return
		}
		var httpURL = null.NewString(httpURLStr, true)
		if httpURLStr == "" {
			httpURL = null.NewString(httpURLStr, false)
		}

		var wsURL null.String
		if ws != "" {
			wsURL = null.StringFrom(ws)
		}

		node = evmtypes.Node{
			Name:       name,
			EVMChainID: *utils.NewBigI(chainID),
			WSURL:      wsURL,
			HTTPURL:    httpURL,
			SendOnly:   t == "sendonly",
		}
		return
	}
	return newNodeClient[evmtypes.Node, presenters.EVMNodeResource, EVMNodePresenter, EVMNodePresenters](c, "evm", createNode)
}

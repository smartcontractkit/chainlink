package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/urfave/cli"
	"go.uber.org/multierr"
)

type EVMForwarderPresenter struct {
	JAID // This is needed to render the id for a JSONAPI Resource as normal JSON
	presenters.EVMForwarderResource
}

var evmFwdsHeaders = []string{"ID", "Address", "Chain ID", "Created At"}

// ToRow presents the EVMForwarderResource as a slice of strings.
func (p *EVMForwarderPresenter) ToRow() []string {
	row := []string{
		p.GetID(),
		p.Address.String(),
		p.EVMChainID.ToInt().String(),
		p.CreatedAt.Format(time.RFC3339),
	}
	return row
}

// RenderTable implements TableRenderer
func (p *EVMForwarderPresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(evmFwdsHeaders, rows, rt.Writer)

	return nil
}

// EVMForwarderPresenters implements TableRenderer for a slice of EVMForwarderPresenter.
type EVMForwarderPresenters []EVMForwarderPresenter

// RenderTable implements TableRenderer
func (ps EVMForwarderPresenters) RenderTable(rt RendererTable) error {
	var rows [][]string

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(evmFwdsHeaders, rows, rt.Writer)

	return nil
}

// ListForwarders list all forwarder addresses tracked by node
func (cli *Client) ListForwarders(c *cli.Context) (err error) {
	return cli.getPage("/v2/nodes/evm/forwarders", c.Int("page"), &EVMForwarderPresenters{})
}

// DeleteForwarder deletes forwarder address from node db by id.
func (cli *Client) DeleteForwarder(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass the forwarder id to be archived"))
	}
	resp, err := cli.HTTP.Delete("/v2/nodes/evm/forwarders/" + c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}
	_, err = cli.parseResponse(resp)
	if err != nil {
		return cli.errorOut(err)
	}

	fmt.Printf("Forwarder %v Deleted\n", c.Args().First())
	return nil
}

// AddForwarder adds forwarder address to node db.
func (cli *Client) CreateForwarder(c *cli.Context) (err error) {
	addressHex := c.String("address")
	chainIDStr := c.String("evmChainID")

	addressBytes, err := hexutil.Decode(addressHex)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "could not decode address"))
	}
	address := gethCommon.BytesToAddress(addressBytes)

	var chainID *big.Int
	if chainIDStr != "" {
		var ok bool
		chainID, ok = big.NewInt(0).SetString(chainIDStr, 10)
		if !ok {
			return cli.errorOut(errors.Wrap(err, "invalid evmChainID"))
		}
	}

	request, err := json.Marshal(web.CreateEVMForwarderRequest{
		EVMChainID: (*utils.Big)(chainID),
		Address:    address,
	})
	if err != nil {
		return cli.errorOut(err)
	}

	resp, err := cli.HTTP.Post("/v2/nodes/evm/forwarders", bytes.NewReader(request))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode >= 400 {
		body, rerr := ioutil.ReadAll(resp.Body)
		if err != nil {
			err = multierr.Append(err, rerr)
			return cli.errorOut(err)
		}
		fmt.Printf("Response: '%v', Status: %d\n", string(body), resp.StatusCode)
		return cli.errorOut(err)
	}

	err = cli.renderAPIResponse(resp, &EVMForwarderPresenter{}, "Forwarder created")
	return err
}

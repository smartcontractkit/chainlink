package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/urfave/cli"
	"go.uber.org/multierr"
)

type EVMChainPresenter struct {
	presenters.EVMChainResource
}

func (p *EVMChainPresenter) ToRow() []string {
	// NOTE: it's impossible to omitempty null fields when serializing to JSON: https://github.com/golang/go/issues/11939
	config, err := json.MarshalIndent(p.Config, "", "    ")
	if err != nil {
		panic(err)
	}

	row := []string{
		p.GetID(),
		strconv.FormatBool(p.Enabled),
		string(config),
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
	}
	return row
}

// RenderTable implements TableRenderer
// Just renders a single row
func (p EVMChainPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Enabled", "Config", "Created", "Updated"}
	rows := [][]string{}
	rows = append(rows, p.ToRow())

	renderList(headers, rows, rt.Writer)

	return nil
}

type EVMChainPresenters []EVMChainPresenter

// RenderTable implements TableRenderer
func (ps EVMChainPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Enabled", "Config", "Created", "Updated"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(headers, rows, rt.Writer)

	return nil
}

// IndexEVMChains returns all EVM chains.
func (cli *Client) IndexEVMChains(c *cli.Context) (err error) {
	return cli.getPage("/v2/chains/evm", c.Int("page"), &EVMChainPresenters{})
}

// CreateEVMChain adds a new EVM chain.
func (cli *Client) CreateEVMChain(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass in the chain's parameters [-id integer] [JSON blob | JSON filepath]"))
	}
	chainID := c.Int64("id")
	if chainID == 0 {
		return cli.errorOut(errors.New("missing chain ID [-id integer]"))
	}

	buf, err := getBufferFromJSON(c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}

	params := map[string]interface{}{
		"chainID": chainID,
		"config":  json.RawMessage(buf.Bytes()),
	}

	body, err := json.Marshal(params)
	if err != nil {
		return cli.errorOut(err)
	}

	var resp *http.Response
	resp, err = cli.HTTP.Post("/v2/chains/evm", bytes.NewBuffer(body))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &EVMChainPresenter{})
}

// RemoveEVMChain removes a specific EVM Chain by id.
func (cli *Client) RemoveEVMChain(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass the id of the chain to be removed"))
	}
	chainID := c.Args().First()
	resp, err := cli.HTTP.Delete("/v2/chains/evm/" + chainID)
	if err != nil {
		return cli.errorOut(err)
	}
	_, err = cli.parseResponse(resp)
	if err != nil {
		return cli.errorOut(err)
	}

	fmt.Printf("Chain %v deleted\n", c.Args().First())
	return nil
}

func (cli *Client) ConfigureEVMChain(c *cli.Context) (err error) {
	chainID := c.Int64("id")
	if chainID == 0 {
		return cli.errorOut(errors.New("missing chain ID (usage: chainlink evm chains configure [-id integer] [key1=value1 key2=value2 ...])"))
	}

	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass in at least one chain configuration parameters (usage: chainlink evm chains configure [-id integer] [key1=value1 key2=value2 ...])"))
	}

	// Fetch existing config
	resp, err := cli.HTTP.Get(fmt.Sprintf("/v2/chains/evm/%v", chainID))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	var chain presenters.EVMChainResource
	if err = cli.deserializeAPIResponse(resp, &chain, &jsonapi.Links{}); err != nil {
		return cli.errorOut(err)
	}
	config := chain.Config

	// Parse new key-value pairs
	params := map[string]interface{}{}
	for _, arg := range c.Args() {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return cli.errorOut(errors.Errorf("invalid parameter: %v", arg))
		}

		var value interface{}
		if err = json.Unmarshal([]byte(parts[1]), &value); err != nil {
			// treat it as a string
			value = parts[1]
		}
		// TODO: handle `key=nil` and `key=` besides just null?
		params[parts[0]] = value
	}

	// Combine new values with the existing config
	// (serialize to a partial JSON map, deserialize to the old config struct)
	rawUpdates, err := json.Marshal(params)
	if err != nil {
		return cli.errorOut(err)
	}
	err = json.Unmarshal(rawUpdates, &config)
	if err != nil {
		return cli.errorOut(err)
	}

	// Send the new config
	params = map[string]interface{}{
		"enabled": chain.Enabled,
		"config":  config,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return cli.errorOut(err)
	}
	resp, err = cli.HTTP.Patch(fmt.Sprintf("/v2/chains/evm/%v", chainID), bytes.NewBuffer(body))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	return cli.renderAPIResponse(resp, &EVMChainPresenter{})
}

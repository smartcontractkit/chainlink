package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

var chainHeaders = []string{"ID", "Enabled", "Config", "Created", "Updated"}

// chainCommand returns a cli.Command with subcommands for the given ChainClient.
// The chainId cli.Flag must be named "id", but may be String or Int.
// All four types can be inferred from the ChainClient - see newChainClient for details.
func chainCommand[C chains.Config, R presenters.ChainResource[C], P TableRenderer, P2 ~[]P](
	typ string, client ChainClient[C, R, P, P2], chainID cli.Flag) cli.Command {
	if flagName := chainID.GetName(); flagName != "id" {
		panic(fmt.Errorf("chainID flag name must be 'id', got: %s", flagName))
	}
	lower := strings.ToLower(typ)
	return cli.Command{
		Name:  lower,
		Usage: fmt.Sprintf("Commands for handling %s chains", typ),
		Subcommands: cli.Commands{
			{
				Name:   "create",
				Usage:  fmt.Sprintf("Create a new %s chain", typ),
				Action: client.CreateChain,
				Flags:  []cli.Flag{chainID},
			},
			{
				Name:   "delete",
				Usage:  fmt.Sprintf("Delete an existing %s chain", typ),
				Action: client.RemoveChain,
			},
			{
				Name:   "list",
				Usage:  fmt.Sprintf("List all existing %s chains", typ),
				Action: client.IndexChains,
			},
			{
				Name:   "configure",
				Usage:  fmt.Sprintf("Configure an existing %s chain", typ),
				Action: client.ConfigureChain,
				Flags:  []cli.Flag{chainID},
			},
		},
	}
}

// ChainClient is a generic client interface for any type of chain.
type ChainClient[C chains.Config, R presenters.ChainResource[C], P TableRenderer, P2 ~[]P] interface {
	IndexChains(c *cli.Context) error
	CreateChain(c *cli.Context) error
	RemoveChain(c *cli.Context) error
	ConfigureChain(c *cli.Context) error
}

type chainClient[C chains.Config, R presenters.ChainResource[C], P TableRenderer, P2 ~[]P] struct {
	*Client
	path string
}

// newChainClient returns a new ChainClient for a particular type of chains.Config.
// P is a TableRenderer corresponding to R, and P2 is the slice variant (type P2 []P).
func newChainClient[C chains.Config, R presenters.ChainResource[C], P TableRenderer, P2 ~[]P](c *Client, name string) ChainClient[C, R, P, P2] {
	return &chainClient[C, R, P, P2]{
		Client: c,
		path:   "/v2/chains/" + name,
	}
}

// IndexChains returns all chains.
func (cli *chainClient[C, R, P, P2]) IndexChains(c *cli.Context) (err error) {
	var p P2
	return cli.getPage(cli.path, c.Int("page"), &p)
}

// CreateChain adds a new  chain.
func (cli *chainClient[C, R, P, P2]) CreateChain(c *cli.Context) (err error) {
	chainID := c.String("id")
	if chainID == "" {
		return cli.errorOut(errors.New("missing chain ID [-id]"))
	}
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass in the chain config parameters [JSON blob | JSON filepath]"))
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
	resp, err = cli.HTTP.Post(cli.path, bytes.NewBuffer(body))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var p P
	return cli.renderAPIResponse(resp, &p)
}

// RemoveChain removes a specific  Chain by id.
func (cli *chainClient[C, R, P, P2]) RemoveChain(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass the id of the chain to be removed"))
	}
	chainID := c.Args().First()
	resp, err := cli.HTTP.Delete(cli.path + "/" + chainID)
	if err != nil {
		return cli.errorOut(err)
	}
	_, err = cli.parseResponse(resp)
	if err != nil {
		return cli.errorOut(err)
	}

	fmt.Printf("Chain %v deleted\n", chainID)
	return nil
}

// ConfigureChain configures an existing  chain.
func (cli *chainClient[C, R, P, P2]) ConfigureChain(c *cli.Context) (err error) {
	chainID := c.String("id")
	if chainID == "" {
		return cli.errorOut(errors.New("missing chain ID [-id]"))
	}

	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass in at least one chain configuration parameters (usage: ... chains configure [key1=value1 key2=value2 ...])"))
	}

	// Fetch existing config
	resp, err := cli.HTTP.Get(fmt.Sprintf("%s/%v", cli.path, chainID))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	var chain R
	if err = cli.deserializeAPIResponse(resp, &chain, &jsonapi.Links{}); err != nil {
		return cli.errorOut(err)
	}
	config := chain.GetConfig()

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
		"enabled": chain.IsEnabled(),
		"config":  config,
	}
	body, err := json.Marshal(params)
	if err != nil {
		return cli.errorOut(err)
	}
	resp, err = cli.HTTP.Patch(fmt.Sprintf("%s/%v", cli.path, chainID), bytes.NewBuffer(body))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	var p P
	return cli.renderAPIResponse(resp, &p)
}

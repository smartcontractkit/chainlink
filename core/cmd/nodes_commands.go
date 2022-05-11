package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/chains"
)

type NodeResource interface {
	jsonapi.EntityNamer
}

// nodeCommand returns a cli.Command with subcommands for the given NodeClient.
// A string cli.Flag for "name" is automatically included.
// All four types can be inferred from the NodeClient - see newNodeClient for details.
func nodeCommand[N chains.Node, R NodeResource, P TableRenderer, P2 ~[]P](typ string, client NodeClient[N, R, P, P2], flags ...cli.Flag) cli.Command {
	lower := strings.ToLower(typ)
	return cli.Command{
		Name:  lower,
		Usage: fmt.Sprintf("Commands for handling %s node configuration", typ),
		Subcommands: cli.Commands{
			{
				Name:   "create",
				Usage:  fmt.Sprintf("Create a new %s node", typ),
				Action: client.CreateNode,
				Flags: append([]cli.Flag{
					cli.StringFlag{Name: "name", Usage: "node name"},
				}, flags...),
			},
			{
				Name:   "delete",
				Usage:  fmt.Sprintf("Delete an existing %s node", typ),
				Action: client.RemoveNode,
			},
			{
				Name:   "list",
				Usage:  fmt.Sprintf("List all existing %s nodes", typ),
				Action: client.IndexNodes,
			},
		},
	}
}

// NodeClient is a generic client interface for any of node.
type NodeClient[N chains.Node, R NodeResource, P TableRenderer, P2 ~[]P] interface {
	IndexNodes(c *cli.Context) error
	CreateNode(c *cli.Context) error
	RemoveNode(c *cli.Context) error
}

type nodeClient[N chains.Node, R NodeResource, P TableRenderer, P2 ~[]P] struct {
	*Client
	path       string
	createNode func(c *cli.Context) (N, error)
}

// newNodeClient returns a new NodeClient for a particular type of chains.Node.
// P is a TableRenderer corresponding to R, and P2 is the slice variant (type P2 []P).
func newNodeClient[N chains.Node, R NodeResource, P TableRenderer, P2 ~[]P](c *Client, name string, createNode func(*cli.Context) (N, error)) NodeClient[N, R, P, P2] {
	return &nodeClient[N, R, P, P2]{
		Client:     c,
		path:       "/v2/nodes/" + name,
		createNode: createNode,
	}
}

// IndexNodes returns all nodes.
func (cli *nodeClient[N, R, P, P2]) IndexNodes(c *cli.Context) (err error) {
	var p P2
	return cli.getPage(cli.path, c.Int("page"), &p)
}

// CreateNode adds a new node to the chainlink node
func (cli *nodeClient[N, R, P, P2]) CreateNode(c *cli.Context) (err error) {
	name := c.String("name")
	if name == "" {
		return cli.errorOut(errors.New("missing --name"))
	}
	newNode, err := cli.createNode(c)
	if err != nil {
		return cli.errorOut(err)
	}

	body, err := json.Marshal(newNode)
	if err != nil {
		return cli.errorOut(err)
	}

	resp, err := cli.HTTP.Post(cli.path, bytes.NewBuffer(body))
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

// RemoveNode removes a specific Node by name.
func (cli *nodeClient[N, R, P, P2]) RemoveNode(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass the id of the node to be removed"))
	}
	nodeID := c.Args().First()
	resp, err := cli.HTTP.Delete(cli.path + "/" + nodeID)
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

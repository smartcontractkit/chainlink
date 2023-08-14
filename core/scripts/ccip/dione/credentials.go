package dione

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
)

type DonCredentials struct {
	Env       Environment
	Bootstrap client.ChainlinkConfig
	Nodes     []client.ChainlinkConfig
}

func (dc *DonCredentials) WriteToFile() error {
	path := getFileLocation(dc.Env, CREDENTIALS_FOLDER)
	file, err := json.MarshalIndent(dc, "", "  ")
	if err != nil {
		return err
	}
	return WriteJSON(path+".NEW", file)
}

func (dc *DonCredentials) DialNodes() (nodes []*client.Chainlink, bootstrap *client.Chainlink, err error) {
	for _, config := range dc.Nodes {
		cfg := config
		chainlinkNode, err2 := client.NewChainlink(&cfg)
		if err2 != nil {
			return []*client.Chainlink{}, &client.Chainlink{}, err2
		}
		nodes = append(nodes, chainlinkNode)
	}

	bootstrap, err = client.NewChainlink(&dc.Bootstrap)
	if err != nil {
		return []*client.Chainlink{}, &client.Chainlink{}, err
	}

	return nodes, bootstrap, nil
}

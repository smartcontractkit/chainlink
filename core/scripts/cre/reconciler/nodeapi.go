package reconciler

import (
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
)

type clNodeClient struct{ c *clclient.ChainlinkClient }

type clNodeDialer struct{}

func (clNodeDialer) Dial(apiURL, email, password string) (NodeClient, error) {
	c, err := clclient.NewChainlinkClient(&clclient.Config{URL: apiURL, Email: email, Password: password})
	if err != nil {
		return nil, err
	}
	return &clNodeClient{c: c}, nil
}

func (n *clNodeClient) ReadCSAKey() (string, error) {
	keys, _, err := n.c.ReadCSAKeys()
	if err != nil || len(keys.Data) == 0 {
		return "", err
	}
	return strings.TrimPrefix(keys.Data[0].Attributes.PublicKey, "csa_"), nil
}

func (n *clNodeClient) ReadPeerID() (string, error) {
	keys, err := n.c.MustReadP2PKeys()
	if err != nil || len(keys.Data) == 0 {
		return "", err
	}
	return keys.Data[0].Attributes.PeerID, nil
}

func (n *clNodeClient) ReadEVMAddresses() (map[string]string, error) {
	ethKeys, err := n.c.MustReadETHKeys()
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, ek := range ethKeys.Data {
		if ek.Attributes.ChainID == "" {
			continue
		}
		result[ek.Attributes.ChainID] = ek.Attributes.Address
	}
	return result, nil
}

func (n *clNodeClient) ReadOCR2BundleIDs() (map[string]string, error) {
	ocr2Keys, err := n.c.MustReadOCR2Keys()
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, k := range ocr2Keys.Data {
		family := strings.ToLower(k.Attributes.ChainType)
		if family == "" {
			continue
		}
		result[family] = k.ID
	}
	return result, nil
}

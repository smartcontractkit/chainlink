// Package client enables interaction with APIs of test components like the mockserver and Chainlink nodes
package client

import (
	"fmt"
	"os"
	"regexp"

	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-env/environment"
)

type ChainlinkK8sClient struct {
	ChartName string
	PodName   string
	*ChainlinkClient
}

// NewChainlink creates a new Chainlink model using a provided config
func NewChainlinkK8sClient(c *ChainlinkConfig, podName, chartName string) (*ChainlinkK8sClient, error) {
	rc, err := initRestyClient(c.URL, c.Email, c.Password, c.HTTPTimeout)
	if err != nil {
		return nil, err
	}
	_, isSet := os.LookupEnv("CL_CLIENT_DEBUG")
	if isSet {
		rc.SetDebug(true)
	}
	return &ChainlinkK8sClient{
		ChainlinkClient: &ChainlinkClient{
			APIClient: rc,
			pageSize:  25,
			Config:    c,
		},
		ChartName: chartName,
		PodName:   podName,
	}, nil
}

// UpgradeVersion upgrades the chainlink node to the new version
// Note: You need to call Run() on the test environment for changes to take effect
// Note: This function is not thread safe, call from a single thread
func (c *ChainlinkK8sClient) UpgradeVersion(testEnvironment *environment.Environment, newImage, newVersion string) error {
	if newVersion == "" {
		return fmt.Errorf("new version is empty")
	}
	if newImage == "" {
		newImage = os.Getenv("CHAINLINK_IMAGE")
	}
	log.Info().
		Str("Chart Name", c.ChartName).
		Str("Old Image", os.Getenv("CHAINLINK_IMAGE")).
		Str("Old Version", os.Getenv("CHAINLINK_VERSION")).
		Str("New Image", newImage).
		Str("New Version", newVersion).
		Msg("Upgrading Chainlink Node")
	upgradeVals := map[string]any{
		"chainlink": map[string]any{
			"image": map[string]any{
				"image":   newImage,
				"version": newVersion,
			},
		},
	}
	testEnvironment, err := testEnvironment.UpdateHelm(c.ChartName, upgradeVals)
	return err
}

// Name Chainlink instance chart/service name
func (c *ChainlinkK8sClient) Name() string {
	return c.ChartName
}

func parseHostname(s string) string {
	r := regexp.MustCompile(`://(?P<Host>.*):`)
	return r.FindStringSubmatch(s)[1]
}

// ConnectChainlinkNodes creates new Chainlink clients
func ConnectChainlinkNodes(e *environment.Environment) ([]*ChainlinkK8sClient, error) {
	var clients []*ChainlinkK8sClient
	for _, nodeDetails := range e.ChainlinkNodeDetails {
		c, err := NewChainlinkK8sClient(&ChainlinkConfig{
			URL:        nodeDetails.LocalIP,
			Email:      "notreal@fakeemail.ch",
			Password:   "fj293fbBnlQ!f9vNs",
			InternalIP: parseHostname(nodeDetails.InternalIP),
		}, nodeDetails.PodName, nodeDetails.ChartName)
		if err != nil {
			return nil, err
		}
		log.Debug().
			Str("URL", c.Config.URL).
			Str("Internal IP", c.Config.InternalIP).
			Str("Chart Name", nodeDetails.ChartName).
			Str("Pod Name", nodeDetails.PodName).
			Msg("Connected to Chainlink node")
		clients = append(clients, c)
	}
	return clients, nil
}

// ReconnectChainlinkNodes reconnects to Chainlink nodes after they have been modified, say through a Helm upgrade
// Note: Experimental as of now, will likely not work predictably.
func ReconnectChainlinkNodes(testEnvironment *environment.Environment, nodes []*ChainlinkK8sClient) (err error) {
	for _, node := range nodes {
		for _, details := range testEnvironment.ChainlinkNodeDetails {
			if details.ChartName == node.ChartName { // Make the link from client to pod consistent
				node, err = NewChainlinkK8sClient(&ChainlinkConfig{
					URL:        details.LocalIP,
					Email:      "notreal@fakeemail.ch",
					Password:   "fj293fbBnlQ!f9vNs",
					InternalIP: parseHostname(details.InternalIP),
				}, details.PodName, details.ChartName)
				if err != nil {
					return err
				}
				log.Debug().
					Str("URL", node.Config.URL).
					Str("Internal IP", node.Config.InternalIP).
					Str("Chart Name", node.ChartName).
					Str("Pod Name", node.PodName).
					Msg("Reconnected to Chainlink node")
			}
		}
	}
	return nil
}

// ConnectChainlinkNodeURLs creates new Chainlink clients based on just URLs, should only be used inside K8s tests
func ConnectChainlinkNodeURLs(urls []string) ([]*ChainlinkK8sClient, error) {
	var clients []*ChainlinkK8sClient
	for _, url := range urls {
		c, err := ConnectChainlinkNodeURL(url)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}

// ConnectChainlinkNodeURL creates a new Chainlink client based on just a URL, should only be used inside K8s tests
func ConnectChainlinkNodeURL(url string) (*ChainlinkK8sClient, error) {
	return NewChainlinkK8sClient(&ChainlinkConfig{
		URL:        url,
		Email:      "notreal@fakeemail.ch",
		Password:   "fj293fbBnlQ!f9vNs",
		InternalIP: parseHostname(url),
	},
		parseHostname(url),   // a decent guess
		"connectedNodeByURL", // an intentionally bad decision
	)
}

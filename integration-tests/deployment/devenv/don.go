package devenv

import (
	"context"
	"fmt"
	"strconv"

	"github.com/AlekSi/pointer"
	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog"

	clclient "github.com/smartcontractkit/chainlink/integration-tests/client"
	csav1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/csa/v1"
	nodev1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/node/v1"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/shared/ptypes"
	"github.com/smartcontractkit/chainlink/integration-tests/web/sdk/client"
)

type NodeInfo struct {
	CLConfig    clclient.ChainlinkConfig
	IsBootstrap bool
	Name        string
	AdminAddr   string
}
type DON struct {
	Bootstrap []Node
	Nodes     []Node
}

func NewRegisteredDON(ctx context.Context, logger zerolog.Logger, nodeInfo []NodeInfo, jd JobDistributor) (*DON, error) {
	don := &DON{
		Bootstrap: make([]Node, 0),
		Nodes:     make([]Node, 0),
	}
	for i, info := range nodeInfo {
		if info.Name == "" {
			info.Name = fmt.Sprintf("node-%d", i)
		}
		node, err := NewNode(info)
		if err != nil {
			return nil, fmt.Errorf("failed to create node %d: %w", i, err)
		}
		// node Labels so that it's easier to query them
		if info.IsBootstrap {
			node.labels = append(node.labels, &ptypes.Label{
				Key:   "bootstrap",
				Value: pointer.ToString("true"),
			})
		} else {
			node.labels = append(node.labels, &ptypes.Label{
				Key:   "bootstrap",
				Value: pointer.ToString("false"),
			})
		}
		// Set up Job distributor in node
		jdId, err := node.SetUpJobDistributor(ctx, jd)
		if err != nil {
			return nil, fmt.Errorf("failed to set up job distributor in node %s: %w", info.Name, err)
		}
		if info.IsBootstrap {
			don.Bootstrap = append(don.Bootstrap, *node)
		} else {
			don.Nodes = append(don.Nodes, *node)
		}
	}
}

func NewNode(nodeInfo NodeInfo) (*Node, error) {
	gqlClient, err := client.New(nodeInfo.CLConfig.URL, client.Credentials{
		Email:    nodeInfo.CLConfig.Email,
		Password: nodeInfo.CLConfig.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create FMS client: %w", err)
	}
	return &Node{
		gqlClient: gqlClient,
		name:      nodeInfo.Name,
		adminAddr: nodeInfo.AdminAddr,
	}, nil
}

type Node struct {
	gqlClient client.Client
	NodeId    string
	labels    []*ptypes.Label
	name      string
	adminAddr string
}

func (n *Node) CreateCCIPOCR2SupportedChains(ctx context.Context, jdId string, chains []ChainConfig) error {
	var multiErr multierror.Error
	for _, chain := range chains {
		chainId := strconv.FormatUint(chain.ChainId, 10)
		accountAddr, err := n.gqlClient.FetchAccountAddress(ctx, chainId)
		if err != nil {
			multiErr.Errors = append(multiErr.Errors, err)
			continue
		}
		if accountAddr == nil {
			multiErr.Errors = append(multiErr.Errors, fmt.Errorf("no account found for chain %s", chain))
			continue
		}
		n.gqlClient.CreateJobDistributorChainConfig(ctx, client.JobDistributorChainConfigInput{
			JobDistributorID:     jdId,
			ChainID:              chainId,
			ChainType:            chain.ChainType,
			AccountAddr:          pointer.GetString(accountAddr),
			AdminAddr:            "",
			Ocr2Enabled:          true,
			Ocr2IsBootstrap:      false,
			Ocr2Multiaddr:        "",
			Ocr2ForwarderAddress: "",
			Ocr2P2PPeerID:        "",
			Ocr2KeyBundleID:      "",
			Ocr2Plugins:          `{"commit":true,"execute":true,"median":false,"mercury":false}`,
		})
	}
}

func (n *Node) SetUpJobDistributor(ctx context.Context, jd JobDistributor) (string, error) {
	// Get the public key of the node
	csaKey, err := n.gqlClient.FetchCSAPublicKey(ctx)
	if err != nil {
		return "", err
	}
	if csaKey == nil {
		return "", fmt.Errorf("no csa key found for node %s", n.name)
	}

	// register the node in the job distributor
	registerResponse, err := jd.RegisterNode(ctx, &nodev1.RegisterNodeRequest{
		PublicKey: *csaKey,
		Labels:    n.labels,
		Name:      n.name,
	})

	if err != nil {
		return "", fmt.Errorf("failed to register node %s: %w", n.name, err)
	}
	if registerResponse.GetNode().GetId() == "" {
		return "", fmt.Errorf("no node id returned from job distributor for node %s", n.name)
	}
	n.NodeId = registerResponse.GetNode().GetId()

	// Get the keypairs from the job distributor
	keypairs, err := jd.ListKeypairs(ctx, &csav1.ListKeypairsRequest{})
	if err != nil {
		return "", err
	}
	if len(keypairs.Keypairs) == 0 {
		return "", fmt.Errorf("no keypairs found from job distributor running at %s", jd.URL)
	}
	// now create the job distributor in the node
	return n.gqlClient.CreateJobDistributor(ctx, client.JobDistributorInput{
		Name:      "Job Distributor",
		Uri:       jd.URL,
		PublicKey: keypairs.Keypairs[0].PublicKey,
	})
}

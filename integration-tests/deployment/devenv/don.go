package devenv

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AlekSi/pointer"
	"github.com/rs/zerolog"

	clclient "github.com/smartcontractkit/chainlink/integration-tests/client"
	csav1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/csa/v1"
	nodev1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/node/v1"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/shared/ptypes"
	"github.com/smartcontractkit/chainlink/integration-tests/web/sdk/client"
)

type ChainInfo struct {
	ChainId   string
	ChainType string
}

type NodeInfo struct {
	CLConfig    clclient.ChainlinkConfig
	IsBootstrap bool
	Name        string
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
		node, err := NewNode(logger, info.CLConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create node %d: %w", i, err)
		}
		if info.Name == "" {
			info.Name = fmt.Sprintf("node-%d", i)
		}
		// node Labels so that it's easier to query them
		nodeLabels := make([]*ptypes.Label, 0)
		if info.IsBootstrap {
			nodeLabels = append(nodeLabels, &ptypes.Label{
				Key:   "bootstrap",
				Value: pointer.ToString("true"),
			})
		} else {
			nodeLabels = append(nodeLabels, &ptypes.Label{
				Key:   "bootstrap",
				Value: pointer.ToString("false"),
			})
		}
		// Register the node in Job distributor
		registerResponse, err := node.RegisterNodeToJobDistributor(ctx, jd.NodeServiceClient, nodeLabels, info.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to register node %w", err)
		}
		node.NodeId = registerResponse.GetId()
		if info.IsBootstrap {
			don.Bootstrap = append(don.Bootstrap, *node)
		} else {
			don.Nodes = append(don.Nodes, *node)
		}
	}
}

type Node struct {
	gqlClient client.Client
	clClient  *clclient.ChainlinkClient
	NodeId    string
}

func NewNode(logger zerolog.Logger, nodeInfo clclient.ChainlinkConfig) (*Node, error) {
	gqlClient, err := client.New(nodeInfo.URL, client.Credentials{
		Email:    nodeInfo.Email,
		Password: nodeInfo.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create FMS client: %w", err)
	}
	chainlinkClient, err := clclient.NewChainlinkClient(&nodeInfo, logger)
	if err != nil {
		return nil, err
	}

	if _, raw, err := chainlinkClient.Health(); err != nil || raw.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to connect to chainlink node: %w", err)
	}
	return &Node{
		gqlClient: gqlClient,
		clClient:  chainlinkClient,
	}, nil
}

func (n *Node) CreateChainConfig(ctx context.Context, isBootstrap bool, chainInfo ChainInfo) error {
	nodeKeyBundles, _, err := clclient.CreateNodeKeysBundle(
		[]*clclient.ChainlinkClient{n.clClient},
		chainInfo.ChainType,
		chainInfo.ChainId,
	)
	if err != nil {
		return fmt.Errorf("failed to create node key bundle: %w for node id %d", err, n.NodeId)
	}
	if nodeKeyBundles == nil || len(nodeKeyBundles) == 0 {
		return fmt.Errorf("failed to create node key bundle for node id %d", n.NodeId)
	}
	nodeKeyBundle := nodeKeyBundles[0]
	return n.gqlClient.CreateJobDistributorChainConfig(ctx, client.JobDistributorChainConfigInput{
		// TODO : check if the config is correct
		JobDistributorID:  "",
		ChainID:           chainInfo.ChainId,
		ChainType:         chainInfo.ChainType,
		AccountAddr:       nodeKeyBundle.EthAddress,
		AccountAddrPubKey: nodeKeyBundle.TXKey.Data.Attributes.PublicKey,
		AdminAddr:         nodeKeyBundle.EthAddress,
		Ocr2Enabled:       true,
		Ocr2IsBootstrap:   isBootstrap,
		Ocr2P2PPeerID:     nodeKeyBundle.PeerID,
		Ocr2KeyBundleID:   nodeKeyBundle.OCR2Key.Data.ID,
		Ocr2Plugins:       `{"commit":true,"execute":true,"median":false,"mercury":false,"rebalancer":false}`,
		Ocr2Multiaddr:     n.clClient.URL(),
	})
}

func (n *Node) CreateFeedsManager(ctx context.Context, jd JobDistributor) error {
	keypairs, err := jd.ListKeypairs(ctx, &csav1.ListKeypairsRequest{})
	if err != nil {
		return err
	}
	if len(keypairs.Keypairs) == 0 {
		return fmt.Errorf("no keypairs found from job distributor running at %s", jd.URL)
	}
	err = n.gqlClient.CreateJobDistributor(ctx, client.JobDistributorInput{
		Name:      "Job Distributor",
		Uri:       jd.URL,
		PublicKey: keypairs.Keypairs[0].PublicKey,
	})
	if err != nil {
		return fmt.Errorf("failed to create feeds manager: %w", err)
	}

	return nil
}

func (n *Node) GetCSAKeys(ctx context.Context) (*string, error) {
	nodeCSAResult, err := n.gqlClient.GetCSAKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get csa keypair for node %w", err)
	}
	if nodeCSAResult.GetCsaKeys().Results == nil || len(nodeCSAResult.GetCsaKeys().Results) == 0 {
		return nil, fmt.Errorf("failed to get csa keypair for node: %w", err)
	}
	nodeCSA := nodeCSAResult.GetCsaKeys().Results[0].GetPublicKey()
	return &nodeCSA, nil
}

func (n *Node) RegisterNodeToJobDistributor(ctx context.Context, jd nodev1.NodeServiceClient, labels []*ptypes.Label, name string) (*nodev1.Node, error) {
	nodeCSAResult, err := n.GetCSAKeys(ctx)
	if err != nil || nodeCSAResult == nil {
		return nil, fmt.Errorf("failed to get csa keypair for node %s: %w", name, err)
	}

	nodeCSA := *nodeCSAResult
	registerResponse, err := jd.RegisterNode(ctx, &nodev1.RegisterNodeRequest{
		PublicKey: nodeCSA,
		Labels:    labels,
		Name:      name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register node %s:%w", name, err)
	}
	if registerResponse.GetNode() == nil {
		return nil, fmt.Errorf("failed to register node %s returned null response", name)
	}
	return registerResponse.GetNode(), nil
}

package devenv

import (
	"context"
	"fmt"
	"strconv"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-multierror"

	clclient "github.com/smartcontractkit/chainlink/integration-tests/client"
	csav1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/csa/v1"
	nodev1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/node/v1"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/shared/ptypes"
	"github.com/smartcontractkit/chainlink/integration-tests/web/sdk/client"
)

const (
	NodeLabelType      = "type"
	NodeLabelBootstrap = "bootstrap"
	NodeLabelPlugin    = "plugin"
)

type NodeInfo struct {
	CLConfig    clclient.ChainlinkConfig
	IsBootstrap bool
	Name        string
	AdminAddr   string
}

func (info NodeInfo) Validate() error {
	var err error
	if info.CLConfig.URL == "" {
		err = multierror.Append(err, fmt.Errorf("chainlink url is required"))
	}
	if info.CLConfig.Email == "" {
		err = multierror.Append(err, fmt.Errorf("chainlink email is required"))
	}
	if info.CLConfig.Password == "" {
		err = multierror.Append(err, fmt.Errorf("chainlink password is required"))
	}
	if !info.IsBootstrap && !common.IsHexAddress(info.AdminAddr) {
		err = multierror.Append(err, fmt.Errorf("admin address is required for payment if node is not bootstrap"))
	}
	return err
}

type DON struct {
	Nodes []Node
	JDId  string
}

func (don *DON) AllNodeIds() []string {
	var nodeIds []string
	for _, node := range don.Nodes {
		nodeIds = append(nodeIds, node.NodeId)
	}
	return nodeIds
}

func (don *DON) CreateSupportedChains(ctx context.Context, chains []ChainConfig) error {
	var err error
	for _, node := range don.Nodes {
		err = multierror.Append(err, node.CreateCCIPOCR2SupportedChains(ctx, don.JDId, chains))
	}
	return err
}

func NewRegisteredDON(ctx context.Context, nodeInfo []NodeInfo, jd JobDistributor) (*DON, error) {
	don := &DON{
		Nodes: make([]Node, 0),
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
			// create multi address for OCR2, applicable only for bootstrap nodes
			node.multiAddr = info.CLConfig.URL
			// no need to set admin address for bootstrap nodes, as there will be no payment
			node.adminAddr = ""
			node.labels = append(node.labels, &ptypes.Label{
				Key:   NodeLabelType,
				Value: pointer.ToString(NodeLabelBootstrap),
			})
		} else {
			// multi address is not applicable for non-bootstrap nodes
			// explicitly set it to empty string to denote that
			node.multiAddr = ""
			node.labels = append(node.labels, &ptypes.Label{
				Key:   NodeLabelType,
				Value: pointer.ToString(NodeLabelPlugin),
			})
		}
		// Set up Job distributor in node and register node with the job distributor
		jdId, err := node.SetUpAndLinkJobDistributor(ctx, jd)
		if err != nil {
			return nil, fmt.Errorf("failed to set up job distributor in node %s: %w", info.Name, err)
		}

		don.Nodes = append(don.Nodes, *node)
		don.JDId = jdId
	}
	return don, nil
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
	NodeId    string
	gqlClient client.Client
	labels    []*ptypes.Label
	name      string
	adminAddr string
	multiAddr string
}

// CreateCCIPOCR2SupportedChains creates JobDistributorChainConfig for the node
// it works under assumption that the node is already registered with the job distributor
// expects bootstrap nodes to have type label set as bootstrap
// It fetches the account address, peer id, OCR2 key bundle id and creates the JobDistributorChainConfig
func (n *Node) CreateCCIPOCR2SupportedChains(ctx context.Context, jdId string, chains []ChainConfig) error {
	for _, chain := range chains {
		chainId := strconv.FormatUint(chain.ChainId, 10)
		accountAddr, err := n.gqlClient.FetchAccountAddress(ctx, chainId)
		if err != nil {
			return fmt.Errorf("failed to fetch account address for node %s: %w", n.name, err)
		}
		if accountAddr == nil {
			return fmt.Errorf("no account address found for node %s", n.name)
		}
		peerID, err := n.gqlClient.FetchP2PPeerID(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch peer id for node %s: %w", n.name, err)
		}
		if peerID == nil {
			return fmt.Errorf("no peer id found for node %s", n.name)
		}

		ocr2BundleId, err := n.gqlClient.FetchOCR2KeyBundleID(ctx, chain.ChainType)
		if err != nil {
			return fmt.Errorf("failed to fetch OCR2 key bundle id for node %s: %w", n.name, err)
		}
		if ocr2BundleId == "" {
			return fmt.Errorf("no OCR2 key bundle id found for node %s", n.name)
		}
		// fetch node labels to know if the node is bootstrap or plugin
		isBootstrap := false
		for _, label := range n.labels {
			if label.Key == NodeLabelType && pointer.GetString(label.Value) == NodeLabelBootstrap {
				isBootstrap = true
				break
			}
		}
		err = n.gqlClient.CreateJobDistributorChainConfig(ctx, client.JobDistributorChainConfigInput{
			JobDistributorID: jdId,
			ChainID:          chainId,
			ChainType:        chain.ChainType,
			AccountAddr:      pointer.GetString(accountAddr),
			AdminAddr:        n.adminAddr,
			Ocr2Enabled:      true,
			Ocr2IsBootstrap:  isBootstrap,
			Ocr2Multiaddr:    n.multiAddr,
			Ocr2P2PPeerID:    pointer.GetString(peerID),
			Ocr2KeyBundleID:  ocr2BundleId,
			Ocr2Plugins:      `{"commit":true,"execute":true,"median":false,"mercury":false}`,
		})
		if err != nil {
			return fmt.Errorf("failed to create CCIPOCR2SupportedChains for node %s: %w", n.name, err)
		}
	}
	return nil
}

// SetUpAndLinkJobDistributor sets up the job distributor in the node and registers the node with the job distributor
func (n *Node) SetUpAndLinkJobDistributor(ctx context.Context, jd JobDistributor) (string, error) {
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

package devenv

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-multierror"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	clclient "github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	nodev1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/node/v1"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/shared/ptypes"
	"github.com/smartcontractkit/chainlink/integration-tests/web/sdk/client"
)

const (
	NodeLabelKeyType        = "type"
	NodeLabelValueBootstrap = "bootstrap"
	NodeLabelValuePlugin    = "plugin"
)

// NodeInfo holds the information required to create a node
type NodeInfo struct {
	CLConfig    clclient.ChainlinkConfig // config to connect to chainlink node via API
	P2PPort     string                   // port for P2P communication
	IsBootstrap bool                     // denotes if the node is a bootstrap node
	Name        string                   // name of the node, used to identify the node, helpful in logs
	AdminAddr   string                   // admin address to send payments to, applicable only for non-bootstrap nodes
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
}

func (don *DON) NodeIds() []string {
	var nodeIds []string
	for _, node := range don.Nodes {
		nodeIds = append(nodeIds, node.NodeId)
	}
	return nodeIds
}

func (don *DON) FundNodes(ctx context.Context, amount *big.Int, chains map[uint64]deployment.Chain) error {
	var err error
	for sel, chain := range chains {
		for _, node := range don.Nodes {
			// if node is bootstrap, no need to fund it
			if node.multiAddr != "" {
				continue
			}
			accountAddr, ok := node.AccountAddr[sel]
			if !ok {
				err = multierror.Append(err, fmt.Errorf("node %s has no account address for chain %d", node.Name, sel))
				continue
			}
			if err1 := FundAddress(ctx, chain.DeployerKey, common.HexToAddress(accountAddr), amount, chain); err1 != nil {
				err = multierror.Append(err, err1)
			}
		}
	}
	return err
}

func (don *DON) CreateSupportedChains(ctx context.Context, chains []ChainConfig) error {
	var err error
	for i, node := range don.Nodes {
		if err1 := node.CreateCCIPOCRSupportedChains(ctx, chains); err1 != nil {
			err = multierror.Append(err, err1)
		}
		don.Nodes[i] = node
	}
	return err
}

// NewRegisteredDON creates a DON with the given node info, registers the nodes with the job distributor
// and sets up the job distributor in the nodes
func NewRegisteredDON(ctx context.Context, nodeInfo []NodeInfo, jd JobDistributor) (*DON, error) {
	don := &DON{
		Nodes: make([]Node, 0),
	}
	for i, info := range nodeInfo {
		if info.Name == "" {
			info.Name = fmt.Sprintf("node-%d", i)
		}
		if err := info.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate node info %d: %w", i, err)
		}
		node, err := NewNode(info)
		if err != nil {
			return nil, fmt.Errorf("failed to create node %d: %w", i, err)
		}
		// node Labels so that it's easier to query them
		if info.IsBootstrap {
			// create multi address for OCR2, applicable only for bootstrap nodes

			node.multiAddr = fmt.Sprintf("%s:%s", info.CLConfig.InternalIP, info.P2PPort)
			// no need to set admin address for bootstrap nodes, as there will be no payment
			node.adminAddr = ""
			node.labels = append(node.labels, &ptypes.Label{
				Key:   NodeLabelKeyType,
				Value: pointer.ToString(NodeLabelValueBootstrap),
			})
		} else {
			// multi address is not applicable for non-bootstrap nodes
			// explicitly set it to empty string to denote that
			node.multiAddr = ""
			node.labels = append(node.labels, &ptypes.Label{
				Key:   NodeLabelKeyType,
				Value: pointer.ToString(NodeLabelValuePlugin),
			})
		}
		// Set up Job distributor in node and register node with the job distributor
		err = node.SetUpAndLinkJobDistributor(ctx, jd)
		if err != nil {
			return nil, fmt.Errorf("failed to set up job distributor in node %s: %w", info.Name, err)
		}

		don.Nodes = append(don.Nodes, *node)
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
		Name:      nodeInfo.Name,
		adminAddr: nodeInfo.AdminAddr,
	}, nil
}

type Node struct {
	NodeId      string            // node id returned by job distributor after node is registered with it
	JDId        string            // job distributor id returned by node after Job distributor is created in node
	Name        string            // name of the node
	AccountAddr map[uint64]string // chain selector to node's account address mapping for supported chains
	gqlClient   client.Client     // graphql client to interact with the node
	labels      []*ptypes.Label   // labels with which the node is registered with the job distributor
	adminAddr   string            // admin address to send payments to, applicable only for non-bootstrap nodes
	multiAddr   string            // multi address denoting node's FQN (needed for deriving P2PBootstrappers in OCR), applicable only for bootstrap nodes
}

// CreateCCIPOCRSupportedChains creates a JobDistributorChainConfig for the node.
// It works under assumption that the node is already registered with the job distributor.
// It expects bootstrap nodes to have label with key "type" and value as "bootstrap".
// It fetches the account address, peer id, and OCR2 key bundle id and creates the JobDistributorChainConfig.
func (n *Node) CreateCCIPOCRSupportedChains(ctx context.Context, chains []ChainConfig) error {
	for _, chain := range chains {
		chainId := strconv.FormatUint(chain.ChainID, 10)
		selector, err := chainselectors.SelectorFromChainId(chain.ChainID)
		if err != nil {
			return fmt.Errorf("failed to get selector from chain id %d: %w", chain.ChainID, err)
		}
		accountAddr, err := n.gqlClient.FetchAccountAddress(ctx, chainId)
		if err != nil {
			return fmt.Errorf("failed to fetch account address for node %s: %w", n.Name, err)
		}
		if accountAddr == nil {
			return fmt.Errorf("no account address found for node %s", n.Name)
		}
		if n.AccountAddr == nil {
			n.AccountAddr = make(map[uint64]string)
		}
		n.AccountAddr[selector] = *accountAddr
		peerID, err := n.gqlClient.FetchP2PPeerID(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch peer id for node %s: %w", n.Name, err)
		}
		if peerID == nil {
			return fmt.Errorf("no peer id found for node %s", n.Name)
		}

		ocr2BundleId, err := n.gqlClient.FetchOCR2KeyBundleID(ctx, chain.ChainType)
		if err != nil {
			return fmt.Errorf("failed to fetch OCR2 key bundle id for node %s: %w", n.Name, err)
		}
		if ocr2BundleId == "" {
			return fmt.Errorf("no OCR2 key bundle id found for node %s", n.Name)
		}
		// fetch node labels to know if the node is bootstrap or plugin
		isBootstrap := false
		for _, label := range n.labels {
			if label.Key == NodeLabelKeyType && pointer.GetString(label.Value) == NodeLabelValueBootstrap {
				isBootstrap = true
				break
			}
		}
		err = n.gqlClient.CreateJobDistributorChainConfig(ctx, client.JobDistributorChainConfigInput{
			JobDistributorID: n.JDId,
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
			return fmt.Errorf("failed to create CCIPOCR2SupportedChains for node %s: %w", n.Name, err)
		}
	}
	return nil
}

func (n *Node) AcceptJob(ctx context.Context, id string) error {
	spec, err := n.gqlClient.ApproveJobProposalSpec(ctx, id, false)
	if err != nil {
		return err
	}
	if spec == nil {
		return fmt.Errorf("no job proposal spec found for job id %s", id)
	}
	return nil
}

// RegisterNodeToJobDistributor fetches the CSA public key of the node and registers the node with the job distributor
// it sets the node id returned by JobDistributor as a result of registration in the node struct
func (n *Node) RegisterNodeToJobDistributor(ctx context.Context, jd JobDistributor) error {
	// Get the public key of the node
	csaKeyRes, err := n.gqlClient.FetchCSAPublicKey(ctx)
	if err != nil {
		return err
	}
	if csaKeyRes == nil {
		return fmt.Errorf("no csa key found for node %s", n.Name)
	}
	csaKey := strings.TrimPrefix(*csaKeyRes, "csa_")
	// register the node in the job distributor
	registerResponse, err := jd.RegisterNode(ctx, &nodev1.RegisterNodeRequest{
		PublicKey: csaKey,
		Labels:    n.labels,
		Name:      n.Name,
	})

	if err != nil {
		return fmt.Errorf("failed to register node %s: %w", n.Name, err)
	}
	if registerResponse.GetNode().GetId() == "" {
		return fmt.Errorf("no node id returned from job distributor for node %s", n.Name)
	}
	n.NodeId = registerResponse.GetNode().GetId()
	return nil
}

// CreateJobDistributor fetches the keypairs from the job distributor and creates the job distributor in the node
// and returns the job distributor id
func (n *Node) CreateJobDistributor(ctx context.Context, jd JobDistributor) (string, error) {
	// Get the keypairs from the job distributor
	csaKey, err := jd.GetCSAPublicKey(ctx)
	if err != nil {
		return "", err
	}
	// create the job distributor in the node with the csa key
	return n.gqlClient.CreateJobDistributor(ctx, client.JobDistributorInput{
		Name:      "Job Distributor",
		Uri:       jd.WSRPC,
		PublicKey: csaKey,
	})
}

// SetUpAndLinkJobDistributor sets up the job distributor in the node and registers the node with the job distributor
// it sets the job distributor id for node
func (n *Node) SetUpAndLinkJobDistributor(ctx context.Context, jd JobDistributor) error {
	// register the node in the job distributor
	err := n.RegisterNodeToJobDistributor(ctx, jd)
	if err != nil {
		return err
	}
	// now create the job distributor in the node
	id, err := n.CreateJobDistributor(ctx, jd)
	if err != nil {
		return err
	}
	n.JDId = id
	return nil
}

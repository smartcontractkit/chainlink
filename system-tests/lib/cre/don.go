package cre

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sethvargo/go-retry"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	clclient "github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/deployment/environment/web/sdk/client"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
)

const (
	LabelNodeTypeKey            = "type"
	LabelNodeTypeValueBootstrap = "bootstrap"
	LabelNodeTypeValuePlugin    = "plugin"

	LabelNodeP2PIDKey = "p2p_id"
)

type DON struct {
	Name string `toml:"name" json:"name"`
	ID   uint64 `toml:"id" json:"id"`

	Nodes []*Node `toml:"nodes" json:"nodes"`

	Flags             []CapabilityFlag `toml:"flags" json:"flags"` // capabilities and roles
	chainCapabilities map[string]*ChainCapabilityConfig

	gh GatewayHelper
}

func (m *DON) ToMetadata() *DonMetadata {
	dm := &DonMetadata{
		Name:          m.Name,
		ID:            m.ID,
		Flags:         m.Flags,
		NodesMetadata: make([]*NodeMetadata, len(m.Nodes)),
		// caution: missing NodeSet field, since we don't have it here
	}

	for i, node := range m.Nodes {
		dm.NodesMetadata[i] = node.ToMetadata()
	}

	return dm
}

// copied from flags.go to avoid import cycle
func (m *DON) HasFlag(flag CapabilityFlag) bool {
	if slices.Contains(m.Flags, flag) {
		return true
	}

	for _, value := range m.Flags {
		if strings.HasPrefix(value, flag+"-") {
			return true
		}
	}

	return false
}

func (m *DON) ContainsBootstrapNode() bool {
	for _, node := range m.Nodes {
		if slices.Contains(node.Roles, BootstrapNode) {
			return true
		}
	}

	return false
}

func (m *DON) ContainsGatewayNode() bool {
	for _, node := range m.Nodes {
		if slices.Contains(node.Roles, GatewayNode) {
			return true
		}
	}

	return false
}

func (m *DON) GatewayNode() (*Node, error) {
	if !m.ContainsGatewayNode() {
		return nil, errors.New("don does not contain a gateway node")
	}

	for _, node := range m.Nodes {
		if slices.Contains(node.Roles, GatewayNode) {
			return node, nil
		}
	}

	return nil, errors.New("no gateway node found in don")
}

// Currently only one bootstrap node is supported.
func (m *DON) BootstrapNode() (*Node, error) {
	if !m.ContainsBootstrapNode() {
		return nil, errors.New("don does not contain a bootstrap node")
	}

	for _, node := range m.Nodes {
		if slices.Contains(node.Roles, BootstrapNode) {
			return node, nil
		}
	}

	return nil, errors.New("no bootstrap node found in don")
}

func (m *DON) WorkerNodes() ([]*Node, error) {
	workers := make([]*Node, 0)
	for _, node := range m.Nodes {
		if slices.Contains(node.Roles, WorkerNode) {
			workers = append(workers, node)
		}
	}

	if len(workers) == 0 {
		return nil, errors.New("don does not contain any worker nodes")
	}

	return workers, nil
}

func (m *DON) JDNodeIDs() []string {
	nodeIDs := []string{}
	for _, n := range m.Nodes {
		nodeIDs = append(nodeIDs, n.JobDistributorDetails.NodeID)
	}
	return nodeIDs
}

func (m *DON) NeedsAnyGateway() bool {
	return m.gh.NeedsAnyGateway(m.Flags)
}

func (m *DON) NeedsWebAPIGateway() bool {
	return m.gh.NeedsWebAPIGateway(m.Flags)
}

func (m *DON) ChainCapabilities() map[string]*ChainCapabilityConfig {
	return m.chainCapabilities
}

func NewDON(ctx context.Context, donMetadata *DonMetadata, ctfNodes []*clnode.Output, supportedChains []ChainConfig, jd *jd.JobDistributor) (*DON, error) {
	don := &DON{
		Nodes:             make([]*Node, 0),
		Name:              donMetadata.Name,
		ID:                donMetadata.ID,
		Flags:             donMetadata.Flags,
		chainCapabilities: donMetadata.ns.ChainCapabilities,
	}
	mu := &sync.Mutex{}

	errgroup := errgroup.Group{}
	for idx, nodeMetadata := range donMetadata.NodesMetadata {
		errgroup.Go(func() error {
			node, err := NewNode(ctx, fmt.Sprintf("%s-node%d", donMetadata.Name, idx), nodeMetadata, ctfNodes[idx])
			if err != nil {
				return fmt.Errorf("failed to create node %d: %w", idx, err)
			}

			// Set up Job distributor in node and register node with the job distributor
			setupErr := node.SetUpAndLinkJobDistributor(ctx, jd)
			if setupErr != nil {
				return fmt.Errorf("failed to set up job distributor in node %s: %w", node.Name, setupErr)
			}

			for _, role := range node.Roles {
				switch role {
				case WorkerNode, BootstrapNode:
					if err := node.CreateSupportedChains(ctx, supportedChains, jd); err != nil {
						return fmt.Errorf("failed to create supported chains: %w", err)
					}
				case GatewayNode:
					// no chains configuration needed for gateway nodes
				default:
					return fmt.Errorf("unknown node role: %s", role)
				}
			}

			mu.Lock()
			don.Nodes = append(don.Nodes, node)
			mu.Unlock()

			return nil
		})
	}

	if err := errgroup.Wait(); err != nil {
		return nil, fmt.Errorf("failed to create new nodes in DON: %w", err)
	}

	return don, nil
}

func (n *Node) CreateSupportedChains(ctx context.Context, chains []ChainConfig, jd *jd.JobDistributor) error {
	jdChains := []JDChainConfigInput{}
	for _, chain := range chains {
		jdChains = append(jdChains, JDChainConfigInput{
			ChainID:   chain.ChainID,
			ChainType: chain.ChainType,
		})
	}

	if err := n.CreateJDChainConfigs(ctx, jdChains, jd); err != nil {
		return fmt.Errorf("failed to create JD chain configs for node %s: %w", n.Name, err)
	}

	return nil
}

type Node struct {
	Name                  string                 `toml:"name" json:"name"`
	Host                  string                 `toml:"host" json:"host"`
	Index                 int                    `toml:"index" json:"index"`
	Keys                  *secrets.NodeKeys      `toml:"-" json:"-"`
	Addresses             Addresses              `toml:"addresses" json:"addresses"`
	JobDistributorDetails *JobDistributorDetails `toml:"job_distributor_details" json:"job_distributor_details"`
	Roles                 []string               `toml:"roles" json:"roles"`

	Clients NodeClients `toml:"-" json:"-"`
	DON     DON         `toml:"-" json:"-"`
}

func (n *Node) ToMetadata() *NodeMetadata {
	node := &NodeMetadata{
		Index: n.Index,
		Keys:  n.Keys,
		Roles: n.Roles,
		Host:  n.Host,
	}

	if node.Keys == nil {
		node.Keys = &secrets.NodeKeys{}
	}

	return node
}

func (n *Node) GetHost() string {
	return n.Host
}

// CleansedPeerID returns the PeerID without the "p2p_" prefix, or an empty string if P2PKey is nil
func (n *Node) CleansedPeerID() string {
	return n.Keys.CleansedPeerID()
}

func (n *Node) HasRole(role string) bool {
	for _, r := range n.Roles {
		if strings.EqualFold(r, role) {
			return true
		}
	}

	return false
}

func NewNode(ctx context.Context, name string, nodeMetadata *NodeMetadata, ctfNode *clnode.Output) (*Node, error) {
	gqlClient, gqErr := client.NewWithContext(ctx, ctfNode.Node.ExternalURL, client.Credentials{
		Email:    ctfNode.Node.APIAuthUser,
		Password: ctfNode.Node.APIAuthPassword,
	})
	if gqErr != nil {
		return nil, fmt.Errorf("failed to create node graphql client: %w", gqErr)
	}

	chainlinkClient, cErr := clclient.NewChainlinkClient(&clclient.ChainlinkConfig{
		URL:         ctfNode.Node.ExternalURL,
		Email:       ctfNode.Node.APIAuthUser,
		Password:    ctfNode.Node.APIAuthPassword,
		InternalIP:  ctfNode.Node.InternalIP,
		HTTPTimeout: ptr.Ptr(10 * time.Second),
	}, framework.L)
	if cErr != nil {
		return nil, fmt.Errorf("failed to create node rest client: %w", cErr)
	}

	node := &Node{
		Clients: NodeClients{
			GQLClient:  gqlClient,
			RestClient: chainlinkClient,
		},
		Name:  name,
		Index: nodeMetadata.Index,
		Keys:  nodeMetadata.Keys,
		Roles: nodeMetadata.Roles,
		Host:  nodeMetadata.Host,
	}

	for _, role := range node.Roles {
		switch role {
		case WorkerNode:
			// multi address is not applicable for non-bootstrap nodes; explicitly set it to empty string to denote that
			node.Addresses.MultiAddress = ""

			// set admin address for non-bootstrap nodes (capability registry requires non-null admin address; use arbitrary default value if node is not configured)
			node.Addresses.AdminAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
		case BootstrapNode:
			// create multi address for OCR2; applicable only for bootstrap nodes
			p2pURL, err := url.Parse(ctfNode.Node.InternalP2PUrl)
			if err != nil {
				return nil, fmt.Errorf("failed to parse p2p url: %w", err)
			}
			node.Addresses.MultiAddress = fmt.Sprintf("%s:%s", ctfNode.Node.InternalIP, p2pURL.Port())

			// no need to set admin address for bootstrap nodes, as there will be no payment
			node.Addresses.AdminAddress = ""
		case GatewayNode:
			// no specific data to set for gateway nodes yet
		default:
			return nil, fmt.Errorf("unknown node role: %s", role)
		}
	}

	return node, nil
}

type JobDistributorDetails struct {
	NodeID string `toml:"node_id" json:"node_id"` // nodeID returned by JD after node is registered with it
	JDID   string `toml:"jd_id" json:"jd_id"`     // JD ID returned by node after Job distributor is created in the node
}

type Addresses struct {
	AdminAddress string `toml:"admin_address" json:"admin_address"` // address used to pay for transactions, applicable only for worker nodes
	MultiAddress string `toml:"multi_address" json:"multi_address"` // multi address used by OCR2, applicable only for bootstrap nodes

	// maybe in the future add public addresses per chain to avoid the need to access node's keys every time?
}

type NodeClients struct {
	GQLClient  client.Client             // graphql client to interact with the node
	RestClient *clclient.ChainlinkClient // rest client to interact with the node
}

type JDChainConfigInput struct {
	ChainID   string
	ChainType string
}

func (n *Node) CreateJDChainConfigs(ctx context.Context, chains []JDChainConfigInput, jd *jd.JobDistributor) error {
	for _, chain := range chains {
		var account string

		switch strings.ToLower(chain.ChainType) {
		case chainselectors.FamilyEVM, chainselectors.FamilyTron:
			chainID, parseErr := strconv.ParseUint(chain.ChainID, 10, 64)
			if parseErr != nil {
				return fmt.Errorf("failed to parse chain id %s: %w", chain.ChainID, parseErr)
			}

			if chainID == 0 {
				return fmt.Errorf("invalid chain id: %s", chain.ChainID)
			}

			evmKey, ok := n.Keys.EVM[chainID]
			if ok {
				account = evmKey.PublicAddress.Hex()
			} else {
				var fetchErr error
				accountAddr, fetchErr := n.Clients.GQLClient.FetchAccountAddress(ctx, chain.ChainID)
				if fetchErr != nil {
					return fmt.Errorf("failed to fetch account address for node %s: %w", n.Name, fetchErr)
				}
				if accountAddr == nil {
					return fmt.Errorf("no account address found for node %s", n.Name)
				}
				account = *accountAddr
			}
		case chainselectors.FamilySolana:
			solKey, ok := n.Keys.Solana[chain.ChainID]
			if ok {
				account = solKey.PublicAddress.String()
			} else {
				accounts, fetchErr := n.Clients.GQLClient.FetchKeys(ctx, chain.ChainType)
				if fetchErr != nil {
					return fmt.Errorf("failed to fetch account address for node %s and chain %s: %w", n.Name, chain.ChainType, fetchErr)
				}
				if len(accounts) == 0 {
					return fmt.Errorf("failed to fetch account address for node %s and chain %s", n.Name, chain.ChainType)
				}
				account = accounts[0]
			}
		case chainselectors.FamilyAptos:
			// always fetch; currently Node doesn't have Aptos keys
			accounts, err := n.Clients.GQLClient.FetchKeys(ctx, chain.ChainType)
			if err != nil {
				return fmt.Errorf("failed to fetch account address for node %s and chain %s: %w", n.Name, chain.ChainType, err)
			}
			if len(accounts) == 0 {
				return fmt.Errorf("failed to fetch account address for node %s and chain %s", n.Name, chain.ChainType)
			}
			account = accounts[0]
		default:
			return fmt.Errorf("unsupported chainType %v", chain.ChainType)
		}

		chainType := chain.ChainType
		if strings.EqualFold(chain.ChainType, blockchain.FamilyTron) {
			chainType = strings.ToUpper(blockchain.FamilyEVM)
		}
		ocr2BundleID, createErr := n.Clients.GQLClient.FetchOCR2KeyBundleID(ctx, chainType)
		if createErr != nil {
			return fmt.Errorf("failed to fetch OCR2 key bundle id for node %s: %w", n.Name, createErr)
		}
		if ocr2BundleID == "" {
			return fmt.Errorf("no OCR2 key bundle id found for node %s", n.Name)
		}

		if n.Keys.OCR2BundleIDs == nil {
			n.Keys.OCR2BundleIDs = make(map[string]string)
		}

		n.Keys.OCR2BundleIDs[strings.ToLower(chainType)] = ocr2BundleID

		// retry twice with 5 seconds interval to create JobDistributorChainConfig
		retryErr := retry.Do(ctx, retry.WithMaxDuration(10*time.Second, retry.NewConstant(3*time.Second)), func(ctx context.Context) error {
			// check the node chain config to see if this chain already exists
			nodeChainConfigs, listErr := jd.ListNodeChainConfigs(context.Background(), &nodev1.ListNodeChainConfigsRequest{
				Filter: &nodev1.ListNodeChainConfigsRequest_Filter{
					NodeIds: []string{n.JobDistributorDetails.NodeID},
				}})
			if listErr != nil {
				return retry.RetryableError(fmt.Errorf("failed to list node chain configs for node %s, retrying..: %w", n.Name, listErr))
			}
			if nodeChainConfigs != nil {
				for _, chainConfig := range nodeChainConfigs.ChainConfigs {
					if chainConfig.Chain.Id == chain.ChainID {
						return nil
					}
				}
			}

			// we need to create JD chain config for each chain, because later on changestes ask the node for that chain data
			// each node needs to have OCR2 enabled, because p2pIDs are used by some contracts to identify nodes (e.g. capability registry)
			_, createErr = n.Clients.GQLClient.CreateJobDistributorChainConfig(ctx, client.JobDistributorChainConfigInput{
				JobDistributorID: n.JobDistributorDetails.JDID,
				ChainID:          chain.ChainID,
				ChainType:        chainType,
				AccountAddr:      account,
				AdminAddr:        n.Addresses.AdminAddress,
				Ocr2Enabled:      true,
				Ocr2IsBootstrap:  n.HasRole(BootstrapNode),
				Ocr2Multiaddr:    n.Addresses.MultiAddress,
				Ocr2P2PPeerID:    n.Keys.P2PKey.PeerID.String(),
				Ocr2KeyBundleID:  ocr2BundleID,
				Ocr2Plugins:      `{}`,
			})
			// TODO: add a check if the chain config failed because of a duplicate in that case, should we update or return success?
			if createErr != nil {
				return fmt.Errorf("failed to create JD chain configuration for node %s: %w", n.Name, createErr)
			}

			// JD silently fails to update nodeChainConfig. Therefore, we fetch the node config and
			// if it's not updated , throw an error
			return retry.RetryableError(errors.New("retrying CreateChainConfig in JD"))
		})

		if retryErr != nil {
			return fmt.Errorf("failed to create  create JD chain configuration for node %s: %w", n.Name, retryErr)
		}
	}
	return nil
}

// AcceptJob accepts the job proposal for the given job proposal spec
func (n *Node) AcceptJob(ctx context.Context, spec string) error {
	// fetch JD to get the job proposals
	jd, err := n.Clients.GQLClient.GetJobDistributor(ctx, n.JobDistributorDetails.JDID)
	if err != nil {
		return err
	}
	if jd.GetJobProposals() == nil {
		return fmt.Errorf("no job proposals found for node %s", n.Name)
	}
	// locate the job proposal id for the given job spec
	var idToAccept string
	for _, jp := range jd.JobProposals {
		if jp.LatestSpec.Definition == spec {
			idToAccept = jp.Id
			break
		}
	}
	if idToAccept == "" {
		return fmt.Errorf("no job proposal found for job spec %s", spec)
	}
	approvedSpec, err := n.Clients.GQLClient.ApproveJobProposalSpec(ctx, idToAccept, false)
	if err != nil {
		return err
	}
	if approvedSpec == nil {
		return fmt.Errorf("no job proposal spec found for job id %s", idToAccept)
	}
	return nil
}

// RegisterNodeToJobDistributor fetches the CSA public key of the node and registers the node with the job distributor
// it sets the node id returned by JobDistributor as a result of registration in the node struct
func (n *Node) RegisterNodeToJobDistributor(ctx context.Context, jd *jd.JobDistributor, labels []*ptypes.Label) error {
	// Get the public key of the node
	if n.Keys.CSAKey == nil {
		csaKeyRes, err := n.Clients.GQLClient.FetchCSAPublicKey(ctx)
		if err != nil {
			return err
		}
		if csaKeyRes == nil {
			return fmt.Errorf("no csa key found for node %s", n.Name)
		}

		n.Keys.CSAKey = &crypto.CSAKey{
			Key: *csaKeyRes,
		}
	}

	labels = append(labels, &ptypes.Label{
		Key:   LabelNodeP2PIDKey,
		Value: ptr.Ptr(n.Keys.P2PKey.PeerID.String()),
	})

	// register the node in the job distributor
	registerResponse, err := jd.RegisterNode(ctx, &nodev1.RegisterNodeRequest{
		PublicKey: n.Keys.CSAKey.CleansedKey(),
		Labels:    labels,
		Name:      n.Name,
	})

	if n.JobDistributorDetails == nil {
		n.JobDistributorDetails = &JobDistributorDetails{}
	}

	// node already registered, fetch it's id
	if err != nil && strings.Contains(err.Error(), "AlreadyExists") {
		nodesResponse, listErr := jd.ListNodes(ctx, &nodev1.ListNodesRequest{
			Filter: &nodev1.ListNodesRequest_Filter{
				Selectors: []*ptypes.Selector{
					{
						Key:   LabelNodeP2PIDKey,
						Op:    ptypes.SelectorOp_EQ,
						Value: ptr.Ptr(n.Keys.P2PKey.PeerID.String()),
					},
				},
			},
		})
		if listErr != nil {
			return listErr
		}
		nodes := nodesResponse.GetNodes()
		if len(nodes) == 0 {
			return fmt.Errorf("failed to find node: %v", n.Name)
		}
		n.JobDistributorDetails.NodeID = nodes[0].Id
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to register node %s: %w", n.Name, err)
	}
	if registerResponse.GetNode().GetId() == "" {
		return fmt.Errorf("no node id returned from job distributor for node %s", n.Name)
	}
	n.JobDistributorDetails.NodeID = registerResponse.GetNode().GetId()

	return nil
}

// CreateJobDistributor fetches the keypairs from the job distributor and creates the job distributor in the node
// and returns the job distributor id
func (n *Node) CreateJobDistributor(ctx context.Context, jd *jd.JobDistributor) (string, error) {
	// Get the keypairs from the job distributor
	csaKey, err := jd.GetCSAPublicKey(ctx)
	if err != nil {
		return "", err
	}
	// create the job distributor in the node with the csa key
	resp, err := n.Clients.GQLClient.ListJobDistributors(ctx)
	if err != nil {
		return "", fmt.Errorf("could not list job distributors: %w", err)
	}
	if len(resp.FeedsManagers.Results) > 0 {
		for _, fm := range resp.FeedsManagers.Results {
			if fm.GetPublicKey() == csaKey {
				return fm.GetId(), nil
			}
		}
	}
	return n.Clients.GQLClient.CreateJobDistributor(ctx, client.JobDistributorInput{
		Name:      "Job Distributor",
		Uri:       jd.WSRPC,
		PublicKey: csaKey,
	})
}

// SetUpAndLinkJobDistributor sets up the job distributor in the node and registers the node with the job distributor
// it sets the job distributor id for node
func (n *Node) SetUpAndLinkJobDistributor(ctx context.Context, jd *jd.JobDistributor) error {
	labels := make([]*ptypes.Label, 0)

	for _, role := range n.Roles {
		switch role {
		case WorkerNode:
			labels = append(labels, &ptypes.Label{
				Key:   LabelNodeTypeKey,
				Value: ptr.Ptr(LabelNodeTypeValuePlugin),
			})
		case BootstrapNode:
			labels = append(labels, &ptypes.Label{
				Key:   LabelNodeTypeKey,
				Value: ptr.Ptr(LabelNodeTypeValueBootstrap),
			})
		case GatewayNode:
			// no specific data to set for gateway nodes yet
		default:
			return fmt.Errorf("unknown node role: %s", role)
		}
	}

	// register the node in the job distributor
	err := n.RegisterNodeToJobDistributor(ctx, jd, labels)
	if err != nil {
		return err
	}
	// now create the job distributor in the node
	id, err := n.CreateJobDistributor(ctx, jd)
	if err != nil &&
		!strings.Contains(err.Error(), "DuplicateFeedsManagerError") {
		return fmt.Errorf("failed to create job distributor in node %s: %w", n.Name, err)
	}
	// wait for the node to connect to the job distributor
	err = retry.Do(ctx, retry.WithMaxDuration(1*time.Minute, retry.NewFibonacci(1*time.Second)), func(ctx context.Context) error {
		getRes, getErr := jd.GetNode(ctx, &nodev1.GetNodeRequest{
			Id: n.JobDistributorDetails.NodeID,
		})
		if getErr != nil {
			return retry.RetryableError(fmt.Errorf("failed to get node %s: %w", n.Name, getErr))
		}
		if getRes.GetNode() == nil {
			return fmt.Errorf("no node found for node id %s", n.JobDistributorDetails.NodeID)
		}
		if !getRes.GetNode().IsConnected {
			return retry.RetryableError(fmt.Errorf("node %s not connected to job distributor", n.Name))
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to connect node %s to job distributor: %w", n.Name, err)
	}
	n.JobDistributorDetails.JDID = id
	return nil
}

func (n *Node) CancelProposalsByExternalJobID(ctx context.Context, externalJobIDs []string) ([]string, error) {
	jd, err := n.Clients.GQLClient.GetJobDistributor(ctx, n.JobDistributorDetails.JDID)
	if err != nil {
		return nil, err
	}
	if jd.GetJobProposals() == nil {
		return nil, fmt.Errorf("no job proposals found for node %s", n.Name)
	}

	proposalIDs := []string{}
	for _, jp := range jd.JobProposals {
		if !slices.Contains(externalJobIDs, jp.ExternalJobID) {
			continue
		}

		proposalIDs = append(proposalIDs, jp.Id)
		spec, err := n.Clients.GQLClient.CancelJobProposalSpec(ctx, jp.Id)
		if err != nil {
			return nil, err
		}

		if spec == nil {
			return nil, fmt.Errorf("no job proposal spec found for id %s", jp.Id)
		}
	}

	return proposalIDs, nil
}

func (n *Node) ApproveProposals(ctx context.Context, proposalIDs []string) error {
	for _, proposalID := range proposalIDs {
		spec, err := n.Clients.GQLClient.ApproveJobProposalSpec(ctx, proposalID, false)
		if err != nil {
			return err
		}
		if spec == nil {
			return fmt.Errorf("no job proposal spec found for id %s", proposalID)
		}
	}
	return nil
}

func (n *Node) ExportOCR2Keys(id string) (*clclient.OCR2ExportKey, error) {
	keys, _, err := n.Clients.RestClient.ExportOCR2Key(id)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func LinkToJobDistributor(ctx context.Context, input *LinkDonsToJDInput) (*cldf.Environment, []*DON, error) {
	if input == nil {
		return nil, nil, errors.New("input is nil")
	}
	if err := input.Validate(); err != nil {
		return nil, nil, errors.Wrap(err, "input validation failed")
	}

	jdConfig := jd.JDConfig{
		GRPC:  input.JdOutput.ExternalGRPCUrl,
		WSRPC: input.JdOutput.InternalWSRPCUrl,
		Creds: insecure.NewCredentials(),
	}

	jdClient, jdErr := jd.NewJDClient(jdConfig)
	if jdErr != nil {
		return nil, nil, errors.Wrap(jdErr, "failed to create JD client")
	}

	dons := make([]*DON, len(input.NodeSetOutput))

	for idx, nodeOutput := range input.NodeSetOutput {
		supportedChains, schErr := FindDONsSupportedChains(input.Topology.DonsMetadata.List()[idx], input.BlockchainOutputs)
		if schErr != nil {
			return nil, nil, errors.Wrap(schErr, "failed to find supported chains for DON")
		}

		don, regErr := NewDON(ctx, input.Topology.DonsMetadata.List()[idx], nodeOutput.CLNodes, supportedChains, jdClient)
		if regErr != nil {
			return nil, nil, fmt.Errorf("failed to create a new DON: %w", regErr)
		}

		dons[idx] = don
	}

	var nodeIDs []string
	for _, don := range dons {
		nodeIDs = append(nodeIDs, don.JDNodeIDs()...)
	}

	input.CldfEnvironment.Offchain = jdClient
	input.CldfEnvironment.NodeIDs = nodeIDs

	return input.CldfEnvironment, dons, nil
}

// copied from flags package to avoid circular dependency
func HasFlagForAnyChain(values []string, capability string) bool {
	if slices.Contains(values, capability) {
		return true
	}

	for _, value := range values {
		if strings.HasPrefix(value, capability+"-") {
			return true
		}
	}

	return false
}

// TODO do we need to use metadata here? maybe actually some interface that both DON and metadata would implement?
func FindDONsSupportedChains(donMetadata *DonMetadata, blockchainOutputs []*WrappedBlockchainOutput) ([]ChainConfig, error) {
	chains := make([]ChainConfig, 0)

	for chainSelector, bcOut := range blockchainOutputs {
		hasEVMChainEnabled := slices.Contains(donMetadata.EVMChains(), bcOut.ChainID)
		hasSolanaWriteCapability := HasFlagForAnyChain(donMetadata.Flags, WriteSolanaCapability)
		chainIsSolana := strings.EqualFold(bcOut.BlockchainOutput.Family, chainselectors.FamilySolana)

		if !hasEVMChainEnabled && (!hasSolanaWriteCapability || !chainIsSolana) {
			continue
		}

		cfg, cfgErr := ChainConfigFromWrapped(bcOut)
		if cfgErr != nil {
			return nil, errors.Wrapf(cfgErr, "failed to build chain config for chain selector %d", chainSelector)
		}

		chains = append(chains, cfg)
	}

	return chains, nil
}

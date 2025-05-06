package devenv

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/rs/zerolog"
	"github.com/sethvargo/go-retry"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	clclient "github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/deployment/environment/web/sdk/client"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
)

const (
	LabelNodeTypeKey            = "type"
	LabelNodeTypeValueBootstrap = "bootstrap"
	LabelNodeTypeValuePlugin    = "plugin"

	LabelNodeP2PIDKey = "p2p_id"

	LabelJobTypeKey         = "jobType"
	LabelJobTypeValueLLO    = "llo"
	LabelJobTypeValueStream = "stream"

	LabelEnvironmentKey = "environment"
	LabelStreamIDKey    = "streamID"

	LabelProductKey = "product"
)

// NodeInfo holds the information required to create a node
type NodeInfo struct {
	CLConfig      clclient.ChainlinkConfig // config to connect to chainlink node via API
	P2PPort       string                   // port for P2P communication
	IsBootstrap   bool                     // denotes if the node is a bootstrap node
	Name          string                   // name of the node, used to identify the node, helpful in logs
	AdminAddr     string                   // admin address to send payments to, applicable only for non-bootstrap nodes
	MultiAddr     string                   // multi address denoting node's FQN (needed for deriving P2PBootstrappers in OCR), applicable only for bootstrap nodes
	Labels        map[string]string        // labels to use when registering the node with job distributor
	ContainerName string                   // name of Docker container
}

type DON struct {
	Nodes []Node
}

func (don *DON) PluginNodes() []Node {
	var pluginNodes []Node
	for _, node := range don.Nodes {
		for _, label := range node.labels {
			if label.Key == LabelNodeTypeKey && value(label.Value) == LabelNodeTypeValuePlugin {
				pluginNodes = append(pluginNodes, node)
			}
		}
	}
	return pluginNodes
}

// ReplayAllLogs replays all logs for the chains on all nodes for given block numbers for each chain
func (don *DON) ReplayAllLogs(blockbyChain map[uint64]uint64) error {
	for _, node := range don.Nodes {
		if err := node.ReplayLogs(blockbyChain); err != nil {
			return err
		}
	}
	return nil
}

func (don *DON) NodeIds() []string {
	var nodeIds []string
	for _, node := range don.Nodes {
		nodeIds = append(nodeIds, node.NodeID)
	}
	return nodeIds
}

func (don *DON) CreateSupportedChains(ctx context.Context, chains []ChainConfig, jd JobDistributor) error {
	g := new(errgroup.Group)
	for i := range don.Nodes {
		i := i
		g.Go(func() error {
			node := &don.Nodes[i]
			var jdChains []JDChainConfigInput
			for _, chain := range chains {
				jdChains = append(jdChains, JDChainConfigInput{
					ChainID:   strconv.FormatUint(chain.ChainID, 10),
					ChainType: chain.ChainType,
				})
			}
			if err1 := node.CreateCCIPOCRSupportedChains(ctx, jdChains, jd); err1 != nil {
				return err1
			}
			don.Nodes[i] = *node
			return nil
		})
	}
	return g.Wait()
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
		node, err := NewNode(info)
		if err != nil {
			return nil, fmt.Errorf("failed to create node %d: %w", i, err)
		}
		if info.IsBootstrap {
			// create multi address for OCR2, applicable only for bootstrap nodes
			if info.MultiAddr == "" {
				node.multiAddr = fmt.Sprintf("%s:%s", info.CLConfig.InternalIP, info.P2PPort)
			} else {
				node.multiAddr = info.MultiAddr
			}
			// no need to set admin address for bootstrap nodes, as there will be no payment
			node.adminAddr = ""
			node.labels = append(node.labels, &ptypes.Label{
				Key:   LabelNodeTypeKey,
				Value: ptr(LabelNodeTypeValueBootstrap),
			})
		} else {
			// multi address is not applicable for non-bootstrap nodes
			// explicitly set it to empty string to denote that
			node.multiAddr = ""

			// set admin address for non-bootstrap nodes
			node.adminAddr = info.AdminAddr

			// capability registry requires non-null admin address; use arbitrary default value if node is not configured
			if info.AdminAddr == "" {
				node.adminAddr = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
			}

			node.labels = append(node.labels, &ptypes.Label{
				Key:   LabelNodeTypeKey,
				Value: ptr(LabelNodeTypeValuePlugin),
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
		return nil, fmt.Errorf("failed to create node graphql client: %w", err)
	}
	chainlinkClient, err := clclient.NewChainlinkClient(&nodeInfo.CLConfig, zerolog.Logger{})
	if err != nil {
		return nil, fmt.Errorf("failed to create node rest client: %w", err)
	}
	// node Labels so that it's easier to query them
	labels := make([]*ptypes.Label, 0)
	for key, value := range nodeInfo.Labels {
		labels = append(labels, &ptypes.Label{
			Key:   key,
			Value: &value,
		})
	}
	return &Node{
		gqlClient:  gqlClient,
		restClient: chainlinkClient,
		Name:       nodeInfo.Name,
		adminAddr:  nodeInfo.AdminAddr,
		multiAddr:  nodeInfo.MultiAddr,
		labels:     labels,
	}, nil
}

type Node struct {
	NodeID          string                    // node id returned by job distributor after node is registered with it
	JDId            string                    // job distributor id returned by node after Job distributor is created in node
	Name            string                    // name of the node
	AccountAddr     map[string]string         // chain id to node's account address mapping for supported chains
	Ocr2KeyBundleID string                    // OCR2 key bundle id of the node
	gqlClient       client.Client             // graphql client to interact with the node
	restClient      *clclient.ChainlinkClient // rest client to interact with the node
	labels          []*ptypes.Label           // labels with which the node is registered with the job distributor
	adminAddr       string                    // admin address to send payments to, applicable only for non-bootstrap nodes
	multiAddr       string                    // multi address denoting node's FQN (needed for deriving P2PBootstrappers in OCR), applicable only for bootstrap nodes
}

type JDChainConfigInput struct {
	ChainID   string
	ChainType string
}

func (n *Node) Labels() []*ptypes.Label {
	return n.labels
}

func (n *Node) AddLabel(label *ptypes.Label) {
	n.labels = append(n.labels, label)
}

// CreateCCIPOCRSupportedChains creates a JobDistributorChainConfig for the node.
// It works under assumption that the node is already registered with the job distributor.
// It expects bootstrap nodes to have label with key "type" and value as "bootstrap".
// It fetches the account address, peer id, and OCR2 key bundle id and creates the JobDistributorChainConfig.
func (n *Node) CreateCCIPOCRSupportedChains(ctx context.Context, chains []JDChainConfigInput, jd JobDistributor) error {
	for _, chain := range chains {
		var account string
		switch chain.ChainType {
		case "EVM":
			accountAddr, err := n.gqlClient.FetchAccountAddress(ctx, chain.ChainID)
			if err != nil {
				return fmt.Errorf("failed to fetch account address for node %s: %w", n.Name, err)
			}
			if accountAddr == nil {
				return fmt.Errorf("no account address found for node %s", n.Name)
			}
			if n.AccountAddr == nil {
				n.AccountAddr = make(map[string]string)
			}
			n.AccountAddr[chain.ChainID] = *accountAddr
			account = *accountAddr
		case "APTOS", "SOLANA":
			accounts, err := n.gqlClient.FetchKeys(ctx, chain.ChainType)
			if err != nil {
				return fmt.Errorf("failed to fetch account address for node %s: %w", n.Name, err)
			}
			if len(accounts) == 0 {
				return fmt.Errorf("no account address found for node %s", n.Name)
			}

			account = accounts[0]
		default:
			return fmt.Errorf("unsupported chainType %v", chain.ChainType)
		}

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
		n.Ocr2KeyBundleID = ocr2BundleId

		// fetch node labels to know if the node is bootstrap or plugin
		// if multi address is set, then it's a bootstrap node
		isBootstrap := n.multiAddr != ""
		for _, label := range n.labels {
			if label.Key == LabelNodeTypeKey && value(label.Value) == LabelNodeTypeValueBootstrap {
				isBootstrap = true
				break
			}
		}

		// retry twice with 5 seconds interval to create JobDistributorChainConfig
		err = retry.Do(ctx, retry.WithMaxDuration(10*time.Second, retry.NewConstant(3*time.Second)), func(ctx context.Context) error {
			// check the node chain config to see if this chain already exists
			nodeChainConfigs, err := jd.ListNodeChainConfigs(context.Background(), &nodev1.ListNodeChainConfigsRequest{
				Filter: &nodev1.ListNodeChainConfigsRequest_Filter{
					NodeIds: []string{n.NodeID},
				}})
			if err != nil {
				return retry.RetryableError(fmt.Errorf("failed to list node chain configs for node %s, retrying..: %w", n.Name, err))
			}
			if nodeChainConfigs != nil {
				for _, chainConfig := range nodeChainConfigs.ChainConfigs {
					if chainConfig.Chain.Id == chain.ChainID {
						return nil
					}
				}
			}

			// JD silently fails to update nodeChainConfig. Therefore, we fetch the node config and
			// if it's not updated , throw an error
			_, err = n.gqlClient.CreateJobDistributorChainConfig(ctx, client.JobDistributorChainConfigInput{
				JobDistributorID: n.JDId,
				ChainID:          chain.ChainID,
				ChainType:        chain.ChainType,
				AccountAddr:      account,
				AdminAddr:        n.adminAddr,
				Ocr2Enabled:      true,
				Ocr2IsBootstrap:  isBootstrap,
				Ocr2Multiaddr:    n.multiAddr,
				Ocr2P2PPeerID:    value(peerID),
				Ocr2KeyBundleID:  ocr2BundleId,
				Ocr2Plugins:      `{"commit":true,"execute":true,"median":false,"mercury":false}`,
			})
			// todo: add a check if the chain config failed because of a duplicate in that case, should we update or return success?
			if err != nil {
				return fmt.Errorf("failed to create CCIPOCR2SupportedChains for node %s: %w", n.Name, err)
			}

			return retry.RetryableError(errors.New("retrying CreateChainConfig in JD"))
		})

		if err != nil {
			return fmt.Errorf("failed to create CCIPOCR2SupportedChains for node %s: %w", n.Name, err)
		}
	}
	return nil
}

// AcceptJob accepts the job proposal for the given job proposal spec
func (n *Node) AcceptJob(ctx context.Context, spec string) error {
	// fetch JD to get the job proposals
	jd, err := n.gqlClient.GetJobDistributor(ctx, n.JDId)
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
	approvedSpec, err := n.gqlClient.ApproveJobProposalSpec(ctx, idToAccept, false)
	if err != nil {
		return err
	}
	if approvedSpec == nil {
		return fmt.Errorf("no job proposal spec found for job id %s", idToAccept)
	}
	return nil
}

func (n *Node) DeleteJob(ctx context.Context, jobID string) error {
	jobs, err := n.gqlClient.ListJobs(ctx, 0, 1000)
	if err != nil {
		return err
	}
	if jobs == nil || len(jobs.Jobs.Results) == 0 {
		return fmt.Errorf("no jobs found for node %s", n.Name)
	}
	for _, job := range jobs.Jobs.Results {
		if job.ExternalJobID == jobID {
			spec, err := n.gqlClient.CancelJobProposalSpec(ctx, job.Id)
			if err != nil {
				return err
			}
			if spec == nil {
				return fmt.Errorf("on deletion response no job proposal spec found for job id %s", jobID)
			}
			return nil
		}
	}
	return fmt.Errorf("no job found for job id %s", jobID)
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

	// tag nodes with p2p_id for easy lookup
	peerID, err := n.gqlClient.FetchP2PPeerID(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch peer id for node %s: %w", n.Name, err)
	}
	if peerID == nil {
		return fmt.Errorf("no peer id found for node %s", n.Name)
	}
	n.labels = append(n.labels, &ptypes.Label{
		Key:   LabelNodeP2PIDKey,
		Value: peerID,
	})

	// register the node in the job distributor
	registerResponse, err := jd.RegisterNode(ctx, &nodev1.RegisterNodeRequest{
		PublicKey: csaKey,
		Labels:    n.labels,
		Name:      n.Name,
	})
	// node already registered, fetch it's id
	if err != nil && strings.Contains(err.Error(), "AlreadyExists") {
		nodesResponse, err := jd.ListNodes(ctx, &nodev1.ListNodesRequest{
			Filter: &nodev1.ListNodesRequest_Filter{
				Selectors: []*ptypes.Selector{
					{
						Key:   LabelNodeP2PIDKey,
						Op:    ptypes.SelectorOp_EQ,
						Value: peerID,
					},
				},
			},
		})
		if err != nil {
			return err
		}
		nodes := nodesResponse.GetNodes()
		if len(nodes) == 0 {
			return fmt.Errorf("failed to find node: %v", n.Name)
		}
		n.NodeID = nodes[0].Id
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to register node %s: %w", n.Name, err)
	}
	if registerResponse.GetNode().GetId() == "" {
		return fmt.Errorf("no node id returned from job distributor for node %s", n.Name)
	}
	n.NodeID = registerResponse.GetNode().GetId()
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
	resp, err := n.gqlClient.ListJobDistributors(ctx)
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
	if err != nil &&
		!strings.Contains(err.Error(), "DuplicateFeedsManagerError") {
		return fmt.Errorf("failed to create job distributor in node %s: %w", n.Name, err)
	}
	// wait for the node to connect to the job distributor
	err = retry.Do(ctx, retry.WithMaxDuration(1*time.Minute, retry.NewFibonacci(1*time.Second)), func(ctx context.Context) error {
		getRes, err := jd.GetNode(ctx, &nodev1.GetNodeRequest{
			Id: n.NodeID,
		})
		if err != nil {
			return retry.RetryableError(fmt.Errorf("failed to get node %s: %w", n.Name, err))
		}
		if getRes.GetNode() == nil {
			return fmt.Errorf("no node found for node id %s", n.NodeID)
		}
		if !getRes.GetNode().IsConnected {
			return retry.RetryableError(fmt.Errorf("node %s not connected to job distributor", n.Name))
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to connect node %s to job distributor: %w", n.Name, err)
	}
	n.JDId = id
	return nil
}

func (n *Node) ExportEVMKeysForChain(chainId string) ([]*clclient.ExportedEVMKey, error) {
	return n.restClient.ExportEVMKeysForChain(chainId)
}

// ReplayLogs replays logs for the chains on the node for given block numbers for each chain
func (n *Node) ReplayLogs(blockByChain map[uint64]uint64) error {
	for sel, block := range blockByChain {
		chainID, err := chainsel.ChainIdFromSelector(sel)
		if err != nil {
			return err
		}
		response, _, err := n.restClient.ReplayLogPollerFromBlock(int64(block), int64(chainID))
		if err != nil {
			return err
		}
		if response.Data.Attributes.Message != "Replay started" {
			return fmt.Errorf("unexpected response message from log poller's replay: %s", response.Data.Attributes.Message)
		}
	}
	return nil
}

func (n *Node) ExportOCR2Keys(id string) (*clclient.OCR2ExportKey, error) {
	keys, _, err := n.restClient.ExportOCR2Key(id)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func ptr[T any](v T) *T {
	return &v
}

func value[T any](v *T) T {
	zero := new(T)
	if v == nil {
		return *zero
	}
	return *v
}

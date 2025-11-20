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
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	cre_offchain "github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	offchain_ops "github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain/changeset/operations"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink/deployment/environment/web/sdk/client"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/solana"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
)

const (
	LabelNodeTypeKey            = "type"
	LabelNodeTypeValueBootstrap = "bootstrap"
	LabelNodeTypeValuePlugin    = "plugin"

	LabelNodeP2PIDKey = "p2p_id"
)

type Role string

const (
	RoleBootstrap Role = "bootstrap"
	RoleWorker    Role = "plugin" // label value used by chainlink-deployments-framework to denote worker nodes
	RoleGateway   Role = "gateway"
)

func NewRole(role string) (Role, error) {
	switch strings.ToLower(role) {
	case "bootstrap":
		return RoleBootstrap, nil
	case "worker", "plugin":
		return RoleWorker, nil
	case "gateway":
		return RoleGateway, nil
	default:
		return "", fmt.Errorf("unknown role: %s", role)
	}
}

type Roles []Role

func (r Roles) Contains(role Role) bool {
	return slices.Contains(r, role)
}

func (r Roles) Strings() []string {
	result := make([]string, len(r))
	for i, role := range r {
		result[i] = string(role)
	}

	return result
}

func MustNewRoles(roles []string) Roles {
	r, err := NewRoles(roles)
	if err != nil {
		panic(err)
	}

	return r
}

func NewRoles(roles []string) (Roles, error) {
	result := make(Roles, len(roles))
	for i, role := range roles {
		r, err := NewRole(role)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}

	return result, nil
}

type Don struct {
	Name string `toml:"name" json:"name"`
	ID   uint64 `toml:"id" json:"id"`
	F    uint8  `toml:"f" json:"f"` // max faulty nodes

	Nodes []*Node `toml:"nodes" json:"nodes"`

	Flags                     []CapabilityFlag `toml:"flags" json:"flags"` // capabilities and roles
	chainCapabilityConfigs    map[string]*ChainCapabilityConfig
	capabilityConfigOverrides map[string]map[string]any
}

func (d *Don) Metadata() *DonMetadata {
	dm := &DonMetadata{
		Name:          d.Name,
		ID:            d.ID,
		Flags:         d.Flags,
		NodesMetadata: make([]*NodeMetadata, len(d.Nodes)),
		// caution: missing NodeSet field, since we don't have it here
	}

	for i, node := range d.Nodes {
		dm.NodesMetadata[i] = node.Metadata()
	}

	return dm
}

// copied from flags.go to avoid import cycle
func (d *Don) HasFlag(flag CapabilityFlag) bool {
	if slices.Contains(d.Flags, flag) {
		return true
	}

	for _, value := range d.Flags {
		if strings.HasPrefix(value, flag+"-") {
			return true
		}
	}

	return false
}

func (d *Don) Gateway() (*Node, bool) {
	for _, node := range d.Nodes {
		if node.Roles.Contains(RoleGateway) {
			return node, true
		}
	}

	return nil, false
}

// Currently only one bootstrap node is supported.
func (d *Don) Bootstrap() (*Node, bool) {
	for _, node := range d.Nodes {
		if node.Roles.Contains(RoleBootstrap) {
			return node, true
		}
	}

	return nil, false
}

func (d *Don) WorkersCount() int {
	workers, wErr := d.Workers()
	if wErr != nil {
		return 0
	}

	return len(workers)
}

func (d *Don) Workers() ([]*Node, error) {
	workers := make([]*Node, 0)
	for _, node := range d.Nodes {
		if node.Roles.Contains(RoleWorker) {
			workers = append(workers, node)
		}
	}

	if len(workers) == 0 {
		return nil, errors.New("don does not contain any worker nodes")
	}

	return workers, nil
}

func (d *Don) JDNodeIDs() []string {
	nodeIDs := []string{}
	for _, n := range d.Nodes {
		nodeIDs = append(nodeIDs, n.JobDistributorDetails.NodeID)
	}
	return nodeIDs
}

func (d *Don) GetChainCapabilityConfigs() map[string]*ChainCapabilityConfig {
	return d.chainCapabilityConfigs
}

func (d *Don) GetCapabilityConfigOverrides() map[string]map[string]any {
	return d.capabilityConfigOverrides
}

func (d *Don) GetCapabilityFlags() []string {
	return d.Flags
}

func (d *Don) GetName() string {
	return d.Name
}

func NewDON(ctx context.Context, donMetadata *DonMetadata, ctfNodes []*clnode.Output) (*Don, error) {
	don := &Don{
		Nodes:                     make([]*Node, 0),
		Name:                      donMetadata.Name,
		ID:                        donMetadata.ID,
		Flags:                     donMetadata.Flags,
		chainCapabilityConfigs:    donMetadata.ns.ChainCapabilities,
		capabilityConfigOverrides: donMetadata.ns.CapabilityOverrides,
	}

	mu := &sync.Mutex{}

	errgroup := errgroup.Group{}
	for idx, nodeMetadata := range donMetadata.NodesMetadata {
		errgroup.Go(func() error {
			node, err := NewNode(ctx, fmt.Sprintf("%s-node%d", donMetadata.Name, idx), nodeMetadata, ctfNodes[idx])
			if err != nil {
				return fmt.Errorf("failed to create node %d: %w", idx, err)
			}
			node.DON = don

			mu.Lock()
			don.Nodes = append(don.Nodes, node)
			mu.Unlock()

			return nil
		})
	}

	if err := errgroup.Wait(); err != nil {
		return nil, fmt.Errorf("failed to create new nodes in DON: %w", err)
	}

	forwarderF := (don.WorkersCount() - 1) / 3
	if forwarderF == 0 {
		if don.HasFlag(ConsensusCapability) || don.HasFlag(ConsensusCapabilityV2) {
			return nil, fmt.Errorf("incorrect number of worker nodes: %d. Resulting F must conform to formula: mod((N-1)/3) > 0", don.WorkersCount())
		}
		// for other capabilities, we can use 1 as F
		forwarderF = 1
	}

	don.F = uint8(forwarderF) //nolint:gosec //will never happen, we don't use more than 31 nodes

	return don, nil
}

func registerWithJD(ctx context.Context, d *Don, supportedChains []blockchains.Blockchain, cldfEnv *cldf.Environment) error {
	mu := &sync.Mutex{}

	jd, ok := cldfEnv.Offchain.(*jd.JobDistributor)
	if !ok {
		return fmt.Errorf("offchain environment is not a *.jd.JobDistributor, but %T", cldfEnv.Offchain)
	}

	errgroup := errgroup.Group{}
	for idx, node := range d.Nodes {
		errgroup.Go(func() error {
			// Set up Job distributor in node and register node with the job distributor
			setupErr := node.setUpAndLinkJobDistributor(ctx, cldfEnv)
			if setupErr != nil {
				return fmt.Errorf("failed to set up job distributor in node %s: %w", node.Name, setupErr)
			}

			for _, role := range node.Roles {
				switch role {
				case RoleWorker, RoleBootstrap:
					if err := createJDChainConfigs(ctx, node, supportedChains, jd); err != nil {
						return fmt.Errorf("failed to create supported chains in node %s: %w", node.Name, err)
					}
				case RoleGateway:
					// no chains configuration needed for gateway nodes
				default:
					return fmt.Errorf("unknown node role: %s", role)
				}
			}

			mu.Lock()
			d.Nodes[idx] = node
			mu.Unlock()

			return nil
		})
	}

	if err := errgroup.Wait(); err != nil {
		return fmt.Errorf("failed to create new nodes in DON: %w", err)
	}

	return nil
}

type Node struct {
	Name                  string                 `toml:"name" json:"name"`
	Host                  string                 `toml:"host" json:"host"`
	Index                 int                    `toml:"index" json:"index"`
	UUID                  string                 `toml:"uuid" json:"uuid"`
	Keys                  *secrets.NodeKeys      `toml:"-" json:"-"`
	Addresses             Addresses              `toml:"addresses" json:"addresses"`
	JobDistributorDetails *JobDistributorDetails `toml:"job_distributor_details" json:"job_distributor_details"`
	Roles                 Roles                  `toml:"roles" json:"roles"`

	Clients NodeClients `toml:"-" json:"-"`
	DON     *Don        `toml:"-" json:"-"`
}

func (n *Node) Metadata() *NodeMetadata {
	node := &NodeMetadata{
		Index: n.Index,
		Keys:  n.Keys,
		Roles: n.Roles.Strings(),
		Host:  n.Host,
		UUID:  n.UUID,
	}

	if node.Keys == nil {
		node.Keys = &secrets.NodeKeys{}
	}

	return node
}

func (n *Node) GetHost() string {
	return n.Host
}

func (n *Node) PeerID() string {
	return n.Keys.PeerID()
}

func (n *Node) HasRole(role Role) bool {
	return slices.Contains(n.Roles, role)
}

func NewNode(ctx context.Context, name string, nodeMetadata *NodeMetadata, ctfNode *clnode.Output) (*Node, error) {
	gqlClient, gqErr := client.NewWithContext(ctx, ctfNode.Node.ExternalURL, client.Credentials{
		Email:    ctfNode.Node.APIAuthUser,
		Password: ctfNode.Node.APIAuthPassword,
	})
	if gqErr != nil {
		return nil, fmt.Errorf("failed to create node graphql client: %w", gqErr)
	}

	chainlinkClient, cErr := clclient.NewChainlinkClient(&clclient.Config{
		URL:         ctfNode.Node.ExternalURL,
		Email:       ctfNode.Node.APIAuthUser,
		Password:    ctfNode.Node.APIAuthPassword,
		InternalIP:  ctfNode.Node.InternalIP,
		HTTPTimeout: ptr.Ptr(10 * time.Second),
	})
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
		Roles: MustNewRoles(nodeMetadata.Roles),
		Host:  nodeMetadata.Host,
		UUID:  nodeMetadata.UUID,
	}

	for i, role := range nodeMetadata.Roles {
		r, err := NewRole(role)
		if err != nil {
			return nil, fmt.Errorf("failed to parse role %s: %w", role, err)
		}
		node.Roles[i] = r
	}

	for _, role := range node.Roles {
		switch role {
		case RoleWorker:
			// multi address is not applicable for non-bootstrap nodes; explicitly set it to empty string to denote that
			node.Addresses.MultiAddress = ""

			// set admin address for non-bootstrap nodes (capability registry requires non-null admin address; use arbitrary default value if node is not configured)
			node.Addresses.AdminAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
		case RoleBootstrap:
			// create multi address for OCR2; applicable only for bootstrap nodes
			p2pURL, err := url.Parse(ctfNode.Node.InternalP2PUrl)
			if err != nil {
				return nil, fmt.Errorf("failed to parse p2p url: %w", err)
			}
			node.Addresses.MultiAddress = fmt.Sprintf("%s:%s", ctfNode.Node.InternalIP, p2pURL.Port())

			// no need to set admin address for bootstrap nodes, as there will be no payment
			node.Addresses.AdminAddress = ""
		case RoleGateway:
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

func createJDChainConfigs(ctx context.Context, n *Node, supportedChains []blockchains.Blockchain, jd *jd.JobDistributor) error {
	for _, chain := range supportedChains {
		var account string
		chainIDStr := strconv.FormatUint(chain.ChainID(), 10)

		switch strings.ToLower(chain.ChainFamily()) {
		case chainselectors.FamilyEVM, chainselectors.FamilyTron:
			evmKey, ok := n.Keys.EVM[chain.ChainID()]
			if ok {
				account = evmKey.PublicAddress.Hex()
			} else {
				var fetchErr error
				accountAddr, fetchErr := n.Clients.GQLClient.FetchAccountAddress(ctx, chainIDStr)
				if fetchErr != nil {
					return fmt.Errorf("failed to fetch account address for node %s: %w", n.Name, fetchErr)
				}
				if accountAddr == nil {
					return fmt.Errorf("no account address found for node %s", n.Name)
				}
				account = *accountAddr
			}
		case chainselectors.FamilySolana:
			// solana chainID is a string, so we need to use it directly
			solChain := chain.(*solana.Blockchain)
			chainIDStr = solChain.SolanaChainID
			solKey, ok := n.Keys.Solana[chainIDStr]
			if ok {
				account = solKey.PublicAddress.String()
			} else {
				accounts, fetchErr := n.Clients.GQLClient.FetchKeys(ctx, strings.ToUpper(chain.ChainFamily()))
				if fetchErr != nil {
					return fmt.Errorf("failed to fetch account address for node %s and chain %s: %w", n.Name, chain.ChainFamily(), fetchErr)
				}
				if len(accounts) == 0 {
					return fmt.Errorf("failed to fetch account address for node %s and chain %s", n.Name, chain.ChainFamily())
				}
				account = accounts[0]
			}
		case chainselectors.FamilyAptos:
			// always fetch; currently Node doesn't have Aptos keys
			accounts, err := n.Clients.GQLClient.FetchKeys(ctx, strings.ToUpper(chain.ChainFamily()))
			if err != nil {
				return fmt.Errorf("failed to fetch account address for node %s and chain %s: %w", n.Name, chain.ChainFamily(), err)
			}
			if len(accounts) == 0 {
				return fmt.Errorf("failed to fetch account address for node %s and chain %s", n.Name, chain.ChainFamily())
			}
			account = accounts[0]
		default:
			return fmt.Errorf("unsupported chainType %v", chain.ChainFamily())
		}

		chainType := strings.ToUpper(chain.ChainFamily())
		if chain.IsFamily(blockchain.FamilyTron) {
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
					if chainConfig.Chain.Id == chainIDStr {
						return nil
					}
				}
			}

			// we need to create JD chain config for each chain, because later on changestes ask the node for that chain data
			// each node needs to have OCR2 enabled, because p2pIDs are used by some contracts to identify nodes (e.g. capability registry)
			_, createErr = n.Clients.GQLClient.CreateJobDistributorChainConfig(ctx, client.JobDistributorChainConfigInput{
				JobDistributorID: n.JobDistributorDetails.JDID,
				ChainID:          chainIDStr,
				ChainType:        chainType,
				AccountAddr:      account,
				AdminAddr:        n.Addresses.AdminAddress,
				Ocr2Enabled:      true,
				Ocr2IsBootstrap:  n.HasRole(RoleBootstrap),
				Ocr2Multiaddr:    n.Addresses.MultiAddress,
				Ocr2P2PPeerID:    n.Keys.P2PKey.PeerID.String(),
				Ocr2KeyBundleID:  ocr2BundleID,
				Ocr2Plugins:      `{}`,
			})
			// TODO: add a check if the chain config failed because of a duplicate in that case, should we update or return success?
			if createErr != nil {
				return createErr
			}

			// JD silently fails to update nodeChainConfig. Therefore, we fetch the node config and
			// if it's not updated , throw an error
			return retry.RetryableError(errors.New("retrying CreateChainConfig in JD"))
		})

		if retryErr != nil {
			return fmt.Errorf("failed to create JD chain configuration for node %s: %w", n.Name, retryErr)
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
func (n *Node) RegisterNodeToJobDistributor(ctx context.Context, cldfEnv *cldf.Environment) error {
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

	labels := make(map[string]string)
	labels[cre_offchain.FilterKeyDONName] = n.DON.Name

	in := offchain_ops.JDRegisterNodeOpInput{
		Domain:      cre_offchain.ProductLabel,
		Name:        n.Name,
		CSAKey:      strings.TrimPrefix(n.Keys.CSAKey.Key, "csa_"),
		P2PID:       n.PeerID(),
		DONName:     n.DON.Name,
		Labels:      labels,
		IsBootstrap: n.HasRole(RoleBootstrap),
	}

	regOut, regErr := operations.ExecuteOperation(cldfEnv.OperationsBundle, offchain_ops.JDUpsertNodeOp, offchain_ops.JDRegisterNodeOpDeps{
		Env: *cldfEnv,
	}, in)
	if regErr != nil {
		return fmt.Errorf("failed to register node in job distributor: %w", regErr)
	}

	if n.JobDistributorDetails == nil {
		n.JobDistributorDetails = &JobDistributorDetails{}
	}

	n.JobDistributorDetails.NodeID = regOut.Output.Node.ID

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

// setUpAndLinkJobDistributor sets up the job distributor in the node and registers the node with the job distributor
// it sets the job distributor id for node
func (n *Node) setUpAndLinkJobDistributor(ctx context.Context, cldfEnv *cldf.Environment) error {
	err := n.RegisterNodeToJobDistributor(ctx, cldfEnv)
	if err != nil {
		return err
	}

	jd, ok := cldfEnv.Offchain.(*jd.JobDistributor)
	if !ok {
		return fmt.Errorf("offchain environment is not a *.jd.JobDistributor, but %T", cldfEnv.Offchain)
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

func (n *Node) ExportOCR2Keys(id string) (*clclient.ExportedOCR2Key, error) {
	keys, _, err := n.Clients.RestClient.ExportOCR2Key(id)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func LinkToJobDistributor(ctx context.Context, input *LinkDonsToJDInput) error {
	if input == nil {
		return errors.New("input is nil")
	}

	var nodeIDs []string

	for idx, don := range input.Dons.List() {
		supportedChains, schErr := findDonSupportedChains(input.Topology.DonsMetadata.List()[idx], input.Blockchains)
		if schErr != nil {
			return errors.Wrap(schErr, "failed to find supported chains for DON")
		}

		if err := registerWithJD(ctx, don, supportedChains, input.CldfEnvironment); err != nil {
			return fmt.Errorf("failed to register DON with JD: %w", err)
		}
		nodeIDs = append(nodeIDs, don.JDNodeIDs()...)
	}

	input.CldfEnvironment.NodeIDs = nodeIDs

	return nil
}

// copied from flags package to avoid circular dependency
func HasFlag(values []string, capability string) bool {
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

func findDonSupportedChains(donMetadata *DonMetadata, bcs []blockchains.Blockchain) ([]blockchains.Blockchain, error) {
	chains := make([]blockchains.Blockchain, 0)

	for _, bc := range bcs {
		hasEVMChainEnabled := slices.Contains(donMetadata.EVMChains(), bc.ChainID())
		hasSolanaWriteCapability := donMetadata.HasFlag(WriteSolanaCapability)
		chainIsSolana := bc.IsFamily(chainselectors.FamilySolana)

		if !hasEVMChainEnabled && (!hasSolanaWriteCapability || !chainIsSolana) {
			continue
		}

		chains = append(chains, bc)
	}

	return chains, nil
}

// Make DonMetadata also implement it, just in case?
type KeystoneDON interface {
	KeystoneDONConfig() ks_contracts_op.ConfigureKeystoneDON
	ResolveORC3Config(config *keystone_changeset.OracleConfig) *keystone_changeset.OracleConfig
}

func (d *Don) KeystoneDONConfig() ks_contracts_op.ConfigureKeystoneDON {
	don := ks_contracts_op.ConfigureKeystoneDON{
		Name: d.Name,
	}

	for _, node := range d.Nodes {
		if node.HasRole(RoleWorker) {
			don.NodeIDs = append(don.NodeIDs, node.Keys.P2PKey.PeerID.String())
		}
	}

	return don
}

func (d *Don) ResolveORC3Config(config *keystone_changeset.OracleConfig) *keystone_changeset.OracleConfig {
	config.TransmissionSchedule = []int{d.WorkersCount()}

	return config
}

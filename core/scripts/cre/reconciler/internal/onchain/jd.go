package onchain

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	cldfjd "github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	webclient "github.com/smartcontractkit/chainlink/deployment/environment/web/sdk/client"
	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

func (d *Deployer) jdChainConfigSummary(desired *domain.DesiredState, topology *cre.Topology, state *domain.StateFile) string {
	nodes := jdChainConfigNodes(topology, state.NodeRuntime)
	var b strings.Builder
	registryChainVal, _ := desired.RegistryChain()
	fmt.Fprintf(&b, "Registry chain selector: %d\n", registryChainVal.ChainID)
	fmt.Fprintf(&b, "Nodes requiring JD chain config: %d\n", len(nodes))
	for _, n := range nodes {
		role := "worker"
		if n.nodeMeta != nil && n.nodeMeta.HasRole(cre.BootstrapNode) {
			role = "bootstrap"
		}
		fmt.Fprintf(&b, "  - %s (DON %s, %s)\n", n.nodeName, n.donName, role)
	}
	return b.String()
}

type jdChainConfigNode struct {
	donName  string
	nodeName string
	nodeMeta *cre.NodeMetadata
	donMeta  *cre.DonMetadata
}

// jdChainConfigNodes returns worker and bootstrap nodes that need JD chain configs,
// mirroring Local CRE registerWithJD behavior (gateway nodes are excluded).
func jdChainConfigNodes(topology *cre.Topology, runtime map[string]domain.NodeRuntimeInfo) []jdChainConfigNode {
	if topology == nil || topology.DonsMetadata == nil {
		return nil
	}

	var nodes []jdChainConfigNode
	seenPeer := make(map[string]struct{})

	for _, donMeta := range topology.DonsMetadata.List() {
		if flags.HasNoOtherFlags(donMeta.Flags, []string{cre.GatewayDON}) {
			continue
		}
		for _, nodeMeta := range donMeta.NodesMetadata {
			if nodeMeta.HasRole(cre.GatewayNode) {
				continue
			}
			if !nodeMeta.HasRole(cre.WorkerNode) && !nodeMeta.HasRole(cre.BootstrapNode) {
				continue
			}
			if nodeMeta.Keys == nil {
				continue
			}
			peerID := nodeMeta.Keys.PeerID()
			if peerID == "" {
				continue
			}
			if _, exists := seenPeer[peerID]; exists {
				continue
			}
			seenPeer[peerID] = struct{}{}

			nodes = append(nodes, jdChainConfigNode{
				donName:  donMeta.Name,
				nodeMeta: nodeMeta,
				donMeta:  donMeta,
				nodeName: chartNodeNameForWorker(nodeMeta, runtime),
			})
		}
	}
	return nodes
}

func chartNodeNameForWorker(worker *cre.NodeMetadata, runtime map[string]domain.NodeRuntimeInfo) string {
	if worker == nil || worker.Keys == nil {
		return ""
	}
	want := strings.TrimPrefix(worker.Keys.PeerID(), "p2p_")
	for name, info := range runtime {
		peer := strings.TrimPrefix(info.PeerID, "p2p_")
		if peer == want {
			return name
		}
	}
	return worker.Host
}

func (d *Deployer) prepareJDChainConfigsForCapReg(
	ctx context.Context,
	topology *cre.Topology,
	creEnv *cre.Environment,
	chainSelector uint64,
	desired *domain.DesiredState,
	cv *domain.ChartValues,
	state *domain.StateFile,
) error {
	d.log.Info().Uint64("chainSelector", chainSelector).Msg("Preparing JD node chain configs for CapReg")

	if creEnv.CldfEnvironment.Offchain == nil {
		return errors.New("JD client is required to prepare node chain configs for CapReg")
	}
	jd, ok := creEnv.CldfEnvironment.Offchain.(*cldfjd.JobDistributor)
	if !ok {
		return fmt.Errorf("offchain client must be a Job Distributor, got %T", creEnv.CldfEnvironment.Offchain)
	}

	nodes := jdChainConfigNodes(topology, state.NodeRuntime)

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(min(len(nodes), 8))
	for _, n := range nodes {
		g.Go(func() error {
			if n.nodeName == "" {
				return fmt.Errorf("DON %s: could not resolve chart node name for JD chain config", n.donName)
			}

			supportedChains, err := cre.FindDonSupportedChains(n.donMeta, creEnv.Blockchains)
			if err != nil {
				return errors.Wrapf(err, "DON %s node %s", n.donName, n.nodeName)
			}
			if len(supportedChains) == 0 {
				return fmt.Errorf("DON %s node %s: no supported chains", n.donName, n.nodeName)
			}

			creNode, err := d.buildCapRegCRENode(gctx, n.nodeName, n.nodeMeta, jd, desired, cv, state)
			if err != nil {
				return errors.Wrapf(err, "DON %s node %s", n.donName, n.nodeName)
			}

			if err := cre.CreateJDChainConfigs(gctx, creNode, supportedChains, jd); err != nil {
				return errors.Wrapf(err, "DON %s node %s", n.donName, n.nodeName)
			}

			role := "worker"
			if n.nodeMeta.HasRole(cre.BootstrapNode) {
				role = "bootstrap"
			}
			d.log.Info().
				Str("don", n.donName).
				Str("node", n.nodeName).
				Str("role", role).
				Int("chains", len(supportedChains)).
				Msg("JD chain configs prepared")
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	return verifyCapRegNodeInfo(creEnv.CldfEnvironment.Offchain, chainSelector, topology, state.NodeRuntime)
}

func (d *Deployer) buildCapRegCRENode(
	ctx context.Context,
	nodeName string,
	nodeMeta *cre.NodeMetadata,
	jd *cldfjd.JobDistributor,
	desired *domain.DesiredState,
	cv *domain.ChartValues,
	state *domain.StateFile,
) (*cre.Node, error) {
	runtime, ok := state.NodeRuntime[nodeName]
	if !ok {
		return nil, errors.New("no discovered runtime info")
	}
	if runtime.APIURL == "" {
		return nil, errors.New("missing API URL in runtime state")
	}
	if runtime.CSAKey == "" {
		return nil, errors.New("missing CSA key in runtime state")
	}

	apiInfo, err := d.k8s.GetNodeAPIInfo(ctx, nodeName, cv.GetNodeNamespace(nodeName))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get API credentials")
	}

	gqlClient, err := webclient.NewWithContext(ctx, runtime.APIURL, webclient.Credentials{
		Email:    apiInfo.Email,
		Password: apiInfo.Password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create GraphQL client")
	}

	nodeID, err := d.resolveJDNodeID(ctx, nodeName, runtime.CSAKey, jd, state)
	if err != nil {
		return nil, err
	}

	jdID, err := resolveNodeJobDistributorID(ctx, gqlClient, jd)
	if err != nil {
		return nil, err
	}

	node := &cre.Node{
		Name:  nodeName,
		Host:  nodeMeta.Host,
		Index: nodeMeta.Index,
		UUID:  nodeMeta.UUID,
		Keys:  nodeMeta.Keys,
		Roles: cre.MustNewRoles(nodeMeta.Roles),
		Clients: cre.NodeClients{
			GQLClient: gqlClient,
		},
		JobDistributorDetails: &cre.JobDistributorDetails{
			NodeID: nodeID,
			JDID:   jdID,
		},
	}

	if nodeMeta.HasRole(cre.BootstrapNode) {
		node.Addresses.AdminAddress = ""
		node.Addresses.MultiAddress = fmt.Sprintf(
			"%s.%s.svc.cluster.local:%d",
			nodeName,
			cv.GetNodeNamespace(nodeName),
			cre.OCRPeeringPort,
		)
	} else {
		registryChainVal, _ := desired.RegistryChain()
		registryChain := strconv.FormatUint(registryChainVal.ChainID, 10)
		if addr, ok := runtime.EVMAddress[registryChain]; ok && addr != "" {
			node.Addresses.AdminAddress = addr
		} else {
			node.Addresses.AdminAddress = FallbackDeployerAddress
		}
		node.Addresses.MultiAddress = ""
	}

	return node, nil
}

func (d *Deployer) resolveJDNodeID(ctx context.Context, nodeName, csaKey string, jd *cldfjd.JobDistributor, state *domain.StateFile) (string, error) {
	d.mu.Lock()
	if state.JDNodeIDs != nil {
		if id, ok := state.JDNodeIDs[nodeName]; ok && id != "" {
			d.mu.Unlock()
			return id, nil
		}
	}
	d.mu.Unlock()

	resp, err := jd.ListNodes(ctx, &nodev1.ListNodesRequest{
		Filter: &nodev1.ListNodesRequest_Filter{
			PublicKeys: []string{csaKey},
		},
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to list JD node for CSA key")
	}
	if len(resp.GetNodes()) == 0 {
		return "", fmt.Errorf("node not found in JD for CSA key %s", truncStr(csaKey, 16))
	}

	nodeID := resp.GetNodes()[0].GetId()
	d.mu.Lock()
	if state.JDNodeIDs == nil {
		state.JDNodeIDs = make(map[string]string)
	}
	state.JDNodeIDs[nodeName] = nodeID
	d.mu.Unlock()
	return nodeID, nil
}

func resolveNodeJobDistributorID(ctx context.Context, gql webclient.Client, jd *cldfjd.JobDistributor) (string, error) {
	csaKey, err := jd.GetCSAPublicKey(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to get JD CSA public key")
	}

	resp, err := gql.ListJobDistributors(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to list job distributors in node")
	}
	for _, fm := range resp.FeedsManagers.Results {
		if fm.GetPublicKey() == csaKey {
			return fm.GetId(), nil
		}
	}
	return "", fmt.Errorf("no job distributor linked in node matching JD CSA key %s", truncStr(csaKey, 16))
}

// truncStr truncates s to n characters, appending "..." if it was longer.
func truncStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

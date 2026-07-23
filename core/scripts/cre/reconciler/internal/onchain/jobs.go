package onchain

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	cldfjd "github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	griddleinfra "github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/gateway"
	feature_set "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/sets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

// buildDonsForJobs constructs cre.Dons with live node clients for PostEnvStartup
// and gateway job creation. Topology must already include DonMetadata with NodeSet
// references (from buildTopology).
func (d *Deployer) buildDonsForJobs(
	ctx context.Context,
	topology *cre.Topology,
	jd *cldfjd.JobDistributor,
	desired *domain.DesiredState,
	cv *domain.ChartValues,
	state *domain.StateFile,
) (*cre.Dons, error) {
	if topology == nil || topology.DonsMetadata == nil {
		return nil, errors.New("topology is required to build DONs for jobs")
	}

	var dons []*cre.Don
	for _, donMeta := range topology.DonsMetadata.List() {
		ctfNodes := make([]*clnode.Output, len(donMeta.NodesMetadata))
		for i, nodeMeta := range donMeta.NodesMetadata {
			nodeName := chartNodeNameForDonNode(donMeta, nodeMeta, state.NodeRuntime)
			if nodeName == "" {
				return nil, fmt.Errorf("DON %s: could not resolve chart node name for node index %d", donMeta.Name, i)
			}

			ctfNode, err := d.buildCLNodeOutput(ctx, nodeName, cv, state)
			if err != nil {
				return nil, errors.Wrapf(err, "DON %s node %s", donMeta.Name, nodeName)
			}
			ctfNodes[i] = ctfNode
		}

		don, err := cre.NewDON(ctx, donMeta, ctfNodes)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to build DON %s", donMeta.Name)
		}

		for i, node := range don.Nodes {
			nodeName := chartNodeNameForDonNode(donMeta, donMeta.NodesMetadata[i], state.NodeRuntime)
			enriched, enrichErr := d.buildCapRegCRENode(ctx, nodeName, donMeta.NodesMetadata[i], jd, desired, cv, state)
			if enrichErr != nil {
				return nil, errors.Wrapf(enrichErr, "DON %s node %s", donMeta.Name, nodeName)
			}
			node.JobDistributorDetails = enriched.JobDistributorDetails
			node.Addresses = enriched.Addresses
			don.Nodes[i] = node
		}

		if id, ok := state.DONIDs[donMeta.Name]; ok {
			don.ID = id
		}

		dons = append(dons, don)
	}

	return cre.NewDons(dons, topology.GatewayConnectors), nil
}

func chartNodeNameForDonNode(donMeta *cre.DonMetadata, nodeMeta *cre.NodeMetadata, runtime map[string]domain.NodeRuntimeInfo) string {
	if flags.HasNoOtherFlags(donMeta.Flags, []string{cre.GatewayDON}) {
		for name, info := range runtime {
			if info.NodeType == string(domain.RoleGateway) {
				want := strings.TrimPrefix(nodeMeta.Keys.PeerID(), "p2p_")
				peer := strings.TrimPrefix(info.PeerID, "p2p_")
				if peer == want {
					return name
				}
			}
		}
	}
	return chartNodeNameForWorker(nodeMeta, runtime)
}

func (d *Deployer) buildCLNodeOutput(ctx context.Context, nodeName string, cv *domain.ChartValues, state *domain.StateFile) (*clnode.Output, error) {
	runtime, ok := state.NodeRuntime[nodeName]
	if !ok {
		return nil, errors.New("no discovered runtime info")
	}
	if runtime.APIURL == "" {
		return nil, errors.New("missing API URL in runtime state")
	}

	apiInfo, err := d.k8s.GetNodeAPIInfo(ctx, nodeName, cv.GetNodeNamespace(nodeName))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get API credentials")
	}

	namespace := cv.GetNodeNamespace(nodeName)
	internalURL := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", nodeName, namespace, griddleinfra.NodeAPIPort)
	internalP2PURL := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", nodeName, namespace, griddleinfra.OCRPeeringPort)

	return &clnode.Output{
		UseCache: true,
		Node: &clnode.NodeOut{
			APIAuthUser:     apiInfo.Email,
			APIAuthPassword: apiInfo.Password,
			ExternalURL:     runtime.APIURL,
			InternalURL:     internalURL,
			InternalP2PUrl:  internalP2PURL,
			InternalIP:      nodeName,
		},
	}, nil
}

// runPostEnvStartup runs gateway job creation and each Feature's PostEnvStartup,
// mirroring local CRE environment.go after nodes have synced the registry.
func (d *Deployer) runPostEnvStartup(ctx context.Context, desired *domain.DesiredState, creEnv *cre.Environment, topology *cre.Topology, dons *cre.Dons) error {
	if desired.NeedsGateway() {
		d.log.Info().Msg("J1a: Creating gateway jobs")
		if err := gateway.CreateJobs(ctx, creEnv, dons, topology, topology.GatewayServiceConfigs, gateway.WhitelistConfig{
			ExtraAllowedIPsCIDR: []string{AllowAllIPsCIDR},
		}); err != nil {
			return errors.Wrap(err, "failed to create gateway jobs")
		}
	}

	features := feature_set.New()
	for _, feature := range features.List() {
		for _, don := range dons.DonsWithFlag(feature.Flag()) {
			d.log.Info().Str("feature", feature.Flag()).Str("don", don.Name).Msg("Running PostEnvStartup")
			if err := feature.PostEnvStartup(ctx, d.log, don, dons, creEnv); err != nil {
				return errors.Wrapf(err, "PostEnvStartup failed for %s on DON %s", feature.Flag(), don.Name)
			}
		}
	}

	return nil
}

func (d *Deployer) jobCreationSummary(desired *domain.DesiredState, topology *cre.Topology) string {
	var b strings.Builder
	fmt.Fprintf(&b, "DONs in topology: %d\n", len(topology.DonsMetadata.List()))
	for _, donMeta := range topology.DonsMetadata.List() {
		fmt.Fprintf(&b, "  • %s: %d nodes, caps=%v\n", donMeta.Name, len(donMeta.NodesMetadata), donMeta.Flags)
	}
	if desired.NeedsGateway() {
		fmt.Fprintf(&b, "Gateway jobs: yes (%d connector configs)\n", len(topology.GatewayServiceConfigs))
	} else {
		b.WriteString("Gateway jobs: no\n")
	}
	return b.String()
}

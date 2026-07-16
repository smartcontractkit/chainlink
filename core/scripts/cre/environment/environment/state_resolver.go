package environment

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type LocalCREStateResolver struct {
	configPath string
	cfg        *envconfig.Config
	topology   *cre.Topology // rebuilt from saved env start state; includes don_family pairing when enabled
}

func LoadLocalCREStateResolver() (*LocalCREStateResolver, error) {
	return NewLocalCREStateResolver(envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot))
}

func TryLoadLocalCREStateResolver() (*LocalCREStateResolver, error) {
	statePath := envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot)
	if _, err := os.Stat(statePath); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to stat local CRE state file")
	}

	return NewLocalCREStateResolver(statePath)
}

func NewLocalCREStateResolver(configPath string) (*LocalCREStateResolver, error) {
	cfg := &envconfig.Config{}
	if err := cfg.Load(configPath); err != nil {
		return nil, errors.Wrap(err, "failed to load local CRE state")
	}

	topology, err := cre.NewTopology(cfg.NodeSets, *cfg.Infra, toCapabilityConfigMap(cfg.CapabilityConfigs))
	if err != nil {
		return nil, errors.Wrap(err, "failed to rebuild topology from local CRE state")
	}

	return &LocalCREStateResolver{
		configPath: configPath,
		cfg:        cfg,
		topology:   topology,
	}, nil
}

func toCapabilityConfigMap(in map[string]cre.CapabilityConfig) map[cre.CapabilityFlag]cre.CapabilityConfig {
	out := make(map[cre.CapabilityFlag]cre.CapabilityConfig, len(in))
	for key, value := range in {
		out[key] = value
	}

	return out
}

func (r *LocalCREStateResolver) RegistryRPC() (string, error) {
	if len(r.cfg.Blockchains) == 0 {
		return "", errors.New("no blockchains found in local CRE state")
	}

	if r.cfg.Blockchains[0] == nil || r.cfg.Blockchains[0].Out == nil {
		return "", errors.New("registry blockchain output missing from local CRE state")
	}

	if len(r.cfg.Blockchains[0].Out.Nodes) == 0 {
		return "", errors.New("registry blockchain has no nodes in local CRE state")
	}

	return r.cfg.Blockchains[0].Out.Nodes[0].ExternalHTTPUrl, nil
}

func (r *LocalCREStateResolver) AddressRef(contractType deployment.ContractType) (*datastore.AddressRef, error) {
	addresses, err := r.cfg.GetAddresses()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get addresses from local CRE state")
	}

	for _, addrRef := range addresses {
		if datastore.ContractType(contractType) == addrRef.Type {
			return &addrRef, nil
		}
	}

	return nil, fmt.Errorf("did not find any address for %s contract", contractType)
}

// WorkflowDONMetadata returns the first workflow DON in topology order.
//
// Legacy helper for single-DON topologies. Multi-DON deploy and vault polling should use
// ResolveWorkflowDONMetadata with an explicit selector (workflow_don_resolver.go).
func (r *LocalCREStateResolver) WorkflowDONMetadata() (*cre.DonMetadata, error) {
	workflowDONs, err := r.topology.DonsMetadata.WorkflowDONs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get workflow DONs from local CRE state")
	}

	if len(workflowDONs) == 0 {
		return nil, errors.New("no workflow DON found in local CRE state")
	}

	return workflowDONs[0], nil
}

func (r *LocalCREStateResolver) WorkflowDONID() (uint32, error) {
	workflowDON, err := r.WorkflowDONMetadata()
	if err != nil {
		return 0, err
	}

	return libc.MustSafeUint32FromUint64(workflowDON.ID), nil
}

// GatewayURLForDonFamily returns the gateway URL for a specific don_family.
//
// Used when workflow deploy omits --gateway-url and encrypts vault secrets — secrets must go
// through the gateway paired to that family, not the first gateway in the topology.
func (r *LocalCREStateResolver) GatewayURLForDonFamily(donFamily string) (string, error) {
	connectors := r.topology.GatewayConnectorsForDonFamily(donFamily)
	if len(connectors.Configurations) == 0 {
		return "", fmt.Errorf("no gateway connector found for don_family %q", donFamily)
	}

	// Several gateway nodesets may share one don_family; any connector in that family works for deploy.
	return r.formatGatewayURL(connectors.Configurations[0])
}

func (r *LocalCREStateResolver) WorkflowDONFamily() (string, error) {
	workflowDON, err := r.WorkflowDONMetadata()
	if err != nil {
		return "", err
	}

	return workflowDON.DonFamily, nil
}

func (r *LocalCREStateResolver) WorkflowDONName() (string, error) {
	workflowDON, err := r.WorkflowDONMetadata()
	if err != nil {
		return "", err
	}

	return workflowDON.Name, nil
}

func (r *LocalCREStateResolver) GatewayURL() (string, error) {
	if r.topology.GatewayConnectors == nil || len(r.topology.GatewayConnectors.Configurations) == 0 {
		return "", errors.New("no gateway connectors found in local CRE state")
	}

	return r.formatGatewayURL(r.topology.GatewayConnectors.Configurations[0])
}

func (r *LocalCREStateResolver) WorkflowRegistryOutput() (*cre.WorkflowRegistryOutput, error) {
	path := envconfig.MustWorkflowRegistryStateFileAbsPath(relativePathToRepoRoot)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to read workflow registry state file")
	}

	var out cre.WorkflowRegistryOutput
	if err := toml.Unmarshal(data, &out); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal workflow registry state file")
	}

	return &out, nil
}

// WorkflowDONNodeInfo returns the shared PostgreSQL port and worker node count for the
// workflow DON as recorded in the local CRE state file. These values are used by
// waitForVaultConfigPropagation to poll each node's registry_syncer_states table.
func (r *LocalCREStateResolver) WorkflowDONNodeInfo() (dbPort int, nodeCount int, err error) {
	donMeta, err := r.WorkflowDONMetadata()
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to get workflow DON metadata")
	}

	return r.workflowDONNodeInfoFor(donMeta)
}

// formatGatewayURL builds an external gateway URL from a gateway connector config.
func (r *LocalCREStateResolver) formatGatewayURL(cfg *cre.DonGatewayConfiguration) (string, error) {
	if cfg == nil || cfg.GatewayConfiguration == nil {
		return "", errors.New("gateway connector config is nil")
	}

	var provider infra.Provider
	if r.cfg != nil && r.cfg.Infra != nil {
		provider = *r.cfg.Infra
	} else if cfg.Incoming.Host == "" {
		return "", errors.New("gateway connector has no host and infra section is missing from local CRE state")
	}

	return cfg.ExternalHTTPURL(provider), nil
}

// workflowDONNodeInfoFor returns DB port and worker count for a specific workflow DON.
// Used by waitForVaultConfigPropagation when deploy targets one DON in a multi-DON topology.
func (r *LocalCREStateResolver) workflowDONNodeInfoFor(donMeta *cre.DonMetadata) (dbPort int, nodeCount int, err error) {
	if r.cfg.Infra == nil {
		return 0, 0, errors.New("infra section is missing from local CRE state file")
	}
	if r.cfg.Infra.IsKubernetes() {
		return 0, 0, errors.New("direct DB polling is not supported for Kubernetes provider; vault config propagation requires a static wait on Kubernetes")
	}

	nodeSet, err := r.nodeSetForDON(donMeta.Name)
	if err != nil {
		return 0, 0, err
	}
	if nodeSet.DbInput == nil {
		return 0, 0, fmt.Errorf("nodeset %q has no DbInput in local CRE state", donMeta.Name)
	}

	workers, err := donMeta.Workers()
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to get workflow DON workers")
	}

	return nodeSet.DbInput.Port, len(workers), nil
}

// nodeSetForDON looks up the saved env TOML nodeset block for a workflow DON name.
func (r *LocalCREStateResolver) nodeSetForDON(donName string) (*cre.NodeSet, error) {
	for _, ns := range r.cfg.NodeSets {
		if ns.Name == donName {
			return ns, nil
		}
	}
	return nil, fmt.Errorf("no nodeset found for workflow DON %q in local CRE state", donName)
}

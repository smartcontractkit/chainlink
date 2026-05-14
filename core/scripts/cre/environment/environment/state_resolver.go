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
)

type LocalCREStateResolver struct {
	configPath string
	cfg        *envconfig.Config
	topology   *cre.Topology
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

	cfg := r.topology.GatewayConnectors.Configurations[0]
	host := cfg.Incoming.Host
	if host == "" {
		host = r.cfg.Infra.ExternalGatewayHost()
	}

	return fmt.Sprintf("%s://%s:%d%s", cfg.Incoming.Protocol, host, cfg.Incoming.ExternalPort, cfg.Incoming.Path), nil
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
	if r.cfg.Infra == nil {
		return 0, 0, errors.New("infra section is missing from local CRE state file")
	}
	if r.cfg.Infra.IsKubernetes() {
		return 0, 0, errors.New("direct DB polling is not supported for Kubernetes provider; vault config propagation requires a static wait on Kubernetes")
	}

	donMeta, err := r.WorkflowDONMetadata()
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to get workflow DON metadata")
	}

	// Find the NodeSet whose name matches the workflow DON name.
	var nodeSet *cre.NodeSet
	for _, ns := range r.cfg.NodeSets {
		if ns.Name == donMeta.Name {
			nodeSet = ns
			break
		}
	}
	if nodeSet == nil {
		return 0, 0, fmt.Errorf("no nodeset found for workflow DON %q in local CRE state", donMeta.Name)
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

package environment

import (
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"
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

func semverFromFlag(version string) (*semver.Version, error) {
	parsed, err := semver.NewVersion(version)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid contract version %q", version)
	}

	return parsed, nil
}

package environment

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pkgerrors "github.com/pkg/errors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf_deployment "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"

	"github.com/smartcontractkit/chainlink/deployment"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

const (
	ArtifactFileName = "env_artifact.json"
	NOPAdminPrefix   = "0xaadd000000000000000000000000000000"
)

type EnvArtifact struct {
	RegistryChainSelector uint64                                               `json:"home_chain_selector"`
	AddressRefs           []datastore.AddressRef                               `json:"address_refs"`
	AddressBook           map[uint64]map[string]cldf_deployment.TypeAndVersion `json:"address_book"`
	JdConfig              jd.Output                                            `json:"jd_config"`
	Nodes                 map[string]NodesArtifact                             `json:"nodes"`
	DONs                  []DonArtifact                                        `json:"dons"`
	Bootstrappers         []BootstrapNodeArtifact                              `json:"bootstrappers"`
	NOPs                  []NOPArtifact                                        `json:"nops"`
}

type NodesArtifact struct {
	Nodes map[string]SimpleNodeArtifact `json:"nodes"`
}

type SimpleNodeArtifact struct {
	Name string `json:"name"`
}

type DonArtifact struct {
	DonName        string                  `json:"don_name"`
	DonID          uint64                  `json:"don_id"`
	F              uint8                   `json:"f"`
	BootstrapNodes []string                `json:"bootstrap_nodes"`
	Capabilities   []DONCapabilityArtifact `json:"capabilities,omitempty"`
	Nodes          []FullNodeArtifact      `json:"nodes"`
}

type FullNodeArtifact struct {
	Name   string `json:"name"`
	NOP    string `json:"nop"`
	CSAKey string `json:"csa_key"`
}

type DONCapabilityArtifact struct {
	Capability capabilities_registry.CapabilitiesRegistryCapability `json:"capability"`
	Config     *DONCapabilityConfig                                 `json:"config,omitempty"`
}

type DONCapabilityConfig struct {
	*capabilitiespb.CapabilityConfig
}

func (c *DONCapabilityConfig) UnmarshalJSON(data []byte) error {
	if c.CapabilityConfig == nil {
		c.CapabilityConfig = &capabilitiespb.CapabilityConfig{}
	}

	type Alias DONCapabilityConfig
	var aux struct {
		// allow standard JSON unmarshalling of all fields except RemoteConfig
		*Alias

		// use a map to hold any nested shape: RemoteTriggerConfig/RemoteTargetConfig/RemoteExecutableConfig
		RemoteConfig map[string]json.RawMessage `json:"RemoteConfig,omitempty"`
		// use a map to hold any methods, if present, to iterate later
		MethodConfigs map[string]json.RawMessage `json:"method_configs,omitempty"`
	}

	aux.Alias = (*Alias)(c)

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.RemoteConfig != nil {
		// parse the remote config based on the key
		switch {
		case aux.RemoteConfig["RemoteTriggerConfig"] != nil:
			var rt capabilitiespb.RemoteTriggerConfig
			if err := json.Unmarshal(aux.RemoteConfig["RemoteTriggerConfig"], &rt); err != nil {
				return err
			}
			c.RemoteConfig = &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
				RemoteTriggerConfig: &rt,
			}
		case aux.RemoteConfig["RemoteTargetConfig"] != nil:
			var tgt capabilitiespb.RemoteTargetConfig
			if err := json.Unmarshal(aux.RemoteConfig["RemoteTargetConfig"], &tgt); err != nil {
				return err
			}
			c.RemoteConfig = &capabilitiespb.CapabilityConfig_RemoteTargetConfig{
				RemoteTargetConfig: &tgt,
			}
		case aux.RemoteConfig["RemoteExecutableConfig"] != nil:
			var ex capabilitiespb.RemoteExecutableConfig
			if err := json.Unmarshal(aux.RemoteConfig["RemoteExecutableConfig"], &ex); err != nil {
				return err
			}
			c.RemoteConfig = &capabilitiespb.CapabilityConfig_RemoteExecutableConfig{
				RemoteExecutableConfig: &ex,
			}
		default:
			keys := make([]string, 0, len(aux.RemoteConfig))
			for k := range aux.RemoteConfig {
				keys = append(keys, k)
			}
			return fmt.Errorf("unknown remote config type in capability config, keys: %v", keys)
		}
	}

	if aux.MethodConfigs != nil {
		methodConfigs := make(map[string]*capabilitiespb.CapabilityMethodConfig, len(aux.MethodConfigs))
		for methodName, methodConfig := range aux.MethodConfigs {
			var methodRemoteConfig map[string]json.RawMessage
			if err := json.Unmarshal(methodConfig, &methodRemoteConfig); err != nil {
				return err
			}

			var innerRemoteConfig map[string]json.RawMessage
			if err := json.Unmarshal(methodRemoteConfig["RemoteConfig"], &innerRemoteConfig); err != nil {
				return err
			}
			switch {
			case innerRemoteConfig["RemoteTriggerConfig"] != nil:
				var rt capabilitiespb.RemoteTriggerConfig
				if err := json.Unmarshal(innerRemoteConfig["RemoteTriggerConfig"], &rt); err != nil {
					return err
				}
				methodConfigs[methodName] = &capabilitiespb.CapabilityMethodConfig{
					RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteTriggerConfig{
						RemoteTriggerConfig: &rt,
					},
				}
			case innerRemoteConfig["RemoteExecutableConfig"] != nil:
				var ex capabilitiespb.RemoteExecutableConfig
				if err := json.Unmarshal(innerRemoteConfig["RemoteExecutableConfig"], &ex); err != nil {
					return err
				}
				methodConfigs[methodName] = &capabilitiespb.CapabilityMethodConfig{
					RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
						RemoteExecutableConfig: &ex,
					},
				}
			default:
				keys := make([]string, 0, len(innerRemoteConfig))
				for k := range innerRemoteConfig {
					keys = append(keys, k)
				}
				return fmt.Errorf("unknown method config type for method %s, unknown config value keys: %s", methodName, strings.Join(keys, ","))
			}
		}

		c.MethodConfigs = methodConfigs
	}

	return nil
}

type BootstrapNodeArtifact struct {
	Name       string `json:"name"`
	NOP        string `json:"nop"`
	CSAKey     string `json:"csa_key"`
	P2PID      string `json:"p2p_id"`
	OCRUrl     string `json:"ocr_url"`
	DON2DONUrl string `json:"don2d_url"`
}

type NOPArtifact struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Admin string `json:"admin"`
}

func DumpArtifact(
	absPath string,
	datastore datastore.AddressRefStore,
	addressBook cldf_deployment.AddressBook,
	jdOutput jd.Output,
	donTopology cre.DonTopology,
	offchainClient cldf_offchain.Client,
	capabilityRegistryFns []cre.CapabilityRegistryConfigFn,
	nodeSets []*cre.CapabilitiesAwareNodeSet,
) (string, error) {
	artifact, err := GenerateArtifact(datastore, addressBook, jdOutput, donTopology, offchainClient, capabilityRegistryFns, nodeSets)
	if err != nil {
		return "", pkgerrors.Wrap(err, "failed to generate environment artifact")
	}

	// Let's save the artifact to disk
	artifactPath, err := persistArtifact(absPath, artifact)
	if err != nil {
		return "", pkgerrors.Wrap(err, "failed to persist environment artifact")
	}
	return artifactPath, nil
}

func GenerateArtifact(
	ds datastore.AddressRefStore,
	addressBook cldf_deployment.AddressBook,
	jdOutput jd.Output,
	donTopology cre.DonTopology,
	offchainClient cldf_offchain.Client,
	capabilityRegistryFns []cre.CapabilityRegistryConfigFn,
	nodeSets []*cre.CapabilitiesAwareNodeSet,
) (*EnvArtifact, error) {
	var err error

	addresses, err := addressBook.Addresses()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to get addresses from address book")
	}

	addressRecords, err := ds.Fetch()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to fetch address records from datastore")
	}

	artifact := EnvArtifact{
		RegistryChainSelector: donTopology.HomeChainSelector,
		JdConfig:              jdOutput,
		AddressBook:           addresses,
		AddressRefs:           addressRecords,
		Nodes:                 make(map[string]NodesArtifact),
		DONs:                  make([]DonArtifact, 0),
		Bootstrappers:         make([]BootstrapNodeArtifact, 0),
		NOPs:                  make([]NOPArtifact, 0),
	}

	for donIdx, don := range donTopology.Dons.List() {
		donArtifact := DonArtifact{
			DonName:        don.Name,
			DonID:          don.ID,
			F:              0, // F will be calculated based on the number of worker nodes
			BootstrapNodes: make([]string, 0),
			Nodes:          make([]FullNodeArtifact, 0),
			Capabilities:   make([]DONCapabilityArtifact, 0),
		}

		workerNodes, wErr := don.Workers()
		if wErr != nil {
			return nil, pkgerrors.Wrap(wErr, "failed to find worker nodes")
		}
		donArtifact.F = libc.MustSafeUint8((len(workerNodes) - 1) / 3)

		for _, capabilityFn := range capabilityRegistryFns {
			if capabilityFn == nil {
				continue
			}

			capabilitiesFn, capabilitiesFnErr := capabilityFn(don.Flags, nodeSets[donIdx])
			if capabilitiesFnErr != nil {
				return nil, pkgerrors.Wrap(capabilitiesFnErr, "failed to get capabilities from capability registry function")
			}

			for _, capability := range capabilitiesFn {
				donArtifact.Capabilities = append(donArtifact.Capabilities, DONCapabilityArtifact{
					Capability: capabilities_registry.CapabilitiesRegistryCapability{
						Version:        capability.Capability.Version,
						LabelledName:   capability.Capability.LabelledName,
						CapabilityType: capability.Capability.CapabilityType,
					},
					Config: &DONCapabilityConfig{capability.Config},
				})
			}
		}

		nop := NOPArtifact{
			ID:    donIdx + 1, // NOP IDs start from 1
			Name:  fmt.Sprintf("NOP for %s DON", don.Name),
			Admin: fmt.Sprintf("%s%06d", NOPAdminPrefix, donIdx+1),
		}

		var nodeIDs []string
		for _, node := range donTopology.Dons.List()[donIdx].Nodes {
			nodeIDs = append(nodeIDs, node.JobDistributorDetails.NodeID)
		}

		artifact.Nodes[don.Name] = NodesArtifact{
			Nodes: make(map[string]SimpleNodeArtifact),
		}

		artifact.NOPs = append(artifact.NOPs, nop)
		artifact.DONs = append(artifact.DONs, donArtifact)

		nodeInfo, nodeInfoErr := deployment.NodeInfo(nodeIDs, offchainClient)
		if nodeInfoErr != nil {
			if !strings.Contains(nodeInfoErr.Error(), "missing node metadata") {
				return nil, pkgerrors.Wrapf(nodeInfoErr, "failed to get node info for DON %s", don.Name)
			}
			framework.L.Warn().Msgf("Metadata is missing for some nodes in DON %s: %s", don.Name, nodeInfoErr.Error())
		}

		for _, node := range nodeInfo {
			if node.IsBootstrap {
				donArtifact.BootstrapNodes = append(donArtifact.BootstrapNodes, node.Name)
				artifact.Bootstrappers = append(artifact.Bootstrappers, BootstrapNodeArtifact{
					NOP:        nop.Name,
					Name:       node.Name,
					CSAKey:     node.CSAKey,
					P2PID:      node.PeerID.Raw(),
					OCRUrl:     "", // TODO: this will be needed to distribute job specs
					DON2DONUrl: "",
				})
				artifact.Nodes[don.Name].Nodes[node.NodeID] = SimpleNodeArtifact{Name: node.Name}
				continue
			}

			artifact.Nodes[don.Name].Nodes[node.NodeID] = SimpleNodeArtifact{Name: node.Name}
			donArtifact.Nodes = append(donArtifact.Nodes, FullNodeArtifact{
				NOP:    nop.Name,
				Name:   node.Name,
				CSAKey: node.CSAKey,
			})
		}
	}

	return &artifact, nil
}

func persistArtifact(absPath string, artifact *EnvArtifact) (string, error) {
	err := os.MkdirAll(filepath.Dir(absPath), 0755)
	if err != nil {
		return "", pkgerrors.Wrap(err, "failed to create directory for the environment artifact")
	}

	err = WriteJSONFile(absPath, artifact)
	if err != nil {
		return "", pkgerrors.Wrap(err, "failed to write environment artifact to file")
	}

	return absPath, nil
}

func ReadEnvArtifact(absPath string) (*EnvArtifact, error) {
	var artifact EnvArtifact

	content, readErr := os.ReadFile(absPath)
	if readErr != nil {
		return nil, pkgerrors.Wrapf(readErr, "failed to read environment artifact from %s. Make sure that local CRE environment is running", absPath)
	}

	if err := json.Unmarshal(content, &artifact); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to unmarshal environment artifact")
	}

	return &artifact, nil
}

func MustEnvArtifactAbsPath(relativePathToRepoRoot string) string {
	path, err := filepath.Abs(filepath.Join(relativePathToRepoRoot, envconfig.StateDirname, ArtifactFileName))
	if err != nil {
		panic(err)
	}

	return path
}

func EnvArtifactFileExists(relativePathToRepoRoot string) bool {
	_, statErr := os.Stat(MustEnvArtifactAbsPath(relativePathToRepoRoot))
	return statErr == nil
}

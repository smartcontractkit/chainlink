package environment

import (
	"fmt"
	"os"
	"path/filepath"

	pkgerrors "github.com/pkg/errors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf_deployment "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"

	"github.com/smartcontractkit/chainlink/deployment"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
)

const (
	artifactDirName = "env_artifact"
	NOPAdminPrefix  = "0xaadd000000000000000000000000000000"
)

type EnvArtifact struct {
	AddressRefs   []datastore.AddressRef                               `json:"address_refs"`
	AddressBook   map[uint64]map[string]cldf_deployment.TypeAndVersion `json:"address_book"`
	JdConfig      jd.Output                                            `json:"jd_config"`
	Nodes         NodesArtifact                                        `json:"nodes"`
	DONs          []DonArtifact                                        `json:"dons"`
	Bootstrappers []BootstrapNodeArtifact                              `json:"bootstrappers"`
	NOPs          []NOPArtifact                                        `json:"nops"`
}

type NodesArtifact struct {
	Nodes map[string]SimpleNodeArtifact `json:"nodes"`
}

type SimpleNodeArtifact struct {
	Name string `json:"name"`
}

type DonArtifact struct {
	DonName        string                  `json:"don_name"`
	DonID          int                     `json:"don_id"`
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
	datastore datastore.AddressRefStore,
	addressBook cldf_deployment.AddressBook,
	jdOutput jd.Output,
	donTopology cre.DonTopology,
	offchainClient cldf_deployment.OffchainClient,
	capabilityFactoryFns []cre.DONCapabilityWithConfigFactoryFn,
) (string, error) {
	artifact, err := GenerateArtifact(datastore, addressBook, jdOutput, donTopology, offchainClient, capabilityFactoryFns)
	if err != nil {
		return "", pkgerrors.Wrap(err, "failed to generate environment artifact")
	}

	// Let's save the artifact to disk
	artifactPath, err := persistArtifact(artifact)
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
	offchainClient cldf_deployment.OffchainClient,
	capabilityFactoryFns []cre.DONCapabilityWithConfigFactoryFn,
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
		JdConfig:    jdOutput,
		AddressBook: addresses,
		AddressRefs: addressRecords,
		Nodes: NodesArtifact{
			Nodes: make(map[string]SimpleNodeArtifact),
		},
		DONs:          make([]DonArtifact, 0),
		Bootstrappers: make([]BootstrapNodeArtifact, 0),
		NOPs:          make([]NOPArtifact, 0),
	}

	for i, don := range donTopology.DonsWithMetadata {
		donArtifact := DonArtifact{
			DonName:        don.Name,
			DonID:          int(don.ID),
			F:              0, // F will be calculated based on the number of worker nodes
			BootstrapNodes: make([]string, 0),
			Nodes:          make([]FullNodeArtifact, 0),
			Capabilities:   make([]DONCapabilityArtifact, 0),
		}

		workerNodes, workerNodesErr := crenode.FindManyWithLabel(don.NodesMetadata, &cre.Label{
			Key:   crenode.NodeTypeKey,
			Value: cre.WorkerNode,
		}, crenode.EqualLabels)
		if workerNodesErr != nil {
			return nil, pkgerrors.Wrap(workerNodesErr, "failed to find worker nodes")
		}

		donArtifact.F = libc.MustSafeUint8((len(workerNodes) - 1) / 3)

		for _, factoryFn := range capabilityFactoryFns {
			capabilities := factoryFn(don.Flags)
			for _, capability := range capabilities {
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
			ID:    i + 1, // NOP IDs start from 1
			Name:  fmt.Sprintf("NOP for %s DON", don.Name),
			Admin: fmt.Sprintf("%s%06d", NOPAdminPrefix, i+1),
		}

		var nodeIDs []string
		for _, node := range don.DON.Nodes {
			nodeIDs = append(nodeIDs, node.NodeID)
		}

		nodeInfo, nodeInfoErr := deployment.NodeInfo(nodeIDs, offchainClient)
		if nodeInfoErr != nil {
			return nil, pkgerrors.Wrapf(nodeInfoErr, "failed to get node info for DON %s", don.Name)
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
				continue
			}

			artifact.Nodes.Nodes[node.NodeID] = SimpleNodeArtifact{Name: node.Name}
			donArtifact.Nodes = append(donArtifact.Nodes, FullNodeArtifact{
				NOP:    nop.Name,
				Name:   node.Name,
				CSAKey: node.CSAKey,
			})
		}

		artifact.NOPs = append(artifact.NOPs, nop)
		artifact.DONs = append(artifact.DONs, donArtifact)
	}

	return &artifact, nil
}

func persistArtifact(artifact *EnvArtifact) (string, error) {
	err := os.MkdirAll(artifactDirName, 0755)
	if err != nil {
		return "", pkgerrors.Wrap(err, "failed to create directory for the environment artifact")
	}
	err = WriteJSONFile(artifactDirName+"/env_artifact.json", artifact)
	if err != nil {
		return "", pkgerrors.Wrap(err, "failed to write environment artifact to file")
	}

	absPath, absPathErr := filepath.Abs(artifactDirName + "/env_artifact.json")
	if absPathErr != nil {
		return "", pkgerrors.Wrap(absPathErr, "failed to get absolute path for the environment artifact")
	}

	return absPath, nil
}

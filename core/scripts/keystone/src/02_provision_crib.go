package src

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	"gopkg.in/yaml.v3"

	evmcfg "github.com/smartcontractkit/chainlink-integrations/evm/config/toml"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type Helm struct {
	Helm Chart `yaml:"helm"`
}

type Chart struct {
	HelmValues HelmValues `yaml:"values"`
}

type HelmValues struct {
	Chainlink Chainlink `yaml:"chainlink,omitempty"`
	Ingress   Ingress   `yaml:"ingress,omitempty"`
}

type Ingress struct {
	Hosts []Host `yaml:"hosts,omitempty"`
}

type Host struct {
	Host string `yaml:"host,omitempty"`
	HTTP HTTP   `yaml:"http,omitempty"`
}

type HTTP struct {
	Paths []Path `yaml:"paths,omitempty"`
}

type Path struct {
	Path    string  `yaml:"path,omitempty"`
	Backend Backend `yaml:"backend,omitempty"`
}

type Backend struct {
	Service Service `yaml:"service,omitempty"`
}

type Service struct {
	Name string `yaml:"name,omitempty"`
	Port Port   `yaml:"port,omitempty"`
}

type Port struct {
	Number int `yaml:"number,omitempty"`
}

type Chainlink struct {
	Nodes map[string]Node `yaml:"nodes,omitempty"`
}

type Node struct {
	Image         string `yaml:"image,omitempty"`
	OverridesToml string `yaml:"overridesToml,omitempty"`
}

func writePreprovisionConfig(nodeSetSize int, outputPath string) {
	chart := generatePreprovisionConfig(nodeSetSize)

	writeCribConfig(chart, outputPath)
}

func writeCribConfig(chart Helm, outputPath string) {
	yamlData, err := yaml.Marshal(chart)
	helpers.PanicErr(err)

	if outputPath == "-" {
		_, err = os.Stdout.Write(yamlData)
		helpers.PanicErr(err)
	} else {
		ensureArtefactsDir(filepath.Dir(outputPath))
		err = os.WriteFile(outputPath, yamlData, 0600)
		helpers.PanicErr(err)
	}
}

func generatePreprovisionConfig(nodeSetSize int) Helm {
	nodeSets := []string{"ks-wf-", "ks-str-trig-"}
	nodes := make(map[string]Node)
	nodeNames := []string{}

	for nodeSetIndex, prefix := range nodeSets {
		// Bootstrap node
		btNodeName := fmt.Sprintf("%d-%sbt-node1", nodeSetIndex, prefix)
		nodeNames = append(nodeNames, btNodeName)
		nodes[btNodeName] = Node{
			Image: "${runtime.images.app}",
		}

		// Other nodes
		for i := 2; i <= nodeSetSize; i++ {
			nodeName := fmt.Sprintf("%d-%snode%d", nodeSetIndex, prefix, i)
			nodeNames = append(nodeNames, nodeName)
			nodes[nodeName] = Node{
				Image: "${runtime.images.app}",
			}
		}
	}

	ingress := generateIngress(nodeNames)

	helm := Helm{
		Chart{
			HelmValues: HelmValues{
				Chainlink: Chainlink{
					Nodes: nodes,
				},
				Ingress: ingress,
			},
		},
	}

	return helm
}

func writePostProvisionConfig(
	nodeSets NodeSets,
	chainID int64,
	capabilitiesP2PPort int64,
	forwarderAddress string,
	capabilitiesRegistryAddress string,
	outputPath string,
) {
	chart := generatePostprovisionConfig(
		nodeSets,
		chainID,
		capabilitiesP2PPort,
		forwarderAddress,
		capabilitiesRegistryAddress,
	)

	writeCribConfig(chart, outputPath)
}

func generatePostprovisionConfig(
	nodeSets NodeSets,
	chainID int64,
	capabilitiesP2PPort int64,
	forwarderAddress string,
	capabillitiesRegistryAddress string,
) Helm {
	nodes := make(map[string]Node)
	nodeNames := []string{}
	var capabilitiesBootstrapper *ocrcommontypes.BootstrapperLocator

	// Build nodes for each NodeSet
	for nodeSetIndex, nodeSet := range []NodeSet{nodeSets.Workflow, nodeSets.StreamsTrigger} {
		// Bootstrap node
		btNodeName := fmt.Sprintf("%d-%sbt-node1", nodeSetIndex, nodeSet.Prefix)
		// Note this line ordering is important,
		// we assign capabilitiesBootstrapper after we generate overrides so that
		// we do not include the bootstrapper config to itself
		overridesToml := generateOverridesToml(
			chainID,
			capabilitiesP2PPort,
			capabillitiesRegistryAddress,
			"",
			"",
			capabilitiesBootstrapper,
			nodeSet.Name,
		)
		nodes[btNodeName] = Node{
			Image:         "${runtime.images.app}",
			OverridesToml: overridesToml,
		}
		if nodeSet.Name == WorkflowNodeSetName {
			workflowBtNodeKey := nodeSets.Workflow.NodeKeys[0] // First node key as bootstrapper
			wfBt, err := ocrcommontypes.NewBootstrapperLocator(workflowBtNodeKey.P2PPeerID, []string{fmt.Sprintf("%s:%d", nodeSets.Workflow.Nodes[0].ServiceName, capabilitiesP2PPort)})
			helpers.PanicErr(err)
			capabilitiesBootstrapper = wfBt
		}
		nodeNames = append(nodeNames, btNodeName)

		// Other nodes
		for i, nodeKey := range nodeSet.NodeKeys[1:] { // Start from second key
			nodeName := fmt.Sprintf("%d-%snode%d", nodeSetIndex, nodeSet.Prefix, i+2)
			nodeNames = append(nodeNames, nodeName)
			overridesToml := generateOverridesToml(
				chainID,
				capabilitiesP2PPort,
				capabillitiesRegistryAddress,
				nodeKey.EthAddress,
				forwarderAddress,
				capabilitiesBootstrapper,
				nodeSet.Name,
			)
			nodes[nodeName] = Node{
				Image:         "${runtime.images.app}",
				OverridesToml: overridesToml,
			}
		}
	}

	ingress := generateIngress(nodeNames)

	helm := Helm{
		Chart{
			HelmValues: HelmValues{
				Chainlink: Chainlink{
					Nodes: nodes,
				},
				Ingress: ingress,
			},
		},
	}

	return helm
}

func generateOverridesToml(
	chainID int64,
	capabilitiesP2PPort int64,
	externalRegistryAddress string,
	fromAddress string,
	forwarderAddress string,
	capabilitiesBootstrapper *ocrcommontypes.BootstrapperLocator,
	nodeSetName string,
) string {
	evmConfig := &evmcfg.EVMConfig{
		ChainID: big.NewI(chainID),
		Nodes:   nil, // We have the rpc nodes set globally
	}

	conf := chainlink.Config{
		Core: toml.Core{
			Capabilities: toml.Capabilities{
				ExternalRegistry: toml.ExternalRegistry{
					Address:   ptr(externalRegistryAddress),
					NetworkID: ptr("evm"),
					ChainID:   ptr(strconv.FormatInt(chainID, 10)),
				},
				Peering: toml.P2P{
					V2: toml.P2PV2{
						Enabled:         ptr(true),
						ListenAddresses: ptr([]string{fmt.Sprintf("0.0.0.0:%d", capabilitiesP2PPort)}),
					},
				},
			},
		},
	}

	if capabilitiesBootstrapper != nil {
		conf.Core.Capabilities.Peering.V2.DefaultBootstrappers = ptr([]ocrcommontypes.BootstrapperLocator{*capabilitiesBootstrapper})

		if nodeSetName == WorkflowNodeSetName {
			evmConfig.Workflow = evmcfg.Workflow{
				FromAddress:      ptr(evmtypes.MustEIP55Address(fromAddress)),
				ForwarderAddress: ptr(evmtypes.MustEIP55Address(forwarderAddress)),
			}
		}
	}

	conf.EVM = evmcfg.EVMConfigs{
		evmConfig,
	}

	confStr, err := conf.TOMLString()
	helpers.PanicErr(err)

	return confStr
}

// New function to generate Ingress
func generateIngress(nodeNames []string) Ingress {
	hosts := make([]Host, 0, len(nodeNames))

	for _, nodeName := range nodeNames {
		host := Host{
			Host: fmt.Sprintf("${DEVSPACE_NAMESPACE}-%s.${DEVSPACE_INGRESS_BASE_DOMAIN}", nodeName),
			HTTP: HTTP{
				Paths: []Path{
					{
						Path: "/",
						Backend: Backend{
							Service: Service{
								Name: "app-" + nodeName,
								Port: Port{
									Number: 6688,
								},
							},
						},
					},
				},
			},
		}
		hosts = append(hosts, host)
	}

	return Ingress{
		Hosts: hosts,
	}
}

func ptr[T any](t T) *T { return &t }

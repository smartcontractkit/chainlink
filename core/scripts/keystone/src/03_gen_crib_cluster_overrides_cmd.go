package src

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	"gopkg.in/yaml.v3"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type generateCribClusterOverridesPreprovision struct{}
type generateCribClusterOverridesPostprovision struct{}

func NewGenerateCribClusterOverridesPreprovisionCommand() *generateCribClusterOverridesPreprovision {
	return &generateCribClusterOverridesPreprovision{}
}

func NewGenerateCribClusterOverridesPostprovisionCommand() *generateCribClusterOverridesPostprovision {
	return &generateCribClusterOverridesPostprovision{}
}

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

func (g *generateCribClusterOverridesPreprovision) Name() string {
	return "crib-config-preprovision"
}

func (g *generateCribClusterOverridesPreprovision) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	nodeSetSize := fs.Int("nodeSetSize", 4, "number of nodes in a nodeset")
	outputPath := fs.String("outpath", "-", "the path to output the generated overrides (use '-' for stdout)")

	err := fs.Parse(args)
	if err != nil || outputPath == nil || *outputPath == "" {
		fs.Usage()
		os.Exit(1)
	}

	chart := generatePreprovisionConfig(*nodeSetSize)

	yamlData, err := yaml.Marshal(chart)
	helpers.PanicErr(err)

	if *outputPath == "-" {
		_, err = os.Stdout.Write(yamlData)
		helpers.PanicErr(err)
	} else {
		err = os.WriteFile(filepath.Join(*outputPath, "crib-preprovision.yaml"), yamlData, 0600)
		helpers.PanicErr(err)
	}
}

func generatePreprovisionConfig(nodeSetSize int) Helm {
	nodeSets := []string{"ks-wf-", "ks-str-trig-"}
	nodes := make(map[string]Node)
	nodeNames := []string{}
	globalIndex := 0

	for _, prefix := range nodeSets {
		// Bootstrap node
		btNodeName := fmt.Sprintf("%d-%sbt-node1", globalIndex, prefix)
		nodeNames = append(nodeNames, btNodeName)
		nodes[btNodeName] = Node{
			Image: "${runtime.images.app}",
		}

		// Other nodes
		for i := 2; i <= nodeSetSize; i++ {
			globalIndex++
			nodeName := fmt.Sprintf("%d-%snode%d", globalIndex, prefix, i)
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

func (g *generateCribClusterOverridesPostprovision) Name() string {
	return "crib-config-postprovision"
}

func (g *generateCribClusterOverridesPostprovision) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	chainID := fs.Int64("chainid", 1337, "chain id")
	outputPath := fs.String("outpath", "-", "the path to output the generated overrides (use '-' for stdout)")
	nodeSetSize := fs.Int("nodeSetSize", 4, "number of nodes in a nodeset")
	publicKeys := fs.String("publickeys", "", "Custom public keys json location")
	nodeList := fs.String("nodes", "", "Custom node list location")
	artefactsDir := fs.String("artefacts", "", "Custom artefacts directory location")

	err := fs.Parse(args)
	if err != nil || outputPath == nil || *outputPath == "" || chainID == nil || *chainID == 0 {
		fs.Usage()
		os.Exit(1)
	}

	if *artefactsDir == "" {
		*artefactsDir = defaultArtefactsDir
	}
	if *publicKeys == "" {
		*publicKeys = defaultPublicKeys
	}
	if *nodeList == "" {
		*nodeList = defaultNodeList
	}

	contracts, err := LoadDeployedContracts(*artefactsDir)
	helpers.PanicErr(err)

	chart := generatePostprovisionConfig(nodeList, chainID, publicKeys, contracts, *nodeSetSize)

	yamlData, err := yaml.Marshal(chart)
	helpers.PanicErr(err)

	if *outputPath == "-" {
		_, err = os.Stdout.Write(yamlData)
		helpers.PanicErr(err)
	} else {
		err = os.WriteFile(filepath.Join(*outputPath, "crib-postprovision.yaml"), yamlData, 0600)
		helpers.PanicErr(err)
	}
}

func generatePostprovisionConfig(nodeList *string, chainID *int64, publicKeys *string, contracts deployedContracts, nodeSetSize int) Helm {
	nodeSets := downloadNodeSets(*nodeList, *chainID, *publicKeys, nodeSetSize)

	nodes := make(map[string]Node)
	nodeNames := []string{}
	// FIXME: Ideally we just save the node list as a map and don't need to sort
	globalIndex := 0

	// Prepare the bootstrapper locator from the workflow NodeSet
	workflowBtNodeName := fmt.Sprintf("%d-%sbt-node1", globalIndex, nodeSets.Workflow.Prefix)
	nodeNames = append(nodeNames, workflowBtNodeName)
	workflowBtNodeKey := nodeSets.Workflow.NodeKeys[0] // First node key as bootstrapper
	wfBt, err := ocrcommontypes.NewBootstrapperLocator(workflowBtNodeKey.P2PPeerID, []string{fmt.Sprintf("%s:6691", workflowBtNodeName)})
	helpers.PanicErr(err)

	// Build nodes for each NodeSet
	for _, nodeSet := range []NodeSet{nodeSets.Workflow, nodeSets.StreamsTrigger} {
		// Bootstrap node
		btNodeName := fmt.Sprintf("%d-%sbt-node1", globalIndex, nodeSet.Prefix)
		nodeNames = append(nodeNames, btNodeName)
		overridesToml := generateOverridesToml(
			*chainID,
			contracts.CapabilityRegistry.Hex(),
			"",
			"",
			nil,
			nodeSet.Name,
		)
		nodes[btNodeName] = Node{
			Image:         "${runtime.images.app}",
			OverridesToml: overridesToml,
		}

		// Other nodes
		for i, nodeKey := range nodeSet.NodeKeys[1:] { // Start from second key
			globalIndex++
			nodeName := fmt.Sprintf("%d-%snode%d", globalIndex, nodeSet.Prefix, i+2)
			nodeNames = append(nodeNames, nodeName)
			overridesToml := generateOverridesToml(
				*chainID,
				contracts.CapabilityRegistry.Hex(),
				nodeKey.EthAddress,
				contracts.ForwarderContract.Hex(),
				wfBt,
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
					ChainID:   ptr(fmt.Sprintf("%d", chainID)),
				},
				Peering: toml.P2P{
					V2: toml.P2PV2{
						Enabled:         ptr(true),
						ListenAddresses: ptr([]string{"0.0.0.0:6691"}),
					},
				},
			},
		},
	}

	if capabilitiesBootstrapper != nil {
		conf.Core.Capabilities.Peering.V2.DefaultBootstrappers = ptr([]ocrcommontypes.BootstrapperLocator{*capabilitiesBootstrapper})

		// FIXME: Use const for names
		if nodeSetName == "workflow" {
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
	var hosts []Host

	for _, nodeName := range nodeNames {
		host := Host{
			Host: fmt.Sprintf("${DEVSPACE_NAMESPACE}-%s.${DEVSPACE_INGRESS_BASE_DOMAIN}", nodeName),
			HTTP: HTTP{
				Paths: []Path{
					{
						Path: "/",
						Backend: Backend{
							Service: Service{
								Name: fmt.Sprintf("app-%s", nodeName),
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

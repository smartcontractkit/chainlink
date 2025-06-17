package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"

	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// NodeSet represents a set of nodes with associated metadata.
// NodeKeys are indexed by the same order as Nodes.
type NodeSet struct {
	Name     string
	Prefix   string
	Nodes    []NodeWithCreds
	NodeKeys []NodeKeys
}

var (
	WorkflowNodeSetName         = "workflow"
	WorkflowNodeSetPrefix       = "ks-wf-"
	StreamsTriggerNodeSetName   = "streams-trigger"
	StreamsTriggerNodeSetPrefix = "ks-str-trig-"
)

// NodeSets holds the two NodeSets: Workflow and StreamsTrigger.
type NodeSets struct {
	Workflow       NodeSet
	StreamsTrigger NodeSet
}

func downloadNodeSets(chainID int64, nodeSetPath string, nodeSetSize int) NodeSets {
	if _, err := os.Stat(nodeSetPath); err == nil {
		fmt.Println("Loading existing nodesets at:", nodeSetPath)
		nodeSets := mustReadJSON[NodeSets](nodeSetPath)
		return nodeSets
	}

	fmt.Println("Connecting to Kubernetes to fetch node credentials...")
	crib := NewCribClient()
	nodes, err := crib.getCLNodes()
	PanicErr(err)

	totalNodes := len(nodes)
	// Workflow and StreamsTrigger nodeSets should have the same number of nodes
	// hence we need at least 2 * nodeSetSize nodes
	requiredNodes := nodeSetSize * 2
	if totalNodes < requiredNodes {
		panic(fmt.Errorf("not enough nodes to populate both nodeSets: required %d, got %d", requiredNodes, totalNodes))
	}

	nodeSets := NodeSets{
		Workflow: NodeSet{
			Name:   WorkflowNodeSetName,
			Prefix: WorkflowNodeSetPrefix,
			Nodes:  nodes[:nodeSetSize],
		},
		StreamsTrigger: NodeSet{
			Name:   StreamsTriggerNodeSetName,
			Prefix: StreamsTriggerNodeSetPrefix,
			Nodes:  nodes[nodeSetSize : nodeSetSize*2],
		},
	}

	nodeSets.Workflow.NodeKeys = mustFetchNodeKeys(chainID, nodeSets.Workflow.Nodes, true)
	nodeSets.StreamsTrigger.NodeKeys = mustFetchNodeKeys(chainID, nodeSets.StreamsTrigger.Nodes, false)
	mustWriteJSON(nodeSetPath, nodeSets)

	return nodeSets
}

// NodeKeys represents the keys for a single node.
// If there are multiple OCR2KBs or OCR2AptosKBs, only the first one is used.
type NodeKeys struct {
	AptosAccount string `json:"AptosAccount"`
	EthAddress   string `json:"EthAddress"`
	P2PPeerID    string `json:"P2PPeerID"`
	CSAPublicKey string `json:"CSAPublicKey"`
	OCR2KBTrimmed
	OCR2AptosKBTrimmed
}

// This is an OCR key bundle with the prefixes on each respective key
// trimmed off
type OCR2KBTrimmed struct {
	OCR2BundleID          string `json:"OCR2BundleID"`          // used only in job spec
	OCR2OnchainPublicKey  string `json:"OCR2OnchainPublicKey"`  // ocr2on_evm_<key>
	OCR2OffchainPublicKey string `json:"OCR2OffchainPublicKey"` // ocr2off_evm_<key>
	OCR2ConfigPublicKey   string `json:"OCR2ConfigPublicKey"`   // ocr2cfg_evm_<key>
}

// This is an Aptos key bundle with the prefixes on each respective key
// trimmed off
type OCR2AptosKBTrimmed struct {
	AptosBundleID         string `json:"AptosBundleID"`
	AptosOnchainPublicKey string `json:"AptosOnchainPublicKey"` // ocr2on_aptos_<key>
}

func mustFetchNodeKeys(chainID int64, nodes []NodeWithCreds, createAptosKeys bool) []NodeKeys {
	nodeKeys := []NodeKeys{}

	for _, n := range nodes {
		api := newNodeAPI(n)
		// Get eth key
		fmt.Printf("Fetching ETH keys for node %s\n", n.ServiceName)
		eKey := api.mustExec(api.methods.ListETHKeys)
		ethKeys := mustJSON[[]presenters.ETHKeyResource](eKey)
		ethAddress, err := findFirstGoodEthKeyAddress(chainID, *ethKeys)
		helpers.PanicErr(err)

		var aptosAccount string
		if createAptosKeys {
			aptosAccount = getOrCreateAptosKey(api)
		}

		// Get p2p key
		fmt.Printf("Fetching P2P key for node %s\n", n.ServiceName)
		p2pKeys := api.mustExec(api.methods.ListP2PKeys)
		p2pKey := mustJSON[[]presenters.P2PKeyResource](p2pKeys)
		if len(*p2pKey) != 1 {
			helpers.PanicErr(errors.New("node must have single p2p key"))
		}
		peerID := strings.TrimPrefix((*p2pKey)[0].PeerID, "p2p_")

		// Get OCR2 key bundles for both EVM and Aptos chains
		bundles := api.mustExec(api.methods.ListOCR2KeyBundles)
		ocr2Bundles := mustJSON[cmd.OCR2KeyBundlePresenters](bundles)

		expectedBundleLen := 1

		// evm key bundles
		fmt.Printf("Fetching OCR2 EVM key bundles for node %s\n", n.ServiceName)
		ocr2EvmBundles := getTrimmedEVMOCR2KBs(*ocr2Bundles)
		evmBundleLen := len(ocr2EvmBundles)
		if evmBundleLen < expectedBundleLen {
			fmt.Printf("WARN: node has %d EVM OCR2 bundles when it should have at least %d, creating bundles...\n", evmBundleLen, expectedBundleLen)
			for i := evmBundleLen; i < expectedBundleLen; i++ {
				cBundle := api.withArg("evm").mustExec(api.methods.CreateOCR2KeyBundle)
				createdBundle := mustJSON[cmd.OCR2KeyBundlePresenter](cBundle)
				fmt.Printf("Created OCR2 EVM key bundle %s\n", string(cBundle))
				ocr2EvmBundles = append(ocr2EvmBundles, trimmedOCR2KB(*createdBundle))
			}
		}

		// aptos key bundles
		var ocr2AptosBundles []OCR2AptosKBTrimmed
		if createAptosKeys {
			fmt.Printf("Fetching OCR2 Aptos key bundles for node %s\n", n.ServiceName)
			ocr2AptosBundles = createAptosOCR2KB(ocr2Bundles, expectedBundleLen, api)
		}

		fmt.Printf("Fetching CSA keys for node %s\n", n.ServiceName)
		csaKeys := api.mustExec(api.methods.ListCSAKeys)
		csaKeyResources := mustJSON[[]presenters.CSAKeyResource](csaKeys)
		csaPubKey, err := findFirstCSAPublicKey(*csaKeyResources)
		helpers.PanicErr(err)

		// We can handle multiple OCR bundles in the future
		// but for now we only support a single bundle per node
		keys := NodeKeys{
			OCR2KBTrimmed: ocr2EvmBundles[0],
			EthAddress:    ethAddress,
			AptosAccount:  aptosAccount,
			P2PPeerID:     peerID,
			CSAPublicKey:  strings.TrimPrefix(csaPubKey, "csa_"),
		}
		if createAptosKeys {
			keys.OCR2AptosKBTrimmed = ocr2AptosBundles[0]
		}

		nodeKeys = append(nodeKeys, keys)
	}

	return nodeKeys
}

func trimmedOCR2KB(ocr2Bndl cmd.OCR2KeyBundlePresenter) OCR2KBTrimmed {
	return OCR2KBTrimmed{
		OCR2BundleID:          ocr2Bndl.ID,
		OCR2ConfigPublicKey:   strings.TrimPrefix(ocr2Bndl.ConfigPublicKey, "ocr2cfg_evm_"),
		OCR2OnchainPublicKey:  strings.TrimPrefix(ocr2Bndl.OnchainPublicKey, "ocr2on_evm_"),
		OCR2OffchainPublicKey: strings.TrimPrefix(ocr2Bndl.OffChainPublicKey, "ocr2off_evm_"),
	}
}

func trimmedAptosOCR2KB(ocr2Bndl cmd.OCR2KeyBundlePresenter) OCR2AptosKBTrimmed {
	return OCR2AptosKBTrimmed{
		AptosBundleID:         ocr2Bndl.ID,
		AptosOnchainPublicKey: strings.TrimPrefix(ocr2Bndl.OnchainPublicKey, "ocr2on_aptos_"),
	}
}

func createAptosOCR2KB(ocr2Bundles *cmd.OCR2KeyBundlePresenters, expectedBundleLen int, api *nodeAPI) []OCR2AptosKBTrimmed {
	ocr2AptosBundles := getTrimmedAptosOCR2KBs(*ocr2Bundles)
	aptosBundleLen := len(ocr2AptosBundles)

	if aptosBundleLen < expectedBundleLen {
		fmt.Printf("WARN: node has %d Aptos OCR2 bundles when it should have at least %d, creating bundles...\n", aptosBundleLen, expectedBundleLen)
		for i := aptosBundleLen; i < expectedBundleLen; i++ {
			cBundle := api.withArg("aptos").mustExec(api.methods.CreateOCR2KeyBundle)
			createdBundle := mustJSON[cmd.OCR2KeyBundlePresenter](cBundle)
			fmt.Println("Created OCR2 Aptos key bundle", string(cBundle))
			ocr2AptosBundles = append(ocr2AptosBundles, trimmedAptosOCR2KB(*createdBundle))
		}
	}

	return ocr2AptosBundles
}

// getOrCreateAptosKey returns the Aptos account of the node.
//
// If the node has no Aptos keys, it creates one and returns the account.
func getOrCreateAptosKey(api *nodeAPI) string {
	api.output.Reset()
	aKeysClient := cmd.NewAptosKeysClient(api.methods)
	err := aKeysClient.ListKeys(&cli.Context{App: api.app})
	helpers.PanicErr(err)
	var aptosKeys []presenters.AptosKeyResource
	helpers.PanicErr(json.Unmarshal(api.output.Bytes(), &aptosKeys))
	if len(aptosKeys) == 0 {
		api.output.Reset()
		fmt.Printf("WARN: node has no aptos keys, creating one...\n")
		err = aKeysClient.CreateKey(&cli.Context{App: api.app})
		helpers.PanicErr(err)
		api.output.Reset()
		err = aKeysClient.ListKeys(&cli.Context{App: api.app})
		helpers.PanicErr(err)
		helpers.PanicErr(json.Unmarshal(api.output.Bytes(), &aptosKeys))
		api.output.Reset()
	}

	if len(aptosKeys) != 1 {
		fmt.Printf("Node has %d aptos keys\n", len(aptosKeys))
		PanicErr(errors.New("node must have single aptos key"))
	}

	aptosAccount := aptosKeys[0].Account
	api.output.Reset()

	return aptosAccount
}

func getTrimmedAptosOCR2KBs(ocr2Bundles cmd.OCR2KeyBundlePresenters) []OCR2AptosKBTrimmed {
	aptosBundles := []OCR2AptosKBTrimmed{}
	for _, b := range ocr2Bundles {
		if b.ChainType == "aptos" {
			aptosBundles = append(aptosBundles, trimmedAptosOCR2KB(b))
		}
	}
	return aptosBundles
}

func getTrimmedEVMOCR2KBs(ocr2Bundles cmd.OCR2KeyBundlePresenters) []OCR2KBTrimmed {
	evmBundles := []OCR2KBTrimmed{}
	for _, b := range ocr2Bundles {
		if b.ChainType == "evm" {
			evmBundles = append(evmBundles, trimmedOCR2KB(b))
		}
	}
	return evmBundles
}

func findFirstCSAPublicKey(csaKeyResources []presenters.CSAKeyResource) (string, error) {
	for _, r := range csaKeyResources {
		return r.PubKey, nil
	}
	return "", errors.New("did not find any CSA Key Resources")
}

func findFirstGoodEthKeyAddress(chainID int64, ethKeys []presenters.ETHKeyResource) (string, error) {
	for _, ethKey := range ethKeys {
		if ethKey.EVMChainID.Equal(ubig.NewI(chainID)) && !ethKey.Disabled {
			return ethKey.Address, nil
		}
	}
	return "", errors.New("did not find an enabled ETH key for the given chain ID")
}

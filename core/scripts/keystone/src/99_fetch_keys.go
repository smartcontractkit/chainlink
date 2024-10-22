package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// KeylessNodeSet represents a set of nodes without NodeKeys.
type KeylessNodeSet struct {
	Name   string
	Prefix string
	Nodes  []*NodeWthCreds
}

// NodeSet represents a set of nodes with associated metadata.
// It embeds KeylessNodeSet and includes NodeKeys.
// NodeKeys are indexed by the same order as Nodes.
type NodeSet struct {
	KeylessNodeSet
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

// downloadKeylessNodeSets downloads the node API credentials or loads them from disk if they already exist.
// It returns a NodeSets struct without NodeKeys.
func downloadKeylessNodeSets(keylessNodeSetPath string, nodeSetSize int) NodeSets {
	if _, err := os.Stat(keylessNodeSetPath); err == nil {
		fmt.Println("Loading existing keyless nodesets at:", keylessNodeSetPath)
		nodeSets := mustReadJSON[NodeSets](keylessNodeSetPath)

		return nodeSets
	}

	fmt.Println("Connecting to Kubernetes to fetch node credentials...")
	crib := NewCribClient()
	clNodesWithCreds, err := crib.GetCLNodeCredentials()
	PanicErr(err)

	nodesList := clNodesWithCredsToNodes(clNodesWithCreds)
	if len(nodesList) == 0 {
		panic("no nodes found")
	}
	keylessNodeSets, err := splitNodesIntoNodeSets(nodesList, nodeSetSize)
	PanicErr(err)

	mustWriteJSON(keylessNodeSetPath, keylessNodeSets)

	return keylessNodeSets
}

func downloadNodeSets(keylessNodeSetsPath string, chainID int64, nodeSetPath string, nodeSetSize int) NodeSets {
	if _, err := os.Stat(nodeSetPath); err == nil {
		fmt.Println("Loading existing nodesets at", nodeSetPath)
		nodeSets := mustReadJSON[NodeSets](nodeSetPath)

		return nodeSets
	}

	nodeSetsWithoutKeys := downloadKeylessNodeSets(keylessNodeSetsPath, nodeSetSize)
	nodeSets := populateNodeKeys(chainID, nodeSetsWithoutKeys)
	mustWriteJSON(nodeSetPath, nodeSets)

	return nodeSets
}

// splitNodesIntoNodeSets splits the nodes into NodeSets for 'workflow' and 'streams-trigger' nodeSets.
func splitNodesIntoNodeSets(nodes []*NodeWthCreds, nodeSetSize int) (NodeSets, error) {
	totalNodes := len(nodes)
	requiredNodes := nodeSetSize * 2
	if totalNodes < requiredNodes {
		return NodeSets{}, fmt.Errorf("not enough nodes to populate both nodeSets: required %d, got %d", requiredNodes, totalNodes)
	}

	return NodeSets{
		Workflow: NodeSet{
			KeylessNodeSet: KeylessNodeSet{
				Name:   WorkflowNodeSetName,
				Prefix: WorkflowNodeSetPrefix,
				Nodes:  nodes[:nodeSetSize],
			},
			// NodeKeys will be populated later
		},
		StreamsTrigger: NodeSet{
			KeylessNodeSet: KeylessNodeSet{
				Name:   StreamsTriggerNodeSetName,
				Prefix: StreamsTriggerNodeSetPrefix,
				Nodes:  nodes[nodeSetSize : nodeSetSize*2],
			},
			// NodeKeys will be populated later
		},
	}, nil
}

// populateNodeKeys fetches and assigns NodeKeys to each NodeSet in NodeSets.
func populateNodeKeys(chainID int64, nodeSetsWithoutKeys NodeSets) NodeSets {
	var nodeSets NodeSets

	nodeSets.Workflow = NodeSet{
		KeylessNodeSet: nodeSetsWithoutKeys.Workflow.KeylessNodeSet,
	}
	workflowKeys := mustFetchAllNodeKeys(chainID, nodeSets.Workflow.Nodes)
	nodeSets.Workflow.NodeKeys = convertAllKeysToNodeKeys(workflowKeys)

	nodeSets.StreamsTrigger = NodeSet{
		KeylessNodeSet: nodeSetsWithoutKeys.StreamsTrigger.KeylessNodeSet,
	}
	streamsTriggerKeys := mustFetchAllNodeKeys(chainID, nodeSets.StreamsTrigger.Nodes)
	nodeSets.StreamsTrigger.NodeKeys = convertAllKeysToNodeKeys(streamsTriggerKeys)

	return nodeSets
}

// gatherAllNodeKeys aggregates all NodeKeys from NodeSets.
func gatherAllNodeKeys(nodeSets NodeSets) []AllNodeKeys {
	var allKeys []AllNodeKeys
	for _, nodeSet := range []NodeSet{nodeSets.Workflow, nodeSets.StreamsTrigger} {
		for _, key := range nodeSet.NodeKeys {
			allKeys = append(allKeys, key.toAllNodeKeys())
		}
	}
	return allKeys
}

// convertAllKeysToNodeKeys converts AllNodeKeys to NodeKeys.
func convertAllKeysToNodeKeys(allKeys []AllNodeKeys) []NodeKeys {
	nodeKeys := []NodeKeys{}
	for _, k := range allKeys {
		nodeKeys = append(nodeKeys, k.toNodeKeys())
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

type AllNodeKeys struct {
	AptosAccount string               `json:"AptosAccount"`
	OCR2AptosKBs []OCR2AptosKBTrimmed `json:"OCR2AptosKBs"`
	EthAddress   string               `json:"EthAddress"`
	P2PPeerID    string               `json:"P2PPeerID"` // p2p_<key>
	OCR2KBs      []OCR2KBTrimmed      `json:"OCR2KBs"`
	CSAPublicKey string               `json:"CSAPublicKey"`
}

func (a AllNodeKeys) toNodeKeys() NodeKeys {
	return NodeKeys{
		AptosAccount: a.AptosAccount,
		OCR2AptosKBTrimmed: OCR2AptosKBTrimmed{
			AptosBundleID:         a.OCR2AptosKBs[0].AptosBundleID,
			AptosOnchainPublicKey: a.OCR2AptosKBs[0].AptosOnchainPublicKey,
		},
		OCR2KBTrimmed: OCR2KBTrimmed{
			OCR2BundleID:          a.OCR2KBs[0].OCR2BundleID,
			OCR2ConfigPublicKey:   a.OCR2KBs[0].OCR2ConfigPublicKey,
			OCR2OnchainPublicKey:  a.OCR2KBs[0].OCR2OnchainPublicKey,
			OCR2OffchainPublicKey: a.OCR2KBs[0].OCR2OffchainPublicKey,
		},
		EthAddress:   a.EthAddress,
		P2PPeerID:    a.P2PPeerID,
		CSAPublicKey: a.CSAPublicKey,
	}
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

type NodeKeys struct {
	EthAddress string `json:"EthAddress"`
	OCR2KBTrimmed
	AptosAccount string `json:"AptosAccount"`
	OCR2AptosKBTrimmed
	P2PPeerID    string `json:"P2PPeerID"`
	CSAPublicKey string `json:"CSAPublicKey"`
}

func (n NodeKeys) toAllNodeKeys() AllNodeKeys {
	return AllNodeKeys{
		EthAddress:   n.EthAddress,
		AptosAccount: n.AptosAccount,
		P2PPeerID:    n.P2PPeerID,
		OCR2KBs: []OCR2KBTrimmed{
			{
				OCR2BundleID:          n.OCR2BundleID,
				OCR2ConfigPublicKey:   n.OCR2ConfigPublicKey,
				OCR2OnchainPublicKey:  n.OCR2OnchainPublicKey,
				OCR2OffchainPublicKey: n.OCR2OffchainPublicKey,
			},
		},
		OCR2AptosKBs: []OCR2AptosKBTrimmed{
			{
				AptosBundleID:         n.AptosBundleID,
				AptosOnchainPublicKey: n.AptosOnchainPublicKey,
			},
		},
		CSAPublicKey: n.CSAPublicKey,
	}
}

func mustFetchAllNodeKeys(chainId int64, nodes []*NodeWthCreds) []AllNodeKeys {
	allNodeKeys := []AllNodeKeys{}

	for _, n := range nodes {
		api := newNodeAPI(n)
		// Get eth key
		eKey := api.mustExec(api.methods.ListETHKeys)
		ethKeys := mustJSON[[]presenters.ETHKeyResource](eKey)
		ethAddress, err := findFirstGoodEthKeyAddress(chainId, *ethKeys)
		helpers.PanicErr(err)

		// Get aptos account key
		api.output.Reset()
		aKeysClient := cmd.NewAptosKeysClient(api.methods)
		err = aKeysClient.ListKeys(&cli.Context{App: api.app})
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
			// list number of keys
			fmt.Printf("Node has %d aptos keys\n", len(aptosKeys))
			PanicErr(errors.New("node must have single aptos key"))
		}
		fmt.Printf("Node has aptos account %s\n", aptosKeys[0].Account)
		aptosAccount := aptosKeys[0].Account
		api.output.Reset()

		// Get p2p key
		p2pKeys := api.mustExec(api.methods.ListP2PKeys)
		p2pKey := mustJSON[[]presenters.P2PKeyResource](p2pKeys)
		if len(*p2pKey) != 1 {
			helpers.PanicErr(errors.New("node must have single p2p key"))
		}
		peerID := strings.TrimPrefix((*p2pKey)[0].PeerID, "p2p_")

		// Get OCR2 key bundles for both EVM and Aptos chains
		bundles := api.mustExec(api.methods.ListOCR2KeyBundles)
		ocr2Bundles := mustJSON[cmd.OCR2KeyBundlePresenters](bundles)

		// We use the same bundle length for each chain since
		// we marhshall them together into a single multichain key
		// via ocrcommon.MarshalMultichainPublicKey
		expectedBundleLen := 2

		// evm key bundles
		ocr2EvmBundles := getTrimmedEVMOCR2KBs(*ocr2Bundles)
		evmBundleLen := len(ocr2EvmBundles)
		if evmBundleLen < expectedBundleLen {
			fmt.Printf("WARN: node has %d EVM OCR2 bundles when it should have at least 2, creating bundles...\n", evmBundleLen)
			for i := evmBundleLen; i < expectedBundleLen; i++ {
				cBundle := api.withArg("evm").mustExec(api.methods.CreateOCR2KeyBundle)
				createdBundle := mustJSON[cmd.OCR2KeyBundlePresenter](cBundle)
				fmt.Printf("Created OCR2 EVM key bundle %s\n", string(cBundle))
				ocr2EvmBundles = append(ocr2EvmBundles, trimmedOCR2KB(*createdBundle))
			}
		}

		// aptos key bundles
		ocr2AptosBundles := getTrimmedAptosOCR2KBs(*ocr2Bundles)
		aptosBundleLen := len(ocr2AptosBundles)
		if aptosBundleLen < expectedBundleLen {
			fmt.Printf("WARN: node has %d Aptos OCR2 bundles when it should have at least 2, creating bundles...\n", aptosBundleLen)
			for i := aptosBundleLen; i < expectedBundleLen; i++ {
				cBundle := api.withArg("aptos").mustExec(api.methods.CreateOCR2KeyBundle)
				createdBundle := mustJSON[cmd.OCR2KeyBundlePresenter](cBundle)
				fmt.Println("Created OCR2 Aptos key bundle", string(cBundle))
				ocr2AptosBundles = append(ocr2AptosBundles, trimmedAptosOCR2KB(*createdBundle))
			}
		}

		csaKeys := api.mustExec(api.methods.ListCSAKeys)
		csaKeyResources := mustJSON[[]presenters.CSAKeyResource](csaKeys)
		csaPubKey, err := findFirstCSAPublicKey(*csaKeyResources)
		helpers.PanicErr(err)

		nodeKeys := AllNodeKeys{
			EthAddress:   ethAddress,
			AptosAccount: aptosAccount,
			P2PPeerID:    peerID,
			OCR2KBs:      ocr2EvmBundles,
			OCR2AptosKBs: ocr2AptosBundles,
			CSAPublicKey: strings.TrimPrefix(csaPubKey, "csa_"),
		}

		allNodeKeys = append(allNodeKeys, nodeKeys)
	}

	return allNodeKeys
}

func findFirstCSAPublicKey(csaKeyResources []presenters.CSAKeyResource) (string, error) {
	for _, r := range csaKeyResources {
		return r.PubKey, nil
	}
	return "", errors.New("did not find any CSA Key Resources")
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

func findFirstGoodEthKeyAddress(chainID int64, ethKeys []presenters.ETHKeyResource) (string, error) {
	for _, ethKey := range ethKeys {
		if ethKey.EVMChainID.Equal(ubig.NewI(chainID)) && !ethKey.Disabled {
			if ethKey.EthBalance == nil || ethKey.EthBalance.IsZero() {
				fmt.Println("WARN: selected ETH address has zero balance", ethKey.Address)
			}
			return ethKey.Address, nil
		}
	}
	return "", errors.New("did not find an enabled ETH key for the given chain ID")
}

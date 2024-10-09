package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/urfave/cli"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func downloadAllNodeKeys(nodeList string, chainID int64, pubKeysPath string) []AllNodeKeys {
	if _, err := os.Stat(pubKeysPath); err == nil {
		fmt.Println("Loading existing public keys at:", pubKeysPath)
		allKeys := mustParseJSON[[]AllNodeKeys](pubKeysPath)
		return allKeys
	}

	nodes := downloadNodeAPICredentials(nodeList)
	allKeys := mustFetchAllNodeKeys(chainID, nodes)

	marshalledNodeKeys, err := json.MarshalIndent(allKeys, "", " ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(pubKeysPath, marshalledNodeKeys, 0600)
	if err != nil {
		panic(err)
	}
	fmt.Println("Keystone OCR2 public keys have been saved to:", pubKeysPath)

	return allKeys
}

func downloadNodePubKeys(nodeList string, chainID int64, pubKeysPath string, index ...int) []NodeKeys {
	keys := []NodeKeys{}
	allKeys := downloadAllNodeKeys(nodeList, chainID, pubKeysPath)

	for _, k := range allKeys {
		keys = append(keys, k.toNodeKeys(index...))
	}

	return keys
}

// downloadNodeAPICredentials downloads the node API credentials, or loads them from disk if they already exist
//
// The nodes are sorted by URL. In the case of crib, the bootstrap node is the first node in the list.
func downloadNodeAPICredentials(nodeListPath string) []*node {
	if _, err := os.Stat(nodeListPath); err == nil {
		fmt.Println("Loading existing node host list at:", nodeListPath)
		nodesList := mustReadNodesList(nodeListPath)
		return nodesList
	}

	fmt.Println("Connecting to Kubernetes to fetch node credentials...")
	crib := NewCribClient()
	clNodesWithCreds, err := crib.GetCLNodeCredentials()

	if err != nil {
		panic(err)
	}

	nodesList := clNodesWithCredsToNodes(clNodesWithCreds)
	err = writeNodesList(nodeListPath, nodesList)
	if err != nil {
		panic(err)
	}
	if len(nodesList) == 0 {
		panic("No nodes found")
	}
	return nodesList
}

func clNodesWithCredsToNodes(clNodesWithCreds []CLNodeCredentials) []*node {
	nodes := []*node{}
	for _, cl := range clNodesWithCreds {
		n := node{
			url:      cl.URL,
			password: cl.Password,
			login:    cl.Username,
		}
		nodes = append(nodes, &n)
	}

	// sort nodes by URL
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].url.String() < nodes[j].url.String()
	})
	return nodes
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

func (a AllNodeKeys) toNodeKeys(index ...int) NodeKeys {
	i := 0
	if len(index) > 0 {
		i = index[0]
	}
	if i >= len(a.OCR2KBs) {
		panic("index out of range")
	}

	return NodeKeys{
		OCR2KBTrimmed: OCR2KBTrimmed{
			OCR2BundleID:          a.OCR2KBs[i].OCR2BundleID,
			OCR2ConfigPublicKey:   a.OCR2KBs[i].OCR2ConfigPublicKey,
			OCR2OnchainPublicKey:  a.OCR2KBs[i].OCR2OnchainPublicKey,
			OCR2OffchainPublicKey: a.OCR2KBs[i].OCR2OffchainPublicKey,
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

func mustFetchAllNodeKeys(chainId int64, nodes []*node) []AllNodeKeys {
	allNodeKeys := []AllNodeKeys{}

	for _, n := range nodes {
		api := newNodeAPI(n)

		// Get eth key
		eKey := api.mustExec(api.methods.ListETHKeys)
		ethKeys := mustJSON[[]presenters.ETHKeyResource](eKey)
		ethAddress, err := findFirstGoodEthKeyAddress(chainId, *ethKeys)
		helpers.PanicErr(err)

		// Get aptos account key
		output := &bytes.Buffer{}
		aKeysClient := cmd.NewAptosKeysClient(api.methods)
		err = aKeysClient.ListKeys(&cli.Context{App: api.app})
		helpers.PanicErr(err)
		var aptosKeys []presenters.AptosKeyResource
		helpers.PanicErr(json.Unmarshal(output.Bytes(), &aptosKeys))
		if len(aptosKeys) != 1 {
			helpers.PanicErr(errors.New("node must have single aptos key"))
		}
		aptosAccount := aptosKeys[0].Account
		output.Reset()

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
				fmt.Println("Created OCR2 EVM key bundle", string(cBundle))
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

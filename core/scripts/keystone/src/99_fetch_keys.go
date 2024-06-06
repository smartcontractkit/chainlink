package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/urfave/cli"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func downloadNodePubKeys(chainID int64, pubKeysPath string) []NodeKeys {
	// Check if file exists already, and if so, return the keys
	if _, err := os.Stat(pubKeysPath); err == nil {
		fmt.Println("Loading existing public keys at:", pubKeysPath)
		return mustParseJSON[[]NodeKeys](pubKeysPath)
	}

	nodes := downloadNodeAPICredentialsDefault()
	nodesKeys := mustFetchNodesKeys(chainID, nodes)

	marshalledNodeKeys, err := json.MarshalIndent(nodesKeys, "", " ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(pubKeysPath, marshalledNodeKeys, 0600)
	if err != nil {
		panic(err)
	}
	fmt.Println("Keystone OCR2 public keys have been saved to:", pubKeysPath)

	return nodesKeys
}

// downloadNodeAPICredentialsDefault downloads the node API credentials, or loads them from disk if they already exist
//
// The nodes are sorted by URL. In the case of crib, the bootstrap node is the first node in the list.
func downloadNodeAPICredentialsDefault() []*node {
	return downloadNodeAPICredentials(".cache/NodeList.txt")
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

type ocr2Bundle struct {
	ID                string `json:"id"`
	ChainType         string `json:"chainType"`
	OnchainPublicKey  string `json:"onchainPublicKey"`
	OffchainPublicKey string `json:"offchainPublicKey"`
	ConfigPublicKey   string `json:"configPublicKey"`
}

func mustFetchNodesKeys(chainID int64, nodes []*node) (nca []NodeKeys) {
	for _, n := range nodes {
		output := &bytes.Buffer{}
		client, app := newApp(n, output)

		fmt.Println("Logging in:", n.url)
		loginFs := flag.NewFlagSet("test", flag.ContinueOnError)
		loginFs.Bool("bypass-version-check", true, "")
		loginCtx := cli.NewContext(app, loginFs, nil)
		err := client.RemoteLogin(loginCtx)
		helpers.PanicErr(err)
		output.Reset()

		err = client.ListETHKeys(&cli.Context{
			App: app,
		})
		helpers.PanicErr(err)
		var ethKeys []presenters.ETHKeyResource
		helpers.PanicErr(json.Unmarshal(output.Bytes(), &ethKeys))
		ethAddress, err := findFirstGoodEthKeyAddress(chainID, ethKeys)
		helpers.PanicErr(err)
		output.Reset()

		err = client.ListP2PKeys(&cli.Context{
			App: app,
		})
		helpers.PanicErr(err)
		var p2pKeys []presenters.P2PKeyResource
		helpers.PanicErr(json.Unmarshal(output.Bytes(), &p2pKeys))
		if len(p2pKeys) != 1 {
			helpers.PanicErr(errors.New("node must have single p2p key"))
		}
		peerID := strings.TrimPrefix(p2pKeys[0].PeerID, "p2p_")
		output.Reset()

		var ocr2Bundles []ocr2Bundle
		err = client.ListOCR2KeyBundles(&cli.Context{
			App: app,
		})
		helpers.PanicErr(err)
		helpers.PanicErr(json.Unmarshal(output.Bytes(), &ocr2Bundles))
		ocr2BundleIndex := findEvmOCR2Bundle(ocr2Bundles)
		output.Reset()
		if ocr2BundleIndex == -1 {
			fmt.Println("WARN: node does not have EVM OCR2 bundle, creating one")
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			err = fs.Parse([]string{"evm"})
			helpers.PanicErr(err)
			ocr2CreateBundleCtx := cli.NewContext(app, fs, nil)
			err = client.CreateOCR2KeyBundle(ocr2CreateBundleCtx)
			helpers.PanicErr(err)
			output.Reset()

			err = client.ListOCR2KeyBundles(&cli.Context{
				App: app,
			})
			helpers.PanicErr(err)
			helpers.PanicErr(json.Unmarshal(output.Bytes(), &ocr2Bundles))
			ocr2BundleIndex = findEvmOCR2Bundle(ocr2Bundles)
			output.Reset()
		}

		ocr2Bndl := ocr2Bundles[ocr2BundleIndex]

		err = client.ListCSAKeys(&cli.Context{
			App: app,
		})
		helpers.PanicErr(err)
		var csaKeys []presenters.CSAKeyResource
		helpers.PanicErr(json.Unmarshal(output.Bytes(), &csaKeys))
		csaPubKey, err := findFirstCSAPublicKey(csaKeys)
		helpers.PanicErr(err)
		output.Reset()

		nc := NodeKeys{
			EthAddress:            ethAddress,
			P2PPeerID:             peerID,
			OCR2BundleID:          ocr2Bndl.ID,
			OCR2ConfigPublicKey:   strings.TrimPrefix(ocr2Bndl.ConfigPublicKey, "ocr2cfg_evm_"),
			OCR2OnchainPublicKey:  strings.TrimPrefix(ocr2Bndl.OnchainPublicKey, "ocr2on_evm_"),
			OCR2OffchainPublicKey: strings.TrimPrefix(ocr2Bndl.OffchainPublicKey, "ocr2off_evm_"),
			CSAPublicKey:          csaPubKey,
		}

		nca = append(nca, nc)
	}
	return
}

func findFirstCSAPublicKey(csaKeyResources []presenters.CSAKeyResource) (string, error) {
	for _, r := range csaKeyResources {
		return r.PubKey, nil
	}
	return "", errors.New("did not find any CSA Key Resources")
}

func findEvmOCR2Bundle(ocr2Bundles []ocr2Bundle) int {
	for i, b := range ocr2Bundles {
		if b.ChainType == "evm" {
			return i
		}
	}
	return -1
}

func findFirstGoodEthKeyAddress(chainID int64, ethKeys []presenters.ETHKeyResource) (string, error) {
	for _, ethKey := range ethKeys {
		if ethKey.EVMChainID.Equal(ubig.NewI(chainID)) && !ethKey.Disabled {
			if ethKey.EthBalance.IsZero() {
				fmt.Println("WARN: selected ETH address has zero balance", ethKey.Address)
			}
			return ethKey.Address, nil
		}
	}
	return "", errors.New("did not find an enabled ETH key for the given chain ID")
}

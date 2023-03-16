package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/urfave/cli"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

type ocr2Bundle struct {
	ID                string `json:"id"`
	ChainType         string `json:"chainType"`
	OnchainPublicKey  string `json:"onchainPublicKey"`
	OffchainPublicKey string `json:"offchainPublicKey"`
	ConfigPublicKey   string `json:"configPublicKey"`
}

func mustFetchNodesKeys(chainID int64, nodes []*node) (nca []NodeKeys) {
	for _, node := range nodes {
		output := &bytes.Buffer{}
		client, app := newApp(node, output)

		fmt.Println("Logging in:", node.url)
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
		if ocr2BundleIndex == -1 {
			helpers.PanicErr(errors.New("node must have EVM OCR2 bundle"))
		}
		ocr2Bundle := ocr2Bundles[ocr2BundleIndex]

		nc := NodeKeys{
			EthAddress:            ethAddress,
			P2PPeerID:             peerID,
			OCR2BundleID:          ocr2Bundle.ID,
			OCR2ConfigPublicKey:   strings.TrimPrefix(ocr2Bundle.ConfigPublicKey, "ocr2cfg_evm_"),
			OCR2OnchainPublicKey:  strings.TrimPrefix(ocr2Bundle.OnchainPublicKey, "ocr2on_evm_"),
			OCR2OffchainPublicKey: strings.TrimPrefix(ocr2Bundle.OffchainPublicKey, "ocr2off_evm_"),
		}

		nca = append(nca, nc)
	}
	return
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
		if ethKey.EVMChainID.Equal(utils.NewBigI(chainID)) && !ethKey.Disabled {
			if ethKey.EthBalance.IsZero() {
				fmt.Println("WARN: selected ETH address has zero balance", ethKey.Address)
			}
			return ethKey.Address, nil
		}
	}
	return "", errors.New("did not find an enabled ETH key for the given chain ID")
}

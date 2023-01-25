package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli"

	clcmd "github.com/smartcontractkit/chainlink/core/cmd"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func newApp(remoteNodeURL string, writer io.Writer) (*clcmd.Client, *cli.App) {
	prompter := clcmd.NewTerminalPrompter()
	client := &clcmd.Client{
		Renderer:                       clcmd.RendererJSON{Writer: writer},
		AppFactory:                     clcmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:          clcmd.TerminalKeyStoreAuthenticator{Prompter: prompter},
		FallbackAPIInitializer:         clcmd.NewPromptingAPIInitializer(prompter),
		Runner:                         clcmd.ChainlinkRunner{},
		PromptingSessionRequestBuilder: clcmd.NewPromptingSessionRequestBuilder(prompter),
		ChangePasswordPrompter:         clcmd.NewChangePasswordPrompter(),
		PasswordPrompter:               clcmd.NewPasswordPrompter(),
	}
	app := clcmd.NewApp(client)
	fs := flag.NewFlagSet("blah", flag.ContinueOnError)
	fs.Bool("json", true, "")
	fs.String("remote-node-url", remoteNodeURL, "")
	helpers.PanicErr(app.Before(cli.NewContext(nil, fs, nil)))
	// overwrite renderer since it's set to stdout after Before() is called
	client.Renderer = clcmd.RendererJSON{Writer: writer}
	return client, app
}

var (
	remoteNodeURLs = flag.String("remote-node-urls", "", "remote node URL")
	checkMarkEmoji = "✅"
	xEmoji         = "❌"
	infoEmoji      = "ℹ️"
)

type ocr2Bundle struct {
	ID                string `json:"id"`
	ChainType         string `json:"chainType"`
	OnchainPublicKey  string `json:"onchainPublicKey"`
	OffchainPublicKey string `json:"offchainPublicKey"`
	ConfigPublicKey   string `json:"configPublicKey"`
}

func main() {
	flag.Parse()

	if remoteNodeURLs == nil {
		fmt.Println("flag -remote-node-urls required")
		os.Exit(1)
	}

	urls := strings.Split(*remoteNodeURLs, ",")
	var (
		allDKGSignKeys         []string
		allDKGEncryptKeys      []string
		allOCR2KeyIDs          []string
		allOCR2OffchainPubkeys []string
		allOCR2OnchainPubkeys  []string
		allOCR2ConfigPubkeys   []string
		allETHKeys             []string
		allPeerIDs             []string
	)
	for _, remoteNodeURL := range urls {
		output := &bytes.Buffer{}
		client, app := newApp(remoteNodeURL, output)

		// login first to establish the session
		fmt.Println("logging in to:", remoteNodeURL)
		loginFs := flag.NewFlagSet("test", flag.ContinueOnError)
		loginFs.String("file", "", "")
		loginFs.Bool("bypass-version-check", true, "")
		loginCtx := cli.NewContext(app, loginFs, nil)
		err := client.RemoteLogin(loginCtx)
		helpers.PanicErr(err)
		output.Reset()
		fmt.Println()

		// check for DKG signing keys
		err = clcmd.NewDKGSignKeysClient(client).ListKeys(&cli.Context{
			App: app,
		})
		helpers.PanicErr(err)
		var dkgSignKeys []presenters.DKGSignKeyResource
		helpers.PanicErr(json.Unmarshal(output.Bytes(), &dkgSignKeys))
		switch len(dkgSignKeys) {
		case 1:
			fmt.Println(checkMarkEmoji, "found 1 DKG sign key on", remoteNodeURL)
		case 0:
			fmt.Println(xEmoji, "did not find any DKG sign keys on", remoteNodeURL, ", please create one")
		default:
			fmt.Println(infoEmoji, "found more than 1 DKG sign key on", remoteNodeURL, ", consider removing all but one")
		}
		output.Reset()
		fmt.Println()

		// check for DKG encryption keys
		err = clcmd.NewDKGEncryptKeysClient(client).ListKeys(&cli.Context{
			App: app,
		})
		helpers.PanicErr(err)
		var dkgEncryptKeys []presenters.DKGEncryptKeyResource
		helpers.PanicErr(json.Unmarshal(output.Bytes(), &dkgEncryptKeys))
		switch len(dkgEncryptKeys) {
		case 1:
			fmt.Println(checkMarkEmoji, "found 1 DKG encrypt key on", remoteNodeURL)
		case 0:
			fmt.Println(xEmoji, "did not find any DKG encrypt keys on", remoteNodeURL, ", please create one")
		default:
			fmt.Println(infoEmoji, "found more than 1 DKG encrypt key on", remoteNodeURL, ", consider removing all but one")
		}
		output.Reset()
		fmt.Println()

		// check for OCR2 keys
		err = client.ListOCR2KeyBundles(&cli.Context{
			App: app,
		})
		helpers.PanicErr(err)
		var ocr2Keys []ocr2Bundle
		helpers.PanicErr(json.Unmarshal(output.Bytes(), &ocr2Keys))
		ethBundle := func() *ocr2Bundle {
			for _, b := range ocr2Keys {
				if b.ChainType == "evm" {
					return &b
				}
			}
			return nil
		}()
		if ethBundle != nil {
			fmt.Println(checkMarkEmoji, "found ocr evm key bundle on", remoteNodeURL)
		} else {
			fmt.Println(xEmoji, "did not find ocr evm key bundle on", remoteNodeURL, ", please create one")
		}
		output.Reset()
		fmt.Println()

		// check for ETH keys
		err = client.ListETHKeys(&cli.Context{
			App: app,
		})
		helpers.PanicErr(err)
		var ethKeys []presenters.ETHKeyResource
		helpers.PanicErr(json.Unmarshal(output.Bytes(), &ethKeys))
		switch {
		case len(ethKeys) >= 5:
			fmt.Println(checkMarkEmoji, "found", len(ethKeys), "eth keys on", remoteNodeURL)
		case len(ethKeys) < 5:
			fmt.Println(xEmoji, "found only", len(ethKeys), "eth keys on", remoteNodeURL, ", consider creating more")
		}
		output.Reset()
		fmt.Println()

		// check for peer ids
		err = client.ListP2PKeys(&cli.Context{
			App: app,
		})
		helpers.PanicErr(err)
		var p2pKeys []presenters.P2PKeyResource
		helpers.PanicErr(json.Unmarshal(output.Bytes(), &p2pKeys))
		switch len(p2pKeys) {
		case 1:
			fmt.Println(checkMarkEmoji, "found P2P key on", remoteNodeURL)
		case 0:
			fmt.Println(xEmoji, "no P2P keys found on", remoteNodeURL, ", please create one")
		default:
			fmt.Println(infoEmoji, "found", len(p2pKeys), "P2P keys on", remoteNodeURL, ", consider removing all but one")
		}
		output.Reset()
		fmt.Println()

		for _, dkgSign := range dkgSignKeys {
			allDKGSignKeys = append(allDKGSignKeys, dkgSign.PublicKey)
		}
		for _, dkgEncrypt := range dkgEncryptKeys {
			allDKGEncryptKeys = append(allDKGEncryptKeys, dkgEncrypt.PublicKey)
		}
		for _, ocr2Bundle := range ocr2Keys {
			if ocr2Bundle.ChainType == "evm" {
				allOCR2KeyIDs = append(allOCR2KeyIDs, ocr2Bundle.ID)
				allOCR2ConfigPubkeys = append(allOCR2ConfigPubkeys, strings.TrimPrefix(ocr2Bundle.ConfigPublicKey, "ocr2cfg_evm_"))
				allOCR2OffchainPubkeys = append(allOCR2OffchainPubkeys, strings.TrimPrefix(ocr2Bundle.OffchainPublicKey, "ocr2off_evm_"))
				allOCR2OnchainPubkeys = append(allOCR2OnchainPubkeys, strings.TrimPrefix(ocr2Bundle.OnchainPublicKey, "ocr2on_evm_"))
			}
		}
		for _, ethKey := range ethKeys {
			allETHKeys = append(allETHKeys, ethKey.Address)
		}
		for _, peerKey := range p2pKeys {
			allPeerIDs = append(allPeerIDs, strings.TrimPrefix(peerKey.PeerID, "p2p_"))
		}
	}

	fmt.Println("------------- NODE INFORMATION -------------")
	fmt.Println("DKG sign keys:", strings.Join(allDKGSignKeys, ","))
	fmt.Println("DKG encrypt keys:", strings.Join(allDKGEncryptKeys, ","))
	fmt.Println("OCR2 key IDs:", strings.Join(allOCR2KeyIDs, ","))
	fmt.Println("OCR2 config public keys:", strings.Join(allOCR2ConfigPubkeys, ","))
	fmt.Println("OCR2 onchain public keys:", strings.Join(allOCR2OnchainPubkeys, ","))
	fmt.Println("OCR2 offchain public keys:", strings.Join(allOCR2OffchainPubkeys, ","))
	fmt.Println("ETH addresses:", strings.Join(allETHKeys, ","))
	fmt.Println("Peer IDs:", strings.Join(allPeerIDs, ","))
}

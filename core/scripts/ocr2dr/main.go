package main

// Inspired by ocr2vrf/readiness tool.

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v3"

	clcmd "github.com/smartcontractkit/chainlink/core/cmd"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func newApp(r remote, writer io.Writer) (*clcmd.Client, *cli.App) {
	client := &clcmd.Client{
		Renderer:                       clcmd.RendererJSON{Writer: writer},
		AppFactory:                     clcmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:          clcmd.TerminalKeyStoreAuthenticator{Prompter: r},
		FallbackAPIInitializer:         clcmd.NewPromptingAPIInitializer(r),
		Runner:                         clcmd.ChainlinkRunner{},
		PromptingSessionRequestBuilder: clcmd.NewPromptingSessionRequestBuilder(r),
		ChangePasswordPrompter:         clcmd.NewChangePasswordPrompter(),
		PasswordPrompter:               clcmd.NewPasswordPrompter(),
	}
	app := clcmd.NewApp(client)
	fs := flag.NewFlagSet("blah", flag.ContinueOnError)
	fs.String("remote-node-url", fmt.Sprintf("https://%s", r.host), "")
	helpers.PanicErr(app.Before(cli.NewContext(nil, fs, nil)))
	// overwrite renderer since it's set to stdout after Before() is called
	client.Renderer = clcmd.RendererJSON{Writer: writer}
	return client, app
}

func mustReadRemotes(path string) []remote {
	remotesInput, err := readLines(path)
	if err != nil {
		fmt.Printf("Failed to read remotes from file: %s, err: %v\n", path, err)
		os.Exit(1)
	}
	var remotes []remote
	for ln, r := range remotesInput {
		rr := strings.TrimSpace(r)
		if len(rr) == 0 {
			continue
		}
		s := strings.Split(rr, " ")
		if len(s) != 3 {
			fmt.Printf("Failed to parse remote: %s, at line: %d\n", r, ln)
			os.Exit(1)
		}
		remotes = append(remotes, remote{
			host:     s[0],
			login:    s[1],
			password: s[2],
		})
	}
	return remotes
}

func mustReadConfig() *config {
	data, err := os.ReadFile(configFile)
	helpers.PanicErr(err)
	c := &config{}
	err = yaml.Unmarshal(data, c)
	helpers.PanicErr(err)
	return c
}

var (
	checkMarkEmoji = "✅"
	xEmoji         = "❌"
	infoEmoji      = "ℹ️"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Check USAGE.md")
		os.Exit(1)
	}

	config := mustReadConfig()
	fmt.Printf("Config chain ID: %d\n", config.ChainID)

	remotes := mustReadRemotes(os.Args[1])
	var nodes []Node

	for _, remote := range remotes {
		var (
			allOCR2KeyIDs          []string
			allOCR2OffchainPubkeys []string
			allOCR2OnchainPubkeys  []string
			allOCR2ConfigPubkeys   []string
			allETHKeys             []string
			allPeerIDs             []string
		)

		remoteNodeURL := fmt.Sprintf("https://%s", remote.host)
		output := &bytes.Buffer{}
		client, app := newApp(remote, output)

		// login first to establish the session
		fmt.Println("Logging in to:", remoteNodeURL)
		loginFs := flag.NewFlagSet("test", flag.ContinueOnError)
		//loginFs.String("file", "", "")
		loginFs.Bool("bypass-version-check", true, "")
		loginCtx := cli.NewContext(app, loginFs, nil)
		err := client.RemoteLogin(loginCtx)
		helpers.PanicErr(err)
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

		node := Node{
			Host:                remoteNodeURL,
			ETHKeys:             allETHKeys,
			P2PPeerIDS:          allPeerIDs,
			OCR2KeyIDs:          allOCR2KeyIDs,
			OCR2ConfigPubKeys:   allOCR2ConfigPubkeys,
			OCR2OffchainPubKeys: allOCR2OffchainPubkeys,
			OCR2OnchainPubKeys:  allOCR2OnchainPubkeys,
		}
		nodes = append(nodes, node)
	}

	processJobSpecs(config, nodes)

	js, err := json.Marshal(nodes)
	helpers.PanicErr(err)

	err = os.WriteFile(filepath.Join(artefactsDir, clusterFile), js, 0644)
	helpers.PanicErr(err)

	fmt.Println(string(js))
}

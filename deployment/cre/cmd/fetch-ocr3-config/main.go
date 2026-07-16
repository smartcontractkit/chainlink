// fetch-ocr3-config reads the current OCR3 configuration from an on-chain
// OCR3Capability contract and prints it in a human-readable format that
// matches the offchainConfig block used in CRE pipeline input YAMLs.
//
// Usage:
//
//	go run ./deployment/cre/cmd/fetch-ocr3-config \
//	  --rpc-url  https://mainnet.infura.io/v3/<key> \
//	  --contract 0xABCD... \
//	  --plugin   consensus|dontime|chain-cap \
//	  [--format  yaml|json] \
//	  [--out     /path/to/output.yaml]
//
// --rpc-url may also be supplied via the ETH_RPC_URL environment variable.
// The plugin type can be found from the qualifier field in address_refs.json.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/yaml.v3"

	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"

	creview "github.com/smartcontractkit/chainlink/deployment/cre/view"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	rpcURL := flag.String("rpc-url", os.Getenv("ETH_RPC_URL"), "EVM JSON-RPC endpoint URL (or set ETH_RPC_URL)")
	contractAddr := flag.String("contract", "", "OCR3Capability contract address (0x hex)")
	plugin := flag.String("plugin", "", fmt.Sprintf("Plugin type: %s | %s | %s",
		creview.PluginTypeConsensus, creview.PluginTypeDontime, creview.PluginTypeChainCap))
	format := flag.String("format", "yaml", "Output format: yaml or json")
	outPath := flag.String("out", "", "Write output to this file instead of stdout")
	flag.Parse()

	var errs []string
	if *rpcURL == "" {
		errs = append(errs, "--rpc-url is required (or set ETH_RPC_URL)")
	}
	if *contractAddr == "" {
		errs = append(errs, "--contract is required")
	}
	if *plugin == "" {
		errs = append(errs, fmt.Sprintf("--plugin is required: %s | %s | %s",
			creview.PluginTypeConsensus, creview.PluginTypeDontime, creview.PluginTypeChainCap))
	}
	if *format != "yaml" && *format != "json" {
		errs = append(errs, fmt.Sprintf("--format must be yaml or json, got %q", *format))
	}
	if len(errs) > 0 {
		flag.Usage()
		return errors.New(strings.Join(errs, "; "))
	}

	ctx := context.Background()

	client, err := ethclient.DialContext(ctx, *rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC %q: %w", *rpcURL, err)
	}
	defer client.Close()

	addr := common.HexToAddress(*contractAddr)
	contract, err := ocr3_capability.NewOCR3Capability(addr, client)
	if err != nil {
		return fmt.Errorf("failed to instantiate OCR3Capability at %s: %w", addr.Hex(), err)
	}

	view, err := creview.GenerateOCR3ConfigViewForPlugin(ctx, *contract, creview.PluginType(*plugin))
	if err != nil {
		return fmt.Errorf("failed to fetch OCR3 config: %w", err)
	}

	var output []byte
	switch *format {
	case "yaml":
		output, err = yaml.Marshal(view)
	case "json":
		output, err = json.MarshalIndent(view, "", "  ")
	}
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	w := os.Stdout
	if *outPath != "" {
		f, ferr := os.Create(*outPath)
		if ferr != nil {
			return fmt.Errorf("failed to create output file %q: %w", *outPath, ferr)
		}
		defer f.Close()
		w = f
	}

	if _, err = w.Write(output); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

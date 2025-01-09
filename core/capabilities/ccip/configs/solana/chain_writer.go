package solana

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainwriter"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/codec"
)

//go:embed ccip_router.json
var ccipRouter string

func GetSolanaChainWriterConfig(fromAddress string) (chainwriter.ChainWriterConfig, error) {
	// TODO once on-chain account lookup address are available, the config will be updated

	// check fromAddress
	_, err := solana.PublicKeyFromBase58(fromAddress)
	if err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("invalid from address %s: %w", fromAddress, err)
	}

	// validate CCIP Router IDL, errors not expected
	var idl codec.IDL
	if err = json.Unmarshal([]byte(ccipRouter), &idl); err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("unexpected error: invalid CCIP Router IDL, error: %w", err)
	}

	solConfig := chainwriter.ChainWriterConfig{
		Programs: map[string]chainwriter.ProgramConfig{
			"ccip-router": {
				Methods: map[string]chainwriter.MethodConfig{
					"execute": {
						FromAddress:        fromAddress,
						InputModifications: nil,
						ChainSpecificName:  "execute"},
					"commit": {
						FromAddress:        fromAddress,
						InputModifications: nil,
						ChainSpecificName:  "commit"},
				},
				IDL: ccipRouter},
		},
	}

	return solConfig, nil
}

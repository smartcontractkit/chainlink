package aptosconfig

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"

	"github.com/smartcontractkit/chainlink-aptos/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-aptos/relayer/chainwriter"
	"github.com/smartcontractkit/chainlink-aptos/relayer/utils"
)

func GetChainWriterConfig(publicKeyStr string) (chainwriter.ChainWriterConfig, error) {
	fromAddress, err := utils.HexPublicKeyToAddress(publicKeyStr)
	if err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("failed to parse Aptos address from public key %s: %w", publicKeyStr, err)
	}

	return chainwriter.ChainWriterConfig{
		Modules: map[string]*chainwriter.ChainWriterModule{
			consts.ContractNameOffRamp: {
				Name: "offramp",
				Functions: map[string]*chainwriter.ChainWriterFunction{
					consts.MethodCommit: {
						Name:        "commit",
						PublicKey:   publicKeyStr,
						FromAddress: fromAddress.String(),
						Params: []config.AptosFunctionParam{
							{
								Name:     "ReportContext",
								Type:     "vector<vector<u8>>",
								Required: true,
							},
							{
								Name:     "Report",
								Type:     "vector<u8>",
								Required: true,
							},
							{
								Name:     "Signatures",
								Type:     "vector<vector<u8>>",
								Required: true,
							},
						},
					},
					consts.MethodExecute: {
						Name:        "execute",
						PublicKey:   publicKeyStr,
						FromAddress: fromAddress.String(),
						Params: []config.AptosFunctionParam{
							{
								Name:     "ReportContext",
								Type:     "vector<vector<u8>>",
								Required: true,
							},
							{
								Name:     "Report",
								Type:     "vector<u8>",
								Required: true,
							},
						},
					},
				},
			},
		},
		FeeStrategy: chainwriter.DefaultFeeStrategy,
	}, nil
}

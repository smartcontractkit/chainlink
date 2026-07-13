package aptosconfig

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-aptos/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/types/aptos"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip/consts"
)

func GetChainWriterConfig(publicKeyStr string) (aptos.ContractWriterConfig, error) {
	fromAddress, err := codec.HexPublicKeyToAddress(publicKeyStr)
	if err != nil {
		return aptos.ContractWriterConfig{}, fmt.Errorf("failed to parse Aptos address from public key %s: %w", publicKeyStr, err)
	}

	return aptos.ContractWriterConfig{
		Modules: map[string]*aptos.ContractWriterModule{
			consts.ContractNameOffRamp: {
				Name: "offramp",
				Functions: map[string]*aptos.ContractWriterFunction{
					consts.MethodCommit: {
						Name:        "commit",
						PublicKey:   publicKeyStr,
						FromAddress: fromAddress,
						Params: []aptos.FunctionParam{
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
						FromAddress: fromAddress,
						Params: []aptos.FunctionParam{
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
		FeeStrategy: aptos.DefaultFeeStrategy,
	}, nil
}

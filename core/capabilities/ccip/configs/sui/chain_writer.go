package suiconfig

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"

	_ "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter" // Register Sui chainwriter
	chainwriter "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

func GetChainWriterConfig(publicKeyStr string) (chainwriter.ChainWriterConfig, error) {
	// returns 32 byte pubKey
	rawPubKey, err := hex.DecodeString(publicKeyStr)
	if err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("invalid public key hex %q: %w", publicKeyStr, err)
	}

	pubKeyBytes := ed25519.PublicKey(rawPubKey)
	nonMutable := false

	return chainwriter.ChainWriterConfig{
		Modules: map[string]*chainwriter.ChainWriterModule{
			consts.ContractNameOffRamp: {
				Name: "offramp",
				Functions: map[string]*chainwriter.ChainWriterFunction{
					consts.MethodCommit: {
						Name:      "commit",
						PublicKey: pubKeyBytes,
						Params:    []codec.SuiFunctionParam{},
						PTBCommands: []chainwriter.ChainWriterPTBCommand{
							{
								Type:     codec.SuiPTBCommandMoveCall,
								ModuleId: strPtr("offramp"),
								Function: strPtr("commit"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "ref",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:     "state",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:      "clock",
										Type:      "object_id",
										Required:  true,
										IsMutable: &nonMutable,
									},
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
						},
					},
					consts.MethodExecute: {
						Name:      "execute",
						PublicKey: pubKeyBytes,
						Params:    []codec.SuiFunctionParam{},
						PTBCommands: []chainwriter.ChainWriterPTBCommand{
							{
								Type:     codec.SuiPTBCommandMoveCall,
								ModuleId: strPtr("offramp"),
								Function: strPtr("init_execute"),
								Params: []codec.SuiFunctionParam{
									{
										Name:      "ref",
										Type:      "object_id",
										Required:  true,
										IsMutable: &nonMutable,
									},
									{
										Name:     "state",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:      "clock",
										Type:      "object_id",
										Required:  true,
										IsMutable: &nonMutable,
									},
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
										Name:     "token_receiver",
										Type:     "vector<u8>",
										Required: true,
									},
								},
							},
							{
								Type:     codec.SuiPTBCommandMoveCall,
								ModuleId: strPtr("offramp"),
								Function: strPtr("finish_execute"),
								Params: []codec.SuiFunctionParam{
									{
										Name:      "ref",
										Type:      "object_id",
										Required:  true,
										IsMutable: &nonMutable,
									},
									{
										Name:     "state",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:     "receiver_params",
										Type:     "ptb_dependency",
										Required: true,
										PTBDependency: &codec.PTBCommandDependency{
											CommandIndex: uint16(0),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		// TODO: come back to it
		// FeeStrategy: chainwriter.DefaultFeeStrategy,
	}, nil
}

// Helper function to convert a string to a string pointer
func strPtr(s string) *string {
	return &s
}

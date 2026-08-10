package suiconfig

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip/consts"
	types "github.com/smartcontractkit/chainlink-common/pkg/types/sui"
)

func GetChainWriterConfig(publicKeyStr string) (types.ChainWriterConfig, error) {
	// returns 32 byte pubKey
	rawPubKey, err := hex.DecodeString(publicKeyStr)
	if err != nil {
		return types.ChainWriterConfig{}, fmt.Errorf("invalid public key hex %q: %w", publicKeyStr, err)
	}

	pubKeyBytes := ed25519.PublicKey(rawPubKey)
	nonMutable := false

	return types.ChainWriterConfig{
		Modules: map[string]*types.ChainWriterModule{
			consts.ContractNameOffRamp: {
				Name: "offramp",
				Functions: map[string]*types.ChainWriterFunction{
					consts.MethodCommit: {
						Name:      "commit",
						PublicKey: pubKeyBytes,
						Params:    []types.SuiFunctionParam{},
						PTBCommands: []types.ChainWriterPTBCommand{
							{
								Type:     types.SuiPTBCommandMoveCall,
								ModuleId: new("offramp"),
								Function: new("commit"),
								Params: []types.SuiFunctionParam{
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
						Params:    []types.SuiFunctionParam{},
						PTBCommands: []types.ChainWriterPTBCommand{
							{
								Type:     types.SuiPTBCommandMoveCall,
								ModuleId: new("offramp"),
								Function: new("init_execute"),
								Params: []types.SuiFunctionParam{
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
								Type:     types.SuiPTBCommandMoveCall,
								ModuleId: new("offramp"),
								Function: new("finish_execute"),
								Params: []types.SuiFunctionParam{
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
										PTBDependency: &types.PTBCommandDependency{
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

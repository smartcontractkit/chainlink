package proposalutils

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_home"
)

func Test_AnalyzeRmnHomeSetConfig(t *testing.T) {
	dataEncoded := "EY26xQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAaAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAASAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA+mHAtKPceHWHsUU/gZ3gbQUSqxoUZSI6A1wsufMtT1Ds0lEhXpCRE0bKF15QNbgZoIwngeB5DppZ27p+FxzwtFJl4urMXKidA5eY8JSHVPtI1h4ZiTOwbcUx4ju2t4QG2jZ8/J045hVKKkjqbrOPTnFXdd4sYe0Egs4TMSMiShZgAHkr6j6Bhgz3+4p0ob33+C3inZ/H6HxUrzaTVNBUGFa+O4yMWpkJ/FppF/cGzqRqFrEWePBy5HGnhxR8MH8IQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA3kG6T8nZGtkAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAHAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAzPCjGiIfPJsAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAHAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMEYRtq/7p2oAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAHAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
	data := make([]byte, base64.StdEncoding.DecodedLen(len(dataEncoded)))
	_, err := base64.StdEncoding.Decode(data, []byte(dataEncoded))
	require.NoError(t, err)
	_abi, err := abi.JSON(strings.NewReader(rmn_home.RMNHomeABI))
	require.NoError(t, err)
	decoder := NewTxCallDecoder(nil)
	analyzeResult, err := decoder.Analyze(common.Address{}.Hex(), &_abi, data)
	require.NoError(t, err)
	assert.Equal(t, _abi.Methods["setCandidate"].String(), analyzeResult.Method)
	assert.Equal(t, "staticConfig", analyzeResult.Inputs[0].Name)
	assert.Equal(t, StructArgument{
		Fields: []NamedArgument{
			{
				Name: "Nodes",
				Value: ArrayArgument{
					Elements: []Argument{
						StructArgument{
							Fields: []NamedArgument{
								{
									Name:  "PeerId",
									Value: BytesArgument{Value: hexutil.MustDecode("0xe98702d28f71e1d61ec514fe067781b4144aac68519488e80d70b2e7ccb53d43")},
								},
								{
									Name:  "OffchainPublicKey",
									Value: BytesArgument{Value: hexutil.MustDecode("0xb34944857a42444d1b285d7940d6e06682309e0781e43a69676ee9f85c73c2d1")},
								},
							},
						},
						StructArgument{
							Fields: []NamedArgument{
								{
									Name:  "PeerId",
									Value: BytesArgument{Value: hexutil.MustDecode("0x49978bab3172a2740e5e63c2521d53ed2358786624cec1b714c788eedade101b")},
								},
								{
									Name:  "OffchainPublicKey",
									Value: BytesArgument{Value: hexutil.MustDecode("0x68d9f3f274e3985528a923a9bace3d39c55dd778b187b4120b384cc48c892859")},
								},
							},
						},
						StructArgument{
							Fields: []NamedArgument{
								{
									Name:  "PeerId",
									Value: BytesArgument{Value: hexutil.MustDecode("0x8001e4afa8fa061833dfee29d286f7dfe0b78a767f1fa1f152bcda4d53415061")},
								},
								{
									Name:  "OffchainPublicKey",
									Value: BytesArgument{Value: hexutil.MustDecode("0x5af8ee32316a6427f169a45fdc1b3a91a85ac459e3c1cb91c69e1c51f0c1fc21")},
								},
							},
						},
					},
				},
			},
			{
				Name:  "OffchainConfig",
				Value: BytesArgument{Value: hexutil.MustDecode("0x")},
			},
		},
	}, analyzeResult.Inputs[0].Value)
	assert.Equal(t, "dynamicConfig", analyzeResult.Inputs[1].Name)
	assert.Equal(t, StructArgument{
		Fields: []NamedArgument{
			{
				Name: "SourceChains",
				Value: ArrayArgument{
					Elements: []Argument{
						StructArgument{
							Fields: []NamedArgument{
								{
									Name:  "ChainSelector",
									Value: ChainSelectorArgument{Value: 16015286601757825753},
								},
								{
									Name:  "FObserve",
									Value: SimpleArgument{Value: "1"},
								},
								{
									Name:  "ObserverNodesBitmap",
									Value: SimpleArgument{Value: "7"},
								},
							},
						},
						StructArgument{
							Fields: []NamedArgument{
								{
									Name:  "ChainSelector",
									Value: ChainSelectorArgument{Value: 14767482510784806043},
								},
								{
									Name:  "FObserve",
									Value: SimpleArgument{Value: "1"},
								},
								{
									Name:  "ObserverNodesBitmap",
									Value: SimpleArgument{Value: "7"},
								},
							},
						},
						StructArgument{
							Fields: []NamedArgument{
								{
									Name:  "ChainSelector",
									Value: ChainSelectorArgument{Value: 3478487238524512106},
								},
								{
									Name:  "FObserve",
									Value: SimpleArgument{Value: "1"},
								},
								{
									Name:  "ObserverNodesBitmap",
									Value: SimpleArgument{Value: "7"},
								},
							},
						},
					},
				},
			},
			{
				Name:  "OffchainConfig",
				Value: BytesArgument{Value: hexutil.MustDecode("0x")},
			},
		},
	}, analyzeResult.Inputs[1].Value)
	assert.Equal(t, "digestToOverwrite", analyzeResult.Inputs[2].Name)
	assert.Equal(t, BytesArgument{Value: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000")}, analyzeResult.Inputs[2].Value)
	assert.Equal(t, "newConfigDigest", analyzeResult.Outputs[0].Name)
	assert.Equal(t, BytesArgument{Value: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000060")}, analyzeResult.Outputs[0].Value)
}

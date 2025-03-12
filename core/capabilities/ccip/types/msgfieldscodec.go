package types

import cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

// AddressCodec is an interface that defines the methods for encoding and decoding addresses
type AddressCodec interface {
	AddressBytesToString([]byte) (string, error)
	AddressStringToBytes(string) ([]byte, error)
}

// ExtraDataCodec is an interface for decoding extra args and dest exec data into a chain-agnostic map[string]any representation
type ExtraDataCodec interface {
	// DecodeExtraArgs reformat bytes into a chain agnostic map[string]any representation for extra args
	DecodeExtraArgs(extraArgs cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error)
	// DecodeTokenAmountDestExecData reformat bytes to chain-agnostic map[string]any for tokenAmount DestExecData field
	DecodeTokenAmountDestExecData(destExecData cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error)
}

// ChainSpecificExtraDataDecoder is an interface for decoding chain specific extra args and dest exec data into a map[string]any representation
type ChainSpecificExtraDataDecoder interface {
	// DecodeExtraArgsToMap reformat bytes into a chain agnostic map[string]any representation for extra args
	DecodeExtraArgsToMap(extraArgs cciptypes.Bytes) (map[string]any, error)
	// DecodeDestExecDataToMap reformat bytes into a chain agnostic map[string]interface{} representation for dest exec data
	DecodeDestExecDataToMap(destExecData cciptypes.Bytes) (map[string]any, error)
}

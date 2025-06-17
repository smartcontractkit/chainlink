package common

import cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

// ChainSpecificAddressCodec is an interface that defines the methods for encoding and decoding addresses for a specific chain
type ChainSpecificAddressCodec interface {
	// AddressBytesToString converts an address from bytes to string
	AddressBytesToString([]byte) (string, error)
	// AddressStringToBytes converts an address from string to bytes
	AddressStringToBytes(string) ([]byte, error)
	// OracleIDAsAddressBytes returns a valid address for this chain family with the bytes set to the given oracle ID.
	OracleIDAsAddressBytes(oracleID uint8) ([]byte, error)
}

// SourceChainExtraDataCodec is an interface for decoding source chain specific extra args and dest exec data into a map[string]any representation for a specific chain
// For chain A to chain B message, this interface will be the chain A specific codec
type SourceChainExtraDataCodec interface {
	// DecodeExtraArgsToMap reformat bytes into a chain agnostic map[string]any representation for extra args
	DecodeExtraArgsToMap(extraArgs cciptypes.Bytes) (map[string]any, error)
	// DecodeDestExecDataToMap reformat bytes into a chain agnostic map[string]interface{} representation for dest exec data
	DecodeDestExecDataToMap(destExecData cciptypes.Bytes) (map[string]any, error)
}

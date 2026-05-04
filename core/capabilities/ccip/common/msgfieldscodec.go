package common

import cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

// ProtocolAddressCodec defines the canonical OCR3 ContractConfig encoding for a
// chain family. It is used by the bootstrap path, where chain SDKs and LOOP
// relayers may be unavailable.
type ProtocolAddressCodec interface {
	// OracleIDAsAddressBytes returns a valid address for this chain family with the bytes set to the given oracle ID.
	OracleIDAsAddressBytes(oracleID uint8) ([]byte, error)
	// TransmitterBytesToString converts a transmitter account from bytes to string
	TransmitterBytesToString([]byte) (string, error)
}

// ChainSpecificAddressCodec is the full chain-integration codec used by the
// plugin path. For LOOPP-backed chains, it flows through CCIPProvider.Codec().
type ChainSpecificAddressCodec interface {
	// AddressBytesToString converts an address from bytes to string
	AddressBytesToString([]byte) (string, error)
	// AddressStringToBytes converts an address from string to bytes
	AddressStringToBytes(string) ([]byte, error)
	// OracleIDAsAddressBytes returns a valid address for this chain family with the bytes set to the given oracle ID.
	OracleIDAsAddressBytes(oracleID uint8) ([]byte, error)
	// TransmitterBytesToString converts a transmitter account from bytes to string
	TransmitterBytesToString([]byte) (string, error)
}

// SourceChainExtraDataCodec is an interface for decoding source chain specific extra args and dest exec data into a map[string]any representation for a specific chain
// For chain A to chain B message, this interface will be the chain A specific codec
type SourceChainExtraDataCodec interface {
	// DecodeExtraArgsToMap reformat bytes into a chain agnostic map[string]any representation for extra args
	DecodeExtraArgsToMap(extraArgs cciptypes.Bytes) (map[string]any, error)
	// DecodeDestExecDataToMap reformat bytes into a chain agnostic map[string]interface{} representation for dest exec data
	DecodeDestExecDataToMap(destExecData cciptypes.Bytes) (map[string]any, error)
}

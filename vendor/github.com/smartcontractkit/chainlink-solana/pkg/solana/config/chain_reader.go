package config

import (
	"encoding/json"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings/binary"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type ChainReader struct {
	Namespaces map[string]ChainReaderMethods `json:"namespaces" toml:"namespaces"`
}

type ChainReaderMethods struct {
	Methods map[string]ChainDataReader `json:"methods" toml:"methods"`
}

type ChainDataReader struct {
	AnchorIDL string `json:"anchorIDL" toml:"anchorIDL"`
	// Encoding defines the type of encoding used for on-chain data. Currently supported
	// are 'borsh' and 'bincode'.
	Encoding   EncodingType           `json:"encoding" toml:"encoding"`
	Procedures []ChainReaderProcedure `json:"procedures" toml:"procedures"`
}

type EncodingType int

const (
	EncodingTypeBorsh EncodingType = iota
	EncodingTypeBincode

	encodingTypeBorshStr   = "borsh"
	encodingTypeBincodeStr = "bincode"
)

func (t EncodingType) MarshalJSON() ([]byte, error) {
	switch t {
	case EncodingTypeBorsh:
		return json.Marshal(encodingTypeBorshStr)
	case EncodingTypeBincode:
		return json.Marshal(encodingTypeBincodeStr)
	default:
		return nil, fmt.Errorf("%w: unrecognized encoding type: %d", types.ErrInvalidConfig, t)
	}
}

func (t *EncodingType) UnmarshalJSON(data []byte) error {
	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("%w: %s", types.ErrInvalidConfig, err.Error())
	}

	switch str {
	case encodingTypeBorshStr:
		*t = EncodingTypeBorsh
	case encodingTypeBincodeStr:
		*t = EncodingTypeBincode
	default:
		return fmt.Errorf("%w: unrecognized encoding type: %s", types.ErrInvalidConfig, str)
	}

	return nil
}

type RPCOpts struct {
	Encoding   *solana.EncodingType `json:"encoding,omitempty"`
	Commitment *rpc.CommitmentType  `json:"commitment,omitempty"`
	DataSlice  *rpc.DataSlice       `json:"dataSlice,omitempty"`
}

type ChainReaderProcedure chainDataProcedureFields

type chainDataProcedureFields struct {
	// IDLAccount refers to the account defined in the IDL.
	IDLAccount string `json:"idlAccount,omitempty"`
	// OutputModifications provides modifiers to convert chain data format to custom
	// output formats.
	OutputModifications codec.ModifiersConfig `json:"outputModifications,omitempty"`
	// RPCOpts provides optional configurations for commitment, encoding, and data
	// slice offsets.
	RPCOpts *RPCOpts `json:"rpcOpts,omitempty"`
}

// BuilderForEncoding returns a builder for the encoding configuration. Defaults to little endian.
func BuilderForEncoding(eType EncodingType) encodings.Builder {
	switch eType {
	case EncodingTypeBorsh:
		return binary.LittleEndian()
	case EncodingTypeBincode:
		return binary.BigEndian()
	default:
		return binary.LittleEndian()
	}
}

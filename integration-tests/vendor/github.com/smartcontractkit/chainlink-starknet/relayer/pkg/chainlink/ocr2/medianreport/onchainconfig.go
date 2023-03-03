package medianreport

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"

	caigotypes "github.com/dontpanicdao/caigo/types"
	"github.com/smartcontractkit/libocr/bigbigendian"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
)

const (
	OnchainConfigVersion = 1
	byteWidth            = 32
	length               = 3 * byteWidth
)

// report format
// 32 bytes - version
// 32 bytes - min
// 32 bytes - max

type OnchainConfigCodec struct{}

var _ median.OnchainConfigCodec = &OnchainConfigCodec{}

// DecodeToFelts decodes the onchainconfig into felt values (used in config digest hashing)
func (codec OnchainConfigCodec) DecodeToFelts(b []byte) ([]*big.Int, error) {
	if len(b) != length {
		return []*big.Int{}, fmt.Errorf("unexpected length of OnchainConfig, expected %v, got %v", length, len(b))
	}

	configVersion, err := bigbigendian.DeserializeSigned(byteWidth, b[:32])
	if err != nil {
		return []*big.Int{}, fmt.Errorf("unable to decode version: %s", err)
	}
	if OnchainConfigVersion != configVersion.Int64() {
		return []*big.Int{}, fmt.Errorf("unexpected version of OnchainConfig, expected %v, got %v", OnchainConfigVersion, configVersion.Int64())
	}

	min, err := bigbigendian.DeserializeSigned(byteWidth, b[byteWidth:2*byteWidth])
	if err != nil {
		return []*big.Int{}, err
	}
	max, err := bigbigendian.DeserializeSigned(byteWidth, b[2*byteWidth:])
	if err != nil {
		return []*big.Int{}, err
	}

	// ensure felts (used in config digester)
	min = starknet.SignedBigToFelt(min)
	max = starknet.SignedBigToFelt(max)

	return []*big.Int{configVersion, min, max}, nil
}

// Decode converts the onchainconfig via the outputs of DecodeToFelts into signed big.Ints that libocr expects
func (codec OnchainConfigCodec) Decode(b []byte) (median.OnchainConfig, error) {
	felts, err := codec.DecodeToFelts(b)
	if err != nil {
		return median.OnchainConfig{}, err
	}

	// convert felts to big.Ints
	min := starknet.FeltToSignedBig(&caigotypes.Felt{Int: felts[1]})
	max := starknet.FeltToSignedBig(&caigotypes.Felt{Int: felts[2]})

	if !(min.Cmp(max) <= 0) {
		return median.OnchainConfig{}, fmt.Errorf("OnchainConfig min (%v) should not be greater than max(%v)", min, max)
	}

	return median.OnchainConfig{Min: min, Max: max}, nil
}

// TODO: both 'EncodeFromBigInt' and 'EncodeFromFelt' have the same signature - we need a custom type to represent Felts
// EncodeFromBigInt encodes the config where min & max are big Ints with positive or negative values
func (codec OnchainConfigCodec) EncodeFromBigInt(version, min, max *big.Int) ([]byte, error) {
	return codec.EncodeFromFelt(version, starknet.SignedBigToFelt(min), starknet.SignedBigToFelt(max))
}

// EncodeFromFelt encodes the config where min & max are big.Int representations of a felt
// Cairo has no notion of signed values: negative values have to be wrapped into the upper half of PRIME (so 0 to PRIME/2 is positive, PRIME/2 to PRIME is negative)
func (codec OnchainConfigCodec) EncodeFromFelt(version, min, max *big.Int) ([]byte, error) {
	if version.Uint64() != OnchainConfigVersion {
		return nil, fmt.Errorf("unexpected version of OnchainConfig, expected %v, got %v", OnchainConfigVersion, version.Int64())
	}

	versionBytes, err := bigbigendian.SerializeSigned(byteWidth, version)
	if err != nil {
		return nil, err
	}

	minBytes, err := bigbigendian.SerializeSigned(byteWidth, min)
	if err != nil {
		return nil, err
	}

	maxBytes, err := bigbigendian.SerializeSigned(byteWidth, max)
	if err != nil {
		return nil, err
	}
	result := []byte{}
	result = append(result, versionBytes...)
	result = append(result, minBytes...)
	result = append(result, maxBytes...)

	return result, nil
}

// Encode takes the interface that libocr uses (+/- big.Ints) and serializes it into 3 felts
func (codec OnchainConfigCodec) Encode(c median.OnchainConfig) ([]byte, error) {
	return codec.EncodeFromBigInt(big.NewInt(OnchainConfigVersion), c.Min, c.Max)
}

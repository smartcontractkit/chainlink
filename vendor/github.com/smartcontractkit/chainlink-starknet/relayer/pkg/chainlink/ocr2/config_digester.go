package ocr2

import (
	"encoding/binary"
	"errors"
	"math/big"
	"strings"

	"github.com/NethermindEth/starknet.go/curve"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/ocr2/medianreport"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"
)

// TODO: use libocr constant
const ConfigDigestPrefixStarknet types.ConfigDigestPrefix = 4

var _ types.OffchainConfigDigester = (*offchainConfigDigester)(nil)

type offchainConfigDigester struct {
	chainID  string
	contract string
}

func NewOffchainConfigDigester(chainID, contract string) offchainConfigDigester {
	return offchainConfigDigester{
		chainID:  chainID,
		contract: contract,
	}
}

// TODO: ConfigDigest is byte[32] but what we really want here is a felt
func (d offchainConfigDigester) ConfigDigest(cfg types.ContractConfig) (types.ConfigDigest, error) {
	configDigest := types.ConfigDigest{}

	contractAddress, valid := new(big.Int).SetString(strings.TrimPrefix(d.contract, "0x"), 16)
	if !valid {
		return configDigest, errors.New("invalid contract address")
	}

	if len(d.chainID) > 31 {
		return configDigest, errors.New("chainID exceeds max length")
	}

	if len(cfg.Signers) != len(cfg.Transmitters) {
		return configDigest, errors.New("must have equal number of signers and transmitters")
	}

	if len(cfg.Signers) <= 3*int(cfg.F) {
		return configDigest, errors.New("number of oracles must be greater than 3*f")
	}

	oracles := []*big.Int{}
	for i := range cfg.Signers {
		signer := new(big.Int).SetBytes(cfg.Signers[i])
		transmitter, valid := new(big.Int).SetString(strings.TrimPrefix(string(cfg.Transmitters[i]), "0x"), 16)
		if !valid {
			return configDigest, errors.New("invalid transmitter")
		}
		oracles = append(oracles, signer, transmitter)
	}

	offchainConfig := starknet.EncodeFelts(cfg.OffchainConfig)

	onchainConfig, err := medianreport.OnchainConfigCodec{}.DecodeToFelts(cfg.OnchainConfig)
	if err != nil {
		return configDigest, err
	}

	// golang... https://stackoverflow.com/questions/28625546/mixing-exploded-slices-and-regular-parameters-in-variadic-functions
	msg := []*big.Int{
		new(big.Int).SetBytes([]byte(d.chainID)),       // chain_id
		contractAddress,                                // contract_address
		new(big.Int).SetUint64(cfg.ConfigCount),        // config_count
		new(big.Int).SetInt64(int64(len(cfg.Signers))), // oracles_len
	}
	msg = append(msg, oracles...)
	msg = append(
		msg,
		big.NewInt(int64(cfg.F)),              // f
		big.NewInt(int64(len(onchainConfig))), // onchain_config_len
	)
	msg = append(msg, onchainConfig...)
	msg = append(
		msg,
		new(big.Int).SetUint64(cfg.OffchainConfigVersion), // offchain_config_version
		big.NewInt(int64(len(offchainConfig))),            // offchain_config_len
	)
	msg = append(msg, offchainConfig...) // offchain_config

	digest, err := curve.Curve.ComputeHashOnElements(msg)
	if err != nil {
		return configDigest, err
	}
	digest.FillBytes(configDigest[:])

	// set first two bytes to the digest prefix
	pre, err := d.ConfigDigestPrefix()
	if err != nil {
		return configDigest, err
	}
	binary.BigEndian.PutUint16(configDigest[:2], uint16(pre))

	return configDigest, nil
}

func (offchainConfigDigester) ConfigDigestPrefix() (types.ConfigDigestPrefix, error) {
	return ConfigDigestPrefixStarknet, nil
}

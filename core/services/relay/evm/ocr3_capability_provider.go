package evm

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/ethereum/go-ethereum/crypto"

	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

type ocr3CapabilityProvider struct {
	types.PluginProvider
	transmitter ocr3types.ContractTransmitter[[]byte]
}

func (o *ocr3CapabilityProvider) OCR3ContractTransmitter() ocr3types.ContractTransmitter[[]byte] {
	return o.transmitter
}

var _ LogDecoder = &ocr3CapabilityLogDecoder{}

type ocr3CapabilityLogDecoder struct {
	eventName string
	eventSig  common.Hash
	abi       *abi.ABI
}

// Modified newOCR2AggregatorLogDecoder to use OCR3Capability ABI
func newOCR3CapabilityLogDecoder() (*ocr3CapabilityLogDecoder, error) {
	const eventName = "ConfigSet"
	abi, err := ocr3_capability.OCR3CapabilityMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return &ocr3CapabilityLogDecoder{
		eventName: eventName,
		eventSig:  abi.Events[eventName].ID,
		abi:       abi,
	}, nil
}

func (d *ocr3CapabilityLogDecoder) Decode(rawLog []byte) (ocrtypes.ContractConfig, error) {
	unpacked := new(ocr3_capability.OCR3CapabilityConfigSet)
	err := d.abi.UnpackIntoInterface(unpacked, d.eventName, rawLog)
	if err != nil {
		return ocrtypes.ContractConfig{}, fmt.Errorf("failed to unpack log data: %w", err)
	}

	var transmitAccounts []ocrtypes.Account
	for _, addr := range unpacked.Transmitters {
		transmitAccounts = append(transmitAccounts, ocrtypes.Account(addr.Hex()))
	}
	var signers []ocrtypes.OnchainPublicKey
	allPubKeys := map[string]any{}
	for _, pubKey := range unpacked.Signers {
		pubKey := pubKey

		// validate uniqueness of each individual key
		pubKeys, err := ocrcommon.UnmarshalMultichainPublicKey(pubKey)
		if err != nil {
			return ocrtypes.ContractConfig{}, err
		}
		for _, key := range pubKeys {
			raw := hex.EncodeToString(key)
			_, exists := allPubKeys[raw]
			if exists {
				return ocrtypes.ContractConfig{}, fmt.Errorf("Duplicate onchain public key: %v", raw)
			}
			allPubKeys[raw] = struct{}{}
		}

		signers = append(signers, pubKey)
	}

	return ocrtypes.ContractConfig{
		ConfigDigest:          unpacked.ConfigDigest,
		ConfigCount:           unpacked.ConfigCount,
		Signers:               signers,
		Transmitters:          transmitAccounts,
		F:                     unpacked.F,
		OnchainConfig:         unpacked.OnchainConfig,
		OffchainConfigVersion: unpacked.OffchainConfigVersion,
		OffchainConfig:        unpacked.OffchainConfig,
	}, nil
}

func (d *ocr3CapabilityLogDecoder) EventSig() common.Hash {
	return d.eventSig
}

var _ ocrtypes.OffchainConfigDigester = OCR3CapabilityOffchainConfigDigester{}

// EVMOffchainConfigDigester forked to not assume signers are 20 byte addresses
type OCR3CapabilityOffchainConfigDigester struct {
	ChainID         uint64
	ContractAddress common.Address
}

func (d OCR3CapabilityOffchainConfigDigester) ConfigDigest(ctx context.Context, cc ocrtypes.ContractConfig) (ocrtypes.ConfigDigest, error) {
	signers := [][]byte{}
	for _, signer := range cc.Signers {
		signers = append(signers, signer)
	}
	transmitters := []common.Address{}
	for i, transmitter := range cc.Transmitters {
		if !strings.HasPrefix(string(transmitter), "0x") || len(transmitter) != 42 || !common.IsHexAddress(string(transmitter)) {
			return ocrtypes.ConfigDigest{}, fmt.Errorf("%v-th evm transmitter should be a 42 character Ethereum address string, but got '%v'", i, transmitter)
		}
		a := common.HexToAddress(string(transmitter))
		transmitters = append(transmitters, a)
	}

	return ocr3CapabilityConfigDigest(
		d.ChainID,
		d.ContractAddress,
		cc.ConfigCount,
		signers,
		transmitters,
		cc.F,
		cc.OnchainConfig,
		cc.OffchainConfigVersion,
		cc.OffchainConfig,
	), nil
}

const ConfigDigestPrefixKeystoneOCR3Capability ocrtypes.ConfigDigestPrefix = 0x000e

func (d OCR3CapabilityOffchainConfigDigester) ConfigDigestPrefix(ctx context.Context) (ocrtypes.ConfigDigestPrefix, error) {
	return ConfigDigestPrefixKeystoneOCR3Capability, nil
}

func makeOCR3CapabilityConfigDigestArgs() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{Name: "chainId", Type: mustNewType("uint256")},
		{Name: "contractAddress", Type: mustNewType("address")},
		{Name: "configCount", Type: mustNewType("uint64")},
		{Name: "signers", Type: mustNewType("bytes[]")},
		{Name: "transmitters", Type: mustNewType("address[]")},
		{Name: "f", Type: mustNewType("uint8")},
		{Name: "onchainConfig", Type: mustNewType("bytes")},
		{Name: "encodedConfigVersion", Type: mustNewType("uint64")},
		{Name: "encodedConfig", Type: mustNewType("bytes")},
	})
}

var ocr3CapabilityConfigDigestArgs = makeOCR3CapabilityConfigDigestArgs()

func ocr3CapabilityConfigDigest(
	chainID uint64,
	contractAddress common.Address,
	configCount uint64,
	oracles [][]byte,
	transmitters []common.Address,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) ocrtypes.ConfigDigest {
	chainIDBig := new(big.Int)
	chainIDBig.SetUint64(chainID)
	msg, err := ocr3CapabilityConfigDigestArgs.Pack(
		chainIDBig,
		contractAddress,
		configCount,
		oracles,
		transmitters,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)
	if err != nil {
		// assertion
		panic(err)
	}
	rawHash := crypto.Keccak256(msg)
	configDigest := ocrtypes.ConfigDigest{}
	if n := copy(configDigest[:], rawHash); n != len(configDigest) {
		// assertion
		panic("copy too little data")
	}
	if ConfigDigestPrefixKeystoneOCR3Capability != 0x000e {
		// assertion
		panic("wrong ConfigDigestPrefix")
	}
	configDigest[0] = 0
	configDigest[1] = 0x0e
	return configDigest
}

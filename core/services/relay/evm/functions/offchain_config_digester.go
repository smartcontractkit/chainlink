package functions

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/libocr/gethwrappers2/exposedocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var (
	_                     types.OffchainConfigDigester = FunctionsOffchainConfigDigester{}
	configDigestArgs                                   = makeConfigDigestArgs()
	FunctionsDigestPrefix                              = types.ConfigDigestPrefixEVM
	ThresholdDigestPrefix                              = types.ConfigDigestPrefix(7)
	S4DigestPrefix                                     = types.ConfigDigestPrefix(8)
)

type FunctionsOffchainConfigDigester struct {
	ChainID         uint64
	ContractAddress common.Address
	PluginType      FunctionsPluginType
}

func (d FunctionsOffchainConfigDigester) ConfigDigest(cc types.ContractConfig) (types.ConfigDigest, error) {
	signers := []common.Address{}
	for i, signer := range cc.Signers {
		if len(signer) != 20 {
			return types.ConfigDigest{}, fmt.Errorf("%v-th evm signer should be a 20 byte address, but got %x", i, signer)
		}
		a := common.BytesToAddress(signer)
		signers = append(signers, a)
	}
	transmitters := []common.Address{}
	for i, transmitter := range cc.Transmitters {
		if !strings.HasPrefix(string(transmitter), "0x") || len(transmitter) != 42 || !common.IsHexAddress(string(transmitter)) {
			return types.ConfigDigest{}, fmt.Errorf("%v-th evm transmitter should be a 42 character Ethereum address string, but got '%v'", i, transmitter)
		}
		a := common.HexToAddress(string(transmitter))
		transmitters = append(transmitters, a)
	}

	return configDigest(
		d.PluginType,
		d.ChainID,
		d.ContractAddress,
		cc.ConfigCount,
		signers,
		transmitters,
		cc.F,
		cc.OnchainConfig,
		cc.OffchainConfigVersion,
		cc.OffchainConfig,
	)
}

func (d FunctionsOffchainConfigDigester) ConfigDigestPrefix() (types.ConfigDigestPrefix, error) {
	switch d.PluginType {
	case FunctionsPlugin:
		return FunctionsDigestPrefix, nil
	case ThresholdPlugin:
		return ThresholdDigestPrefix, nil
	case S4Plugin:
		return S4DigestPrefix, nil
	default:
		return 0, fmt.Errorf("unknown plugin type: %v", d.PluginType)
	}
}

func configDigest(
	pluginType FunctionsPluginType,
	chainID uint64,
	contractAddress common.Address,
	configCount uint64,
	oracles []common.Address,
	transmitters []common.Address,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) (configDigest types.ConfigDigest, err error) {
	chainIDBig := new(big.Int)
	chainIDBig.SetUint64(chainID)
	fmt.Printf(
		"configDigest()\nchainIDBig: %v\ncontractAddress %v\nconfigCount %v\noracles %v\ntransmitters %v\nf %v\nonchainConfig %v\noffchainConfigVersion %v\noffchainConfig %v\n",
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
	msg, err := configDigestArgs.Pack(
		chainIDBig,
		contractAddress,
		uint64(4),
		oracles,
		transmitters,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)
	if err != nil {
		return configDigest, err
	}

	rawHash := crypto.Keccak256(msg)
	if n := copy(configDigest[:], rawHash); n != len(configDigest) {
		return configDigest, errors.New("copy too little data")
	}

	var prefix types.ConfigDigestPrefix
	switch pluginType {
	case FunctionsPlugin:
		prefix = FunctionsDigestPrefix
	case ThresholdPlugin:
		prefix = ThresholdDigestPrefix
	case S4Plugin:
		prefix = S4DigestPrefix
	default:
		return configDigest, errors.New("unknown plugin type")
	}

	fmt.Printf("configDigest() config digest: %v", rawHash)
	binary.BigEndian.PutUint16(configDigest[:2], uint16(prefix))
	return configDigest, nil
}

func makeConfigDigestArgs() abi.Arguments {
	abi, err := abi.JSON(strings.NewReader(
		exposedocr2aggregator.ExposedOCR2AggregatorABI))
	if err != nil {
		// assertion
		panic(fmt.Sprintf("could not parse aggregator ABI: %s", err.Error()))
	}
	return abi.Methods["exposedConfigDigestFromConfigData"].Inputs
}

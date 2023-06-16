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

func makeConfigDigestArgs() abi.Arguments {
	abi, err := abi.JSON(strings.NewReader(
		exposedocr2aggregator.ExposedOCR2AggregatorABI))
	if err != nil {
		// assertion
		panic(fmt.Sprintf("could not parse aggregator ABI: %s", err.Error()))
	}
	return abi.Methods["exposedConfigDigestFromConfigData"].Inputs
}

var (
	configDigestArgs      = makeConfigDigestArgs()
	functionsDigestPrefix = types.ConfigDigestPrefixEVM
	thresholdDigestPrefix = types.ConfigDigestPrefix(7)
	s4DigestPrefix        = types.ConfigDigestPrefix(8)
)

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
	msg, err := configDigestArgs.Pack(
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
		return configDigest, err
	}
	rawHash := crypto.Keccak256(msg)
	if n := copy(configDigest[:], rawHash); n != len(configDigest) {
		return configDigest, errors.New("copy too little data")
	}

	var prefix types.ConfigDigestPrefix
	switch pluginType {
	case FunctionsPlugin:
		prefix = functionsDigestPrefix
	case ThresholdPlugin:
		prefix = thresholdDigestPrefix
	case S4Plugin:
		prefix = s4DigestPrefix
	default:
		return configDigest, errors.New("unknown plugin type")
	}

	binary.BigEndian.PutUint16(configDigest[:2], uint16(prefix))
	return configDigest, nil
}

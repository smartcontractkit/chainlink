package config

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink/libocr/gethwrappers/exposedoffchainaggregator"

	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

func makeConfigDigestArgs() abi.Arguments {
	abi, err := abi.JSON(strings.NewReader(
		exposedoffchainaggregator.ExposedOffchainAggregatorABI))
	if err != nil {
		panic(fmt.Sprintf("could not parse aggregator ABI: %s", err.Error()))
	}
	return abi.Methods["exposedConfigDigestFromConfigData"].Inputs
}

var configDigestArgs = makeConfigDigestArgs()

func ConfigDigest(
	contractAddress common.Address,
	configCount uint64,
	oracles []common.Address,
	transmitters []common.Address,
	threshold uint8,
	encodedConfigVersion uint64,
	config []byte,
) types.ConfigDigest {
	msg, err := configDigestArgs.Pack(
		contractAddress,
		configCount,
		oracles,
		transmitters,
		threshold,
		encodedConfigVersion,
		config,
	)
	if err != nil {
		panic(err)
	}
	rawHash := crypto.Keccak256(msg)
	return types.BytesToConfigDigest(rawHash)
}

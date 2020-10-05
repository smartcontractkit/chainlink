package confighelper

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

func makeConfigDigestArgs() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{Name: "contractAddress", Type: mustNewType("address")},
		{Name: "configCount", Type: mustNewType("uint64")},
		{Name: "signers", Type: mustNewType("address[]")},
		{Name: "transmitters", Type: mustNewType("address[]")},
		{Name: "threshold", Type: mustNewType("uint8")},
		{Name: "encodedConfigVersion", Type: mustNewType("uint64")},
		{Name: "encodedConfig", Type: mustNewType("bytes")},
	})
}

var configDigestArgs = makeConfigDigestArgs()

func configDigest(
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

func ContractConfigFromConfigSetEvent(changed offchainaggregator.OffchainAggregatorConfigSet) types.ContractConfig {
	return types.ContractConfig{
		configDigest(
			changed.Raw.Address,
			changed.ConfigCount,
			changed.Signers,
			changed.Transmitters,
			changed.Threshold,
			changed.EncodedConfigVersion,
			changed.Encoded,
		),
		changed.Signers,
		changed.Transmitters,
		changed.Threshold,
		changed.EncodedConfigVersion,
		changed.Encoded,
	}
}

package mercury

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ava-labs/coreth/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/wsrpc/credentials"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mercury_exposed_verifier"
)

func makeConfigDigestArgs() abi.Arguments {
	abi, err := abi.JSON(strings.NewReader(mercury_exposed_verifier.MercuryExposedVerifierABI))
	if err != nil {
		// assertion
		panic(fmt.Sprintf("could not parse aggregator ABI: %s", err.Error()))
	}
	return abi.Methods["exposedConfigDigestFromConfigData"].Inputs
}

var configDigestArgs = makeConfigDigestArgs()

func configDigest(
	feedID common.Hash,
	chainID uint64,
	contractAddress common.Address,
	configCount uint64,
	oracles []common.Address,
	transmitters []credentials.StaticSizedPublicKey,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) types.ConfigDigest {
	chainIDBig := new(big.Int)
	chainIDBig.SetUint64(chainID)
	msg, err := configDigestArgs.Pack(
		feedID,
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
	configDigest := types.ConfigDigest{}
	if n := copy(configDigest[:], rawHash); n != len(configDigest) {
		// assertion
		panic("copy too little data")
	}
	if types.ConfigDigestPrefixEVM != 1 {
		// assertion
		panic("wrong ConfigDigestPrefix")
	}
	configDigest[0] = 0
	configDigest[1] = 1
	return configDigest
}

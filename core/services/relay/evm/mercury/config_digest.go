package mercury

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
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
	binary.BigEndian.PutUint16(configDigest[:2], uint16(types.ConfigDigestPrefixMercuryV02))
	if !(configDigest[0] == 0 || configDigest[1] == 6) {
		// assertion
		panic("unexpected mismatch")
	}
	return configDigest
}

package mercury

import (
	"math/big"
	"reflect"
	"testing"
	"unsafe"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/wsrpc/credentials"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/exposed_verifier"
)

// Adapted from: https://github.com/smartcontractkit/offchain-reporting/blob/991ebe1462fd56826a1ddfb34287d542acb2baee/lib/offchainreporting2/chains/evmutil/config_digest_test.go

func TestConfigCalculationMatches(t *testing.T) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "could not make private key for EOA owner")
	owner, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)
	backend := backends.NewSimulatedBackend(
		core.GenesisAlloc{owner.From: {Balance: new(big.Int).Lsh(big.NewInt(1), 60)}},
		ethconfig.Defaults.Miner.GasCeil,
	)
	_, _, eoa, err := exposed_verifier.DeployExposedVerifier(
		owner, backend,
	)
	backend.Commit()
	require.NoError(t, err, "could not deploy test EOA")
	p := gopter.NewProperties(nil)
	p.Property("onchain/offchain config digests match", prop.ForAll(
		func(
			feedID [32]byte,
			chainID uint64,
			contractAddress common.Address,
			configCount uint64,
			oracles []common.Address,
			transmitters [][32]byte,
			f uint8,
			onchainConfig []byte,
			offchainConfigVersion uint64,
			offchainConfig []byte,
		) bool {
			chainIDBig := new(big.Int).SetUint64(chainID)
			golangDigest := configDigest(
				feedID,
				chainIDBig,
				contractAddress,
				configCount,
				oracles,
				*(*[]credentials.StaticSizedPublicKey)(unsafe.Pointer(&transmitters)),
				f,
				onchainConfig,
				offchainConfigVersion,
				offchainConfig,
				ocrtypes.ConfigDigestPrefixMercuryV02,
			)

			bigChainID := new(big.Int)
			bigChainID.SetUint64(chainID)

			solidityDigest, err := eoa.ExposedConfigDigestFromConfigData(nil,
				feedID,
				bigChainID,
				contractAddress,
				configCount,
				oracles,
				transmitters,
				f,
				onchainConfig,
				offchainConfigVersion,
				offchainConfig,
			)
			require.NoError(t, err, "could not compute solidity version of config digest")
			return golangDigest == solidityDigest
		},
		GenHash(t),
		gen.UInt64(),
		GenAddress(t),
		gen.UInt64(),
		GenAddressArray(t),
		GenClientPubKeyArray(t),
		gen.UInt8(),
		GenBytes(t),
		gen.UInt64(),
		GenBytes(t),
	))
	p.TestingRun(t)
}

func GenHash(t *testing.T) gopter.Gen {
	var byteGens []gopter.Gen
	for i := 0; i < 32; i++ {
		byteGens = append(byteGens, gen.UInt8())
	}
	return gopter.CombineGens(byteGens...).Map(
		func(byteArray interface{}) (rv common.Hash) {
			array, ok := byteArray.(*gopter.GenResult).Retrieve()
			require.True(t, ok, "failed to retrieve gen result")
			for i, byteVal := range array.([]interface{}) {
				rv[i] = byteVal.(uint8)
			}
			return rv
		},
	)
}

func GenHashArray(t *testing.T) gopter.Gen {
	return gen.UInt8Range(0, 31).FlatMap(
		func(length interface{}) gopter.Gen {
			var hashGens []gopter.Gen
			for i := uint8(0); i < length.(uint8); i++ {
				hashGens = append(hashGens, GenHash(t))
			}
			return gopter.CombineGens(hashGens...).Map(
				func(hashArray interface{}) (rv []common.Hash) {
					array, ok := hashArray.(*gopter.GenResult).Retrieve()
					require.True(t, ok, "could not extract hash array")
					for _, hashVal := range array.([]interface{}) {
						rv = append(rv, hashVal.(common.Hash))
					}
					return rv
				},
			)
		},
		reflect.ValueOf([]common.Hash{}).Type(),
	)
}

func GenAddress(t *testing.T) gopter.Gen {
	return GenHash(t).Map(
		func(hash interface{}) common.Address {
			iHash, ok := hash.(*gopter.GenResult).Retrieve()
			require.True(t, ok, "failed to retrieve hash")
			return common.BytesToAddress(iHash.(common.Hash).Bytes())
		},
	)
}

func GenAddressArray(t *testing.T) gopter.Gen {
	return GenHashArray(t).Map(
		func(hashes interface{}) (rv []common.Address) {
			hashArray, ok := hashes.(*gopter.GenResult).Retrieve()
			require.True(t, ok, "failed to retrieve hashes")
			for _, hash := range hashArray.([]common.Hash) {
				rv = append(rv, common.BytesToAddress(hash.Bytes()))
			}
			return rv
		},
	)
}

func GenClientPubKey(t *testing.T) gopter.Gen {
	return GenHash(t).Map(
		func(hash interface{}) (pk [32]byte) {
			iHash, ok := hash.(*gopter.GenResult).Retrieve()
			require.True(t, ok, "failed to retrieve hash")
			copy(pk[:], (iHash.(common.Hash).Bytes()))
			return
		},
	)
}

func GenClientPubKeyArray(t *testing.T) gopter.Gen {
	return GenHashArray(t).Map(
		func(hashes interface{}) (rv [][32]byte) {
			hashArray, ok := hashes.(*gopter.GenResult).Retrieve()
			require.True(t, ok, "failed to retrieve hashes")
			for _, hash := range hashArray.([]common.Hash) {
				pk := [32]byte{}
				copy(pk[:], hash.Bytes())
				rv = append(rv, pk)
			}
			return rv
		},
	)
}

func GenBytes(t *testing.T) gopter.Gen {
	return gen.UInt16Range(0, 2000).FlatMap(
		func(length interface{}) gopter.Gen {
			var byteGens []gopter.Gen
			for i := uint16(0); i < length.(uint16); i++ {
				byteGens = append(byteGens, gen.UInt8())
			}
			return gopter.CombineGens(byteGens...).Map(
				func(byteArray interface{}) []byte {
					array, ok := byteArray.(*gopter.GenResult).Retrieve()
					require.True(t, ok, "failed to retrieve gen result")
					iArray := array.([]interface{})
					rv := make([]byte, len(iArray))
					for i, byteVal := range iArray {
						rv[i] = byteVal.(uint8)
					}
					return rv
				},
			)
		},
		reflect.ValueOf([]byte{}).Type(),
	)
}

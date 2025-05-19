package llo

import (
	"context"
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/wsrpc/credentials"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/exposed_configurator"
)

// Originally sourced from: https://github.com/smartcontractkit/offchain-reporting/blob/991ebe1462fd56826a1ddfb34287d542acb2baee/lib/offchainreporting2/chains/evmutil/offchain_config_digester.go

var _ ocrtypes.OffchainConfigDigester = OffchainConfigDigester{}

func NewOffchainConfigDigester(configID common.Hash, chainID *big.Int, contractAddress common.Address, prefix ocrtypes.ConfigDigestPrefix) OffchainConfigDigester {
	return OffchainConfigDigester{configID, chainID, contractAddress, prefix}
}

type OffchainConfigDigester struct {
	ConfigID        common.Hash
	ChainID         *big.Int
	ContractAddress common.Address
	Prefix          ocrtypes.ConfigDigestPrefix
}

func (d OffchainConfigDigester) ConfigDigest(_ context.Context, cc ocrtypes.ContractConfig) (ocrtypes.ConfigDigest, error) {
	onchainPubKeys := make([][]byte, len(cc.Signers))
	for i, signer := range cc.Signers {
		// Onchainpubkeys can be anything
		// TODO: Implement and enforce MultiChainKeyBundle format?
		// MERC-3594
		onchainPubKeys[i] = signer
	}
	transmitters := []credentials.StaticSizedPublicKey{}
	for i, transmitter := range cc.Transmitters {
		if len(transmitter) != 2*ed25519.PublicKeySize {
			return ocrtypes.ConfigDigest{}, errors.Errorf("%v-th evm transmitter should be a 64 character hex-encoded ed25519 public key, but got '%v' (%d chars)", i, transmitter, len(transmitter))
		}
		var t credentials.StaticSizedPublicKey
		b, err := hex.DecodeString(string(transmitter))
		if err != nil {
			return ocrtypes.ConfigDigest{}, errors.Wrapf(err, "%v-th evm transmitter is not valid hex, got: %q", i, transmitter)
		}
		copy(t[:], b)

		transmitters = append(transmitters, t)
	}

	return configDigest(
		d.ConfigID,
		d.ChainID,
		d.ContractAddress,
		cc.ConfigCount,
		onchainPubKeys,
		transmitters,
		cc.F,
		cc.OnchainConfig,
		cc.OffchainConfigVersion,
		cc.OffchainConfig,
		d.Prefix,
	), nil
}

func (d OffchainConfigDigester) ConfigDigestPrefix(context.Context) (ocrtypes.ConfigDigestPrefix, error) {
	return d.Prefix, nil
}

func makeConfigDigestArgs() abi.Arguments {
	abi, err := abi.JSON(strings.NewReader(exposed_configurator.ExposedConfiguratorABI))
	if err != nil {
		// assertion
		panic(fmt.Sprintf("could not parse configurator ABI: %s", err.Error()))
	}
	return abi.Methods["exposedConfigDigestFromConfigData"].Inputs
}

var configDigestArgs = makeConfigDigestArgs()

func configDigest(
	feedID common.Hash,
	chainID *big.Int,
	contractAddress common.Address,
	configCount uint64,
	onchainPubKeys [][]byte,
	transmitters []credentials.StaticSizedPublicKey,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	prefix types.ConfigDigestPrefix,
) types.ConfigDigest {
	msg, err := configDigestArgs.Pack(
		feedID,
		chainID,
		contractAddress,
		configCount,
		onchainPubKeys,
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
	binary.BigEndian.PutUint16(configDigest[:2], uint16(prefix))
	return configDigest
}

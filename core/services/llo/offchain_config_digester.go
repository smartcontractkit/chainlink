package llo

import (
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

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/wsrpc/credentials"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/exposed_channel_verifier"
)

// Originally sourced from: https://github.com/smartcontractkit/offchain-reporting/blob/991ebe1462fd56826a1ddfb34287d542acb2baee/lib/offchainreporting2/chains/evmutil/offchain_config_digester.go

var _ ocrtypes.OffchainConfigDigester = OffchainConfigDigester{}

func NewOffchainConfigDigester(chainID *big.Int, contractAddress common.Address) OffchainConfigDigester {
	return OffchainConfigDigester{chainID, contractAddress}
}

type OffchainConfigDigester struct {
	ChainID         *big.Int
	ContractAddress common.Address
}

func (d OffchainConfigDigester) ConfigDigest(cc ocrtypes.ContractConfig) (ocrtypes.ConfigDigest, error) {
	signers := []common.Address{}
	for i, signer := range cc.Signers {
		if len(signer) != 20 {
			return ocrtypes.ConfigDigest{}, errors.Errorf("%v-th evm signer should be a 20 byte address, but got %x", i, signer)
		}
		a := common.BytesToAddress(signer)
		signers = append(signers, a)
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

func (d OffchainConfigDigester) ConfigDigestPrefix() (ocrtypes.ConfigDigestPrefix, error) {
	return ocrtypes.ConfigDigestPrefixLLO, nil
}

func makeConfigDigestArgs() abi.Arguments {
	abi, err := abi.JSON(strings.NewReader(exposed_channel_verifier.ExposedChannelVerifierABI))
	if err != nil {
		// assertion
		panic(fmt.Sprintf("could not parse aggregator ABI: %s", err.Error()))
	}
	return abi.Methods["exposedConfigDigestFromConfigData"].Inputs
}

var configDigestArgs = makeConfigDigestArgs()

func configDigest(
	chainID *big.Int,
	contractAddress common.Address,
	configCount uint64,
	oracles []common.Address,
	transmitters []credentials.StaticSizedPublicKey,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) (types.ConfigDigest, error) {
	msg, err := configDigestArgs.Pack(
		chainID,
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
		return types.ConfigDigest{}, fmt.Errorf("could not pack config digest args: %v", err)
	}
	rawHash := crypto.Keccak256(msg)
	configDigest := types.ConfigDigest{}
	if n := copy(configDigest[:], rawHash); n != len(configDigest) {
		return types.ConfigDigest{}, fmt.Errorf("copied too little data: %d/%d", n, len(configDigest))
	}
	binary.BigEndian.PutUint16(configDigest[:2], uint16(ocrtypes.ConfigDigestPrefixLLO))
	if !(configDigest[0] == 0 && configDigest[1] == 9) {
		return types.ConfigDigest{}, fmt.Errorf("wrong ConfigDigestPrefix; got: %x, expected: %d", configDigest[:2], ocrtypes.ConfigDigestPrefixLLO)
	}
	return configDigest, nil
}

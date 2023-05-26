package mercury

import (
	"crypto/ed25519"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/wsrpc/credentials"
)

// Originally sourced from: https://github.com/smartcontractkit/offchain-reporting/blob/991ebe1462fd56826a1ddfb34287d542acb2baee/lib/offchainreporting2/chains/evmutil/offchain_config_digester.go

var _ ocrtypes.OffchainConfigDigester = OffchainConfigDigester{}

func NewOffchainConfigDigester(feedID [32]byte, chainID uint64, contractAddress common.Address) OffchainConfigDigester {
	return OffchainConfigDigester{feedID, chainID, contractAddress}
}

type OffchainConfigDigester struct {
	FeedID          [32]byte
	ChainID         uint64
	ContractAddress common.Address
}

func (d OffchainConfigDigester) ConfigDigest(cc types.ContractConfig) (types.ConfigDigest, error) {
	signers := []common.Address{}
	for i, signer := range cc.Signers {
		if len(signer) != 20 {
			return types.ConfigDigest{}, errors.Errorf("%v-th evm signer should be a 20 byte address, but got %x", i, signer)
		}
		a := common.BytesToAddress(signer)
		signers = append(signers, a)
	}
	transmitters := []credentials.StaticSizedPublicKey{}
	for i, transmitter := range cc.Transmitters {
		if len(transmitter) != 2*ed25519.PublicKeySize {
			return types.ConfigDigest{}, errors.Errorf("%v-th evm transmitter should be a 64 character hex-encoded ed25519 public key, but got '%v' (%d chars)", i, transmitter, len(transmitter))
		}
		var t credentials.StaticSizedPublicKey
		b, err := hex.DecodeString(string(transmitter))
		if err != nil {
			return types.ConfigDigest{}, errors.Wrapf(err, "%v-th evm transmitter is not valid hex, got: %q", i, transmitter)
		}
		copy(t[:], b)

		transmitters = append(transmitters, t)
	}

	return configDigest(
		d.FeedID,
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

func (d OffchainConfigDigester) ConfigDigestPrefix() types.ConfigDigestPrefix {
	return types.ConfigDigestPrefixEVM
}

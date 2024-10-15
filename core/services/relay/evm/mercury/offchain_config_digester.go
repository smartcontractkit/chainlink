package mercury

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/wsrpc/credentials"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
)

// Originally sourced from: https://github.com/smartcontractkit/offchain-reporting/blob/991ebe1462fd56826a1ddfb34287d542acb2baee/lib/offchainreporting2/chains/evmutil/offchain_config_digester.go

var _ ocrtypes.OffchainConfigDigester = OffchainConfigDigester{}

func NewOffchainConfigDigester(feedID [32]byte, chainID *big.Int, contractAddress common.Address, prefix ocrtypes.ConfigDigestPrefix) OffchainConfigDigester {
	return OffchainConfigDigester{feedID, chainID, contractAddress, prefix}
}

type OffchainConfigDigester struct {
	FeedID          utils.FeedID
	ChainID         *big.Int
	ContractAddress common.Address
	Prefix          ocrtypes.ConfigDigestPrefix
}

func (d OffchainConfigDigester) ConfigDigest(ctx context.Context, cc ocrtypes.ContractConfig) (ocrtypes.ConfigDigest, error) {
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
		common.Hash(d.FeedID),
		d.ChainID,
		d.ContractAddress,
		cc.ConfigCount,
		signers,
		transmitters,
		cc.F,
		cc.OnchainConfig,
		cc.OffchainConfigVersion,
		cc.OffchainConfig,
		d.Prefix,
	), nil
}

func (d OffchainConfigDigester) ConfigDigestPrefix(ctx context.Context) (ocrtypes.ConfigDigestPrefix, error) {
	return d.Prefix, nil
}

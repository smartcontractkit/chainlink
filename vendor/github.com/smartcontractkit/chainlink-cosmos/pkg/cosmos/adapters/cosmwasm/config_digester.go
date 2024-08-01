package cosmwasm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	cosmosSDK "github.com/cosmos/cosmos-sdk/types"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/crypto/blake2s"
)

const ConfigDigestPrefixCosmos types.ConfigDigestPrefix = 2

var _ types.OffchainConfigDigester = (*OffchainConfigDigester)(nil)

type OffchainConfigDigester struct {
	chainID  string
	contract cosmosSDK.AccAddress
}

func NewOffchainConfigDigester(chainID string, contract cosmosSDK.AccAddress) OffchainConfigDigester {
	return OffchainConfigDigester{
		chainID:  chainID,
		contract: contract,
	}
}

func (cd OffchainConfigDigester) ConfigDigest(cfg types.ContractConfig) (types.ConfigDigest, error) {
	digest := types.ConfigDigest{}
	buf := bytes.NewBuffer([]byte{})

	if len(cd.chainID) > math.MaxUint8 {
		return digest, errors.New("chainID exceeds max uint8 length")
	}

	if err := binary.Write(buf, binary.BigEndian, uint8(len(cd.chainID))); err != nil {
		return digest, err
	}

	if _, err := buf.Write([]byte(cd.chainID)); err != nil {
		return digest, err
	}

	if _, err := buf.Write([]byte(cd.contract.String())); err != nil {
		return digest, err
	}

	if err := binary.Write(buf, binary.BigEndian, uint32(cfg.ConfigCount)); err != nil {
		return digest, err
	}

	if err := binary.Write(buf, binary.BigEndian, uint8(len(cfg.Signers))); err != nil {
		return digest, err
	}

	for _, signer := range cfg.Signers {
		if _, err := buf.Write(signer); err != nil {
			return digest, err
		}
	}

	// NOTE: We assume that signers and transmitters have the same length, currently
	// enforced onchain https://github.com/smartcontractkit/chainlink-cosmos/blob/9465a4ace6954b8647869d279363a25d1ae1b934/contracts/ocr2/src/contract.rs#L508
	// and that they have fixed sizes, thus we don't need a transmitter length prefix
	// here to avoid config digest collisions. Should that enforcement change
	// we'll need to add a length prefix here.
	for _, transmitter := range cfg.Transmitters {
		if _, err := buf.Write([]byte(transmitter)); err != nil {
			return digest, err
		}
	}

	if err := binary.Write(buf, binary.BigEndian, cfg.F); err != nil {
		return digest, err
	}

	if err := binary.Write(buf, binary.BigEndian, uint32(len(cfg.OnchainConfig))); err != nil {
		return digest, err
	}

	if _, err := buf.Write(cfg.OnchainConfig); err != nil {
		return digest, err
	}

	if err := binary.Write(buf, binary.BigEndian, cfg.OffchainConfigVersion); err != nil {
		return digest, err
	}

	if err := binary.Write(buf, binary.BigEndian, uint32(len(cfg.OffchainConfig))); err != nil {
		return digest, err
	}

	if _, err := buf.Write(cfg.OffchainConfig); err != nil {
		return digest, err
	}

	rawHash := blake2s.Sum256(buf.Bytes())
	if n := copy(digest[:], rawHash[:]); n != len(digest) {
		return digest, fmt.Errorf("incorrect hash size %d, expected %d", n, len(digest))
	}

	pre, err := cd.ConfigDigestPrefix()
	if err != nil {
		return digest, err
	}

	digest[0] = 0x00
	digest[1] = uint8(pre)

	return digest, nil
}

// This should return the same constant value on every invocation
func (OffchainConfigDigester) ConfigDigestPrefix() (types.ConfigDigestPrefix, error) {
	return ConfigDigestPrefixCosmos, nil
}

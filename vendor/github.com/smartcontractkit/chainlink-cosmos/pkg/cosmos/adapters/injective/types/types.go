package types

import (
	fmt "fmt"
	"strings"

	errors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"
	"golang.org/x/crypto/sha3"
)

const FeedIDMaxLength = 20

var digestPrefixCosmos = []byte("\x00\x02")
var digestSeparator = []byte("\x00\x00")

func (cfg *ContractConfig) Digest(chainID, feedID string) []byte {
	data, err := proto.Marshal(cfg)
	if err != nil {
		panic("unmarshable")
	}

	w := sha3.NewLegacyKeccak256()
	if _, err := w.Write(data); err != nil {
		panic(err)
	}
	if _, err := w.Write(digestSeparator); err != nil {
		panic(err)
	}
	if _, err := w.Write([]byte(chainID)); err != nil {
		panic(err)
	}
	if _, err := w.Write(digestSeparator); err != nil {
		panic(err)
	}
	if _, err := w.Write([]byte(feedID)); err != nil {
		panic(err)
	}

	configDigest := w.Sum(nil)
	configDigest[0] = digestPrefixCosmos[0]
	configDigest[1] = digestPrefixCosmos[1]

	return configDigest
}

func (cfg *FeedConfig) ValidTransmitters() map[string]struct{} {
	transmitters := make(map[string]struct{})
	for _, transmitter := range cfg.Transmitters {
		transmitters[transmitter] = struct{}{}
	}
	return transmitters
}

func (cfg *FeedConfig) TransmitterFromSigner() map[string]sdk.AccAddress {
	transmitterFromSigner := make(map[string]sdk.AccAddress)
	for idx, signer := range cfg.Signers {
		addr, _ := sdk.AccAddressFromBech32(cfg.Transmitters[idx])
		transmitterFromSigner[signer] = addr
	}
	return transmitterFromSigner
}

func (cfg *FeedConfig) ValidateBasic() error {
	if err := checkConfigValid(
		len(cfg.Signers),
		len(cfg.Transmitters),
		int(cfg.F),
	); err != nil {
		return err
	}

	if cfg.ModuleParams == nil {
		return errors.Wrap(ErrIncorrectConfig, "onchain config is not specified")
	}

	// TODO: determine whether this is a sensible enough limitation
	if len(cfg.ModuleParams.FeedId) == 0 || len(cfg.ModuleParams.FeedId) > FeedIDMaxLength {
		return errors.Wrap(ErrIncorrectConfig, "feed_id is missing or incorrect length")
	}

	if strings.TrimSpace(cfg.ModuleParams.FeedId) != cfg.ModuleParams.FeedId {
		return errors.Wrap(ErrIncorrectConfig, "feed_id cannot have leading or trailing space characters")
	}

	if len(cfg.ModuleParams.FeedAdmin) > 0 {
		if _, err := sdk.AccAddressFromBech32(cfg.ModuleParams.FeedAdmin); err != nil {
			return err
		}
	}

	if len(cfg.ModuleParams.BillingAdmin) > 0 {
		if _, err := sdk.AccAddressFromBech32(cfg.ModuleParams.BillingAdmin); err != nil {
			return err
		}
	}

	if cfg.ModuleParams.MinAnswer.IsNil() || cfg.ModuleParams.MaxAnswer.IsNil() {
		return errors.Wrap(ErrIncorrectConfig, "MinAnswer and MaxAnswer cannot be nil")
	}

	if cfg.ModuleParams.LinkPerTransmission.IsNil() || !cfg.ModuleParams.LinkPerTransmission.IsPositive() {
		return errors.Wrap(ErrIncorrectConfig, "LinkPerTransmission must be positive")
	}

	if cfg.ModuleParams.LinkPerObservation.IsNil() || !cfg.ModuleParams.LinkPerObservation.IsPositive() {
		return errors.Wrap(ErrIncorrectConfig, "LinkPerObservation must be positive")
	}

	seenTransmitters := make(map[string]struct{}, len(cfg.Transmitters))
	for _, transmitter := range cfg.Transmitters {
		addr, err := sdk.AccAddressFromBech32(transmitter)
		if err != nil {
			return err
		}

		if _, ok := seenTransmitters[addr.String()]; ok {
			return ErrRepeatedAddress
		}
		seenTransmitters[addr.String()] = struct{}{}
	}

	seenSigners := make(map[string]struct{}, len(cfg.Signers))
	for _, signer := range cfg.Signers {
		addr, err := sdk.AccAddressFromBech32(signer)
		if err != nil {
			return err
		}

		if _, ok := seenSigners[addr.String()]; ok {
			return ErrRepeatedAddress
		}
		seenSigners[addr.String()] = struct{}{}
	}

	if len(cfg.ModuleParams.LinkDenom) == 0 {
		return sdkerrors.ErrInvalidCoins
	}

	return nil
}

func checkConfigValid(
	numSigners, numTransmitters, f int,
) error {
	if numSigners > MaxNumOracles {
		return ErrTooManySigners
	}

	if f <= 0 {
		return errors.Wrap(ErrIncorrectConfig, "f must be positive")
	}

	if numSigners != numTransmitters {
		return errors.Wrap(ErrIncorrectConfig, "oracle addresses out of registration")
	}

	if numSigners <= 3*f {
		return errors.Wrapf(ErrIncorrectConfig, "faulty-oracle f too high: %d", f)
	}

	return nil
}

func ReportFromBytes(buf []byte) (*ReportToSign, error) {
	var r ReportToSign
	if err := proto.Unmarshal(buf, &r); err != nil {
		err = fmt.Errorf("failed to proto-decode ReportToSign from bytes: %w", err)
		return nil, err
	}

	return &r, nil
}

func (r *ReportToSign) Bytes() []byte {
	data, err := proto.Marshal(r)
	if err != nil {
		panic("unmarshable")
	}

	return data
}

func (r *ReportToSign) Digest() []byte {
	w := sha3.NewLegacyKeccak256()
	w.Write(r.Bytes())
	return w.Sum(nil)
}

type Reward struct {
	Addr   sdk.AccAddress
	Amount sdk.Coin
}

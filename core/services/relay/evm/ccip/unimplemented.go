package ccip

import (
	"context"
	"errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// The current PluginProvider interface does not support an error return. This was fine up until CCIP.
// CCIP is the first product to introduce the idea of incomplete implementations of a provider based on
// what chain (for CCIP, src or dest) the provider is created for. The Unimplemented* implementations below allow us to return
// a non nil value, which is hopefully a better developer experience should you find yourself using the right methods
// but on the *wrong* provider.

// [UnimplementedOffchainConfigDigester] satisfies the OCR OffchainConfigDigester interface
type UnimplementedOffchainConfigDigester struct{}

func (e UnimplementedOffchainConfigDigester) ConfigDigest(ctx context.Context, config ocrtypes.ContractConfig) (ocrtypes.ConfigDigest, error) {
	return ocrtypes.ConfigDigest{}, errors.New("unimplemented for this relayer")
}

func (e UnimplementedOffchainConfigDigester) ConfigDigestPrefix(ctx context.Context) (ocrtypes.ConfigDigestPrefix, error) {
	return 0, errors.New("unimplemented for this relayer")
}

// [UnimplementedContractConfigTracker] satisfies the OCR ContractConfigTracker interface
type UnimplementedContractConfigTracker struct{}

func (u UnimplementedContractConfigTracker) Notify() <-chan struct{} {
	return nil
}

func (u UnimplementedContractConfigTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	return 0, ocrtypes.ConfigDigest{}, errors.New("unimplemented for this relayer")
}

func (u UnimplementedContractConfigTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	return ocrtypes.ContractConfig{}, errors.New("unimplemented for this relayer")
}

func (u UnimplementedContractConfigTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	return 0, errors.New("unimplemented for this relayer")
}

// [UnimplementedContractTransmitter] satisfies the OCR ContractTransmitter interface
type UnimplementedContractTransmitter struct{}

func (u UnimplementedContractTransmitter) Transmit(context.Context, ocrtypes.ReportContext, ocrtypes.Report, []ocrtypes.AttributedOnchainSignature) error {
	return errors.New("unimplemented for this relayer")
}

func (u UnimplementedContractTransmitter) FromAccount(ctx context.Context) (ocrtypes.Account, error) {
	return "", errors.New("unimplemented for this relayer")
}

func (u UnimplementedContractTransmitter) LatestConfigDigestAndEpoch(ctx context.Context) (configDigest ocrtypes.ConfigDigest, epoch uint32, err error) {
	return ocrtypes.ConfigDigest{}, 0, errors.New("unimplemented for this relayer")
}

package generic

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ ocr3types.ContractTransmitter[[]byte] = (*contractTransmitter)(nil)

type contractTransmitter struct {
	impl          ocr3types.ContractTransmitter[[]byte]
	transmitterID string
}

func NewContractTransmitter(transmitterID string, impl ocr3types.ContractTransmitter[[]byte]) *contractTransmitter {
	return &contractTransmitter{
		impl:          impl,
		transmitterID: transmitterID,
	}
}

func (ct *contractTransmitter) Transmit(
	ctx context.Context,
	configDigest types.ConfigDigest,
	seqNr uint64,
	reportWithInfo ocr3types.ReportWithInfo[[]byte],
	attributedOnchainSignature []types.AttributedOnchainSignature,
) error {
	return ct.impl.Transmit(ctx, configDigest, seqNr, reportWithInfo, attributedOnchainSignature)
}

func (ct *contractTransmitter) FromAccount(ctx context.Context) (types.Account, error) {
	return types.Account(ct.transmitterID), nil
}

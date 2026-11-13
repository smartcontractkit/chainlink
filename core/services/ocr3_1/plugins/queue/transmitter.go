package queue

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type noopTransmitter struct {
	fromAccount ocrtypes.Account
}

func NewTransmitter(fromAccount ocrtypes.Account) ocr3types.ContractTransmitter[[]byte] {
	return &noopTransmitter{fromAccount: fromAccount}
}

func (c *noopTransmitter) Transmit(ctx context.Context, digest ocrtypes.ConfigDigest, seqNum uint64, rpi ocr3types.ReportWithInfo[[]byte], signatures []ocrtypes.AttributedOnchainSignature) error {
	return nil
}

func (c *noopTransmitter) FromAccount(ctx context.Context) (ocrtypes.Account, error) {
	return c.fromAccount, nil
}

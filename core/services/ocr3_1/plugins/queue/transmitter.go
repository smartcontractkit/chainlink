package queue

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type contractTransmitter struct {
	lggr        logger.Logger
	fromAccount ocrtypes.Account
}

func NewTransmitter(lggr logger.Logger, fromAccount ocrtypes.Account) ocr3types.ContractTransmitter[[]byte] {
	return &contractTransmitter{
		lggr:        lggr,
		fromAccount: fromAccount,
	}
}

func (c *contractTransmitter) Transmit(ctx context.Context, digest ocrtypes.ConfigDigest, seqNum uint64, rpi ocr3types.ReportWithInfo[[]byte], signatures []ocrtypes.AttributedOnchainSignature) error {
	return nil
}

func (c *contractTransmitter) FromAccount(ctx context.Context) (ocrtypes.Account, error) {
	return c.fromAccount, nil
}

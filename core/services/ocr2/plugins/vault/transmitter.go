package vault

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type Transmitter struct {
	lggr        logger.Logger
	store       *requests.Store[*Request]
	fromAccount types.Account
}

func NewTransmitter(lggr logger.Logger, fromAccount types.Account, store *requests.Store[*Request]) *Transmitter {
	return &Transmitter{
		lggr:        lggr.Named("VaultTransmitter"),
		store:       store,
		fromAccount: fromAccount,
	}
}

func (c *Transmitter) Transmit(ctx context.Context, cd types.ConfigDigest, seqNr uint64, rwi ocr3types.ReportWithInfo[[]byte], sigs []types.AttributedOnchainSignature) error {
	info := &vault.ReportInfo{}
	err := proto.Unmarshal(rwi.Info, info)
	if err != nil {
		return err
	}

	req := c.store.Get(info.Id)
	if req == nil {
		return fmt.Errorf("request with ID %s not found", info.Id)
	}

	// Convert the sequence number to the epoch + round number.
	// We convert as follows:
	// - epoch = seqNr
	// - round number = 0
	seqToEpoch := make([]byte, 32)
	binary.BigEndian.PutUint32(seqToEpoch[32-5:32-1], uint32(seqNr)) //nolint:gosec // the unsafe cast mirrors the OCR3OnchainKeyringAdapter implementation
	zeros := make([]byte, 32)
	responseCtx := append(append(cd[:], seqToEpoch...), zeros...)

	signatures := make([][]byte, len(sigs))
	for i, s := range sigs {
		signatures[i] = s.Signature
	}

	c.lggr.Debugw("transmitting report", "requestID", info.Id, "requestType", info.Format.String())
	req.SendResponse(ctx, &Response{
		ID:         info.Id,
		Payload:    rwi.Report,
		Format:     info.Format.String(),
		Context:    responseCtx,
		Signatures: signatures,
	})

	return nil
}

func (c *Transmitter) FromAccount(ctx context.Context) (types.Account, error) {
	return c.fromAccount, nil
}

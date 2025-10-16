package vault

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type Transmitter struct {
	lggr        logger.Logger
	handler     *requests.Handler[*vaulttypes.Request, *vaulttypes.Response]
	fromAccount types.Account
}

func NewTransmitter(lggr logger.Logger, fromAccount types.Account, handler *requests.Handler[*vaulttypes.Request, *vaulttypes.Response]) *Transmitter {
	return &Transmitter{
		lggr:        lggr.Named("VaultTransmitter"),
		handler:     handler,
		fromAccount: fromAccount,
	}
}

func extractReportInfo(rwi ocr3types.ReportWithInfo[[]byte]) (*vault.ReportInfo, error) {
	infoWrapper := &structpb.Struct{}
	err := proto.Unmarshal(rwi.Info, infoWrapper)
	if err != nil {
		return nil, err
	}

	infoWrapper.AsMap()
	reportInfoString, ok := infoWrapper.AsMap()["reportInfo"]
	if !ok {
		return nil, errors.New("reportInfo not found in report info struct")
	}

	ris, ok := reportInfoString.(string)
	if !ok {
		return nil, errors.New("reportInfo is not bytes")
	}

	rib, err := base64.StdEncoding.DecodeString(ris)
	if err != nil {
		return nil, err
	}

	reportInfo := &vault.ReportInfo{}
	err = proto.Unmarshal(rib, reportInfo)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal ReportInfo: %w", err)
	}
	return reportInfo, nil
}

func (c *Transmitter) Transmit(ctx context.Context, cd types.ConfigDigest, seqNr uint64, rwi ocr3types.ReportWithInfo[[]byte], sigs []types.AttributedOnchainSignature) error {
	info, err := extractReportInfo(rwi)
	if err != nil {
		return fmt.Errorf("could not extract report info: %w", err)
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
	c.handler.SendResponse(ctx, &vaulttypes.Response{
		ID:         info.Id,
		Payload:    rwi.Report,
		Format:     info.Format.String(),
		Context:    responseCtx,
		Signatures: signatures,
	})

	return nil
}

func (c *Transmitter) FromAccount(_ context.Context) (types.Account, error) {
	return c.fromAccount, nil
}

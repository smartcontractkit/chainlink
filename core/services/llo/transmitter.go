package llo

import (
	"context"
	"crypto/ed25519"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/services/llo/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

// LLO Transmitter implementation, based on
// core/services/relay/evm/mercury/transmitter.go

// TODO: prom metrics (common with mercury/transmitter.go?)
// https://smartcontract-it.atlassian.net/browse/MERC-3659

const (
	// Mercury server error codes
	DuplicateReport = 2
	// TODO: revisit these values in light of parallel composition
	// https://smartcontract-it.atlassian.net/browse/MERC-3659
	// maxTransmitQueueSize = 10_000
	// maxDeleteQueueSize   = 10_000
	// transmitTimeout      = 5 * time.Second
)

type Transmitter interface {
	llotypes.Transmitter
	services.Service
}

type transmitter struct {
	services.StateMachine
	lggr        logger.Logger
	rpcClient   wsrpc.Client
	fromAccount string
}

func NewTransmitter(lggr logger.Logger, rpcClient wsrpc.Client, fromAccount ed25519.PublicKey) Transmitter {
	return &transmitter{
		services.StateMachine{},
		lggr,
		rpcClient,
		fmt.Sprintf("%x", fromAccount),
	}
}

func (t *transmitter) Start(ctx context.Context) error {
	return nil
}

func (t *transmitter) Close() error {
	return nil
}

func (t *transmitter) HealthReport() map[string]error {
	report := map[string]error{t.Name(): t.Healthy()}
	services.CopyHealth(report, t.rpcClient.HealthReport())
	return report
}

func (t *transmitter) Name() string { return t.lggr.Name() }

func (t *transmitter) Transmit(
	ctx context.Context,
	digest types.ConfigDigest,
	seqNr uint64,
	report ocr3types.ReportWithInfo[llotypes.ReportInfo],
	sigs []types.AttributedOnchainSignature,
) (err error) {
	var payload []byte

	switch report.Info.ReportFormat {
	case llotypes.ReportFormatJSON:
		// TODO: exactly how to handle JSON here?
		// https://smartcontract-it.atlassian.net/browse/MERC-3659
		fallthrough
	case llotypes.ReportFormatEVMPremiumLegacy:
		payload, err = evm.ReportCodecPremiumLegacy{}.Pack(digest, seqNr, report.Report, sigs)
	default:
		return fmt.Errorf("Transmit failed; unsupported report format: %q", report.Info.ReportFormat)
	}

	if err != nil {
		return fmt.Errorf("Transmit: encode failed; %w", err)
	}

	req := &pb.TransmitRequest{
		Payload:      payload,
		ReportFormat: uint32(report.Info.ReportFormat),
	}

	// TODO: persistenceManager and queueing, error handling, retry etc
	// https://smartcontract-it.atlassian.net/browse/MERC-3659
	_, err = t.rpcClient.Transmit(ctx, req)
	return err
}

// FromAccount returns the stringified (hex) CSA public key
func (t *transmitter) FromAccount() (ocr2types.Account, error) {
	return ocr2types.Account(t.fromAccount), nil
}

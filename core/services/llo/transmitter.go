package llo

import (
	"context"
	"crypto/ed25519"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
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

var PayloadTypes = getPayloadTypes()

func getPayloadTypes() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{Name: "reportContext", Type: mustNewType("bytes32[2]")},
		{Name: "report", Type: mustNewType("bytes")},
		{Name: "rawRs", Type: mustNewType("bytes32[]")},
		{Name: "rawSs", Type: mustNewType("bytes32[]")},
		{Name: "rawVs", Type: mustNewType("bytes32")},
	})
}

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
		fallthrough
	case llotypes.ReportFormatEVM:
		payload, err = encodeEVM(digest, seqNr, report.Report, sigs)
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

func encodeEVM(digest types.ConfigDigest, seqNr uint64, report ocr2types.Report, sigs []types.AttributedOnchainSignature) ([]byte, error) {
	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
	for i, as := range sigs {
		r, s, v, err := evmutil.SplitSignature(as.Signature)
		if err != nil {
			return nil, fmt.Errorf("eventTransmit(ev): error in SplitSignature: %w", err)
		}
		rs = append(rs, r)
		ss = append(ss, s)
		vs[i] = v
	}
	rawReportCtx := ocr2key.RawReportContext3(digest, seqNr)

	payload, err := PayloadTypes.Pack(rawReportCtx, []byte(report), rs, ss, vs)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack failed; %w", err)
	}
	return payload, nil
}

// FromAccount returns the stringified (hex) CSA public key
func (t *transmitter) FromAccount() (ocr2types.Account, error) {
	return ocr2types.Account(t.fromAccount), nil
}

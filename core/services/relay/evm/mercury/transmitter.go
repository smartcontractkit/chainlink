package mercury

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

type Transmitter interface {
	relaymercury.Transmitter
	services.ServiceCtx
}

type ConfigTracker interface {
	LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error)
}

var _ Transmitter = &mercuryTransmitter{}

type mercuryTransmitter struct {
	lggr       logger.Logger
	rpcClient  wsrpc.Client
	cfgTracker ConfigTracker

	feedID      [32]byte
	fromAccount string
}

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
		{Name: "reportContext", Type: mustNewType("bytes32[3]")},
		{Name: "report", Type: mustNewType("bytes")},
		{Name: "rawRs", Type: mustNewType("bytes32[]")},
		{Name: "rawSs", Type: mustNewType("bytes32[]")},
		{Name: "rawVs", Type: mustNewType("bytes32")},
	})
}

func NewTransmitter(lggr logger.Logger, cfgTracker ConfigTracker, rpcClient wsrpc.Client, fromAccount ed25519.PublicKey, feedID [32]byte) *mercuryTransmitter {
	return &mercuryTransmitter{lggr.Named("MercuryTransmitter"), rpcClient, cfgTracker, feedID, fmt.Sprintf("%x", fromAccount)}
}

func (mt *mercuryTransmitter) Start(ctx context.Context) error { return mt.rpcClient.Start(ctx) }
func (mt *mercuryTransmitter) Close() error                    { return mt.rpcClient.Close() }
func (mt *mercuryTransmitter) Ready() error                    { return mt.rpcClient.Ready() }
func (mt *mercuryTransmitter) Name() string {
	return mt.lggr.Name()
}

func (mt *mercuryTransmitter) HealthReport() map[string]error {
	return mt.rpcClient.HealthReport()
}

// Transmit sends the report to the on-chain smart contract's Transmit method.
func (mt *mercuryTransmitter) Transmit(ctx context.Context, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signatures []ocrtypes.AttributedOnchainSignature) error {
	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
	for i, as := range signatures {
		r, s, v, err := evmutil.SplitSignature(as.Signature)
		if err != nil {
			panic("eventTransmit(ev): error in SplitSignature")
		}
		rs = append(rs, r)
		ss = append(ss, s)
		vs[i] = v
	}
	rawReportCtx := evmutil.RawReportContext(reportCtx)

	payload, err := PayloadTypes.Pack(rawReportCtx, []byte(report), rs, ss, vs)
	if err != nil {
		return errors.Wrap(err, "abi.Pack failed")
	}

	req := &pb.TransmitRequest{
		Payload: payload,
	}

	mt.lggr.Debugw("Transmitting report", "transmitRequest", req, "report", report, "reportCtx", reportCtx, "signatures", signatures)

	res, err := mt.rpcClient.Transmit(ctx, req)
	if err != nil {
		return errors.Wrap(err, "Transmit report to Mercury server failed")
	}

	if res.Error == "" {
		mt.lggr.Debugw("Transmit report success", "response", res, "reportCtx", reportCtx)
	} else {
		err := errors.New(res.Error)
		mt.lggr.Errorw("Transmit report failed; mercury server returned error", "response", res, "reportCtx", reportCtx, "err", err)
		return err
	}

	return nil
}

// FromAccount returns the stringified (hex) CSA public key
func (mt *mercuryTransmitter) FromAccount() ocrtypes.Account {
	return ocrtypes.Account(mt.fromAccount)
}

// LatestConfigDigestAndEpoch retrieves the latest config digest and epoch from the OCR2 contract.
func (mt *mercuryTransmitter) LatestConfigDigestAndEpoch(ctx context.Context) (cd ocrtypes.ConfigDigest, epoch uint32, err error) {
	mt.lggr.Debug("LatestConfigDigestAndEpoch")
	req := &pb.LatestReportRequest{
		FeedId: mt.feedID[:],
	}
	resp, err := mt.rpcClient.LatestReport(ctx, req)
	if err != nil {
		mt.lggr.Errorw("LatestConfigDigestAndEpoch failed", "err", err)
		return cd, epoch, errors.Wrap(err, "LatestConfigDigestAndEpoch failed to fetch LatestReport")
	}
	if resp == nil {
		return cd, epoch, errors.New("LatestConfigDigestAndEpoch expected LatestReport to return non-nil response")
	}
	if resp.Error != "" {
		err = errors.New(resp.Error)
		mt.lggr.Errorw("LatestConfigDigestAndEpoch failed; mercury server returned error", "err", err)
		return cd, epoch, err
	}
	if resp.Report == nil {
		_, cd, err = mt.cfgTracker.LatestConfigDetails(ctx)
		mt.lggr.Info("LatestConfigDigestAndEpoch returned empty LatestReport, this is a brand new feed")
		return cd, epoch, errors.Wrap(err, "fallback to LatestConfigDetails on empty LatestReport failed")
	}
	cd, err = ocrtypes.BytesToConfigDigest(resp.Report.ConfigDigest)
	if err != nil {
		return cd, epoch, errors.Wrapf(err, "LatestConfigDigestAndEpoch failed; response contained invalid config digest, got: 0x%x", resp.Report.ConfigDigest)
	}
	if !bytes.Equal(resp.Report.FeedId, mt.feedID[:]) {
		return cd, epoch, errors.Errorf("LatestConfigDigestAndEpoch failed; mismatched feed IDs, expected: 0x%x, got: 0x%x", mt.feedID, resp.Report.FeedId)
	}

	mt.lggr.Debugw("LatestConfigDigestAndEpoch success", "cd", cd, "epoch", epoch)

	return cd, resp.Report.Epoch, nil
}

func (mt *mercuryTransmitter) FetchInitialMaxFinalizedBlockNumber(ctx context.Context) (int64, error) {
	mt.lggr.Debug("FetchInitialMaxFinalizedBlockNumber")
	req := &pb.LatestReportRequest{
		FeedId: mt.feedID[:],
	}
	resp, err := mt.rpcClient.LatestReport(ctx, req)
	if err != nil {
		mt.lggr.Errorw("FetchInitialMaxFinalizedBlockNumber failed", "err", err)
		return 0, errors.Wrap(err, "FetchInitialMaxFinalizedBlockNumber failed to fetch LatestReport")
	}
	if resp == nil {
		return 0, errors.New("FetchInitialMaxFinalizedBlockNumber expected LatestReport to return non-nil response")
	}
	if resp.Error != "" {
		err = errors.New(resp.Error)
		mt.lggr.Errorw("FetchInitialMaxFinalizedBlockNumber failed; mercury server returned error", "err", err)
		return 0, err
	}
	if resp.Report == nil {
		mt.lggr.Infow("FetchInitialMaxFinalizedBlockNumber returned empty LatestReport; this is a new feed so initial block number is 0", "currentBlockNum", 0)
		return 0, nil
	} else if !bytes.Equal(resp.Report.FeedId, mt.feedID[:]) {
		return 0, errors.Errorf("FetchInitialMaxFinalizedBlockNumber failed; mismatched feed IDs, expected: 0x%x, got: 0x%x", mt.feedID, resp.Report.FeedId)
	}

	mt.lggr.Debugw("FetchInitialMaxFinalizedBlockNumber success", "currentBlockNum", resp.Report.CurrentBlockNumber)

	return resp.Report.CurrentBlockNumber, nil
}

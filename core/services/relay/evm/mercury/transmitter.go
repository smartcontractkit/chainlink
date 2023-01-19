package mercury

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm/mercury/wsrpc/pb"
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
func (mt *mercuryTransmitter) Healthy() error                  { return mt.rpcClient.Healthy() }
func (mt *mercuryTransmitter) Ready() error                    { return mt.rpcClient.Ready() }
func (mt *mercuryTransmitter) Name() string {
	return mt.lggr.Name()
}

func (mt *mercuryTransmitter) HealthReport() map[string]error {
	return map[string]error{mt.Name(): mt.Healthy()}
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
		mt.lggr.Errorw("Transmit report failed", "response", res, "reportCtx", reportCtx)

	}

	return nil
}

// FromAccount returns the stringified (hex) CSA public key
func (mt *mercuryTransmitter) FromAccount() ocrtypes.Account {
	return ocrtypes.Account(mt.fromAccount)
}

// LatestConfigDigestAndEpoch retrieves the latest config digest and epoch from the OCR2 contract.
func (mt *mercuryTransmitter) LatestConfigDigestAndEpoch(ctx context.Context) (cd ocrtypes.ConfigDigest, epoch uint32, err error) {
	req := &pb.LatestReportRequest{
		FeedId: mt.feedID[:],
	}
	resp, err := mt.rpcClient.LatestReport(ctx, req)
	if err != nil {
		return cd, epoch, errors.Wrap(err, "LatestConfigDigestAndEpoch failed to fetch LatestReport")
	}
	if len(resp.ConfigDigest) == 0 {
		mt.cfgTracker.LatestConfigDetails(ctx)
		mt.lggr.Info("LatestConfigDigestAndEpoch returned nil reponse, this is a brand new feed")
		return cd, epoch, nil
	}
	cd, err = ocrtypes.BytesToConfigDigest(resp.ConfigDigest)
	if err != nil {
		return cd, epoch, errors.Wrapf(err, "LatestConfigDigestAndEpoch failed; response contained invalid config digest, got: 0x%x", resp.ConfigDigest)
	}
	if !bytes.Equal(resp.FeedId, mt.feedID[:]) {
		return cd, epoch, errors.Errorf("LatestConfigDigestAndEpoch failed; mismatched feed IDs, expected: 0x%x, got: 0x%x", mt.feedID, resp.FeedId)
	}

	return cd, resp.Epoch, nil
}

func (mt *mercuryTransmitter) FetchInitialMaxFinalizedBlockNumber(ctx context.Context) (int64, error) {
	req := &pb.LatestReportRequest{
		FeedId: mt.feedID[:],
	}
	resp, err := mt.rpcClient.LatestReport(ctx, req)
	if err != nil {
		return 0, errors.Wrap(err, "FetchInitialMaxFinalizedBlockNumber failed to fetch LatestReport")
	}
	if len(resp.FeedId) == 0 {
		mt.lggr.Infow("FetchInitialMaxFinalizedBlockNumber returned empty LatestReport; this is a new feed so initial block number is 0")
		return 0, nil
	} else if !bytes.Equal(resp.FeedId, mt.feedID[:]) {
		return 0, errors.Errorf("FetchInitialMaxFinalizedBlockNumber failed; mismatched feed IDs, expected: 0x%x, got: 0x%x", mt.feedID, resp.FeedId)
	}

	return resp.BlockNumber, nil
}

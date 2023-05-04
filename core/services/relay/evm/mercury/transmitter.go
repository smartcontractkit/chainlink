package mercury

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/jpillora/backoff"
	pkgerrors "github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/exp/maps"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// TODO: Revisit this choice of value
// TODO: Ought we to have one queue per mercury server URL instead?
const MaxTransmitQueueSize = 10_000

type Transmitter interface {
	relaymercury.Transmitter
	services.ServiceCtx
}

type ConfigTracker interface {
	LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error)
}

var _ Transmitter = &mercuryTransmitter{}

type mercuryTransmitter struct {
	utils.StartStopOnce
	lggr       logger.Logger
	rpcClient  wsrpc.Client
	cfgTracker ConfigTracker

	feedID      [32]byte
	fromAccount string

	stopCh utils.StopChan
	queue  *TransmitQueue
	wg     sync.WaitGroup
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
	return &mercuryTransmitter{
		utils.StartStopOnce{},
		lggr.Named("MercuryTransmitter"),
		rpcClient,
		cfgTracker,
		feedID,
		fmt.Sprintf("%x", fromAccount),
		make(chan (struct{})),
		NewTransmitQueue(lggr, MaxTransmitQueueSize),
		sync.WaitGroup{},
	}
}

func (mt *mercuryTransmitter) Start(ctx context.Context) (err error) {
	return mt.StartOnce("MercuryTransmitter", func() error {
		if err := mt.rpcClient.Start(ctx); err != nil {
			return err
		}
		mt.wg.Add(1)
		go mt.runloop()
		return nil
	})
}

func (mt *mercuryTransmitter) Close() error {
	return mt.StopOnce("MercuryTransmitter", func() error {
		mt.queue.Close()
		close(mt.stopCh)
		mt.wg.Wait()
		return mt.rpcClient.Close()
	})
}
func (mt *mercuryTransmitter) Ready() error { return mt.StartStopOnce.Ready() }
func (mt *mercuryTransmitter) Name() string { return mt.lggr.Name() }

func (mt *mercuryTransmitter) HealthReport() map[string]error {
	report := map[string]error{mt.Name(): mt.StartStopOnce.Healthy()}
	maps.Copy(report, mt.rpcClient.HealthReport())
	maps.Copy(report, mt.queue.HealthReport())
	return report
}

func (mt *mercuryTransmitter) runloop() {
	defer mt.wg.Done()
	// Exponential backoff with very short retry interval (since latency is a priority)
	// 5ms, 10ms, 20ms, 40ms etc
	b := backoff.Backoff{
		Min:    5 * time.Millisecond,
		Max:    1 * time.Second,
		Factor: 2,
		Jitter: true,
	}
	ctx, cancel := mt.stopCh.Ctx(context.Background())
	defer cancel()
	for {
		t := mt.queue.BlockingPop()
		if t == nil {
			// queue was closed
			return
		}
		res, err := mt.rpcClient.Transmit(ctx, t.Req)
		if ctx.Err() != nil {
			// context only canceled on transmitter close so we can exit
			// the runloop here
			return
		} else if err != nil {
			// TODO: log this
			if ok := mt.queue.Push(t.Req, t.ReportCtx); !ok {
				// TODO: log this?
				return
			}
			// Wait a backoff duration before pulling the latest back off
			// the heap
			select {
			case <-time.After(b.Duration()):
				continue
			case <-mt.stopCh:
				return
			}
		} else {
			b.Reset()
			if res.Error == "" {
				mt.lggr.Debugw("Transmit report success", "req", t.Req, "response", res, "reportCtx", t.ReportCtx)
			} else {
				// We don't need to retry here because the mercury server
				// has confirmed it received the report. We only need to retry
				// on networking/unknown errors
				err := errors.New(res.Error)
				mt.lggr.Errorw("Transmit report failed; mercury server returned error", "req", t.Req, "response", res, "reportCtx", t.ReportCtx, "err", err)
			}
			continue
		}
	}
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
		return pkgerrors.Wrap(err, "abi.Pack failed")
	}

	req := &pb.TransmitRequest{
		Payload: payload,
	}

	mt.lggr.Debugw("Transmit enqueue", "req", req, "report", report, "reportCtx", reportCtx, "signatures", signatures)

	if ok := mt.queue.Push(req, reportCtx); !ok {
		return errors.New("transmit queue is closed")
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
		return cd, epoch, pkgerrors.Wrap(err, "LatestConfigDigestAndEpoch failed to fetch LatestReport")
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
		return cd, epoch, pkgerrors.Wrap(err, "fallback to LatestConfigDetails on empty LatestReport failed")
	}
	cd, err = ocrtypes.BytesToConfigDigest(resp.Report.ConfigDigest)
	if err != nil {
		return cd, epoch, pkgerrors.Wrapf(err, "LatestConfigDigestAndEpoch failed; response contained invalid config digest, got: 0x%x", resp.Report.ConfigDigest)
	}
	if !bytes.Equal(resp.Report.FeedId, mt.feedID[:]) {
		return cd, epoch, fmt.Errorf("LatestConfigDigestAndEpoch failed; mismatched feed IDs, expected: 0x%x, got: 0x%x", mt.feedID, resp.Report.FeedId)
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
		return 0, pkgerrors.Wrap(err, "FetchInitialMaxFinalizedBlockNumber failed to fetch LatestReport")
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
		return 0, fmt.Errorf("FetchInitialMaxFinalizedBlockNumber failed; mismatched feed IDs, expected: 0x%x, got: 0x%x", mt.feedID, resp.Report.FeedId)
	}

	mt.lggr.Debugw("FetchInitialMaxFinalizedBlockNumber success", "currentBlockNum", resp.Report.CurrentBlockNumber)

	return resp.Report.CurrentBlockNumber, nil
}

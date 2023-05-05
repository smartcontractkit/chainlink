package mercury

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jpillora/backoff"
	pkgerrors "github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
)

// Mercury-specific reporting plugin, based off of median:
// https://github.com/smartcontractkit/offchain-reporting/blob/master/lib/offchainreporting2/reportingplugin/median/median.go

const onchainConfigVersion = 1
const onchainConfigEncodedLength = 1 + byteWidthInt192 + byteWidthInt192

type OnchainConfig struct {
	// applies to all values: price, bid and ask
	Min *big.Int
	Max *big.Int
}

type OnchainConfigCodec interface {
	Encode(OnchainConfig) ([]byte, error)
	Decode([]byte) (OnchainConfig, error)
}

var _ OnchainConfigCodec = StandardOnchainConfigCodec{}

// StandardOnchainConfigCodec provides a mercury-specific implementation of
// OnchainConfigCodec.
//
// An encoded onchain config is expected to be in the format
// <version><min><max>
// where version is a uint8 and min and max are in the format
// returned by EncodeValueInt192.
type StandardOnchainConfigCodec struct{}

func (StandardOnchainConfigCodec) Decode(b []byte) (OnchainConfig, error) {
	if len(b) != onchainConfigEncodedLength {
		return OnchainConfig{}, pkgerrors.Errorf("unexpected length of OnchainConfig, expected %v, got %v", onchainConfigEncodedLength, len(b))
	}

	if b[0] != onchainConfigVersion {
		return OnchainConfig{}, pkgerrors.Errorf("unexpected version of OnchainConfig, expected %v, got %v", onchainConfigVersion, b[0])
	}

	min, err := DecodeValueInt192(b[1 : 1+byteWidthInt192])
	if err != nil {
		return OnchainConfig{}, err
	}
	max, err := DecodeValueInt192(b[1+byteWidthInt192:])
	if err != nil {
		return OnchainConfig{}, err
	}

	if !(min.Cmp(max) <= 0) {
		return OnchainConfig{}, pkgerrors.Errorf("OnchainConfig min (%v) should not be greater than max(%v)", min, max)
	}

	return OnchainConfig{min, max}, nil
}

func (StandardOnchainConfigCodec) Encode(c OnchainConfig) ([]byte, error) {
	minBytes, err := EncodeValueInt192(c.Min)
	if err != nil {
		return nil, err
	}
	maxBytes, err := EncodeValueInt192(c.Max)
	if err != nil {
		return nil, err
	}
	result := make([]byte, 0, onchainConfigEncodedLength)
	result = append(result, onchainConfigVersion)
	result = append(result, minBytes...)
	result = append(result, maxBytes...)
	return result, nil
}

type OffchainConfig struct{}

func DecodeOffchainConfig(b []byte) (o OffchainConfig, err error) {
	return
}

func (c OffchainConfig) Encode() []byte {
	return []byte{}
}

type ObsResult[T any] struct {
	Val T
	Err error
}

type Observation struct {
	BenchmarkPrice        ObsResult[*big.Int]
	Bid                   ObsResult[*big.Int]
	Ask                   ObsResult[*big.Int]
	CurrentBlockNum       ObsResult[int64]
	CurrentBlockHash      ObsResult[[]byte]
	CurrentBlockTimestamp ObsResult[uint64]
}

// DataSource implementations must be thread-safe. Observe may be called by many
// different threads concurrently.
type DataSource interface {
	// Observe queries the data source. Returns a value or an error. Once the
	// context is expires, Observe may still do cheap computations and return a
	// result, but should return as quickly as possible.
	//
	// More details: In the current implementation, the context passed to
	// Observe will time out after MaxDurationObservation. However, Observe
	// should *not* make any assumptions about context timeout behavior. Once
	// the context times out, Observe should prioritize returning as quickly as
	// possible, but may still perform fast computations to return a result
	// rather than error. For example, if Observe medianizes a number of data
	// sources, some of which already returned a result to Observe prior to the
	// context's expiry, Observe might still compute their median, and return it
	// instead of an error.
	//
	// Important: Observe should not perform any potentially time-consuming
	// actions like database access, once the context passed has expired.
	Observe(context.Context, ocrtypes.ReportTimestamp) (Observation, error)
}

var _ ocrtypes.ReportingPluginFactory = Factory{}

const maxObservationLength = 32 + // feedID
	4 + // timestamp
	byteWidthInt192 + // benchmarkPrice
	byteWidthInt192 + // bid
	byteWidthInt192 + // ask
	8 + // currentBlockNum
	32 + // currentBlockHash
	8 + // currentBlockTimestamp
	8 + // validFromBlockNum
	16 /* overapprox. of protobuf overhead */

// All functions on ReportCodec should be pure and thread-safe.
// Be careful validating and parsing any data passed.
type ReportCodec interface {
	// Implementers may assume that there is at most one
	// ParsedAttributedObservation per observer, and that all observers are
	// valid. However, observation values, timestamps, etc... should all be
	// treated as untrusted.
	BuildReport(paos []ParsedAttributedObservation, f int) (ocrtypes.Report, error)

	// Returns the maximum length of a report based on n, the number of oracles.
	// The output of BuildReport must respect this maximum length.
	MaxReportLength(n int) int

	// CurrentBlockNumFromReport returns the median current block number from a report
	CurrentBlockNumFromReport(types.Report) (int64, error)
}

const unfetchedInitialMaxFinalizedBlockNumber int64 = -1

func newInitialMaxFinalizedBlockNumber() (a *atomic.Int64) {
	a = new(atomic.Int64)
	a.Store(unfetchedInitialMaxFinalizedBlockNumber)
	return
}

type Fetcher interface {
	// FetchInitialMaxFinalizedBlockNumber should fetch the initial max
	// finalized block number from the mercury server.
	FetchInitialMaxFinalizedBlockNumber(context.Context) (int64, error)
}

type Transmitter interface {
	Fetcher
	// NOTE: Mercury doesn't actually transmit on-chain, so there is no
	// "contract" involved with the transmitter.
	// - Transmit should be implemented and send to Mercury server
	// - LatestConfigDigestAndEpoch should be implemented and fetch from Mercury server
	// - FromAccount() should return CSA public key
	ocrtypes.ContractTransmitter
}

type Factory struct {
	dataSource         DataSource
	logger             logger.Logger
	onchainConfigCodec OnchainConfigCodec
	reportCodec        ReportCodec
	fetcher            Fetcher
}

func NewFactory(ds DataSource, lggr logger.Logger, occ OnchainConfigCodec, rc ReportCodec, f Fetcher) Factory {
	return Factory{ds, lggr, occ, rc, f}
}

func (fac Factory) NewReportingPlugin(configuration ocrtypes.ReportingPluginConfig) (ocrtypes.ReportingPlugin, ocrtypes.ReportingPluginInfo, error) {
	offchainConfig, err := DecodeOffchainConfig(configuration.OffchainConfig)
	if err != nil {
		return nil, ocrtypes.ReportingPluginInfo{}, err
	}

	onchainConfig, err := fac.onchainConfigCodec.Decode(configuration.OnchainConfig)
	if err != nil {
		return nil, ocrtypes.ReportingPluginInfo{}, err
	}

	maxReportLength := fac.reportCodec.MaxReportLength(configuration.N)

	wg := sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	r := &reportingPlugin{
		offchainConfig,
		onchainConfig,
		fac.dataSource,
		fac.logger,
		fac.reportCodec,
		configuration.ConfigDigest,
		configuration.F,
		epochRound{},
		new(big.Int),
		maxReportLength,
		newInitialMaxFinalizedBlockNumber(),
		sync.WaitGroup{},
		cancel,
		sync.Once{},
	}

	go func() {
		defer wg.Done()

		b := backoff.Backoff{
			Min: 1 * time.Second,
			Max: 10 * time.Second,
		}
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(b.Duration()):
				initial, err := fac.fetcher.FetchInitialMaxFinalizedBlockNumber(ctx)
				if err != nil {
					fac.logger.Warnw("FetchInitialMaxFinalizedBlockNumber failed", "err", err)
					continue
				}
				r.maxFinalizedBlockNumber.CompareAndSwap(unfetchedInitialMaxFinalizedBlockNumber, initial)
				return
			}
		}
	}()

	return r, ocrtypes.ReportingPluginInfo{
		Name:          "Mercury",
		UniqueReports: false,
		Limits: ocrtypes.ReportingPluginLimits{
			MaxQueryLength:       0,
			MaxObservationLength: maxObservationLength,
			MaxReportLength:      maxReportLength,
		},
	}, nil
}

var _ ocrtypes.ReportingPlugin = (*reportingPlugin)(nil)

type reportingPlugin struct {
	offchainConfig OffchainConfig
	onchainConfig  OnchainConfig
	dataSource     DataSource
	logger         logger.Logger
	reportCodec    ReportCodec

	configDigest             ocrtypes.ConfigDigest
	f                        int
	latestAcceptedEpochRound epochRound
	latestAcceptedMedian     *big.Int
	maxReportLength          int
	maxFinalizedBlockNumber  *atomic.Int64

	// fetch initial finalized block number state management
	wg         sync.WaitGroup
	cancel     context.CancelFunc
	cancelOnce sync.Once
}

func (rp *reportingPlugin) Query(ctx context.Context, repts ocrtypes.ReportTimestamp) (ocrtypes.Query, error) {
	return nil, nil
}

func (rp *reportingPlugin) Observation(ctx context.Context, repts ocrtypes.ReportTimestamp, query ocrtypes.Query) (ocrtypes.Observation, error) {
	if len(query) != 0 {
		return nil, errors.New("expected empty query")
	}

	obs, err := rp.dataSource.Observe(ctx, repts)
	if err != nil {
		return nil, pkgerrors.Errorf("DataSource.Observe returned an error: %s", err)
	}

	p := MercuryObservationProto{Timestamp: uint32(time.Now().Unix())}

	var obsErrors []error
	if maxFinalizedBlockNumber := rp.maxFinalizedBlockNumber.Load(); maxFinalizedBlockNumber == unfetchedInitialMaxFinalizedBlockNumber {
		obsErrors = append(obsErrors, errors.New("failed to observe ValidFromBlockNum; initial maxFinalizedBlockNumber has not yet been fetched from the mercury server"))
	} else if obs.CurrentBlockNum.Err == nil && obs.CurrentBlockNum.Val < maxFinalizedBlockNumber {
		obsErrors = append(obsErrors, pkgerrors.Errorf("failed to observe ValidFromBlockNum; current block number %d (hash: 0x%x) < max finalized block number %d; ignoring observation for out-of-date RPC", obs.CurrentBlockNum, obs.CurrentBlockHash, maxFinalizedBlockNumber))
	} else {
		// NOTE: obs.CurrentBlockNum == maxFinalizedBlockNumber is ok here
		// (multiple observations for the same block number) since it will be
		// de-duplicated in the Report stage
		p.ValidFromBlockNum = maxFinalizedBlockNumber + 1
		p.ValidFromBlockNumValid = true
	}

	if obs.BenchmarkPrice.Err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(obs.BenchmarkPrice.Err, "failed to observe BenchmarkPrice"))
	} else if benchmarkPrice, err := EncodeValueInt192(obs.BenchmarkPrice.Val); err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(err, "failed to observe BenchmarkPrice; encoding failed"))
	} else {
		p.BenchmarkPrice = benchmarkPrice
		p.BenchmarkPriceValid = true
	}

	if obs.Bid.Err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(obs.Bid.Err, "failed to observe Bid"))
	} else if bid, err := EncodeValueInt192(obs.Bid.Val); err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(err, "failed to observe Bid; encoding failed"))
	} else {
		p.Bid = bid
		p.BidValid = true
	}

	if obs.Ask.Err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(obs.Ask.Err, "failed to observe Ask"))
	} else if bid, err := EncodeValueInt192(obs.Ask.Val); err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(err, "failed to observe Ask; encoding failed"))
	} else {
		p.Ask = bid
		p.AskValid = true
	}

	if obs.CurrentBlockNum.Err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(obs.CurrentBlockNum.Err, "failed to observe CurrentBlockNum"))
	} else {
		p.CurrentBlockNum = obs.CurrentBlockNum.Val
		p.CurrentBlockNumValid = true
	}

	if obs.CurrentBlockHash.Err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(obs.CurrentBlockHash.Err, "failed to observe CurrentBlockHash"))
	} else {
		p.CurrentBlockHash = obs.CurrentBlockHash.Val
		p.CurrentBlockHashValid = true
	}

	if obs.CurrentBlockTimestamp.Err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(obs.CurrentBlockTimestamp.Err, "failed to observe CurrentBlockTimestamp"))
	} else {
		p.CurrentBlockTimestamp = obs.CurrentBlockTimestamp.Val
		p.CurrentBlockTimestampValid = true
	}

	if len(obsErrors) > 0 {
		rp.logger.Warnw(fmt.Sprintf("Observe failed %d/6 observations", len(obsErrors)), "err", errors.Join(obsErrors...))
	}

	return proto.Marshal(&p)
}

type ParsedAttributedObservation struct {
	Timestamp             uint32
	BenchmarkPrice        *big.Int
	Bid                   *big.Int
	Ask                   *big.Int
	CurrentBlockNum       int64 // inclusive; current block
	CurrentBlockHash      []byte
	CurrentBlockTimestamp uint64
	ValidFromBlockNum     int64 // exclusive; one above previous upper block
	Observer              commontypes.OracleID
}

func parseAttributedObservation(ao ocrtypes.AttributedObservation) (ParsedAttributedObservation, error) {
	var obs MercuryObservationProto
	if err := proto.Unmarshal(ao.Observation, &obs); err != nil {
		return ParsedAttributedObservation{}, pkgerrors.Errorf("attributed observation cannot be unmarshaled: %s", err)
	}
	benchmarkPrice, err := DecodeValueInt192(obs.BenchmarkPrice)
	if err != nil {
		return ParsedAttributedObservation{}, pkgerrors.Errorf("benchmarkPrice cannot be converted to big.Int: %s", err)
	}
	bid, err := DecodeValueInt192(obs.Bid)
	if err != nil {
		return ParsedAttributedObservation{}, pkgerrors.Errorf("bid cannot be converted to big.Int: %s", err)
	}
	ask, err := DecodeValueInt192(obs.Ask)
	if err != nil {
		return ParsedAttributedObservation{}, pkgerrors.Errorf("ask cannot be converted to big.Int: %s", err)
	}
	if len(obs.CurrentBlockHash) == 0 {
		return ParsedAttributedObservation{}, pkgerrors.Errorf("wrong len for hash: %d", len(obs.CurrentBlockHash))
	}

	return ParsedAttributedObservation{
		obs.Timestamp,
		benchmarkPrice,
		bid,
		ask,
		obs.CurrentBlockNum,
		obs.CurrentBlockHash,
		obs.CurrentBlockTimestamp,
		obs.ValidFromBlockNum,
		ao.Observer,
	}, nil
}

func parseAttributedObservations(lggr logger.Logger, aos []ocrtypes.AttributedObservation) []ParsedAttributedObservation {
	paos := make([]ParsedAttributedObservation, 0, len(aos))
	for i, ao := range aos {
		pao, err := parseAttributedObservation(ao)
		if err != nil {
			lggr.Warnw("parseAttributedObservations: dropping invalid observation",
				"observer", ao.Observer,
				"error", err,
				"i", i,
			)
			continue
		}
		paos = append(paos, pao)
	}
	return paos
}

func (rp *reportingPlugin) Report(ctx context.Context, repts types.ReportTimestamp, query types.Query, aos []types.AttributedObservation) (bool, types.Report, error) {
	if len(query) != 0 {
		return false, nil, pkgerrors.Errorf("expected empty query")
	}

	paos := parseAttributedObservations(rp.logger, aos)

	// By assumption, we have at most f malicious oracles, so there should be at least f+1 valid paos
	if !(rp.f+1 <= len(paos)) {
		return false, nil, pkgerrors.Errorf("only received %v valid attributed observations, but need at least f+1 (%v)", len(paos), rp.f+1)
	}

	should, err := rp.shouldReport(ctx, repts, paos)
	if err != nil {
		return false, nil, err
	}
	if !should {
		return false, nil, nil
	}
	report, err := rp.reportCodec.BuildReport(paos, rp.f)
	if err != nil {
		return false, nil, err
	}
	if !(len(report) <= rp.maxReportLength) {
		return false, nil, pkgerrors.Errorf("report violates MaxReportLength limit set by ReportCodec (%v vs %v)", len(report), rp.maxReportLength)
	}

	return true, report, nil
}

func (rp *reportingPlugin) shouldReport(ctx context.Context, repts types.ReportTimestamp, paos []ParsedAttributedObservation) (bool, error) {
	if !(rp.f+1 <= len(paos)) {
		return false, pkgerrors.Errorf("only received %v valid attributed observations, but need at least f+1 (%v)", len(paos), rp.f+1)
	}

	if err := errors.Join(
		rp.checkBenchmarkPrice(paos),
		rp.checkBid(paos),
		rp.checkAsk(paos),
		rp.checkBlockValues(paos),
	); err != nil {
		rp.logger.Debugw("shouldReport: no", "err", err)
		return false, nil
	}

	rp.logger.Debugw("shouldReport: yes",
		"timestamp", repts,
	)
	return true, nil
}

func (rp *reportingPlugin) checkBenchmarkPrice(paos []ParsedAttributedObservation) error {
	return ValidateBenchmarkPrice(paos, rp.onchainConfig.Min, rp.onchainConfig.Max)
}

func (rp *reportingPlugin) checkBid(paos []ParsedAttributedObservation) error {
	return ValidateBid(paos, rp.onchainConfig.Min, rp.onchainConfig.Max)
}

func (rp *reportingPlugin) checkAsk(paos []ParsedAttributedObservation) error {
	return ValidateAsk(paos, rp.onchainConfig.Min, rp.onchainConfig.Max)
}

func (rp *reportingPlugin) checkBlockValues(paos []ParsedAttributedObservation) error {
	return ValidateBlockValues(paos, rp.f, rp.maxFinalizedBlockNumber.Load())
}

func (rp *reportingPlugin) ShouldAcceptFinalizedReport(ctx context.Context, repts types.ReportTimestamp, report types.Report) (bool, error) {
	reportEpochRound := epochRound{repts.Epoch, repts.Round}
	if !rp.latestAcceptedEpochRound.Less(reportEpochRound) {
		rp.logger.Debugw("ShouldAcceptFinalizedReport() = false, report is stale",
			"latestAcceptedEpochRound", rp.latestAcceptedEpochRound,
			"reportEpochRound", reportEpochRound,
		)
		return false, nil
	}

	if !(len(report) <= rp.maxReportLength) {
		rp.logger.Warnw("report violates MaxReportLength limit set by ReportCodec",
			"reportEpochRound", reportEpochRound,
			"reportLength", len(report),
			"maxReportLength", rp.maxReportLength,
		)
		return false, nil
	}

	currentBlockNum, err := rp.reportCodec.CurrentBlockNumFromReport(report)
	if err != nil {
		return false, pkgerrors.Wrap(err, "error during CurrentBlockNumFromReport")
	}

	rp.logger.Debugw("ShouldAcceptFinalizedReport() = true",
		"reportEpochRound", reportEpochRound,
		"latestAcceptedEpochRound", rp.latestAcceptedEpochRound,
	)

	if currentBlockNum > rp.maxFinalizedBlockNumber.Load() {
		rp.cancelOnce.Do(rp.cancel) // abort fetch because we will store the value from the protocol instead
		rp.maxFinalizedBlockNumber.Store(currentBlockNum)
	}
	rp.latestAcceptedEpochRound = reportEpochRound

	return true, nil
}

func (rp *reportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, repts types.ReportTimestamp, report types.Report) (bool, error) {
	return true, nil
}

func (rp *reportingPlugin) Close() error {
	rp.cancelOnce.Do(rp.cancel)
	rp.wg.Wait()
	return nil
}

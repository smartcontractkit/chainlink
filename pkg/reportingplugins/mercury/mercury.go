package mercury

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	pkgerrors "github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/bigbigendian"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
)

// Mercury-specific reporting plugin, based off of median:
// https://github.com/smartcontractkit/offchain-reporting/blob/master/lib/offchainreporting2/reportingplugin/median/median.go

const onchainConfigVersion = 1

var onchainConfigVersionBig = big.NewInt(onchainConfigVersion)

const onchainConfigEncodedLength = 96 // 3x 32bit evm words, version + min + max

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

	v, err := bigbigendian.DeserializeSigned(32, b[:32])
	if err != nil {
		return OnchainConfig{}, err
	}
	if v.Cmp(onchainConfigVersionBig) != 0 {
		return OnchainConfig{}, pkgerrors.Errorf("unexpected version of OnchainConfig, expected %v, got %v", onchainConfigVersion, v)
	}

	min, err := bigbigendian.DeserializeSigned(32, b[32:64])
	if err != nil {
		return OnchainConfig{}, err
	}
	max, err := bigbigendian.DeserializeSigned(32, b[64:96])
	if err != nil {
		return OnchainConfig{}, err
	}

	if !(min.Cmp(max) <= 0) {
		return OnchainConfig{}, pkgerrors.Errorf("OnchainConfig min (%v) should not be greater than max(%v)", min, max)
	}

	return OnchainConfig{min, max}, nil
}

func (StandardOnchainConfigCodec) Encode(c OnchainConfig) ([]byte, error) {
	verBytes, err := bigbigendian.SerializeSigned(32, onchainConfigVersionBig)
	if err != nil {
		return nil, err
	}
	minBytes, err := bigbigendian.SerializeSigned(32, c.Min)
	if err != nil {
		return nil, err
	}
	maxBytes, err := bigbigendian.SerializeSigned(32, c.Max)
	if err != nil {
		return nil, err
	}
	result := make([]byte, 0, onchainConfigEncodedLength)
	result = append(result, verBytes...)
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
	// MaxFinalizedBlockNumber comes from previous report when present and is
	// only observed from mercury server when previous report is nil
	MaxFinalizedBlockNumber ObsResult[int64]
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
	Observe(ctx context.Context, repts ocrtypes.ReportTimestamp, fetchMaxFinalizedBlockNum bool) (Observation, error)
}

var _ ocr3types.MercuryPluginFactory = Factory{}

// Maximum length in bytes of Observation, Report returned by the
// MercuryPlugin. Used for defending against spam attacks.
const maxObservationLength = 4 + // timestamp
	byteWidthInt192 + // benchmarkPrice
	byteWidthInt192 + // bid
	byteWidthInt192 + // ask
	1 + // pricesValid
	8 + // currentBlockNum
	32 + // currentBlockHash
	8 + // currentBlockTimestamp
	1 + // currentBlockValid
	8 + // maxFinalizedBlockNumber
	1 + // maxFinalizedBlockNumberValid
	32 // [> overapprox. of protobuf overhead <]

// ReportCodec All functions on ReportCodec should be pure and thread-safe.
// Be careful validating and parsing any data passed.
type ReportCodec interface {
	// BuildReport Implementers may assume that there is at most one
	// ParsedAttributedObservation per observer, and that all observers are
	// valid. However, observation values, timestamps, etc... should all be
	// treated as untrusted.
	BuildReport(paos []ParsedAttributedObservation, f int, validFromBlockNum int64) (ocrtypes.Report, error)

	// MaxReportLength Returns the maximum length of a report based on n, the number of oracles.
	// The output of BuildReport must respect this maximum length.
	MaxReportLength(n int) (int, error)

	// CurrentBlockNumFromReport returns the median current block number from a report
	CurrentBlockNumFromReport(types.Report) (int64, error)
}

type Fetcher interface {
	// FetchInitialMaxFinalizedBlockNumber should fetch the initial max
	// finalized block number from the mercury server.
	FetchInitialMaxFinalizedBlockNumber(context.Context) (*int64, error)
}

type Transmitter interface {
	Fetcher
	// NOTE: Mercury doesn't actually transmit on-chain, so there is no
	// "contract" involved with the transmitter.
	// - Transmit should be implemented and send to Mercury server
	// - LatestConfigDigestAndEpoch is a stub method, does not need to do anything
	// - FromAccount() should return CSA public key
	ocrtypes.ContractTransmitter
}

type Factory struct {
	dataSource         DataSource
	logger             logger.Logger
	onchainConfigCodec OnchainConfigCodec
	reportCodec        ReportCodec
}

func NewFactory(ds DataSource, lggr logger.Logger, occ OnchainConfigCodec, rc ReportCodec) Factory {
	return Factory{ds, lggr, occ, rc}
}

func (fac Factory) NewMercuryPlugin(configuration ocr3types.MercuryPluginConfig) (ocr3types.MercuryPlugin, ocr3types.MercuryPluginInfo, error) {
	offchainConfig, err := DecodeOffchainConfig(configuration.OffchainConfig)
	if err != nil {
		return nil, ocr3types.MercuryPluginInfo{}, err
	}

	onchainConfig, err := fac.onchainConfigCodec.Decode(configuration.OnchainConfig)
	if err != nil {
		return nil, ocr3types.MercuryPluginInfo{}, err
	}

	maxReportLength, err := fac.reportCodec.MaxReportLength(configuration.N)
	if err != nil {
		return nil, ocr3types.MercuryPluginInfo{}, err
	}

	r := &reportingPlugin{
		offchainConfig,
		onchainConfig,
		fac.dataSource,
		fac.logger,
		fac.reportCodec,
		configuration.ConfigDigest,
		configuration.F,
		epochRound{},
		maxReportLength,
	}

	return r, ocr3types.MercuryPluginInfo{
		Name: "Mercury",
		Limits: ocr3types.MercuryPluginLimits{
			MaxObservationLength: maxObservationLength,
			MaxReportLength:      maxReportLength,
		},
	}, nil
}

var _ ocr3types.MercuryPlugin = (*reportingPlugin)(nil)

type reportingPlugin struct {
	offchainConfig OffchainConfig
	onchainConfig  OnchainConfig
	dataSource     DataSource
	logger         logger.Logger
	reportCodec    ReportCodec

	configDigest             ocrtypes.ConfigDigest
	f                        int
	latestAcceptedEpochRound epochRound
	maxReportLength          int
}

func (rp *reportingPlugin) Query(ctx context.Context, repts ocrtypes.ReportTimestamp) (ocrtypes.Query, error) {
	return nil, nil
}

func (rp *reportingPlugin) Observation(ctx context.Context, repts ocrtypes.ReportTimestamp, previousReport types.Report) (ocrtypes.Observation, error) {
	obs, err := rp.dataSource.Observe(ctx, repts, previousReport == nil)
	if err != nil {
		return nil, pkgerrors.Errorf("DataSource.Observe returned an error: %s", err)
	}

	p := MercuryObservationProto{Timestamp: uint32(time.Now().Unix())}

	var obsErrors []error
	if previousReport == nil {
		// if previousReport we fall back to the observed MaxFinalizedBlockNumber
		if obs.MaxFinalizedBlockNumber.Err != nil {
			obsErrors = append(obsErrors, err)
		} else if obs.CurrentBlockNum.Err == nil && obs.CurrentBlockNum.Val < obs.MaxFinalizedBlockNumber.Val {
			obsErrors = append(obsErrors, pkgerrors.Errorf("failed to observe ValidFromBlockNum; current block number %d (hash: 0x%x) < max finalized block number %d; ignoring observation for out-of-date RPC", obs.CurrentBlockNum.Val, obs.CurrentBlockHash.Val, obs.MaxFinalizedBlockNumber.Val))
		} else {
			p.MaxFinalizedBlockNumber = obs.MaxFinalizedBlockNumber.Val // MaxFinalizedBlockNumber comes as -1 if unset
			p.MaxFinalizedBlockNumberValid = true
		}
	}

	if obs.BenchmarkPrice.Err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(obs.BenchmarkPrice.Err, "failed to observe BenchmarkPrice"))
	} else if benchmarkPrice, err := EncodeValueInt192(obs.BenchmarkPrice.Val); err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(err, "failed to observe BenchmarkPrice; encoding failed"))
	} else {
		p.BenchmarkPrice = benchmarkPrice
	}

	if obs.Bid.Err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(obs.Bid.Err, "failed to observe Bid"))
	} else if bid, err := EncodeValueInt192(obs.Bid.Val); err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(err, "failed to observe Bid; encoding failed"))
	} else {
		p.Bid = bid
	}

	if obs.Ask.Err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(obs.Ask.Err, "failed to observe Ask"))
	} else if ask, err := EncodeValueInt192(obs.Ask.Val); err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(err, "failed to observe Ask; encoding failed"))
	} else {
		p.Ask = ask
	}

	if obs.BenchmarkPrice.Err == nil && obs.Bid.Err == nil && obs.Ask.Err == nil {
		p.PricesValid = true
	}

	if obs.CurrentBlockNum.Err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(obs.CurrentBlockNum.Err, "failed to observe CurrentBlockNum"))
	} else {
		p.CurrentBlockNum = obs.CurrentBlockNum.Val
	}

	if obs.CurrentBlockHash.Err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(obs.CurrentBlockHash.Err, "failed to observe CurrentBlockHash"))
	} else {
		p.CurrentBlockHash = obs.CurrentBlockHash.Val
	}

	if obs.CurrentBlockTimestamp.Err != nil {
		obsErrors = append(obsErrors, pkgerrors.Wrap(obs.CurrentBlockTimestamp.Err, "failed to observe CurrentBlockTimestamp"))
	} else {
		p.CurrentBlockTimestamp = obs.CurrentBlockTimestamp.Val
	}

	if obs.CurrentBlockNum.Err == nil && obs.CurrentBlockHash.Err == nil && obs.CurrentBlockTimestamp.Err == nil {
		p.CurrentBlockValid = true
	}

	if len(obsErrors) > 0 {
		rp.logger.Warnw(fmt.Sprintf("Observe failed %d/6 observations", len(obsErrors)), "err", errors.Join(obsErrors...))
	}

	return proto.Marshal(&p)
}

type ParsedAttributedObservation struct {
	Timestamp uint32
	Observer  commontypes.OracleID

	BenchmarkPrice *big.Int
	Bid            *big.Int
	Ask            *big.Int
	// All three prices must be valid, or none are (they all should come from one API query and hold invariant bid <= bm <= ask)
	PricesValid bool

	CurrentBlockNum       int64 // inclusive; current block
	CurrentBlockHash      []byte
	CurrentBlockTimestamp uint64
	// All three block observations must be valid, or none are (they all come from the same block)
	CurrentBlockValid bool

	// MaxFinalizedBlockNumber comes from previous report when present and is
	// only observed from mercury server when previous report is nil
	//
	// MaxFinalizedBlockNumber will be -1 if there is none
	MaxFinalizedBlockNumber      int64
	MaxFinalizedBlockNumberValid bool
}

func parseAttributedObservation(ao ocrtypes.AttributedObservation) (pao ParsedAttributedObservation, err error) {
	var obs MercuryObservationProto
	if err = proto.Unmarshal(ao.Observation, &obs); err != nil {
		return ParsedAttributedObservation{}, pkgerrors.Errorf("attributed observation cannot be unmarshaled: %s", err)
	}

	pao.Timestamp = obs.Timestamp
	pao.Observer = ao.Observer

	if obs.PricesValid {
		pao.BenchmarkPrice, err = DecodeValueInt192(obs.BenchmarkPrice)
		if err != nil {
			return ParsedAttributedObservation{}, pkgerrors.Errorf("benchmarkPrice cannot be converted to big.Int: %s", err)
		}
		pao.Bid, err = DecodeValueInt192(obs.Bid)
		if err != nil {
			return ParsedAttributedObservation{}, pkgerrors.Errorf("bid cannot be converted to big.Int: %s", err)
		}
		pao.Ask, err = DecodeValueInt192(obs.Ask)
		if err != nil {
			return ParsedAttributedObservation{}, pkgerrors.Errorf("ask cannot be converted to big.Int: %s", err)
		}
		pao.PricesValid = true
	}

	if obs.CurrentBlockValid {
		if len(obs.CurrentBlockHash) != evmHashLen {
			return ParsedAttributedObservation{}, pkgerrors.Errorf("wrong len for hash: %d (expected: %d)", len(obs.CurrentBlockHash), evmHashLen)
		}
		pao.CurrentBlockHash = obs.CurrentBlockHash
		if obs.CurrentBlockNum < 0 {
			return ParsedAttributedObservation{}, pkgerrors.Errorf("negative block number: %d", obs.CurrentBlockNum)
		}
		pao.CurrentBlockNum = obs.CurrentBlockNum
		pao.CurrentBlockTimestamp = obs.CurrentBlockTimestamp
		pao.CurrentBlockValid = true
	}

	if obs.MaxFinalizedBlockNumberValid {
		pao.MaxFinalizedBlockNumber = obs.MaxFinalizedBlockNumber
		pao.MaxFinalizedBlockNumberValid = true
	}

	return pao, nil
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

func (rp *reportingPlugin) Report(repts types.ReportTimestamp, previousReport types.Report, aos []types.AttributedObservation) (shouldReport bool, report types.Report, err error) {
	paos := parseAttributedObservations(rp.logger, aos)

	// By assumption, we have at most f malicious oracles, so there should be at least f+1 valid paos
	if !(rp.f+1 <= len(paos)) {
		return false, nil, pkgerrors.Errorf("only received %v valid attributed observations, but need at least f+1 (%v)", len(paos), rp.f+1)
	}

	var validFromBlockNum int64
	if previousReport != nil {
		var currentBlockNum int64
		currentBlockNum, err = rp.reportCodec.CurrentBlockNumFromReport(previousReport)
		if err != nil {
			return false, nil, err
		}
		validFromBlockNum = currentBlockNum + 1
	} else {
		var maxFinalizedBlockNumber int64
		maxFinalizedBlockNumber, err = GetConsensusMaxFinalizedBlockNum(paos, rp.f)
		if err != nil {
			return false, nil, err
		}
		validFromBlockNum = maxFinalizedBlockNumber + 1
	}
	should, err := rp.shouldReport(validFromBlockNum, repts, paos)
	if err != nil {
		return false, nil, err
	}
	if !should {
		return false, nil, nil
	}
	report, err = rp.reportCodec.BuildReport(paos, rp.f, validFromBlockNum)
	if err != nil {
		rp.logger.Debugw("failed to BuildReport", "paos", paos, "f", rp.f, "validFromBlockNum", validFromBlockNum, "repts", repts)
		return false, nil, err
	}
	if !(len(report) <= rp.maxReportLength) {
		return false, nil, pkgerrors.Errorf("report with len %d violates MaxReportLength limit set by ReportCodec (%d)", len(report), rp.maxReportLength)
	} else if len(report) == 0 {
		return false, nil, errors.New("report may not have zero length (invariant violation)")
	}

	return true, report, nil
}

func (rp *reportingPlugin) shouldReport(validFromBlockNum int64, repts types.ReportTimestamp, paos []ParsedAttributedObservation) (bool, error) {
	if !(rp.f+1 <= len(paos)) {
		return false, pkgerrors.Errorf("only received %v valid attributed observations, but need at least f+1 (%v)", len(paos), rp.f+1)
	}

	if err := errors.Join(
		rp.checkBenchmarkPrice(paos),
		rp.checkBid(paos),
		rp.checkAsk(paos),
		rp.checkCurrentBlock(paos, validFromBlockNum),
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
	return ValidateBenchmarkPrice(paos, rp.f, rp.onchainConfig.Min, rp.onchainConfig.Max)
}

func (rp *reportingPlugin) checkBid(paos []ParsedAttributedObservation) error {
	return ValidateBid(paos, rp.f, rp.onchainConfig.Min, rp.onchainConfig.Max)
}

func (rp *reportingPlugin) checkAsk(paos []ParsedAttributedObservation) error {
	return ValidateAsk(paos, rp.f, rp.onchainConfig.Min, rp.onchainConfig.Max)
}

func (rp *reportingPlugin) checkCurrentBlock(paos []ParsedAttributedObservation, validFromBlockNum int64) error {
	return ValidateCurrentBlock(paos, rp.f, validFromBlockNum)
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

	rp.logger.Debugw("ShouldAcceptFinalizedReport() = true",
		"reportEpochRound", reportEpochRound,
		"latestAcceptedEpochRound", rp.latestAcceptedEpochRound,
	)

	rp.latestAcceptedEpochRound = reportEpochRound

	return true, nil
}

func (rp *reportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, repts types.ReportTimestamp, report types.Report) (bool, error) {
	return true, nil
}

func (rp *reportingPlugin) Close() error {
	return nil
}

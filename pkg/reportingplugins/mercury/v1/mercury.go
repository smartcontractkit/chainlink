package mercury_v1

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	pkgerrors "github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
)

// Mercury-specific reporting plugin, based off of median:
// https://github.com/smartcontractkit/offchain-reporting/blob/master/lib/offchainreporting2/reportingplugin/median/median.go

type Observation struct {
	BenchmarkPrice        mercury.ObsResult[*big.Int]
	Bid                   mercury.ObsResult[*big.Int]
	Ask                   mercury.ObsResult[*big.Int]
	CurrentBlockNum       mercury.ObsResult[int64]
	CurrentBlockHash      mercury.ObsResult[[]byte]
	CurrentBlockTimestamp mercury.ObsResult[uint64]
	// MaxFinalizedBlockNumber comes from previous report when present and is
	// only observed from mercury server when previous report is nil
	MaxFinalizedBlockNumber mercury.ObsResult[int64]
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
	mercury.ByteWidthInt192 + // benchmarkPrice
	mercury.ByteWidthInt192 + // bid
	mercury.ByteWidthInt192 + // ask
	1 + // pricesValid
	8 + // currentBlockNum
	32 + // currentBlockHash
	8 + // currentBlockTimestamp
	1 + // currentBlockValid
	8 + // maxFinalizedBlockNumber
	1 + // maxFinalizedBlockNumberValid
	32 // [> overapprox. of protobuf overhead <]

type Factory struct {
	dataSource         DataSource
	logger             logger.Logger
	onchainConfigCodec mercury.OnchainConfigCodec
	reportCodec        ReportCodec
}

func NewFactory(ds DataSource, lggr logger.Logger, occ mercury.OnchainConfigCodec, rc ReportCodec) Factory {
	return Factory{ds, lggr, occ, rc}
}

func (fac Factory) NewMercuryPlugin(configuration ocr3types.MercuryPluginConfig) (ocr3types.MercuryPlugin, ocr3types.MercuryPluginInfo, error) {
	offchainConfig, err := mercury.DecodeOffchainConfig(configuration.OffchainConfig)
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
		mercury.EpochRound{},
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
	offchainConfig mercury.OffchainConfig
	onchainConfig  mercury.OnchainConfig
	dataSource     DataSource
	logger         logger.Logger
	reportCodec    ReportCodec

	configDigest             ocrtypes.ConfigDigest
	f                        int
	latestAcceptedEpochRound mercury.EpochRound
	maxReportLength          int
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

	var bpErr, bidErr, askErr error
	if obs.BenchmarkPrice.Err != nil {
		bpErr = pkgerrors.Wrap(obs.BenchmarkPrice.Err, "failed to observe BenchmarkPrice")
		obsErrors = append(obsErrors, bpErr)
	} else if benchmarkPrice, err := mercury.EncodeValueInt192(obs.BenchmarkPrice.Val); err != nil {
		bpErr = pkgerrors.Wrap(err, "failed to observe BenchmarkPrice; encoding failed")
		obsErrors = append(obsErrors, bpErr)
	} else {
		p.BenchmarkPrice = benchmarkPrice
	}

	if obs.Bid.Err != nil {
		bidErr = pkgerrors.Wrap(obs.Bid.Err, "failed to observe Bid")
		obsErrors = append(obsErrors, bidErr)
	} else if bid, err := mercury.EncodeValueInt192(obs.Bid.Val); err != nil {
		bidErr = pkgerrors.Wrap(err, "failed to observe Bid; encoding failed")
		obsErrors = append(obsErrors, bidErr)
	} else {
		p.Bid = bid
	}

	if obs.Ask.Err != nil {
		askErr = pkgerrors.Wrap(obs.Ask.Err, "failed to observe Ask")
		obsErrors = append(obsErrors, askErr)
	} else if ask, err := mercury.EncodeValueInt192(obs.Ask.Val); err != nil {
		askErr = pkgerrors.Wrap(err, "failed to observe Ask; encoding failed")
		obsErrors = append(obsErrors, askErr)
	} else {
		p.Ask = ask
	}

	if bpErr == nil && bidErr == nil && askErr == nil {
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

func parseAttributedObservation(ao ocrtypes.AttributedObservation) (PAO, error) {
	var pao parsedAttributedObservation
	var obs MercuryObservationProto
	if err := proto.Unmarshal(ao.Observation, &obs); err != nil {
		return parsedAttributedObservation{}, pkgerrors.Errorf("attributed observation cannot be unmarshaled: %s", err)
	}

	pao.Timestamp = obs.Timestamp
	pao.Observer = ao.Observer

	if obs.PricesValid {
		var err error
		pao.BenchmarkPrice, err = mercury.DecodeValueInt192(obs.BenchmarkPrice)
		if err != nil {
			return parsedAttributedObservation{}, pkgerrors.Errorf("benchmarkPrice cannot be converted to big.Int: %s", err)
		}
		pao.Bid, err = mercury.DecodeValueInt192(obs.Bid)
		if err != nil {
			return parsedAttributedObservation{}, pkgerrors.Errorf("bid cannot be converted to big.Int: %s", err)
		}
		pao.Ask, err = mercury.DecodeValueInt192(obs.Ask)
		if err != nil {
			return parsedAttributedObservation{}, pkgerrors.Errorf("ask cannot be converted to big.Int: %s", err)
		}
		pao.PricesValid = true
	}

	if obs.CurrentBlockValid {
		if len(obs.CurrentBlockHash) != mercury.EvmHashLen {
			return parsedAttributedObservation{}, pkgerrors.Errorf("wrong len for hash: %d (expected: %d)", len(obs.CurrentBlockHash), mercury.EvmHashLen)
		}
		pao.CurrentBlockHash = obs.CurrentBlockHash
		if obs.CurrentBlockNum < 0 {
			return parsedAttributedObservation{}, pkgerrors.Errorf("negative block number: %d", obs.CurrentBlockNum)
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

func parseAttributedObservations(lggr logger.Logger, aos []ocrtypes.AttributedObservation) []PAO {
	paos := make([]PAO, 0, len(aos))
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

	if len(paos) == 0 {
		return false, nil, errors.New("got zero valid attributed observations")
	}

	// By assumption, we have at most f malicious oracles, so there should be at least f+1 valid paos
	if !(rp.f+1 <= len(paos)) {
		return false, nil, pkgerrors.Errorf("only received %v valid attributed observations, but need at least f+1 (%v)", len(paos), rp.f+1)
	}

	rf, err := rp.buildReportFields(previousReport, paos)
	if err != nil {
		rp.logger.Errorw("failed to build report fields", "paos", paos, "f", rp.f, "reportFields", rf, "repts", repts, "err", err)
		return false, nil, err
	}

	if rf.CurrentBlockNum < rf.ValidFromBlockNum {
		rp.logger.Debugw("shouldReport: no (overlap)", "currentBlockNum", rf.CurrentBlockNum, "validFromBlockNum", rf.ValidFromBlockNum, "repts", repts)
		return false, nil, nil
	}

	if err = rp.validateReport(rf); err != nil {
		rp.logger.Errorw("shouldReport: no (validation error)", "reportFields", rf, "err", err, "repts", repts, "paos", paos)
		return false, nil, err
	}
	rp.logger.Debugw("shouldReport: yes",
		"timestamp", repts,
	)

	report, err = rp.reportCodec.BuildReport(rf)
	if err != nil {
		rp.logger.Debugw("failed to BuildReport", "paos", paos, "f", rp.f, "reportFields", rf, "repts", repts)
		return false, nil, err
	}
	if !(len(report) <= rp.maxReportLength) {
		return false, nil, pkgerrors.Errorf("report with len %d violates MaxReportLength limit set by ReportCodec (%d)", len(report), rp.maxReportLength)
	} else if len(report) == 0 {
		return false, nil, errors.New("report may not have zero length (invariant violation)")
	}

	return true, report, nil
}

func (rp *reportingPlugin) buildReportFields(previousReport types.Report, paos []PAO) (rf ReportFields, merr error) {
	var err error
	if previousReport != nil {
		var maxFinalizedBlockNumber int64
		maxFinalizedBlockNumber, err = rp.reportCodec.CurrentBlockNumFromReport(previousReport)
		if err != nil {
			merr = errors.Join(merr, err)
		} else {
			rf.ValidFromBlockNum = maxFinalizedBlockNumber + 1
		}
	} else {
		var maxFinalizedBlockNumber int64
		maxFinalizedBlockNumber, err = GetConsensusMaxFinalizedBlockNum(paos, rp.f)
		if err != nil {
			merr = errors.Join(merr, err)
		} else {
			rf.ValidFromBlockNum = maxFinalizedBlockNumber + 1
		}
	}

	mPaos := convert(paos)

	rf.Timestamp = mercury.GetConsensusTimestamp(mPaos)

	rf.BenchmarkPrice, err = mercury.GetConsensusBenchmarkPrice(mPaos, rp.f)
	merr = errors.Join(merr, pkgerrors.Wrap(err, "GetConsensusBenchmarkPrice failed"))

	rf.Bid, err = mercury.GetConsensusBid(convertBid(paos), rp.f)
	merr = errors.Join(merr, pkgerrors.Wrap(err, "GetConsensusBid failed"))

	rf.Ask, err = mercury.GetConsensusAsk(convertAsk(paos), rp.f)
	merr = errors.Join(merr, pkgerrors.Wrap(err, "GetConsensusAsk failed"))

	rf.CurrentBlockHash, rf.CurrentBlockNum, rf.CurrentBlockTimestamp, err = GetConsensusCurrentBlock(paos, rp.f)
	merr = errors.Join(merr, pkgerrors.Wrap(err, "GetConsensusCurrentBlock failed"))

	return rf, merr
}

func (rp *reportingPlugin) validateReport(rf ReportFields) error {
	return errors.Join(
		mercury.ValidateBetween("median benchmark price", rf.BenchmarkPrice, rp.onchainConfig.Min, rp.onchainConfig.Max),
		mercury.ValidateBetween("median bid", rf.Bid, rp.onchainConfig.Min, rp.onchainConfig.Max),
		mercury.ValidateBetween("median ask", rf.Ask, rp.onchainConfig.Min, rp.onchainConfig.Max),
		ValidateCurrentBlock(rf),
	)
}

func (rp *reportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, repts types.ReportTimestamp, report types.Report) (bool, error) {
	return true, nil
}

func (rp *reportingPlugin) Close() error {
	return nil
}

// convert funcs are necessary because go is not smart enough to cast
// []interface1 to []interface2 even if interface1 is a superset of interface2
func convert(pao []PAO) (ret []mercury.PAO) {
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}
func convertBid(pao []PAO) (ret []mercury.PAOBid) {
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}
func convertAsk(pao []PAO) (ret []mercury.PAOAsk) {
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}

package mercury_v2

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	pkgerrors "github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
)

type Observation struct {
	BenchmarkPrice mercury.ObsResult[*big.Int]

	MaxFinalizedTimestamp mercury.ObsResult[int64]

	LinkPrice   mercury.ObsResult[*big.Int]
	NativePrice mercury.ObsResult[*big.Int]
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
	Observe(ctx context.Context, repts ocrtypes.ReportTimestamp, fetchMaxFinalizedTimestamp bool) (Observation, error)
}

var _ ocr3types.MercuryPluginFactory = Factory{}

const maxObservationLength = 32 + // feedID
	4 + // timestamp
	mercury.ByteWidthInt192 + // benchmarkPrice
	4 + // validFromTimestamp
	mercury.ByteWidthInt192 + // linkFee
	mercury.ByteWidthInt192 + // nativeFee
	16 /* overapprox. of protobuf overhead */

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
		new(big.Int),
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
	latestAcceptedMedian     *big.Int
	maxReportLength          int
}

var MissingPrice = big.NewInt(-1)

func (rp *reportingPlugin) Observation(ctx context.Context, repts ocrtypes.ReportTimestamp, previousReport ocrtypes.Report) (ocrtypes.Observation, error) {
	obs, err := rp.dataSource.Observe(ctx, repts, previousReport == nil)
	if err != nil {
		return nil, pkgerrors.Errorf("DataSource.Observe returned an error: %s", err)
	}

	observationTimestamp := time.Now()
	if observationTimestamp.Unix() > math.MaxUint32 {
		return nil, fmt.Errorf("current unix epoch %d exceeds max uint32", observationTimestamp.Unix())
	}
	p := MercuryObservationProto{Timestamp: uint32(observationTimestamp.Unix())}
	var obsErrors []error

	var bpErr error
	if obs.BenchmarkPrice.Err != nil {
		bpErr = pkgerrors.Wrap(obs.BenchmarkPrice.Err, "failed to observe BenchmarkPrice")
		obsErrors = append(obsErrors, bpErr)
	} else if benchmarkPrice, err := mercury.EncodeValueInt192(obs.BenchmarkPrice.Val); err != nil {
		bpErr = pkgerrors.Wrapf(err, "failed to encode BenchmarkPrice; val=%s", obs.BenchmarkPrice.Val)
		obsErrors = append(obsErrors, bpErr)
	} else {
		p.BenchmarkPrice = benchmarkPrice
		p.PricesValid = true
	}

	var maxFinalizedTimestampErr error
	if obs.MaxFinalizedTimestamp.Err != nil {
		maxFinalizedTimestampErr = pkgerrors.Wrap(obs.MaxFinalizedTimestamp.Err, "failed to observe MaxFinalizedTimestamp")
		obsErrors = append(obsErrors, maxFinalizedTimestampErr)
	} else {
		p.MaxFinalizedTimestamp = obs.MaxFinalizedTimestamp.Val
		p.MaxFinalizedTimestampValid = true
	}

	var linkErr error
	if obs.LinkPrice.Err != nil {
		linkErr = pkgerrors.Wrap(obs.LinkPrice.Err, "failed to observe LINK price")
		obsErrors = append(obsErrors, linkErr)
	} else if obs.LinkPrice.Val.Cmp(MissingPrice) <= 0 {
		p.LinkFee = mercury.MaxInt192Enc
	} else {
		linkFee := mercury.CalculateFee(obs.LinkPrice.Val, rp.offchainConfig.BaseUSDFee)
		if linkFeeEncoded, err := mercury.EncodeValueInt192(linkFee); err != nil {
			linkErr = pkgerrors.Wrapf(err, "failed to encode LINK fee; val=%s", linkFee)
			obsErrors = append(obsErrors, linkErr)
		} else {
			p.LinkFee = linkFeeEncoded
		}
	}

	if linkErr == nil {
		p.LinkFeeValid = true
	}

	var nativeErr error
	if obs.NativePrice.Err != nil {
		nativeErr = pkgerrors.Wrap(obs.NativePrice.Err, "failed to observe native price")
		obsErrors = append(obsErrors, nativeErr)
	} else if obs.NativePrice.Val.Cmp(MissingPrice) <= 0 {
		p.NativeFee = mercury.MaxInt192Enc
	} else {
		nativeFee := mercury.CalculateFee(obs.NativePrice.Val, rp.offchainConfig.BaseUSDFee)
		if nativeFeeEncoded, err := mercury.EncodeValueInt192(nativeFee); err != nil {
			nativeErr = pkgerrors.Wrapf(err, "failed to encode native fee; val=%s", nativeFee)
			obsErrors = append(obsErrors, nativeErr)
		} else {
			p.NativeFee = nativeFeeEncoded
		}
	}

	if nativeErr == nil {
		p.NativeFeeValid = true
	}

	if len(obsErrors) > 0 {
		rp.logger.Warnw(fmt.Sprintf("Observe failed %d/4 observations", len(obsErrors)), "err", errors.Join(obsErrors...))
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
		pao.PricesValid = true
	}

	if obs.MaxFinalizedTimestampValid {
		pao.MaxFinalizedTimestamp = obs.MaxFinalizedTimestamp
		pao.MaxFinalizedTimestampValid = true
	}

	if obs.LinkFeeValid {
		var err error
		pao.LinkFee, err = mercury.DecodeValueInt192(obs.LinkFee)
		if err != nil {
			return parsedAttributedObservation{}, pkgerrors.Errorf("link price cannot be converted to big.Int: %s", err)
		}
		pao.LinkFeeValid = true
	}
	if obs.NativeFeeValid {
		var err error
		pao.NativeFee, err = mercury.DecodeValueInt192(obs.NativeFee)
		if err != nil {
			return parsedAttributedObservation{}, pkgerrors.Errorf("native price cannot be converted to big.Int: %s", err)
		}
		pao.NativeFeeValid = true
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

func (rp *reportingPlugin) Report(repts ocrtypes.ReportTimestamp, previousReport ocrtypes.Report, aos []ocrtypes.AttributedObservation) (shouldReport bool, report ocrtypes.Report, err error) {
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

	if rf.Timestamp < rf.ValidFromTimestamp {
		rp.logger.Debugw("shouldReport: no (overlap)", "observationTimestamp", rf.Timestamp, "validFromTimestamp", rf.ValidFromTimestamp, "repts", repts)
		return false, nil, nil
	}

	if err = rp.validateReport(rf); err != nil {
		rp.logger.Errorw("shouldReport: no (validation error)", "reportFields", rf, "err", err, "repts", repts, "paos", paos)
		return false, nil, err
	}
	rp.logger.Debugw("shouldReport: yes", "repts", repts)

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

func (rp *reportingPlugin) buildReportFields(previousReport ocrtypes.Report, paos []PAO) (rf ReportFields, merr error) {
	mPaos := convert(paos)
	rf.Timestamp = mercury.GetConsensusTimestamp(mPaos)

	var err error
	if previousReport != nil {
		var maxFinalizedTimestamp uint32
		maxFinalizedTimestamp, err = rp.reportCodec.ObservationTimestampFromReport(previousReport)
		merr = errors.Join(merr, err)
		rf.ValidFromTimestamp = maxFinalizedTimestamp + 1
	} else {
		var maxFinalizedTimestamp int64
		maxFinalizedTimestamp, err = mercury.GetConsensusMaxFinalizedTimestamp(convertMaxFinalizedTimestamp(paos), rp.f)
		if err != nil {
			merr = errors.Join(merr, err)
		} else if maxFinalizedTimestamp < 0 {
			// no previous observation timestamp available, e.g. in case of new
			// feed; use current timestamp as start of range
			rf.ValidFromTimestamp = rf.Timestamp
		} else if maxFinalizedTimestamp+1 > math.MaxUint32 {
			merr = errors.Join(err, fmt.Errorf("maxFinalizedTimestamp is too large, got: %d", maxFinalizedTimestamp))
		} else {
			rf.ValidFromTimestamp = uint32(maxFinalizedTimestamp + 1)
		}
	}

	rf.BenchmarkPrice, err = mercury.GetConsensusBenchmarkPrice(mPaos, rp.f)
	merr = errors.Join(merr, pkgerrors.Wrap(err, "GetConsensusBenchmarkPrice failed"))

	rf.LinkFee, err = mercury.GetConsensusLinkFee(convertLinkFee(paos), rp.f)
	if err != nil {
		// It is better to generate a report that will validate for free,
		// rather than no report at all, if we cannot come to consensus on a
		// valid fee.
		rp.logger.Errorw("Cannot come to consensus on LINK fee, falling back to 0", "err", err, "paos", paos)
		rf.LinkFee = big.NewInt(0)
	}

	rf.NativeFee, err = mercury.GetConsensusNativeFee(convertNativeFee(paos), rp.f)
	if err != nil {
		// It is better to generate a report that will validate for free,
		// rather than no report at all, if we cannot come to consensus on a
		// valid fee.
		rp.logger.Errorw("Cannot come to consensus on Native fee, falling back to 0", "err", err, "paos", paos)
		rf.NativeFee = big.NewInt(0)
	}

	if int64(rf.Timestamp)+int64(rp.offchainConfig.ExpirationWindow) > math.MaxUint32 {
		merr = errors.Join(merr, fmt.Errorf("timestamp %d + expiration window %d overflows uint32", rf.Timestamp, rp.offchainConfig.ExpirationWindow))
	} else {
		rf.ExpiresAt = rf.Timestamp + rp.offchainConfig.ExpirationWindow
	}

	return rf, merr
}

func (rp *reportingPlugin) validateReport(rf ReportFields) error {
	return errors.Join(
		mercury.ValidateBetween("median benchmark price", rf.BenchmarkPrice, rp.onchainConfig.Min, rp.onchainConfig.Max),
		mercury.ValidateFee("median link fee", rf.LinkFee),
		mercury.ValidateFee("median native fee", rf.NativeFee),
		mercury.ValidateValidFromTimestamp(rf.Timestamp, rf.ValidFromTimestamp),
		mercury.ValidateExpiresAt(rf.Timestamp, rf.ExpiresAt),
	)
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
func convertMaxFinalizedTimestamp(pao []PAO) (ret []mercury.PAOMaxFinalizedTimestamp) {
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}
func convertLinkFee(pao []PAO) (ret []mercury.PAOLinkFee) {
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}
func convertNativeFee(pao []PAO) (ret []mercury.PAONativeFee) {
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}

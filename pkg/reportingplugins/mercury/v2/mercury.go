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

	MaxFinalizedTimestamp mercury.ObsResult[uint32]

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
		bpErr = pkgerrors.Wrap(err, "failed to observe BenchmarkPrice; encoding failed")
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
	} else {
		linkFee := mercury.CalculateFee(obs.LinkPrice.Val, rp.offchainConfig.BaseUSDFeeCents)
		if linkFeeEncoded, err := mercury.EncodeValueInt192(linkFee); err != nil {
			linkErr = pkgerrors.Wrap(err, "failed to observe LINK price; encoding failed")
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
	} else {
		nativeFee := mercury.CalculateFee(obs.NativePrice.Val, rp.offchainConfig.BaseUSDFeeCents)
		if nativeFeeEncoded, err := mercury.EncodeValueInt192(nativeFee); err != nil {
			nativeErr = pkgerrors.Wrap(err, "failed to observe native price; encoding failed")
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

func parseAttributedObservation(ao ocrtypes.AttributedObservation) (ParsedAttributedObservation, error) {
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

func (rp *reportingPlugin) Report(repts ocrtypes.ReportTimestamp, previousReport ocrtypes.Report, aos []ocrtypes.AttributedObservation) (shouldReport bool, report ocrtypes.Report, err error) {
	paos := parseAttributedObservations(rp.logger, aos)

	// By assumption, we have at most f malicious oracles, so there should be at least f+1 valid paos
	if !(rp.f+1 <= len(paos)) {
		return false, nil, pkgerrors.Errorf("only received %v valid attributed observations, but need at least f+1 (%v)", len(paos), rp.f+1)
	}

	observationTimestamp := mercury.GetConsensusTimestamp(Convert(paos))

	var validFromTimestamp uint32
	if previousReport != nil {
		validFromTimestamp, err = rp.reportCodec.ObservationTimestampFromReport(previousReport)
		if err != nil {
			return false, nil, err
		}
	} else {
		validFromTimestamp, err = mercury.GetConsensusMaxFinalizedTimestamp(Convert(paos), rp.f)
		if err != nil {
			return false, nil, err
		}

		// no previous observation timestamp available, e.g. in case of new feed
		if validFromTimestamp == 0 {
			validFromTimestamp = observationTimestamp
		}
	}

	should, err := rp.shouldReport(paos, observationTimestamp, validFromTimestamp)
	if err != nil || !should {
		rp.logger.Debugw("shouldReport: no", "err", err)
		return false, nil, err
	}

	rp.logger.Debugw("shouldReport: yes",
		"timestamp", repts,
	)

	expiresAt := observationTimestamp + rp.offchainConfig.ExpirationWindow

	report, err = rp.reportCodec.BuildReport(paos, rp.f, validFromTimestamp, expiresAt)
	if err != nil {
		rp.logger.Debugw("failed to BuildReport", "paos", paos, "f", rp.f, "validFromTimestamp", validFromTimestamp, "repts", repts)
		return false, nil, err
	}

	if !(len(report) <= rp.maxReportLength) {
		return false, nil, pkgerrors.Errorf("report with len %d violates MaxReportLength limit set by ReportCodec (%d)", len(report), rp.maxReportLength)
	} else if len(report) == 0 {
		return false, nil, errors.New("report may not have zero length (invariant violation)")
	}

	return true, report, nil
}

func (rp *reportingPlugin) shouldReport(paos []ParsedAttributedObservation, observationTimestamp uint32, validFromTimestamp uint32) (bool, error) {
	if err := errors.Join(
		rp.checkBenchmarkPrice(paos),
		rp.checkValidFromTimestamp(observationTimestamp, validFromTimestamp),
		rp.checkExpiresAt(observationTimestamp, rp.offchainConfig.ExpirationWindow),
	); err != nil {
		return false, err
	}

	return true, nil
}

func (rp *reportingPlugin) checkBenchmarkPrice(paos []ParsedAttributedObservation) error {
	mPaos := Convert(paos)
	return mercury.ValidateBenchmarkPrice(mPaos, rp.f, rp.onchainConfig.Min, rp.onchainConfig.Max)
}

func (rp *reportingPlugin) checkValidFromTimestamp(observationTimestamp uint32, validFromTimestamp uint32) error {
	return mercury.ValidateValidFromTimestamp(observationTimestamp, validFromTimestamp)
}

func (rp *reportingPlugin) checkExpiresAt(observationTimestamp uint32, expirationWindow uint32) error {
	return mercury.ValidateExpiresAt(observationTimestamp, expirationWindow)
}

func (rp *reportingPlugin) ShouldAcceptFinalizedReport(ctx context.Context, repts ocrtypes.ReportTimestamp, report ocrtypes.Report) (bool, error) {
	reportEpochRound := mercury.EpochRound{Epoch: repts.Epoch, Round: repts.Round}
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

func (rp *reportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, repts ocrtypes.ReportTimestamp, report ocrtypes.Report) (bool, error) {
	return true, nil
}

func (rp *reportingPlugin) Close() error {
	return nil
}

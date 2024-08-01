package median

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"time"

	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

const onchainConfigVersion = 1
const onchainConfigEncodedLength = 1 + byteWidth + byteWidth

type OnchainConfig struct {
	Min *big.Int
	Max *big.Int
}

type OnchainConfigCodec interface {
	Encode(OnchainConfig) ([]byte, error)
	Decode([]byte) (OnchainConfig, error)
}

var _ OnchainConfigCodec = StandardOnchainConfigCodec{}

// StandardOnchainConfigCodec provides a standard implementation of OnchainConfigCodec.
// This is the implementation used by the EVM and Solana integrations.
//
// An encoded onchain config is expected to be in the format
// <version><min><max>
// where version is a uint8 and min and max are in the format
// returned by EncodeValue.
type StandardOnchainConfigCodec struct{}

func (StandardOnchainConfigCodec) Decode(b []byte) (OnchainConfig, error) {
	if len(b) != onchainConfigEncodedLength {
		return OnchainConfig{}, fmt.Errorf("unexpected length of OnchainConfig, expected %v, got %v", onchainConfigEncodedLength, len(b))
	}

	if b[0] != onchainConfigVersion {
		return OnchainConfig{}, fmt.Errorf("unexpected version of OnchainConfig, expected %v, got %v", onchainConfigVersion, b[0])
	}

	min, err := DecodeValue(b[1 : 1+byteWidth])
	if err != nil {
		return OnchainConfig{}, err
	}
	max, err := DecodeValue(b[1+byteWidth:])
	if err != nil {
		return OnchainConfig{}, err
	}

	if !(min.Cmp(max) <= 0) {
		return OnchainConfig{}, fmt.Errorf("OnchainConfig min (%v) should not be greater than max(%v)", min, max)
	}

	return OnchainConfig{min, max}, nil
}

func (StandardOnchainConfigCodec) Encode(c OnchainConfig) ([]byte, error) {
	minBytes, err := EncodeValue(c.Min)
	if err != nil {
		return nil, err
	}
	maxBytes, err := EncodeValue(c.Max)
	if err != nil {
		return nil, err
	}
	result := make([]byte, 0, onchainConfigEncodedLength)
	result = append(result, onchainConfigVersion)
	result = append(result, minBytes...)
	result = append(result, maxBytes...)
	return result, nil
}

type OffchainConfig struct {
	// If AlphaReportInfinite is true, the deviation check parametrized by
	// AlphaReportPPB will never be satisfied.
	AlphaReportInfinite bool
	// AlphaReportPPB determines the relative deviation between the median (i.e.
	// answer) in the contract and the current median of observations (offchain)
	// at which a report should be issued. That is, a report is issued if
	// abs((offchainMedian - contractMedian)/contractMedian) >= alphaReport.
	AlphaReportPPB uint64 // PPB is parts-per-billion
	// If AlphaAcceptInfinite is true, the deviation check parametrized by
	// AlphaAcceptPPB will never be satisfied.
	AlphaAcceptInfinite bool
	// AlphaAcceptPPB determines the relative deviation between the median in a
	// newly generated report considered for transmission and the median of the
	// currently pending report. That is, a report is accepted for transmission
	// if abs((newMedian - pendingMedian)/pendingMedian) >= alphaAccept. If no
	// report is pending, this variable has no effect.
	AlphaAcceptPPB uint64 // PPB is parts-per-billion
	// DeltaC is the maximum age of the latest report in the contract. If the
	// maximum age is exceeded, a new report will be created by the report
	// generation protocol.
	DeltaC time.Duration
}

func DecodeOffchainConfig(b []byte) (OffchainConfig, error) {
	var configProto NumericalMedianConfigProto
	if err := proto.Unmarshal(b, &configProto); err != nil {
		return OffchainConfig{}, err
	}

	deltaC := time.Duration(configProto.GetDeltaCNanoseconds())
	if !(0 <= deltaC) {
		return OffchainConfig{}, fmt.Errorf("DeltaC (%v) must be non-negative", deltaC)
	}

	return OffchainConfig{
		configProto.GetAlphaReportInfinite(),
		configProto.GetAlphaReportPpb(),
		configProto.GetAlphaAcceptInfinite(),
		configProto.GetAlphaAcceptPpb(),
		time.Duration(configProto.GetDeltaCNanoseconds()),
	}, nil
}

func (c OffchainConfig) Encode() []byte {
	configProto := NumericalMedianConfigProto{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		c.AlphaReportInfinite,
		c.AlphaReportPPB,
		c.AlphaAcceptInfinite,
		c.AlphaAcceptPPB,
		uint64(c.DeltaC),
	}
	result, err := proto.Marshal(&configProto)
	if err != nil {
		// assertion
		panic(fmt.Sprintf("unexpected error while encoding Config: %v", err))
	}
	return result
}

type MedianContract interface {
	LatestTransmissionDetails(
		ctx context.Context,
	) (
		configDigest types.ConfigDigest,
		epoch uint32,
		round uint8,
		latestAnswer *big.Int,
		latestTimestamp time.Time,
		err error,
	)

	// LatestRoundRequested returns the configDigest, epoch, and round from the latest
	// RoundRequested event emitted by the contract. LatestRoundRequested may or may not
	// return a result if the latest such event was emitted in a block b such that
	// b.timestamp < tip.timestamp - lookback.
	//
	// If no event is found, LatestRoundRequested should return zero values, not an error.
	// An error should only be returned if an actual error occurred during execution,
	// e.g. because there was an error querying the blockchain or the database.
	//
	// As an optimization, this function may also return zero values, if no
	// RoundRequested event has been emitted after the latest NewTransmission event.
	LatestRoundRequested(
		ctx context.Context,
		lookback time.Duration,
	) (
		configDigest types.ConfigDigest,
		epoch uint32,
		round uint8,
		err error,
	)
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
	Observe(context.Context, types.ReportTimestamp) (*big.Int, error)
}

// All functions on ReportCodec should be pure and thread-safe.
// Be careful validating and parsing any data passed.
type ReportCodec interface {
	// Implementers may assume that there is at most one
	// ParsedAttributedObservation per observer, and that all observers are
	// valid. However, observation values, timestamps, etc... should all be
	// treated as untrusted.
	BuildReport([]ParsedAttributedObservation) (types.Report, error)

	// Gets the "median" (the n//2-th ranked element to be more precise where n
	// is the length of the list) observation from the report. The input to this
	// function should be an output of BuildReport in the benign case.
	// Nevertheless, make sure to treat the input to this function as untrusted.
	MedianFromReport(types.Report) (*big.Int, error)

	// Returns the maximum length of a report based on n, the number of oracles.
	// The output of BuildReport must respect this maximum length.
	MaxReportLength(n int) (int, error)
}

var _ types.ReportingPluginFactory = NumericalMedianFactory{}

const maxObservationLength = 4 /* timestamp */ +
	byteWidth /* observation */ +
	byteWidth /* juelsPerFeeCoin */ +
	byteWidth /* gasPriceSubunits */ +
	16 /* overapprox. of protobuf overhead */

type NumericalMedianFactory struct {
	ContractTransmitter       MedianContract
	DataSource                DataSource
	JuelsPerFeeCoinDataSource DataSource
	// The Observe() function of the following DataSource returns a non-zero value if the underlying
	// chain does not support reading tx.gasPrice during execution. This is useful e.g. on Starknet.
	// The returned price is expected to be in subunits of the coin used for gas. E.g. on chains that
	// use Ether for gas this would be denominated in Wei.
	GasPriceSubunitsDataSource DataSource
	// Set this to false unless you need GasPriceSubunits to be included in reports
	// for the chain you're targeting.
	// Be careful! Older versions of the ReportingPlugin will discard observations
	// made by newer versions of the ReportingPlugin with this value
	// set to true. This could lead to liveness failures. Only set this to true if all
	// oracles in the protocol instance are running the newer version of the
	// ReportingPlugin.
	IncludeGasPriceSubunitsInObservation bool
	Logger                               commontypes.Logger
	OnchainConfigCodec                   OnchainConfigCodec
	ReportCodec                          ReportCodec
}

func (fac NumericalMedianFactory) NewReportingPlugin(configuration types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {

	offchainConfig, err := DecodeOffchainConfig(configuration.OffchainConfig)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}

	onchainConfig, err := fac.OnchainConfigCodec.Decode(configuration.OnchainConfig)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}

	logger := loghelper.MakeRootLoggerWithContext(fac.Logger).MakeChild(commontypes.LogFields{
		"configDigest":    configuration.ConfigDigest,
		"reportingPlugin": "NumericalMedian",
	})

	maxReportLength, err := fac.ReportCodec.MaxReportLength(configuration.N)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}

	return &numericalMedian{
			offchainConfig,
			onchainConfig,
			fac.ContractTransmitter,
			fac.DataSource,
			fac.JuelsPerFeeCoinDataSource,
			fac.GasPriceSubunitsDataSource,
			fac.IncludeGasPriceSubunitsInObservation,
			logger,
			fac.ReportCodec,

			configuration.ConfigDigest,
			configuration.F,
			epochRound{},
			new(big.Int),
			maxReportLength,
		}, types.ReportingPluginInfo{
			"NumericalMedian",
			false,
			types.ReportingPluginLimits{
				0,
				maxObservationLength,
				maxReportLength,
			},
		}, nil
}

func Deviates(thresholdPPB uint64, old *big.Int, new *big.Int) bool {
	if old.Cmp(i(0)) == 0 {
		if new.Cmp(i(0)) == 0 { //nolint:gosimple
			return false // Both values are zero; no deviation
		}
		return true // Any deviation from 0 is significant
	}
	// ||new - old|| / ||old||, approximated by a float
	change := &big.Rat{}
	change.SetFrac(i(0).Sub(new, old), old)
	change.Abs(change)
	threshold := &big.Rat{}
	threshold.SetFrac(
		(&big.Int{}).SetUint64(thresholdPPB),
		(&big.Int{}).SetUint64(1e9),
	)
	return change.Cmp(threshold) >= 0
}

var _ types.ReportingPlugin = (*numericalMedian)(nil)

type numericalMedian struct {
	offchainConfig                       OffchainConfig
	onchainConfig                        OnchainConfig
	contractTransmitter                  MedianContract
	dataSource                           DataSource
	juelsPerFeeCoinDataSource            DataSource
	gasPriceSubunitsDataSource           DataSource
	includeGasPriceSubunitsInObservation bool
	logger                               loghelper.LoggerWithContext
	reportCodec                          ReportCodec

	configDigest             types.ConfigDigest
	f                        int
	latestAcceptedEpochRound epochRound
	latestAcceptedMedian     *big.Int
	maxReportLength          int
}

func (nm *numericalMedian) Query(ctx context.Context, repts types.ReportTimestamp) (types.Query, error) {
	return nil, nil
}

func (nm *numericalMedian) Observation(ctx context.Context, repts types.ReportTimestamp, query types.Query) (types.Observation, error) {
	if len(query) != 0 {
		return nil, fmt.Errorf("expected empty query")
	}

	observe := func(dataSource DataSource, name string) ([]byte, error) {
		value, err := dataSource.Observe(ctx, repts)

		if err != nil {
			return nil, fmt.Errorf("%v.Observe returned an error: %w", name, err)
		}
		if value == nil {
			return nil, fmt.Errorf("%v.Observe returned unexpected nil big.Int", name)
		}
		encoded, err := EncodeValue(value)
		if err != nil {
			return nil, fmt.Errorf("failed to encode output of %v.Observe : %w", name, err)
		}
		return encoded, nil
	}
	var subs subprocesses.Subprocesses
	var value, juelsPerFeeCoin, gasPriceSubunits []byte
	var valueErr, juelsPerFeeCoinErr, gasPriceSubunitsErr error
	subs.Go(func() {
		value, valueErr = observe(nm.dataSource, "DataSource")
	})
	subs.Go(func() {
		juelsPerFeeCoin, juelsPerFeeCoinErr = observe(nm.juelsPerFeeCoinDataSource, "JuelsPerFeeCoinDataSource")
	})
	subs.Go(func() {
		gasPriceSubunits, gasPriceSubunitsErr = observe(nm.gasPriceSubunitsDataSource, "GasPriceSubunitsDataSource")
	})
	subs.Wait()

	err := multierr.Combine(valueErr, juelsPerFeeCoinErr, gasPriceSubunitsErr)
	if err != nil {
		return nil, fmt.Errorf("error in Observation: %w", err)
	}

	if !nm.includeGasPriceSubunitsInObservation {
		gasPriceSubunits = nil
	}

	return proto.Marshal(&NumericalMedianObservationProto{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		uint32(time.Now().Unix()),
		value,
		juelsPerFeeCoin,
		gasPriceSubunits,
	})
}

type ParsedAttributedObservation struct {
	Timestamp        uint32
	Value            *big.Int
	JuelsPerFeeCoin  *big.Int
	GasPriceSubunits *big.Int
	Observer         commontypes.OracleID
}

func parseAttributedObservation(ao types.AttributedObservation) (ParsedAttributedObservation, error) {
	var observationProto NumericalMedianObservationProto
	if err := proto.Unmarshal(ao.Observation, &observationProto); err != nil {
		return ParsedAttributedObservation{}, fmt.Errorf("attributed observation cannot be unmarshaled: %w", err)
	}
	value, err := DecodeValue(observationProto.Value)
	if err != nil {
		return ParsedAttributedObservation{}, fmt.Errorf("attributed observation with value that cannot be converted to big.Int: %w", err)
	}
	juelsPerFeeCoin, err := DecodeValue(observationProto.JuelsPerFeeCoin)
	if err != nil {
		return ParsedAttributedObservation{}, fmt.Errorf("attributed observation with juelsPerFeeCoin that cannot be converted to big.Int: %w", err)
	}
	var gasPriceSubunits *big.Int
	if len(observationProto.GasPriceSubunits) == 0 {
		// "gasPriceSubunits" may not be sent by nodes in the DON
		// if they are using an older version of the median reporting plugin
		// or in newer versions if IncludeGasPriceSubunitsInObservation is false
		gasPriceSubunits = new(big.Int)
	} else {
		gasPriceSubunits, err = DecodeValue(observationProto.GasPriceSubunits)
		if err != nil {
			return ParsedAttributedObservation{}, fmt.Errorf("attributed observation with gasPriceSubunits that cannot be converted to big.Int: %w", err)
		}
	}

	return ParsedAttributedObservation{
		observationProto.Timestamp,
		value,
		juelsPerFeeCoin,
		gasPriceSubunits,
		ao.Observer,
	}, nil
}

func parseAttributedObservations(logger loghelper.LoggerWithContext, aos []types.AttributedObservation) []ParsedAttributedObservation {
	paos := make([]ParsedAttributedObservation, 0, len(aos))
	for i, ao := range aos {
		pao, err := parseAttributedObservation(ao)
		if err != nil {
			logger.Warn("parseAttributedObservations: dropping invalid observation", commontypes.LogFields{
				"observer": ao.Observer,
				"error":    err,
				"i":        i,
			})
			continue
		}
		paos = append(paos, pao)
	}
	return paos
}

func (nm *numericalMedian) Report(ctx context.Context, repts types.ReportTimestamp, query types.Query, aos []types.AttributedObservation) (bool, types.Report, error) {
	if len(query) != 0 {
		return false, nil, fmt.Errorf("expected empty query")
	}

	paos := parseAttributedObservations(nm.logger, aos)

	// The Report function is guaranteed to receive at least 2f+1 distinct attributed
	// observations. By assumption, up to f of these may be faulty, which includes
	// being malformed. Conversely, there have to be at least f+1 valid paos.
	if !(nm.f+1 <= len(paos)) {
		return false, nil, fmt.Errorf("only received %v valid attributed observations, but need at least f+1 (%v)", len(paos), nm.f+1)
	}

	should, err := nm.shouldReport(ctx, repts, paos)
	if err != nil {
		return false, nil, err
	}
	if !should {
		return false, nil, nil
	}
	report, err := nm.reportCodec.BuildReport(paos)
	if err != nil {
		return false, nil, err
	}
	if !(len(report) <= nm.maxReportLength) {
		return false, nil, fmt.Errorf("report violates MaxReportLength limit set by ReportCodec (%v vs %v)", len(report), nm.maxReportLength)
	}

	return true, report, nil
}

func (nm *numericalMedian) shouldReport(ctx context.Context, repts types.ReportTimestamp, paos []ParsedAttributedObservation) (bool, error) {
	if len(paos) == 0 {
		return false, fmt.Errorf("cannot handle empty attributed observations")
	}

	var resultTransmissionDetails struct {
		configDigest    types.ConfigDigest
		epoch           uint32
		round           uint8
		latestAnswer    *big.Int
		latestTimestamp time.Time
		err             error
	}
	var resultRoundRequested struct {
		configDigest types.ConfigDigest
		epoch        uint32
		round        uint8
		err          error
	}

	var subs subprocesses.Subprocesses
	subs.Go(func() {
		resultTransmissionDetails.configDigest,
			resultTransmissionDetails.epoch,
			resultTransmissionDetails.round,
			resultTransmissionDetails.latestAnswer,
			resultTransmissionDetails.latestTimestamp,
			resultTransmissionDetails.err =
			nm.contractTransmitter.LatestTransmissionDetails(ctx)
	})
	subs.Go(func() {
		resultRoundRequested.configDigest,
			resultRoundRequested.epoch,
			resultRoundRequested.round,
			resultRoundRequested.err =
			nm.contractTransmitter.LatestRoundRequested(ctx, nm.offchainConfig.DeltaC)
	})
	subs.Wait()

	if err := multierr.Combine(resultTransmissionDetails.err, resultRoundRequested.err); err != nil {
		return false, fmt.Errorf("error during LatestTransmissionDetails/LatestRoundRequested: %w", err)
	}

	if resultTransmissionDetails.latestAnswer == nil {
		return false, fmt.Errorf("nil latestAnswer was returned by LatestTransmissionDetails. This should never happen")
	}

	// sort by values
	sort.Slice(paos, func(i, j int) bool {
		return paos[i].Value.Cmp(paos[j].Value) < 0
	})

	answer := paos[len(paos)/2].Value

	if !(nm.onchainConfig.Min.Cmp(answer) <= 0 && answer.Cmp(nm.onchainConfig.Max) <= 0) {
		nm.logger.Warn("shouldReport: no, answer is outside of min/max configured for contract", commontypes.LogFields{
			"result": false,
			"answer": answer,
			"min":    nm.onchainConfig.Min,
			"max":    nm.onchainConfig.Max,
		})
		return false, nil
	}

	initialRound := // Is this the first round for this configuration?
		resultTransmissionDetails.configDigest == repts.ConfigDigest &&
			resultTransmissionDetails.epoch == 0 &&
			resultTransmissionDetails.round == 0
	deviation := // Has the result changed enough to merit a new report?
		!nm.offchainConfig.AlphaReportInfinite &&
			Deviates(nm.offchainConfig.AlphaReportPPB, resultTransmissionDetails.latestAnswer, answer)

	deltaCTimeout := // Has enough time passed since the last report, to merit a new one?
		resultTransmissionDetails.latestTimestamp.Add(nm.offchainConfig.DeltaC).
			Before(time.Now())
	unfulfilledRequest := // Has a new report been requested explicitly?
		resultRoundRequested.configDigest == repts.ConfigDigest &&
			!(epochRound{resultRoundRequested.epoch, resultRoundRequested.round}).
				Less(epochRound{resultTransmissionDetails.epoch, resultTransmissionDetails.round})

	logger := nm.logger.MakeChild(commontypes.LogFields{
		"timestamp":                 repts,
		"initialRound":              initialRound,
		"alphaReportInfinite":       nm.offchainConfig.AlphaReportInfinite,
		"alphaReportPPB":            nm.offchainConfig.AlphaReportPPB,
		"deviation":                 deviation,
		"deltaC":                    nm.offchainConfig.DeltaC,
		"deltaCTimeout":             deltaCTimeout,
		"lastTransmissionTimestamp": resultTransmissionDetails.latestTimestamp,
		"unfulfilledRequest":        unfulfilledRequest,
	})

	// The following is more succinctly expressed as a disjunction, but breaking
	// the branches up into their own conditions makes it easier to check that
	// each branch is tested, and also allows for more expressive log messages
	if initialRound {
		logger.Info("shouldReport: yes, because it's the first round of the first epoch", commontypes.LogFields{
			"result": true,
		})
		return true, nil
	}
	if deviation {
		logger.Info("shouldReport: yes, because new median deviates sufficiently from current onchain value", commontypes.LogFields{
			"result": true,
		})
		return true, nil
	}
	if deltaCTimeout {
		logger.Info("shouldReport: yes, because deltaC timeout since last onchain report", commontypes.LogFields{
			"result": true,
		})
		return true, nil
	}
	if unfulfilledRequest {
		logger.Info("shouldReport: yes, because a new report has been explicitly requested", commontypes.LogFields{
			"result": true,
		})
		return true, nil
	}
	logger.Info("shouldReport: no", commontypes.LogFields{"result": false})
	return false, nil
}

func (nm *numericalMedian) ShouldAcceptFinalizedReport(ctx context.Context, repts types.ReportTimestamp, report types.Report) (bool, error) {
	reportEpochRound := epochRound{repts.Epoch, repts.Round}
	if !nm.latestAcceptedEpochRound.Less(reportEpochRound) {
		nm.logger.Debug("ShouldAcceptFinalizedReport() = false, report is stale", commontypes.LogFields{
			"latestAcceptedEpochRound": nm.latestAcceptedEpochRound,
			"reportEpochRound":         reportEpochRound,
		})
		return false, nil
	}

	contractConfigDigest, contractEpoch, contractRound, _, _, err := nm.contractTransmitter.LatestTransmissionDetails(ctx)
	if err != nil {
		return false, err
	}

	contractEpochRound := epochRound{contractEpoch, contractRound}

	if contractConfigDigest != nm.configDigest {
		nm.logger.Debug("ShouldAcceptFinalizedReport() = false, config digest mismatch", commontypes.LogFields{
			"contractConfigDigest": contractConfigDigest,
			"reportConfigDigest":   nm.configDigest,
			"reportEpochRound":     reportEpochRound,
		})
		return false, nil
	}

	if !contractEpochRound.Less(reportEpochRound) {
		nm.logger.Debug("ShouldAcceptFinalizedReport() = false, report is stale", commontypes.LogFields{
			"contractEpochRound": contractEpochRound,
			"reportEpochRound":   reportEpochRound,
		})
		return false, nil
	}

	if !(len(report) <= nm.maxReportLength) {
		nm.logger.Warn("report violates MaxReportLength limit set by ReportCodec", commontypes.LogFields{
			"reportEpochRound": reportEpochRound,
			"reportLength":     len(report),
			"maxReportLength":  nm.maxReportLength,
		})
		return false, nil
	}

	reportMedian, err := nm.reportCodec.MedianFromReport(report)
	if err != nil {
		return false, fmt.Errorf("error during MedianFromReport: %w", err)
	}

	deviates := !nm.offchainConfig.AlphaAcceptInfinite && Deviates(nm.offchainConfig.AlphaAcceptPPB, nm.latestAcceptedMedian, reportMedian)
	nothingPending := !contractEpochRound.Less(nm.latestAcceptedEpochRound)
	result := deviates || nothingPending

	nm.logger.Debug("ShouldAcceptFinalizedReport() = result", commontypes.LogFields{
		"contractEpochRound":       contractEpochRound,
		"reportEpochRound":         reportEpochRound,
		"latestAcceptedEpochRound": nm.latestAcceptedEpochRound,
		"alphaAcceptInfinite":      nm.offchainConfig.AlphaAcceptInfinite,
		"alphaAcceptPPB":           nm.offchainConfig.AlphaAcceptPPB,
		"deviates":                 deviates,
		"result":                   result,
	})

	if result {
		nm.latestAcceptedEpochRound = reportEpochRound
		nm.latestAcceptedMedian = reportMedian
	}

	return result, nil
}

func (nm *numericalMedian) ShouldTransmitAcceptedReport(ctx context.Context, repts types.ReportTimestamp, report types.Report) (bool, error) {
	reportEpochRound := epochRound{repts.Epoch, repts.Round}

	contractConfigDigest, contractEpoch, contractRound, _, _, err := nm.contractTransmitter.LatestTransmissionDetails(ctx)
	if err != nil {
		return false, err
	}

	contractEpochRound := epochRound{contractEpoch, contractRound}

	if contractConfigDigest != nm.configDigest {
		nm.logger.Debug("ShouldTransmitAcceptedReport() = false, config digest mismatch", commontypes.LogFields{
			"contractConfigDigest": contractConfigDigest,
			"reportConfigDigest":   nm.configDigest,
			"reportEpochRound":     reportEpochRound,
		})
		return false, nil
	}

	if !contractEpochRound.Less(reportEpochRound) {
		nm.logger.Debug("ShouldTransmitAcceptedReport() = false, report is stale", commontypes.LogFields{
			"contractEpochRound": contractEpochRound,
			"reportEpochRound":   reportEpochRound,
		})
		return false, nil
	}

	return true, nil
}

func (nm *numericalMedian) Close() error {
	return nil
}

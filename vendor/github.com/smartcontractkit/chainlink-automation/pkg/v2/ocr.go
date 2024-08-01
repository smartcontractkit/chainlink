package ocr2keepers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-automation/pkg/v2/config"
)

const (
	// ObservationUpkeepsLimit is the maximum number of upkeeps that should be
	// in a single observation
	ObservationUpkeepsLimit = 1
	// ReportKeysLimit is the maximum number of upkeeps that should be added to
	// a single report regardless of report capacity
	ReportKeysLimit = 10
)

var ErrNotEnoughInputs = fmt.Errorf("not enough inputs")

// Encoder provides functions to correctly encode a report from a list of keys.
// This is a really big interface. This should probably be reduced by either
// using reflection at runtime to determine interface compliance or something
// else.
type Encoder interface {
	// GetMedian returns the median BlockKey for the array of provided keys
	GetMedian([]BlockKey) BlockKey
	// ValidateUpkeepKey provides validation for upkeep keys
	// this might go away as we transition to a more structed static type
	ValidateUpkeepKey(UpkeepKey) (bool, error)
	// ValidateUpkeepIdentifier provides validation for upkeep ids. In most
	// cases these ids wil be big.Int but this function allows that validation
	// to be abstracted from the plugin
	ValidateUpkeepIdentifier(UpkeepIdentifier) (bool, error)
	// ValidateBlockKey provides validation for block ids. In most cases this
	// will by uint32 or uint64, but this allows the plugin to not care
	ValidateBlockKey(BlockKey) (bool, error)
	// MakeUpkeepKey combines a block and upkeep id into an upkeep key. This
	// will probably go away with a more structured static upkeep type.
	MakeUpkeepKey(BlockKey, UpkeepIdentifier) UpkeepKey
	// EncodeReport produces a fully encoded report from a slice of upkeeps. The
	// result should be encoded for the chain and fully executable, respecting
	// chain block limits and contract specific validation. The result is
	// expressed in bytes.
	EncodeReport([]UpkeepResult) ([]byte, error)
	// KeysFromReport extracts upkeep keys from an encoded report. This will
	// also change with a more structed static upkeep type. The key detail is
	// the block and upkeep id returned from this function. Both values are
	// needed for a coordinator of conditional upkeeps.
	KeysFromReport([]byte) ([]UpkeepKey, error)
	// Eligible determines if an upkeep is eligible or not. This allows an
	// upkeep result to be abstract and only the encoder is able and responsible
	// for decoding it.
	Eligible(UpkeepResult) (bool, error)
	// Detail is a temporary value that provides upkeep key and gas to perform.
	// A better approach might be needed here.
	Detail(UpkeepResult) (UpkeepKey, uint32, error)
}

// Coordinator provides functions to track in-flight status of upkeeps as they
// move through the OCR process
type Coordinator interface {
	IsPending(UpkeepKey) (bool, error)
	Accept(key UpkeepKey) error
	IsTransmissionConfirmed(key UpkeepKey) bool
}

// ConditionalObserver provides observations queued up conditional upkeeps. This
// type is distinctly different from a log observer
type ConditionalObserver interface {
	Observe() (BlockKey, []UpkeepIdentifier, error)
}

// Runner is the interface for an object that should determine eligibility state
type Runner interface {
	CheckUpkeep(context.Context, bool, ...UpkeepKey) ([]UpkeepResult, error)
}

type ocrPlugin struct {
	// id      commontypes.OracleID
	encoder      Encoder
	coordinator  Coordinator
	condObserver ConditionalObserver
	runner       Runner
	logger       *log.Logger
	subProcs     []PluginStarterCloser
	// configuration vars
	conf           config.OffchainConfig
	mercuryEnabled bool
}

// Query implements the types.ReportingPlugin interface in OCR2. The query
// produced from this method is intended to be empty
func (p *ocrPlugin) Query(_ context.Context, _ types.ReportTimestamp) (types.Query, error) {
	return types.Query{}, nil
}

// Observation implements the types.ReportingPlugin interface in OCR2. This
// method pulls observations from multiple sources and produces an amalgamation
// of individual observations to the libOCR protocol
func (p *ocrPlugin) Observation(_ context.Context, t types.ReportTimestamp, _ types.Query) (types.Observation, error) {
	lCtx := newOcrLogContext(t)

	allIDs := make([]UpkeepIdentifier, 0)

	// naive implementation of getting observations
	// Observer may be too simple and we need a queue mechanism to distribute
	// pulling from multiple observers and their respective queues
	// estimates as items are popped from the queue
	block, ids, err := p.condObserver.Observe()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to sample upkeeps for observation: %s", err, lCtx)
	}

	allIDs = append(allIDs, ids...)

	allIDs = shuffleObservations(allIDs, getRandomKeySource(t))

	// limit the total number of observations if over the limit
	if len(allIDs) > ObservationUpkeepsLimit {
		allIDs = allIDs[:ObservationUpkeepsLimit]
	}

	// build the observation using the median block and all ids
	observation := Observation{
		BlockKey:          block,
		UpkeepIdentifiers: allIDs,
	}

	p.logger.Printf("observation: %v", observation)

	// observations can only be a limited size in bytes after encoding
	// encode the observation and remove ids until encoded bytes is under the
	// limit
	b, err := limitedLengthEncode(observation, MaxObservationLength)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to encode upkeep keys for observation: %s", err, lCtx)
	}

	// write the number of keys returned from sampling to the debug log
	// this offers a record of the number of performs the node has visibility
	// of for each epoch/round
	p.logger.Printf("OCR observation completed successfully with block number %s, %d eligible upkeeps(%s): %s", string(block), len(allIDs), upkeepIdentifiersToString(allIDs), lCtx)

	return b, nil
}

// Report implements the types.ReportingPlugin interface in OC2. This method
// chooses a single upkeep from the provided observations by the earliest block
// number, checks the upkeep, and builds a report. Multiple upkeeps in a single
// report is supported by how the data is abi encoded, but no gas estimations
// exist yet.
func (p *ocrPlugin) Report(ctx context.Context, t types.ReportTimestamp, _ types.Query, attributed []types.AttributedObservation) (bool, types.Report, error) {
	lCtx := newOcrLogContext(t)
	ctx = context.WithValue(ctx, ocrLogContextKey{}, lCtx)

	// --------- Condition Upkeep OCR Logic to Collect Keys from Observations -----------
	var keysToCheck []UpkeepKey
	{
		// Must not be empty
		if len(attributed) == 0 {
			return false, nil, fmt.Errorf("%w: must provide at least 1 observation", ErrNotEnoughInputs)
		}

		var (
			keys [][]UpkeepKey
			err  error
		)

		if keys, err = ObservationsToUpkeepKeys(
			attributed,
			p.encoder,
			p.encoder,
			p.encoder,
			p.logger,
		); err != nil {
			return false, nil, fmt.Errorf("%w: failed to build upkeep keys from the given observations", err)
		}

		upkeepKeysStr := make([]string, len(keys))
		for i, uk := range keys {
			upkeepKeysStr[i] = upkeepKeysToString(uk)
		}
		p.logger.Printf("Parsed observation keys to check in report %s: %s", strings.Join(upkeepKeysStr, ", "), lCtx)

		// pass the filter to the dedupe function
		// ensure no locked keys come through
		keysToCheck, err = filterDedupeShuffleObservations(keys, getRandomKeySource(t), p.coordinator.IsPending)
		if err != nil {
			return false, nil, fmt.Errorf("%w: failed to sort/dedupe attributed observations: %s", err, lCtx)
		}

		p.logger.Printf("Post filtering, deduping and shuffling, keys to check in report %s: %s", upkeepKeysToString(keysToCheck), lCtx)
	}
	// ------------- End getting keys from conditional Upkeeps ----------

	// ------------- Length check added to limit check load from malicious observations
	{
		if len(keysToCheck) > ReportKeysLimit {
			keysToCheck = keysToCheck[:ReportKeysLimit]
		}

		if len(keysToCheck) == 0 {
			p.logger.Printf("OCR report completed successfully with no eligible keys: %s", lCtx)
			return false, nil, nil
		}
	}
	// ------------- End length check ----------

	// -------------- Check Process Specific to Conditional Upkeeps ---------
	checkedUpkeeps, err := p.runner.CheckUpkeep(ctx, p.mercuryEnabled, keysToCheck...)
	if err != nil {
		return false, nil, fmt.Errorf("%w: failed to check upkeeps from attributed observation: %s", err, lCtx)
	}

	// No upkeeps found for the given keys
	if len(checkedUpkeeps) == 0 {
		p.logger.Printf("OCR report completed successfully with no successfully checked upkeeps: %s", lCtx)
		return false, nil, nil
	}

	if len(checkedUpkeeps) > len(keysToCheck) {
		return false, nil, fmt.Errorf("unexpected number of upkeeps returned expected max %d but given %d", len(keysToCheck), len(checkedUpkeeps))
	}
	// ------------ End Conditional Upkeep Check ----------------

	// ------------ Begin Report Building -----------
	// Conditional upkeeps have a different result than log triggered upkeeps.
	// We may need to determine eligibility in a separate loop and have a
	// common type between both the UpkeepResult and the struct derived from a
	// log triggered observation.

	var totalReportGas uint32
	toPerform := make([]UpkeepResult, 0, len(checkedUpkeeps))

	for _, result := range checkedUpkeeps {
		if ok, err := p.encoder.Eligible(result); err != nil && ok {
			continue
		}

		// key is only needed for logging. maybe have a string representation
		// on an upkeep result??
		// gas is necessary to calculate the total number of upkeeps
		// that can be included in the next result
		key, gas, err := p.encoder.Detail(result)
		if err != nil {
			continue
		}

		upkeepMaxGas := gas + p.conf.GasOverheadPerUpkeep
		if totalReportGas+upkeepMaxGas > p.conf.GasLimitPerReport {
			// We don't break here since there could be an upkeep with the lower
			// gas limit so there could be a space for it in the report.
			p.logger.Printf("skipping upkeep %s due to report limit, current capacity is %d, upkeep gas is %d with %d overhead", key, totalReportGas, gas, p.conf.GasOverheadPerUpkeep)
			continue
		}

		p.logger.Printf("reporting %s to be performed with gas limit %d and %d overhead: %s", key, gas, p.conf.GasOverheadPerUpkeep, lCtx.Short())

		toPerform = append(toPerform, result)
		totalReportGas += upkeepMaxGas

		// Don't exceed specified maxUpkeepBatchSize value in offchain config
		if len(toPerform) >= p.conf.MaxUpkeepBatchSize {
			break
		}
	}

	// if nothing to report, return false with no error
	if len(toPerform) == 0 {
		p.logger.Printf("OCR report completed successfully with no eligible upkeeps: %s", lCtx)
		return false, nil, nil
	}

	b, err := p.encoder.EncodeReport(toPerform)
	if err != nil {
		return false, nil, fmt.Errorf("%w: failed to encode OCR report: %s", err, lCtx)
	}

	p.logger.Printf("OCR report completed successfully with %d upkeep added to the report: %s", len(toPerform), lCtx)

	return true, b, nil
}

// ShouldAcceptFinalizedReport implements the types.ReportingPlugin interface
// from OCR2. The implementation checks the length of the report and the number
// of keys in the report. Finally it applies a lockout to all keys in the report
func (p *ocrPlugin) ShouldAcceptFinalizedReport(_ context.Context, rt types.ReportTimestamp, r types.Report) (bool, error) {
	lCtx := newOcrLogContext(rt)

	if len(r) == 0 {
		p.logger.Printf("finalized report is empty; not accepting: %s", lCtx)
		return false, nil
	}

	keys, err := p.encoder.KeysFromReport(r)
	if err != nil {
		return false, fmt.Errorf("%w: failed to decode report: %s", err, lCtx)
	}

	if len(keys) == 0 {
		p.logger.Printf("no upkeeps in report; not accepting: %s", lCtx)
		return false, fmt.Errorf("no ids in report: %s", lCtx)
	}

	for _, key := range keys {
		// indicate to the filter that the key has been accepted for transmit
		if err = p.coordinator.Accept(key); err != nil {
			return false, fmt.Errorf("%w: failed to accept key: %s", err, lCtx)
		}
		p.logger.Printf("accepting key %s: %s", key, lCtx.Short())
	}

	p.logger.Printf("OCR should accept completed successfully: %s", lCtx)

	return true, nil
}

// ShouldTransmitAcceptedReport implements the types.ReportingPlugin interface
// from OCR2. The implementation essentially draws straws on which node should
// be the transmitter.
func (p *ocrPlugin) ShouldTransmitAcceptedReport(_ context.Context, rt types.ReportTimestamp, r types.Report) (bool, error) {
	lCtx := newOcrLogContext(rt)

	keys, err := p.encoder.KeysFromReport(r)
	if err != nil {
		return false, fmt.Errorf("%w: failed to get ids from report: %s", err, lCtx)
	}

	if len(keys) == 0 {
		return false, fmt.Errorf("no ids in report: %s", lCtx)
	}

	for _, key := range keys {
		transmitConfirmed := p.coordinator.IsTransmissionConfirmed(key)
		// multiple keys can be in a single report. if one has a non confirmed transmission
		// (while others may not have), try to transmit anyway
		if !transmitConfirmed {
			p.logger.Printf("upkeep '%s' transmit not confirmed, transmitting whole report: %s", key, lCtx.Short())
			p.logger.Printf("OCR should transmit completed successfully with result true: %s", lCtx)
			return true, nil
		}
		p.logger.Printf("upkeep '%s' was already transmitted: %s", key, lCtx)
	}

	p.logger.Printf("OCR should transmit completed successfully with result false: %s", lCtx)

	return false, nil
}

// Close implements the types.ReportingPlugin interface in OCR2.
// internal services before or after libOCR calls this function. Also, does this
// function get called before or after a new instance is created.
func (p *ocrPlugin) Close() error {
	var finalErr error

	// need to close dependent services first
	for _, proc := range p.subProcs {
		if err := proc.Close(); err != nil {
			finalErr = errors.Join(finalErr, err)
		}
	}

	return finalErr
}

func upkeepIdentifiersToString(ids []UpkeepIdentifier) string {
	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = string(id)
	}

	return strings.Join(idsStr, ", ")
}

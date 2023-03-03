package keepers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/ocr2keepers/pkg/chain"
	ktypes "github.com/smartcontractkit/ocr2keepers/pkg/types"
)

const (
	// observationUpkeepsLimit is the max number of upkeeps that Observation could return.
	observationUpkeepsLimit = 1

	// reportKeysLimit is the maximum number of upkeep keys checked during the report phase
	reportKeysLimit = 10
)

type ocrLogContextKey struct{}

type ocrLogContext struct {
	Epoch     uint32
	Round     uint8
	StartTime time.Time
}

func newOcrLogContext(rt types.ReportTimestamp) ocrLogContext {
	return ocrLogContext{
		Epoch:     rt.Epoch,
		Round:     rt.Round,
		StartTime: time.Now(),
	}
}

func (c ocrLogContext) String() string {
	return fmt.Sprintf("[epoch=%d, round=%d, completion=%dms]", c.Epoch, c.Round, time.Since(c.StartTime)/time.Millisecond)
}

func (c ocrLogContext) Short() string {
	return fmt.Sprintf("[epoch=%d, round=%d]", c.Epoch, c.Round)
}

// Query implements the types.ReportingPlugin interface in OCR2. The query produced from this
// method is intended to be empty.
func (k *keepers) Query(_ context.Context, _ types.ReportTimestamp) (types.Query, error) {
	return types.Query{}, nil
}

// Observation implements the types.ReportingPlugin interface in OCR2. This method samples a set
// of upkeeps available in and UpkeepService and produces an observation containing upkeeps that
// need to be executed.
func (k *keepers) Observation(ctx context.Context, rt types.ReportTimestamp, _ types.Query) (types.Observation, error) {
	lCtx := newOcrLogContext(rt)
	ctx = context.WithValue(ctx, ocrLogContextKey{}, lCtx)

	blockKey, results, err := k.service.SampleUpkeeps(ctx, k.filter.Filter())
	if err != nil {
		return nil, fmt.Errorf("%w: failed to sample upkeeps for observation: %s", err, lCtx)
	}

	// keyList produces a sorted result so the following reduction of keys
	// should be more uniform for all nodes
	keys := keyList(filterUpkeeps(results, ktypes.Eligible))

	obs := &chain.UpkeepObservation{
		BlockKey: chain.BlockKey(blockKey.String()),
	}

	identifiers := make([]ktypes.UpkeepIdentifier, 0)
	for _, upkeepKey := range keys {
		_, upkeepID, _ := upkeepKey.BlockKeyAndUpkeepID()
		identifiers = append(identifiers, upkeepID)
	}

	// Shuffle the observations before we limit it to observationUpkeepsLimit
	keyRandSource := getRandomKeySource(rt)
	identifiers = shuffleObservations(identifiers, keyRandSource)
	// Check limit
	if len(identifiers) > observationUpkeepsLimit {
		identifiers = identifiers[:observationUpkeepsLimit]
	}

	obs.UpkeepIdentifiers = identifiers

	b, err := limitedLengthEncode(obs, maxObservationLength)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to encode upkeep keys for observation: %s", err, lCtx)
	}

	// write the number of keys returned from sampling to the debug log
	// this offers a record of the number of performs the node has visibility
	// of for each epoch/round
	k.logger.Printf("OCR observation completed successfully with block number %s, %d eligible upkeeps(%s): %s", blockKey, len(identifiers), upkeepIdentifiersToString(identifiers), lCtx)

	return b, nil
}

// Report implements the types.ReportingPlugin interface in OC2. This method chooses a single upkeep
// from the provided observations by the earliest block number, checks the upkeep, and builds a
// report. Multiple upkeeps in a single report is supported by how the data is abi encoded, but
// no gas estimations exist yet.
func (k *keepers) Report(ctx context.Context, rt types.ReportTimestamp, _ types.Query, attributed []types.AttributedObservation) (bool, types.Report, error) {
	var err error

	lCtx := newOcrLogContext(rt)
	ctx = context.WithValue(ctx, ocrLogContextKey{}, lCtx)

	// Must not be empty
	if len(attributed) == 0 {
		return false, nil, fmt.Errorf("%w: must provide at least 1 observation", ErrNotEnoughInputs)
	}

	// Build upkeep keys from the given observations
	upkeepKeys, err := observationsToUpkeepKeys(k.logger, attributed, k.reportBlockLag)
	if err != nil {
		return false, nil, fmt.Errorf("%w: failed to build upkeep keys from the given observations", err)
	}

	upkeepKeysStr := make([]string, len(upkeepKeys))
	for i, uk := range upkeepKeys {
		upkeepKeysStr[i] = upkeepKeysToString(uk)
	}
	k.logger.Printf("Parsed observation keys to check in report %s: %s", strings.Join(upkeepKeysStr, ", "), lCtx)

	// pass the filter to the dedupe function
	// ensure no locked keys come through
	keyRandSource := getRandomKeySource(rt)
	keysToCheck, err := filterDedupeShuffleObservations(upkeepKeys, keyRandSource, k.filter.Filter())
	if err != nil {
		return false, nil, fmt.Errorf("%w: failed to sort/dedupe attributed observations: %s", err, lCtx)
	}
	k.logger.Printf("Post filtering, deduping and shuffling, keys to check in report %s: %s", upkeepKeysToString(keysToCheck), lCtx)

	// Check the limit
	if len(keysToCheck) > reportKeysLimit {
		keysToCheck = keysToCheck[:reportKeysLimit]
	}

	// No keys found for the given keys
	if len(keysToCheck) == 0 {
		k.logger.Printf("OCR report completed successfully with no eligible keys: %s", lCtx)
		return false, nil, nil
	}

	// Check all upkeeps from the given observation
	checkedUpkeeps, err := k.service.CheckUpkeep(ctx, keysToCheck...)
	if err != nil {
		return false, nil, fmt.Errorf("%w: failed to check upkeeps from attributed observation: %s", err, lCtx)
	}

	// No upkeeps found for the given keys
	if len(checkedUpkeeps) == 0 {
		k.logger.Printf("OCR report completed successfully with no successfully checked upkeeps: %s", lCtx)
		return false, nil, nil
	}

	if len(checkedUpkeeps) > len(keysToCheck) {
		return false, nil, fmt.Errorf("unexpected number of upkeeps returned expected max %d but given %d", len(keysToCheck), len(checkedUpkeeps))
	}

	// Collect eligible upkeeps
	var reportCapacity uint32
	toPerform := make([]ktypes.UpkeepResult, 0, len(checkedUpkeeps))
	for _, checkedUpkeep := range checkedUpkeeps {
		if checkedUpkeep.State != ktypes.Eligible {
			continue
		}

		upkeepMaxGas := checkedUpkeep.ExecuteGas + k.upkeepGasOverhead
		if reportCapacity+upkeepMaxGas > k.reportGasLimit {
			// We don't break here since there could be an upkeep with the lower
			// gas limit so there could be a space for it in the report.
			k.logger.Printf("skipping upkeep %s due to report limit, current capacity is %d, upkeep gas is %d with %d overhead", checkedUpkeep.Key, reportCapacity, checkedUpkeep.ExecuteGas, k.upkeepGasOverhead)
			continue
		}

		k.logger.Printf("reporting %s to be performed with gas limit %d and %d overhead: %s", checkedUpkeep.Key, checkedUpkeep.ExecuteGas, k.upkeepGasOverhead, lCtx.Short())

		toPerform = append(toPerform, checkedUpkeep)
		reportCapacity += upkeepMaxGas

		// Don't exceed specified maxUpkeepBatchSize value in offchain config
		if len(toPerform) >= k.maxUpkeepBatchSize {
			break
		}
	}

	// if nothing to report, return false with no error
	if len(toPerform) == 0 {
		k.logger.Printf("OCR report completed successfully with no eligible upkeeps: %s", lCtx)
		return false, nil, nil
	}

	b, err := k.encoder.EncodeReport(toPerform)
	if err != nil {
		return false, nil, fmt.Errorf("%w: failed to encode OCR report: %s", err, lCtx)
	}

	k.logger.Printf("OCR report completed successfully with %d upkeep added to the report: %s", len(toPerform), lCtx)

	return true, b, err
}

// ShouldAcceptFinalizedReport implements the types.ReportingPlugin interface
// from OCR2. The implementation checks the length of the report and the number
// of keys in the report. Finally it applies a lockout to all keys in the report
func (k *keepers) ShouldAcceptFinalizedReport(_ context.Context, rt types.ReportTimestamp, r types.Report) (bool, error) {
	lCtx := newOcrLogContext(rt)

	if len(r) == 0 {
		k.logger.Printf("finalized report is empty; not accepting: %s", lCtx)
		return false, nil
	}

	results, err := k.encoder.DecodeReport(r)
	if err != nil {
		return false, fmt.Errorf("%w: failed to decode report: %s", err, lCtx)
	}

	if len(results) == 0 {
		k.logger.Printf("no upkeeps in report; not accepting: %s", lCtx)
		return false, fmt.Errorf("no ids in report: %s", lCtx)
	}

	for _, r := range results {
		// indicate to the filter that the key has been accepted for transmit
		if err = k.filter.Accept(r.Key); err != nil {
			return false, fmt.Errorf("%w: failed to accept key: %s", err, lCtx)
		}
		k.logger.Printf("accepting key %s: %s", r.Key, lCtx.Short())
	}

	k.logger.Printf("OCR should accept completed successfully: %s", lCtx)

	return true, nil
}

// ShouldTransmitAcceptedReport implements the types.ReportingPlugin interface
// from OCR2. The implementation essentially draws straws on which node should
// be the transmitter.
func (k *keepers) ShouldTransmitAcceptedReport(_ context.Context, rt types.ReportTimestamp, r types.Report) (bool, error) {
	lCtx := newOcrLogContext(rt)

	results, err := k.encoder.DecodeReport(r)
	if err != nil {
		return false, fmt.Errorf("%w: failed to get ids from report: %s", err, lCtx)
	}

	if len(results) == 0 {
		return false, fmt.Errorf("no ids in report: %s", lCtx)
	}

	for _, id := range results {
		transmitConfirmed := k.filter.IsTransmissionConfirmed(id.Key)
		// multiple keys can be in a single report. if one has a non confirmed transmission
		// (while others may not have), try to transmit anyway
		if !transmitConfirmed {
			k.logger.Printf("upkeep '%s' transmit not confirmed, transmitting whole report: %s", id.Key, lCtx.Short())
			k.logger.Printf("OCR should transmit completed successfully with result true: %s", lCtx)
			return true, nil
		}
		k.logger.Printf("upkeep '%s' was already transmitted: %s", id.Key, lCtx)
	}

	k.logger.Printf("OCR should transmit completed successfully with result false: %s", lCtx)
	return false, nil
}

// Close implements the types.ReportingPlugin interface in OCR2.
func (k *keepers) Close() error {
	return nil
}

package ocr2keepers

import (
	"fmt"
	"log"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var (
	ErrBlockKeyNotParsable     = fmt.Errorf("block identifier not parsable")
	ErrUpkeepKeyNotParsable    = fmt.Errorf("upkeep key not parsable")
	ErrInvalidBlockKey         = fmt.Errorf("invalid block key")
	ErrInvalidUpkeepIdentifier = fmt.Errorf("invalid upkeep identifier")
	ErrTooManyErrors           = fmt.Errorf("too many errors in parallel worker process")
)

type Validator interface {
	ValidateUpkeepKey(UpkeepKey) (bool, error)
	ValidateUpkeepIdentifier(UpkeepIdentifier) (bool, error)
	ValidateBlockKey(BlockKey) (bool, error)
}

type Builder interface {
	MakeUpkeepKey(BlockKey, UpkeepIdentifier) UpkeepKey
}

type MedianCalculator interface {
	// GetMedian returns the median BlockKey for the array of provided keys
	GetMedian([]BlockKey) BlockKey
}

// Observation defines the data structure that nodes use to communication the
// details of observed upkeeps
type Observation struct {
	BlockKey          BlockKey           `json:"1"`
	UpkeepIdentifiers []UpkeepIdentifier `json:"2"`
}

func (u Observation) Validate(v Validator) error {
	if ok, err := v.ValidateBlockKey(u.BlockKey); !ok || err != nil {
		if err != nil {
			return err
		}

		return ErrInvalidBlockKey
	}

	for _, ui := range u.UpkeepIdentifiers {
		if ok, err := v.ValidateUpkeepIdentifier(ui); !ok || err != nil {
			if err != nil {
				return err
			}

			return ErrInvalidUpkeepIdentifier
		}
	}

	return nil
}

// ObservationsToUpkeepKeys loops through all observations, collects the
// UpkeepIdentifier list from each one, calculates the median block number, and
// constructs upkeep keys from each identifier with the median block number
func ObservationsToUpkeepKeys(
	attr []types.AttributedObservation,
	v Validator,
	e MedianCalculator,
	b Builder,
	logger *log.Logger,
) ([][]UpkeepKey, error) {
	var (
		parseErrors  int
		allBlockKeys []BlockKey
	)

	upkeepIDs := make([][]UpkeepIdentifier, 0, len(attr))

	for _, obs := range attr {
		// a single observation returning an error here can void all other
		// good observations. ensure this loop continues on error, but collect
		// them and throw an error if ALL observations fail at this point.
		var upkeepObservation Observation
		if err := decode(obs.Observation, &upkeepObservation); err != nil {
			logger.Printf("unable to decode observation: %s", err.Error())
			parseErrors++
			continue
		}

		// validate the observation using the provided validator
		if err := upkeepObservation.Validate(v); err != nil {
			logger.Printf("failed to validate observation: %s", err.Error())
			parseErrors++
			continue
		}

		allBlockKeys = append(allBlockKeys, upkeepObservation.BlockKey)

		// if we have a non-empty list of upkeep identifiers, limit the upkeeps
		// we take to observationUpkeepsLimit
		if len(upkeepObservation.UpkeepIdentifiers) > 0 {
			ids := upkeepObservation.UpkeepIdentifiers[:]

			if len(ids) > ObservationUpkeepsLimit {
				ids = ids[:ObservationUpkeepsLimit]
			}

			upkeepIDs = append(upkeepIDs, ids)
		}
	}

	if parseErrors == len(attr) {
		return nil, fmt.Errorf("%w: cannot prepare sorted key list; observations not properly encoded", ErrTooManyErrors)
	}

	// Here we calculate the median block that will be applied for all upkeep keys.
	// reportBlockLag is subtracted from the median block to ensure enough nodes have that block in their blockchain
	medianBlock := e.GetMedian(allBlockKeys)
	// logger.Printf("calculated median block %s, accounting for reportBlockLag of %d", medianBlock, reportBlockLag)

	upkeepKeys, err := createKeysWithMedianBlock(b, medianBlock, upkeepIDs)
	if err != nil {
		return nil, err
	}

	return upkeepKeys, nil
}

func createKeysWithMedianBlock(b Builder, medianBlock BlockKey, upkeepIDLists [][]UpkeepIdentifier) ([][]UpkeepKey, error) {
	var res = make([][]UpkeepKey, len(upkeepIDLists))

	for i, upkeepIDs := range upkeepIDLists {
		var keys []UpkeepKey

		for _, upkeepID := range upkeepIDs {
			keys = append(keys, b.MakeUpkeepKey(medianBlock, upkeepID))
		}

		res[i] = keys
	}

	return res, nil
}

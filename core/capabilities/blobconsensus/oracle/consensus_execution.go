package oracle

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"slices"
	"time"

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

var (
	ErrNoValuesMetThreshold                         = errors.New("no values met f+1 threshold")
	ErrMoreThanOneValidOutcomeForIdenticalConsensus = errors.New("not identical, multiple values with f+1 occurrences")
	ErrInsufficientObservations                     = errors.New("insufficient observations to reach consensus")
)

// Constants for type names used in aggregation logic.
var (
	typeInt64   = reflect.TypeOf((*valuespb.Value_Int64Value)(nil))
	typeUint64  = reflect.TypeOf((*valuespb.Value_Uint64Value)(nil))
	typeFloat64 = reflect.TypeOf((*valuespb.Value_Float64Value)(nil))
	typeDecimal = reflect.TypeOf((*valuespb.Value_DecimalValue)(nil))
	typeBigInt  = reflect.TypeOf((*valuespb.Value_BigintValue)(nil))
	typeString  = reflect.TypeOf((*valuespb.Value_StringValue)(nil))
	typeTime    = reflect.TypeOf((*valuespb.Value_TimeValue)(nil))
)

// CalculateOutcomeForObservations determines the outcome for a set of observations based on a consensus descriptor.
// It now supports median aggregation for Int64, Uint64, Float64, Decimal, BigInt, and Time types. It assumes that the observationProtos
// are already validated to ensure they all correctly unmarshal to a values.Value
func CalculateOutcomeForObservations(
	lggr logger.Logger,
	observations []*valuespb.Value,
	consensusDescriptor *sdk.ConsensusDescriptor,
	defaultValue *valuespb.Value,
	f int,
) (*valuespb.Value, error) {
	switch desc := consensusDescriptor.GetDescriptor_().(type) {
	case *sdk.ConsensusDescriptor_Aggregation:
		aggregation := consensusDescriptor.GetAggregation()
		switch aggregation {
		case sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL:
			return handleIdenticalAggregation(lggr, observations, f)
		case sdk.AggregationType_AGGREGATION_TYPE_MEDIAN:
			return handleMedianAggregation(lggr, observations, f)
		case sdk.AggregationType_AGGREGATION_TYPE_COMMON_PREFIX:
			return handleCommonPrefixAggregation(lggr, observations, f)
		case sdk.AggregationType_AGGREGATION_TYPE_COMMON_SUFFIX:
			return handleCommonSuffixAggregation(lggr, observations, f)
		default:
			return nil, fmt.Errorf("unknown aggregation type: %s", aggregation)
		}
	case *sdk.ConsensusDescriptor_FieldsMap:
		return handleFieldsMapAggregation(lggr, observations, desc.FieldsMap.GetFields(), defaultValue, f)
	default:
		return nil, fmt.Errorf("unknown consensus descriptor type: %T", desc)
	}
}

func handleFieldsMapAggregation(
	lggr logger.Logger,
	observations []*valuespb.Value,
	desc map[string]*sdk.ConsensusDescriptor,
	defaultValue *valuespb.Value,
	f int,
) (*valuespb.Value, error) {
	if len(observations) < f+1 {
		return nil, ErrInsufficientObservations
	}

	sortedKeys := make([]string, 0, len(desc))
	for k := range desc {
		sortedKeys = append(sortedKeys, k)
	}
	slices.Sort(sortedKeys)

	result := make(map[string]*valuespb.Value, 0)
	for _, key := range sortedKeys {
		d := desc[key]
		var (
			aggregated *valuespb.Value
			err        error
			obsForKey  = make([]*valuespb.Value, 0, len(observations))
		)

		for i, obs := range observations {
			if obs != nil {
				switch obs.Value.(type) {
				case *valuespb.Value_MapValue:
					fields := obs.GetMapValue().GetFields()
					obsForKey = append(obsForKey, fields[key])
				default:
					lggr.Debugf("unsupported observation type at index %d for key %s: %T", i, key, obs.Value)
					continue
				}
			} else {
				lggr.Debugf("ignoring nil observation at index %d for key %s", i, key)
			}
		}

		var defaultForKey *valuespb.Value
		if defaultValue != nil {
			switch defaultValue.Value.(type) {
			case *valuespb.Value_MapValue:
				defaultForKey = defaultValue.GetMapValue().GetFields()[key]
			default:
				lggr.Debugf("missing default for key: %s", key)
			}
		}

		aggregated, err = CalculateOutcomeForObservations(lggr, obsForKey, d, defaultForKey, f)
		if err == nil {
			result[key] = aggregated
			continue
		}

		if defaultForKey != nil {
			result[key] = defaultForKey
			continue
		}

		return nil, fmt.Errorf("aggregation for field failed '%s': %w", key, err)
	}

	return valuespb.NewMapValue(result), nil
}

func handleMedianAggregation(
	lggr logger.Logger,
	observations []*valuespb.Value,
	f int,
) (*valuespb.Value, error) {
	var (
		medianResult *valuespb.Value
		err          error
	)

	// The Report function is guaranteed to receive at least 2f+1 distinct attributed
	// observations. By assumption, up to f of these may be faulty, which includes
	// being malformed. Conversely, there have to be at least f+1 valid observations.
	filtered, medianType, err := filterObservations(observations, f+1)
	if err != nil {
		return nil, err
	}

	switch medianType {
	case typeUint64:
		medianResult, err = getMedian(lggr,
			filtered,
			func(val *valuespb.Value) (uint64, error) {
				return val.GetUint64Value(), nil
			},
			func(a, b uint64) int {
				if a < b {
					return -1
				}
				if a > b {
					return 1
				}
				return 0
			},
			f,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate uint64 median: %w", err)
		}

	case typeInt64:
		medianResult, err = getMedian(lggr,
			filtered,
			func(val *valuespb.Value) (int64, error) {
				return val.GetInt64Value(), nil
			},
			func(a, b int64) int {
				if a < b {
					return -1
				}
				if a > b {
					return 1
				}
				return 0
			},
			f,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate int64 median: %w", err)
		}

	case typeFloat64:
		medianResult, err = getMedian(lggr,
			filtered,
			func(val *valuespb.Value) (float64, error) {
				return val.GetFloat64Value(), nil
			},
			func(a, b float64) int {
				if a < b {
					return -1
				}
				if a > b {
					return 1
				}
				return 0
			},
			f,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate float64 median: %w", err)
		}

	case typeDecimal:
		medianResult, err = getMedian(lggr,
			filtered,
			func(val *valuespb.Value) (decimal.Decimal, error) {
				var d decimal.Decimal
				v, err := values.FromProto(val)
				if err != nil {
					return d, err
				}
				return d, v.UnwrapTo(&d)
			},
			func(a, b decimal.Decimal) int {
				return a.Cmp(b)
			},
			f,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate decimal median: %w", err)
		}

	case typeBigInt:
		medianResult, err = getMedian(lggr,
			filtered,
			func(val *valuespb.Value) (*big.Int, error) {
				var got big.Int
				v, err := values.FromProto(val)
				if err != nil {
					return nil, err
				}
				return &got, v.UnwrapTo(&got)
			},
			func(a, b *big.Int) int {
				return a.Cmp(b)
			},
			f,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate big.Int median: %w", err)
		}

	case typeTime:
		medianResult, err = getMedian(lggr,
			filtered,
			func(val *valuespb.Value) (time.Time, error) {
				var got time.Time
				v, err := values.FromProto(val)
				if err != nil {
					return got, err
				}
				return got, v.UnwrapTo(&got)
			},
			func(a, b time.Time) int {
				return a.Compare(b)
			},
			f,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate time median: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported type for median aggregation: %s", medianType)
	}

	return medianResult, nil
}

func handleIdenticalAggregation(_ logger.Logger, values []*valuespb.Value, f int) (*valuespb.Value, error) {
	n := len(values)
	if n == 0 {
		return nil, errors.New("input slice cannot be empty for identical aggregation")
	}

	type valueOccurrence struct {
		count int
		value *valuespb.Value
	}

	var (
		marshaler       = &proto.MarshalOptions{Deterministic: true}
		identityMap     = make(map[string]valueOccurrence)
		uniqueCandidate *valuespb.Value
	)

	for _, currentValue := range values {
		b, err := marshaler.Marshal(currentValue)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal value: %w", err)
		}
		key := string(b)

		observation := identityMap[key]
		observation.count++
		if observation.value == nil {
			observation.value = currentValue
		}
		identityMap[key] = observation

		if observation.count == f+1 {
			if uniqueCandidate != nil {
				return nil, ErrMoreThanOneValidOutcomeForIdenticalConsensus
			}
			uniqueCandidate = observation.value
		}
	}

	if uniqueCandidate == nil {
		return nil, ErrNoValuesMetThreshold
	}

	return uniqueCandidate, nil
}

// handleCommonSuffixAggregation reverses the underlying lists in the slice of
// observations and delegates logic to handleCommonPrefixAggregation and then
// reverses the result a final time.
func handleCommonSuffixAggregation(lggr logger.Logger, observationSlices []*valuespb.Value, f int) (*valuespb.Value, error) {
	var reversedObservations []*valuespb.Value
	for i, obsProto := range observationSlices {
		reversed, err := reverseListValue(obsProto)
		if err != nil {
			lggr.Warnf("skipping observations at index %d: %s", i, err)
			continue
		}
		reversedObservations = append(reversedObservations, reversed)
	}

	commonPrefixOfReversed, err := handleCommonPrefixAggregation(lggr, reversedObservations, f)
	if err != nil {
		return nil, fmt.Errorf("failed to find common prefix of reversed lists: %w", err)
	}

	return reverseListValue(commonPrefixOfReversed)
}

// handleCommonPrefixAggregation finds the longest common prefix among a set of
// observations.
//
// It iterates through the sets index by index. At each index,
// it performs an identical aggregation check on the corresponding observations
// aggregated from all sets.
//
// A set of observations is only carried forward to the next index check if its
// observation at the current index matches the consensus value.
//
// The process stops when either:
// 1. The identical aggregation fails to reach the f+1 threshold (ErrNoValuesMetThreshold).
// 2. The number of remaining consensus-matching lists drops below the f+1 threshold.
//
// This ensures that the common prefix is robustly determined by the longest sequence
// of elements that at least f+1 lists agree on.
func handleCommonPrefixAggregation(
	lggr logger.Logger,
	observations []*valuespb.Value,
	f int,
) (*valuespb.Value, error) {
	// Transform slice of values to slice of lists
	var currentLists []*valuespb.List
	var maxListLength int
	for i, obsProto := range observations {
		if obsProto != nil {
			switch obsProto.Value.(type) {
			case *valuespb.Value_ListValue:
				list := obsProto.GetListValue()
				currentLists = append(currentLists, list)
				if len(list.GetFields()) > maxListLength {
					maxListLength = len(list.GetFields())
				}
			default:
				lggr.Warnf("value at index %d is of type %T", i, obsProto.Value)
				continue
			}
		}
	}

	if len(currentLists) < f+1 {
		return nil, ErrInsufficientObservations
	}

	var commonPrefixElements []*valuespb.Value
	var nextLists []*valuespb.List

	for i := range maxListLength {
		// Aggregate observations at index
		var elementsAtIndex []*valuespb.Value
		for _, list := range currentLists {
			if len(list.GetFields()) > i {
				elementsAtIndex = append(elementsAtIndex, list.GetFields()[i])
			}
		}

		identicalValue, err := handleIdenticalAggregation(lggr, elementsAtIndex, f)
		if err != nil {
			// Consensus failed at this index, so the common prefix ends here.
			break
		}

		// Update the set of lists to select from for the next index (i+1)
		nextLists = make([]*valuespb.List, 0, len(currentLists))
		for _, list := range currentLists {
			if len(list.GetFields()) > i {
				// Check if the current list element matches the identical consensus value
				if proto.Equal(list.GetFields()[i], identicalValue) {
					nextLists = append(nextLists, list)
				}
			}
		}

		// If the filtered set of lists no longer meets the f+1 threshold,
		// the common prefix must end here.
		if len(nextLists) < f+1 {
			break
		}

		commonPrefixElements = append(commonPrefixElements, identicalValue)

		// Update the list of observations for the next iteration
		currentLists = nextLists
	}

	return valuespb.NewListValue(commonPrefixElements), nil
}

// filterObservations returns all the observations that meet the minimum observation
// threshold of the same underlying type.  Errors if no single type meets the
// threshold.
func filterObservations(observationProtos []*valuespb.Value, minObservations int) ([]*valuespb.Value, reflect.Type, error) {
	if len(observationProtos) < minObservations {
		return nil, nil, fmt.Errorf("insufficient observations (%d) to meet minimum (%d)", len(observationProtos), minObservations)
	}

	observationsByType := map[reflect.Type][]*valuespb.Value{}
	for _, observation := range observationProtos {
		if observation.Value == nil {
			continue
		}

		tpe := reflect.TypeOf(observation.Value)
		observationsByType[tpe] = append(observationsByType[tpe], observation)
	}

	var dominantType reflect.Type
	for tpe, obsOfType := range observationsByType {
		if len(obsOfType) >= minObservations {
			if dominantType != nil {
				// More than one type meets the threshold
				return nil, nil, ErrMoreThanOneValidOutcomeForIdenticalConsensus
			}
			dominantType = tpe
		}
	}

	if dominantType == nil {
		return nil, nil, ErrNoValuesMetThreshold
	}

	return observationsByType[dominantType], dominantType, nil
}

// getMedian is a generic helper function that calculates the median
// for a slice of values.Value that can be unwrapped to type T.
// It accepts functions for unwrapping and comparing the values.
//
// For an even number of elements, we take the left of the two middle elements.
func getMedian[T any](
	lggr logger.Logger,
	observations []*valuespb.Value,
	unwrap func(val *valuespb.Value) (T, error),
	compare func(a, b T) int,
	f int,
) (*valuespb.Value, error) {
	if len(observations) < f+1 {
		return nil, ErrInsufficientObservations
	}

	var unwrappedValues []T
	for _, v := range observations {
		unwrapped, err := unwrap(v)
		if err != nil {
			// It's possible the value could be corrupt and fail to unwrap, so skip and log warning as this should not happen
			lggr.Warnf("failed to unwrap observation during median calculation: %s", err)
		} else {
			unwrappedValues = append(unwrappedValues, unwrapped)
		}
	}

	// As values are filtered for unwrapping errors, need to re-check the number of observations is still sufficient for consensus
	if len(unwrappedValues) < f+1 {
		return nil, ErrInsufficientObservations
	}

	slices.SortFunc(unwrappedValues, compare)

	medianVal := unwrappedValues[len(unwrappedValues)/2]
	if len(unwrappedValues)%2 == 0 && len(unwrappedValues) > 0 {
		medianVal = unwrappedValues[len(unwrappedValues)/2-1]
	}

	v, err := values.Wrap(medianVal)
	if err != nil {
		return nil, err
	}

	return values.Proto(v), nil
}

func reverseListValue(list *valuespb.Value) (*valuespb.Value, error) {
	if list != nil {
		switch list.Value.(type) {
		case *valuespb.Value_ListValue:
			reversed := list.GetListValue().GetFields()
			reverse(reversed)
			return valuespb.NewListValue(reversed), nil
		default:
			return nil, fmt.Errorf("cannot reverse value of type %T", list.Value)
		}
	}
	return new(valuespb.Value), nil
}

func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

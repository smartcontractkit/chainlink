package read

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
)

type EventQuery struct {
	Filter           query.KeyFilter
	EventBinding     *EventBinding
	SequenceDataType any
	IsValuePtr       bool
	Address          common.Address
}

func MultiEventTypeQuery(ctx context.Context, lp logpoller.LogPoller, eventQueries []EventQuery, limitAndSort query.LimitAndSort) (iter.Seq2[string, commontypes.Sequence], error) {
	seqIter, err := multiEventTypeQueryWithoutErrorWrapping(ctx, lp, eventQueries, limitAndSort)
	if err != nil {
		if len(eventQueries) > 0 {
			var calls []Call
			for _, eq := range eventQueries {
				calls = append(calls, Call{
					ContractAddress: eq.Address,
					ContractName:    eq.EventBinding.contractName,
					ReadName:        eq.EventBinding.eventName,
					ReturnVal:       eq.SequenceDataType,
				})
			}

			err = newErrorFromCalls(err, calls, "", eventReadType)
		} else {
			err = fmt.Errorf("no event queries provided: %w", err)
		}
	}

	return seqIter, err
}

func multiEventTypeQueryWithoutErrorWrapping(ctx context.Context, lp logpoller.LogPoller, eventQueries []EventQuery, limitAndSort query.LimitAndSort) (iter.Seq2[string, commontypes.Sequence], error) {
	if err := validateEventQueries(eventQueries); err != nil {
		return nil, fmt.Errorf("error validating event queries: %w", err)
	}

	for _, eq := range eventQueries {
		if err := eq.EventBinding.validateBound(eq.Address); err != nil {
			return nil, err
		}
	}

	allFilterExpressions := make([]query.Expression, 0, len(eventQueries))
	for _, eq := range eventQueries {
		var expressions []query.Expression

		defaultExpressions := []query.Expression{
			logpoller.NewAddressFilter(eq.Address),
			logpoller.NewEventSigFilter(eq.EventBinding.hash),
		}
		expressions = append(expressions, defaultExpressions...)

		remapped, remapErr := eq.EventBinding.remap(eq.Filter)
		if remapErr != nil {
			return nil, fmt.Errorf("error remapping filter: %w", remapErr)
		}
		expressions = append(expressions, remapped.Expressions...)

		filterExpression := query.And(expressions...)

		allFilterExpressions = append(allFilterExpressions, filterExpression)
	}

	eventQuery := query.Or(allFilterExpressions...)

	queryName := createQueryName(eventQueries)

	logs, err := lp.FilteredLogs(ctx, []query.Expression{eventQuery}, limitAndSort, queryName)
	if err != nil {
		return nil, wrapInternalErr(err)
	}

	seqIter, err := decodeMultiEventTypeLogsIntoSequences(ctx, logs, eventQueries)
	if err != nil {
		return nil, wrapInternalErr(err)
	}

	return seqIter, nil
}

func createQueryName(eventQueries []EventQuery) string {
	queryName := ""
	contractToEvents := map[string][]string{}
	for _, eq := range eventQueries {
		contractName := eq.EventBinding.contractName + "-" + eq.Address.String()

		if _, exists := contractToEvents[contractName]; !exists {
			contractToEvents[contractName] = []string{}
		}
		contractToEvents[contractName] = append(contractToEvents[contractName], eq.EventBinding.eventName)
	}

	contractNames := make([]string, 0, len(contractToEvents))
	for contractName := range contractToEvents {
		contractNames = append(contractNames, contractName)
	}

	sort.Strings(contractNames)

	for _, contractName := range contractNames {
		queryName += contractName + "-"
		for _, event := range contractToEvents[contractName] {
			queryName += event + "-"
		}
	}

	queryName = strings.TrimSuffix(queryName, "-")
	return queryName
}

func validateEventQueries(eventQueries []EventQuery) error {
	duplicateCheck := map[common.Hash]EventQuery{}
	for _, eq := range eventQueries {
		if eq.EventBinding == nil {
			return errors.New("event binding is nil")
		}

		if eq.SequenceDataType == nil {
			return errors.New("sequence data type is nil")
		}

		// TODO support queries with the same event signature but different filters
		if _, exists := duplicateCheck[eq.EventBinding.hash]; exists {
			return fmt.Errorf("duplicate event query for event signature %s, event name %s", eq.EventBinding.hash, eq.EventBinding.eventName)
		}
		duplicateCheck[eq.EventBinding.hash] = eq
	}
	return nil
}

func decodeMultiEventTypeLogsIntoSequences(ctx context.Context, logs []logpoller.Log, eventQueries []EventQuery) (iter.Seq2[string, commontypes.Sequence], error) {
	type sequenceWithKey struct {
		Key      string
		Sequence commontypes.Sequence
	}
	sequenceWithKeys := make([]sequenceWithKey, 0, len(logs))
	eventSigToEventQuery := map[common.Hash]EventQuery{}
	for _, eq := range eventQueries {
		eventSigToEventQuery[eq.EventBinding.hash] = eq
	}

	for _, logEntry := range logs {
		eventSignatureHash := logEntry.EventSig

		eq, exists := eventSigToEventQuery[eventSignatureHash]
		if !exists {
			return nil, fmt.Errorf("no event query found for log with event signature %s", eventSignatureHash)
		}

		seqWithKey := sequenceWithKey{
			Key: eq.Filter.Key,
			Sequence: commontypes.Sequence{
				Cursor: logpoller.FormatContractReaderCursor(logEntry),
				Head: commontypes.Head{
					Height:    strconv.FormatInt(logEntry.BlockNumber, 10),
					Hash:      logEntry.BlockHash.Bytes(),
					Timestamp: uint64(logEntry.BlockTimestamp.Unix()), //nolint:gosec // G115 false positive
				},
			},
		}

		var typeVal reflect.Value

		typeInto := reflect.TypeOf(eq.SequenceDataType)
		if typeInto.Kind() == reflect.Pointer {
			typeVal = reflect.New(typeInto.Elem())
		} else {
			typeVal = reflect.Indirect(reflect.New(typeInto))
		}

		// create a new value of the same type as 'into' for the data to be extracted to
		seqWithKey.Sequence.Data = typeVal.Interface()

		if err := eq.EventBinding.decodeLog(ctx, &logEntry, seqWithKey.Sequence.Data); err != nil {
			return nil, err
		}

		if eq.IsValuePtr {
			wrappedValue, err := values.Wrap(seqWithKey.Sequence.Data)
			if err != nil {
				return nil, err
			}
			seqWithKey.Sequence.Data = &wrappedValue
		}

		sequenceWithKeys = append(sequenceWithKeys, seqWithKey)
	}

	return func(yield func(string, commontypes.Sequence) bool) {
		for _, s := range sequenceWithKeys {
			if !yield(s.Key, s.Sequence) {
				return
			}
		}
	}, nil
}

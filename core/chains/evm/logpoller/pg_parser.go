package logpoller

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

func orderBy(sortBy []SortBy) string {
	if len(sortBy) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(" ORDER BY ")
	for i, sort := range sortBy {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s %s", sort.field, sort.dir.pgString()))
	}
	return sb.String()
}

func limitBy(limit int) string {
	if limit == 0 {
		return ""
	}
	return fmt.Sprintf(" LIMIT %d", limit)
}

// PgParserVisitor is a visitor that builds a postgres query and arguments from a QFilter
// Keep in mind it's not designed to be thread safe so should not be used concurrently.
type PgParserVisitor struct {
	evmChainId *big.Int
	query      strings.Builder
	args       map[string]any
	errors     []error
	// filterNameIndex is an increment used for avoiding name collision when filtering the same column by multiple values
	filterNameIndex int
}

func NewPgParser(evmChainID *big.Int) *PgParserVisitor {
	return &PgParserVisitor{
		evmChainId: evmChainID,
		args:       map[string]any{},
	}
}

func (v *PgParserVisitor) Parse(f QFilter) (string, map[string]any, error) {
	f = v.addEvmChainIdFilter(f)
	f.accept(v)
	return v.query.String(), v.args, errors.Join(v.errors...)
}

func (v *PgParserVisitor) addEvmChainIdFilter(filter QFilter) QFilter {
	switch casted := filter.(type) {
	case *AndFilter:
		return AppendedNewFilter(casted, NewEvmChainIdFilter(v.evmChainId))
	default:
		return NewAndFilter(casted, NewEvmChainIdFilter(v.evmChainId))
	}
}

func (op ComparisonOperator) pgString() string {
	switch op {
	case Eq:
		return "="
	case Neq:
		return "<>"
	case Gt:
		return ">"
	case Lt:
		return "<"
	case Gte:
		return ">="
	case Lte:
		return "<="
	default:
		panic("unknown comparison operator")
	}
}

func (sd SortDirection) pgString() string {
	switch sd {
	case Asc:
		return "asc"
	case Desc:
		return "desc"
	default:
		panic("unknown sort direction")
	}
}

func (v *PgParserVisitor) VisitAndFilter(node AndFilter) {
	for i, filter := range node.filters {
		if i > 0 {
			v.query.WriteString(" AND ")
		}
		filter.accept(v)
	}
}

func (v *PgParserVisitor) VisitAddressFilter(node AddressFilter) {
	switch len(node.address) {
	case 0:
		v.errors = append(v.errors, errors.New("address filter must have at least one address signature"))
	case 1:
		v.query.WriteString("address = :address")
		v.args["address"] = node.address[0]
	default:
		v.query.WriteString("address = ANY(:address)")
		v.args["address"] = concatBytes(node.address)
	}
}

func (v *PgParserVisitor) VisitEventSigFilter(node EventSigFilter) {
	switch len(node.eventSig) {
	case 0:
		v.errors = append(v.errors, errors.New("event sig filter must have at least one event signature"))
	case 1:
		v.query.WriteString("event_sig = :event_sig")
		v.args["event_sig"] = node.eventSig[0].Bytes()
	default:
		v.query.WriteString("event_sig = ANY(:event_sig)")
		v.args["event_sig"] = concatBytes(node.eventSig)
	}
}

func (v *PgParserVisitor) VisitDataWordFilter(node DataWordFilter) {
	argumentName, valueName := v.wordValueVariables()

	v.query.WriteString(fmt.Sprintf(
		"substring(data from 32*:%s+1 for 32) %s :%s",
		argumentName,
		node.operator.pgString(),
		valueName,
	))
	v.args[argumentName] = node.index
	v.args[valueName] = node.value.Bytes()
}

func (v *PgParserVisitor) VisitConfirmationFilter(node ConfirmationFilter) {
	v.query.WriteString("block_number <= ")
	v.query.WriteString(nestedBlockNumberQuery(node.confs))
	v.args["confs"] = node.confs
}

func (v *PgParserVisitor) VisitEvmChainIdFilter(node EvmChainIdFilter) {
	v.query.WriteString("evm_chain_id = :evm_chain_id")
	v.args["evm_chain_id"] = node.chainId
}

func (v *PgParserVisitor) VisitTopicFilter(node TopicFilter) {
	topicIndex, err := sanitizeTopicIndex(node.index)
	if err != nil {
		v.errors = append(v.errors, err)
		return
	}

	topicIndexName, topicValueName := v.topicVariables()
	v.query.WriteString(fmt.Sprintf("topics[:%s] %s :%s", topicIndexName, node.operator.pgString(), topicValueName))
	v.args[topicIndexName] = topicIndex
	v.args[topicValueName] = node.value.Bytes()
}

func (v *PgParserVisitor) VisitTopicsFilter(node TopicsFilter) {
	topicIndex, err := sanitizeTopicIndex(node.index)
	if err != nil {
		v.errors = append(v.errors, err)
		return
	}

	topicIndexName, topicValueName := v.topicVariables()
	v.query.WriteString(fmt.Sprintf("topics[:%s] = ANY(:%s)", topicIndexName, topicValueName))
	v.args[topicIndexName] = topicIndex
	v.args[topicValueName] = concatBytes(node.values)
}

func (v *PgParserVisitor) topicVariables() (string, string) {
	v.filterNameIndex++
	argumentName := fmt.Sprintf("topic_index_%d", v.filterNameIndex)
	valueName := fmt.Sprintf("topic_value_%d", v.filterNameIndex)
	return argumentName, valueName
}

func (v *PgParserVisitor) wordValueVariables() (string, string) {
	v.filterNameIndex++
	argumentName := fmt.Sprintf("word_index_%d", v.filterNameIndex)
	valueName := fmt.Sprintf("word_value_%d", v.filterNameIndex)
	return argumentName, valueName
}

func (v *PgParserVisitor) VisitBlockFilter(node BlockFilter) {
	v.visitSimpleNameFilter("block_number", node.operator, node.block)
}

func (v *PgParserVisitor) VisitBlockTimestampFilter(node BlockTimestampFilter) {
	v.visitSimpleNameFilter("block_timestamp", node.operator, node.timestamp)
}

func (v *PgParserVisitor) VisitTxHashFilter(node TxHashFilter) {
	v.visitSimpleNameFilter("tx_hash", Eq, node.txHash.Bytes())
}

func (v *PgParserVisitor) visitSimpleNameFilter(fieldName string, operator ComparisonOperator, value interface{}) {
	v.filterNameIndex++
	argumentValue := fmt.Sprintf("%s_%d", fieldName, v.filterNameIndex)

	v.query.WriteString(fmt.Sprintf("%s %s :%s", fieldName, operator.pgString(), argumentValue))
	v.args[argumentValue] = value
}

func sanitizeTopicIndex(index int) (int, error) {
	// Only topicIndex 1 through 3 is valid. 0 is the event sig and only 4 total topics are allowed
	if index < 1 || index > 3 {
		return 0, fmt.Errorf("invalid index for topic: %d", index)
	}
	return index + 1, nil
}

func nestedBlockNumberQuery(confs Confirmations) string {
	if confs == Finalized {
		return `
				(SELECT finalized_block_number 
				FROM evm.log_poller_blocks 
				WHERE evm_chain_id = :evm_chain_id 
				ORDER BY block_number DESC LIMIT 1) `
	}
	// Intentionally wrap with greatest() function and don't return negative block numbers when :confs > :block_number
	// It doesn't impact logic of the outer query, because block numbers are never less or equal to 0 (guarded by log_poller_blocks_block_number_check)
	return `
			(SELECT greatest(block_number - :confs, 0) 
			FROM evm.log_poller_blocks 	
			WHERE evm_chain_id = :evm_chain_id 
			ORDER BY block_number DESC LIMIT 1) `

}

package logpoller

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

const (
	blockFieldName     = "block_number"
	timestampFieldName = "block_timestamp"
	txHashFieldName    = "tx_hash"
	eventSigFieldName  = "event_sig"
)

var (
	ErrUnexpectedCursorFormat = errors.New("unexpected cursor format")
)

// The parser builds SQL expressions piece by piece for each Accept function call and resets the error and expression
// values after every call.
type pgDSLParser struct {
	args *queryArgs

	// transient properties expected to be set and reset with every expression
	expression string
	err        error
}

var _ primitives.Visitor = (*pgDSLParser)(nil)

func (v *pgDSLParser) Comparator(_ primitives.Comparator) {}

func (v *pgDSLParser) Block(p primitives.Block) {
	cmp, err := cmpOpToString(p.Operator)
	if err != nil {
		v.err = err

		return
	}

	v.expression = fmt.Sprintf(
		"%s %s :%s",
		blockFieldName,
		cmp,
		v.args.withIndexedField(blockFieldName, p.Block),
	)
}

func (v *pgDSLParser) Confirmations(p primitives.Confirmations) {
	switch p.ConfirmationLevel {
	case primitives.Finalized:
		v.expression = v.nestedConfQuery(true, 0)
	case primitives.Unconfirmed:
		// Unconfirmed in the evm relayer is an alias to the case of 0 confirmations
		// set the level to the number 0 and fallthrough to the default case
		p.ConfirmationLevel = primitives.ConfirmationLevel(0)

		fallthrough
	default:
		// the default case passes the confirmation level as a number directly to a subquery
		v.expression = v.nestedConfQuery(false, uint64(evmtypes.Confirmations(p.ConfirmationLevel)))
	}
}

func (v *pgDSLParser) Timestamp(p primitives.Timestamp) {
	cmp, err := cmpOpToString(p.Operator)
	if err != nil {
		v.err = err

		return
	}

	v.expression = fmt.Sprintf(
		"%s %s :%s",
		timestampFieldName,
		cmp,
		v.args.withIndexedField(timestampFieldName, time.Unix(int64(p.Timestamp), 0)),
	)
}

func (v *pgDSLParser) TxHash(p primitives.TxHash) {
	bts, err := hexutil.Decode(p.TxHash)
	if errors.Is(err, hexutil.ErrMissingPrefix) {
		bts, err = hexutil.Decode("0x" + p.TxHash)
	}

	if err != nil {
		v.err = err

		return
	}

	txHash := common.BytesToHash(bts)

	v.expression = fmt.Sprintf(
		"%s = :%s",
		txHashFieldName,
		v.args.withIndexedField(txHashFieldName, txHash),
	)
}

func (v *pgDSLParser) VisitAddressFilter(p *addressFilter) {
	v.expression = fmt.Sprintf(
		"address = :%s",
		v.args.withIndexedField("address", p.address),
	)
}

func (v *pgDSLParser) VisitEventSigFilter(p *eventSigFilter) {
	v.expression = fmt.Sprintf(
		"%s = :%s",
		eventSigFieldName,
		v.args.withIndexedField(eventSigFieldName, p.eventSig),
	)
}

func (v *pgDSLParser) nestedConfQuery(finalized bool, confs uint64) string {
	var (
		from     = "FROM evm.log_poller_blocks "
		where    = "WHERE evm_chain_id = :evm_chain_id "
		order    = "ORDER BY block_number DESC LIMIT 1"
		selector string
	)

	if finalized {
		selector = "SELECT finalized_block_number "
	} else {
		selector = fmt.Sprintf("SELECT greatest(block_number - :%s, 0) ",
			v.args.withIndexedField("confs", confs),
		)
	}

	var builder strings.Builder

	builder.WriteString(selector)
	builder.WriteString(from)
	builder.WriteString(where)
	builder.WriteString(order)

	return fmt.Sprintf("%s <= (%s)", blockFieldName, builder.String())
}

func (v *pgDSLParser) VisitEventByWordFilter(p *eventByWordFilter) {
	if len(p.ValueComparers) > 0 {
		wordIdx := v.args.withIndexedField("word_index", p.WordIndex)

		comps := make([]string, len(p.ValueComparers))
		for idx, comp := range p.ValueComparers {
			comps[idx], v.err = makeComp(comp, v.args, "word_value", wordIdx, "substring(data from 32*:%s+1 for 32) %s :%s")
			if v.err != nil {
				return
			}
		}

		v.expression = strings.Join(comps, " AND ")
	}
}

func (v *pgDSLParser) VisitEventTopicsByValueFilter(p *eventByTopicFilter) {
	if len(p.ValueComparers) > 0 {
		if !(p.Topic == 1 || p.Topic == 2 || p.Topic == 3) {
			v.err = fmt.Errorf("invalid index for topic: %d", p.Topic)

			return
		}

		// Add 1 since postgresql arrays are 1-indexed.
		topicIdx := v.args.withIndexedField("topic_index", p.Topic+1)

		comps := make([]string, len(p.ValueComparers))
		for idx, comp := range p.ValueComparers {
			comps[idx], v.err = makeComp(comp, v.args, "topic_value", topicIdx, "topics[:%s] %s :%s")
			if v.err != nil {
				return
			}
		}

		v.expression = strings.Join(comps, " AND ")
	}
}

func makeComp(comp primitives.ValueComparator, args *queryArgs, field, subfield, pattern string) (string, error) {
	cmp, err := cmpOpToString(comp.Operator)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		pattern,
		subfield,
		cmp,
		args.withIndexedField(field, common.HexToHash(comp.Value)),
	), nil
}

func (v *pgDSLParser) buildQuery(chainID *big.Int, expressions []query.Expression, limiter query.LimitAndSort) (string, *queryArgs, error) {
	// reset transient properties
	v.args = newQueryArgs(chainID)
	v.expression = ""
	v.err = nil

	// build the query string
	clauses := []string{"SELECT evm.logs.* FROM evm.logs"}

	where, err := v.whereClause(expressions, limiter)
	if err != nil {
		return "", nil, err
	}

	clauses = append(clauses, where)

	order, err := v.orderClause(limiter)
	if err != nil {
		return "", nil, err
	}

	if len(order) > 0 {
		clauses = append(clauses, order)
	}

	limit := v.limitClause(limiter)
	if len(limit) > 0 {
		clauses = append(clauses, limit)
	}

	return strings.Join(clauses, " "), v.args, nil
}

func (v *pgDSLParser) whereClause(expressions []query.Expression, limiter query.LimitAndSort) (string, error) {
	segment := "WHERE evm_chain_id = :evm_chain_id"

	if len(expressions) > 0 {
		exp, err := v.combineExpressions(expressions, query.AND)
		if err != nil {
			return "", err
		}

		segment = fmt.Sprintf("%s AND %s", segment, exp)
	}

	if limiter.HasCursorLimit() {
		var op string
		switch limiter.Limit.CursorDirection {
		case query.CursorFollowing:
			op = ">"
		case query.CursorPrevious:
			op = "<"
		default:
			return "", errors.New("invalid cursor direction")
		}

		block, txHash, logIdx, err := valuesFromCursor(limiter.Limit.Cursor)
		if err != nil {
			return "", err
		}

		segment = fmt.Sprintf("%s AND block_number %s= :cursor_block AND tx_hash %s= :cursor_txhash AND log_index %s :cursor_log_index", segment, op, op, op)

		v.args.withField("cursor_block_number", block).
			withField("cursor_txhash", common.HexToHash(txHash)).
			withField("cursor_log_index", logIdx)
	}

	return segment, nil
}

func (v *pgDSLParser) orderClause(limiter query.LimitAndSort) (string, error) {
	sorting := limiter.SortBy

	if limiter.HasCursorLimit() && !limiter.HasSequenceSort() {
		var dir query.SortDirection

		switch limiter.Limit.CursorDirection {
		case query.CursorFollowing:
			dir = query.Asc
		case query.CursorPrevious:
			dir = query.Desc
		default:
			return "", errors.New("unexpected cursor direction")
		}

		sorting = append(sorting, query.NewSortBySequence(dir))
	}

	if len(sorting) == 0 {
		return "", nil
	}

	sort := make([]string, len(sorting))

	for idx, sorted := range sorting {
		var name string

		order, err := orderToString(sorted.GetDirection())
		if err != nil {
			return "", err
		}

		switch sorted.(type) {
		case query.SortByBlock:
			name = blockFieldName
		case query.SortBySequence:
			sort[idx] = fmt.Sprintf("block_number %s, tx_hash %s, log_index %s", order, order, order)

			continue
		case query.SortByTimestamp:
			name = timestampFieldName
		default:
			return "", errors.New("unexpected sort by")
		}

		sort[idx] = fmt.Sprintf("%s %s", name, order)
	}

	return fmt.Sprintf("ORDER BY %s", strings.Join(sort, ", ")), nil
}

func (v *pgDSLParser) limitClause(limiter query.LimitAndSort) string {
	if !limiter.HasCursorLimit() && limiter.Limit.Count == 0 {
		return ""
	}

	return fmt.Sprintf("LIMIT %d", limiter.Limit.Count)
}

func (v *pgDSLParser) getLastExpression() (string, error) {
	exp := v.expression
	err := v.err

	v.expression = ""
	v.err = nil

	return exp, err
}

func (v *pgDSLParser) combineExpressions(expressions []query.Expression, op query.BoolOperator) (string, error) {
	grouped := len(expressions) > 1
	clauses := make([]string, len(expressions))

	for idx, exp := range expressions {
		if exp.IsPrimitive() {
			exp.Primitive.Accept(v)

			clause, err := v.getLastExpression()
			if err != nil {
				return "", err
			}

			clauses[idx] = clause
		} else {
			clause, err := v.combineExpressions(exp.BoolExpression.Expressions, exp.BoolExpression.BoolOperator)
			if err != nil {
				return "", err
			}

			clauses[idx] = clause
		}
	}

	output := strings.Join(clauses, fmt.Sprintf(" %s ", op.String()))

	if grouped {
		output = fmt.Sprintf("(%s)", output)
	}

	return output, nil
}

func cmpOpToString(op primitives.ComparisonOperator) (string, error) {
	switch op {
	case primitives.Eq:
		return "=", nil
	case primitives.Neq:
		return "!=", nil
	case primitives.Gt:
		return ">", nil
	case primitives.Gte:
		return ">=", nil
	case primitives.Lt:
		return "<", nil
	case primitives.Lte:
		return "<=", nil
	default:
		return "", errors.New("invalid comparison operator")
	}
}

func orderToString(dir query.SortDirection) (string, error) {
	switch dir {
	case query.Asc:
		return "ASC", nil
	case query.Desc:
		return "DESC", nil
	default:
		return "", errors.New("invalid sort direction")
	}
}

func valuesFromCursor(cursor string) (int64, string, int, error) {
	parts := strings.Split(cursor, "-")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("%w: must be composed as block-txhash-logindex", ErrUnexpectedCursorFormat)
	}

	block, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", 0, fmt.Errorf("%w: block number not parsable as int64", ErrUnexpectedCursorFormat)
	}

	logIdx, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		return 0, "", 0, fmt.Errorf("%w: log index not parsable as int", ErrUnexpectedCursorFormat)
	}

	return block, parts[1], int(logIdx), nil
}

type addressFilter struct {
	address common.Address
}

func NewAddressFilter(address common.Address) query.Expression {
	return query.Expression{
		Primitive: &addressFilter{address: address},
	}
}

func (f *addressFilter) Accept(visitor primitives.Visitor) {
	switch v := visitor.(type) {
	case *pgDSLParser:
		v.VisitAddressFilter(f)
	}
}

type eventSigFilter struct {
	eventSig common.Hash
}

func NewEventSigFilter(hash common.Hash) query.Expression {
	return query.Expression{
		Primitive: &eventSigFilter{eventSig: hash},
	}
}

func (f *eventSigFilter) Accept(visitor primitives.Visitor) {
	switch v := visitor.(type) {
	case *pgDSLParser:
		v.VisitEventSigFilter(f)
	}
}

type eventByWordFilter struct {
	EventSig       common.Hash
	WordIndex      uint8
	ValueComparers []primitives.ValueComparator
}

func NewEventByWordFilter(eventSig common.Hash, wordIndex uint8, valueComparers []primitives.ValueComparator) query.Expression {
	return query.Expression{Primitive: &eventByWordFilter{
		EventSig:       eventSig,
		WordIndex:      wordIndex,
		ValueComparers: valueComparers,
	}}
}

func (f *eventByWordFilter) Accept(visitor primitives.Visitor) {
	switch v := visitor.(type) {
	case *pgDSLParser:
		v.VisitEventByWordFilter(f)
	}
}

type eventByTopicFilter struct {
	Topic          uint64
	ValueComparers []primitives.ValueComparator
}

func NewEventByTopicFilter(topicIndex uint64, valueComparers []primitives.ValueComparator) query.Expression {
	return query.Expression{Primitive: &eventByTopicFilter{
		Topic:          topicIndex,
		ValueComparers: valueComparers,
	}}
}

func (f *eventByTopicFilter) Accept(visitor primitives.Visitor) {
	switch v := visitor.(type) {
	case *pgDSLParser:
		v.VisitEventTopicsByValueFilter(f)
	}
}

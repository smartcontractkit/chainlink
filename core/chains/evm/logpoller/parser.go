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
	defaultSort        = "block_number DESC, log_index DESC"
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

func (v *pgDSLParser) Confidence(p primitives.Confidence) {
	switch p.ConfidenceLevel {
	case primitives.Finalized:
		// the highest level of confidence maps to finalized
		v.expression = v.nestedConfQuery(true, 0)
	case primitives.Unconfirmed:
		v.expression = v.nestedConfQuery(false, 0)
	default:
		v.err = errors.New("unrecognized confidence level; use confidence to confirmations mappings instead")

		return
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
	if len(p.HashedValueComparers) > 0 {
		wordIdx := v.args.withIndexedField("word_index", p.WordIndex)

		comps := make([]string, len(p.HashedValueComparers))
		for idx, comp := range p.HashedValueComparers {
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

func (v *pgDSLParser) VisitConfirmationsFilter(p *confirmationsFilter) {
	switch p.Confirmations {
	case evmtypes.Finalized:
		// the highest level of confidence maps to finalized
		v.expression = v.nestedConfQuery(true, 0)
	default:
		v.expression = v.nestedConfQuery(false, uint64(p.Confirmations))
	}
}

func makeComp(comp HashedValueComparator, args *queryArgs, field, subfield, pattern string) (string, error) {
	cmp, err := cmpOpToString(comp.Operator)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		pattern,
		subfield,
		cmp,
		args.withIndexedField(field, comp.Value),
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
		exp, hasFinalized, err := v.combineExpressions(expressions, query.AND)
		if err != nil {
			return "", err
		}

		if limiter.HasCursorLimit() && !hasFinalized {
			return "", errors.New("cursor-base queries limited to only finalized blocks")
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

		block, logIdx, _, err := valuesFromCursor(limiter.Limit.Cursor)
		if err != nil {
			return "", err
		}

		segment = fmt.Sprintf("%s AND (block_number %s :cursor_block_number OR (block_number = :cursor_block_number AND log_index %s :cursor_log_index))", segment, op, op)

		v.args.withField("cursor_block_number", block).
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
		return fmt.Sprintf("ORDER BY %s", defaultSort), nil
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
			sort[idx] = fmt.Sprintf("block_number %s, log_index %s, tx_hash %s", order, order, order)

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

func (v *pgDSLParser) combineExpressions(expressions []query.Expression, op query.BoolOperator) (string, bool, error) {
	grouped := len(expressions) > 1
	clauses := make([]string, len(expressions))

	var isFinalized bool

	for idx, exp := range expressions {
		if exp.IsPrimitive() {
			exp.Primitive.Accept(v)

			switch prim := exp.Primitive.(type) {
			case *primitives.Confidence:
				isFinalized = prim.ConfidenceLevel == primitives.Finalized
			case *confirmationsFilter:
				isFinalized = prim.Confirmations == evmtypes.Finalized
			}

			clause, err := v.getLastExpression()
			if err != nil {
				return "", isFinalized, err
			}

			clauses[idx] = clause
		} else {
			clause, fin, err := v.combineExpressions(exp.BoolExpression.Expressions, exp.BoolExpression.BoolOperator)
			if err != nil {
				return "", isFinalized, err
			}

			if fin {
				isFinalized = fin
			}

			clauses[idx] = clause
		}
	}

	output := strings.Join(clauses, fmt.Sprintf(" %s ", op.String()))

	if grouped {
		output = fmt.Sprintf("(%s)", output)
	}

	return output, isFinalized, nil
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

func valuesFromCursor(cursor string) (int64, int, []byte, error) {
	partCount := 3

	parts := strings.Split(cursor, "-")
	if len(parts) != partCount {
		return 0, 0, nil, fmt.Errorf("%w: must be composed as block-logindex-txHash", ErrUnexpectedCursorFormat)
	}

	block, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("%w: block number not parsable as int64", ErrUnexpectedCursorFormat)
	}

	logIdx, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("%w: log index not parsable as int", ErrUnexpectedCursorFormat)
	}

	txHash, err := hexutil.Decode(parts[2])
	if err != nil {
		return 0, 0, nil, fmt.Errorf("%w: invalid transaction hash: %s", ErrUnexpectedCursorFormat, err.Error())
	}

	return block, int(logIdx), txHash, nil
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

type HashedValueComparator struct {
	Value    common.Hash
	Operator primitives.ComparisonOperator
}

type eventByWordFilter struct {
	WordIndex            uint8
	HashedValueComparers []HashedValueComparator
}

func NewEventByWordFilter(wordIndex uint8, valueComparers []HashedValueComparator) query.Expression {
	return query.Expression{Primitive: &eventByWordFilter{
		WordIndex:            wordIndex,
		HashedValueComparers: valueComparers,
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
	ValueComparers []HashedValueComparator
}

func NewEventByTopicFilter(topicIndex uint64, valueComparers []HashedValueComparator) query.Expression {
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

type confirmationsFilter struct {
	Confirmations evmtypes.Confirmations
}

func NewConfirmationsFilter(confirmations evmtypes.Confirmations) query.Expression {
	return query.Expression{Primitive: &confirmationsFilter{
		Confirmations: confirmations,
	}}
}

func (f *confirmationsFilter) Accept(visitor primitives.Visitor) {
	switch v := visitor.(type) {
	case *pgDSLParser:
		v.VisitConfirmationsFilter(f)
	}
}

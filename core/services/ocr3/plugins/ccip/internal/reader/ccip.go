package reader

import "C"
import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/ccipocr3/internal/model"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

var (
	ErrChainReaderNotFound = errors.New("chain reader not found")
)

type CCIP interface {
	// MsgsBetweenSeqNums reads the provided chains.
	// Finds and returns ccip messages submitted between the provided sequence numbers.
	// Messages are sorted ascending based on their timestamp and limited up to the provided limit.
	MsgsBetweenSeqNums(ctx context.Context, chain model.ChainSelector, seqNumRange model.SeqNumRange) ([]model.CCIPMsg, error)

	// NextSeqNum reads the destination chain.
	// Returns the next expected sequence number for each one of the provided chains.
	NextSeqNum(ctx context.Context, chains []model.ChainSelector) (seqNum []model.SeqNum, err error)

	// GasPrices reads the provided chains gas prices.
	GasPrices(ctx context.Context, chains []model.ChainSelector) ([]model.BigInt, error)

	// Close closes any open resources.
	Close(ctx context.Context) error
}

type CCIPChainReader struct {
	chainReaders map[model.ChainSelector]types.ChainReader
	destChain    model.ChainSelector
}

func (r *CCIPChainReader) MsgsBetweenSeqNums(ctx context.Context, chain model.ChainSelector, seqNumRange model.SeqNumRange) ([]model.CCIPMsg, error) {
	if err := r.validateReaderExistence(chain); err != nil {
		return nil, err
	}

	const (
		contractName       = "OnRamp"
		eventName          = "CCIPSendRequested"
		eventAttributeName = "SequenceNumber"
	)

	seq, err := r.chainReaders[chain].QueryKey(
		ctx,
		contractName,
		query.KeyFilter{
			Key: eventName,
			Expressions: []query.Expression{
				{
					Primitive: &primitives.Comparator{
						Name: eventAttributeName,
						ValueComparators: []primitives.ValueComparator{
							{
								Value:    seqNumRange.Start().String(),
								Operator: primitives.Gte,
							},
							{
								Value:    seqNumRange.End().String(),
								Operator: primitives.Lte,
							},
						},
					},
					BoolExpression: query.BoolExpression{},
				},
			},
		},
		query.LimitAndSort{
			SortBy: []query.SortBy{
				query.NewSortByTimestamp(query.Asc),
			},
			Limit: query.Limit{
				Count: uint64(seqNumRange.End() - seqNumRange.Start() + 1),
			},
		},
		&model.CCIPMsg{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query onRamp: %w", err)
	}

	msgs := make([]model.CCIPMsg, 0)
	for _, item := range seq {
		msg, ok := item.Data.(model.CCIPMsg)
		if !ok {
			return nil, fmt.Errorf("failed to cast %v to CCIPMsg", item.Data)
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (r *CCIPChainReader) NextSeqNum(ctx context.Context, chains []model.ChainSelector) ([]model.SeqNum, error) {
	if err := r.validateReaderExistence(r.destChain); err != nil {
		return nil, err
	}

	const (
		contractName = "OffRamp"
		funcName     = "getExpectedNextSequenceNumbers"
	)

	seqNums := make([]model.SeqNum, 0)
	err := r.chainReaders[r.destChain].GetLatestValue(
		ctx,
		contractName,
		funcName,
		map[string]any{
			"chains": chains,
		},
		&seqNums,
	)
	return seqNums, err
}

func (r *CCIPChainReader) GasPrices(ctx context.Context, chains []model.ChainSelector) ([]model.BigInt, error) {
	if err := r.validateReaderExistence(chains...); err != nil {
		return nil, err
	}
	panic("implement me")
}

func (r *CCIPChainReader) Close(ctx context.Context) error {
	return nil
}

func (r *CCIPChainReader) validateReaderExistence(chains ...model.ChainSelector) error {
	for _, ch := range chains {
		_, exists := r.chainReaders[ch]
		if !exists {
			return fmt.Errorf("chain %d: %w", ch, ErrChainReaderNotFound)
		}
	}
	return nil
}

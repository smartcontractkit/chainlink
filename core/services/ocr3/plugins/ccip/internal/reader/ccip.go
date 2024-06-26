package reader

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"golang.org/x/sync/errgroup"
)

var (
	ErrContractReaderNotFound = errors.New("contract reader not found")
	ErrContractWriterNotFound = errors.New("contract writer not found")
)

// TODO: unit test the implementation when the actual contract reader and writer interfaces are finalized and mocks
// can be generated.
type CCIPChainReader struct {
	contractReaders map[cciptypes.ChainSelector]types.ContractReader
	contractWriters map[cciptypes.ChainSelector]types.ChainWriter
	destChain       cciptypes.ChainSelector
}

func NewCCIPChainReader(
	contractReaders map[cciptypes.ChainSelector]types.ContractReader,
	contractWriters map[cciptypes.ChainSelector]types.ChainWriter,
	destChain cciptypes.ChainSelector,
) *CCIPChainReader {
	return &CCIPChainReader{
		contractReaders: contractReaders,
		contractWriters: contractWriters,
		destChain:       destChain,
	}
}

func (r *CCIPChainReader) CommitReportsGTETimestamp(
	ctx context.Context, dest cciptypes.ChainSelector, ts time.Time, limit int,
) ([]cciptypes.CommitPluginReportWithMeta, error) {
	if err := r.validateReaderExistence(dest); err != nil {
		return nil, err
	}
	panic("implement me")
}

func (r *CCIPChainReader) ExecutedMessageRanges(
	ctx context.Context, source, dest cciptypes.ChainSelector, seqNumRange cciptypes.SeqNumRange,
) ([]cciptypes.SeqNumRange, error) {
	if err := r.validateReaderExistence(source, dest); err != nil {
		return nil, err
	}
	panic("implement me")
}

func (r *CCIPChainReader) MsgsBetweenSeqNums(
	ctx context.Context, chain cciptypes.ChainSelector, seqNumRange cciptypes.SeqNumRange,
) ([]cciptypes.CCIPMsg, error) {
	if err := r.validateReaderExistence(chain); err != nil {
		return nil, err
	}

	const (
		contractName       = "OnRamp"
		eventName          = "CCIPSendRequested"
		eventAttributeName = "SequenceNumber"
	)

	seq, err := r.contractReaders[chain].QueryKey(
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
		&cciptypes.CCIPMsg{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query onRamp: %w", err)
	}

	msgs := make([]cciptypes.CCIPMsg, 0)
	for _, item := range seq {
		msg, ok := item.Data.(cciptypes.CCIPMsg)
		if !ok {
			return nil, fmt.Errorf("failed to cast %v to CCIPMsg", item.Data)
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (r *CCIPChainReader) NextSeqNum(
	ctx context.Context, chains []cciptypes.ChainSelector,
) ([]cciptypes.SeqNum, error) {
	if err := r.validateReaderExistence(r.destChain); err != nil {
		return nil, err
	}

	const (
		contractName = "OffRamp"
		funcName     = "getExpectedNextSequenceNumbers"
	)

	seqNums := make([]cciptypes.SeqNum, 0)
	err := r.contractReaders[r.destChain].GetLatestValue(
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

func (r *CCIPChainReader) GasPrices(ctx context.Context, chains []cciptypes.ChainSelector) ([]cciptypes.BigInt, error) {
	if err := r.validateWriterExistence(chains...); err != nil {
		return nil, err
	}

	eg := new(errgroup.Group)
	gasPrices := make([]cciptypes.BigInt, len(chains))
	for i, chain := range chains {
		i, chain := i, chain
		eg.Go(func() error {
			gasPrice, err := r.contractWriters[chain].GetFeeComponents(ctx)
			if err != nil {
				return fmt.Errorf("failed to get gas price: %w", err)
			}
			gasPrices[i] = cciptypes.NewBigInt(gasPrice.ExecutionFee)
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return gasPrices, nil
}

func (r *CCIPChainReader) Close(ctx context.Context) error {
	return nil
}

func (r *CCIPChainReader) validateReaderExistence(chains ...cciptypes.ChainSelector) error {
	for _, ch := range chains {
		_, exists := r.contractReaders[ch]
		if !exists {
			return fmt.Errorf("chain %d: %w", ch, ErrContractReaderNotFound)
		}
	}
	return nil
}

func (r *CCIPChainReader) validateWriterExistence(chains ...cciptypes.ChainSelector) error {
	for _, ch := range chains {
		_, exists := r.contractWriters[ch]
		if !exists {
			return fmt.Errorf("chain %d: %w", ch, ErrContractWriterNotFound)
		}
	}
	return nil
}

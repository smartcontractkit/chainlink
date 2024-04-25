package reader

import "C"
import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/ccipocr3/internal/model"
)

var (
	ErrChainReaderNotFound = errors.New("chain reader not found")
)

type CCIP interface {
	// MsgsAfterTimestamp reads the provided chains.
	// Finds and returns ccip messages submitted after the target time.
	// Messages are sorted ascending based on their timestamp and limited up to the provided limit.
	MsgsAfterTimestamp(ctx context.Context, chains []model.ChainSelector, ts time.Time, limit int) ([]model.CCIPMsg, error)

	// MsgsBetweenSeqNums reads the provided chains.
	// Finds and returns ccip messages submitted between the provided sequence numbers.
	// Messages are sorted ascending based on their timestamp and limited up to the provided limit.
	MsgsBetweenSeqNums(ctx context.Context, chains []model.ChainSelector, seqNumRange model.SeqNumRange) ([]model.CCIPMsg, error)

	// NextSeqNum reads the destination chain.
	// Returns the next expected sequence number for each one of the provided chains.
	NextSeqNum(ctx context.Context, chains []model.ChainSelector) (seqNum []model.SeqNum, err error)
}

type ChainReader interface{} // TODO: Imported from chainlink-common

type CCIPChainReader struct {
	chainReaders map[model.ChainSelector]ChainReader
	destChain    model.ChainSelector
}

func (r *CCIPChainReader) MsgsAfterTimestamp(ctx context.Context, chains []model.ChainSelector, ts time.Time, limit int) ([]model.CCIPMsg, error) {
	if err := r.validateReaderExistence(chains...); err != nil {
		return nil, err
	}
	panic("implement me")
}

func (r *CCIPChainReader) MsgsBetweenSeqNums(ctx context.Context, chains []model.ChainSelector, seqNumRange model.SeqNumRange) ([]model.CCIPMsg, error) {
	if err := r.validateReaderExistence(chains...); err != nil {
		return nil, err
	}
	panic("implement me")
}

func (r *CCIPChainReader) NextSeqNum(ctx context.Context, chains []model.ChainSelector) (seqNum []model.SeqNum, err error) {
	if err := r.validateReaderExistence(r.destChain); err != nil {
		return nil, err
	}
	panic("implement me")
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

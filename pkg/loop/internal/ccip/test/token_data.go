package test

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"time"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

var TokenDataReader = staticTokenDataReader{
	staticTokenDataReaderConfig{
		readTokenDataRequest: readTokenDataRequest{
			msg: cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: cciptypes.EVM2EVMMessage{
					SequenceNumber:      1,
					GasLimit:            big.NewInt(1),
					Nonce:               1,
					MessageID:           ccip.Hash{1},
					SourceChainSelector: 1,
					Sender:              ccip.Address("token data reader sender"),
					Receiver:            ccip.Address("token data reader receiver"),
					Strict:              true,
					FeeToken:            ccip.Address("token data reader feeToken"),
					FeeTokenAmount:      big.NewInt(1),
					Data:                []byte("token data reader data"),
					TokenAmounts: []ccip.TokenAmount{
						{
							Token:  ccip.Address("token data reader token"),
							Amount: big.NewInt(1),
						},
					},
					SourceTokenData: [][]byte{
						[]byte("token data reader sourceTokenData"),
						[]byte("token data reader sourceTokenData2"),
					},
				},
				BlockTimestamp: time.Unix(17779, 0).UTC(),
				Executed:       true,
				Finalized:      true,
				LogIndex:       1,
				TxHash:         "0x123",
			},
		},
		readTokenDataResponse: []byte("read token data response"),
	},
}

type TokenDataReaderEvaluator interface {
	cciptypes.TokenDataReader
	testtypes.Evaluator[cciptypes.TokenDataReader]
}
type staticTokenDataReader struct {
	staticTokenDataReaderConfig
}

var _ TokenDataReaderEvaluator = staticTokenDataReader{}

// Close implements ccip.TokenDataReader.
func (s staticTokenDataReader) Close() error {
	return nil
}

// Evaluate implements types_test.Evaluator.
func (s staticTokenDataReader) Evaluate(ctx context.Context, other cciptypes.TokenDataReader) error {
	got, err := other.ReadTokenData(ctx, s.readTokenDataRequest.msg, s.readTokenDataRequest.tokenIndex)
	if err != nil {
		return fmt.Errorf("failed to get other ReadTokenData: %w", err)
	}
	if !reflect.DeepEqual(got, s.readTokenDataResponse) {
		return fmt.Errorf("unexpected token data: wanted %v got %v", s.readTokenDataResponse, got)
	}
	return nil
}

// ReadTokenData implements ccip.TokenDataReader.
func (s staticTokenDataReader) ReadTokenData(ctx context.Context, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta, tokenIndex int) (tokenData []byte, err error) {
	if !reflect.DeepEqual(s.readTokenDataRequest.msg, msg) {
		return nil, fmt.Errorf("unexpected msg: wanted %v got %v", s.readTokenDataRequest.msg, msg)
	}
	if tokenIndex != s.readTokenDataRequest.tokenIndex {
		return nil, fmt.Errorf("unexpected tokenIndex: wanted %v got %v", s.readTokenDataRequest.tokenIndex, tokenIndex)
	}
	return s.readTokenDataResponse, nil
}

var _ cciptypes.TokenDataReader = staticTokenDataReader{}

type staticTokenDataReaderConfig struct {
	readTokenDataRequest  readTokenDataRequest
	readTokenDataResponse []byte
}

type readTokenDataRequest struct {
	msg        cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta
	tokenIndex int
}

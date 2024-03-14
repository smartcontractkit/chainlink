package test

import (
	"context"
	"fmt"
	"math/big"

	"github.com/stretchr/testify/assert"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

var OnRamp = staticOnRamp{
	staticOnRampConfig: staticOnRampConfig{
		addressResponse: ccip.Address("some-address"),
		routerResponse:  ccip.Address("some-router"),
		configResponse: ccip.OnRampDynamicConfig{
			Router:                            "some-router",
			MaxNumberOfTokensPerMsg:           11,
			DestGasOverhead:                   13,
			DestGasPerPayloadByte:             17,
			DestDataAvailabilityOverheadGas:   23,
			DestGasPerDataAvailabilityByte:    29,
			DestDataAvailabilityMultiplierBps: 31,
			PriceRegistry:                     "some-price-registry",
			MaxDataBytes:                      37,
			MaxPerMsgGasLimit:                 41,
		},
		getSendRequestsBetweenSeqNumsResponse: getSendRequestsBetweenSeqNumsResponse{
			EVM2EVMMessageWithTxMeta: []ccip.EVM2EVMMessageWithTxMeta{
				{
					TxMeta: ccip.TxMeta{
						BlockNumber:             1,
						BlockTimestampUnixMilli: 2,
						TxHash:                  "tx-hash",
						LogIndex:                3,
					},
					EVM2EVMMessage: ccip.EVM2EVMMessage{
						SequenceNumber:      5,
						GasLimit:            big.NewInt(7),
						Nonce:               11,
						MessageID:           ccip.Hash{0: 1, 31: 7},
						SourceChainSelector: 13,
						Sender:              "sender",
						Receiver:            "receiver",
						Strict:              true,
						FeeToken:            "fee-token",
						FeeTokenAmount:      big.NewInt(17),
						Data:                []byte{19},
						TokenAmounts: []ccip.TokenAmount{
							{
								Token:  "token-1",
								Amount: big.NewInt(23),
							},
							{
								Token:  "token-2",
								Amount: big.NewInt(29),
							},
						},
					},
				},
			},
		},
		getSendRequestsBetweenSeqNums: getSendRequestsBetweenSeqNums{
			SeqNumMin: 1,
			SeqNumMax: 2,
			Finalized: true,
		},
	},
}

type OnRampEvaluator interface {
	ccip.OnRampReader
	testtypes.Evaluator[ccip.OnRampReader]
}

var _ OnRampEvaluator = staticOnRamp{}

type staticOnRampConfig struct {
	addressResponse ccip.Address
	routerResponse  ccip.Address
	configResponse  ccip.OnRampDynamicConfig
	getSendRequestsBetweenSeqNums
	getSendRequestsBetweenSeqNumsResponse
}

type staticOnRamp struct {
	staticOnRampConfig
}

// Address implements OnRampEvaluator.
func (s staticOnRamp) Address() (ccip.Address, error) {
	return s.addressResponse, nil
}

// Evaluate implements OnRampEvaluator.
func (s staticOnRamp) Evaluate(ctx context.Context, other ccip.OnRampReader) error {
	address, err := other.Address()
	if err != nil {
		return fmt.Errorf("failed to get address: %w", err)
	}
	if address != s.addressResponse {
		return fmt.Errorf("expected address %s but got %s", s.addressResponse, address)
	}

	router, err := other.RouterAddress()
	if err != nil {
		return fmt.Errorf("failed to get router: %w", err)
	}
	if router != s.routerResponse {
		return fmt.Errorf("expected router %s but got %s", s.routerResponse, router)
	}

	config, err := other.GetDynamicConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}
	if config != s.configResponse {
		return fmt.Errorf("expected config %v but got %v", s.configResponse, config)
	}

	sendRequests, err := other.GetSendRequestsBetweenSeqNums(ctx, s.getSendRequestsBetweenSeqNums.SeqNumMin, s.getSendRequestsBetweenSeqNums.SeqNumMax, s.getSendRequestsBetweenSeqNums.Finalized)
	if err != nil {
		return fmt.Errorf("failed to get send requests: %w", err)
	}
	if !assert.ObjectsAreEqual(s.getSendRequestsBetweenSeqNumsResponse.EVM2EVMMessageWithTxMeta, sendRequests) {
		return fmt.Errorf("expected send requests %v but got %v", s.getSendRequestsBetweenSeqNumsResponse.EVM2EVMMessageWithTxMeta, sendRequests)
	}

	return nil
}

// GetDynamicConfig implements OnRampEvaluator.
func (s staticOnRamp) GetDynamicConfig() (ccip.OnRampDynamicConfig, error) {
	return s.configResponse, nil
}

// GetSendRequestsBetweenSeqNums implements OnRampEvaluator.
func (s staticOnRamp) GetSendRequestsBetweenSeqNums(ctx context.Context, seqNumMin uint64, seqNumMax uint64, finalized bool) ([]ccip.EVM2EVMMessageWithTxMeta, error) {
	if seqNumMin != s.getSendRequestsBetweenSeqNums.SeqNumMin {
		return nil, fmt.Errorf("expected seqNumMin %d but got %d", s.getSendRequestsBetweenSeqNums.SeqNumMin, seqNumMin)
	}
	if seqNumMax != s.getSendRequestsBetweenSeqNums.SeqNumMax {
		return nil, fmt.Errorf("expected seqNumMax %d but got %d", s.getSendRequestsBetweenSeqNums.SeqNumMax, seqNumMax)
	}
	if finalized != s.getSendRequestsBetweenSeqNums.Finalized {
		return nil, fmt.Errorf("expected finalized %t but got %t", s.getSendRequestsBetweenSeqNums.Finalized, finalized)
	}
	return s.getSendRequestsBetweenSeqNumsResponse.EVM2EVMMessageWithTxMeta, nil
}

// RouterAddress implements OnRampEvaluator.
func (s staticOnRamp) RouterAddress() (ccip.Address, error) {
	return s.routerResponse, nil
}

type getSendRequestsBetweenSeqNums struct {
	SeqNumMin uint64
	SeqNumMax uint64
	Finalized bool
}

type getSendRequestsBetweenSeqNumsResponse struct {
	EVM2EVMMessageWithTxMeta []ccip.EVM2EVMMessageWithTxMeta
}

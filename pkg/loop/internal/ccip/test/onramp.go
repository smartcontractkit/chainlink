package test

import (
	"context"
	"fmt"
	"math/big"

	"github.com/stretchr/testify/assert"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// OnRampReader is a static test implementation of [testtypes.Evaluator] for [ccip.OnRampReader].
// The implementation is a simple struct that returns predefined responses.
var OnRampReader = staticOnRamp{
	staticOnRampConfig: staticOnRampConfig{
		addressResponse: ccip.Address("some-address"),
		routerResponse:  ccip.Address("some-router"),
		dynamicConfigResponse: ccip.OnRampDynamicConfig{
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

// ensure the types are satisfied
var _ OnRampEvaluator = staticOnRamp{}

type staticOnRampConfig struct {
	addressResponse              ccip.Address
	routerResponse               ccip.Address
	dynamicConfigResponse        ccip.OnRampDynamicConfig
	isSourceChainHealthyResponse bool
	isSourceCursedResponse       bool
	sourcePriceRegistryResponse  ccip.Address
	getSendRequestsBetweenSeqNums
	getSendRequestsBetweenSeqNumsResponse
}

type staticOnRamp struct {
	staticOnRampConfig
}

// Address implements OnRampEvaluator.
func (s staticOnRamp) Address(context.Context) (ccip.Address, error) {
	return s.addressResponse, nil
}

// Close implements OnRampEvaluator.
func (s staticOnRamp) Close() error {
	return nil
}

// Evaluate implements OnRampEvaluator. It checks that the responses match the expected values.
func (s staticOnRamp) Evaluate(ctx context.Context, other ccip.OnRampReader) error {
	address, err := other.Address(ctx)
	if err != nil {
		return fmt.Errorf("failed to get address: %w", err)
	}
	if address != s.addressResponse {
		return fmt.Errorf("expected address %s but got %s", s.addressResponse, address)
	}

	router, err := other.RouterAddress(ctx)
	if err != nil {
		return fmt.Errorf("failed to get router: %w", err)
	}
	if router != s.routerResponse {
		return fmt.Errorf("expected router %s but got %s", s.routerResponse, router)
	}

	config, err := other.GetDynamicConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}
	if config != s.dynamicConfigResponse {
		return fmt.Errorf("expected config %v but got %v", s.dynamicConfigResponse, config)
	}

	sendRequests, err := other.GetSendRequestsBetweenSeqNums(ctx, s.getSendRequestsBetweenSeqNums.SeqNumMin, s.getSendRequestsBetweenSeqNums.SeqNumMax, s.getSendRequestsBetweenSeqNums.Finalized)
	if err != nil {
		return fmt.Errorf("failed to get send requests: %w", err)
	}
	if !assert.ObjectsAreEqual(s.getSendRequestsBetweenSeqNumsResponse.EVM2EVMMessageWithTxMeta, sendRequests) {
		return fmt.Errorf("expected send requests %v but got %v", s.getSendRequestsBetweenSeqNumsResponse.EVM2EVMMessageWithTxMeta, sendRequests)
	}

	isSourceChainHealthy, err := other.IsSourceChainHealthy(ctx)
	if err != nil {
		return fmt.Errorf("is source chain healthy: %w", err)
	}
	if isSourceChainHealthy != s.isSourceChainHealthyResponse {
		return fmt.Errorf("expected is source chain healthy to be: %v", s.isSourceChainHealthyResponse)
	}

	isSourceCursed, err := other.IsSourceCursed(ctx)
	if err != nil {
		return fmt.Errorf("is source cursed: %w", err)
	}
	if isSourceCursed != s.isSourceCursedResponse {
		return fmt.Errorf("expected is source cursed to be: %v", s.isSourceCursedResponse)
	}

	sourcePriceRegistryAddress, err := other.SourcePriceRegistryAddress(ctx)
	if err != nil {
		return fmt.Errorf("get source price registry address: %w", err)
	}
	if sourcePriceRegistryAddress != s.sourcePriceRegistryResponse {
		return fmt.Errorf("expected source price registry address to be: %v", s.sourcePriceRegistryResponse)
	}

	return nil
}

// GetDynamicConfig implements OnRampEvaluator.
func (s staticOnRamp) GetDynamicConfig(context.Context) (ccip.OnRampDynamicConfig, error) {
	return s.dynamicConfigResponse, nil
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
func (s staticOnRamp) RouterAddress(context.Context) (ccip.Address, error) {
	return s.routerResponse, nil
}

func (s staticOnRamp) IsSourceChainHealthy(ctx context.Context) (bool, error) {
	return s.isSourceChainHealthyResponse, nil
}

func (s staticOnRamp) IsSourceCursed(ctx context.Context) (bool, error) {
	return s.isSourceCursedResponse, nil
}

func (s staticOnRamp) SourcePriceRegistryAddress(ctx context.Context) (ccip.Address, error) {
	return s.sourcePriceRegistryResponse, nil
}

type getSendRequestsBetweenSeqNums struct {
	SeqNumMin uint64
	SeqNumMax uint64
	Finalized bool
}

type getSendRequestsBetweenSeqNumsResponse struct {
	EVM2EVMMessageWithTxMeta []ccip.EVM2EVMMessageWithTxMeta
}

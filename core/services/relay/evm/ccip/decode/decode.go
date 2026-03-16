package decode

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

func DecodeExecReport(ctx context.Context, args abi.Arguments, report []byte) (ccip.ExecReport, error) {
	unpacked, err := args.Unpack(report)
	if err != nil {
		return ccip.ExecReport{}, err
	}
	if len(unpacked) == 0 {
		return ccip.ExecReport{}, errors.New("assumptionViolation: expected at least one element")
	}
	// Must be anonymous struct here
	erStruct, ok := unpacked[0].(struct {
		Messages []struct {
			SourceChainSelector uint64         `json:"sourceChainSelector"`
			Sender              common.Address `json:"sender"`
			Receiver            common.Address `json:"receiver"`
			SequenceNumber      uint64         `json:"sequenceNumber"`
			GasLimit            *big.Int       `json:"gasLimit"`
			Strict              bool           `json:"strict"`
			Nonce               uint64         `json:"nonce"`
			FeeToken            common.Address `json:"feeToken"`
			FeeTokenAmount      *big.Int       `json:"feeTokenAmount"`
			Data                []uint8        `json:"data"`
			TokenAmounts        []struct {
				Token  common.Address `json:"token"`
				Amount *big.Int       `json:"amount"`
			} `json:"tokenAmounts"`
			SourceTokenData [][]uint8 `json:"sourceTokenData"`
			MessageId       [32]uint8 `json:"messageId"`
		} `json:"messages"`
		OffchainTokenData [][][]uint8 `json:"offchainTokenData"`
		Proofs            [][32]uint8 `json:"proofs"`
		ProofFlagBits     *big.Int    `json:"proofFlagBits"`
	})
	if !ok {
		return ccip.ExecReport{}, fmt.Errorf("got %T", unpacked[0])
	}
	messages := make([]ccip.EVM2EVMMessage, 0, len(erStruct.Messages))
	for _, msg := range erStruct.Messages {
		var tokensAndAmounts []ccip.TokenAmount
		for _, tokenAndAmount := range msg.TokenAmounts {
			tokensAndAmounts = append(tokensAndAmounts, ccip.TokenAmount{
				Token:  ccip.Address(tokenAndAmount.Token.String()),
				Amount: tokenAndAmount.Amount,
			})
		}
		messages = append(messages, ccip.EVM2EVMMessage{
			SequenceNumber:      msg.SequenceNumber,
			GasLimit:            msg.GasLimit,
			Nonce:               msg.Nonce,
			MessageID:           msg.MessageId,
			SourceChainSelector: msg.SourceChainSelector,
			Sender:              ccip.Address(msg.Sender.String()),
			Receiver:            ccip.Address(msg.Receiver.String()),
			Strict:              msg.Strict,
			FeeToken:            ccip.Address(msg.FeeToken.String()),
			FeeTokenAmount:      msg.FeeTokenAmount,
			Data:                msg.Data,
			TokenAmounts:        tokensAndAmounts,
			SourceTokenData:     msg.SourceTokenData,
			// TODO: Not needed for plugins, but should be recomputed for consistency.
			// Requires the offramp knowing about onramp version
			Hash: [32]byte{},
		})
	}

	// Unpack will populate with big.Int{false, <allocated empty nat>} for 0 values,
	// which is different from the expected big.NewInt(0). Rebuild to the expected value for this case.
	return ccip.ExecReport{
		Messages:          messages,
		OffchainTokenData: erStruct.OffchainTokenData,
		Proofs:            erStruct.Proofs,
		ProofFlagBits:     new(big.Int).SetBytes(erStruct.ProofFlagBits.Bytes()),
	}, nil
}

func DecodeCommitReport(commitReportArgs abi.Arguments, report []byte) (ccip.CommitStoreReport, error) {
	unpacked, err := commitReportArgs.Unpack(report)
	if err != nil {
		return ccip.CommitStoreReport{}, err
	}
	if len(unpacked) != 1 {
		return ccip.CommitStoreReport{}, errors.New("expected single struct value")
	}

	commitReport, ok := unpacked[0].(struct {
		PriceUpdates struct {
			TokenPriceUpdates []struct {
				SourceToken common.Address `json:"sourceToken"`
				UsdPerToken *big.Int       `json:"usdPerToken"`
			} `json:"tokenPriceUpdates"`
			GasPriceUpdates []struct {
				DestChainSelector uint64   `json:"destChainSelector"`
				UsdPerUnitGas     *big.Int `json:"usdPerUnitGas"`
			} `json:"gasPriceUpdates"`
		} `json:"priceUpdates"`
		Interval struct {
			Min uint64 `json:"min"`
			Max uint64 `json:"max"`
		} `json:"interval"`
		MerkleRoot [32]byte `json:"merkleRoot"`
	})
	if !ok {
		return ccip.CommitStoreReport{}, errors.Errorf("invalid commit report got %T", unpacked[0])
	}

	var tokenPriceUpdates []ccip.TokenPrice
	for _, u := range commitReport.PriceUpdates.TokenPriceUpdates {
		tokenPriceUpdates = append(tokenPriceUpdates, ccip.TokenPrice{
			Token: ccip.Address(u.SourceToken.String()),
			Value: u.UsdPerToken,
		})
	}

	var gasPrices []ccip.GasPrice
	for _, u := range commitReport.PriceUpdates.GasPriceUpdates {
		gasPrices = append(gasPrices, ccip.GasPrice{
			DestChainSelector: u.DestChainSelector,
			Value:             u.UsdPerUnitGas,
		})
	}

	return ccip.CommitStoreReport{
		TokenPrices: tokenPriceUpdates,
		GasPrices:   gasPrices,
		Interval: ccip.CommitStoreInterval{
			Min: commitReport.Interval.Min,
			Max: commitReport.Interval.Max,
		},
		MerkleRoot: commitReport.MerkleRoot,
	}, nil
}

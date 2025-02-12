package ccipsolana

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// MessageHasherV1 implements the MessageHasher interface.
// Compatible with:
// - "OnRamp 1.6.0-dev"
type MessageHasherV1 struct {
	lggr logger.Logger
}

func NewMessageHasherV1(lggr logger.Logger) *MessageHasherV1 {
	return &MessageHasherV1{
		lggr: lggr,
	}
}

// Hash implements the MessageHasher interface.
func (h *MessageHasherV1) Hash(_ context.Context, msg cciptypes.Message) (cciptypes.Bytes32, error) {
	h.lggr.Debugw("hashing message", "msg", msg)

	anyToSolanaMessage := ccip_offramp.Any2SVMRampMessage{}
	anyToSolanaMessage.Header = ccip_offramp.RampMessageHeader{
		SourceChainSelector: uint64(msg.Header.SourceChainSelector),
		DestChainSelector:   uint64(msg.Header.DestChainSelector),
		SequenceNumber:      uint64(msg.Header.SequenceNumber),
		MessageId:           msg.Header.MessageID,
		Nonce:               msg.Header.Nonce,
	}
	anyToSolanaMessage.TokenReceiver = solana.PublicKeyFromBytes(msg.Receiver)
	anyToSolanaMessage.Sender = msg.Sender
	anyToSolanaMessage.Data = msg.Data
	for _, ta := range msg.TokenAmounts {
		destGasAmount, err := extractDestGasAmountFromMap(ta.DestExecDataDecoded)
		if err != nil {
			return [32]byte{}, err
		}

		anyToSolanaMessage.TokenAmounts = append(anyToSolanaMessage.TokenAmounts, ccip_offramp.Any2SVMTokenTransfer{
			SourcePoolAddress: ta.SourcePoolAddress,
			DestTokenAddress:  solana.PublicKeyFromBytes(ta.DestTokenAddress),
			ExtraData:         ta.ExtraData,
			DestGasAmount:     destGasAmount,
			Amount:            ccip_offramp.CrossChainAmount{LeBytes: tokens.ToLittleEndianU256(ta.Amount.Int.Uint64())},
		})
	}

	var err error
	var msgAccounts []solana.PublicKey
	anyToSolanaMessage.ExtraArgs, msgAccounts, err = parseExtraArgsMapWithAccounts(msg.ExtraArgsDecoded)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to decode ExtraArgs: %w", err)
	}

	hash, err := ccip.HashAnyToSVMMessage(anyToSolanaMessage, msg.Header.OnRamp, msgAccounts)
	return [32]byte(hash), err
}

func parseExtraArgsMapWithAccounts(input map[string]any) (ccip_offramp.Any2SVMRampExtraArgs, []solana.PublicKey, error) {
	// Parse input map into SolanaExtraArgs
	var out ccip_offramp.Any2SVMRampExtraArgs
	var accounts []solana.PublicKey

	// Iterate through the expected fields in the struct
	// the field name should match with the one in SVMExtraArgsV1
	// https://github.com/smartcontractkit/chainlink/blob/33c0bda696b0ed97f587a46eacd5c65bed9fb2c1/contracts/src/v0.8/ccip/libraries/Client.sol#L57
	for fieldName, fieldValue := range input {
		lowercase := strings.ToLower(fieldName)
		switch lowercase {
		case "computeunits":
			// Expect uint32
			if v, ok := fieldValue.(uint32); ok {
				out.ComputeUnits = v
			} else {
				return out, accounts, errors.New("invalid type for ComputeUnits, expected uint32")
			}
		case "accountiswritablebitmap":
			// Expect uint64
			if v, ok := fieldValue.(uint64); ok {
				out.IsWritableBitmap = v
			} else {
				return out, accounts, errors.New("invalid type for IsWritableBitmap, expected uint64")
			}
		case "accounts":
			// Expect [][32]byte
			if v, ok := fieldValue.([][32]byte); ok {
				a := make([]solana.PublicKey, len(v))
				for i, val := range v {
					a[i] = solana.PublicKeyFromBytes(val[:])
				}
				accounts = a
			} else {
				return out, accounts, errors.New("invalid type for Accounts, expected [][32]byte")
			}
		default:
			// no error here, aswe only need the keys to construct SVMExtraArgs, other keys can be skipped without
			// return errors because there's no guarantee SVMExtraArgs will match with SVMExtraArgsV1
		}
	}
	return out, accounts, nil
}

// Interface compliance check
var _ cciptypes.MessageHasher = (*MessageHasherV1)(nil)

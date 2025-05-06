package ccipsolana

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

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
	lggr           logger.Logger
	extraDataCodec common.ExtraDataCodec
}

type extraData struct {
	extraArgs     ccip_offramp.Any2SVMRampExtraArgs
	accounts      []solana.PublicKey
	tokenReceiver solana.PublicKey
}

func NewMessageHasherV1(lggr logger.Logger, extraDataCodec common.ExtraDataCodec) *MessageHasherV1 {
	return &MessageHasherV1{
		lggr:           lggr,
		extraDataCodec: extraDataCodec,
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
	if solana.PublicKeyLength != len(msg.Receiver) {
		return [32]byte{}, fmt.Errorf("invalid receiver length: %d", len(msg.Receiver))
	}

	anyToSolanaMessage.Sender = msg.Sender
	anyToSolanaMessage.Data = msg.Data
	for _, ta := range msg.TokenAmounts {
		destExecDataDecodedMap, err := h.extraDataCodec.DecodeTokenAmountDestExecData(ta.DestExecData, msg.Header.SourceChainSelector)
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to decode dest exec data: %w", err)
		}

		destGasAmount, err := extractDestGasAmountFromMap(destExecDataDecodedMap)
		if err != nil {
			return [32]byte{}, err
		}

		if solana.PublicKeyLength != len(ta.DestTokenAddress) {
			return [32]byte{}, fmt.Errorf("invalid DestTokenAddress length: %d", len(ta.DestTokenAddress))
		}
		anyToSolanaMessage.TokenAmounts = append(anyToSolanaMessage.TokenAmounts, ccip_offramp.Any2SVMTokenTransfer{
			SourcePoolAddress: ta.SourcePoolAddress,
			DestTokenAddress:  solana.PublicKeyFromBytes(ta.DestTokenAddress),
			ExtraData:         ta.ExtraData,
			DestGasAmount:     destGasAmount,
			Amount:            ccip_offramp.CrossChainAmount{LeBytes: tokens.ToLittleEndianU256(ta.Amount.Int.Uint64())},
		})
	}

	extraDataDecodedMap, err := h.extraDataCodec.DecodeExtraArgs(msg.ExtraArgs, msg.Header.SourceChainSelector)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to decode extra args: %w", err)
	}

	ed, err := parseExtraDataMap(extraDataDecodedMap)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to decode ExtraArgs: %w", err)
	}

	anyToSolanaMessage.TokenReceiver = ed.tokenReceiver
	anyToSolanaMessage.ExtraArgs = ed.extraArgs
	accounts := ed.accounts
	// if logical receiver is empty, don't prepend it to the accounts list
	if !msg.Receiver.IsZeroOrEmpty() {
		accounts = append([]solana.PublicKey{solana.PublicKeyFromBytes(msg.Receiver)}, accounts...)
	}

	hash, err := ccip.HashAnyToSVMMessage(anyToSolanaMessage, msg.Header.OnRamp, accounts)
	return [32]byte(hash), err
}

func parseExtraDataMap(input map[string]any) (extraData, error) {
	// Parse input map into SolanaExtraArgs
	var out extraData
	var extraArgs ccip_offramp.Any2SVMRampExtraArgs
	var accounts []solana.PublicKey
	var tokenReceiver solana.PublicKey

	// Iterate through the expected fields in the struct
	// the field name should match with the one in SVMExtraArgsV1
	// https://github.com/smartcontractkit/chainlink/blob/33c0bda696b0ed97f587a46eacd5c65bed9fb2c1/contracts/src/v0.8/ccip/libraries/Client.sol#L57
	for fieldName, fieldValue := range input {
		lowercase := strings.ToLower(fieldName)
		switch lowercase {
		case "computeunits":
			// Expect uint32
			if v, ok := fieldValue.(uint32); ok {
				extraArgs.ComputeUnits = v
			} else {
				return out, errors.New("invalid type for ComputeUnits, expected uint32")
			}
		case "accountiswritablebitmap":
			// Expect uint64
			if v, ok := fieldValue.(uint64); ok {
				extraArgs.IsWritableBitmap = v
			} else {
				return out, errors.New("invalid type for IsWritableBitmap, expected uint64")
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
				return out, errors.New("invalid type for Accounts, expected [][32]byte")
			}
		case "tokenreceiver":
			// Expect [32]byte
			v, ok := fieldValue.([32]byte)
			if !ok {
				return out, errors.New("invalid type for TokenReceiver, expected [32]byte")
			}
			tokenReceiver = solana.PublicKeyFromBytes(v[:])
		default:
			// no error here, unneeded keys can be skipped without return errors
		}
	}

	out.extraArgs = extraArgs
	out.accounts = accounts
	out.tokenReceiver = tokenReceiver
	return out, nil
}

func SerializeExtraArgs(tag []byte, data any) ([]byte, error) {
	return ccip.SerializeExtraArgs(data, strings.TrimPrefix(hexutil.Encode(tag), "0x"))
}

// Interface compliance check
var _ cciptypes.MessageHasher = (*MessageHasherV1)(nil)

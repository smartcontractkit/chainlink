package ccipsolana

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// MessageHasherV1 implements the MessageHasher interface.
// Compatible with:
// - "OnRamp 1.6.0-dev"
type MessageHasherV1 struct {
	lggr           logger.Logger
	extraDataCodec ccipocr3.ExtraDataCodecBundle
}

type extraData struct {
	extraArgs     ccip_offramp.Any2SVMRampExtraArgs
	accounts      []solana.PublicKey
	tokenReceiver solana.PublicKey
}

func NewMessageHasherV1(lggr logger.Logger, extraDataCodec ccipocr3.ExtraDataCodecBundle) *MessageHasherV1 {
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
			Amount:            ccip_offramp.CrossChainAmount{LeBytes: [32]uint8(encodeBigIntToFixedLengthLE(ta.Amount.Int, 32))},
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

	// Iterate through the expected fields in the struct.
	// The field names should match SVMExtraArgsV1 from the EVM Client library:
	// https://github.com/smartcontractkit/chainlink/blob/33c0bda696b0ed97f587a46eacd5c65bed9fb2c1/contracts/src/v0.8/ccip/libraries/Client.sol#L57
	//
	// Note: when the source chain ExtraDataCodec runs in a LOOP plugin (e.g. TON relay),
	// the map[string]any values go through gRPC protobuf serialization which changes types:
	//   uint32 -> int64, [32]byte -> []byte, [][32]byte -> []interface{}
	// Each case below handles both the native type and the LOOP-converted type.
	for fieldName, fieldValue := range input {
		lowercase := strings.ToLower(fieldName)
		switch lowercase {
		case "computeunits":
			switch v := fieldValue.(type) {
			case uint32:
				extraArgs.ComputeUnits = v
			case int64: // LOOP gRPC converts uint32 -> int64
				if v < 0 || v > math.MaxUint32 {
					return out, fmt.Errorf("ComputeUnits out of uint32 range: %d", v)
				}
				extraArgs.ComputeUnits = uint32(v)
			default:
				return out, fmt.Errorf("invalid type for ComputeUnits, expected uint32, got %T", fieldValue)
			}
		case "accountiswritablebitmap":
			switch v := fieldValue.(type) {
			case uint64:
				extraArgs.IsWritableBitmap = v
			case int64: // LOOP gRPC may convert uint64 -> int64
				extraArgs.IsWritableBitmap = uint64(v) //nolint:gosec // bitmap, interpret bits as-is
			default:
				return out, fmt.Errorf("invalid type for IsWritableBitmap, expected uint64, got %T", fieldValue)
			}
		case "accounts":
			switch v := fieldValue.(type) {
			case [][32]byte:
				a := make([]solana.PublicKey, len(v))
				for i, val := range v {
					a[i] = solana.PublicKeyFromBytes(val[:])
				}
				accounts = a
			case []interface{}: // LOOP gRPC converts [][32]byte -> []interface{}
				a := make([]solana.PublicKey, len(v))
				for i, elem := range v {
					bs, ok := elem.([]byte)
					if !ok {
						return out, fmt.Errorf("invalid type for Accounts[%d], expected []byte, got %T", i, elem)
					}
					if len(bs) != 32 {
						return out, fmt.Errorf("invalid length for Accounts[%d]: expected 32, got %d", i, len(bs))
					}
					a[i] = solana.PublicKeyFromBytes(bs)
				}
				accounts = a
			case [][]byte: // alternative LOOP representation
				a := make([]solana.PublicKey, len(v))
				for i, bs := range v {
					if len(bs) != 32 {
						return out, fmt.Errorf("invalid length for Accounts[%d]: expected 32, got %d", i, len(bs))
					}
					a[i] = solana.PublicKeyFromBytes(bs)
				}
				accounts = a
			default:
				return out, fmt.Errorf("invalid type for Accounts, expected [][32]byte, got %T", fieldValue)
			}
		case "tokenreceiver":
			switch v := fieldValue.(type) {
			case [32]byte:
				tokenReceiver = solana.PublicKeyFromBytes(v[:])
			case []byte: // LOOP gRPC converts [32]byte -> []byte
				if len(v) != 32 {
					return out, fmt.Errorf("invalid length for TokenReceiver: expected 32, got %d", len(v))
				}
				tokenReceiver = solana.PublicKeyFromBytes(v)
			default:
				return out, fmt.Errorf("invalid type for TokenReceiver, expected [32]byte, got %T", fieldValue)
			}
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

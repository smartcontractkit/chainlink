package rpc

import (
	"context"
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

type GetParsedTransactionOpts struct {
	Commitment CommitmentType `json:"commitment,omitempty"`
}

type GetParsedTransactionResult struct {
	Slot        uint64
	BlockTime   *solana.UnixTimeSeconds
	Transaction *ParsedTransaction
	Meta        *ParsedTransactionMeta
}

func (cl *Client) GetParsedTransaction(
	ctx context.Context,
	txSig solana.Signature,
	opts *GetParsedTransactionOpts,
) (out *GetParsedTransactionResult, err error) {
	params := []interface{}{txSig}
	obj := M{}
	if opts != nil {
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
	}
	obj["encoding"] = solana.EncodingJSONParsed
	params = append(params, obj)
	err = cl.rpcClient.CallForInto(ctx, &out, "getTransaction", params)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, ErrNotFound
	}
	return
}

func (wrap InstructionInfoEnvelope) MarshalJSON() ([]byte, error) {
	if wrap.asString != "" {
		return json.Marshal(wrap.asString)
	}
	return json.Marshal(wrap.asInstructionInfo)
}

func (wrap *InstructionInfoEnvelope) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || (len(data) == 4 && string(data) == "null") {
		// TODO: is this an error?
		return nil
	}

	firstChar := data[0]

	switch firstChar {
	// Check if first character is `[`, standing for a JSON array.
	case '"':
		// It's base64 (or similar)
		{
			err := json.Unmarshal(data, &wrap.asString)
			if err != nil {
				return err
			}
		}
	case '{':
		// It's JSON, most likely.
		{
			return json.Unmarshal(data, &wrap.asInstructionInfo)
		}
	default:
		return fmt.Errorf("Unknown kind: %v", data)
	}

	return nil
}

func (obj GetParsedTransactionResult) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	err = encoder.WriteUint64(obj.Slot, bin.LE)
	if err != nil {
		return err
	}
	{
		if obj.BlockTime == nil {
			err = encoder.WriteBool(false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteBool(true)
			if err != nil {
				return err
			}
			err = encoder.WriteInt64(int64(*obj.BlockTime), bin.LE)
			if err != nil {
				return err
			}
		}
	}
	{
		if obj.Transaction == nil {
			err = encoder.WriteBool(false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteBool(true)
			if err != nil {
				return err
			}
			err = encoder.Encode(obj.Transaction)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (obj GetParsedTransactionResult) UnmarshalWithDecoder(decoder *bin.Decoder) (err error) {
	// Deserialize `Slot`:
	obj.Slot, err = decoder.ReadUint64(bin.LE)
	if err != nil {
		return err
	}
	// Deserialize `BlockTime` (optional):
	{
		ok, err := decoder.ReadBool()
		if err != nil {
			return err
		}
		if ok {
			err = decoder.Decode(&obj.BlockTime)
			if err != nil {
				return err
			}
		}
	}
	{
		ok, err := decoder.ReadBool()
		if err != nil {
			return err
		}
		if ok {
			// NOTE: storing as JSON bytes:
			buf, err := decoder.ReadByteSlice()
			if err != nil {
				return err
			}
			err = json.Unmarshal(buf, &obj.Transaction)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

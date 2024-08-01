// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"context"
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

type GetTransactionOpts struct {
	Encoding solana.EncodingType `json:"encoding,omitempty"`

	// Desired commitment. "processed" is not supported. If parameter not provided, the default is "finalized".
	Commitment CommitmentType `json:"commitment,omitempty"`

	// Max transaction version to return in responses.
	// If the requested block contains a transaction with a higher version, an error will be returned.
	MaxSupportedTransactionVersion *uint64
}

// GetTransaction returns transaction details for a confirmed transaction.
//
// NEW: This method is only available in solana-core v1.7 or newer.
// Please use `getConfirmedTransaction` for solana-core v1.6
func (cl *Client) GetTransaction(
	ctx context.Context,
	txSig solana.Signature, // transaction signature
	opts *GetTransactionOpts,
) (out *GetTransactionResult, err error) {
	params := []interface{}{txSig}
	if opts != nil {
		obj := M{}
		if opts.Encoding != "" {
			if !solana.IsAnyOfEncodingType(
				opts.Encoding,
				// Valid encodings:
				// solana.EncodingJSON, // TODO
				// solana.EncodingJSONParsed, // TODO
				solana.EncodingBase58,
				solana.EncodingBase64,
				solana.EncodingBase64Zstd,
			) {
				return nil, fmt.Errorf("provided encoding is not supported: %s", opts.Encoding)
			}
			obj["encoding"] = opts.Encoding
		}
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		if opts.MaxSupportedTransactionVersion != nil {
			obj["maxSupportedTransactionVersion"] = *opts.MaxSupportedTransactionVersion
		}
		if len(obj) > 0 {
			params = append(params, obj)
		}
	}
	err = cl.rpcClient.CallForInto(ctx, &out, "getTransaction", params)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, ErrNotFound
	}
	return
}

type GetTransactionResult struct {
	// The slot this transaction was processed in.
	Slot uint64 `json:"slot"`

	// Estimated production time, as Unix timestamp (seconds since the Unix epoch)
	// of when the transaction was processed.
	// Nil if not available.
	BlockTime *solana.UnixTimeSeconds `json:"blockTime" bin:"optional"`

	Transaction *TransactionResultEnvelope `json:"transaction" bin:"optional"`
	Meta        *TransactionMeta           `json:"meta,omitempty" bin:"optional"`
	Version     TransactionVersion         `json:"version"`
}

// TransactionResultEnvelope will contain a *solana.Transaction if the requested encoding is `solana.EncodingJSON`
// (which is also the default when the encoding is not specified),
// or a `solana.Data` in case of EncodingBase58, EncodingBase64.
type TransactionResultEnvelope struct {
	asDecodedBinary     solana.Data
	asParsedTransaction *solana.Transaction
}

func (wrap TransactionResultEnvelope) MarshalJSON() ([]byte, error) {
	if wrap.asParsedTransaction != nil {
		return json.Marshal(wrap.asParsedTransaction)
	}
	return json.Marshal(wrap.asDecodedBinary)
}

func (wrap *TransactionResultEnvelope) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || (len(data) == 4 && string(data) == "null") {
		// TODO: is this an error?
		return nil
	}

	firstChar := data[0]

	switch firstChar {
	// Check if first character is `[`, standing for a JSON array.
	case '[':
		// It's base64 (or similar)
		{
			err := wrap.asDecodedBinary.UnmarshalJSON(data)
			if err != nil {
				return err
			}
		}
	case '{':
		// It's JSON, most likely.
		{
			return json.Unmarshal(data, &wrap.asParsedTransaction)
		}
	default:
		return fmt.Errorf("Unknown kind: %v", data)
	}

	return nil
}

// GetBinary returns the decoded bytes if the encoding is
// "base58", "base64".
func (dt *TransactionResultEnvelope) GetBinary() []byte {
	return dt.asDecodedBinary.Content
}

func (dt *TransactionResultEnvelope) GetData() solana.Data {
	return dt.asDecodedBinary
}

// GetRawJSON returns a *solana.Transaction when the data
// encoding is EncodingJSON.
func (dt *TransactionResultEnvelope) GetTransaction() (*solana.Transaction, error) {
	if dt.asDecodedBinary.Content != nil {
		tx := new(solana.Transaction)
		err := tx.UnmarshalWithDecoder(bin.NewBinDecoder(dt.asDecodedBinary.Content))
		if err != nil {
			return nil, err
		}
		return tx, nil
	}
	return dt.asParsedTransaction, nil
}

func (obj TransactionResultEnvelope) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	return encoder.Encode(obj.asDecodedBinary)
}

func (obj *TransactionResultEnvelope) UnmarshalWithDecoder(decoder *bin.Decoder) (err error) {
	return decoder.Decode(&obj.asDecodedBinary)
}

func (obj GetTransactionResult) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
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
			err = obj.Transaction.MarshalWithEncoder(encoder)
			if err != nil {
				return err
			}
		}
	}
	{
		if obj.Meta == nil {
			err = encoder.WriteBool(false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteBool(true)
			if err != nil {
				return err
			}
			// NOTE: storing as JSON bytes:
			buf, err := json.Marshal(obj.Meta)
			if err != nil {
				return err
			}
			err = encoder.WriteBytes(buf, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (obj *GetTransactionResult) UnmarshalWithDecoder(decoder *bin.Decoder) (err error) {
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
			obj.Transaction = new(TransactionResultEnvelope)
			err = obj.Transaction.UnmarshalWithDecoder(decoder)
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
			err = json.Unmarshal(buf, &obj.Meta)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

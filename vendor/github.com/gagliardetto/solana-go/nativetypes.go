// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
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

package solana

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"io"

	bin "github.com/gagliardetto/binary"
	"github.com/mostynb/zstdpool-freelist"
	"github.com/mr-tron/base58"
)

type Padding []byte

type Hash PublicKey

func MustHashFromBase58(in string) Hash {
	return Hash(MustPublicKeyFromBase58(in))
}

func HashFromBase58(in string) (Hash, error) {
	tmp, err := PublicKeyFromBase58(in)
	if err != nil {
		return Hash{}, err
	}
	return Hash(tmp), nil
}

func HashFromBytes(in []byte) Hash {
	return Hash(PublicKeyFromBytes(in))
}

func (ha Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(base58.Encode(ha[:]))
}

func (ha *Hash) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	tmp, err := PublicKeyFromBase58(s)
	if err != nil {
		return fmt.Errorf("invalid public key %q: %w", s, err)
	}
	*ha = Hash(tmp)
	return
}

func (ha Hash) Equals(pb Hash) bool {
	return ha == pb
}

var zeroHash = Hash{}

func (ha Hash) IsZero() bool {
	return ha == zeroHash
}

func (ha Hash) String() string {
	return base58.Encode(ha[:])
}

type Signature [64]byte

var zeroSignature = Signature{}

func (sig Signature) IsZero() bool {
	return sig == zeroSignature
}
func (sig Signature) Equals(pb Signature) bool {
	return sig == pb
}

func SignatureFromBase58(in string) (out Signature, err error) {
	val, err := base58.Decode(in)
	if err != nil {
		return
	}

	if len(val) != SignatureLength {
		err = fmt.Errorf("invalid length, expected 64, got %d", len(val))
		return
	}
	copy(out[:], val)
	return
}

func MustSignatureFromBase58(in string) Signature {
	out, err := SignatureFromBase58(in)
	if err != nil {
		panic(err)
	}
	return out
}

func SignatureFromBytes(in []byte) (out Signature) {
	byteCount := len(in)
	if byteCount == 0 {
		return
	}

	max := SignatureLength
	if byteCount < max {
		max = byteCount
	}

	copy(out[:], in[0:max])
	return
}

func (p Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(base58.Encode(p[:]))
}

func (p *Signature) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	dat, err := base58.Decode(s)
	if err != nil {
		return err
	}

	if len(dat) != SignatureLength {
		return fmt.Errorf("invalid length for Signature, expected 64, got %d", len(dat))
	}

	target := Signature{}
	copy(target[:], dat)
	*p = target
	return
}

func (s Signature) Verify(pubkey PublicKey, msg []byte) bool {
	return ed25519.Verify(pubkey[:], msg, s[:])
}

func (p Signature) String() string {
	return base58.Encode(p[:])
}

type Base58 []byte

func (t Base58) MarshalJSON() ([]byte, error) {
	return json.Marshal(base58.Encode(t))
}

func (t *Base58) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	if s == "" {
		*t = []byte{}
		return nil
	}

	*t, err = base58.Decode(s)
	return
}

func (t Base58) String() string {
	return base58.Encode(t)
}

type Data struct {
	Content  []byte
	Encoding EncodingType
}

func (t Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		[]interface{}{
			t.String(),
			t.Encoding,
		})
}

var zstdDecoderPool = zstdpool.NewDecoderPool()

func (t *Data) UnmarshalJSON(data []byte) (err error) {
	var in []string
	if err := json.Unmarshal(data, &in); err != nil {
		return err
	}

	if len(in) != 2 {
		return fmt.Errorf("invalid length for solana.Data, expected 2, found %d", len(in))
	}

	contentString := in[0]
	encodingString := in[1]
	t.Encoding = EncodingType(encodingString)

	if contentString == "" {
		t.Content = []byte{}
		return nil
	}

	switch t.Encoding {
	case EncodingBase58:
		var err error
		t.Content, err = base58.Decode(contentString)
		if err != nil {
			return err
		}
	case EncodingBase64:
		var err error
		t.Content, err = base64.StdEncoding.DecodeString(contentString)
		if err != nil {
			return err
		}
	case EncodingBase64Zstd:
		rawBytes, err := base64.StdEncoding.DecodeString(contentString)
		if err != nil {
			return err
		}
		dec, err := zstdDecoderPool.Get(nil)
		if err != nil {
			return err
		}
		defer zstdDecoderPool.Put(dec)

		t.Content, err = dec.DecodeAll(rawBytes, nil)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported encoding %s", encodingString)
	}
	return
}

var zstdEncoderPool = zstdpool.NewEncoderPool()

func (t Data) String() string {
	switch EncodingType(t.Encoding) {
	case EncodingBase58:
		return base58.Encode(t.Content)
	case EncodingBase64:
		return base64.StdEncoding.EncodeToString(t.Content)
	case EncodingBase64Zstd:
		enc, err := zstdEncoderPool.Get(nil)
		if err != nil {
			// TODO: remove panic?
			panic(err)
		}
		defer zstdEncoderPool.Put(enc)
		return base64.StdEncoding.EncodeToString(enc.EncodeAll(t.Content, nil))
	default:
		// TODO
		return ""
	}
}

func (obj Data) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	err = encoder.WriteBytes(obj.Content, true)
	if err != nil {
		return err
	}
	err = encoder.WriteString(string(obj.Encoding))
	if err != nil {
		return err
	}
	return nil
}

func (obj *Data) UnmarshalWithDecoder(decoder *bin.Decoder) (err error) {
	obj.Content, err = decoder.ReadByteSlice()
	if err != nil {
		return err
	}
	{
		enc, err := decoder.ReadString()
		if err != nil {
			return err
		}
		obj.Encoding = EncodingType(enc)
	}
	return nil
}

type ByteWrapper struct {
	io.Reader
}

func (w *ByteWrapper) ReadByte() (byte, error) {
	var b [1]byte
	// NOTE: w.Read() gives no guaranties about the number of bytes actually read.
	// Using io.ReadFull reads exactly len(buf) bytes from r into buf.
	_, err := io.ReadFull(w, b[:])
	return b[0], err
}

type EncodingType string

const (
	EncodingBase58     EncodingType = "base58"      // limited to Account data of less than 129 bytes
	EncodingBase64     EncodingType = "base64"      // will return base64 encoded data for Account data of any size
	EncodingBase64Zstd EncodingType = "base64+zstd" // compresses the Account data using Zstandard and base64-encodes the result

	// attempts to use program-specific state parsers to
	// return more human-readable and explicit account state data.
	// If "jsonParsed" is requested but a parser cannot be found,
	// the field falls back to "base64" encoding, detectable when the data field is type <string>.
	// Cannot be used if specifying dataSlice parameters (offset, length).
	EncodingJSONParsed EncodingType = "jsonParsed"

	EncodingJSON EncodingType = "json" // NOTE: you're probably looking for EncodingJSONParsed
)

// IsAnyOfEncodingType checks whether the provided `candidate` is any of the `allowed`.
func IsAnyOfEncodingType(candidate EncodingType, allowed ...EncodingType) bool {
	for _, v := range allowed {
		if candidate == v {
			return true
		}
	}
	return false
}

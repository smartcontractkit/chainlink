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

package bin

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type SafeString string

func (ss SafeString) MarshalWithEncoder(encoder *Encoder) error {
	return encoder.WriteString(string(ss))
}

func (ss *SafeString) UnmarshalWithDecoder(d *Decoder) error {
	s, e := d.SafeReadUTF8String()
	if e != nil {
		return e
	}

	*ss = SafeString(s)
	return nil
}

type Bool bool

func (b *Bool) UnmarshalJSON(data []byte) error {
	var num int
	err := json.Unmarshal(data, &num)
	if err == nil {
		*b = Bool(num != 0)
		return nil
	}

	var boolVal bool
	if err := json.Unmarshal(data, &boolVal); err != nil {
		return fmt.Errorf("couldn't unmarshal bool as int or true/false: %s", err)
	}

	*b = Bool(boolVal)
	return nil
}

func (b *Bool) UnmarshalWithDecoder(decoder *Decoder) error {
	value, err := decoder.ReadBool()
	if err != nil {
		return err
	}

	*b = Bool(value)
	return nil
}

func (b Bool) MarshalWithEncoder(encoder *Encoder) error {
	return encoder.WriteBool(bool(b))
}

type HexBytes []byte

func (t HexBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(t))
}

func (t *HexBytes) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	*t, err = hex.DecodeString(s)
	return
}

func (t HexBytes) String() string {
	return hex.EncodeToString(t)
}

func (o *HexBytes) UnmarshalWithDecoder(decoder *Decoder) error {
	value, err := decoder.ReadByteSlice()
	if err != nil {
		return fmt.Errorf("hex bytes: %s", err)
	}

	*o = HexBytes(value)
	return nil
}

func (o HexBytes) MarshalWithEncoder(encoder *Encoder) error {
	return encoder.WriteBytes([]byte(o), true)
}

type Varint16 int16

func (o *Varint16) UnmarshalWithDecoder(decoder *Decoder) error {
	value, err := decoder.ReadVarint16()
	if err != nil {
		return fmt.Errorf("varint16: %s", err)
	}

	*o = Varint16(value)
	return nil
}

func (o Varint16) MarshalWithEncoder(encoder *Encoder) error {
	return encoder.WriteVarInt(int(o))
}

type Varuint16 uint16

func (o *Varuint16) UnmarshalWithDecoder(decoder *Decoder) error {
	value, err := decoder.ReadUvarint16()
	if err != nil {
		return fmt.Errorf("varuint16: %s", err)
	}

	*o = Varuint16(value)
	return nil
}

func (o Varuint16) MarshalWithEncoder(encoder *Encoder) error {
	return encoder.WriteUVarInt(int(o))
}

type Varuint32 uint32

func (o *Varuint32) UnmarshalWithDecoder(decoder *Decoder) error {
	value, err := decoder.ReadUvarint64()
	if err != nil {
		return fmt.Errorf("varuint32: %s", err)
	}

	*o = Varuint32(value)
	return nil
}

func (o Varuint32) MarshalWithEncoder(encoder *Encoder) error {
	return encoder.WriteUVarInt(int(o))
}

type Varint32 int32

func (o *Varint32) UnmarshalWithDecoder(decoder *Decoder) error {
	value, err := decoder.ReadVarint32()
	if err != nil {
		return err
	}

	*o = Varint32(value)
	return nil
}

func (o Varint32) MarshalWithEncoder(encoder *Encoder) error {
	return encoder.WriteVarInt(int(o))
}

type JSONFloat64 float64

func (f *JSONFloat64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}

		*f = JSONFloat64(val)

		return nil
	}

	var fl float64
	if err := json.Unmarshal(data, &fl); err != nil {
		return err
	}

	*f = JSONFloat64(fl)

	return nil
}

func (f *JSONFloat64) UnmarshalWithDecoder(dec *Decoder) error {
	value, err := dec.ReadFloat64(dec.currentFieldOpt.Order)
	if err != nil {
		return err
	}

	*f = JSONFloat64(value)
	return nil
}

func (f JSONFloat64) MarshalWithEncoder(enc *Encoder) error {
	return enc.WriteFloat64(float64(f), enc.currentFieldOpt.Order)
}

type Int64 int64

func (i Int64) MarshalJSON() (data []byte, err error) {
	if i > 0xffffffff || i < -0xffffffff {
		encodedInt, err := json.Marshal(int64(i))
		if err != nil {
			return nil, err
		}
		data = append([]byte{'"'}, encodedInt...)
		data = append(data, '"')
		return data, nil
	}
	return json.Marshal(int64(i))
}

func (i *Int64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}

		*i = Int64(val)

		return nil
	}

	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = Int64(v)

	return nil
}

func (i *Int64) UnmarshalWithDecoder(dec *Decoder) error {
	value, err := dec.ReadInt64(dec.currentFieldOpt.Order)
	if err != nil {
		return err
	}

	*i = Int64(value)
	return nil
}

func (i Int64) MarshalWithEncoder(enc *Encoder) error {
	return enc.WriteInt64(int64(i), enc.currentFieldOpt.Order)
}

type Uint64 uint64

func (i Uint64) MarshalJSON() (data []byte, err error) {
	if i > 0xffffffff {
		encodedInt, err := json.Marshal(uint64(i))
		if err != nil {
			return nil, err
		}
		data = append([]byte{'"'}, encodedInt...)
		data = append(data, '"')
		return data, nil
	}
	return json.Marshal(uint64(i))
}

func (i *Uint64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		val, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}

		*i = Uint64(val)

		return nil
	}

	var v uint64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = Uint64(v)

	return nil
}

func (i *Uint64) UnmarshalWithDecoder(dec *Decoder) error {
	value, err := dec.ReadUint64(dec.currentFieldOpt.Order)
	if err != nil {
		return err
	}

	*i = Uint64(value)
	return nil
}

func (i Uint64) MarshalWithEncoder(enc *Encoder) error {
	return enc.WriteUint64(uint64(i), enc.currentFieldOpt.Order)
}

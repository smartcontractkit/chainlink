package types

import (
	"math/big"
	"reflect"

	"github.com/fxamacker/cbor/v2"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type SizedBigInt interface {
	Verify() error
	private()
}

var sizedBigIntType = reflect.TypeOf((*SizedBigInt)(nil)).Elem()

func SizedBigIntType() reflect.Type {
	return sizedBigIntType
}

type int24 big.Int

func (i *int24) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int24) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int24) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int24) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int24) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(24, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int24) private() {}

func init() {
	typeMap["int24"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int24)(nil)),
	}
}

type uint24 big.Int

func (i *uint24) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint24) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint24) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint24) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint24) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 24 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint24) private() {}

func init() {
	typeMap["uint24"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint24)(nil)),
	}
}

type int40 big.Int

func (i *int40) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int40) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int40) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int40) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int40) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(40, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int40) private() {}

func init() {
	typeMap["int40"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int40)(nil)),
	}
}

type uint40 big.Int

func (i *uint40) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint40) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint40) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint40) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint40) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 40 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint40) private() {}

func init() {
	typeMap["uint40"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint40)(nil)),
	}
}

type int48 big.Int

func (i *int48) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int48) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int48) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int48) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int48) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(48, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int48) private() {}

func init() {
	typeMap["int48"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int48)(nil)),
	}
}

type uint48 big.Int

func (i *uint48) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint48) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint48) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint48) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint48) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 48 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint48) private() {}

func init() {
	typeMap["uint48"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint48)(nil)),
	}
}

type int56 big.Int

func (i *int56) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int56) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int56) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int56) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int56) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(56, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int56) private() {}

func init() {
	typeMap["int56"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int56)(nil)),
	}
}

type uint56 big.Int

func (i *uint56) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint56) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint56) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint56) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint56) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 56 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint56) private() {}

func init() {
	typeMap["uint56"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint56)(nil)),
	}
}

type int72 big.Int

func (i *int72) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int72) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int72) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int72) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int72) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(72, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int72) private() {}

func init() {
	typeMap["int72"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int72)(nil)),
	}
}

type uint72 big.Int

func (i *uint72) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint72) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint72) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint72) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint72) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 72 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint72) private() {}

func init() {
	typeMap["uint72"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint72)(nil)),
	}
}

type int80 big.Int

func (i *int80) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int80) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int80) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int80) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int80) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(80, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int80) private() {}

func init() {
	typeMap["int80"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int80)(nil)),
	}
}

type uint80 big.Int

func (i *uint80) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint80) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint80) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint80) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint80) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 80 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint80) private() {}

func init() {
	typeMap["uint80"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint80)(nil)),
	}
}

type int88 big.Int

func (i *int88) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int88) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int88) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int88) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int88) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(88, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int88) private() {}

func init() {
	typeMap["int88"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int88)(nil)),
	}
}

type uint88 big.Int

func (i *uint88) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint88) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint88) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint88) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint88) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 88 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint88) private() {}

func init() {
	typeMap["uint88"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint88)(nil)),
	}
}

type int96 big.Int

func (i *int96) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int96) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int96) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int96) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int96) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(96, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int96) private() {}

func init() {
	typeMap["int96"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int96)(nil)),
	}
}

type uint96 big.Int

func (i *uint96) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint96) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint96) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint96) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint96) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 96 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint96) private() {}

func init() {
	typeMap["uint96"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint96)(nil)),
	}
}

type int104 big.Int

func (i *int104) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int104) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int104) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int104) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int104) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(104, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int104) private() {}

func init() {
	typeMap["int104"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int104)(nil)),
	}
}

type uint104 big.Int

func (i *uint104) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint104) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint104) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint104) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint104) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 104 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint104) private() {}

func init() {
	typeMap["uint104"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint104)(nil)),
	}
}

type int112 big.Int

func (i *int112) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int112) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int112) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int112) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int112) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(112, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int112) private() {}

func init() {
	typeMap["int112"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int112)(nil)),
	}
}

type uint112 big.Int

func (i *uint112) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint112) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint112) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint112) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint112) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 112 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint112) private() {}

func init() {
	typeMap["uint112"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint112)(nil)),
	}
}

type int120 big.Int

func (i *int120) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int120) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int120) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int120) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int120) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(120, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int120) private() {}

func init() {
	typeMap["int120"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int120)(nil)),
	}
}

type uint120 big.Int

func (i *uint120) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint120) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint120) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint120) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint120) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 120 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint120) private() {}

func init() {
	typeMap["uint120"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint120)(nil)),
	}
}

type int128 big.Int

func (i *int128) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int128) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int128) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int128) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int128) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(128, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int128) private() {}

func init() {
	typeMap["int128"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int128)(nil)),
	}
}

type uint128 big.Int

func (i *uint128) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint128) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint128) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint128) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint128) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 128 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint128) private() {}

func init() {
	typeMap["uint128"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint128)(nil)),
	}
}

type int136 big.Int

func (i *int136) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int136) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int136) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int136) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int136) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(136, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int136) private() {}

func init() {
	typeMap["int136"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int136)(nil)),
	}
}

type uint136 big.Int

func (i *uint136) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint136) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint136) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint136) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint136) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 136 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint136) private() {}

func init() {
	typeMap["uint136"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint136)(nil)),
	}
}

type int144 big.Int

func (i *int144) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int144) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int144) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int144) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int144) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(144, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int144) private() {}

func init() {
	typeMap["int144"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int144)(nil)),
	}
}

type uint144 big.Int

func (i *uint144) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint144) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint144) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint144) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint144) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 144 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint144) private() {}

func init() {
	typeMap["uint144"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint144)(nil)),
	}
}

type int152 big.Int

func (i *int152) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int152) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int152) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int152) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int152) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(152, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int152) private() {}

func init() {
	typeMap["int152"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int152)(nil)),
	}
}

type uint152 big.Int

func (i *uint152) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint152) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint152) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint152) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint152) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 152 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint152) private() {}

func init() {
	typeMap["uint152"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint152)(nil)),
	}
}

type int160 big.Int

func (i *int160) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int160) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int160) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int160) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int160) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(160, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int160) private() {}

func init() {
	typeMap["int160"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int160)(nil)),
	}
}

type uint160 big.Int

func (i *uint160) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint160) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint160) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint160) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint160) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 160 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint160) private() {}

func init() {
	typeMap["uint160"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint160)(nil)),
	}
}

type int168 big.Int

func (i *int168) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int168) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int168) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int168) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int168) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(168, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int168) private() {}

func init() {
	typeMap["int168"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int168)(nil)),
	}
}

type uint168 big.Int

func (i *uint168) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint168) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint168) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint168) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint168) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 168 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint168) private() {}

func init() {
	typeMap["uint168"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint168)(nil)),
	}
}

type int176 big.Int

func (i *int176) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int176) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int176) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int176) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int176) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(176, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int176) private() {}

func init() {
	typeMap["int176"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int176)(nil)),
	}
}

type uint176 big.Int

func (i *uint176) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint176) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint176) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint176) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint176) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 176 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint176) private() {}

func init() {
	typeMap["uint176"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint176)(nil)),
	}
}

type int184 big.Int

func (i *int184) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int184) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int184) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int184) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int184) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(184, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int184) private() {}

func init() {
	typeMap["int184"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int184)(nil)),
	}
}

type uint184 big.Int

func (i *uint184) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint184) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint184) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint184) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint184) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 184 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint184) private() {}

func init() {
	typeMap["uint184"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint184)(nil)),
	}
}

type int192 big.Int

func (i *int192) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int192) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int192) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int192) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int192) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(192, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int192) private() {}

func init() {
	typeMap["int192"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int192)(nil)),
	}
}

type uint192 big.Int

func (i *uint192) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint192) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint192) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint192) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint192) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 192 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint192) private() {}

func init() {
	typeMap["uint192"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint192)(nil)),
	}
}

type int200 big.Int

func (i *int200) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int200) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int200) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int200) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int200) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(200, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int200) private() {}

func init() {
	typeMap["int200"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int200)(nil)),
	}
}

type uint200 big.Int

func (i *uint200) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint200) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint200) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint200) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint200) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 200 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint200) private() {}

func init() {
	typeMap["uint200"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint200)(nil)),
	}
}

type int208 big.Int

func (i *int208) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int208) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int208) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int208) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int208) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(208, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int208) private() {}

func init() {
	typeMap["int208"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int208)(nil)),
	}
}

type uint208 big.Int

func (i *uint208) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint208) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint208) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint208) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint208) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 208 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint208) private() {}

func init() {
	typeMap["uint208"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint208)(nil)),
	}
}

type int216 big.Int

func (i *int216) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int216) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int216) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int216) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int216) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(216, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int216) private() {}

func init() {
	typeMap["int216"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int216)(nil)),
	}
}

type uint216 big.Int

func (i *uint216) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint216) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint216) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint216) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint216) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 216 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint216) private() {}

func init() {
	typeMap["uint216"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint216)(nil)),
	}
}

type int224 big.Int

func (i *int224) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int224) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int224) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int224) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int224) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(224, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int224) private() {}

func init() {
	typeMap["int224"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int224)(nil)),
	}
}

type uint224 big.Int

func (i *uint224) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint224) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint224) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint224) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint224) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 224 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint224) private() {}

func init() {
	typeMap["uint224"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint224)(nil)),
	}
}

type int232 big.Int

func (i *int232) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int232) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int232) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int232) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int232) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(232, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int232) private() {}

func init() {
	typeMap["int232"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int232)(nil)),
	}
}

type uint232 big.Int

func (i *uint232) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint232) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint232) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint232) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint232) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 232 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint232) private() {}

func init() {
	typeMap["uint232"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint232)(nil)),
	}
}

type int240 big.Int

func (i *int240) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int240) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int240) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int240) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int240) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(240, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int240) private() {}

func init() {
	typeMap["int240"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int240)(nil)),
	}
}

type uint240 big.Int

func (i *uint240) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint240) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint240) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint240) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint240) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 240 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint240) private() {}

func init() {
	typeMap["uint240"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint240)(nil)),
	}
}

type int248 big.Int

func (i *int248) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int248) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int248) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int248) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int248) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(248, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int248) private() {}

func init() {
	typeMap["int248"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int248)(nil)),
	}
}

type uint248 big.Int

func (i *uint248) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint248) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint248) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint248) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint248) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 248 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint248) private() {}

func init() {
	typeMap["uint248"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint248)(nil)),
	}
}

type int256 big.Int

func (i *int256) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *int256) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *int256) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *int256) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *int256) Verify() error {
	bi := (*big.Int)(i)

	if !codec.FitsInNBitsSigned(256, bi) {
		return types.ErrInvalidType
	}

	return nil
}

func (i *int256) private() {}

func init() {
	typeMap["int256"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*int256)(nil)),
	}
}

type uint256 big.Int

func (i *uint256) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *uint256) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal((*big.Int)(i))
}

func (i *uint256) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.ErrInvalidType
	}

	return i.Verify()
}

func (i *uint256) MarshalText() ([]byte, error) {
	bi := (*big.Int)(i)
	return []byte(bi.String()), nil
}

func (i *uint256) Verify() error {
	bi := (*big.Int)(i)

	if bi.BitLen() > 256 || bi.Sign() < 0 {
		return types.ErrInvalidType
	}

	return nil
}

func (i *uint256) private() {}

func init() {
	typeMap["uint256"] = &ABIEncodingType{
		native:  reflect.TypeOf((*big.Int)(nil)),
		checked: reflect.TypeOf((*uint256)(nil)),
	}
}

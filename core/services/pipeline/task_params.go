package pipeline

import (
	"encoding/hex"
	"encoding/json"
	"math"
	"math/big"
	"net/url"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name PipelineParamUnmarshaler --output ./mocks/ --case=underscore

type PipelineParamUnmarshaler interface {
	UnmarshalPipelineParam(val interface{}) error
}

func ResolveParam(out PipelineParamUnmarshaler, getters []GetterFunc) error {
	var val interface{}
	var err error
	var found bool
	for _, get := range getters {
		val, err = get()
		if errors.Is(errors.Cause(err), ErrParameterEmpty) {
			continue
		} else if err != nil {
			return err
		}
		found = true
		break
	}
	if !found {
		return ErrParameterEmpty
	}

	err = out.UnmarshalPipelineParam(val)
	if err != nil {
		return err
	}
	return nil
}

type StringParam string

func (s *StringParam) UnmarshalPipelineParam(val interface{}) error {
	switch v := val.(type) {
	case string:
		*s = StringParam(v)
		return nil
	case []byte:
		*s = StringParam(string(v))
		return nil
	case ObjectParam:
		if v.Type == StringType {
			*s = v.StringValue
			return nil
		}
	case *ObjectParam:
		if v.Type == StringType {
			*s = v.StringValue
			return nil
		}
	}
	return errors.Wrapf(ErrBadInput, "expected string, got %T", val)
}

type BytesParam []byte

func (b *BytesParam) UnmarshalPipelineParam(val interface{}) error {
	switch v := val.(type) {
	case string:
		// first check if this is a valid hex-encoded string
		if utils.HasHexPrefix(v) {
			noHexPrefix := utils.RemoveHexPrefix(v)
			bs, err := hex.DecodeString(noHexPrefix)
			if err == nil {
				*b = bs
				return nil
			}
		}
		*b = BytesParam(v)
		return nil
	case []byte:
		*b = v
		return nil
	case nil:
		*b = BytesParam(nil)
		return nil
	case ObjectParam:
		if v.Type == StringType {
			*b = BytesParam(v.StringValue)
			return nil
		}
	case *ObjectParam:
		if v.Type == StringType {
			*b = BytesParam(v.StringValue)
			return nil
		}
	}

	return errors.Wrapf(ErrBadInput, "expected array of bytes, got %T", val)
}

type Uint64Param uint64

func (u *Uint64Param) UnmarshalPipelineParam(val interface{}) error {
	switch v := val.(type) {
	case uint:
		*u = Uint64Param(v)
	case uint8:
		*u = Uint64Param(v)
	case uint16:
		*u = Uint64Param(v)
	case uint32:
		*u = Uint64Param(v)
	case uint64:
		*u = Uint64Param(v)
	case int:
		*u = Uint64Param(v)
	case int8:
		*u = Uint64Param(v)
	case int16:
		*u = Uint64Param(v)
	case int32:
		*u = Uint64Param(v)
	case int64:
		*u = Uint64Param(v)
	case float64: // when decoding from db: JSON numbers are floats
		*u = Uint64Param(v)
	case string:
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return errors.Wrap(ErrBadInput, err.Error())
		}
		*u = Uint64Param(n)
	default:
		return errors.Wrapf(ErrBadInput, "expected unsiend integer, got %T", val)
	}
	return nil
}

type MaybeUint64Param struct {
	n     uint64
	isSet bool
}

// NewMaybeUint64Param creates new instance of MaybeUint64Param
func NewMaybeUint64Param(n uint64, isSet bool) MaybeUint64Param {
	return MaybeUint64Param{
		n:     n,
		isSet: isSet,
	}
}

func (p *MaybeUint64Param) UnmarshalPipelineParam(val interface{}) error {
	var n uint64
	switch v := val.(type) {
	case uint:
		n = uint64(v)
	case uint8:
		n = uint64(v)
	case uint16:
		n = uint64(v)
	case uint32:
		n = uint64(v)
	case uint64:
		n = v
	case int:
		n = uint64(v)
	case int8:
		n = uint64(v)
	case int16:
		n = uint64(v)
	case int32:
		n = uint64(v)
	case int64:
		n = uint64(v)
	case float64: // when decoding from db: JSON numbers are floats
		n = uint64(v)
	case string:
		if strings.TrimSpace(v) == "" {
			*p = MaybeUint64Param{0, false}
			return nil
		}
		var err error
		n, err = strconv.ParseUint(v, 10, 64)
		if err != nil {
			return errors.Wrap(ErrBadInput, err.Error())
		}

	default:
		return errors.Wrapf(ErrBadInput, "expected unsigned integer or nil, got %T", val)
	}

	*p = MaybeUint64Param{n, true}
	return nil
}

func (p MaybeUint64Param) Uint64() (uint64, bool) {
	return p.n, p.isSet
}

type MaybeInt32Param struct {
	n     int32
	isSet bool
}

// NewMaybeInt32Param creates new instance of MaybeInt32Param
func NewMaybeInt32Param(n int32, isSet bool) MaybeInt32Param {
	return MaybeInt32Param{
		n:     n,
		isSet: isSet,
	}
}

func (p *MaybeInt32Param) UnmarshalPipelineParam(val interface{}) error {
	var n int32
	switch v := val.(type) {
	case uint:
		if v > math.MaxInt32 {
			return errors.Wrap(ErrBadInput, "overflows int32")
		}
		n = int32(v)
	case uint8:
		n = int32(v)
	case uint16:
		n = int32(v)
	case uint32:
		if v > math.MaxInt32 {
			return errors.Wrap(ErrBadInput, "overflows int32")
		}
		n = int32(v)
	case uint64:
		if v > math.MaxInt32 {
			return errors.Wrap(ErrBadInput, "overflows int32")
		}
		n = int32(v)
	case int:
		if v > math.MaxInt32 || v < math.MinInt32 {
			return errors.Wrap(ErrBadInput, "overflows int32")
		}
		n = int32(v)
	case int8:
		n = int32(v)
	case int16:
		n = int32(v)
	case int32:
		n = int32(v)
	case int64:
		if v > math.MaxInt32 || v < math.MinInt32 {
			return errors.Wrap(ErrBadInput, "overflows int32")
		}
		n = int32(v)
	case float64: // when decoding from db: JSON numbers are floats
		if v > math.MaxInt32 || v < math.MinInt32 {
			return errors.Wrap(ErrBadInput, "overflows int32")
		}
		n = int32(v)
	case string:
		if strings.TrimSpace(v) == "" {
			*p = MaybeInt32Param{0, false}
			return nil
		}
		i, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return errors.Wrap(ErrBadInput, err.Error())
		}
		n = int32(i)

	default:
		return errors.Wrapf(ErrBadInput, "expected signed integer or nil, got %T", val)
	}

	*p = MaybeInt32Param{n, true}
	return nil
}

func (p MaybeInt32Param) Int32() (int32, bool) {
	return p.n, p.isSet
}

type BoolParam bool

func (b *BoolParam) UnmarshalPipelineParam(val interface{}) error {
	switch v := val.(type) {
	case string:
		theBool, err := strconv.ParseBool(v)
		if err != nil {
			return errors.Wrap(ErrBadInput, err.Error())
		}
		*b = BoolParam(theBool)
		return nil
	case bool:
		*b = BoolParam(v)
		return nil
	case ObjectParam:
		if v.Type == BoolType {
			*b = v.BoolValue
			return nil
		}
	case *ObjectParam:
		if v.Type == BoolType {
			*b = v.BoolValue
			return nil
		}
	}

	return errors.Wrapf(ErrBadInput, "expected true or false, got %T", val)
}

type DecimalParam decimal.Decimal

func (d *DecimalParam) UnmarshalPipelineParam(val interface{}) error {
	if v, ok := val.(ObjectParam); ok && v.Type == DecimalType {
		*d = v.DecimalValue
		return nil
	} else if v, ok := val.(*ObjectParam); ok && v.Type == DecimalType {
		*d = v.DecimalValue
		return nil
	}
	x, err := utils.ToDecimal(val)
	if err != nil {
		return errors.Wrap(ErrBadInput, err.Error())
	}
	*d = DecimalParam(x)
	return nil
}

func (d DecimalParam) Decimal() decimal.Decimal {
	return decimal.Decimal(d)
}

type URLParam url.URL

func (u *URLParam) UnmarshalPipelineParam(val interface{}) error {
	switch v := val.(type) {
	case string:
		theURL, err := url.ParseRequestURI(v)
		if err != nil {
			return errors.Wrap(ErrBadInput, err.Error())
		}
		*u = URLParam(*theURL)
		return nil
	default:
		return ErrBadInput
	}
}

func (u *URLParam) String() string {
	return (*url.URL)(u).String()
}

type AddressParam common.Address

func (a *AddressParam) UnmarshalPipelineParam(val interface{}) error {
	switch v := val.(type) {
	case string:
		return a.UnmarshalPipelineParam([]byte(v))
	case []byte:
		if utils.HasHexPrefix(string(v)) && len(v) == 42 {
			*a = AddressParam(common.HexToAddress(string(v)))
			return nil
		} else if len(v) == 20 {
			copy((*a)[:], v)
			return nil
		}
	case common.Address:
		*a = AddressParam(v)
		return nil
	}

	return errors.Wrapf(ErrBadInput, "expected common.Address, got %T", val)
}

// MapParam accepts maps or JSON-encoded strings
type MapParam map[string]interface{}

func (m *MapParam) UnmarshalPipelineParam(val interface{}) error {
	switch v := val.(type) {
	case nil:
		*m = nil
		return nil

	case MapParam:
		*m = v
		return nil

	case map[string]interface{}:
		*m = MapParam(v)
		return nil

	case string:
		return m.UnmarshalPipelineParam([]byte(v))

	case []byte:
		var theMap map[string]interface{}
		err := json.Unmarshal(v, &theMap)
		if err != nil {
			return err
		}
		*m = MapParam(theMap)
		return nil

	case ObjectParam:
		if v.Type == MapType {
			*m = v.MapValue
			return nil
		}

	case *ObjectParam:
		if v.Type == MapType {
			*m = v.MapValue
			return nil
		}

	}

	return errors.Wrapf(ErrBadInput, "expected map, got %T", val)
}

func (m MapParam) Map() map[string]interface{} {
	return (map[string]interface{})(m)
}

type SliceParam []interface{}

func (s *SliceParam) UnmarshalPipelineParam(val interface{}) error {
	switch v := val.(type) {
	case nil:
		*s = nil
		return nil
	case []interface{}:
		*s = v
		return nil
	case string:
		return s.UnmarshalPipelineParam([]byte(v))

	case []byte:
		var theSlice []interface{}
		err := json.Unmarshal(v, &theSlice)
		if err != nil {
			return err
		}
		*s = SliceParam(theSlice)
		return nil
	}

	return errors.Wrapf(ErrBadInput, "expected slice, got %T", val)
}

func (s SliceParam) FilterErrors() (SliceParam, int) {
	var s2 SliceParam
	var errs int
	for _, x := range s {
		if _, is := x.(error); is {
			errs++
		} else {
			s2 = append(s2, x)
		}
	}
	return s2, errs
}

type DecimalSliceParam []decimal.Decimal

func (s *DecimalSliceParam) UnmarshalPipelineParam(val interface{}) error {
	var dsp DecimalSliceParam
	switch v := val.(type) {
	case nil:
		dsp = nil
	case []decimal.Decimal:
		dsp = v
	case []interface{}:
		return s.UnmarshalPipelineParam(SliceParam(v))
	case SliceParam:
		for _, x := range v {
			var d DecimalParam
			err := d.UnmarshalPipelineParam(x)
			if err != nil {
				return err
			}
			dsp = append(dsp, d.Decimal())
		}
	case string:
		return s.UnmarshalPipelineParam([]byte(v))

	case []byte:
		var theSlice []interface{}
		err := json.Unmarshal(v, &theSlice)
		if err != nil {
			return err
		}
		return s.UnmarshalPipelineParam(SliceParam(theSlice))

	default:
		return errors.Wrapf(ErrBadInput, "expected number, got %T", val)
	}
	*s = dsp
	return nil
}

type HashSliceParam []common.Hash

func (s *HashSliceParam) UnmarshalPipelineParam(val interface{}) error {
	var dsp HashSliceParam
	switch v := val.(type) {
	case nil:
		dsp = nil
	case []common.Hash:
		dsp = v
	case string:
		err := json.Unmarshal([]byte(v), &dsp)
		if err != nil {
			return errors.Wrapf(ErrBadInput, "HashSliceParam: %v", err)
		}
	case []byte:
		err := json.Unmarshal(v, &dsp)
		if err != nil {
			return errors.Wrapf(ErrBadInput, "HashSliceParam: %v", err)
		}
	case []interface{}:
		for _, h := range v {
			if s, is := h.(string); is {
				var hash common.Hash
				err := hash.UnmarshalText([]byte(s))
				if err != nil {
					return errors.Wrapf(ErrBadInput, "HashSliceParam: %v", err)
				}
				dsp = append(dsp, hash)
			} else if b, is := h.([]byte); is {
				// same semantic as AddressSliceParam
				var hash common.Hash
				err := hash.UnmarshalText(b)
				if err != nil {
					return errors.Wrapf(ErrBadInput, "HashSliceParam: %v", err)
				}
				dsp = append(dsp, hash)
			} else if h, is := h.(common.Hash); is {
				dsp = append(dsp, h)
			} else {
				return errors.Wrap(ErrBadInput, "HashSliceParam")
			}
		}
	default:
		return errors.Wrap(ErrBadInput, "HashSliceParam")
	}
	*s = dsp
	return nil
}

type AddressSliceParam []common.Address

func (s *AddressSliceParam) UnmarshalPipelineParam(val interface{}) error {
	var asp AddressSliceParam
	switch v := val.(type) {
	case nil:
		asp = nil
	case []common.Address:
		asp = v
	case string:
		err := json.Unmarshal([]byte(v), &asp)
		if err != nil {
			return errors.Wrapf(ErrBadInput, "AddressSliceParam: %v", err)
		}
	case []byte:
		err := json.Unmarshal(v, &asp)
		if err != nil {
			return errors.Wrapf(ErrBadInput, "AddressSliceParam: %v", err)
		}
	case []interface{}:
		for _, a := range v {
			var addr AddressParam
			err := addr.UnmarshalPipelineParam(a)
			if err != nil {
				return errors.Wrapf(ErrBadInput, "AddressSliceParam: %v", err)
			}
			asp = append(asp, common.Address(addr))
		}
	default:
		return errors.Wrapf(ErrBadInput, "AddressSliceParam: cannot convert %T", val)
	}
	*s = asp
	return nil
}

type JSONPathParam []string

// NewJSONPathParam returns a new JSONPathParam using the given separator, or the default if empty.
func NewJSONPathParam(sep string) JSONPathParam {
	if len(sep) == 0 {
		return nil
	}
	return []string{sep}
}

// UnmarshalPipelineParam unmarshals a slice of strings from val.
// If val is a string or []byte, it is split on a separator.
// The default separator is ',' but can be overridden by initializing via NewJSONPathParam.
func (p *JSONPathParam) UnmarshalPipelineParam(val interface{}) error {
	sep := ","
	if len(*p) > 0 {
		// custom separator
		sep = (*p)[0]
	}
	var ssp JSONPathParam
	switch v := val.(type) {
	case nil:
		ssp = nil
	case []string:
		ssp = v
	case []interface{}:
		for _, x := range v {
			if as, is := x.(string); is {
				ssp = append(ssp, as)
			} else {
				return ErrBadInput
			}
		}
	case string:
		if len(v) == 0 {
			return nil
		}
		ssp = strings.Split(v, sep)
	case []byte:
		if len(v) == 0 {
			return nil
		}
		ssp = strings.Split(string(v), sep)
	default:
		return ErrBadInput
	}
	*p = ssp
	return nil
}

type MaybeBigIntParam struct {
	n *big.Int
}

// NewMaybeBigIntParam creates a new instance of MaybeBigIntParam
func NewMaybeBigIntParam(n *big.Int) MaybeBigIntParam {
	return MaybeBigIntParam{
		n: n,
	}
}

func (p *MaybeBigIntParam) UnmarshalPipelineParam(val interface{}) error {
	var n *big.Int
	switch v := val.(type) {
	case uint:
		n = big.NewInt(0).SetUint64(uint64(v))
	case uint8:
		n = big.NewInt(0).SetUint64(uint64(v))
	case uint16:
		n = big.NewInt(0).SetUint64(uint64(v))
	case uint32:
		n = big.NewInt(0).SetUint64(uint64(v))
	case uint64:
		n = big.NewInt(0).SetUint64(v)
	case int:
		n = big.NewInt(int64(v))
	case int8:
		n = big.NewInt(int64(v))
	case int16:
		n = big.NewInt(int64(v))
	case int32:
		n = big.NewInt(int64(v))
	case int64:
		n = big.NewInt(int64(v))
	case float64: // when decoding from db: JSON numbers are floats
		n = big.NewInt(0).SetUint64(uint64(v))
	case string:
		if strings.TrimSpace(v) == "" {
			*p = MaybeBigIntParam{n: nil}
			return nil
		}
		var ok bool
		n, ok = big.NewInt(0).SetString(v, 10)
		if !ok {
			return errors.Wrapf(ErrBadInput, "unable to convert %s to big.Int", v)
		}
	case *big.Int:
		n = v
	case nil:
		*p = MaybeBigIntParam{n: nil}
		return nil
	default:
		return ErrBadInput
	}
	*p = MaybeBigIntParam{n: n}
	return nil
}

func (p MaybeBigIntParam) BigInt() *big.Int {
	return p.n
}

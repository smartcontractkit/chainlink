package pipeline

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
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
		if errors.Cause(err) == ErrParameterEmpty {
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

type GetterFunc func() (interface{}, error)

func From(getters ...interface{}) []GetterFunc {
	var gfs []GetterFunc
	for _, g := range getters {
		switch v := g.(type) {
		case GetterFunc:
			gfs = append(gfs, v)

		default:
			// If a bare value is passed in, create a simple getter
			gfs = append(gfs, func() (interface{}, error) {
				return v, nil
			})
		}
	}
	return gfs
}

func VarExpr(s string, vars Vars) GetterFunc {
	return func() (interface{}, error) {
		trimmed := strings.TrimSpace(s)
		if len(trimmed) == 0 {
			return nil, ErrParameterEmpty
		}
		isVariableExpr := strings.Count(trimmed, "$") == 1 && trimmed[:2] == "$(" && trimmed[len(trimmed)-1] == ')'
		if !isVariableExpr {
			return nil, ErrParameterEmpty
		}
		keypath := strings.TrimSpace(trimmed[2 : len(trimmed)-1])
		val, err := vars.Get(keypath)
		if err != nil {
			return nil, err
		} else if as, is := val.(error); is {
			return nil, errors.Wrapf(ErrTooManyErrors, "VarExpr: %v", as)
		}
		return val, nil
	}
}

func JSONWithVarExprs(s string, vars Vars, allowErrors bool) GetterFunc {
	return func() (interface{}, error) {
		if strings.TrimSpace(s) == "" {
			return nil, ErrParameterEmpty
		}
		replaced := variableRegexp.ReplaceAllFunc([]byte(s), func(expr []byte) []byte {
			keypathStr := strings.TrimSpace(string(expr[2 : len(expr)-1]))
			return []byte(fmt.Sprintf(`{ "__chainlink_var_expr__": "%v" }`, keypathStr))
		})
		var val interface{}
		err := json.Unmarshal(replaced, &val)
		if err != nil {
			return nil, errors.Wrapf(ErrBadInput, "while interpolating variables in JSON payload: %v", err)
		}
		return mapGoValue(val, func(val interface{}) (interface{}, error) {
			if m, is := val.(map[string]interface{}); is {
				maybeKeypath, exists := m["__chainlink_var_expr__"]
				if !exists {
					return val, nil
				}
				keypath, is := maybeKeypath.(string)
				if !is {
					return nil, errors.New("you cannot use __chainlink_var_expr__ in your JSON")
				}
				newVal, err := vars.Get(keypath)
				if err != nil {
					return nil, err
				} else if err, is := newVal.(error); is && !allowErrors {
					return nil, errors.Wrapf(ErrBadInput, "JSONWithVarExprs: %v", err)
				}
				return newVal, nil
			}
			return val, nil
		})
	}
}

func mapGoValue(v interface{}, fn func(val interface{}) (interface{}, error)) (x interface{}, err error) {
	type item struct {
		val         interface{}
		parentMap   map[string]interface{}
		parentKey   string
		parentSlice []interface{}
		parentIdx   int
	}

	stack := []item{{val: v}}
	var current item

	for len(stack) > 0 {
		current = stack[0]
		stack = stack[1:]

		val, err := fn(current.val)
		if err != nil {
			return nil, err
		}

		if current.parentMap != nil {
			current.parentMap[current.parentKey] = val
		} else if current.parentSlice != nil {
			current.parentSlice[current.parentIdx] = val
		}

		if asMap, isMap := val.(map[string]interface{}); isMap {
			for key := range asMap {
				stack = append(stack, item{val: asMap[key], parentMap: asMap, parentKey: key})
			}
		} else if asSlice, isSlice := val.([]interface{}); isSlice {
			for i := range asSlice {
				stack = append(stack, item{val: asSlice[i], parentSlice: asSlice, parentIdx: i})
			}
		}
	}
	return v, nil
}

func NonemptyString(s string) GetterFunc {
	return func() (interface{}, error) {
		trimmed := strings.TrimSpace(s)
		if len(trimmed) == 0 {
			return nil, ErrParameterEmpty
		}
		return trimmed, nil
	}
}

func Input(inputs []Result, index int) GetterFunc {
	return func() (interface{}, error) {
		if len(inputs)-1 < index {
			return nil, ErrParameterEmpty
		}
		return inputs[index].Value, inputs[index].Error
	}
}

func Inputs(inputs []Result) GetterFunc {
	return func() (interface{}, error) {
		var vals []interface{}
		for _, input := range inputs {
			if input.Error != nil {
				vals = append(vals, input.Error)
			} else {
				vals = append(vals, input.Value)
			}
		}
		return vals, nil
	}
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
	default:
		return ErrBadInput
	}
}

type BytesParam []byte

func (b *BytesParam) UnmarshalPipelineParam(val interface{}) error {
	switch v := val.(type) {
	case string:
		if len(v) >= 2 && v[:2] == "0x" {
			bs, err := hex.DecodeString(v[2:])
			if err != nil {
				return err
			}
			*b = BytesParam(bs)
			return nil
		}
		*b = BytesParam(v)
	case []byte:
		*b = BytesParam(v)
	case nil:
		*b = BytesParam(nil)
	default:
		return ErrBadInput
	}
	return nil
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
	case string:
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return errors.Wrap(ErrBadInput, err.Error())
		}
		*u = Uint64Param(n)
	default:
		return ErrBadInput
	}
	return nil
}

type MaybeUint64Param struct {
	n     uint64
	isSet bool
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
		return ErrBadInput
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
		return ErrBadInput
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
	default:
		return ErrBadInput
	}
}

type DecimalParam decimal.Decimal

func (d *DecimalParam) UnmarshalPipelineParam(val interface{}) error {
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
		if bytes.Equal(v[:2], []byte("0x")) && len(v) == 42 {
			*a = AddressParam(common.HexToAddress(string(v)))
			return nil
		} else if len(v) == 20 {
			copy((*a)[:], v)
			return nil
		}
		return ErrBadInput
	case common.Address:
		*a = AddressParam(v)
	default:
		return ErrBadInput
	}
	return nil
}

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

	default:
		return ErrBadInput
	}
}

type SliceParam []interface{}

func (s *SliceParam) UnmarshalPipelineParam(val interface{}) error {
	switch v := val.(type) {
	case nil:
		*s = nil
	case []interface{}:
		*s = v
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

	default:
		return ErrBadInput
	}
	return nil
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
			d, err := utils.ToDecimal(x)
			if err != nil {
				return errors.Wrapf(ErrBadInput, "DecimalSliceParam: wrong type of value while decoding decimals: %v", err.Error())
			}
			dsp = append(dsp, d)
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
		return errors.Wrap(ErrBadInput, "DecimalSliceParam")
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
			return err
		}
	case []byte:
		err := json.Unmarshal(v, &dsp)
		if err != nil {
			return err
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

func (p *JSONPathParam) UnmarshalPipelineParam(val interface{}) error {
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
		ssp = strings.Split(v, ",")
	case []byte:
		if len(v) == 0 {
			return nil
		}
		ssp = strings.Split(string(v), ",")
	default:
		return ErrBadInput
	}
	*p = ssp
	return nil
}

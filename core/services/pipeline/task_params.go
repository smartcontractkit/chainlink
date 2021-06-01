package pipeline

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

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

type BytesParam string

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
		return nil
	default:
		return ErrBadInput
	}
	return nil
}

type MaybeUint64Param struct {
	n     uint64
	isSet bool
}

func (u *MaybeUint64Param) UnmarshalPipelineParam(val interface{}) error {
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
			*u = MaybeUint64Param{0, false}
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

	*u = MaybeUint64Param{n, true}
	return nil
}

func (u MaybeUint64Param) Uint64() (uint64, bool) {
	return u.n, u.isSet
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

type MapParam map[string]interface{}

func (m *MapParam) UnmarshalPipelineParam(val interface{}) error {
	switch v := val.(type) {
	case nil:
		*m = nil
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

type StringSliceParam []string

func (p *StringSliceParam) UnmarshalPipelineParam(val interface{}) error {
	var ssp StringSliceParam
	switch v := val.(type) {
	case nil:
		ssp = nil
	case []string:
		ssp = v
	case []interface{}:
		for _, x := range v {
			if as, is := x.(string); is {
				ssp = append(ssp, strings.TrimSpace(as))
			} else {
				return ErrBadInput
			}
		}
	case string:
		ssp = strings.Split(v, ",")
	case []byte:
		ssp = strings.Split(string(v), ",")
	default:
		return ErrBadInput
	}
	*p = ssp
	return nil
}

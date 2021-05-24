package pipeline

import (
	"encoding/hex"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/multierr"
)

type StringParam string

func (s *StringParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
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

func (b *BytesParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
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

func (u *Uint64Param) UnmarshalPipelineParam(val interface{}, vars Vars) error {
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

func (u *MaybeUint64Param) UnmarshalPipelineParam(val interface{}, vars Vars) error {
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

func (b *BoolParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
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

type MaybeBoolParam string

const (
	MaybeBoolTrue  = MaybeBoolParam("true")
	MaybeBoolFalse = MaybeBoolParam("false")
	MaybeBoolNull  = MaybeBoolParam("")
)

func (m MaybeBoolParam) Bool() (b bool, isSet bool) {
	switch m {
	case MaybeBoolTrue:
		return true, true
	case MaybeBoolFalse:
		return false, true
	default:
		return false, false
	}
}

func (m *MaybeBoolParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
	switch val {
	case "true":
		*m = MaybeBoolTrue
	case "false":
		*m = MaybeBoolFalse
	case "":
		*m = MaybeBoolNull
	case true:
		*m = MaybeBoolTrue
	case false:
		*m = MaybeBoolFalse
	default:
		return ErrBadInput
	}
	return nil
}

type DecimalParam decimal.Decimal

func (d *DecimalParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
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

func (u *URLParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
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

func (m *MapParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
	switch v := val.(type) {
	case nil:
		*m = nil
		return nil

	case map[string]interface{}:
		*m = MapParam(v)
		return nil

	case string:
		return m.UnmarshalPipelineParam([]byte(v), vars)

	case []byte:
		resolved, err := expandVariables(v, vars)
		if err != nil {
			return errors.Wrapf(ErrBadInput, "MapParam: %v", err)
		}

		var theMap map[string]interface{}
		err = json.Unmarshal(resolved, &theMap)
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

func (s *SliceParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
	switch v := val.(type) {
	case nil:
		*s = nil
	case []interface{}:
		*s = v
	case string:
		return s.UnmarshalPipelineParam([]byte(v), vars)

	case []byte:
		resolved, err := expandVariables(v, vars)
		if err != nil {
			return errors.Wrapf(ErrBadInput, "MapParam: %v", err)
		}
		var theSlice []interface{}
		err = json.Unmarshal(resolved, &theSlice)
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

func (s *DecimalSliceParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
	var dsp DecimalSliceParam
	switch v := val.(type) {
	case nil:
		dsp = nil
	case []decimal.Decimal:
		dsp = v
	case []interface{}:
		return s.UnmarshalPipelineParam(SliceParam(v), vars)
	case SliceParam:
		for _, x := range v {
			d, err := utils.ToDecimal(x)
			if err != nil {
				return errors.Wrapf(ErrBadInput, "DecimalSliceParam: wrong type of value while decoding decimals: %v", err.Error())
			}
			dsp = append(dsp, d)
		}
	case string:
		return s.UnmarshalPipelineParam([]byte(v), vars)

	case []byte:
		resolved, err := expandVariables(v, vars)
		if err != nil {
			return errors.Wrapf(ErrBadInput, "DecimalSliceParam: %v", err)
		}
		var theSlice []interface{}
		err = json.Unmarshal(resolved, &theSlice)
		if err != nil {
			return err
		}
		return s.UnmarshalPipelineParam(SliceParam(theSlice), vars)

	default:
		return errors.Wrap(ErrBadInput, "DecimalSliceParam")
	}
	*s = dsp
	return nil
}

type StringSliceParam []string

func (p *StringSliceParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
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

func expandVariables(v []byte, vars Vars) ([]byte, error) {
	var err error
	resolved := variableRegexp.ReplaceAllFunc(v, func(keypath []byte) []byte {
		val, err2 := vars.Get(string(keypath[2 : len(keypath)-1]))
		if err2 != nil {
			err = multierr.Append(err, err2)
			return nil
		}
		bs, err2 := json.Marshal(val)
		if err2 != nil {
			err = multierr.Append(err, err2)
			return nil
		}
		return bs
	})
	if err != nil {
		return nil, err
	}
	return resolved, nil
}

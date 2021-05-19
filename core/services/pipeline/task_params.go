package pipeline

import (
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
	case map[string]interface{}:
		*m = MapParam(v)
		return nil

	case string:
		var err error
		resolved := variableRegexp.ReplaceAllFunc([]byte(v), func(keypath []byte) []byte {
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
			return err
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
	case []interface{}:
		*s = v
	case string:
		return json.Unmarshal([]byte(v), s)
	case []byte:
		return json.Unmarshal(v, s)
	default:
		return ErrBadInput
	}
	return nil
}

type DecimalSliceParam []decimal.Decimal

func (s *DecimalSliceParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
	switch v := val.(type) {
	case []decimal.Decimal:
		*s = v
	case []interface{}:
		var dsp DecimalSliceParam
		for _, x := range v {
			d, err := utils.ToDecimal(x)
			if err != nil {
				return errors.Wrap(ErrBadInput, err.Error())
			}
			dsp = append(dsp, d)
		}
		*s = dsp
	case string:
		return json.Unmarshal([]byte(v), s)
	case []byte:
		return json.Unmarshal(v, s)
	default:
		return ErrBadInput
	}
	return nil
}

type StringSliceParam []string

func (p *StringSliceParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
	switch v := val.(type) {
	case []string:
		*p = v
	case []interface{}:
		var ssp StringSliceParam
		for _, x := range v {
			if as, is := x.(string); is {
				ssp = append(ssp, as)
			} else {
				return ErrBadInput
			}
		}
		*p = ssp
	case string:
		*p = strings.Split(v, ",")
	default:
		return ErrBadInput
	}
	return nil
}

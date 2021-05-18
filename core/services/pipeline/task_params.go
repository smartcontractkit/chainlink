package pipeline

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"go.uber.org/multierr"
)

type StringParam string

func (s *StringParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
	switch v := val.(type) {
	case string:
		*s = StringParam(v)
		return nil
	default:
		return ErrBadInput
	}
}

type BoolParam bool

func (b *BoolParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
	switch v := val.(type) {
	case string:
		theBool, err := strconv.ParseBool(v)
		if err != nil {
			return err
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
		return nil
	case "false":
		*m = MaybeBoolFalse
		return nil
	case "":
		*m = MaybeBoolNull
		return nil
	default:
		return ErrBadInput
	}
}

type URLParam url.URL

func (u *URLParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
	switch v := val.(type) {
	case string:
		theURL, err := url.ParseRequestURI(v)
		if err != nil {
			return err
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

type JSONPathParam []string

func (p *JSONPathParam) UnmarshalPipelineParam(val interface{}, vars Vars) error {
	// var s string
	// switch v := val.(type) {
	// case string:
	//  s = v
	// default:
	//  return nil, ErrBadInput
	// }
	return nil

	// trimmed := strings.TrimSpace(s)
	// if len(trimmed) == 0 {
	//  return nil, ErrBadInput
	// }
	// if trimmed[0] == "[" {
	//  if trimmed[len(trimmed)-1] != "]" {
	//      return nil, ErrBadInput
	//  }
	//  elems := strings.Split(trimmed[1:len(trimmed)-1], ",")
	//  elems = trimStrings(elems)
	//  for _, elem := range elems {
	//      t.Resolve(elem, vars, nil)
	//  }
	// }
}

func (p *JSONPathParam) UnmarshalText(bs []byte) error {
	*p = strings.Split(string(bs), ",")
	return nil
}

func (p *JSONPathParam) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), p)
}
func (p JSONPathParam) Value() (driver.Value, error) {
	return json.Marshal(p)
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
		fmt.Println(string(resolved))

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
	return nil
}

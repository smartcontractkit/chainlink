package adapters

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	simplejson "github.com/bitly/go-simplejson"
	gjson "github.com/tidwall/gjson"
)

// JSONParse holds a path to the desired field in a JSON object,
// made up of an array of strings.
type JSONParse struct {
	Path JSONPath `json:"path"`
}

// TaskType returns the type of Adapter.
func (jpa *JSONParse) TaskType() models.TaskType {
	return TaskTypeJSONParse
}

// Perform returns the value associated to the desired field for a
// given JSON object.
//
// For example, if the JSON data looks like this:
//   {
//     "data": [
//       {"last": "1111"},
//       {"last": "2222"}
//     ]
//   }
//
// Then ["0","last"] would be the path, and "1111" would be the returned value
func (jpa *JSONParse) Perform(input models.RunInput, _ *store.Store, _ *keystore.Master) models.RunOutput {
	var val string
	var err error

	if input.Result().Type == gjson.JSON {
		// Handle case where JSON comes "pre-packaged" as gjson e.g. from bridge (external adapters)
		val = input.Result().Raw
	} else {
		val, err = input.ResultString()
	}

	if err != nil {
		return models.NewRunOutputError(err)
	}

	js, err := simplejson.NewJson([]byte(val))
	if err != nil {
		return models.NewRunOutputError(err)
	}

	last, err := dig(js, jpa.Path)
	if err != nil {
		return moldErrorOutput(js, jpa.Path, input)
	}

	return models.NewRunOutputCompleteWithResult(last.Interface(), input.ResultCollection())
}

func dig(js *simplejson.Json, path []string) (*simplejson.Json, error) {
	var ok bool
	for _, k := range path[:] {
		if isArray(js, k) {
			js, ok = arrayGet(js, k)
		} else {
			js, ok = js.CheckGet(k)
		}
		if !ok {
			return js, errors.New("No value could be found for the key '" + k + "'")
		}
	}
	return js, nil
}

// only error if any keys prior to the last one in the path are nonexistent.
// i.e. Path = ["errorIfNonExistent", "nullIfNonExistent"]
func moldErrorOutput(js *simplejson.Json, path []string, input models.RunInput) models.RunOutput {
	if _, err := getEarlyPath(js, path); err != nil {
		return models.NewRunOutputError(err)
	}
	return models.NewRunOutputCompleteWithResult(nil, input.ResultCollection())
}

func getEarlyPath(js *simplejson.Json, path []string) (*simplejson.Json, error) {
	var ok bool
	for _, k := range path[:len(path)-1] {
		if isArray(js, k) {
			js, ok = arrayGet(js, k)
		} else {
			js, ok = js.CheckGet(k)
		}
		if !ok {
			return js, errors.New("No value could be found for the key '" + k + "'")
		}
	}
	return js, nil
}

func arrayGet(js *simplejson.Json, key string) (*simplejson.Json, bool) {
	input, err := strconv.ParseInt(key, 10, 32)
	if err != nil {
		return js, false
	}
	a, err := js.Array()
	if err != nil {
		return js, false
	}

	index := int(input)
	if index < 0 {
		index = len(a) + index
	}

	if index >= len(a) || index < 0 {
		return js, false
	}
	return js.GetIndex(index), true
}

func isArray(js *simplejson.Json, key string) bool {
	if _, err := js.Array(); err != nil {
		return false
	}
	return true
}

// JSONPath is a path to a value in a JSON object
type JSONPath []string

// UnmarshalJSON implements the Unmarshaler interface
func (jp *JSONPath) UnmarshalJSON(b []byte) error {
	strs := []string{}
	var err error
	if utils.IsQuoted(b) {
		strs = strings.Split(string(utils.RemoveQuotes(b)), ".")
	} else {
		err = json.Unmarshal(b, &strs)
	}
	*jp = JSONPath(strs)
	return err
}

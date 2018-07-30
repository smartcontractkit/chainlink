package adapters

import (
	"encoding/json"
	"errors"
	"strconv"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// JSONParse holds a path to the desired field in a JSON object,
// made up of an array of strings.
type JSONParse struct {
	Path []string `json:"path"`
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
// Then ["0","last"] would be the path, and "111" would be the returned value
func (jpa *JSONParse) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	val, err := input.Value()
	if err != nil {
		return input.WithError(err)
	}

	js, err := simplejson.NewJson([]byte(val))
	if err != nil {
		return input.WithError(err)
	}

	last, err := dig(js, jpa.Path)
	if err != nil {
		return moldErrorOutput(js, jpa.Path, input)
	}

	rval, err := getStringValue(last)
	if err != nil {
		return input.WithError(err)
	}
	return input.WithValue(rval)
}

func dig(js *simplejson.Json, path []string) (*simplejson.Json, error) {
	var ok bool
	for _, k := range path[:len(path)] {
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
func moldErrorOutput(js *simplejson.Json, path []string, input models.RunResult) models.RunResult {
	if _, err := getEarlyPath(js, path); err != nil {
		return input.WithError(err)
	}
	return input.WithNull()
}

func getStringValue(js *simplejson.Json) (string, error) {
	str, err := js.String()
	if err != nil {
		b, err := json.Marshal(js)
		if err != nil {
			return str, err
		}
		str = string(b)
	}
	return str, nil
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
	i, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		return js, false
	}
	a, err := js.Array()
	if err != nil || len(a) < int(i-1) {
		return js, false
	}
	return js.GetIndex(int(i)), true
}

func isArray(js *simplejson.Json, key string) bool {
	if _, err := js.Array(); err != nil {
		return false
	}
	return true
}

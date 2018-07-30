package adapters

import (
	"errors"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/tidwall/gjson"
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

	if !gjson.Valid(val) {
		return input.WithError(errors.New("Invalid json"))
	}

	parsed := gjson.Parse(val)
	js := parsed.Get(jpa.pathSyntax())
	if !js.Exists() {
		return moldErrorOrNullResult(parsed, jpa, input)
	}
	return input.WithValue(js.String())
}

// moldErrorOrNullResult Only error if any keys prior to the last one in the
// path are nonexistent.
// i.e. Path = ["errorsIfNonExistent", "nullIfNonExistent"]
func moldErrorOrNullResult(parsed gjson.Result, jpa *JSONParse, input models.RunResult) models.RunResult {
	prefixPath := strings.Join(jpa.Path[:len(jpa.Path)-1], ".")
	if len(prefixPath) == 0 {
		return input.WithNull()
	}
	prefix := parsed.Get(prefixPath)
	if prefix.Exists() {
		return input.WithNull()
	}
	return input.WithError(fmt.Errorf("No value could be found for the path %s", jpa.Path))
}

func (jpa *JSONParse) pathSyntax() string {
	return strings.Join(jpa.Path, ".")
}

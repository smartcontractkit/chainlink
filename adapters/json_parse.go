package adapters

import (
	"errors"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// JsonParse holds the path to the desired field in a JSON object
type JsonParse struct {
	Path []string `json:"path"`
}

// Perform returns the value associated to the desired field for a
// given JSON object.
//
// For example, if the JSON data looks like this:
//   {
//     "last": "1400"
//   }
//
// Then "last" would be the path, and "1400" would be the returned value
func (jpa *JsonParse) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	// Attempt to store the JSON data given as input
	js, err := simplejson.NewJson([]byte(input.Value()))
	// Return the error if present
	if err != nil {
		return models.RunResultWithError(err)
	}
	// Check if the desired field is available as a decendent
	js, err = checkEarlyPath(js, jpa.Path)
	// Return the error if the path isn't present
	if err != nil {
		return models.RunResultWithError(err)
	}
	// Get the value within the JSON object for the desired path field
	rval, ok := js.CheckGet(jpa.Path[len(jpa.Path)-1])
	// If CheckGet couldn't find the value or path, return an error
	if !ok {
		return models.RunResult{}
	}
	// Return the value of the desired path field
	return models.RunResultWithValue(rval.MustString())
}

// Ensures that the given path for the task is present within
// the JSON object
func checkEarlyPath(js *simplejson.Json, path []string) (*simplejson.Json, error) {
	var ok bool
	for _, k := range path[:len(path)-1] {
		js, ok = js.CheckGet(k)
		if !ok {
			return js, errors.New("No value could be found for the key '" + k + "'")
		}
	}
	return js, nil
}

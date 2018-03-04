package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/araddon/dateparse"
	"github.com/mrwonko/cron"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/tidwall/gjson"
)

// JSON stores the json types string, number, bool, and null.
// Arrays and Objects are returned as their raw json types.
type JSON struct {
	gjson.Result
}

// ParseJSON attempts to coerce the input byte array into valid JSON
// and parse it into a JSON object.
func ParseJSON(b []byte) (JSON, error) {
	var j JSON
	str := string(b)
	if len(str) == 0 {
		str = `{}`
	}
	return j, json.Unmarshal([]byte(str), &j)
}

// UnmarshalJSON parses the JSON bytes and stores in the *JSON pointer.
func (j *JSON) UnmarshalJSON(b []byte) error {
	str := string(b)
	if !gjson.Valid(str) {
		return fmt.Errorf("invalid JSON: %v", str)
	}
	*j = JSON{gjson.Parse(str)}
	return nil
}

// MarshalJSON returns the JSON data if it already exists, returns
// an empty JSON object as bytes if not.
func (j JSON) MarshalJSON() ([]byte, error) {
	if j.Exists() {
		return j.Bytes(), nil
	}
	return []byte("{}"), nil
}

// Merge combines the given JSON with the existing JSON.
func (j JSON) Merge(j2 JSON) (JSON, error) {
	body := j.Map()
	for key, value := range j2.Map() {
		body[key] = value
	}

	cleaned := map[string]interface{}{}
	for k, v := range body {
		cleaned[k] = v.Value()
	}

	b, err := json.Marshal(cleaned)
	if err != nil {
		return JSON{}, err
	}

	var rval JSON
	return rval, gjson.Unmarshal(b, &rval)
}

// Empty returns true if the JSON does not exist.
func (j JSON) Empty() bool {
	return !j.Exists()
}

// Bytes returns the raw JSON.
func (j JSON) Bytes() []byte {
	return []byte(j.String())
}

// Add returns a new instance of JSON with the new value added.
func (j JSON) Add(key string, val interface{}) (JSON, error) {
	var j2 JSON
	b, err := json.Marshal(val)
	if err != nil {
		return j2, err
	}
	str := fmt.Sprintf(`{"%v":%v}`, key, string(b))
	if err = json.Unmarshal([]byte(str), &j2); err != nil {
		return j2, err
	}
	return j.Merge(j2)
}

// WebURL contains the URL of the endpoint.
type WebURL struct {
	*url.URL
}

// UnmarshalJSON parses the raw URL stored in JSON-encoded
// data to a URL structure and sets it to the URL field.
func (w *WebURL) UnmarshalJSON(j []byte) error {
	var v string
	err := json.Unmarshal(j, &v)
	if err != nil {
		return err
	}
	u, err := url.ParseRequestURI(v)
	if err != nil {
		return err
	}
	w.URL = u
	return nil
}

// MarshalJSON returns the JSON-encoded string of the given data.
func (w *WebURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.String())
}

// Time holds a common field for time.
type Time struct {
	time.Time
}

// UnmarshalJSON parses the raw time stored in JSON-encoded
// data and stores it to the Time field.
func (t *Time) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	newTime, err := dateparse.ParseAny(s)
	t.Time = newTime
	return err
}

// ISO8601 formats and returns the time in ISO 8601 standard.
func (t *Time) ISO8601() string {
	return t.UTC().Format("2006-01-02T15:04:05Z07:00")
}

// DurationFromNow returns the amount of time since the Time
// field was last updated.
func (t *Time) DurationFromNow() time.Duration {
	return t.Time.Sub(time.Now())
}

// HumanString formats and returns the time in RFC 3339 standard.
func (t *Time) HumanString() string {
	return utils.ISO8601UTC(t.Time)
}

// Cron holds the string that will represent the spec of the cron-job.
// It uses 6 fields to represent the seconds (1), minutes (2), hours (3),
// day of the month (4), month (5), and day of the week (6).
type Cron string

// UnmarshalJSON parses the raw spec stored in JSON-encoded
// data and stores it to the Cron string.
func (c *Cron) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("Cron: %v", err)
	}
	if s == "" {
		return nil
	}

	_, err = cron.Parse(s)
	if err != nil {
		return fmt.Errorf("Cron: %v", err)
	}
	*c = Cron(s)
	return nil
}

// String returns the current Cron spec string.
func (c Cron) String() string {
	return string(c)
}

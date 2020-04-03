package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"chainlink/core/assets"

	"github.com/araddon/dateparse"
	"github.com/jinzhu/gorm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
	"github.com/mrwonko/cron"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// RunStatus is a string that represents the run status
type RunStatus string

const (
	// RunStatusUnstarted is the default state of any run status.
	RunStatusUnstarted = RunStatus("")
	// RunStatusInProgress is used for when a run is actively being executed.
	RunStatusInProgress = RunStatus("in_progress")
	// RunStatusPendingConfirmations is used for when a run is awaiting for block confirmations.
	RunStatusPendingConfirmations = RunStatus("pending_confirmations")
	// RunStatusPendingConnection states that the run is waiting on a connection to the block chain.
	RunStatusPendingConnection = RunStatus("pending_connection")
	// RunStatusPendingBridge is used for when a run is waiting on the completion
	// of another event.
	RunStatusPendingBridge = RunStatus("pending_bridge")
	// RunStatusPendingSleep is used for when a run is waiting on a sleep function to finish.
	RunStatusPendingSleep = RunStatus("pending_sleep")
	// RunStatusErrored is used for when a run has errored and will not complete.
	RunStatusErrored = RunStatus("errored")
	// RunStatusCompleted is used for when a run has successfully completed execution.
	RunStatusCompleted = RunStatus("completed")
	// RunStatusCancelled is used to indicate a run is no longer desired.
	RunStatusCancelled = RunStatus("cancelled")
)

// Unstarted returns true if the status is the initial state.
func (s RunStatus) Unstarted() bool {
	return s == RunStatusUnstarted
}

// PendingBridge returns true if the status is pending_bridge.
func (s RunStatus) PendingBridge() bool {
	return s == RunStatusPendingBridge
}

// PendingConfirmations returns true if the status is pending_confirmations.
func (s RunStatus) PendingConfirmations() bool {
	return s == RunStatusPendingConfirmations
}

// PendingConnection returns true if the status is pending_connection.
func (s RunStatus) PendingConnection() bool {
	return s == RunStatusPendingConnection
}

// PendingSleep returns true if the status is pending_sleep.
func (s RunStatus) PendingSleep() bool {
	return s == RunStatusPendingSleep
}

// Completed returns true if the status is RunStatusCompleted.
func (s RunStatus) Completed() bool {
	return s == RunStatusCompleted
}

// Cancelled returns true if the status is RunStatusCancelled.
func (s RunStatus) Cancelled() bool {
	return s == RunStatusCancelled
}

// Errored returns true if the status is RunStatusErrored.
func (s RunStatus) Errored() bool {
	return s == RunStatusErrored
}

// Pending returns true if the status is pending external or confirmations.
func (s RunStatus) Pending() bool {
	return s.PendingBridge() || s.PendingConfirmations() || s.PendingSleep() || s.PendingConnection()
}

// Finished returns true if the status is final and can't be changed.
func (s RunStatus) Finished() bool {
	return s.Completed() || s.Errored() || s.Cancelled()
}

// Runnable returns true if the status is ready to be run.
func (s RunStatus) Runnable() bool {
	return !s.Errored() && !s.Pending()
}

// CanStart returns true if the run is ready to begin processed.
func (s RunStatus) CanStart() bool {
	return s.Pending() || s.Unstarted()
}

// Value returns this instance serialized for database storage.
func (s RunStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan reads the database value and returns an instance.
func (s *RunStatus) Scan(value interface{}) error {
	temp, ok := value.(string)
	if !ok {
		return fmt.Errorf("Unable to convert %v of %T to RunStatus", value, value)
	}

	*s = RunStatus(temp)
	return nil
}

// JSON stores the json types string, number, bool, and null.
// Arrays and Objects are returned as their raw json types.
type JSON struct {
	gjson.Result
}

// Value returns this instance serialized for database storage.
func (j JSON) Value() (driver.Value, error) {
	s := j.Bytes()
	if len(s) == 0 {
		return nil, nil
	}
	return s, nil
}

// Scan reads the database value and returns an instance.
func (j *JSON) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		*j = JSON{Result: gjson.Parse(v)}
	case []byte:
		*j = JSON{Result: gjson.ParseBytes(v)}
	default:
		return fmt.Errorf("Unable to convert %v of %T to JSON", value, value)
	}
	return nil
}

// ParseJSON attempts to coerce the input byte array into valid JSON
// and parse it into a JSON object.
func ParseJSON(b []byte) (JSON, error) {
	var j JSON
	str := string(b)
	if len(str) == 0 {
		return j, nil
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

// Bytes returns the raw JSON.
func (j JSON) Bytes() []byte {
	return []byte(j.String())
}

// AsMap returns j as a map
func (j JSON) AsMap() (map[string]interface{}, error) {
	output := make(map[string]interface{})
	switch v := j.Result.Value().(type) {
	case map[string]interface{}:
		for key, value := range v {
			output[key] = value
		}
	case nil:
	default:
		return nil, errors.New("can only add to JSON objects or null")
	}
	return output, nil
}

// mapToJSON returns m as a JSON object, or errors
func mapToJSON(m map[string]interface{}) (JSON, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return JSON{}, err
	}
	return JSON{Result: gjson.ParseBytes(bytes)}, nil
}

// Add returns a new instance of JSON with the new value added.
func (j JSON) Add(insertKey string, insertValue interface{}) (JSON, error) {
	return j.MultiAdd(KV{insertKey: insertValue})
}

// KV represents a key/value pair to be added to a JSON object
type KV map[string]interface{}

// MultiAdd returns a new instance of j with the new values added.
func (j JSON) MultiAdd(keyValues KV) (JSON, error) {
	output, err := j.AsMap()
	if err != nil {
		return JSON{}, err
	}
	for key, value := range keyValues {
		output[key] = value
	}
	return mapToJSON(output)
}

// Delete returns a new instance of JSON with the specified key removed.
func (j JSON) Delete(key string) (JSON, error) {
	js, err := sjson.Delete(j.String(), key)
	if err != nil {
		return j, err
	}
	return ParseJSON([]byte(js))
}

// CBOR returns a bytes array of the JSON map or array encoded to CBOR.
func (j JSON) CBOR() ([]byte, error) {
	switch v := j.Result.Value().(type) {
	case map[string]interface{}, []interface{}, nil:
		return cbor.Marshal(v)
	default:
		var b []byte
		return b, fmt.Errorf("Unable to coerce JSON to CBOR for type %T", v)
	}
}

// WebURL contains the URL of the endpoint.
type WebURL url.URL

// UnmarshalJSON parses the raw URL stored in JSON-encoded
// data to a URL structure and sets it to the URL field.
func (w *WebURL) UnmarshalJSON(j []byte) error {
	var v string
	err := json.Unmarshal(j, &v)
	if err != nil {
		return err
	}
	// handle no url case
	if len(v) == 0 {
		return nil
	}

	u, err := url.ParseRequestURI(v)
	if err != nil {
		return err
	}
	*w = WebURL(*u)
	return nil
}

// MarshalJSON returns the JSON-encoded string of the given data.
func (w WebURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.String())
}

// String delegates to the wrapped URL struct or an empty string when it is nil
func (w WebURL) String() string {
	url := url.URL(w)
	return url.String()
}

// Value returns this instance serialized for database storage.
func (w WebURL) Value() (driver.Value, error) {
	return w.String(), nil
}

// Scan reads the database value and returns an instance.
func (w *WebURL) Scan(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("Unable to convert %v of %T to WebURL", value, value)
	}

	u, err := url.ParseRequestURI(s)
	if err != nil {
		return err
	}
	*w = WebURL(*u)
	return nil
}

// AnyTime holds a common field for time, and serializes it as
// a json number.
type AnyTime struct {
	time.Time
	Valid bool
}

// NewAnyTime creates a new Time.
func NewAnyTime(t time.Time) AnyTime {
	return AnyTime{Time: t, Valid: true}
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (t AnyTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.UTC().MarshalJSON()
}

// UnmarshalJSON parses the raw time stored in JSON-encoded
// data and stores it to the Time field.
func (t *AnyTime) UnmarshalJSON(b []byte) error {
	var str string

	var n json.Number
	if err := json.Unmarshal(b, &n); err == nil {
		str = n.String()
	} else if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	if len(str) == 0 {
		t.Valid = false
		return nil
	}

	newTime, err := dateparse.ParseAny(str)
	t.Time = newTime.UTC()
	t.Valid = true
	return err
}

// MarshalText returns null if not set, or the time.
func (t AnyTime) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.MarshalText()
}

// UnmarshalText parses null or a valid time.
func (t *AnyTime) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	if err := t.Time.UnmarshalText(text); err != nil {
		return err
	}
	t.Valid = true
	return nil
}

// Value returns this instance serialized for database storage.
func (t AnyTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// Scan reads the database value and returns an instance.
func (t *AnyTime) Scan(value interface{}) error {
	switch temp := value.(type) {
	case time.Time:
		t.Time = temp.UTC()
		t.Valid = true
		return nil
	case nil:
		t.Valid = false
		return nil
	default:
		return fmt.Errorf("Unable to convert %v of %T to Time", value, value)
	}
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

// Duration is a time duration.
type Duration time.Duration

// Duration returns the value as the standard time.Duration value.
func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

// String returns a string representing the duration in the form "72h3m0.5s".
// Leading zero units are omitted. As a special case, durations less than one
// second format use a smaller unit (milli-, micro-, or nanoseconds) to ensure
// that the leading digit is non-zero. The zero duration formats as 0s.
func (d Duration) String() string {
	return time.Duration(d).String()
}

// MarshalJSON implements the json.Marshaler interface.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Duration) UnmarshalJSON(input []byte) error {
	var txt string
	err := json.Unmarshal(input, &txt)
	if err != nil {
		return err
	}
	v, err := time.ParseDuration(string(txt))
	if err != nil {
		return err
	}
	*d = Duration(v)
	return nil
}

// WithdrawalRequest request to withdraw LINK.
type WithdrawalRequest struct {
	DestinationAddress common.Address `json:"address"`
	ContractAddress    common.Address `json:"contractAddress"`
	Amount             *assets.Link   `json:"amount"`
}

// SendEtherRequest represents a request to transfer ETH.
type SendEtherRequest struct {
	DestinationAddress common.Address `json:"address"`
	FromAddress        common.Address `json:"from"`
	Amount             *assets.Eth    `json:"amount"`
}

// CreateKeyRequest represents a request to add an ethereum key.
type CreateKeyRequest struct {
	CurrentPassword string `json:"current_password"`
}

// AddressCollection is an array of common.Address
// serializable to and from a database.
type AddressCollection []common.Address

// ToStrings returns this address collection as an array of strings.
func (r AddressCollection) ToStrings() []string {
	// Unable to convert copy-free without unsafe:
	// https://stackoverflow.com/a/48554123/639773
	converted := make([]string, len(r))
	for i, e := range r {
		converted[i] = e.Hex()
	}
	return converted
}

// Value returns the string value to be written to the database.
func (r AddressCollection) Value() (driver.Value, error) {
	return strings.Join(r.ToStrings(), ","), nil
}

// Scan parses the database value as a string.
func (r *AddressCollection) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("Unable to convert %v of %T to AddressCollection", value, value)
	}

	if len(str) == 0 {
		return nil
	}

	arr := strings.Split(str, ",")
	collection := make(AddressCollection, len(arr))
	for i, a := range arr {
		collection[i] = common.HexToAddress(a)
	}
	*r = collection
	return nil
}

// Configuration stores key value pairs for overriding global configuration
type Configuration struct {
	gorm.Model
	Name  string `gorm:"not null;unique;index"`
	Value string `gorm:"not null"`
}

// Merge returns a new map with all keys merged from right to left
func Merge(inputs ...JSON) (JSON, error) {
	output := make(map[string]interface{})

	for _, input := range inputs {
		switch v := input.Result.Value().(type) {
		case map[string]interface{}:
			for key, value := range v {
				output[key] = value
			}
		case nil:
		default:
			return JSON{}, errors.New("can only merge JSON objects")
		}
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		return JSON{}, err
	}

	return JSON{Result: gjson.ParseBytes(bytes)}, nil
}

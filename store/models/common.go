package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/store/assets"

	"github.com/araddon/dateparse"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mrwonko/cron"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/ugorji/go/codec"
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
	return s.Completed() || s.Errored()
}

// Runnable returns true if the status is ready to be run.
func (s RunStatus) Runnable() bool {
	return !s.Errored() && !s.Pending()
}

// CanStart returns true if the run is ready to begin processed.
func (s RunStatus) CanStart() bool {
	return s.Pending() || s.Unstarted()
}

// Automatically add missing start map and end map to a CBOR encoded buffer
func autoAddMapDelimiters(b []byte) []byte {
	if len(b) < 2 {
		return b
	}

	var buffer bytes.Buffer
	if b[0] != 0xbf {
		buffer.Write([]byte{0xbf})
	}
	buffer.Write(b)

	if b[len(b)-1] != 0xff {
		buffer.Write([]byte{0xff})
	}
	return buffer.Bytes()
}

// ParseCBOR attempts to coerce the input byte array into valid CBOR
// and then coerces it into a JSON object.
func ParseCBOR(b []byte) (JSON, error) {
	var m map[interface{}]interface{}

	cbor := codec.NewDecoderBytes(autoAddMapDelimiters(b), new(codec.CborHandle))
	if err := cbor.Decode(&m); err != nil {
		return JSON{}, err
	}

	coerced, err := utils.CoerceInterfaceMapToStringMap(m)
	if err != nil {
		return JSON{}, err
	}

	jsb, err := json.Marshal(coerced)
	if err != nil {
		return JSON{}, err
	}

	var js JSON
	return js, json.Unmarshal(jsb, &js)
}

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
	var b []byte
	cbor := codec.NewEncoderBytes(&b, new(codec.CborHandle))

	switch v := j.Value().(type) {
	case map[string]interface{}, []interface{}, nil:
		return b, cbor.Encode(v)
	default:
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

// Time holds a common field for time.
type Time struct {
	time.Time
}

// UnmarshalJSON parses the raw time stored in JSON-encoded
// data and stores it to the Time field.
func (t *Time) UnmarshalJSON(b []byte) error {
	var n json.Number
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	newTime, err := dateparse.ParseAny(n.String())
	t.Time = newTime.UTC()
	return err
}

// ISO8601 formats and returns the time in ISO 8601 standard.
func (t Time) ISO8601() string {
	return t.UTC().Format("2006-01-02T15:04:05Z07:00")
}

// DurationFromNow returns the amount of time since the Time
// field was last updated.
func (t Time) DurationFromNow() time.Duration {
	return t.Time.Sub(time.Now())
}

// HumanString formats and returns the time in RFC 3339 standard.
func (t Time) HumanString() string {
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

// WithdrawalRequest request to withdraw LINK.
type WithdrawalRequest struct {
	DestinationAddress common.Address `json:"address"`
	ContractAddress    common.Address `json:"contractAddress"`
	Amount             *assets.Link   `json:"amount"`
}

// SendEtherRequest represents a request to transfer ETH.
type SendEtherRequest struct {
	DestinationAddress common.Address `json:"address"`
	Amount             *assets.Eth    `json:"amount"`
}

// Int stores large integers and can deserialize a variety of inputs.
type Int big.Int

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *Int) UnmarshalText(input []byte) error {
	input = utils.RemoveQuotes(input)
	str := string(input)
	var ok bool
	if utils.HasHexPrefix(str) {
		i, ok = i.setString(utils.RemoveHexPrefix(str), 16)
	} else {
		i, ok = i.setString(str, 10)
	}

	if !ok {
		return fmt.Errorf("could not unmarshal %s to Int", str)
	}
	return nil
}

// UnmarshalJSON implements encoding.JSONUnmarshaler.
func (i *Int) UnmarshalJSON(input []byte) error {
	return i.UnmarshalText(input)
}

// ToBig converts *Int to *big.Int.
func (i *Int) ToBig() *big.Int {
	return (*big.Int)(i)
}

func (i *Int) setString(s string, base int) (*Int, bool) {
	w, ok := (*big.Int)(i).SetString(s, base)
	return (*Int)(w), ok
}
